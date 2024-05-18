package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/metrics"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
func run() error {
	fmt.Println("Build version:", buildVersion)
	fmt.Println("Build date:", buildDate)
	fmt.Println("Build commit:", buildCommit)

	log.Println("Start agent")
	log.Println("PublicKey=", CryptoKey)
	log.Println("pollInterval=", pollInterval)
	log.Println("reportInterval=", reportInterval)
	log.Println("URLserver=", URLserver)
	log.Println("KeyAuth=", KeyAuth)

	metrics := metrics.NewMetriсs()
	var publicKey *rsa.PublicKey
	if CryptoKey != "" {
		// Декодируем из формата Pem
		block, _ := pem.Decode([]byte(CryptoKey))
		if block == nil {
			err := fmt.Errorf("failed to parse PEM block containing the key: %s", CryptoKey)
			logger.Log.Error("failed decode crypto key", zap.Error(err))
			return err
		}
		pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			logger.Log.Error("failed parse private key", zap.Error(err))
			return err
		}
		publicKey = pubKey
	}

	postreq := server.NewPostRequest(KeyAuth, publicKey)

	n := notifier.NewNotifier(metrics, postreq, time.Duration(time.Duration(pollInterval)*time.Second), time.Duration(time.Duration(reportInterval)*time.Second), URLserver)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	err := n.StartNotifyCron(ctx, RateLimit)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

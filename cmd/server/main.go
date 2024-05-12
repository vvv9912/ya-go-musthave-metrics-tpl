// Package main is the entry point for the programm.
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/postgresql"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"time"
)

// godoc http://localhost:8080/pkg/github.com/vvv9912/ya-go-musthave-metrics-tpl.git/?m=all
// Variables for server settings, set by flag or environment variable.

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

	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	log.Println("Start server")
	log.Println("KeyAuth=", KeyAuth)
	fmt.Println("CryptoKey=", CryptoKey)
	logger.Log.Info("URLserver=" + URLserver)
	logger.Log.Info("timerSend=", zap.Int("timerSend", timerSend))
	logger.Log.Info("FileStoragePath=" + FileStoragePath)
	logger.Log.Info("Restore=", zap.Bool("RESTORE", RESTORE))

	var Repo *store.Repository

	if DatabaseDsn != "" {
		db, err := sql.Open("pgx", DatabaseDsn)
		if err != nil {
			logger.Log.Error("error open db", zap.Error(err))
			return err
		}
		defer db.Close()
		//миграции
		if err := upGauge(context.Background(), db); err != nil {
			logger.Log.Error("error up gauge", zap.Error(err))
			return err
		}
		if err := upCounter(context.Background(), db); err != nil {
			logger.Log.Error("error up counter", zap.Error(err))
			return err
		}
		database := postgresql.NewDatabase(db)
		Repo = store.NewRepository(database)

	} else {
		stor := storage.NewStorage()
		Repo = store.NewRepository(stor)
	}

	//Если включено восстановление данных
	if RESTORE {

		consumer, err := fileutils.NewConsumer(FileStoragePath)
		if err != nil {
			logger.Log.Error("error consumer", zap.Error(err))
			return err
		}

		event, err := consumer.ReadLastEvent(FileStoragePath)
		if err != nil {
			logger.Log.Info("error read last event", zap.Error(err))
		}

		if event != nil {
			for key, val := range event.Counter {
				err = Repo.UpdateCounter(context.Background(), key, val)
				if err != nil {
					logger.Log.Info("error update counter", zap.Error(err))
				}
			}
			for key, val := range event.Gauge {
				err = Repo.UpdateGauge(context.Background(), key, val)
				if err != nil {
					logger.Log.Info("error update gauge", zap.Error(err))
				}
			}
		}

		err = consumer.Close()
		if err != nil {
			logger.Log.Info("error close", zap.Error(err))
		}
	}

	var privateKey *rsa.PrivateKey
	if CryptoKey != "" {
		block, _ := pem.Decode([]byte(CryptoKey))
		if block == nil {
			err := fmt.Errorf("failed to parse PEM block containing the key: %s", CryptoKey)
			logger.Log.Error("failed decode crypto key", zap.Error(err))
			return err
		}
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			logger.Log.Error("failed parse private key", zap.Error(err))
			return err
		}
		privateKey = privKey
	}

	s := server.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	produce, err := fileutils.NewProducer(FileStoragePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer produce.Close()

	err = s.StartServer(ctx, URLserver, Repo, time.Duration(timerSend)*time.Second, produce, KeyAuth, privateKey)
	if err != nil {
		log.Println(err)
		return err
	}

	<-ctx.Done()
	return nil
}

func upGauge(ctx context.Context, db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS GaugeMetrics (    key text unique not null primary key,    val double precision);"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}
func upCounter(ctx context.Context, db *sql.DB) error {
	query := "CREATE TABLE IF NOT EXISTS CounterMetrics (    key text unique not null primary key, val bigint);"
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

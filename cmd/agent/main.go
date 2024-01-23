package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/metrics"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"log"
	"os"
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
	log.Println("Start agent")
	log.Println("pollInterval=", pollInterval)
	log.Println("reportInterval=", reportInterval)
	log.Println("URLserver=", URLserver)
	if err := logger.Initialize("info"); err != nil {
		return err
	}

	metrics := metrics.NewMetriсs()
	postreq := server.NewPostRequest()

	n := notifier.NewNotifier(metrics, postreq, time.Duration(time.Duration(pollInterval)*time.Second), time.Duration(time.Duration(reportInterval)*time.Second), URLserver)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := n.StartNotifyCron(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

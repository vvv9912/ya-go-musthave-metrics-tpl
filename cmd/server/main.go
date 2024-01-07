package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var URLserver string

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
func run() error {
	log.Println("Start server")
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}
	logger.Log.Info("URLserver=" + URLserver)
	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()

	s := server.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err := s.StartServer(ctx, URLserver, gauge, counter)
	if err != nil {
		log.Println(err)
		return err
	}

	<-ctx.Done()
	return nil
}

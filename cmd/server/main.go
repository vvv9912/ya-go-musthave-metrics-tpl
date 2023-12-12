package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
func run() error {
	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()

	s := server.NewServer()
	ctx := context.TODO()
	err := s.StartServer(ctx, URLserver, gauge, counter)
	if err != nil {
		log.Fatal(err)
		return err
	}
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()
	<-done
	return nil
}

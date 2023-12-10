package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
)

func main() {

	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()
	s := server.NewServer()
	ctx := context.TODO()
	s.StartServer(ctx, "localhost:8080", gauge, counter)
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()
	<-done
}

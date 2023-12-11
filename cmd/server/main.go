package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger := log.New(f, "server: ", log.LstdFlags)
	logger.Println("Server start")

	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()

	s := server.NewServer(logger)
	ctx := context.TODO()
	s.StartServer(ctx, "localhost:8080", gauge, counter)
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()
	<-done
}

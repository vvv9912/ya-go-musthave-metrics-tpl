package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store/storage"
	"log"
	"testing"
	"time"
)

//http://localhost:8080/update/gauge/testGauge/123
//http://localhost:8080/update/counter/testGauge/453
//http://localhost:8080/update/counter/testGauge/none
//http://localhost:8080/update/gauge/testGauge/none
//http://localhost:8080/update/dsad/dsae/none
//http://localhost:8080/update/counter/testGauge/123/123

func TestStartServer(t *testing.T) {

	store := storage.NewStorage()

	s := server.NewServer()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	produce, err := fileutils.NewProducer("test.json")
	if err != nil {
		log.Println(err)
	}
	defer produce.Close()

	err = s.StartServer(ctx, "localhost:8080", store, time.Duration(1*time.Second), produce, "", nil)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

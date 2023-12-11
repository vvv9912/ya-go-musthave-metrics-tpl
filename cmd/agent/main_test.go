package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/metrics"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/server"
	"testing"
	"time"
)

func TestStartNotifyCron(t *testing.T) {
	metrics := metrics.NewMetriсs()
	postreq := server.NewPostRequest()
	n := notifier.NewNotifier(metrics, postreq, time.Duration(2*time.Second), time.Duration(10*time.Second), "http://localhost:8080/update/")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := n.StartNotifyCron(ctx)
	if err != nil {
		return
	}
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()
	<-done
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

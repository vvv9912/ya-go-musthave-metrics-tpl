package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/metrics"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/server"
	"time"
)

func main() {

	metrics := metrics.NewMetriсs()
	postreq := server.NewPostRequest()
	n := notifier.NewNotifier(metrics, postreq, time.Duration(2*time.Second), time.Duration(10*time.Second), "http://localhost:8080/update/")
	ctx := context.Background()
	err := n.StartNotifyCron(ctx)
	if err != nil {
		return
	}
	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		done <- struct{}{}
	}()
	<-done
	//select { //почему лучше так?
	//case <-ctx.Done():
	//	// Обработка завершения программы
	//	return
	//}

}

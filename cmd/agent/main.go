package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/metrics"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Agent/server"
	"log"
	"time"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
	//select { //почему лучше так?
	//case <-ctx.Done():
	//	// Обработка завершения программы
	//	return
	//}

}
func run() error {
	log.Println("poll=", pollInterval)
	log.Println("reportInterval=", reportInterval)
	log.Println("serv=", URLserver)
	metrics := metrics.NewMetriсs()
	postreq := server.NewPostRequest()

	n := notifier.NewNotifier(metrics, postreq, time.Duration(time.Duration(pollInterval)*time.Second), time.Duration(time.Duration(reportInterval)*time.Second), URLserver)
	ctx := context.Background()
	err := n.StartNotifyCron(ctx)
	if err != nil {
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

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
	select {
	case <-ctx.Done():
		// Обработка завершения программы

		return
		//default:
	}
	//metric := metrics.UpdateMetrics()
	//wx := &sync.WaitGroup{}
	//for i := 1; i < 10; i++ {
	//	wx.Add(1)
	//	go func(i int) {
	//
	//		time.Sleep(time.Duration(i) * time.Second)
	//
	//		metric := metrics.UpdateMetricsGauge()
	//		fmt.Println(i)
	//		for key, values := range *metric {
	//			fmt.Println("Key:", key, "Values:", values)
	//		}
	//		wx.Done()
	//	}(i)
	//}
	//wx.Wait()
}

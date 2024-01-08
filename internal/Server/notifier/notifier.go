package notifier

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"log"
	"time"
)

type Writer interface {
	WriteEvent(event *fileutils.Event) error
}

type Notifier struct {
	gauge   model.GaugeStorager
	counter model.CounterStorager
	Writer
	TimerSend time.Duration
}

func NewNotifier(gauge model.GaugeStorager, counter model.CounterStorager, timerSend time.Duration, writer Writer) *Notifier {
	return &Notifier{gauge: gauge, counter: counter, TimerSend: timerSend, Writer: writer}
}

func (n *Notifier) StartNotifier(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			default:

			}

			gauge := n.gauge.GetAllGauge()
			counter := n.counter.GetAllCounter()
			//n.producer.WriteEvent(&fileutils.Event{
			//	Gauge:   gauge,
			//	Counter: counter,
			//})
			err := n.WriteEvent(&fileutils.Event{
				Gauge:   gauge,
				Counter: counter,
			})
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(n.TimerSend)
		}

	}()
}

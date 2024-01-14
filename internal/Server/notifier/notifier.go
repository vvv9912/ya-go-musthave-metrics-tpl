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
type NotifierSend interface {
	NotifierPending() error
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

// Отправка при таймере =0
func (n *Notifier) NotifierPending() error {
	if n.TimerSend != 0 {
		return nil
	}
	gauge := n.gauge.GetAllGauge()
	counter := n.counter.GetAllCounter()
	err := n.WriteEvent(&fileutils.Event{
		Gauge:   gauge,
		Counter: counter,
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (n *Notifier) StartNotifier(ctx context.Context) {
	if n.TimerSend == 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(n.TimerSend)
		for {
			select {
			case <-ctx.Done():
				// Обработка завершения программы
				return
			case <-ticker.C:
				gauge := n.gauge.GetAllGauge()
				counter := n.counter.GetAllCounter()

				err := n.WriteEvent(&fileutils.Event{
					Gauge:   gauge,
					Counter: counter,
				})
				if err != nil {
					log.Println(err)
					return
				}
				//time.Sleep(n.TimerSend)
				return
			default:
				return
			}
		}

	}()

}

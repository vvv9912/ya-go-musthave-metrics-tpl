package notifier

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"log"
	"time"
)

type Writer interface {
	WriteEvent(event *fileutils.Event) error
}
type NotifierSend interface {
	NotifierPending(ctx context.Context) error
}
type Notifier struct {
	gauge   service.GaugeStorager
	counter service.CounterStorager
	Writer
	TimerSend time.Duration
}

func NewNotifier(gauge service.GaugeStorager, counter service.CounterStorager, timerSend time.Duration, writer Writer) *Notifier {
	return &Notifier{gauge: gauge, counter: counter, TimerSend: timerSend, Writer: writer}
}

// Отправка при таймере =0
func (n *Notifier) NotifierPending(ctx context.Context) error {
	if n.TimerSend != 0 {
		return nil
	}
	gauge, err := n.gauge.GetAllGauge(ctx)
	if err != nil {
		return err
	}
	counter, err := n.counter.GetAllCounter(ctx)
	if err != nil {
		return err
	}
	err = n.WriteEvent(&fileutils.Event{
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
				gauge, err := n.gauge.GetAllGauge(ctx)
				if err != nil {
					log.Println(err)
					return
				}
				counter, err := n.counter.GetAllCounter(ctx)
				if err != nil {
					log.Println(err)
					return
				}
				err = n.WriteEvent(&fileutils.Event{
					Gauge:   gauge,
					Counter: counter,
				})
				if err != nil {
					log.Println(err)
					return
				}
				continue
			default:
				continue
			}
		}

	}()

}

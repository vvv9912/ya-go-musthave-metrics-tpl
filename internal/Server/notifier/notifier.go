// Пакет notifier для отправки событий по таймеру.
package notifier

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"log"
	"time"
)

// Writer - определяет метод записи события.
type Writer interface {
	WriteEvent(event *fileutils.Event) error
}

// NotifierSend - определяет метод NotifierPending для отправки уведомлений.
type NotifierSend interface {
	NotifierPending(ctx context.Context) error
}

// Notifier - структура "отправки событий".
type Notifier struct {
	store     store.Storager // хранилище с метриками.
	Writer                   // объект, реализующую запись в бд/кэш.
	TimerSend time.Duration  // таймер, для отправки событий с заданным интервалом.
}

// NewNotifier - конструктор.
func NewNotifier(Storage store.Storager, timerSend time.Duration, writer Writer) *Notifier {
	return &Notifier{store: Storage, TimerSend: timerSend, Writer: writer}
}

// NotifierPending - отправляет событие.
// При таймере нулю, отправка не осуществляется.
func (n *Notifier) NotifierPending(ctx context.Context) error {
	if n.TimerSend != 0 {
		return nil
	}

	gauge, err := n.store.GetAllGauge(ctx)
	if err != nil {
		return err
	}

	counter, err := n.store.GetAllCounter(ctx)
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

// StartNotifier - запускает проверку и отправку на основе заданного интервала TimerSend.
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
				gauge, err := n.store.GetAllGauge(ctx)
				if err != nil {
					log.Println(err)
					return
				}

				counter, err := n.store.GetAllCounter(ctx)
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

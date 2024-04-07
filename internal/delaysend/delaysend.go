// Модуль delaysend повторяет функцию с высталвенной задержкой между повторами.
package delaysend

import (
	"errors"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"time"
)

type DelaySend struct {
	delay []int
	error error
}

// Конструктор.
func NewDelaySend() *DelaySend {
	return &DelaySend{}

}

// Выставление задежки.
func (d *DelaySend) SetDelay(t []int) *DelaySend {
	d.delay = t
	return d
}

// Отправка метрики.
func (d *DelaySend) SendDelayedMetrics(f func() (model.Metrics, error)) (model.Metrics, error) {
	m, err := f()
	if err == nil {
		return m, nil //success
	}

	if !errors.Is(err, d.error) {
		return m, err
	}

	for _, v := range d.delay {
		time.Sleep(time.Duration(v) * time.Second)

		m, err = f()
		if err == nil {
			return m, nil //sucess
		}

		if !errors.Is(err, d.error) {
			return m, err
		}

	}

	return m, err
}

// Отправка события с 1 возвращающей ошибкой.
func (d *DelaySend) SendDelayed(f func() error) error {
	err := f()
	if err == nil {
		return nil //success
	}

	if !errors.Is(err, d.error) {
		return err
	}

	for _, v := range d.delay {
		time.Sleep(time.Duration(v) * time.Second)

		err = f()
		if err == nil {
			return nil //sucess
		}

		if !errors.Is(err, d.error) {
			return err
		}

	}

	return err
}

// Добавление ошибки, которую ожидаем получить для повторной отправки события.
func (d *DelaySend) AddExpectedError(e error) *DelaySend {
	d.error = e
	return d
}

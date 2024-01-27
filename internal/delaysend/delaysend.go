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

func NewDelaySend() *DelaySend {
	return &DelaySend{}

}
func (d *DelaySend) SetDelay(t []int) *DelaySend {
	d.delay = t
	return d
}
func (d *DelaySend) SendDelayedMetrics(f func() (model.Metrics, error)) (model.Metrics, error) {
	m, err := f()
	if err == nil {
		return m, nil //success
	}
	if err != nil {
		if !errors.Is(err, d.error) {
			return m, err
		}
	}
	for _, v := range d.delay {
		time.Sleep(time.Duration(v) * time.Second)

		m, err = f()
		if err == nil {
			return m, nil //sucess
		}
		if err != nil {
			if !errors.Is(err, d.error) {
				return m, err
			}
		}
	}

	return m, err
}

// SendDelayed(f func(...interface{}) error)
func (d *DelaySend) SendDelayed(f func() error) error {
	err := f()
	if err == nil {
		return nil //success
	}
	if err != nil {
		if !errors.Is(err, d.error) {
			return err
		}
	}
	for _, v := range d.delay {
		time.Sleep(time.Duration(v) * time.Second)
		//log.Println("попытка отправить метрику #", v)
		err = f()
		if err == nil {
			return nil //sucess
		}
		if err != nil {
			if !errors.Is(err, d.error) {
				return err
			}
		}
	}

	return err
}

func (d *DelaySend) AddExpectedError(e error) *DelaySend {
	d.error = e
	return d
}

package mw

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"net/http"
	"strconv"
	"strings"
)

type Mw struct {
	GaugeStorage   storage.GaugeStorager
	CounterStorage storage.CounterStorager
}

func (m *Mw) Middlware(next http.Handler) http.Handler {
	// получаем handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// здесь пишем логику обработки
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не соответствует Post; Метод: "+r.Method, http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Mw) MiddlwareGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data := strings.Split(req.URL.Path, "/")
		if len(data) != 5 {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusNotFound)
			return
		}
		name := data[3]
		value, err := strconv.ParseFloat(data[4], 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = m.GaugeStorage.UpdateGauge(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Println("name:", name, "value:", value)
		next.ServeHTTP(res, req)
	})
}
func (m *Mw) MiddlwareCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		data := strings.Split(req.URL.Path, "/")
		if len(data) != 5 {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusNotFound)
			return
		}
		name := data[3]

		value, err := strconv.ParseUint(data[4], 10, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = m.CounterStorage.UpdateCounter(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Println("-------\nname:", name, "value:", value)

		next.ServeHTTP(res, req)
	})
}
func (m *Mw) Middlware2Gauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")
		//name := data[3]
		v := chi.URLParam(req, "Value")
		value, err := strconv.ParseFloat(v, 64)

		err = m.GaugeStorage.UpdateGauge(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//fmt.Println("Сработал MW Gauge")
		fmt.Println("name:", name, "value:", value)
		next.ServeHTTP(res, req) //req.WithContext(ctx) //- передать контекст
		//fmt.Println("офф MW Gauge")
	})
}
func (m *Mw) Middlware2Counter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")
		//if name != "PollCount" {
		//	http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
		//	return
		//}

		v := chi.URLParam(req, "Value")

		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = m.CounterStorage.UpdateCounter(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		//fmt.Println("Сработал MW counter")
		//fmt.Println("-------\n", "до парсинга получили:", req.URL.Path, name, ":", value, "name:", name, "value:", value)
		next.ServeHTTP(res, req)
		//fmt.Println("офф MW counter")
	})
}
func (m *Mw) MiddlwareGetGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")
		//name := data[3]
		val, err := m.GaugeStorage.GetGauge(name)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			return
		}
		value := strconv.FormatFloat(val, 'f', -1, 64)
		ctx := context.WithValue(req.Context(), "val", value)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}
func (m *Mw) MiddlwareGetCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")
		//name := data[3]
		val, err := m.CounterStorage.GetCounter(name)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			return
		}
		value := strconv.FormatUint(val, 10)

		ctx := context.WithValue(req.Context(), "val", value)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}

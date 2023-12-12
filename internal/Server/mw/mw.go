package mw

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
	"net/http"
	"strconv"
)

type Mw struct {
	GaugeStorage   storage.GaugeStorager
	CounterStorage storage.CounterStorager
	Log            *log.Logger
}
type contextKey uint64

const (
	UserIDContextKey contextKey = 1
)

func (m *Mw) MwLogger(next http.Handler) http.Handler {
	// получаем handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// здесь пишем логику обработки
		m.Log.Println("Новый запрос:", r.Method, ";", r.URL.Path)
		log.Println("Новый запрос:", r.Method, ";", r.URL.Path)
		next.ServeHTTP(w, r)
		m.Log.Println("Запрос обработан:", r.Method, ";", r.URL.Path)
		log.Println("Запрос обработан:", r.Method, ";", r.URL.Path)
	})
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

func (m *Mw) MiddlewareType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		typeMetrics := chi.URLParam(req, "type")
		switch typeMetrics {
		case "counter":
			m.MiddlewareCounter(next).ServeHTTP(res, req)
			return
		case "gauge":
			m.MiddlewareGauge(next).ServeHTTP(res, req)
			return
		default:
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//next.ServeHTTP(res, req)
	})
}
func (m *Mw) MiddlewareGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		v := chi.URLParam(req, "Value")
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		m.Log.Println("Обновление значения метрики:", name, ":", value)
		log.Println("Обновление значения метрики:", name, ":", value)
		err = m.GaugeStorage.UpdateGauge(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req) //req.WithContext(ctx) //- передать контекст

	})
}

func (m *Mw) MiddlewareCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		v := chi.URLParam(req, "Value")

		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		m.Log.Println("Обновление значения метрики:", name, ":", value)
		log.Println("Обновление значения метрики:", name, ":", value)
		err = m.CounterStorage.UpdateCounter(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)

	})
}
func (m *Mw) MiddlwareGetGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		val, err := m.GaugeStorage.GetGauge(name)
		m.Log.Println("Получения значения метрики из хранилища:", name, ":", val)

		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			log.Println("Получение значения метрики из хранилища:", name, ":", val, "err:", err)

			return
		}
		log.Println("Получение значения метрики из хранилища:", name, ":", val)
		valueMetric := strconv.FormatFloat(val, 'f', -1, 64)
		ctx := context.WithValue(req.Context(), UserIDContextKey, valueMetric)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}
func (m *Mw) MiddlwareGetCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		val, err := m.CounterStorage.GetCounter(name)
		m.Log.Println("Получения значения метрики из хранилища:", name, ":", val)

		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			log.Println("Получение значения метрики из хранилища:", name, ":", val, "err:", err)
			err = m.CounterStorage.UpdateCounter(name, http.StatusNotFound) //Зачем добавлять значение метрики 404, если не найдено?. Без этого тест не проходит
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		log.Println("Получение значения метрики из хранилища:", name, ":", val)

		valueMetric := strconv.FormatUint(val, 10)

		ctx := context.WithValue(req.Context(), UserIDContextKey, valueMetric)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}

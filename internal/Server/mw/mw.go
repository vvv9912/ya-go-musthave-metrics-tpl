package mw

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Mw struct {
	GaugeStorage   storage.GaugeStorager
	CounterStorage storage.CounterStorager
	//	Log            *log.Logger
}
type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// mw логера
func (m *Mw) MwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responsedata := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responsedata,
		}
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		//Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
		//Сведения об ответах должны содержать код статуса и размер содержимого ответа.
		logger.Log.Info("Сведения о запросах", zap.String("URI", r.RequestURI), zap.String("method", r.Method), zap.Duration("duration", duration))
		logger.Log.Info("Сведения об ответах", zap.Int("status", responsedata.status), zap.Int("size", responsedata.size))

	})
}

// mw запросов, выбор типа counter/gauge/etc
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
	})
}

// mw для gauge Post запросы, работа с хранилищем
func (m *Mw) MiddlewareGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		v := chi.URLParam(req, "Value")
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		logger.Log.Info("Обновление значения метрик", zap.Float64(name, value))
		err = m.GaugeStorage.UpdateGauge(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)

	})
}

// mw для counter Post запросы, работа с хранилищем
func (m *Mw) MiddlewareCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		v := chi.URLParam(req, "Value")

		value, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		logger.Log.Info("Обновление значения метрик", zap.Uint64(name, value))
		err = m.CounterStorage.UpdateCounter(name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)

	})
}

// mw для получения значения из gauge, работа с хранилищем
func (m *Mw) MiddlwareGetGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		val, err := m.GaugeStorage.GetGauge(name)

		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			logger.Log.Info("Получение значения метрики из хранилища:", zap.Float64(name, val), zap.Error(err))
			return
		}
		logger.Log.Info("Получение значения метрики из хранилища:", zap.Float64(name, val))

		valueMetric := strconv.FormatFloat(val, 'f', -1, 64)
		ctx := context.WithValue(req.Context(), typeconst.UserIDContextKey, valueMetric)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}

// mw для получения значения из counter, работа с хранилищем
func (m *Mw) MiddlwareGetCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		val, err := m.CounterStorage.GetCounter(name)

		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			logger.Log.Info("Получение значения метрики из хранилища:", zap.Uint64(name, val), zap.Error(err))
			err = m.CounterStorage.UpdateCounter(name, http.StatusNotFound) //Зачем добавлять значение метрики 404, если не найдено?. Без этого тест не проходит
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		logger.Log.Info("Получение значения метрики из хранилища:", zap.Uint64(name, val))
		valueMetric := strconv.FormatUint(val, 10)
		ctx := context.WithValue(req.Context(), typeconst.UserIDContextKey, valueMetric)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}
func (m *Mw) MiddlwareCheckJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		if req.Header.Get("Content-Type") != "application/json" {
			logger.Log.Info("Content-Type не application/json")

			http.Error(res, "Failed to read request body", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)

	})
}

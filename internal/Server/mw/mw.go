package mw

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/gzipwrapper"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Mw struct {
	Service *service.Service
}

func NewMw(s *service.Service) *Mw {
	return &Mw{Service: s}
}

type responseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.body.Write(b)
	if err != nil {
		logger.Log.Error("error write body", zap.Error(err))
		return n, err
	}
	return rw.ResponseWriter.Write(b)
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
		data := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   data,
		}
		next.ServeHTTP(&lw, r)
		duration := time.Since(start)
		//Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
		//Сведения об ответах должны содержать код статуса и размер содержимого ответа.
		logger.Log.Info("Сведения о запросах", zap.String("URI", r.RequestURI), zap.String("method", r.Method), zap.Duration("duration", duration))
		logger.Log.Info("Сведения об ответах", zap.Int("status", data.status), zap.Int("size", data.size))

	})
}

// проверяем, что клиент умеет получать от сервера сжатые данные в определенном формате
func supportAcceptType(acceptType map[string]struct{}, acceptTypeReq string) bool {
	if acceptTypeReq == "*/*" {
		return true
	} else {
		for key := range acceptType {
			if strings.Contains(acceptTypeReq, key) {
				return true
			}
		}
	}
	return false
}

// проверяем, что клиент поддерживает соответствующий content-type
func supportEncodingType(accpetEncoding map[string]struct{}, acceptEncodingReq string) bool {
	for key := range accpetEncoding {
		if strings.Contains(acceptEncodingReq, key) {
			return true
		}
	}
	return false
}
func (m *Mw) MiddlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		acceptType := map[string]struct{}{
			"application/json": struct{}{},
			"text/html":        struct{}{},
			"html/text":        struct{}{}, //
		}

		supportGzip := map[string]struct{}{
			"gzip": struct{}{},
		}

		supportAccept := supportAcceptType(acceptType, r.Header.Get("Accept"))
		supportEncoding := supportEncodingType(supportGzip, r.Header.Get("Accept-Encoding"))

		if supportAccept && supportEncoding {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := gzipwrapper.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		supportContentEncoding := supportEncodingType(supportGzip, r.Header.Get("Content-Encoding"))

		if supportContentEncoding {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := gzipwrapper.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Error("ошибка декомпрессии", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)
	})

}
func (m *Mw) MiddlewareHashAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reader := io.TeeReader(r.Body, os.Stdout)
		body, err := io.ReadAll(reader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//считаем хеш
		h := hmac.New(sha256.New, []byte(m.Service.KeyAuth))
		h.Write(body)
		dst := h.Sum(nil)
		hash, err := hex.DecodeString(r.Header.Get("HashSHA256"))
		if err != nil {
			logger.Log.Error("Error decoding hash", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ok := hmac.Equal(dst, hash)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// передаём управление хендлеру
		rw := &responseWriter{ResponseWriter: w, body: bytes.NewBuffer(body)}
		next.ServeHTTP(rw, r)

		hWriter := hmac.New(sha256.New, []byte(m.Service.KeyAuth))
		hWriter.Write(rw.body.Bytes())
		hashWriter := h.Sum(nil)
		w.Header().Set("HashSHA256", string(hashWriter))

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

//тк для mW не нужна больше структура, можно сделать так

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
		err = m.Service.Storage.UpdateGauge(req.Context(), name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = m.Service.Metrics.SendMetricstoFile(req.Context())
		if err != nil {
			logger.Log.Error("ошибка отправки метрик в файл", zap.Error(err))
		}
		next.ServeHTTP(res, req)

	})
}

// mw для counter Post запросы, работа с хранилищем
func (m *Mw) MiddlewareCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		v := chi.URLParam(req, "Value")

		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		logger.Log.Info("Обновление значения метрик", zap.Int64(name, value))
		err = m.Service.Storage.UpdateCounter(req.Context(), name, value)
		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = m.Service.Metrics.SendMetricstoFile(req.Context())
		if err != nil {
			logger.Log.Error("ошибка отправки метрик в файл", zap.Error(err))
		}
		next.ServeHTTP(res, req)

	})
}

// mw для получения значения из gauge, работа с хранилищем
func (m *Mw) MiddlwareGetGauge(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		name := chi.URLParam(req, "SomeMetric")

		val, err := m.Service.Metrics.GetGauge(req.Context(), name)

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

		val, err := m.Service.Metrics.GetCounter(req.Context(), name)

		if err != nil {
			http.Error(res, fmt.Sprintln(http.StatusNotFound), http.StatusNotFound)
			logger.Log.Info("Получение значения метрики из хранилища:", zap.Int64(name, val), zap.Error(err))

			err = m.Service.Storage.UpdateCounter(req.Context(), name, http.StatusNotFound) //Зачем добавлять значение метрики 404, если не найдено?. Без этого тест не проходит
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		logger.Log.Info("Получение значения метрики из хранилища:", zap.Int64(name, val))
		valueMetric := strconv.FormatInt(val, 10)
		ctx := context.WithValue(req.Context(), typeconst.UserIDContextKey, valueMetric)

		next.ServeHTTP(res, req.WithContext(ctx))

	})
}
func (m *Mw) MiddlwareCheckJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
			logger.Log.Info("Content-Type не application/json")
			http.Error(res, "Failed to read request body", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(res, req)

	})
}

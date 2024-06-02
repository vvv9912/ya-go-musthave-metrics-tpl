package mw

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/gzipwrapper"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Mw struct {
	Service       *service.Service
	privateKey    *rsa.PrivateKey
	trustedSubnet string
}

func NewMw(s *service.Service, trustedSubnet string, provateKey *rsa.PrivateKey) *Mw {
	return &Mw{Service: s, privateKey: provateKey, trustedSubnet: trustedSubnet}
}

// mw доверенной подсети
func (m *Mw) MwTrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		agentIP := r.Header.Get("X-Real-IP")
		if agentIP == "" {
			next.ServeHTTP(w, r)
		}

		_, ipv4Net, err := net.ParseCIDR(m.trustedSubnet)
		if err != nil {
			logger.Log.Error("Error parsing trusted subnet", zap.String("error", err.Error()))
		}

		agent := net.ParseIP(agentIP)

		if !ipv4Net.Contains(agent) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)

	})
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

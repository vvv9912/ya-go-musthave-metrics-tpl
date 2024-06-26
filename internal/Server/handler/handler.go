// Package handler - обработчик HTTP-запросов.
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/delaysend"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"io"
	"net/http"
	"syscall"
)

// Handler - Структура с сервисным слоем.
type Handler struct {
	Service *service.Service
}

// NewHandler - Конструктор.
func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

// HandlerSucess - обработчик для успешного HTTP-запроса.
func HandlerSucess(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	_, err := res.Write([]byte(body))
	if err != nil {
		logger.Log.Info("Failed to get gauge", zap.Error(err))
		return
	}
}

// HandlerGetCounter - обработчик для получения значения метрики счетчика.
func HandlerGetCounter(res http.ResponseWriter, req *http.Request) {
	valCtx := req.Context().Value(typeconst.UserIDContextKey)
	value := valCtx.(string)
	name := chi.URLParam(req, "SomeMetric")

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	_, err := res.Write([]byte(body))
	if err != nil {
		logger.Log.Info("Failed to get gauge", zap.Error(err))
		return
	}
}

// HandlerGetGauge - обработчик для получения значения метрики Gauge.
func HandlerGetGauge(res http.ResponseWriter, req *http.Request) {
	valCtx := req.Context().Value(typeconst.UserIDContextKey)
	value := valCtx.(string)
	name := chi.URLParam(req, "SomeMetric")

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	_, err := res.Write([]byte(body))
	if err != nil {
		logger.Log.Info("Failed to get gauge", zap.Error(err))
		return
	}
}

// HandlerGetMetrics - функция, которая возвращает обработчик для получения всех метрик.
// Была создана для проверки конструкции, когда нужно использовать дополнительные перменные.
func HandlerGetMetrics(storage store.Storager) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		gauge, err := storage.GetAllGauge(req.Context())
		if err != nil {
			logger.Log.Info("Failed to get gauge", zap.Error(err))
			http.Error(res, "Failed to get gauge", http.StatusInternalServerError)
			return
		}
		body := ""
		for key, value := range gauge {
			body += fmt.Sprintf("%s: %f\n", key, value)
		}
		count, err := storage.GetAllCounter(req.Context())
		if err != nil {
			logger.Log.Info("Failed to get counter", zap.Error(err))
			http.Error(res, "Failed to get counter", http.StatusInternalServerError)
			return
		}

		for key, value := range count {
			body += fmt.Sprintf("%s: %v\n", key, value)
		}
		res.WriteHeader(http.StatusOK)

		_, err = res.Write([]byte(body))
		if err != nil {
			logger.Log.Info("Failed to get gauge", zap.Error(err))
			return
		}
	}
}

// HandlerPostJSON - обработчик JSON запросов с метриками.
func (h *Handler) HandlerPostJSON(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		http.Error(res, "Failed to read request body", http.StatusNotFound)
		return
	}

	var metrics model.Metrics

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		logger.Log.Info("Failed to unmarshal request body to JSON", zap.Error(err))
		http.Error(res, "Failed to unmarshal request body", http.StatusInternalServerError)
		return
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return h.Service.Metrics.PutMetrics(req.Context(), metrics)
		})
	if err != nil {
		logger.Log.Info("Failed to put metrics", zap.Error(err))
		http.Error(res, "Failed to put metrics", http.StatusNotFound)
		return
	}

	err = h.Service.Metrics.SendMetricstoFile(req.Context())
	if err != nil {
		logger.Log.Error("Failed to send metrics to file", zap.Error(err))
		http.Error(res, "Failed to send metrics to file", http.StatusInternalServerError)
		return
	}

	response := body

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(response)
	if err != nil {
		logger.Log.Error("Failed to HandlerPostJSON", zap.Error(err))
		return
	}
}

// HandlerGetJSON - возврат метрики в формате JSON.
func (h *Handler) HandlerGetJSON(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		http.Error(res, "Failed to read request body", http.StatusNotFound)
		return
	}
	var metrics model.Metrics
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		logger.Log.Info("Failed to retrieve metrics", zap.Error(err))
		http.Error(res, "Failed to retrieve metrics", http.StatusInternalServerError)
		return
	}

	metrics, err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).SendDelayedMetrics(func() (model.Metrics, error) {
		return h.Service.Metrics.GetMetrics(req.Context(), metrics)
	})
	if err != nil {
		logger.Log.Info("Failed to get metrics", zap.Error(err))
		http.Error(res, "Failed to get metrics", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Info("Failed to unmarshal metrics", zap.Error(err))
		http.Error(res, "Failed to unmarshal metrics", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write(response)
	if err != nil {
		logger.Log.Error("Failed to write", zap.Error(err))
		return
	}
}

// HandlerGauge - возвращает метрики Gauge.
func (h *Handler) HandlerGauge(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		http.Error(res, "Failed to read request body", http.StatusNotFound)
		return
	}
	var metrics model.Metrics
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		logger.Log.Info("Failed to retrieve metrics", zap.Error(err))
		http.Error(res, "Failed to retrieve metrics", http.StatusInternalServerError)
		return
	}

	metrics, err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).SendDelayedMetrics(func() (model.Metrics, error) {
		return h.Service.Metrics.GetMetrics(req.Context(), metrics)
	})
	if err != nil {
		logger.Log.Info("Failed to get metrics", zap.Error(err))
		http.Error(res, "Failed to get Metrics", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Info("Failed to marshal metrics", zap.Error(err))
		http.Error(res, "Failed to marshal metrics", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write(response)
	if err != nil {
		logger.Log.Error("Failed to write", zap.Error(err))
		return
	}
}

// HandlerPingDatabase - Ping БД.
func (h *Handler) HandlerPingDatabase(res http.ResponseWriter, req *http.Request) {

	err := h.Service.Storage.Ping(req.Context())
	if err != nil {
		logger.Log.Info("Failed to ping database", zap.Error(err))
		http.Error(res, "Failed to ping database", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)

}

// HandlerPostBatched - принимает значения метрик батчами.
func (h *Handler) HandlerPostBatched(res http.ResponseWriter, req *http.Request) {
	var metrics []model.Metrics

	err := json.NewDecoder(req.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		http.Error(res, "Failed to read request body", http.StatusNotFound)
		return
	}

	err = delaysend.NewDelaySend().SetDelay([]int{1, 3, 5}).
		AddExpectedError(syscall.ECONNREFUSED).
		SendDelayed(func() error {
			return h.Service.Storage.UpdateMetricsBatch(req.Context(), metrics)
		})

	if err != nil {
		logger.Log.Info("Failed to send metrics to batch", zap.Error(err))
		http.Error(res, "Failed to send metrics to batch", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

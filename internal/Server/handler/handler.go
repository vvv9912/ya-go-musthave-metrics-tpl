package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/store"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

func HandlerSucess(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
}

func HandlerGetCounter(res http.ResponseWriter, req *http.Request) {
	valCtx := req.Context().Value(typeconst.UserIDContextKey)
	value := valCtx.(string)
	name := chi.URLParam(req, "SomeMetric")

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	res.Write([]byte(body))
}
func HandlerGetGauge(res http.ResponseWriter, req *http.Request) {
	valCtx := req.Context().Value(typeconst.UserIDContextKey)
	value := valCtx.(string)
	name := chi.URLParam(req, "SomeMetric")

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	res.Write([]byte(body))
}

func HandlerGetMetrics(gauger storage.GaugeStorager, counter storage.CounterStorager) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html")
		gauge, err := gauger.GetAllGauge(req.Context())
		if err != nil {
			logger.Log.Info("Failed to get gauge", zap.Error(err))
			http.Error(res, "Failed to get gauge", http.StatusInternalServerError)
			return
		}
		body := ""
		for key, value := range gauge {
			body += fmt.Sprintf("%s: %f\n", key, value)
		}
		count, err := counter.GetAllCounter(req.Context())
		if err != nil {
			logger.Log.Info("Failed to get counter", zap.Error(err))
			http.Error(res, "Failed to get counter", http.StatusInternalServerError)
			return
		}
		for key, value := range count {
			body += fmt.Sprintf("%s: %v\n", key, value)
		}
		res.WriteHeader(http.StatusOK)

		res.Write([]byte(body))
	}
}
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

	err = h.Service.Metrics.PutMetrics(req.Context(), metrics)
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
	res.Write(response)
}

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

	metrics, err = h.Service.Metrics.GetMetrics(req.Context(), metrics)
	if err != nil {
		logger.Log.Info("Failed to get metrics", zap.Error(err))
		http.Error(res, "Failed to get metrics", http.StatusNotFound)
		return
	}
	log.Println(metrics)
	response, err := json.Marshal(metrics)
	if err != nil {
		logger.Log.Info("Failed to unmarshal metrics", zap.Error(err))
		http.Error(res, "Failed to unmarshal metrics", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(response)
}
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

	metrics, err = h.Service.Metrics.GetMetrics(req.Context(), metrics)
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
	res.Write(response)
}
func (h *Handler) HandlerPingDatabase(res http.ResponseWriter, req *http.Request) {
	store2 := h.Service.Store
	if (store2) == (*store.Database)(nil) {
		logger.Log.Info("Failed to ping database")
		http.Error(res, "Failed to ping database", http.StatusInternalServerError)
		return
	}

	err := (store2).Ping(req.Context())
	if err != nil {
		logger.Log.Info("Failed to ping database", zap.Error(err))
		http.Error(res, "Failed to ping database", http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)

}

func (h *Handler) HandlerPostBatched(res http.ResponseWriter, req *http.Request) {
	var metrics []model.Metrics

	err := json.NewDecoder(req.Body).Decode(&metrics)
	if err != nil {
		logger.Log.Info("Failed to read request body", zap.Error(err))
		http.Error(res, "Failed to read request body", http.StatusNotFound)
		return
	}
	log.Println(metrics)
	err = h.Service.Store.UpdateMetricsBatch(req.Context(), metrics)
	log.Println("okay, записалось в бд")
	if err != nil {
		logger.Log.Info("Failed to send metrics to file", zap.Error(err))
		http.Error(res, "Failed to send metrics to file", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

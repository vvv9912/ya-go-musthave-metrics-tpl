package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/project"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/model"
	"io"
	"net/http"
)

type Handler struct {
	P project.Project
}

func NewHandler(p project.Project) *Handler {
	return &Handler{P: p}
}

func HandlerSucess(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
}
func HandlerErrType(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	body := fmt.Sprintf("%v", http.StatusBadRequest)
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
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		gauge := gauger.GetAllGauge()
		body := ""
		for key, value := range gauge {
			body += fmt.Sprintf("%s: %f\n", key, value)
		}
		count := counter.GetAllCounter()
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
		// Обработка ошибки чтения тела запроса
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var metrics model.Metrics
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}
	err = h.P.PutMetrics(metrics)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(body))
}

func (h *Handler) HandlerGetJSON(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		// Обработка ошибки чтения тела запроса
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var metrics model.Metrics
	err = json.Unmarshal(body, &metrics)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}
	metrics, err = h.P.GetMetrics(metrics)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}

	bodyWrite, err := json.Marshal(metrics)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(bodyWrite)
}

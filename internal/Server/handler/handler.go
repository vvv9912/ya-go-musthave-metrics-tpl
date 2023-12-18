package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/typeconst"
	"net/http"
)

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

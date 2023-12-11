package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"net/http"
)

func HandlerSucess(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
}
func HandlerErrType(res http.ResponseWriter, req *http.Request) {

	//res.Header().Set("text/plain", "charset=utf-8")
	res.WriteHeader(http.StatusBadRequest)
	body := fmt.Sprintf("%v", http.StatusBadRequest)
	res.Write([]byte(body))
}
func HandlerGauge(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
	//fmt.Println("Сработал Handler Gauge")
}
func HandlerCounter(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("-------\n", "POST:", req.URL.Path)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
	//fmt.Println("Сработал Handler counter")
}

func HandlerGetCounter(res http.ResponseWriter, req *http.Request) {
	value := req.Context().Value("value_metric").(string)
	name := chi.URLParam(req, "SomeMetric")

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	//fmt.Println("GetRequest=", value)
	//fmt.Println(res.Header())
	res.Write([]byte(body))
	//fmt.Println("Сработал Handler counter")
}
func HandlerGetGauge(res http.ResponseWriter, req *http.Request) {
	value := req.Context().Value("value_metric").(string)
	name := chi.URLParam(req, "SomeMetric")
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Header().Set(name, value)
	res.WriteHeader(http.StatusOK)
	body := value
	//fmt.Println("GetRequest=", value)
	//fmt.Println(res.Header())
	//body := fmt.Sprintf("%v\n%s:%s", http.StatusOK, name, value)
	res.Write([]byte(body))
	//fmt.Println("Сработал Handler counter")
}

func HandlerGetDef(res http.ResponseWriter, req *http.Request, gauger storage.GaugeStorager, counter storage.CounterStorager) {

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
	//fmt.Println("Сработал Handler counter")
}

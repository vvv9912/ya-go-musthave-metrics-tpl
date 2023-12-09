package mw

import (
	"fmt"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"net/http"
	"strconv"
	"strings"
)

type Mw struct {
	GaugeStorage   storage.GaugeStorager
	CounterStorage storage.CounterStorager
}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>,
// Content-Type: text/plain
func (m *Mw) Middlware(next http.Handler) http.Handler {
	// получаем handler приведением типа http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// здесь пишем логику обработки
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не соответствует Post; Метод: "+r.Method, http.StatusBadRequest)
			return
		}
		//if r.Header.Get("Content-Type") != "text/plain" {
		//	http.Error(w, "Content-Type не соответствует text/plain; Метод: "+r.Header.Get("Content-Type"), http.StatusBadRequest)
		//	return
		//}
		//data := strings.Split(r.URL.Path, "/")
		//if len(data) != 5 {
		//	http.Error(w, "Err", http.StatusBadRequest)
		//	return
		//}

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
		value, err := strconv.ParseInt(data[4], 10, 64)
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

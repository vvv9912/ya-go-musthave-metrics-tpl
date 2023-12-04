package handler

import (
	"fmt"
	"net/http"
)

func HandlerGauge(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("text/plain", "charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
}
func HandlerCounter(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("text/plain", "charset=utf-8")
	res.WriteHeader(http.StatusOK)
	body := fmt.Sprintf("%v", http.StatusOK)
	res.Write([]byte(body))
}
func HandlerBase(res http.ResponseWriter, req *http.Request) {

	//res.Header().Set("text/plain", "charset=utf-8")
	res.WriteHeader(http.StatusBadRequest)
	body := fmt.Sprintf("%v", http.StatusBadRequest)
	res.Write([]byte(body))
}

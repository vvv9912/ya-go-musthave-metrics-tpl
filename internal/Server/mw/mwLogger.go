package mw

import (
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

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
func supportEncodingTypeOld(accpetEncoding map[string]struct{}, acceptEncodingReq string) bool {
	for key := range accpetEncoding {
		if strings.Contains(acceptEncodingReq, key) {
			return true
		}
	}
	return false
}

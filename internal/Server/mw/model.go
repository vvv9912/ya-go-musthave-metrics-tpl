package mw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
)

// Logger
type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// для хэша
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	keyAuth    []byte
	hashWriter []byte
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

func (rw *responseWriter) Write(b []byte) (int, error) {
	//Считаем хэш
	hWriter := hmac.New(sha256.New, rw.keyAuth)

	_, err := hWriter.Write(b)
	if err != nil {
		return 0, err
	}

	rw.hashWriter = hWriter.Sum(nil)

	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) GetHash() string {
	return string(rw.hashWriter)
}

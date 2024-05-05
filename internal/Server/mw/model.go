package mw

import (
	"crypto/rand"
	"crypto/rsa"
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

// для шифрования
type responseWriterEncrypt struct {
	http.ResponseWriter
	pk *rsa.PublicKey
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

func (rw *responseWriterEncrypt) Write(b []byte) (int, error) {
	//Шифруем
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rw.pk, b)
	if err != nil {
		return 0, err
	}

	return rw.ResponseWriter.Write(ciphertext)
}

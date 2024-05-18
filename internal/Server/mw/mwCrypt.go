package mw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// responseWriterEncrypt - структура для де/шифрования сообщений.
type responseWriterEncrypt struct {
	http.ResponseWriter
	pk *rsa.PublicKey
}

// Write - шифрует сообщение.
func (rw *responseWriterEncrypt) Write(b []byte) (int, error) {
	//Шифруем
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rw.pk, b)
	if err != nil {
		return 0, err
	}

	return rw.ResponseWriter.Write(ciphertext)
}

// MiddlewareCrypt - middleware для де/шифрования сообщения.
func (m *Mw) MiddlewareCrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Считываем из body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Error("Error read body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		decryptMsg, err := rsa.DecryptPKCS1v15(rand.Reader, m.privateKey, body)
		if err != nil {
			logger.Log.Error("Error decrypt", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Копируем дешифрованное сообщение в body
		r.Body = io.NopCloser(bytes.NewBuffer(decryptMsg))

		// не нужно, тк обратное сообщение не шифруем
		//newWriter := &responseWriterEncrypt{ResponseWriter: w, pk: m.publicKey}

		next.ServeHTTP(w, r)

	})

}

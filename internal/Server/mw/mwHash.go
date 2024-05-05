package mw

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

func (m *Mw) MiddlewareHashAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		hashSha := r.Header.Get("HashSHA256")

		if hashSha != "" {
			hash, err := hex.DecodeString(hashSha)
			if err != nil {
				logger.Log.Error("Error decoding hash", zap.Error(err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// копируем тело запроса для дальнейшей проверки хеша
			reader := io.TeeReader(r.Body, os.Stdout) //

			body, err := io.ReadAll(reader)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			//считаем хеш
			h := hmac.New(sha256.New, []byte(m.Service.KeyAuth))

			h.Write(body)
			dst := h.Sum(nil)

			ok := hmac.Equal(dst, hash)
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		}
		// подменяем метод
		rw := &responseWriter{ResponseWriter: w, body: bytes.NewBuffer(nil)}
		next.ServeHTTP(rw, r)

		w.Header().Set("HashSHA256", rw.GetHash())

	})

}

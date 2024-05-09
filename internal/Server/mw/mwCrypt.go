package mw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"
)

func (m *Mw) MiddlewareCrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		////todo Логика когда будем шифровать возможно ее не будет:)
		//supportEncrypt := true
		//
		//if !supportEncrypt {
		//	next.ServeHTTP(w, r)
		//	return
		//}

		//Считываем из body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			//todo logger
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Рассшифровывем тело
		decryptMsg, err := rsa.DecryptPKCS1v15(rand.Reader, m.privateKey, body)
		if err != nil {
			//todo logger
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Копируем дешифрованное сообщение в body
		r.Body = io.NopCloser(bytes.NewBuffer(decryptMsg))

		newWriter := &responseWriterEncrypt{ResponseWriter: w, pk: m.publicKey}

		next.ServeHTTP(newWriter, r)

	})

}

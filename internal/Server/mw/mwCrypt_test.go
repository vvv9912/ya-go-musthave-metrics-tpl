package mw

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMw_MiddlewareCrypt(t *testing.T) {

	// Create public/private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generate key:", err)
		return
	}

	// Create public key
	publicKey := &privateKey.PublicKey

	// Create a sample request body
	originalBody := []byte("This is a secret message")

	encryptedBody, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, originalBody)
	if err != nil {
		t.Fatalf("Failed to encrypt the request body: %v", err)
	}

	// Create a new HTTP request with the encrypted body
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(encryptedBody))

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create an instance of the Mw struct
	mw := &Mw{
		privateKey: privateKey,
	}

	// Create handler test
	handlerTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			t.Fatalf("Failed to read the request body: %v", err)
			return
		}
		w.Write(data)
	})
	// Call the MiddlewareCrypt function with the dummy handler
	handler := mw.MiddlewareCrypt(handlerTest)

	// Serve the request and record the response
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}
	assert.Equal(t, string(originalBody), rr.Body.String())

}

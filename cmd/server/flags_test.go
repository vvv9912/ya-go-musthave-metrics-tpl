package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

//// todo: Тест с флагами и конфигом. Как переедать флаги?
//func Test_parseFlags(t *testing.T) {
//	//Флаги
//	//flags := map[string]string{
//	//	"k":          "KEY",
//	//	"a":          "ADDRESS",
//	//	"l":          "LOG_LEVEL",
//	//	"f":          "FILE_STORAGE_PATH",
//	//	"i":          "STORE_INTERVAL",
//	//	"r":          "RESTORE",
//	//	"d":          "DATABASE_DSN",
//	//	"crypto-key": "CRYPTO_KEY",
//	//	"c":          "CONFIG",
//	//}
//
//	os.Setenv("KEY", "test_key")
//	os.Setenv("ADDRESS", "test_address")
//	os.Setenv("LOG_LEVEL", "debug")
//	os.Setenv("FILE_STORAGE_PATH", "/tmp/test-db.json")
//	os.Setenv("STORE_INTERVAL", "600")
//	os.Setenv("RESTORE", "false")
//	os.Setenv("DATABASE_DSN", "test_dsn")
//	os.Setenv("CRYPTO_KEY", "test_crypto_key")
//	os.Setenv("CONFIG", "")
//
//	defer func() {
//		os.Unsetenv("KEY")
//		os.Unsetenv("ADDRESS")
//		os.Unsetenv("LOG_LEVEL")
//		os.Unsetenv("FILE_STORAGE_PATH")
//		os.Unsetenv("STORE_INTERVAL")
//		os.Unsetenv("RESTORE")
//		os.Unsetenv("DATABASE_DSN")
//		os.Unsetenv("CRYPTO_KEY")
//		os.Unsetenv("CONFIG")
//	}()
//
//	parseFlags()
//
//	assert.Equal(t, KeyAuth, os.Getenv("KEY"), "failed test_key")
//	assert.Equal(t, URLserver, os.Getenv("ADDRESS"), "failed test_address")
//	assert.Equal(t, FileStoragePath, os.Getenv("FILE_STORAGE_PATH"), "failed test_file")
//	assert.Equal(t, strconv.Itoa(timerSend), os.Getenv("STORE_INTERVAL"), "failed test_store")
//	assert.Equal(t, strconv.FormatBool(RESTORE), os.Getenv("RESTORE"), "failed test_rest")
//	assert.Equal(t, DatabaseDsn, os.Getenv("DATABASE_DSN"), "failed test_database")
//	assert.Equal(t, CryptoKey, os.Getenv("CRYPTO_KEY"), "failed test_crypto_key")
//	assert.Equal(t, Config, os.Getenv("CONFIG"), "failed test_config")
//
//}

// Тест конфига
func Test_parseFlags2(t *testing.T) {

	jsonData := map[string]interface{}{
		"address":        "localhostJson:8099",
		"restore":        true,
		"store_interval": "1s",
		"store_file":     "/path/to/json.db",
		"database_dsn":   "database_json",
		"crypto_key":     "/path/to/json.pem",
	}
	f, err := os.OpenFile("test.json", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ") // Устанавливаем отступы для форматирования JSON

	err = encoder.Encode(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	fileInfo, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("KEY", "test_key")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/test-db.json")
	os.Setenv("RESTORE", "false")
	os.Setenv("CRYPTO_KEY", "test_crypto_key")
	os.Setenv("CONFIG", fileInfo.Name())

	defer func() {
		os.Unsetenv("KEY")
		os.Unsetenv("ADDRESS")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("RESTORE")
		os.Unsetenv("DATABASE_DSN")
		os.Unsetenv("CRYPTO_KEY")
		os.Unsetenv("CONFIG")
	}()

	parseFlags()

	assert.Equal(t, KeyAuth, os.Getenv("KEY"), "failed KeyAuth")
	assert.Equal(t, URLserver, "localhostJson:8099", "failed URLserver")
	assert.Equal(t, FileStoragePath, os.Getenv("FILE_STORAGE_PATH"), "failed FileStoragePath")
	assert.Equal(t, timerSend, 1, "failed timerSend")
	assert.Equal(t, strconv.FormatBool(RESTORE), os.Getenv("RESTORE"), "failed RESTORE")
	assert.Equal(t, DatabaseDsn, "database_json", "failed DatabaseDsn")
	assert.Equal(t, CryptoKey, os.Getenv("CRYPTO_KEY"), "failed CryptoKey")
	assert.Equal(t, Config, os.Getenv("CONFIG"), "failed Config")

}

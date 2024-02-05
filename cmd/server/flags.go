package main

import (
	"errors"
	"flag"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

var (
	flagLogLevel string
	DatabaseDsn  string
)

func (o *NetAddress) String() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}

// Set связывает переменную типа со значением флага
// и устанавливает правила парсинга для пользовательского типа.
func (o *NetAddress) Set(flagValue string) error {
	s := strings.Split(flagValue, ":")
	if len(s) != 2 {
		return errors.New("неправильный формат")
	}
	o.Host = s[0]
	var err error
	o.Port, err = strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	return nil
}

func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	addr := new(NetAddress)
	addr.Host = "localhost"
	addr.Port = 8080
	// если интерфейс не реализован,
	// здесь будет ошибка компиляции
	var _ = flag.Value(addr)
	//var restore string
	flag.StringVar(&KeyAuth, "k", "", "key for auth (по умолчанию пустая)")
	flag.Var(addr, "a", "Net address host:port")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.StringVar(&FileStoragePath, "f", "/tmp/metrics-db.json", "file storage path")
	flag.IntVar(&timerSend, "i", 300, "send timer")
	flag.BoolVar(&RESTORE, "r", true, "restore")
	flag.StringVar(&DatabaseDsn, "d", "", "DATABASE_DSN") //	//postgres://postgres:postgres@localhost:5432/postgres"

	flag.Parse()

	URLserver = addr.String()
	if envKey := os.Getenv("KEY"); envKey != "" {
		log.Println(envKey)
		KeyAuth = envKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		URLserver = envRunAddr
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}
	if envTimerSend := os.Getenv("STORE_INTERVAL"); envTimerSend != "" {
		num, err := strconv.Atoi(envTimerSend)
		if err != nil {
			logger.Log.Panic("timerSend must be int", zap.Error(err))
		}
		timerSend = num
	}

	if envRESTORE := os.Getenv("RESTORE"); envRESTORE != "" {
		boolValue, err := strconv.ParseBool(envRESTORE)
		if err != nil {
			logger.Log.Panic("RESTORE must be bool", zap.Error(err))
		}
		RESTORE = boolValue
	}

	if envDATABASE := os.Getenv("DATABASE_DSN"); envDATABASE != "" {
		DatabaseDsn = envDATABASE
	}

}

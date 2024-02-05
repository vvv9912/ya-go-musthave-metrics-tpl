package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

var URLserver string
var reportInterval uint
var pollInterval uint
var KeyAuth string

func (o *NetAddress) String() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}

// Set связывает переменную типа со значением флага
// и устанавливает правила парсинга для пользовательского типа.
func (o *NetAddress) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("неправильный формат")
	}
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

	flag.Var(addr, "a", "Net address host:port")
	flag.UintVar(&reportInterval, "r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	flag.UintVar(&pollInterval, "p", 2, "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	flag.StringVar(&KeyAuth, "k", "", "key for auth (по умолчанию пустая)")
	flag.Parse()

	flagValid := map[string]struct{}{
		"a": {},
		"r": {},
		"p": {},
	}
	flag.Visit(func(f *flag.Flag) {
		_, ok := flagValid[f.Name]
		if !ok {
			flag.PrintDefaults()
			log.Panic("лишний флаг:" + f.Name)
		}
	})

	URLserver = addr.String()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		URLserver = envRunAddr
	}
	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		uintValue, err := strconv.ParseUint(envReport, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		reportInterval = uint(uintValue)
	}
	if envPoll := os.Getenv("POLL_INTERVAL"); envPoll != "" {
		uintValue, err := strconv.ParseUint(envPoll, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		pollInterval = uint(uintValue)
	}
	if envKey := os.Getenv("KEY"); envKey != "" {
		KeyAuth = envKey
	}
}

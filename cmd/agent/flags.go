package main

import (
	"encoding/json"
	"errors"
	"flag"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

type NetAddress struct {
	Host string
	Port int
}

var URLserver string
var reportInterval uint
var pollInterval uint
var KeyAuth string
var RateLimit uint
var CryptoKey string
var Config string

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

func parseJSON(filePath string, flags map[string]bool) {
	type Conf struct {
		Address        string `json:"address"`
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
		CryptoKey      string `json:"crypto_key"`
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var config Conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range flags {
		if v {
			continue
		}
		switch k {
		case "a":
			URLserver = config.Address
		case "r":
			duration, err := time.ParseDuration(config.ReportInterval)
			if err != nil {
				log.Fatal("Error parse duration", zap.Error(err))
			}
			reportInterval = uint(duration.Seconds())

		case "p":
			duration, err := time.ParseDuration(config.PollInterval)
			if err != nil {
				log.Fatal("Error parse duration", zap.Error(err))
				return
			}
			pollInterval = uint(duration.Seconds())
		case "crypto-key":
			CryptoKey = config.CryptoKey
		}
	}

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

	flag.StringVar(&KeyAuth, "k", "", "key for auth (по умолчанию пустая)")
	flag.Var(addr, "a", "Net address host:port")
	flag.UintVar(&reportInterval, "r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	flag.UintVar(&pollInterval, "p", 2, "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	flag.UintVar(&RateLimit, "l", 1, "одновременно исходящих запросов на сервер (по умолчанию 1)")
	flag.StringVar(&CryptoKey, "crypto-key", "", "crypto-key (по умолчанию, сообщения не шифруются)")
	flag.StringVar(&Config, "c", "", "config")
	flag.Parse()

	flagValid := map[string]bool{
		"a":          false,
		"r":          false,
		"p":          false,
		"k":          false,
		"l":          false,
		"crypto-key": false,
		"c":          false,
	}
	flag.Visit(func(f *flag.Flag) {
		_, ok := flagValid[f.Name]
		if !ok {
			flag.PrintDefaults()
			log.Panic("лишний флаг:" + f.Name)
		}
		flagValid[f.Name] = true
	})

	URLserver = addr.String()
	if envKey := os.Getenv("KEY"); envKey != "" {
		KeyAuth = envKey
	}
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
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		uintValue, err := strconv.ParseUint(envRateLimit, 10, 64)
		if err != nil {
			log.Panic(err)
		}
		RateLimit = uint(uintValue)
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		CryptoKey = envCryptoKey
	}

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		Config = envConfig
	}

	if Config != "" {
		parseJSON(Config, flagValid)
	}
}

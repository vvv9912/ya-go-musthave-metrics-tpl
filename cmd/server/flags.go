package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
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

var (
	URLserver       string // URL of the server
	timerSend       int    // Event sending time
	FileStoragePath string // Path to the temporary file
	RESTORE         bool   // Flag for restoring previous metrics from temporary file
	KeyAuth         string // Authentication key
	flagLogLevel    string
	DatabaseDsn     string
	CryptoKey       string
	Config          string
	trustedSubnet   string
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

func parseJSON(filePath string, flags map[string]bool) {
	type Conf struct {
		Address       string `json:"address"`
		Restore       bool   `json:"restore"`
		StoreInterval string `json:"store_interval"`
		StoreFile     string `json:"store_file"`
		DatabaseDSN   string `json:"database_dsn"`
		CryptoKey     string `json:"crypto_key"`
		TrustedSubnet string `json:"trusted_subnet"`
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
			RESTORE = config.Restore
		case "i":
			duration, err := time.ParseDuration(config.StoreInterval)
			if err != nil {
				log.Fatal("Error parse duration", zap.Error(err))
				return
			}
			timerSend = int(duration.Seconds())

		case "f":
			FileStoragePath = config.StoreFile
		case "d":
			DatabaseDsn = config.DatabaseDSN
		case "crypto-key":
			CryptoKey = config.CryptoKey
		case "t":
			trustedSubnet = config.TrustedSubnet
		}
	}

}

func parseFlags() {
	// Список зарег. переменных
	flagValid := map[string]bool{
		"k":          false,
		"a":          false,
		"l":          false,
		"f":          false,
		"i":          false,
		"r":          false,
		"d":          false,
		"crypto-key": false,
		"c":          false,
		"t":          false,
	}
	// регистрируем переменную flagRunAddr
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
	flag.StringVar(&DatabaseDsn, "d", "", "DATABASE_DSN")
	flag.StringVar(&CryptoKey, "crypto-key", "", "crypto-key (по умолчанию, сообщения не шифруются)")
	flag.StringVar(&Config, "c", "", "config")
	flag.StringVar(&trustedSubnet, "t", "", "TRUSTED_SUBNET")
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		_, ok := flagValid[f.Name]
		if ok {
			flagValid[f.Name] = true
		}
	})

	URLserver = addr.String()
	if envKey := os.Getenv("KEY"); envKey != "" {
		KeyAuth = envKey
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		URLserver = envRunAddr
		flagValid["a"] = true
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
		flagValid["f"] = true
	}
	if envTimerSend := os.Getenv("STORE_INTERVAL"); envTimerSend != "" {
		num, err := strconv.Atoi(envTimerSend)
		if err != nil {
			logger.Log.Panic("timerSend must be int", zap.Error(err))
		}
		timerSend = num
		flagValid["i"] = true
	}
	if envRESTORE := os.Getenv("RESTORE"); envRESTORE != "" {
		boolValue, err := strconv.ParseBool(envRESTORE)
		if err != nil {
			logger.Log.Panic("RESTORE must be bool", zap.Error(err))
		}
		RESTORE = boolValue
		flagValid["r"] = true
	}
	if envDATABASE := os.Getenv("DATABASE_DSN"); envDATABASE != "" {
		DatabaseDsn = envDATABASE
		flagValid["d"] = true
	}
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		CryptoKey = envCryptoKey
		flagValid["crypto-key"] = true
	}
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		Config = envConfig
	}
	if envTrustSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustSubnet != "" {
		trustedSubnet = envTrustSubnet
	}

	if Config != "" {
		parseJSON(Config, flagValid)
	}
}

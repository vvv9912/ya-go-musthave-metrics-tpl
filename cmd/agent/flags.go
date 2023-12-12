package main

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (o *NetAddress) String() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}

var URLserver string
var reportInterval uint
var pollInterval uint

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
	_ = flag.Value(addr)

	flag.Var(addr, "a", "Net address host:port")
	flag.UintVar(&reportInterval, "r", 10, "частота отправки метрик на сервер (по умолчанию 10 секунд)")
	flag.UintVar(&pollInterval, "p", 2, "частота опроса метрик из пакета runtime (по умолчанию 2 секунды)")
	flag.Parse()

	flagValid := map[string]struct{}{
		"a": {},
		"r": {},
		"p": {},
	}
	flag.Visit(func(f *flag.Flag) {
		_, ok := flagValid[f.Name]
		if !ok {
			panic("лишний флаг:" + f.Name)
		}
	})
	//значения по умолчанию
	//fmt.Println(flag.Lookup("a"))
	//if flag.Lookup("a") == nil {
	//	addr.Host = "localhost"
	//	addr.Port = 8080
	//	fmt.Println(addr)
	//}
	URLserver = addr.String()
}

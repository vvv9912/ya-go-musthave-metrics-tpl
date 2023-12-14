package main

import (
	"errors"
	"flag"
	"os"
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
	_ = flag.Value(addr)
	//flag.StringVar(addr, "a", "localhost:8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	//lag.StringVar(addr, "a", "localhost:8080", "address and port to run server")

	flag.Var(addr, "a", "Net address host:port")
	flag.Parse()
	URLserver = addr.String()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		URLserver = envRunAddr
	}
}

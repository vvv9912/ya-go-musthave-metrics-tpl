package server

import (
	"context"
	"fmt"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"net/http"
)

type Server struct {
	s *http.Server
}

func NewServer() *Server {
	s := &http.Server{}
	return &Server{s: s}
}

func (s *Server) StartServer(ctx context.Context, addr string, gaugeStorage storage.GaugeStorager, counterStorage storage.CounterStorager) error {
	mux := http.NewServeMux()
	m := mw.Mw{
		GaugeStorage:   gaugeStorage,
		CounterStorage: counterStorage,
	}
	//http://localhost:8080/update/unknown/testCounter/100
	mux.Handle("/update/gauge/", m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge))))
	mux.Handle("/update/counter/", m.Middlware(m.MiddlwareCounter(http.HandlerFunc(handler.HandlerCounter))))
	s.s.Addr = addr
	s.s.Handler = mux
	ch := make(chan error)
	go func() {
		defer s.s.Close()
		fmt.Println("server start, addr:", addr)
		err := s.s.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			ch <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		return s.s.Close() //s.s.Shutdown(context.Background())
	case err := <-ch:
		return err
	}
	//return nil
}

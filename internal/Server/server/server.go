package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"net/http"
)

//type Server struct {
//	s *http.Server
//}

type Server struct {
	s *chi.Mux
}

func NewServer() *Server {
	s := chi.NewRouter()
	return &Server{s: s}
}

func (s *Server) StartServer(ctx context.Context, addr string, gaugeStorage storage.GaugeStorager, counterStorage storage.CounterStorager) error {
	m := mw.Mw{
		GaugeStorage:   gaugeStorage,
		CounterStorage: counterStorage,
	}
	s.s.With(m.Middlware2Gauge).Post("/update/gauge/{SomeMetric}/{Value}", handler.HandlerGauge)
	s.s.With(m.Middlware2Counter).Post("/update/counter/{SomeMetric}/{Value}", handler.HandlerCounter)
	s.s.With(m.MiddlwareGetGauge).Get("/value/gauge/{SomeMetric}", handler.HandlerGetGauge)
	s.s.With(m.MiddlwareGetCounter).Get("/value/counter/{SomeMetric}", handler.HandlerGetCounter)
	s.s.Get("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.HandlerGetDef(res, req, gaugeStorage, counterStorage)
	}))
	ch := make(chan error)
	go func() {
		//defer s.s.Close()
		fmt.Println("server start, addr:", addr)
		err := http.ListenAndServe(addr, s.s)
		if err != nil {
			fmt.Println(err)
			ch <- err
			return
		}
	}()

	//mux := http.NewServeMux()
	//m := mw.Mw{
	//	GaugeStorage:   gaugeStorage,
	//	CounterStorage: counterStorage,
	//}
	////http://localhost:8080/update/unknown/testCounter/100
	//mux.Handle("/update/gauge/", m.Middlware(m.MiddlwareGauge(http.HandlerFunc(handler.HandlerGauge))))
	//mux.Handle("/update/counter/", m.Middlware(m.MiddlwareCounter(http.HandlerFunc(handler.HandlerCounter))))
	//mux.Handle("/", http.HandlerFunc(handler.HandlerBase))
	//s.s.Addr = addr
	//s.s.Handler = mux
	//ch := make(chan error)
	//go func() {
	//	defer s.s.Close()
	//	fmt.Println("server start, addr:", addr)
	//	err := s.s.ListenAndServe()
	//	if err != nil {
	//		fmt.Println(err)
	//		ch <- err
	//		return
	//	}
	//}()
	//

	//select {
	//case <-ctx.Done():
	//	return s.s.Close() //s.s.Shutdown(context.Background())
	//case err := <-ch:
	//	return err
	//}
	return nil
}

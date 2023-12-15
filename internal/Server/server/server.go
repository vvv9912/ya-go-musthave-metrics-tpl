package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
	"net/http"
	"os"
)

type Server struct {
	s      *chi.Mux
	Logger *log.Logger
}

func NewServer() *Server {
	s := chi.NewRouter()
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	logger := log.New(f, "server: ", log.LstdFlags)
	logger.Println("Server start")
	return &Server{s: s, Logger: logger}
}

func (s *Server) StartServer(ctx context.Context, addr string, gaugeStorage storage.GaugeStorager, counterStorage storage.CounterStorager) error {

	m := mw.Mw{
		GaugeStorage:   gaugeStorage,
		CounterStorage: counterStorage,
		Log:            s.Logger,
	}
	s.s.Use(m.MwLogger)

	s.s.With(m.MiddlewareType).Post("/update/{type}/{SomeMetric}/{Value}", handler.HandlerSucess)
	s.s.With(m.MiddlwareGetCounter).Get("/value/counter/{SomeMetric}", handler.HandlerGetCounter)
	s.s.With(m.MiddlwareGetGauge).Get("/value/gauge/{SomeMetric}", handler.HandlerGetGauge)

	s.s.Post("/update/", handler.HandlerErrType)

	s.s.Get("/", handler.HandlerGetMetrics(gaugeStorage, counterStorage))

	server := http.Server{
		Addr:    addr,
		Handler: s.s,
	}

	ch := make(chan error)

	go func() {

		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
			ch <- err
			return
		}
		log.Println("server start, addr:", addr)
	}()

	select {
	case <-ctx.Done():
		return server.Shutdown(context.Background()) //s.s.Shutdown(context.Background())
	case err := <-ch:
		return err
	}
}

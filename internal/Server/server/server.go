package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
	"net/http"
	"os"
)

//	type Server struct {
//		s *http.Server
//	}

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
	//s.s.Route("/update", func(r chi.Router) {
	//	r.Route("/gauge", func(r chi.Router) {
	//		r.Use(m.Middlware2Gauge)
	//		r.Post("/{SomeMetric}/{Value}", handler.HandlerGauge)
	//	})
	//	r.Route("/counter", func(r chi.Router) {
	//		r.Use(m.Middlware2Counter)
	//		r.Post("/{SomeMetric}/{Value}", handler.HandlerCounter)
	//	})
	//	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
	//		writer.WriteHeader(http.StatusBadRequest)
	//		writer.Write([]byte("Unknown endpoint"))
	//	})
	//})
	s.s.With(m.MiddlewareType).Post("/update/{type}/{SomeMetric}/{Value}", handler.HandlerSucess)
	s.s.With(m.MiddlwareGetCounter).Get("/value/counter/{SomeMetric}", handler.HandlerGetCounter)
	s.s.With(m.MiddlwareGetGauge).Get("/value/gauge/{SomeMetric}", handler.HandlerGetGauge)
	s.s.Post("/update/", handler.HandlerErrType)
	s.s.Get("/", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		handler.HandlerGetDef(res, req, gaugeStorage, counterStorage)
	}))
	ch := make(chan error)
	server := http.Server{}
	server.Handler = s.s
	server.Addr = addr
	go func() {
		log.Println("server start, addr:", addr)
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			ch <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		return server.Shutdown(context.Background()) //s.s.Shutdown(context.Background())
	case err := <-ch:
		return err
	}
}

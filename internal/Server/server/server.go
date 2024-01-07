package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/project"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type Server struct {
	s *chi.Mux
}

func NewServer() *Server {
	s := chi.NewRouter()
	return &Server{s: s}
}

func (s *Server) StartServer(ctx context.Context, addr string, gaugeStorage storage.GaugeStorager, counterStorage storage.CounterStorager) error {
	p := project.NewProject(counterStorage, gaugeStorage)
	h := handler.NewHandler(*p)
	m := mw.Mw{
		GaugeStorage:   gaugeStorage,
		CounterStorage: counterStorage,
		//	Log:            s.Logger,
	}
	s.s.Use(m.MwLogger)

	s.s.With(m.MiddlewareType).Post("/update/{type}/{SomeMetric}/{Value}", handler.HandlerSucess)
	s.s.With(m.MiddlwareGetCounter).Get("/value/counter/{SomeMetric}", handler.HandlerGetCounter)
	s.s.With(m.MiddlwareGetGauge).Get("/value/gauge/{SomeMetric}", handler.HandlerGetGauge)

	s.s.With(m.MiddlwareCheckJson).Post("/update/", h.HandlerGetJSON)

	//	s.s.Post("/update/", handler.HandlerErrType)

	s.s.Get("/", handler.HandlerGetMetrics(gaugeStorage, counterStorage))

	server := http.Server{
		Addr:    addr,
		Handler: s.s,
	}

	ctxServer, cancel := context.WithCancel(ctx)

	go func() {
		logger.Log.Info("server start", zap.String("addr", addr))
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("ctx:", zap.Error(ctx.Err()))
		return server.Shutdown(context.Background())
	case <-ctxServer.Done():
		logger.Log.Info("ctx:", zap.Error(ctxServer.Err()))
		return errors.New("canceled by ctxServer")
	}
}

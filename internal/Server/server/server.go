package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/handler"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/mw"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/notifier"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/service"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Server struct {
	s *chi.Mux
}

func NewServer() *Server {

	s := chi.NewRouter()
	return &Server{s: s}
}

func (s *Server) StartServer(
	ctx context.Context,
	addr string,
	gaugeStorage storage.GaugeStorager,
	counterStorage storage.CounterStorager, timeSend time.Duration, writer notifier.Writer) error {

	var (
		e       = notifier.NewNotifier(gaugeStorage, counterStorage, timeSend, writer)
		Service = service.NewService(counterStorage, gaugeStorage, e)
		h       = handler.NewHandler(Service)
		m       = mw.NewMw(Service)
	)

	s.s.Use(m.MwLogger)
	s.s.Use(m.MiddlewareGzip)

	s.s.With(m.MiddlewareType).Post("/update/{type}/{SomeMetric}/{Value}", handler.HandlerSucess)

	s.s.With(m.MiddlwareGetCounter).Get("/value/counter/{SomeMetric}", handler.HandlerGetCounter)
	s.s.With(m.MiddlwareGetGauge).Get("/value/gauge/{SomeMetric}", handler.HandlerGetGauge)

	s.s.With(m.MiddlwareCheckJSON).Post("/update/", h.HandlerPostJSON)
	s.s.With(m.MiddlwareCheckJSON).Post("/value/", h.HandlerGetJSON)

	s.s.Get("/", handler.HandlerGetMetrics(gaugeStorage, counterStorage))

	server := http.Server{
		Addr:    addr,
		Handler: s.s,
	}

	ctxServer, cancel := context.WithCancel(ctx)

	e.StartNotifier(ctxServer)

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

package main

import (
	"context"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/server"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var URLserver string

var timerSend int
var FileStoragePath string
var RESTORE bool

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}
func run() error {
	log.Println("Start server")
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("URLserver=" + URLserver)
	logger.Log.Info("timerSend=", zap.Int("timerSend", timerSend))
	logger.Log.Info("FileStoragePath=" + FileStoragePath)
	logger.Log.Info("Restore=", zap.Bool("RESTORE", RESTORE))
	counter := storage.NewCounterStorage()
	gauge := storage.NewGaugeStorage()

	//	RESTORE = false
	if RESTORE {

		consumer, err := fileutils.NewConsumer(FileStoragePath)
		if err != nil {
			logger.Log.Info("error consumer", zap.Error(err))
			//return err
		}
		event, err := consumer.ReadLastEvent(FileStoragePath)
		if err != nil {
			logger.Log.Info("error read last event", zap.Error(err))

		}

		if event != nil {
			for key, val := range event.Counter {
				err = counter.UpdateCounter(key, val)
				if err != nil {
					logger.Log.Info("error update counter", zap.Error(err))
				}
			}
			for key, val := range event.Gauge {
				err = gauge.UpdateGauge(key, val)
				if err != nil {
					logger.Log.Info("error update gauge", zap.Error(err))
				}
			}
		}

		defer consumer.Close()
	}
	s := server.NewServer()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	produce, err := fileutils.NewProducer(FileStoragePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer produce.Close()

	err = s.StartServer(ctx, URLserver, gauge, counter, time.Duration(timerSend)*time.Second, produce)
	if err != nil {
		log.Println(err)
		return err
	}

	<-ctx.Done()
	return nil
}

package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/fileutils"
	"github.com/vvv9912/ya-go-musthave-metrics-tpl.git/internal/Server/storage"
	"log"
	"testing"
	"time"
)

func TestServer_StartServer(t *testing.T) {
	type fields struct {
		s *chi.Mux
	}
	type args struct {
		ctx            context.Context
		addr           string
		gaugeStorage   storage.GaugeStorager
		counterStorage storage.CounterStorager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				s: tt.fields.s,
			}
			produce, err := fileutils.NewProducer("test.json")
			if err != nil {
				log.Println(err)
			}
			defer produce.Close()
			if err := s.StartServer(tt.args.ctx, tt.args.addr, tt.args.gaugeStorage, tt.args.counterStorage, time.Duration(1*time.Second), produce); (err != nil) != tt.wantErr {
				t.Errorf("StartServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package fileutils

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

type Event struct {
	Gauge   map[string]float64 `json:"gauge,omitempty"`
	Counter map[string]uint64  `json:"counter,omitempty"`
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	err := os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(fileName string) (*Consumer, error) {
	if fileName[0] == '/' {
		fileName = fileName[1:]
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
	event := &Event{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]uint64),
	}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}
func (c *Consumer) ReadLastEvent(fileName string) (*Event, error) {
	consumer, err := NewConsumer(fileName)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	var lastEvent *Event
	for {
		event, err := consumer.ReadEvent()
		if err != nil {
			if err == io.EOF { // Достигнут конец файла
				break
			}
			return nil, err
		}
		lastEvent = event
	}

	return lastEvent, nil
}
func (c *Consumer) Close() error {
	return c.file.Close()
}

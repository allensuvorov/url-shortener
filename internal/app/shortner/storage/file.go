package storage

import (
	"encoding/json"
	"log"
	"os"
)

// Write to file
type writer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWriter(fileName string) (*writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *writer) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}
func (p *writer) Close() error {
	return p.file.Close()
}

func write(h, u, fsp string) error {
	log.Printf("Storage: saving to path - %s", fsp)
	// path will be /Users/allen/go/src/yandex/projects/urlshortner/internal/app/shortner/storage
	fileName := "urls.txt"
	writer, err := NewWriter(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	if err := writer.WriteEvent(event); err != nil {
		log.Fatal(err)
	}

	return nil
}

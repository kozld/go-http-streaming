package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kozld/go-http-streaming/streaming.go"
)

const (
	endpoint = "http://localhost:9094"
	filePath = "examples/upload/static/book.txt"
)

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("got error: %v", err)
	}

	log.Printf("Filename: %s; Size: %d\n", fileInfo.Name(), fileInfo.Size())

	startTime := time.Now()
	err = streaming.Stream(http.MethodPost, endpoint, file)
	if err != nil {
		log.Fatalf("got error: %v", err)
	}
	endTime := time.Now()

	executionTime := endTime.Sub(startTime)
	log.Printf("Upload completed in %v\n", executionTime)
}

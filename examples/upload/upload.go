package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kozld/go-http-streaming/client.go"
)

const (
	endpoint = "http://localhost:3000"
	filePath = "examples/upload/payload/book.txt"
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
	err = client.Stream(http.MethodPost, endpoint, file)
	endTime := time.Now()

	if err != nil {
		log.Fatalf("got error: %v", err)
	}

	executionTime := endTime.Sub(startTime)
	log.Printf("Upload completed in %v\n", executionTime)
}

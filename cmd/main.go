package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kozld/go-http-streaming/client.go"
)

const (
	defaultUploadHost = "http://localhost:3000"
)

var (
	uploadHost              = os.Getenv("UPLOAD_HOST")
	errFilePathNotSpecified = errors.New("flag \"-filepath\" not specified")
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "filepath", "", "a file path")
	flag.Parse()

	if filepath == "" {
		log.Fatalln(errFilePathNotSpecified)
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Filename: %s; Size: %d\n", fileInfo.Name(), fileInfo.Size())

	endpoint := defaultUploadHost
	if uploadHost != "" {
		endpoint = uploadHost
	}

	c := client.New(http.MethodPost, endpoint)

	startTime := time.Now()
	resp, err := c.Send(file)
	endTime := time.Now()

	if err != nil {
		log.Fatalln(err)
	}

	executionTime := endTime.Sub(startTime)
	log.Printf("Upload completed in %v\n", executionTime)

	bytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", string(bytes))
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
)

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	processStream(req, w)

	fmt.Fprintf(w, "Upload successful")
}

func processStream(r *http.Request, w http.ResponseWriter) {
	scanner := bufio.NewScanner(r.Body)
	defer r.Body.Close()

	numBytes, numChunks := int64(0), int64(0)

	for scanner.Scan() {
		line := scanner.Bytes()
		//log.Println("Received line:", string(line))

		numBytes += int64(len(line))
		numChunks += 1
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading stream:", err)
		http.Error(w, "Error reading stream", http.StatusInternalServerError)
	}

	log.Println("Bytes:", numBytes, "Chunks:", numChunks)
}

func main() {
	http.HandleFunc("/", upload)

	http.ListenAndServe(":9094", nil)
}

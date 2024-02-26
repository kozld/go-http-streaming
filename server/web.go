package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port = 3000
)

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	processStream(req, w)

	log.Println("Upload successful")
}

func processStream(r *http.Request, w http.ResponseWriter) {
	defer r.Body.Close()

	data := make([]byte, 32<<20)
	numBytes := int64(0)

	for {
		n, err := r.Body.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		//fmt.Println(string(data[:n]))
		numBytes += int64(n)
	}

	log.Printf("Recieved bytes: %d\n", numBytes)
	fmt.Fprintf(w, "Upload successful")
}

func main() {
	http.HandleFunc("/", upload)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

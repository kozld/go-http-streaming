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
	defer r.Body.Close()

	reader := bufio.NewReader(r.Body)

	data := make([]byte, 32<<20)
	numBytes := int64(0)

	for {
		n, err := reader.Read(data)
		if err != nil {
			break
		}

		// fmt.Println(string(data[:n]))

		numBytes += int64(n)
	}

	log.Println("Bytes:", numBytes)
}

func main() {
	http.HandleFunc("/", upload)
	http.ListenAndServe(":9094", nil)
}

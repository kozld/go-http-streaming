package client

import (
	"bufio"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	Boundary     = "---XXX---"
	BatchSize    = 500 * 1024
	WorkersCount = 100
)

func Stream(method string, url string, in io.Reader) error {
	ch := Read(in, BatchSize)

	parts := make([]io.Reader, 0, WorkersCount)
	for i := 0; i < WorkersCount; i += 1 {
		parts = append(parts, Do(ch))
	}

	body := Merge(parts...)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)

	return err
}

func Read(in io.Reader, size int) <-chan []byte {
	out := make(chan []byte)

	data := make([]byte, size)
	reader := bufio.NewReader(in)

	go func() {
		for {
			n, err := reader.Read(data)
			if err != nil {
				break
			}

			out <- data[:n]
		}

		close(out)
	}()

	return out
}

func Do(in <-chan []byte) io.Reader {
	reader, writer := io.Pipe()

	go func() {
		for part := range in {
			w := multipart.NewWriter(writer)
			w.SetBoundary(Boundary)

			file, _ := w.CreateFormFile("file", "file")
			file.Write(part)

			w.Close()
		}

		writer.Close()
	}()

	return reader
}

func Merge(parts ...io.Reader) io.Reader {
	return io.MultiReader(parts...)
}

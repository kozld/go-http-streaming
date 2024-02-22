package streaming

import (
	"io"
	"mime/multipart"
	"net/http"
)

const (
	chunkSize    = 500 * 1024
	workersCount = 100
)

func Stream(method string, url string, in io.Reader) error {
	out := read(in, chunkSize)

	parts := make([]io.Reader, 0, workersCount)
	for i := 0; i < workersCount; i += 1 {
		parts = append(parts, do(out))
	}

	body := merge(parts...)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	return err
}

func read(in io.Reader, size int) <-chan []byte {
	out := make(chan []byte)

	go func() {
		for {
			buf := make([]byte, size)
			n, _ := in.Read(buf)
			buf = buf[:n]

			if n == 0 {
				close(out)
				break
			}

			out <- buf
		}
	}()

	return out
}

func do(in <-chan []byte) *io.PipeReader {
	reader, writer := io.Pipe()

	mw := multipart.NewWriter(writer)

	go func() {
		for part := range in {
			mw.WriteField("fieldname", string(part))
		}

		writer.Close()
		mw.Close()
	}()

	return reader
}

func merge(bodies ...io.Reader) io.Reader {
	return io.MultiReader(bodies...)
}

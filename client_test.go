package client_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/kozld/go-http-streaming/client.go"
	"github.com/kozld/go-http-streaming/client.go/testdata"
)

type (
	HTTPClientMock struct{}
)

const (
	endpoint        = "http://upload-server.com"
	fixtureTestFile = "testdata/test_file.txt"
)

func TestSend(t *testing.T) {
	file, err := testdata.Generate(fixtureTestFile)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	expectedFileSize := fileInfo.Size()

	httpMock := &HTTPClientMock{}
	c := client.New(http.MethodPost, endpoint, client.WithHttpClient(httpMock))

	resp, err := c.Send(file)
	if err != nil {
		t.Fatalf(err.Error())
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fileSize := int64(len(bytes))

	if fileSize != expectedFileSize {
		t.Errorf("Result was incorrect, got: %d, want: %d.", fileSize, expectedFileSize)
	}
}

func (c *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	buf := new(bytes.Buffer)
	mr := multipart.NewReader(req.Body, client.Boundary)
	for {
		form, err := mr.ReadForm(32 << 20)
		if err != nil {
			break
		}

		header := form.File["file"]
		file, _ := header[0].Open()

		data := make([]byte, 32<<20)
		for {
			n, err := file.Read(data)
			if err != nil {
				break
			}

			buf.Write(data[:n])
		}
	}

	return &http.Response{Body: io.NopCloser(buf)}, nil
}

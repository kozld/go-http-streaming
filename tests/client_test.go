package client_test

import (
	"mime/multipart"
	"testing"

	"github.com/kozld/go-http-streaming/client.go"
	"github.com/kozld/go-http-streaming/client.go/tests/fixtures"
)

const (
	fixtureTestRead = "fixtures/test_read.txt"
	fixtureTestDo   = "fixtures/test_do.txt"

	batchSize = 200 * 1024
)

func TestRead(t *testing.T) {
	file, err := fixtures.Generate(fixtureTestRead)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	expectedFileSize := fileInfo.Size()

	ch := client.Read(file, batchSize)

	fileSize := int64(0)
	for part := range ch {
		n := len(part)
		fileSize += int64(n)
	}

	if fileSize != expectedFileSize {
		t.Errorf("Result was incorrect, got: %d, want: %d.", fileSize, expectedFileSize)
	}
}

func TestDo(t *testing.T) {
	file, err := fixtures.Generate(fixtureTestDo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	expectedFileSize := fileInfo.Size()

	ch := client.Read(file, batchSize)
	parts := client.Do(ch)

	reader := multipart.NewReader(parts, client.Boundary)

	fileSize := int64(0)
	for {
		form, err := reader.ReadForm(32 << 20)
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

			fileSize += int64(n)
		}
	}

	if fileSize != expectedFileSize {
		t.Errorf("Result was incorrect, got: %d, want: %d.", fileSize, expectedFileSize)
	}
}

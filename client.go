package client

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type (
	Client struct {
		http   HTTPClient
		method string
		url    string

		pr *io.PipeReader
		pw *io.PipeWriter

		chErr  chan error
		chResp chan *http.Response

		closed bool
	}

	Opt func(c *Client)

	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

const (
	Boundary = "---XXX---"
)

var (
	errConnectionClosed = errors.New("error: connection closed")
)

func New(method string, url string, opts ...Opt) *Client {
	c := &Client{
		method: method,
		url:    url,
		chErr:  make(chan error),
		chResp: make(chan *http.Response),
	}

	c.pr, c.pw = io.Pipe()

	for _, opt := range DefaultOpts() {
		opt(c)
	}
	for _, opt := range opts {
		opt(c)
	}

	go c.handleStream()

	return c
}

func DefaultOpts() []Opt {
	return []Opt{
		WithHttpClient(&http.Client{}),
	}
}

// WithHttpClient - Allows set custom HTTP client for testing purposes
func WithHttpClient(client HTTPClient) Opt {
	return func(c *Client) {
		c.http = client
	}
}

// Send - Sends data in multipart/form-data format,
// after sending, closes the HTTP connection (sending an EOF signal) and returns *http.Response
func (c *Client) Send(r io.Reader) (*http.Response, error) {
	err := c.writeFrom(r)
	if err != nil {
		return nil, err
	}

	return c.close()
}

// SendUnclose - Sends data in multipart/form-data format,
// after sending the HTTP connection remains open
func (c *Client) SendUnclose(r io.Reader) error {
	return c.writeFrom(r)
}

func (c *Client) writeFrom(r io.Reader) error {
	mw := multipart.NewWriter(c.pw)
	mw.SetBoundary(Boundary)

	w, err := mw.CreateFormFile("file", "file")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	return mw.Close()
}

// Close - Closes the current HTTP connection (sending an EOF signal) and returns http.Response
func (c *Client) Close() (*http.Response, error) {
	return c.close()
}

func (c *Client) handleStream() {
	for {
		req, _ := http.NewRequest(c.method, c.url, c.pr)
		resp, err := c.http.Do(req)

		if !c.closed {
			continue
		}

		if err != nil {
			c.chErr <- err
		} else {
			c.chResp <- resp
		}

		c.closed = false
	}
}

func (c *Client) close() (*http.Response, error) {
	if c.closed {
		return nil, errConnectionClosed
	}

	c.closed = true
	c.pw.Close()

	select {
	case resp := <-c.chResp:
		return resp, nil
	case err := <-c.chErr:
		return nil, err
	}
}

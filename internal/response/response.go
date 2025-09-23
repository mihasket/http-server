package response

import (
	"fmt"
	"http-server-miha/internal/headers"
	"io"
	"strconv"
)

const HTTP_VERSION = "HTTP/1.1"

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

var statusText = map[StatusCode]string{
	OK:                  "OK",
	BadRequest:          "Bad Request",
	InternalServerError: "Internal Server Error",
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	reasonPhrase := statusText[statusCode]

	message := fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, statusCode, reasonPhrase)

	_, err := w.writer.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	err := h.WriteHeaders(w.writer)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return -1, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) error {
	size := len(p)
	w.WriteBody(fmt.Appendf(nil, "%X\r\n", size))
	w.WriteBody(p[:size])
	w.WriteBody([]byte("\r\n"))

	return nil
}

func (w *Writer) WriteChunkedBodyDone() error {
	w.WriteBody([]byte("0\r\n\r\n"))

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/html")

	return *h
}

func Respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Sucess!</h1>
    <p>You request was sucessful</p>
  </body>
</html>
`)
}

func Respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>You request was not sucessful</p>
  </body>
</html>
`)
}

func Respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>We made a mistake</p>
  </body>
</html>
`)
}

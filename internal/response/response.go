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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase := statusText[statusCode]

	message := fmt.Sprintf("%s %d %s\r\n", HTTP_VERSION, statusCode, reasonPhrase)

	_, err := w.Write([]byte(message))
	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return *h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	err := headers.WriteHeaders(w)
	if err != nil {
		return err
	}

	return nil
}


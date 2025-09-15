package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERR_REQUEST_LINE = errors.New("request line parsing error")
var ERR_HTTP_VERSION = errors.New("request line - incorrect HTTP version")

func hasAllCapital(s string) error {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return ERR_REQUEST_LINE
		}
	}
	return nil
}

// Parse only the first request line
func parseRequestLine(b []byte) (*RequestLine, string, error) {
	reqLineEnd := strings.Index(string(b), "\r\n")
	if reqLineEnd == -1 {
		return nil, "", ERR_REQUEST_LINE
	}

	line := string(b[:reqLineEnd])
	rest := string(b[reqLineEnd+len("\r\n"):])
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return nil, rest, ERR_REQUEST_LINE
	}

	err := hasAllCapital(parts[0])
	if err != nil {
		return nil, rest, err
	}

	httpVersionParts := strings.Split(parts[2], "/")
	if len(httpVersionParts) != 2 || httpVersionParts[0] != "HTTP" || httpVersionParts[1] != "1.1" {
		return nil, rest, ERR_HTTP_VERSION
	}

	reqLine := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpVersionParts[1],
	}

	return reqLine, rest, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, _, err := parseRequestLine(b)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

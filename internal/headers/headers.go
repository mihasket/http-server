package headers

import (
	"bytes"
	"errors"
	"strings"
)

var ERR_WS_BEFORE_COLON = errors.New("White space before colon in field line")
var ERR_NOT_A_HEADER = errors.New("Not a header line")
var ERR_INVALID_CHARACTERS = errors.New("Invalid header characters")

var COLON = []byte(":")
var CRLF = []byte("\r\n")
var WS = []byte(" ")

type headers map[string]string

type Headers struct {
	headers headers
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func isValid(s string) error {
	if len(s) < 1 {
		return ERR_INVALID_CHARACTERS
	}

	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			r == '!' || r == '#' || r == '$' || r == '%' ||
			r == '&' || r == '\'' || r == '*' || r == '+' ||
			r == '-' || r == '.' || r == '^' || r == '_' ||
			r == '`' || r == '|' || r == '~') {
			return ERR_INVALID_CHARACTERS
		}
	}

	return nil
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) {
	h.headers[strings.ToLower(name)] = value
}

func parseHeader(data []byte) (headerName string, headerValue string, err error) {
	headerNameIdx := bytes.Index(data, []byte(COLON))
	if headerNameIdx == -1 {
		return "", "", ERR_NOT_A_HEADER
	}

	if string(data[headerNameIdx-1]) == string(WS) {
		return "", "", ERR_WS_BEFORE_COLON
	}

	headerName = string(bytes.ReplaceAll(data[:headerNameIdx], WS, []byte("")))
	headerValue = string(bytes.ReplaceAll(data[headerNameIdx+len(COLON):], WS, []byte("")))

	err = isValid(headerName)
	if err != nil {
		return "", "", err
	}

	return headerName, headerValue, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	readIndex := 0

	for {
		headerLineEnd := bytes.Index(data[readIndex:], []byte(CRLF))
		if headerLineEnd == -1 {
			break
		}

		// CRLF is at the start, meaning end of headers
		if headerLineEnd == 0 {
			return readIndex, true, nil
		}

		headerName, headerValue, err := parseHeader(data[readIndex : readIndex+headerLineEnd])
		if err != nil {
			return 0, false, err
		}

		readIndex += headerLineEnd + len(CRLF)
		h.Set(headerName, headerValue)
	}

	return readIndex, done, nil
}

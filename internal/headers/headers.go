package headers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ERR_WS_BEFORE_COLON = errors.New("White space before colon in field line")
var ERR_NOT_A_HEADER = errors.New("Not a header line")
var ERR_INVALID_CHARACTERS = errors.New("Invalid header characters")
var ERR_DUPLICATE_HEADERS = errors.New("Duplicate headers in request")

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

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name string, value string) error {
	name = strings.ToLower(name)
	val, ok := h.headers[name]

	if ok && val == value {
		return ERR_DUPLICATE_HEADERS
	}

	if ok {
		h.headers[name] += fmt.Sprintf(", %s", value)
	} else {
		h.headers[name] = value
	}

	return nil
}

func (h *Headers) Output() {
	for key, value := range h.headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
}

func (h *Headers) WriteHeaders(w io.Writer) error {
	for key, value := range h.headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", key, value)

		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}

	return nil
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

func parseHeader(data []byte) (headerName string, headerValue string, err error) {
	headerNameIdx := bytes.Index(data, []byte(COLON))
	if headerNameIdx == -1 {
		return "", "", ERR_NOT_A_HEADER
	}

	if string(data[headerNameIdx-1]) == string(WS) {
		return "", "", ERR_WS_BEFORE_COLON
	}

	// TODO: this
	if bytes.Contains(data[:headerNameIdx], WS) {
		return "", "", ERR_NOT_A_HEADER
	}

	headerName = string(data[:headerNameIdx])
	headerValue = strings.TrimSpace(string(data[headerNameIdx+len(COLON):]))

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
			readIndex += len(CRLF)
			return readIndex, true, nil
		}

		headerName, headerValue, err := parseHeader(data[readIndex : readIndex+headerLineEnd])
		if err != nil {
			return 0, false, err
		}

		readIndex += headerLineEnd + len(CRLF)
		err = h.Set(headerName, headerValue)
		if err != nil {
			return 0, false, err
		}
	}

	return readIndex, false, nil
}

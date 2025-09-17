package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

var ERR_WS_BEFORE_COLON = errors.New("White space before colon in field line")
var ERR_NOT_A_HEADER = errors.New("Not a header line")

var COLON = []byte(":")
var CRLF = []byte("\r\n")
var WS = []byte(" ")

func NewHeaders() Headers {
	return map[string]string{}
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

	return headerName, headerValue, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	done = false
	readIndex := 0

	for {
		headerLineEnd := bytes.Index(data[readIndex:], []byte(CRLF))
		if headerLineEnd == -1 {
			break
		}

		// CRLF is at the start, meaning end of headers
		if headerLineEnd == 0 {
			done = true
			break
		}

		headerName, headerValue, err := parseHeader(data[readIndex : readIndex+headerLineEnd])
		if err != nil {
			return 0, false, err
		}

		readIndex += headerLineEnd + len(CRLF)
		h[headerName] = headerValue
	}

	return readIndex, done, nil
}

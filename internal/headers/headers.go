package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

var ERR_WS_BEFORE_COLON = errors.New("White space before colon in field line")
var ERR_NOT_A_HEADER = errors.New("Not a header line")

const COLON = ":"
const CRLF = "\r\n"
const WS = " "

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Not enough data
	headerLineEnd := bytes.Index(data, []byte(CRLF))
	if headerLineEnd == -1 {
		return 0, false, nil
	}

	// CRLF is at the start, meaning end of headers
	if headerLineEnd == 0 {
		return 0, true, nil
	}

	headerNameIdx := bytes.Index(data, []byte(COLON))
	if headerNameIdx == -1 {
		return 0, false, ERR_NOT_A_HEADER
	}

	if string(data[headerNameIdx-1]) == WS {
		return 0, false, ERR_WS_BEFORE_COLON
	}

	headerName := string(bytes.ReplaceAll(data[:headerNameIdx], []byte(WS), []byte("")))
	headerValue := string(bytes.ReplaceAll(data[headerNameIdx+len(COLON):headerLineEnd], []byte(WS), []byte("")))

	h[headerName] = headerValue

	return headerLineEnd + len(CRLF), false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}

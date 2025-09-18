package request

import (
	"bytes"
	"errors"
	"http-server-miha/internal/headers"
	"io"
)

var ERR_REQUEST_LINE = errors.New("request line parsing error")
var ERR_HTTP_VERSION = errors.New("request line - incorrect HTTP version")
var ERR_READ_DONE_STATE = errors.New("Trying to read data in done state")
var ERR_UNKNOWN_STATE = errors.New("Unknown state in parse")

const CRLF = "\r\n"
const bufferSize = 1024

type ParserState int

const (
	StateInit ParserState = iota
	StateHeader
	StateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       ParserState
}

func hasAllCapital(b []byte) error {
	for _, r := range b {
		if r < 'A' || r > 'Z' {
			return ERR_REQUEST_LINE
		}
	}
	return nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	reqLineEnd := bytes.Index(b, []byte(CRLF))
	if reqLineEnd == -1 {
		return nil, 0, nil
	}

	line := b[:reqLineEnd]
	parts := bytes.Split(line, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ERR_REQUEST_LINE
	}

	err := hasAllCapital(parts[0])
	if err != nil {
		return nil, 0, err
	}

	httpVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpVersionParts) != 2 || string(httpVersionParts[0]) != "HTTP" || string(httpVersionParts[1]) != "1.1" {
		return nil, 0, ERR_HTTP_VERSION
	}

	reqLine := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpVersionParts[1]),
	}

	rest := reqLineEnd + len(CRLF)
	return reqLine, rest, nil
}

func (r *Request) parse(data []byte) (int, error) {
	readIndex := 0

	for {
		currentData := data[readIndex:]

		switch r.State {
		case StateDone:
			return 0, ERR_READ_DONE_STATE
		case StateInit:
			reqLine, consumedBytes, err := parseRequestLine(currentData)
			if err != nil {
				return 0, err
			}

			// Needs more data
			if consumedBytes == 0 {
				return readIndex, nil
			}

			r.RequestLine = *reqLine
			readIndex += consumedBytes
			r.State = StateHeader
		case StateHeader:
			consumedBytes, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}

			// Needs more data
			if consumedBytes == 0 {
				return readIndex, nil
			}

			readIndex += consumedBytes

			if done {
				r.State = StateDone
				return readIndex, nil
			}
		default:
			return 0, ERR_UNKNOWN_STATE
		}
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{
		RequestLine: RequestLine{},
		State:       StateInit,
		Headers:     *headers.NewHeaders(),
	}

	// TODO: Change it back to 8 bytes
	// and make it so that if its not done and the size if full
	// expand the buf * 2 and copy the contents to it
	// you could get a request that is larger than 1024
	buf := make([]byte, bufferSize)
	readIndex := 0

	for r.State != StateDone {
		n, err := reader.Read(buf[readIndex:])
		if err != nil && err != io.EOF {
			return nil, err
		}

		readIndex += n

		if err == io.EOF {
			r.State = StateDone
			break
		}

		consumedBytes, err := r.parse(buf[:readIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[consumedBytes:readIndex])
		readIndex -= consumedBytes
	}

	return r, nil
}

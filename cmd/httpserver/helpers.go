package main

import (
	"crypto/sha256"
	"fmt"
	"http-server-miha/internal/headers"
	"http-server-miha/internal/response"
	"net/http"
)

func toString(bytes []byte) string {
	s := ""

	for _, b := range bytes {
		s += fmt.Sprintf("%02x", b)
	}

	return s
}

func chunkedEncoding(w *response.Writer, h headers.Headers, res *http.Response, status response.StatusCode) {
	w.WriteStatusLine(status)
	h.Remove("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	h.Replace("Content-Type", "text/plain")
	w.WriteHeaders(h)

	body := []byte{}

	for {
		buf := make([]byte, 1024)

		n, err := res.Body.Read(buf)
		if err != nil {
			break
		}

		body = append(body, buf[:n]...)
		w.WriteChunkedBody(buf[:n])
	}

	w.WriteChunkedBodyDone()
	trailers := *headers.NewHeaders()
	outputHash := sha256.Sum256(body)
	trailers.Set("X-Content-SHA256", toString(outputHash[:]))
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeaders(trailers)
}

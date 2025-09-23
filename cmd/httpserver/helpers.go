package main

import (
	"http-server-miha/internal/headers"
	"http-server-miha/internal/response"
	"net/http"
)

func chunkedEncoding(w *response.Writer, h headers.Headers, res *http.Response, status response.StatusCode) {
	w.WriteStatusLine(status)
	h.Remove("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Replace("Content-Type", "text/plain")
	w.WriteHeaders(h)

	for {
		buf := make([]byte, 1024)

		n, err := res.Body.Read(buf)
		if err != nil {
			break
		}

		w.WriteChunkedBody(buf[:n])
	}

	// TODO: Should check for errors on both w.WriteChunked functions
	w.WriteChunkedBodyDone()
}

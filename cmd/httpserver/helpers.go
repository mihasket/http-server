package main

import (
	"fmt"
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

		w.WriteBody(fmt.Appendf(nil, "%X\r\n", n))
		w.WriteBody(buf[:n])
		w.WriteBody([]byte("\r\n"))
	}

	w.WriteBody([]byte("0\r\n\r\n"))
}

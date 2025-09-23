package main

import (
	"fmt"
	"http-server-miha/internal/request"
	"http-server-miha/internal/response"
	"http-server-miha/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 4000

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		body := response.Respond200()
		h := response.GetDefaultHeaders(len(body))
		status := response.OK
		target := req.RequestLine.RequestTarget

		if target == "/badrequest" {
			status = response.BadRequest
			body = response.Respond400()

			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		} else if target == "/servererror" {
			status = response.InternalServerError
			body = response.Respond500()

			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		} else if target == "/video" {
			h.Replace("Content-Type", "video/mp4")

			body, err := os.ReadFile("assets/video.mp4")
			if err != nil {
				log.Fatal(err)
			}
			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

			w.WriteStatusLine(status)
			w.WriteHeaders(h)
			w.WriteBody(body)
		} else if strings.HasPrefix(target, "/httpbin/") {
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				status = response.InternalServerError
				body = response.Respond500()

				h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
				w.WriteStatusLine(status)
				w.WriteHeaders(h)
				w.WriteBody(body)
				return
			}

			chunkedEncoding(w, h, res, status)
			return
		}

		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

package main

import (
	"fmt"
	"http-server-miha/internal/request"
	"http-server-miha/internal/response"
	"http-server-miha/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 4000

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		body := response.Respond200()
		h := response.GetDefaultHeaders(len(body))
		status := response.OK

		if req.RequestLine.RequestTarget == "/badrequest" {
			status = response.BadRequest
			body = response.Respond400()

			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		} else if req.RequestLine.RequestTarget == "/servererror" {
			status = response.InternalServerError
			body = response.Respond500()

			h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
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

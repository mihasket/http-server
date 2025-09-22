package main

import (
	"http-server-miha/internal/request"
	"http-server-miha/internal/response"
	"http-server-miha/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 4000

func main() {
	s, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/badrequest" {
			sErr := &server.HandlerError{
				StatusCode: response.BadRequest,
				Message:    "Bad request\n",
			}

			w.Write([]byte(sErr.Message))

			return sErr
		} else if req.RequestLine.RequestTarget == "/servererror" {
			sErr := &server.HandlerError{
				StatusCode: response.InternalServerError,
				Message:    "Internal server error\n",
			}

			w.Write([]byte(sErr.Message))

			return sErr
		}

		return nil
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

package server

import (
	"fmt"
	"http-server-miha/internal/request"
	"http-server-miha/internal/response"
	"log"
	"net"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
	port     int
	closed   bool
}

func Serve(port int, handler Handler) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
		port:     port,
		closed:   false,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.closed = true
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if s.closed {
			break
		}

		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	defer fmt.Println("Connection to", conn.RemoteAddr(), "closed")

	resWriter := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		resWriter.WriteStatusLine(response.BadRequest)
		h := response.GetDefaultHeaders(0)
		resWriter.WriteHeaders(h)

		return
	}

	s.handler(resWriter, req)
}

package server

import (
	"fmt"
	"http-server-miha/internal/response"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
	port     int
	closed   bool
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: ln,
		port:     port,
		closed:   false,
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

	response.WriteStatusLine(conn, response.OK)
	h := response.GetDefaultHeaders(0)
	response.WriteHeaders(conn, h)
}

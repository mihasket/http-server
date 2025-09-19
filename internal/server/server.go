package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Listener net.Listener
	Port     int
}

func Serve(port int) (*Server, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}

	s := &Server{
		Listener: ln,
		Port:     port,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	return s.Listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		go s.handle(conn)

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func (s *Server) handle(conn net.Conn) {
	output := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!")
	conn.Write(output)
	conn.Close()
}

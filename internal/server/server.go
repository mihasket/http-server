package server

import (
	"bytes"
	"fmt"
	"http-server-miha/internal/request"
	"http-server-miha/internal/response"
	"io"
	"log"
	"net"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (he HandlerError) Write(w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return err
	}

	h := response.GetDefaultHeaders(len(he.Message))
	err = response.WriteHeaders(w, h)
	if err != nil {
		return err
	}

	return nil
}

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.BadRequest,
			Message:    err.Error(),
		}

		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)

		conn.Write(buf.Bytes())
		return
	}

	b := buf.Bytes()

	response.WriteStatusLine(conn, response.OK)
	h := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, h)
	conn.Write(b)
}

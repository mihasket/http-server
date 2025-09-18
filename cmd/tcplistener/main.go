package main

import (
	"fmt"
	"http-server-miha/internal/request"
	"log"
	"net"
)

const port = ":4000"

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}

	defer ln.Close()
	fmt.Println("Listening for TCP traffic on", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error reading request from reader: %s\n", err.Error())
		}

		fmt.Println("Request line:")
		fmt.Println("- Method: ", r.RequestLine.Method)
		fmt.Println("- Target: ", r.RequestLine.RequestTarget)
		fmt.Println("- Version: ", r.RequestLine.HttpVersion)

		r.Headers.Output()

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

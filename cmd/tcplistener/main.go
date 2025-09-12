package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":4000"

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		var line string
		buf := make([]byte, 8)

		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			if err == io.EOF {
				break
			}

			if bytes.Contains(buf, []byte("\n")) {
				parts := strings.SplitN(string(buf[:n]), "\n", 2)
				line = line + parts[0]

				out <- line

				line = parts[1]
				continue
			}

			line = line + string(buf[:n])
		}
	}()

	return out
}

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

		lines := getLinesChannel(conn)

		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

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
	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatal("error", "error", err)
		os.Exit(1)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}

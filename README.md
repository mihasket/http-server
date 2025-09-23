# HTTP/1.1

![image](/docs/image.png)

This project is a minimal implementation of an HTTP/1.1 server built directly on top of [go TCP sockets](https://pkg.go.dev/net#TCPConn).

It follows the specifications in [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) (HTTP Semantics) and [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) (HTTP/1.1).

The goal was to learn how HTTP, more specifically HTTP/1.1, works under the hood, without relying on higher-level libraries.
This is not intended for use as it has a lot of edge cases that were not covered.

## Notable features

- **Custom handler function**: the server exposes a handler interface that allows users to define their own logic for responding to requests
- **Concurrent request handling**: each connection is handled in its own goroutine, allowing multiple requests to be served at once
- **HTTPbin proxy**: requests to `/httpbin/{path}` are proxied to [https://httpbin.org/](https://httpbin.org/) and streamed back to the client
  - **Chunked transfer encoding**: the server can stream responses using `Transfer-Encoding: chunked` when the content length is not known ahead of time
  - Responses are streamed with `Transfer-Encoding: chunked`
  - Trailer fields are added after the stream completes:
    - `X-Content-SHA256`: SHA-256 hash of the proxied response body
    - `X-Content-Length`: total length of the proxied response body

## Run

```bash
go run ./cmd/httpserver/
```

Or if you want live reloading, use [air](https://github.com/air-verse/air):
```bash
air
```

## Testing

1. Status 200
```bash
curl -v http://localhost:4000/
```

Expected output:
```
HTTP/1.1 200 OK
content-length: 137
connection: close
content-type: text/html

<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Sucess!</h1>
    <p>You request was sucessful</p>
  </body>
</html>
```

2. Status 400
```bash
curl -v http://localhost:4000/badrequest
```

Expected output:
```
HTTP/1.1 400 Bad Request
connection: close
content-type: text/html
content-length: 154

<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>You request was not sucessful</p>
  </body>
</html>
```

3. Status 500
```bash
curl -v http://localhost:4000/servererror
```

Expected output:
```
HTTP/1.1 500 Internal Server Error
content-length: 162
connection: close
content-type: text/html

<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>We made a mistake</p>
  </body>
</html>
```

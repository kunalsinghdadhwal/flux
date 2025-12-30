# Flux

A minimal HTTP/1.1 server implementation in Go, built from scratch without using `net/http`.

## About

Flux is an educational HTTP server that implements the HTTP/1.1 protocol from the ground up. It handles raw TCP connections, parses HTTP requests byte-by-byte using a state machine, and constructs HTTP responses manually. The goal is to understand how HTTP works at the protocol level.

## Features

### Request Parsing
- State machine-based parser (Init, Headers, Body, Done states)
- Streaming buffer support for partial reads
- Request line parsing (method, target, HTTP version)
- Header field validation following RFC 7230 token rules

### Header Handling
- Case-insensitive header storage
- Get, Set, Replace, and Delete operations
- Multi-value header support (comma-separated)
- Custom header iteration with ForEach

### Response Writing
- Status line generation (200 OK, 400 Bad Request, 500 Internal Server Error)
- Header serialization with proper CRLF formatting
- Body writing with Content-Length or chunked encoding

### Chunked Transfer Encoding
- Streaming response bodies without knowing total size upfront
- Proper chunk formatting (size in hex + CRLF + data + CRLF)
- Trailer header support (X-Content-SHA256, X-Content-Length)
- Final chunk termination (0 + CRLF + trailers + CRLF)

### Raw HTTP Client
- TCP-based HTTP client without `net/http`
- Manual request line and header construction
- Response parsing with status code extraction
- Connection streaming for large responses

### Server
- TCP listener with concurrent connection handling via goroutines
- Custom handler function support
- Graceful shutdown on SIGINT/SIGTERM signals

## Usage

Build the server:

```bash
go build -o bin/httpserver ./cmd/httpserver
```

Run the server:

```bash
./bin/httpserver
```

Server listens on port 42069.

## Endpoints

| Path | Description |
|------|-------------|
| `/` | Returns 200 OK with HTML success page |
| `/yourproblem` | Returns 400 Bad Request |
| `/myproblem` | Returns 500 Internal Server Error |
| `/httpbin/*` | Reverse proxy to httpbin.org with chunked encoding and SHA256 trailer |
| `/video` | Streams video file from assets/vim.mp4 with chunked encoding |

## Examples

Basic request:

```bash
curl http://localhost:42069/
```

Error responses:

```bash
curl http://localhost:42069/yourproblem
curl http://localhost:42069/myproblem
```

Proxy to httpbin:

```bash
curl http://localhost:42069/httpbin/html
curl http://localhost:42069/httpbin/json
curl http://localhost:42069/httpbin/stream/5
```

Video streaming:

```bash
curl http://localhost:42069/video --output video.mp4
```

View chunked encoding with trailers:

```bash
curl -v --raw http://localhost:42069/httpbin/html
```

Run unit tests:

```bash
go test ./...
```

## How It Works

1. Server starts TCP listener on port 42069
2. Each connection spawns a goroutine
3. Request bytes are read into a buffer
4. State machine parses request line, then headers, then body
5. Handler function processes request and writes response
6. Response is serialized: status line + headers + CRLF + body
7. Connection closes (Connection: close header)

## Dependencies

- Go 1.22.2+
- testify v1.11.1 (testing only)
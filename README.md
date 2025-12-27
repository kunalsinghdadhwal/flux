# Flux

A minimal HTTP/1.1 server implementation in Go, built from scratch following RFC 9110 and RFC 9112 specifications.

## Overview

Flux implements core HTTP/1.1 protocol parsing and response handling without relying on Go's standard library HTTP packages. The server uses a state machine-based request parser with streaming buffer support.

## Features

- HTTP/1.1 request parsing with state machine (Init, Headers, Body, Done)
- Header field parsing with case-insensitive storage and RFC 7230 token validation
- Request body parsing with Content-Length support
- Response status line and header writing (200, 400, 500)
- TCP listener with goroutine-based concurrent connection handling
- Graceful shutdown on SIGINT/SIGTERM

## Components

- `internal/request` - Request parser with sliding buffer for partial reads
- `internal/headers` - Header storage with Get/Set/Replace/Parse operations
- `internal/response` - Response writer with status codes and default headers
- `internal/server` - TCP server with connection handling

## Usage

Build and run the server:

```bash
go build -o bin/httpserver ./cmd/httpserver
./bin/httpserver
```

The server listens on port 42069. Test with:

```bash
curl http://localhost:42069/
curl -X POST -d "body" -H "Content-Length: 4" http://localhost:42069/
```

Run tests:

```bash
go test ./...
```

## Dependencies

- Go 1.22.2+
- testify v1.11.1 (testing only)

## Current State

Implemented:
- Request line parsing (method, target, HTTP version)
- Header parsing and validation
- Body parsing with Content-Length
- Basic response writing
- Concurrent connection handling

Not yet implemented:
- Request routing
- Response body content
- HTTP method-specific handlers
- Connection persistence
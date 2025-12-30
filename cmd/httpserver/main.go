package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kunalsinghdadhwal/flux/internal/headers"
	"github.com/kunalsinghdadhwal/flux/internal/request"
	"github.com/kunalsinghdadhwal/flux/internal/response"
	"github.com/kunalsinghdadhwal/flux/internal/server"
)

const port = 42069

func toStr(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}
	return out
}

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respond200() []byte {
	return []byte(
		`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		body := respond200()
		status := response.StatusOK
		if req.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusBadRequest
			body = respond400()
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInernalServerError
			body = respond500()
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := req.RequestLine.RequestTarget

			res, err := response.Get("http://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				status = response.StatusInernalServerError
			} else {
				defer res.Body.Close()
				w.WriteStatusLine(response.StatusOK)
				h.Delete("Content-length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("content-type", "text/plain")

				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				w.WriteHeaders(*h)

				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}

					fullBody = append(fullBody, data[:n]...)
					w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailers := headers.NewHeaders()
				out := sha256.Sum256(fullBody)

				trailers.Set("X-Content-SHA256", toStr(out[:]))
				trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
				w.WriteHeaders(*trailers)
				return
			}
		} else if req.RequestLine.RequestTarget == "/video" {
			f, err := os.Open("assets/vim.mp4")
			if err != nil {
				body = respond500()
				status = response.StatusInernalServerError
			} else {
				defer f.Close()
				h.Delete("Content-length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("content-type", "video/mp4")
				w.WriteStatusLine(response.StatusOK)
				w.WriteHeaders(*h)

				buf := make([]byte, 4096)
				for {
					n, err := f.Read(buf)
					if n > 0 {
						w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
						w.WriteBody(buf[:n])
						w.WriteBody([]byte("\r\n"))
					}
					if err != nil {
						break
					}
				}
				w.WriteBody([]byte("0\r\n\r\n"))
				return
			}
		}
		h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-type", "text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

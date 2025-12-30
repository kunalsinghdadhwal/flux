package response

import (
	"fmt"
	"io"
	"net"
	"net/url"

	"github.com/kunalsinghdadhwal/flux/internal/headers"
)

type Response struct {
	StatusCode int
	Headers    *headers.Headers
	Body       net.Conn
}

func Get(rawURL string) (*Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "https" {
		return nil, fmt.Errorf("HTTPS not supported, use HTTP")
	}

	host := u.Host
	if u.Port() == "" {
		host = fmt.Sprintf("%s:80", u.Hostname())
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	path := u.Path
	if path == "" {
		path = "/"
	}
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	// Write request
	requestLine := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	conn.Write([]byte(requestLine))
	conn.Write([]byte(fmt.Sprintf("Host: %s\r\n", u.Host)))
	conn.Write([]byte("Connection: close\r\n"))
	conn.Write([]byte("User-Agent: flux/1.0\r\n"))
	conn.Write([]byte("\r\n"))

	// Read and parse status line
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Close()
		return nil, err
	}

	data := buf[:n]
	statusCode := 0
	fmt.Sscanf(string(data), "HTTP/1.1 %d", &statusCode)

	h := headers.NewHeaders()

	// Find headers section
	headerStart := 0
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			headerStart = i + 2
			break
		}
	}

	if headerStart > 0 {
		h.Parse(data[headerStart:])
	}

	return &Response{
		StatusCode: statusCode,
		Headers:    h,
		Body:       conn,
	}, nil
}

type StatusCode int

const (
	StatusOK                 StatusCode = 200
	StatusBadRequest         StatusCode = 400
	StatusInernalServerError StatusCode = 500
)

var UNRECOGNIZED_ERROR_CODE = fmt.Errorf("UNRECOGNIZED ERROR CODE")

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}
	switch statusCode {
	case StatusOK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInernalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return UNRECOGNIZED_ERROR_CODE
	}

	_, err := w.writer.Write(statusLine)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var err error = nil
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})

	b = fmt.Append(b, "\r\n")
	_, err = w.writer.Write(b)

	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)

	return n, err
}

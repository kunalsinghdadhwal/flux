package response

import (
	"fmt"
	"io"

	"github.com/kunalsinghdadhwal/flux/internal/headers"
)

type Response struct {
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
	h.Set("Content-Typr", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, h *headers.Headers) error {
	var err error = nil
	b := []byte{}
	h.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})

	b = fmt.Append(b, "\r\n")
	_, err = w.Write(b)

	return err
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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

	_, err := w.Write(statusLine)
	return err
}

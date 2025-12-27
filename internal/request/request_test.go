package request_test

import (
	"io"
	"strings"
	"testing"

	r "github.com/kunalsinghdadhwal/flux/internal/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}

	req, err := r.RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/", req.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", req.RequestLine.HttpVersion)

	reader = &chunkReader{
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	req, err = r.RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/coffee", req.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", req.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = r.RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)
}

func TestParseHeaders(t *testing.T) {
	// Test: Standard Headers
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	req, err := r.RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)
	host, _ := req.Headers.Get("host")
	assert.Equal(t, "localhost:42069", host)
	userAgent, _ := req.Headers.Get("user-agent")
	assert.Equal(t, "curl/7.81.0", userAgent)
	accept, _ := req.Headers.Get("accept")
	assert.Equal(t, "*/*", accept)

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	req, err = r.RequestFromReader(reader)
	require.Error(t, err)
}

func TestParseBody(t *testing.T) {
	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	req, err := r.RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, "hello world!\n", string(req.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	req, err = r.RequestFromReader(reader)
	require.Error(t, err)
}

package server

import (
	"fmt"
	"io"
	"net"

	"github.com/kunalsinghdadhwal/flux/internal/response"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, headers)

}

func runServer(s *Server, l net.Listener) {
	for {
		conn, err := l.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		closed: false,
	}
	go runServer(server, listener)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}

package request

import (
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERRO_BAD_START_LINE = fmt.Errorf("BAD START LINE IN THE REQUEST") 

func RequestFromReader(r io.Reader) (*Request, error) {
	
}
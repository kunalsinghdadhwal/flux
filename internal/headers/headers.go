package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")
var ERROR_MALFORMED_HEADER = fmt.Errorf("MALFORMED HEADER IN REQUEST")
var ERROR_MALFORMED_FIELD_NAME = fmt.Errorf("MALFORMED FIELD NAME IN HEADER")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERROR_MALFORMED_HEADER
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ERROR_MALFORMED_FIELD_NAME
	}

	return string(name), string(value), nil

}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		// Empty Header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(rn)
		h[name] = value
	}
	return read, done, nil
}

package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const (
	crlf = "\r\n"
)

var (
	separator = []byte(":")
)

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIdx := bytes.Index(data, []byte(crlf))
	if crlfIdx == -1 {
		return 0, false, nil
	}

	if crlfIdx == 0 {
		return 2, true, nil
	}

	sepIdx := bytes.Index(data, separator)
	if sepIdx == -1 {
		return 0, true, errors.New("no appropriate separator found")
	}

	if data[sepIdx-1] == ' ' {
		return 0, false, errors.New("whitespace not allowed between field-name and separator")
	}

	key := string(data[:sepIdx])
	key = strings.TrimSpace(key)
	val := string(data[sepIdx+1 : crlfIdx])
	val = strings.TrimSpace(val)

	h[key] = val

	return crlfIdx + 2, false, nil
}

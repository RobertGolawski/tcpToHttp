package headers

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() *Headers {
	return &Headers{}
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

	firstColon := bytes.Index(data, separator)
	if firstColon == -1 {
		return 0, true, errors.New("no appropriate separator found")
	}

	if data[firstColon-1] == ' ' {
		return 0, true, errors.New("whitespace not allowed between field-name and separator")
	}

	key := string(data[:firstColon])
	key = strings.TrimLeftFunc(key, unicode.IsSpace)

	return 0, false, nil
}

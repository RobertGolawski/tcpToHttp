package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

const (
	crlf         = "\r\n"
	allowedChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
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
	if sepIdx == -1 || sepIdx == 0 {
		return 0, false, errors.New("no appropriate separator found")
	}

	if data[sepIdx-1] == ' ' {
		return 0, false, errors.New("whitespace not allowed between field-name and separator")
	}

	key := string(data[:sepIdx])
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)

	if checkForInvalid(key) {
		return 0, false, errors.New("invalid character")
	}

	val := string(data[sepIdx+1 : crlfIdx])
	val = strings.TrimSpace(val)

	if _, ok := h[key]; ok {
		v := h[key]
		newVal := fmt.Sprintf("%s, %s", v, val)
		h[key] = newVal
	} else {
		h[key] = val
	}

	return crlfIdx + 2, false, nil
}

func checkForInvalid(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(allowedChars, r) {
			return true
		}
	}
	return false
}

func (h Headers) Get(key string) string {
	key = strings.ToLower(key)
	return h[key]
}

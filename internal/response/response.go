package response

import (
	"fmt"
	"io"
	"strconv"
	"tcpToHttp/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func getStatusLine(sc StatusCode) []byte {
	reason := ""
	switch sc {
	case StatusOK:
		reason = "OK"
	case StatusBadRequest:
		reason = "Bad Request"
	case StatusInternalServerError:
		reason = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", sc, reason))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = strconv.Itoa(contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	//I actually wrote this initially to be basically the same
	//as the answer with a loop over k, v and writing the headers one by one but then changed it :madge:
	str := ""
	for k, v := range headers {
		str += fmt.Sprintf("%v: %v\r\n", k, v)
	}
	str += "\r\n"
	_, err := w.Write([]byte(str))
	return err
}

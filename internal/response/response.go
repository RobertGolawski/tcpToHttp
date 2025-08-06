package response

import (
	"errors"
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

const (
	WriterState int = iota
	WriterStatusLine
	WriterHeaders
	WriterBody
)

type Writer struct {
	Con         io.Writer
	WriterState int
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.WriterState != WriterStatusLine {
		return errors.New(fmt.Sprintf("trying to write status like while in state: %v", w.WriterState))
	}
	err := WriteStatusLine(w.Con, statusCode)
	w.WriterState = WriterHeaders
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.WriterState != WriterHeaders {
		return errors.New(fmt.Sprintf("trying to write headers in the wrong state: %v", w.WriterState))
	}

	err := WriteHeaders(w.Con, headers)
	w.WriterState = WriterBody
	return err
}

func (w *Writer) WriteBody(b []byte) error {
	if w.WriterState != WriterBody {
		return errors.New(fmt.Sprintf("trying to write body in the wrong state: %v", w.WriterState))
	}
	_, err := w.Con.Write(b)
	w.WriterState = WriterStatusLine
	return err
}

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
	h["Content-Type"] = "text/html"
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

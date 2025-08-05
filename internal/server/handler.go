package server

import (
	"io"
	"tcpToHttp/internal/request"
)

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

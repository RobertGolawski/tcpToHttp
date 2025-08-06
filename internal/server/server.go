package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"tcpToHttp/internal/request"
	"tcpToHttp/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	Closed   atomic.Bool
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{listener: listener, handler: handler}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.Closed.Store(true)
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			if s.Closed.Load() {
				return
			}
			log.Printf("error accepting connection: %v", err)
			continue
		}

		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("error requesting from reader: %v", err)
		return
	}

	writer := response.Writer{Con: conn, WriterState: response.WriterStatusLine}
	s.handler(&writer, req)

	return
}

func WriteHandlerError(w io.Writer, h *HandlerError) {
	err := response.WriteStatusLine(w, response.StatusCode(h.StatusCode))
	if err != nil {
		log.Printf("error sending response: %v", err)
		return
	}

	body := []byte(h.Message)
	headers := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		log.Printf("error sending headers: %v", err)
	}
	w.Write(body)
}

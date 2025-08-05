package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"tcpToHttp/internal/response"
)

var (
	reply = []byte("HTTP/1.1 200 OK\nContent-Type: text/plain\n\nHello World!")
)

type Server struct {
	listener net.Listener
	Closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{listener: listener}

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
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Printf("error writing statsline: %v", err)
	}
	defaultHeaders := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, defaultHeaders)
	if err != nil {
		log.Printf("error writing headers: %v", err)
	}
	return
}

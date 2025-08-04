package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	serverState int = iota
	serverInitialised
	serverListening
	serverClosing
)

var (
	reply = []byte("HTTP/1.1 200 OK\nContent-Type: text/plain\n\nHello World!")
)

type Server struct {
	state    int
	Port     string
	Listener *net.Listener
}

func Serve(port int) (*Server, error) {
	portString := strconv.Itoa(port)
	s := Server{state: serverInitialised, Port: portString}
	go s.listen()

	return &s, nil
}

func (s *Server) Close() error {
	if s.state == serverClosing {
		return errors.New("Cannot close a server that's closing")
	}

	s.state = serverClosing
	return nil
}

func (s *Server) listen() {
	if s.state == serverInitialised {
		s.state = serverListening
	} else {
		log.Println("Cannot listen to a server that's already listening or closed")
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
	if err != nil {
		log.Printf("error when opening a listener: %v", err)
		return
	}

	defer listener.Close()
	s.Listener = &listener

	for s.state == serverListening {
		con, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}

		go s.handle(con)
	}
}

func (s *Server) handle(conn net.Conn) {
	fmt.Printf("%v", string(reply))
	conn.Write(reply)
	conn.Close()
}

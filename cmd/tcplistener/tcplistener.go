package main

import (
	"fmt"
	"log"
	"net"
	"tcpToHttp/internal/request"
)

func main() {

	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatalf("error making listener: %v", err)
	}

	defer listener.Close()

	for {
		con, err := listener.Accept()
		if err != nil {
			continue
		}

		log.Println("Connection has been accepted")

		req, err := request.RequestFromReader(con)
		if err != nil {
			log.Fatalf("error reading request: %v", err)
		}
		fmt.Printf(`Request line:
- Method: %s
- Target: %s
- Version: %s
`, req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HTTPVersion)

		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Printf("%v", string(req.Body))
		con.Close()
		log.Printf("Connection closed")
	}
}

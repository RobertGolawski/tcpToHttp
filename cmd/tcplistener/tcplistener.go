package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

		myChan := getLinesChannel(con)
		for line := range myChan {
			fmt.Printf("%v\n", line)
		}
		con.Close()
		log.Printf("Connection closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	go func(io.ReadCloser, chan string) {
		data := make([]byte, 200)
		var n int
		var err error
		for {
			n, err = f.Read(data)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("Encountered unexpected error: %v", err)
					close(strChan)
					return
				} else {
					close(strChan)
					return
				}
			}
			str := string(data[:n])
			strChan <- str
		}
	}(f, strChan)

	return strChan

}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	resAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("Error when resolving address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, resAddr)
	if err != nil {
		log.Fatalf("Error when dialling: %v", err)
	}

	defer conn.Close()

	bioReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := bioReader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading line: %v", err)
			continue
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("Error writing to connection: %v", err)
		}
	}
}

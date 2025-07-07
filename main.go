package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// fmt.Println("I hope I get the job!")

	data := make([]byte, 8)

	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}

	for err == nil {
		_, err = file.Read(data)
		fmt.Printf("read: %s\n", string(data))
	}
	if err == io.EOF {
		os.Exit(0)
	} else {
		fmt.Printf("Encountered different error: %v", err)
	}
}

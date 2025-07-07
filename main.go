package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	data := make([]byte, 8)

	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}
	defer file.Close()

	var n int
	for {
		n, err = file.Read(data)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Encountered different error: %v", err)
				return
			} else {
				os.Exit(0)
			}
		}
		fmt.Printf("read: %s\n", string(data[:n]))
	}
}

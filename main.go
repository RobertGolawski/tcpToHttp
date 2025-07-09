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

	var currLine string
	var n int
	for {
		n, err = file.Read(data)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("Encountered different error: %v", err)
				return
			} else {
				break
			}
		}

		str := string(data[:n])

		for _, c := range str {
			currLine += string(c)
			if c == '\n' {
				fmt.Printf("read: %v", currLine)
				currLine = ""
			}

		}
	}
}

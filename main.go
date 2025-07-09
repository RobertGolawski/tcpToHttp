package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {

	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
	}
	defer file.Close()

	myChan := getLinesChannel(file)

	for {
		select {
		case line, ok := <-myChan:
			if !ok {
				return
			}
			fmt.Printf("read: %v", line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	strChan := make(chan string)

	go func(io.ReadCloser, chan string) {
		defer f.Close()
		data := make([]byte, 8)
		var currLine string
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

			for _, c := range str {
				currLine += string(c)
				if c == '\n' {
					strChan <- currLine
					currLine = ""
				}
			}
		}
	}(f, strChan)

	return strChan

}

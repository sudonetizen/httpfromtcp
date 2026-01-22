package main

import (
	"io"
	"os"
	"fmt"
	"log"
	//"strings"
)

const filePath = "messages.txt"

func main() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("os.Open got err: %v for file: %v", err, filePath)	
	}

	fmt.Printf("Reading data from %v\n\n", filePath)

	linesFromChan := getLinesChannel(file)
	for line := range  linesFromChan {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		
		line := ""
		buffer := make([]byte, 8) 

		for {
			n, err := f.Read(buffer)
			if err == io.EOF { break }

			if err != nil {
				fmt.Printf("file.Read got err: %v\n", err.Error())	
				break
			}

			read_bytes := string(buffer[:n])

			for _, char := range read_bytes {
				if char != '\n' {
					line += string(char)
				} else {
					lines <- line
					line = ""
				}
			}

		}
		
	}()
	return lines

}

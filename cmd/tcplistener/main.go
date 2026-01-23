package main

import (
	"io"
	"fmt"
	"log"
	"net"
)

const filePath = "messages.txt"
const listenPort = ":42069"

func main() {
	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatalf("net.Listen got err: %v at port: %v", err, listenPort)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()	
		if err != nil {
			log.Fatalf("l.Accept got err: %v", err)
		}
		fmt.Println("accepted connection from", conn.RemoteAddr())

		linesFromChan := getLinesChannel(conn)
		for line := range  linesFromChan {
			fmt.Printf("read: %s\n", line)
		}

		fmt.Println("connection to", conn.RemoteAddr(), "closed")
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
		
		if line != "" { lines <- line }
		
	}()
	return lines

}

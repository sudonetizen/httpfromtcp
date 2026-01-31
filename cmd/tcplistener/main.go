package main

import (
	"fmt"
	"log"
	"net"

	"httpfromtcp/internal/request"
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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("l.Accept got err: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n", req.RequestLine.Method)
		fmt.Printf("- Target: %v\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))

	}

}


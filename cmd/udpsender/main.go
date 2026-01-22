package main

import (
	"io"
	"os"
	"fmt"
	"net"
	"log"
	"bufio"
)

const serverAddr = "localhost:42069"

func main () {
	// udp end point, remote addr
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("net.ResolveUDPAddr got err: %v at port: %v", err, serverAddr)
	}

	// creating udp conn
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("net.DialUDP got err: %v", err)
	}
	defer udpConn.Close()

	fmt.Printf("sending to %v\ntype message and press enter\n\n", serverAddr)

	// reading from os.stdin 
	r := bufio.NewReader(os.Stdin)

	// loop
	for {
		fmt.Print("> ")

		lineString, err := r.ReadString('\n')
		if err == io.EOF { break }
		
		if err != nil {
			fmt.Printf("r.ReadString got err: %v", err.Error())
			break
		}
		
		n, err := udpConn.Write([]byte(lineString))
		if err != nil {
			fmt.Printf("udpConn.Write got err: %v", err.Error())
			break
		}

		fmt.Printf("sent %v bytes to %v\n", n, udpAddr.String())

	}

}

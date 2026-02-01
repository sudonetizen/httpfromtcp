package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	Listener net.Listener
	closed atomic.Bool
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("net.Listen got err: %v\n", err.Error())
	}

	srv := &Server{Listener: l}

	go srv.listen()

	return srv, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.Listener != nil {return s.Listener.Close() }
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() { return }
			log.Printf("error: acceptiong connection: %v\n", err.Error())
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {	
	defer conn.Close()
	message := "HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/plain\r\n" +
				"Content-Length: 13\r\n" +
				"\r\n" +
				"Hello World!\n"

	n, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("connection write got err: %v\n", err.Error())
		return
	}
	fmt.Printf("sent %v bytes: ->%v<-\n", n, message)
}

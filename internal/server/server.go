package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"httpfromtcp/internal/response"
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

	err := response.WriteStatusLine(conn, response.StatusOk)
	if err != nil {
		fmt.Printf("response.WriteStatusLine got err: %v\n", err.Error())
		return
	}

	h := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, h)
	if err != nil {
		fmt.Printf("response.WriteHeaders got err: %v\n", err.Error())
		return
	}

}

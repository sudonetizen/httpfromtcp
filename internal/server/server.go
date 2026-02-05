package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request) 

type Server struct {
	Listener net.Listener
	closed atomic.Bool
	hfunc Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("net.Listen got err: %v\n", err.Error())
	}

	srv := &Server{Listener: l, hfunc: handlerFunc}

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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("request.RequestFromReader got err: %v\n", err.Error())
		return
	}

	writeStruct := &response.Writer{Writer: conn}	
	s.hfunc(writeStruct, req)

}

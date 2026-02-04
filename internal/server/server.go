package server

import (
	"io"
	"fmt"
	"log"
	"net"
	"bytes"
	"sync/atomic"

	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type Handler func(w io.Writer, req *request.Request) *HandlerError

type Server struct {
	Listener net.Listener
	closed atomic.Bool
	hfunc Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

func writeHandlerError(he *HandlerError, w io.Writer) error {
	err := response.WriteStatusLine(w, he.StatusCode)
	if err != nil {
		return fmt.Errorf("response.WriteStatusLine got err: %v\n", err.Error())
	}

	h := response.GetDefaultHeaders(len([]byte(he.Message)))

	err = response.WriteHeaders(w, h)
	if err != nil {
		return fmt.Errorf("response.WriteHeaders got err: %v\n", err.Error())	
	}

	_, err = w.Write([]byte(he.Message))
	if err != nil {
		return fmt.Errorf("w.Write got err: %v\n", err.Error())
	}

	return nil
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
		err := writeHandlerError(&HandlerError{StatusCode: response.StatusBadReq, Message: err.Error()}, conn)
		if err != nil {
			fmt.Printf("writeHandlerError got err: %v\n", err.Error())
			return
		}
		return
	}

	b := bytes.NewBuffer([]byte{})

	handlerError := s.hfunc(b, req)
	if handlerError != nil {
		err := writeHandlerError(handlerError, conn)

		if err != nil {
			fmt.Printf("writeHandlerError got err: %v\n", err.Error())
			return
		}

		return
	}

	err = response.WriteStatusLine(conn, response.StatusOk)
	if err != nil {
		fmt.Printf("response.WriteStatusLine got err: %v\n", err.Error())
		return
	}

	h := response.GetDefaultHeaders(len(b.Bytes()))

	err = response.WriteHeaders(conn, h)
	if err != nil {
		fmt.Printf("response.WriteHeaders got err: %v\n", err.Error())
		return
	}

	_, err = conn.Write(b.Bytes())
	if err != nil {
		fmt.Printf("b conn.Write got err: %v\n", err.Error())
		return
	}

}

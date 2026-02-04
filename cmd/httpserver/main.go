package main

import (
	"io"
	"os"
	"log"
	"syscall"
	"os/signal"

	"httpfromtcp/internal/server"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

const port = 42069

func Handler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{StatusCode: response.StatusBadReq, Message: "Your problem is not my problem\n"}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{StatusCode: response.StatusSrvErr, Message: "Woopsie, my bad\n"}
	}

	w.Write([]byte("All good, frfr\n"))	

	return nil
}

func main() {
	srv, err := server.Serve(port, Handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

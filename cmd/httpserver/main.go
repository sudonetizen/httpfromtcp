package main

import (
	"os"
	"log"
	"syscall"
	"os/signal"

	"httpfromtcp/internal/server"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

const port = 42069

func Handler(w *response.Writer, req *request.Request) {
	badResp := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
	srvErrResp := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
	okResp := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"

	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadReq)
		h := response.GetDefaultHeaders(len([]byte(badResp)))
		h.OverWrite("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(badResp))
		return
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		w.WriteStatusLine(response.StatusSrvErr)
		h := response.GetDefaultHeaders(len([]byte(srvErrResp)))
		h.OverWrite("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody([]byte(srvErrResp))
		return
	}

	w.WriteStatusLine(response.StatusOk)
	h := response.GetDefaultHeaders(len([]byte(okResp)))
	h.OverWrite("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(okResp))

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

package response

import (
	"io"
	"fmt"

	"httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusOk     StatusCode = 200
	StatusBadReq StatusCode = 400
	StatusSrvErr StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""

	switch statusCode {
	case StatusOk: 
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusBadReq:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case StatusSrvErr:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		statusLine = fmt.Sprintf("HTTP/1.1 %v \r\n", int(statusCode))
	}

	_, err := w.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("w.Write got err: %v\n", err.Error())
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%v", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for v, k := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%v: %v\r\n", v, k)))

		if err != nil {
			return fmt.Errorf("w.Write got err: %v\n", err.Error())
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("w.Write got err: %v\n", err.Error())
	}

	return nil
}

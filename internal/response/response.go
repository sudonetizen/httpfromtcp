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

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	Writer io.Writer
	writerState writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w, writerState: writerStateStatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("not writerStateStatusLine\n")
	}
	defer func() { w.writerState = writerStateHeaders }()

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

	_, err := w.Writer.Write([]byte(statusLine))
	if err != nil {
		return fmt.Errorf("w.Writer.Write got err: %v\n", err.Error())
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

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("not writerStateStatusLine\n")
	}
	defer func() { w.writerState = writerStateBody }()

	for v, k := range headers {
		_, err := w.Writer.Write([]byte(fmt.Sprintf("%v: %v\r\n", v, k)))

		if err != nil {
			return fmt.Errorf("w.Writer.Write got err: %v\n", err.Error())
		}
	}

	_, err := w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return fmt.Errorf("w.Writer.Write got err: %v\n", err.Error())
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("not writerStateStatusLine\n")
	}

	n, err := w.Writer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("w.Writer.Write got err: %v\n", err.Error())
	}

	return n, nil
}

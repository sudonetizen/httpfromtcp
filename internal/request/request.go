package request

import (
	"io"
	"fmt"
	"bytes"
	"strings"

	"httpfromtcp/internal/headers"
)

const bufferSize = 8

type rState int

const (
	initialized rState = iota
	done
	requestStateParsingHeaders
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	state rState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	request := &Request{state: initialized, Headers: headers.Headers{}}

	for request.state != done {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer) * 2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
		
		nBytesRead, err := reader.Read(buffer[readToIndex:])
		if err == io.EOF {
			if request.state != done {
				return nil, fmt.Errorf("request not done yet, state: %v, read: %v bytes, EOF\n", request.state, nBytesRead)
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reader.Read got err: %v\n", err.Error())
		}

		readToIndex += nBytesRead
		nBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("request.Parse got err: %v\n", err.Error())
		}

		copy(buffer, buffer[nBytesParsed:])
		readToIndex -= nBytesParsed

	}
	
	return request, nil

}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	httpMethods := "GET POST PUT DELETE"

	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, nil, nil
	}

	requestLine := string(data[:idx])
	requestLineParts := strings.Split(requestLine, " ")

	// check: request line consists of three parts: method, path, http version
	if len(requestLineParts) != 3 {
		return 0, nil, fmt.Errorf("invalid request line: %v\n", requestLineParts) 
	}

	requestLineMethod := requestLineParts[0]
	requestLineTarget := requestLineParts[1]
	requestLineHttpVersion := requestLineParts[2]

	// check: HTTP method 
	if !strings.Contains(httpMethods, requestLineMethod) {
		return 0, nil, fmt.Errorf("invalid method: %v\n", requestLineMethod)	
	}

	// check: HTTP version, only 1.1 
	httpName := strings.Split(requestLineHttpVersion, "/")[0]
	httpVersion := strings.Split(requestLineHttpVersion, "/")[1]
	if  httpName != "HTTP" || httpVersion != "1.1" {
		return 0, nil, fmt.Errorf("unsupported: %v\n", requestLineHttpVersion)	
	}
	
	return idx+2, &RequestLine{httpVersion, requestLineTarget, requestLineMethod}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	tBytesParsed := 0
	for r.state != done {
		n, err := r.parseSingle(data[tBytesParsed:])
		if err != nil {
			return 0, err
		}
		tBytesParsed += n
		if n == 0 {
			break
		}
	}
	return tBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		n, rl, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("parseRequestLine got err: %v\n", err.Error())
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *rl
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, dne, err := r.Headers.Parse(data)
		if err != nil {
			return 0, fmt.Errorf("r.Headers.Parse got err: %v\n", err.Error())
		}
		if  dne == true {
			r.state = done
		}
		return n, nil
	case done:
		return 0, fmt.Errorf("error: trying to read data in a done state\n")
	default:
		return 0, fmt.Errorf("error: unknown state\n")
		
	}

}

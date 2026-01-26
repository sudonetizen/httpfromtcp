package request

import (
	"io"
	"fmt"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	b, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("io.ReadAll got err: %v\n", err.Error())
		return nil, err
	}

	requestLine, err := parseRequestLine(b)
	if err != nil {
		fmt.Printf("parseRequestLine got err: %v", err.Error())
		return nil, err
	}

	return &Request{RequestLine: *requestLine}, nil

}

func parseRequestLine(data []byte) (*RequestLine, error) {
	httpMethods := "GET POST PUT DELETE"

	requestStr := string(data)
	requestParts := strings.Split(requestStr, "\r\n")
	requestLine := requestParts[0]
	requestLineParts := strings.Split(requestLine, " ")

	// check: request line consists of three parts: method, path, http version
	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("invalid request line: %v\n", requestLineParts) 
	}

	requestLineMethod := requestLineParts[0]
	requestLineTarget := requestLineParts[1]
	requestLineHttpVersion := requestLineParts[2]

	// check: HTTP method 
	if !strings.Contains(httpMethods, requestLineMethod) {
		return nil, fmt.Errorf("invalid method: %v\n", requestLineMethod)	
	}

	// check: HTTP version, only 1.1 
	if requestLineHttpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid : %v\n", requestLineHttpVersion)	
	}
	
	httpVersion := strings.Split(requestLineHttpVersion, "/")[1]
	return &RequestLine{httpVersion, requestLineTarget, requestLineMethod}, nil
}

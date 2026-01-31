package headers

import (
	"fmt"
	"bytes"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Set(key, value string) {
	v, ok := h[key]
	if ok {
		h[key] = v + ", " + value
	} else {
		h[key] = value	
	}
	
}


func (h Headers) Get(key string) (string) {
	v, _ := h[strings.ToLower(key)]
	return v
}

func (h Headers) Okay(key string) bool {
	_, ok := h[strings.ToLower(key)]
	return ok
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {	
	idx := bytes.Index(data, []byte("\r\n"))

	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	fieldLine := string(data[:idx])
	fieldLine = strings.TrimSpace(fieldLine)
	fieldLineParts := strings.SplitN(fieldLine, ":", 2)

	if len(fieldLineParts) != 2 {
		return 0, false, fmt.Errorf("invalid field line: %v\n", fieldLineParts)
	}

	hKey := fieldLineParts[0]

	if strings.Contains(hKey, " ") || len(hKey) < 1 {
		return 0, false, fmt.Errorf("invalid key: %v\n", hKey)
	}

	for _, x := range hKey {
		if  (x >= '0' && x <= '9') || 
			(x >= 'A' && x <= 'Z') || 
			(x >= 'a' && x <= 'z') || 
			x == '!' || x == '#' || x == '$' || x == '%' || x == '&'  || 
			x == '\'' || x == '*' || x == '+' || x == '-' || x == '.' || 
			x == '^' || x == '_' || x == '`' || x == '/' || x == '~' {
			continue
		} else {
			return 0, false, fmt.Errorf("invalid key: %v\n", hKey)
		}
	}

	hKey = strings.ToLower(hKey)

	hVal := fieldLineParts[1]
	hVal = strings.TrimSpace(hVal)

	h.Set(hKey, hVal) 

	return idx+2, false, nil
}

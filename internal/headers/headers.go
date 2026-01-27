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
	fieldLineParts := strings.Split(fieldLine, " ")

	if len(fieldLineParts) != 2 {
		return 0, false, fmt.Errorf("invalid field line: %v\n", fieldLineParts)
	}

	hKey := fieldLineParts[0]
	hKey = hKey[:len(hKey)-1]
	hVal := fieldLineParts[1]

	h[hKey] = hVal

	return idx+2, false, nil
}

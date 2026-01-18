package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if bytes.HasPrefix(data, []byte("\r\n")) {
		return 2, true, nil
	}
	consumed := bytes.Index(data, []byte("\r\n"))
	if consumed == -1 {
		return 0, false, nil
	}
	line := strings.TrimSpace(string(data[:consumed]))
	headerParts := strings.SplitN(line, ": ", 2)
	if len(headerParts) != 2 || strings.HasSuffix(headerParts[0], " ") {
		return 0, false, errors.New("error: Invalid format header")
	}
	headerName := strings.ToLower(headerParts[0])
	headerVal := headerParts[1]
	matched, err := regexp.Match(`^[A-Za-z0-9!#$%&'*+\-.\^_|~]+$`, []byte(headerName))
	if err != nil {
		return 0, false, err
	}
	if !matched {
		return 0, false, errors.New("error: Header name contains not allowed symbols")
	}
	h.Add(headerName, headerVal)
	return consumed + 2, false, nil
}

func (h Headers) Add(key, val string) {
	if _, ok := h[key]; ok == true {
		h[key] += ", " + val
	} else {
		h[key] = val
	}
}

func (h Headers) Set(key, val string) {
	h[key] = val
}

func (h Headers) Get(key string) string {
	lowerKey := strings.ToLower(key)
	return h[lowerKey]
}

func (h Headers) Del(key string) {
	delete(h, key)
}

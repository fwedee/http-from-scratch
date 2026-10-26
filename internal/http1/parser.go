package http1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error not a string: %w", err)
	}

	// Split Statusline, Header and Body
	statusline := strings.Split(strings.TrimSpace(rawResponse), " ")
	if len(statusline) < 3 {
		return nil, errors.New("statusline wrong")
	}
	method, url := statusline[0], statusline[1]

	// Debuggin Purpose
	log.Println(method)
	log.Println(url)

	// Getting Content-Length
	var contentLength int64
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		trimmed := strings.TrimRight(headerLine, "\r\n")
		if trimmed == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "content-length:") {
			val := strings.TrimSpace(trimmed[len("content-length:"):])
			contentLength, _ = strconv.ParseInt(val, 10, 64)
		}
	}

	var body io.Reader
	if contentLength > 0 {
		body = io.LimitReader(reader, int64(contentLength))
	}

	return http.NewRequest(method, url, body)
}

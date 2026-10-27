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
	"unicode"
	"unicode/utf8"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	// Variables for request assignments
	httpHeader := make(map[string][]string)

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error not a string: %w", err)
	}

	// Split Statusline, httpHeader and Body
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

		headerLines := strings.SplitN(trimmed, ":", 2)
		if len(headerLines) < 2 {
			return nil, errors.New("Either key or value is missing in HTTP Header")
		}

		fieldName := headerLines[0]

		// check if is ascii
		if !IsASCII(fieldName) {
			return nil, errors.New("Headerkey needs to be ascii")
		}

		// check if is empty
		if fieldName == "" {
			return nil, errors.New("Headerkey is empty")
		}

		normalizedFieldName := normalizeHeader(fieldName)

		value := strings.TrimSpace(headerLines[1])

		httpHeader[normalizedFieldName] = append(httpHeader[fieldName], value)

	}

	var body io.Reader
	if contentLength > 0 {
		body = io.LimitReader(reader, int64(contentLength))
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header = httpHeader

	return request, err
}

func splitAndTrimSlice(inputs []string) []string {
	var result []string

	for _, item := range inputs {
		for part := range strings.SplitSeq(item, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}

	return result
}

func normalizeHeader(input string) string {
	lowerCase := strings.ToLower(input)
	parts := strings.Split(lowerCase, "-")
	for i, p := range parts {
		parts[i] = uppercaseFirst(p)
	}
	return strings.Join(parts, "-")
}

func uppercaseFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError || size == 0 {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

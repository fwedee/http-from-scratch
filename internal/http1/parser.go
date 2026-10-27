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

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error not a string: %w", err)
	}

	// Split Statusline
	statusline := strings.Split(strings.TrimSpace(rawResponse), " ")
	if len(statusline) < 3 {
		return nil, errors.New("statusline wrong")
	}
	method, url, httpVersion := statusline[0], statusline[1], statusline[2]

	// Parsing HTTPVersion
	validVersion, err := checkHTTPVersion(httpVersion)
	if err != nil {
		return nil, err
	}

	// Debuggin Purpose
	log.Println(method)
	log.Println(httpVersion)
	log.Println(validVersion)
	log.Println(url)

	// Parsing httpHeader and getting Content-Length of body
	httpHeader, contentLength, err := parsingHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("Header parsing went wrong: %w", err)
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

func checkHTTPVersion(version string) (string, error) {
	normversion := strings.TrimSpace(version)
	versionSplit := strings.Split(normversion, "/")
	if len(versionSplit) != 2 {
		return "", errors.New("Either the protocoll or the version is missing")
	}
	protocol, protocollversion := versionSplit[0], versionSplit[1]
	if protocol != "HTTP" {
		return "", errors.New(protocol + " is not supported, just HTTP")
	}

	switch protocollversion {
	case "1.1":
		return normversion, nil
	case "1.0":
		return "", errors.New(protocollversion + " is not supported as HTTP Protocoll version just HTTP/1.1")
	case "2.0":
		return "", errors.New(protocollversion + " is not supported as HTTP Protocoll version just HTTP/1.1")
	default:
		return "", errors.New(protocollversion + " is not supported as HTTP Protocoll version just HTTP/1.1")
	}
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

func parsingHeader(reader *bufio.Reader) (map[string][]string, int64, error) {
	httpHeader := make(map[string][]string)
	var contentLength int64
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, 0, err
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
			return nil, 0, errors.New("Either key or value is missing in HTTP Header")
		}

		fieldName := headerLines[0]

		// check if is ascii
		if !validHeader(fieldName) {
			return nil, 0, errors.New("Header is not valid")
		}

		normalizedFieldName := normalizeHeader(fieldName)

		value := strings.TrimSpace(headerLines[1])

		httpHeader[normalizedFieldName] = append(httpHeader[normalizedFieldName], value)

	}
	return httpHeader, contentLength, nil
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

func validHeader(s string) bool {
	if s == "" {
		return false
	}
	for i := range s {
		l := s[i]

		isLetter := (l >= 'a' && l <= 'z') || (l >= 'A' && l <= 'Z')
		isDigit := l >= '0' && l <= '9'
		isPunctuation := strings.IndexByte("!#$%&'*+-.^_`|~", l) >= 0

		if !isLetter && !isDigit && !isPunctuation {
			return false
		}
	}

	return true
}

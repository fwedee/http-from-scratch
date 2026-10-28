package http1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			if rawResponse == "" {
				return nil, io.EOF
			}
			return nil, &ParseError{
				StatusCode: http.StatusBadRequest,
				Operation:  "read request line",
				Err:        io.ErrUnexpectedEOF,
			}
		}
		return nil, fmt.Errorf("Error not a string: %w", err)
	}

	// Split Statusline
	statusline := strings.Split(strings.TrimSpace(rawResponse), " ")
	if len(statusline) < 3 {
		return nil, errors.New("statusline wrong")
	}
	method, url, httpVersion := statusline[0], statusline[1], statusline[2]

	// Parsing HTTPVersion
	var validHTTPVersion = regexp.MustCompile(`^HTTP/[0-9]\.[0-9]$`)
	validVersion, err := checkHTTPVersion(httpVersion, validHTTPVersion)
	if err != nil {
		return nil, err
	}

	// Parsing HTTP-Method
	supporedHTTPMethods := map[string]struct{}{
		"GET": {}, "POST": {}, "HEAD": {}, "PUT": {}, "DELETE": {}, "TRACE": {}, "CONNECT": {},
	}
	if _, ok := supporedHTTPMethods[strings.ToUpper(method)]; !ok {
		return nil, &ParseError{
			StatusCode: http.StatusMethodNotAllowed,
			Operation:  "check request HTTP method",
			Err:        fmt.Errorf("unsupported method %q", method),
		}
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
	if request.Body == nil {
		request.Body = http.NoBody
	}
	request.Header = httpHeader
	request.Proto = validVersion

	return request, err
}

func checkHTTPVersion(version string, validHTTPVersion *regexp.Regexp) (string, error) {

	if !validHTTPVersion.MatchString(version) {
		return "", &ParseError{
			StatusCode: http.StatusBadRequest,
			Operation:  "parse HTTP version",
			Err:        fmt.Errorf("malformed version %q", version),
		}
	}
	supportedVersion := map[string]struct{}{
		"HTTP/1.1": {},
		"HTTP/1.0": {},
	}

	if _, ok := supportedVersion[version]; !ok {
		return "", &ParseError{
			StatusCode: http.StatusHTTPVersionNotSupported,
			Operation:  "check HTTP version support",
			Err:        fmt.Errorf("unsupported version %q", version),
		}
	}
	return version, nil
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
			if errors.Is(err, io.EOF) {
				return nil, 0, &ParseError{
					StatusCode: http.StatusBadRequest,
					Operation:  "read request header",
					Err:        io.ErrUnexpectedEOF,
				}
			}
			return nil, 0, fmt.Errorf("read request header: %w", err)
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

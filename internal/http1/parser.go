package http1

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxRequestLineBytes = 8 * 2056
const maxHeaderLineBytes = 8 * 2056
const maxHeaderAmount = 100

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

	if !strings.HasSuffix(rawResponse, "\r\n") {
		return nil, &ParseError{
			StatusCode: 400,
			Operation:  "check suffix",
			Err:        fmt.Errorf("suffix is wrong and not rn: %s", rawResponse),
		}
	}
	// TODO: should check before buffering
	if len(rawResponse) > maxRequestLineBytes {
		return nil, &ParseError{
			StatusCode: 400,
			Operation:  "check request size",
			Err:        fmt.Errorf("the request size is too large %d", len(rawResponse)),
		}
	}
	// Split Statusline
	statusline := strings.Split(strings.TrimSuffix(rawResponse, "\r\n"), " ")
	if len(statusline) < 3 {
		return nil, &ParseError{
			StatusCode: 400,
			Operation:  "check statusline length",
			Err:        fmt.Errorf("statusline wrong %v", statusline),
		}
	}
	method, url, httpVersion := statusline[0], statusline[1], statusline[2]

	if !validToken(method) {
		return nil, &ParseError{
			StatusCode: 400,
			Operation:  "check statusline for invalidCharacters",
			Err:        fmt.Errorf("unsupported character %s", method),
		}
	}
	// Parsing HTTPVersion
	var validHTTPVersion = regexp.MustCompile(`^HTTP/[0-9]\.[0-9]$`)
	validVersion, err := checkHTTPVersion(httpVersion, validHTTPVersion)
	if err != nil {
		return nil, err
	}

	// Parsing HTTP-Method
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
	header := http.Header(httpHeader)

	if hasHeaderToken(header.Values("Transfer-Encoding"), "chunked") {
		data, err := readChunkedBody(reader)
		if err != nil {
			return nil, &ParseError{
				StatusCode: 400,
				Operation:  "read chunked body",
				Err:        err,
			}
		}
		body = bytes.NewReader(data)
	} else if contentLength > 0 {
		buf := make([]byte, int(contentLength))
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			return nil, &ParseError{
				StatusCode: 400,
				Operation:  "checking body lenght of request",
				Err:        fmt.Errorf("the request body size and header don't match"),
			}
		}

		body = bytes.NewReader(buf)
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

	if _, ok := supportedHTTPVersion[version]; !ok {
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

func parsingHeader(reader *bufio.Reader) (map[string][]string, int, error) {
	httpHeader := make(map[string][]string)
	var contentLength int
	var seenContentLength bool
	var headAmount int
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

		if (len(headerLine)) >= maxHeaderLineBytes {
			return nil, 0, &ParseError{
				StatusCode: 431,
				Operation:  "check header length",
				Err:        fmt.Errorf("The header is too long %d", len(headerLine)),
			}
		}

		trimmed := strings.TrimSuffix(headerLine, "\r\n")
		if trimmed == "" {
			break
		}

		headAmount += 1
		if headAmount > maxHeaderAmount {
			return nil, 0, &ParseError{
				StatusCode: 431,
				Operation:  "check header amount",
				Err:        fmt.Errorf("Too many header in request %d", headAmount),
			}
		}

		if strings.HasPrefix(strings.ToLower(trimmed), "content-length:") {
			val := strings.TrimSpace(trimmed[len("content-length:"):])
			parsed, err := strconv.ParseUint(val, 10, strconv.IntSize-1)
			length := int(parsed)

			if err != nil || (seenContentLength && parsed != uint64(length)) {
				return nil, 0, &ParseError{
					StatusCode: 400,
					Operation:  "validate Content-Length",
					Err:        fmt.Errorf("invalide or conflicting content-length: %q", val),
				}

			}
			contentLength = length
			seenContentLength = true
		}

		headerLines := strings.SplitN(trimmed, ":", 2)
		if len(headerLines) < 2 {
			return nil, 0, errors.New("Either key or value is missing in HTTP Header")
		}

		fieldName := headerLines[0]

		// check if the header includes only validToken
		if !validToken(fieldName) {
			return nil, 0, &ParseError{
				StatusCode: 400,
				Operation:  "check header for invalidTokens",
				Err:        fmt.Errorf("Header is not valid: %s", fieldName),
			}
		}

		normalizedFieldName := normalizeHeader(fieldName)

		if !validFieldValue(headerLines[1]) {
			return nil, 0, &ParseError{
				StatusCode: 400,
				Operation:  "check if fieldvalue is valid",
				Err:        fmt.Errorf("Fieldvalue is not valid %s", headerLines[1]),
			}
		}
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

func validToken(s string) bool {
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

func validFieldValue(s string) bool {
	for _, r := range s {
		if r < 0x20 && r != '\t' || r == 0x7f {
			return false
		}
	}
	return true
}

func readChunkedBody(reader *bufio.Reader) ([]byte, error) {
	const maxBodyBytes = 1024 * 1024 //1 MiB

	decoded := httputil.NewChunkedReader(reader)
	body, err := io.ReadAll(io.LimitReader(decoded, maxBodyBytes+1))
	if err != nil {
		return nil, fmt.Errorf("chunked body too large")
	}

	_, _, err = parsingHeader(reader)
	if err != nil {
		return nil, err
	}

	if len(body) > maxBodyBytes {
		return nil, fmt.Errorf("chunked body too large")
	}

	return body, nil
}

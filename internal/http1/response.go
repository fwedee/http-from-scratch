package http1

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type Response struct {
	HTTPVersion    string
	StatusCode     int
	ReasonPhrase   string
	Resource       string
	connection     net.Conn
	ResponseHeader http.Header
	wroteHeader    bool
	body           bytes.Buffer
	requestMethod  string
}

func NewResponse(conn net.Conn) *Response {
	return &Response{connection: conn, ResponseHeader: make(http.Header), StatusCode: http.StatusOK}
}

func (r *Response) Header() http.Header {
	return r.ResponseHeader
}

func (r *Response) Write(data []byte) (n int, err error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.body.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.StatusCode = statusCode
}

func (r *Response) Send() error {

	statusText := http.StatusText(r.StatusCode)
	_, err := fmt.Fprintf(r.connection, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText)
	if err != nil {
		return err
	}

	r.ResponseHeader["Content-Length"] = []string{
		strconv.Itoa(r.body.Len()),
	}

	err = r.ResponseHeader.Write(r.connection)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(r.connection, "\r\n")
	if err != nil {
		return err
	}

	if r.requestMethod == http.MethodHead || r.body.Len() == 0 {
		return nil
	}

	_, err = r.connection.Write(r.body.Bytes())
	if err != nil {
		return err
	}
	return nil
}

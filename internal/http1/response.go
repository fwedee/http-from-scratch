package http1

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
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

func (r *Response) Send() {

	statusText := http.StatusText(r.StatusCode)
	fmt.Fprintf(r.connection, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText)
	fmt.Fprintf(r.connection, "Content-Length: %d\r\n", r.body.Len())
	fmt.Fprintf(r.connection, "\r\n")

	_, err := r.connection.Write(r.body.Bytes())
	if err != nil {

	}
	r.ResponseHeader.Write(r.connection)
}

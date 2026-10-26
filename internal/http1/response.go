package http1

import (
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
	return r.connection.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.StatusCode = statusCode
	r.wroteHeader = true

	// ensure connection to be closed
	if r.ResponseHeader.Get("Content-Length") == "" && r.ResponseHeader.Get("Transfer-Encoding") == "" {
		if r.ResponseHeader.Get("Connection") == "" {
			r.ResponseHeader.Set("Connection", "close")
		}
	}

	statusText := http.StatusText(r.StatusCode)
	fmt.Fprintf(r.connection, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText)
	r.ResponseHeader.Write(r.connection)
	fmt.Fprint(r.connection, "\r\n")
}

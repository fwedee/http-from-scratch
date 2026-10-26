package http1

import (
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
}

func NewResponse(conn net.Conn) *Response {
	return &Response{connection: conn, ResponseHeader: make(http.Header), StatusCode: 200}
}

func (r *Response) Header() http.Header {
	return r.ResponseHeader
}

func (r *Response) Write(data []byte) (n int, err error) {
	return r.connection.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

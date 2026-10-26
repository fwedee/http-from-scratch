package http1

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func ListenAndServe(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	fmt.Println("Start ListenAndServe")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()

	if handler == nil {
		handler = http.DefaultServeMux
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection: %w", err)
		}
		go HandleConnection(conn, handler)
	}
}

func HandleConnection(conn net.Conn, mux http.Handler) error {

	reader := bufio.NewReader(conn)

	request, err := ParseRequest(reader)
	if err != nil {
		return err
	}

	response := Response{
		Version:    "test",
		StatusCode: 200,
		Resource:   "/",
		Reason:     ""}

	fmt.Print(response.StatusCode)

	mux.ServeHTTP(&response, request)

	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	return nil
}

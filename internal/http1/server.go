package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

func ListenAndServe(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	log.Println("Start ListenAndServe")
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
		return fmt.Errorf("Error parsing request: %w", err)
	}
	fmt.Println("Requestparsing done")

	response := NewResponse(conn)

	mux.ServeHTTP(response, request)
	fmt.Println(response.HTTPVersion)
	fmt.Println("Success")

	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	return nil
}

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

	defer conn.Close()
	reader := bufio.NewReader(conn)
	request, err := ParseRequest(reader)
	log.Printf("Request: %v", request)

	if err != nil {
		return fmt.Errorf("Error parsing request: %w", err)
	}

	// Create new Response
	response := NewResponse(conn)
	log.Printf("Response: %v", response)

	mux.ServeHTTP(response, request)

	// _, _ = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	return nil
}

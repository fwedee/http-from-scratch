package http1

import (
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
		HandleConnection(conn, handler)
	}
}

func HandleConnection(conn net.Conn, mux http.Handler) error {
	panic("HandleConnection: Not Implemented")
	/*
			Functions that might be useful:
			* ParseRequest()
		  * mux.ServeHTTP()
	*/
}

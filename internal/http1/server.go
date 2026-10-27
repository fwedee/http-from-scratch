package http1

import (
	"bufio"
	"errors"
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
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Error closing connection: %v \n", err)
		}
	}()
	reader := bufio.NewReader(conn)
	request, err := ParseRequest(reader)
	log.Printf("Received connection on address %s\n", conn.LocalAddr().String())

	response := NewResponse(conn)
	if err != nil {
		var parseErr *ParseError
		if errors.As(err, &ParseError) {
			response.WriteHeader(parseErr.StatusCode)
			if sendErr := response.Send(); sendErr != nil {
				return fmt.Errorf("send HTTP error response: %w", sendErr)
			}
			return nil
		}
		return fmt.Errorf("read request error: %w", err)
	}

	// Create new Response
	log.Printf("Response: %v", response)
	mux.ServeHTTP(response, request)
	err = response.Send()
	if err != nil {
		return err
	}

	// for later if I don't want to use ServeHTTP anymore and we need to use our own handler
	// if _ , err := conn.Write([]byte(response)); err != nil {
	// 	return fmt.Errorf("Error writing data: %v", err)
	// }
	//

	return nil
}

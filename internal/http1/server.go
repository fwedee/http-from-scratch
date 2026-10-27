package http1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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
		// Adding error handling
		go HandleConnection(conn, handler)
	}
}

func HandleConnection(conn net.Conn, mux http.Handler) error {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		log.Printf("Received connection on address %s\n", conn.LocalAddr().String())

		request, err := ParseRequest(reader)
		response := NewResponse(conn)
		if err != nil {

			if errors.Is(err, io.EOF) {
				return nil
			}

			if parseErr, ok := errors.AsType[*ParseError](err); ok {
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
		_, readErr := io.Copy(io.Discard, request.Body)
		closeErr := request.Body.Close()

		if readErr != nil {
			if errors.Is(readErr, io.ErrUnexpectedEOF) {
				response = NewResponse(conn)
				response.Header().Set("Connection", "close")
				response.WriteHeader(http.StatusBadRequest)
				return response.Send()
			}
			return readErr
		}

		if closeErr != nil {
			return closeErr
		}

		closeConnection := hasHeaderToken(request.Header.Values("Connection"), "close") || hasHeaderToken(response.Header().Values("Connection"), "close")

		if closeConnection {
			response.Header().Set("Connection", "close")
		}
		if err = response.Send(); err != nil {
			return err
		}
		if closeConnection {
			return nil
		}
		// for later if I don't want to use ServeHTTP anymore and we need to use our own handler
		// if _ , err := conn.Write([]byte(response)); err != nil {
		// 	return fmt.Errorf("Error writing data: %v", err)
		// }
		//
	}
}

func hasHeaderToken(values []string, token string) bool {
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			if strings.EqualFold(strings.TrimSpace(part), token) {
				return true
			}
		}
	}
	return false
}

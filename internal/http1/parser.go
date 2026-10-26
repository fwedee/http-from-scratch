package http1

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error not a string: %w", err)
	}

	// Split Statusline, Header and Body
	splitResponse := strings.Split(rawResponse, "\r\n")
	statusline := strings.Split(splitResponse[0], " ")

	// Debuggin Purpose
	log.Println("PRINT TESTING")
	log.Println(statusline[0])
	log.Println(statusline[1])
	log.Println(statusline[2])

	request, err := http.NewRequest(splitResponse[0], splitResponse[1], reader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Request: %w", err)
	}

	return request, nil
}

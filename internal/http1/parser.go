package http1

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	rawResponse, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("Error not a string: %w", err)
	}
	fmt.Println(rawResponse)

	splitResponse := strings.Split(rawResponse, " ")
	fmt.Println(splitResponse[2])

	method := splitResponse[0]
	var statuscode int

	if strings.ToUpper(method) == "GET" {
		statuscode = 200
	} else {
		return nil, fmt.Errorf("Statuscode is wrong")
	}

	response2 := Response{
		Version:    splitResponse[2],
		StatusCode: statuscode,
		Resource:   splitResponse[1],
		Reason:     ""}

	fmt.Print(response2.StatusCode)

	request, err := http.NewRequest("200", splitResponse[1], reader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Request: %w", err)
	}

	return request, nil
}

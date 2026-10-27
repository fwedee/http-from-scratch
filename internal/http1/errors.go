package http1

import (
	"errors"
	"fmt"
)

type ParseError struct {
	StatusCode int
	Operation  string
	Err        error
}

func (e *ParseError) Error() string {
	if e.Err == nil {
		return e.Operation
	}
	return fmt.Sprintf("%s - %v", e.Operation, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

var ErrBadRequest = errors.New("bad request")
var ErrUnsupportedVersion = errors.New("unsupported HTTP version")

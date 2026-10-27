package http1

import "errors"

var ErrBadRequest = errors.New("bad request")
var ErrUnsupportedVersion = errors.New("unsupported HTTP version")

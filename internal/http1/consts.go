package http1

var supporedHTTPMethods = map[string]struct{}{
	"GET": {}, "POST": {}, "HEAD": {}, "PUT": {}, "DELETE": {}, "TRACE": {}, "CONNECT": {},
}

var supportedHTTPVersion = map[string]struct{}{
	"HTTP/1.1": {},
	"HTTP/1.0": {},
}

var (
	StatusBadRequest              int = 400
	StatusOK                      int = 200
	StatusMethodNotAllowed        int = 405
	StatusHTTPVersionNotSupported int = 505
	StatusRequestHeaderTooLarge   int = 431
)

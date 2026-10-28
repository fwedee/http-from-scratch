package http1

var supporedHTTPMethods = map[string]struct{}{
	"GET": {}, "POST": {}, "HEAD": {}, "PUT": {}, "DELETE": {}, "TRACE": {}, "CONNECT": {},
}

var supportedHTTPVersion = map[string]struct{}{
	"HTTP/1.1": {},
	"HTTP/1.0": {},
}

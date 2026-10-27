## RFC
- stat-line followed by a crlf and some headers then an body
- servers only accep a request
- get each of the fields parse them into a hash table by fieldname until empty line
  - read the body as a stream until or equal to the body lenght or the connection is closed
- MUST parse HTTP messages as a sequence of octets in an encoding that is a superset of ASCII
  - not use Unicode because of security vulns
- CR needs to be correctly handeled if bare
- also \n\r need to be handeled correctly
- should ingore at least one empyt line recieved prior to the request-line
- A recipient that receives whitespace between the start-line and the first header field MUST either reject the message as invalid
- HTTP-version  = HTTP-name "/" DIGIT "." DIGIT
-  A server that receives a request-target longer than any URI it wishes to parse MUST respond with a 414 (URI Too Long) status code
- Recipients of an invalid request-line SHOULD respond with either a 400 (Bad Request) error or a 301 (Moved Permanently) redirect with the request-target properly encoded.
- A server MUST respond with a 400 (Bad Request) status code to any HTTP/1.1 request message that lacks a Host header field and to any

## Links
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Messages#http_responses
- https://pkg.go.dev/net/http#ResponseWriter
- https://datatracker.ietf.org/doc/html/rfc9112

## Needs to be implemented

- [ ] Implementation of proper connection closing (HandleConnection returned (connection closed) but keep-alive was expected: <nil>)
- [ ] write on closed connection (writing to server: http1 test: write on closed connection)
- [ ] error for other versions of http (505)
- [ ] version parsing minor, scheme, nonnumeric version (400)
- [ ] space_in_method, in request target should be 400 not 404
- [ ] unexpected EOF for edgecases like CornerCases

- [ ] writing a handler for other paths
- [ ] /echo/{str} endpoint  return a 200 response, with the response body set to given string, and with a Content-Type and Content-Length header.
- [ ] give back useragent
- [ ] give back a file /files/{filename}
- [ ] get a file back POST /files/{filename}
- [ ] accept Accept-Encoding and Content-Encoding
- [ ] Accept multiple encoding shemas Accept-Encoding: encoding-1, encoding-2, encoding-3
- [ ] encode and send encoding stuff

- [ ] 


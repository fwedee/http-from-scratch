package http1

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantMethod string
		wantTarget string
		wantBody   string
		wantErr    bool
	}{
		{
			name:       "GET without body",
			raw:        "GET /users HTTP/1.1\r\nHost: localhost\r\n\r\n",
			wantMethod: "GET",
			wantTarget: "/users",
			wantBody:   "",
			wantErr:    false,
		},
		{
			name:       "POST with body and extra trailing bytes",
			raw:        "POST /submit HTTP/1.1\r\nHost: localhost\r\nContent-Length: 11\r\n\r\nhello worldsakdad",
			wantMethod: "POST",
			wantTarget: "/submit",
			wantBody:   "hello world",
			wantErr:    false,
		},
		{
			name:    "Nothing",
			raw:     "GARBAGE_NO_SPACES\r\n\r\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tt.raw))
			req, err := ParseRequest(r)

			if (err != nil) != tt.wantErr {
				t.Fatalf("got err = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if req.Method != tt.wantMethod {
				t.Errorf("Method = %q, want %q", req.Method, tt.wantMethod)
			}
			if req.URL.String() != tt.wantTarget {
				t.Errorf("URL = %q, want %q", req.URL.String(), tt.wantTarget)
			}

			if req.Body == nil {
				if tt.wantBody != "" {
					t.Errorf("Body = nil, want %q", tt.wantBody)
				}
				return
			}

			defer req.Body.Close()
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("reading body failed: %v", err)
			}
			if string(bodyBytes) != tt.wantBody {
				t.Errorf("Body = %q, want %q", string(bodyBytes), tt.wantBody)
			}
		})
	}
}

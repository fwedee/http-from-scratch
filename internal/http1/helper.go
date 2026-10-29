package http1

import (
	"strings"
	"unicode/utf8"
)

func rotateLeftEfficient(s string) string {
	if len(s) == 0 {
		return s
	}
	_, size := utf8.DecodeRuneInString(s)
	return s[size:] + s[:size]
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

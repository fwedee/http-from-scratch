package http1

import "unicode/utf8"

func rotateLeftEfficient(s string) string {
	if len(s) == 0 {
		return s
	}
	_, size := utf8.DecodeRuneInString(s)
	return s[size:] + s[:size]
}

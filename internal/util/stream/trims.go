package stream

import "strings"

func TrimQuote() func(string) string {
	return func(s string) string {
		return strings.Trim(s, `"`)
	}
}

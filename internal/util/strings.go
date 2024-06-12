package util

func EndWith(s string, checker string) bool {
	return len(s) >= len(checker) && s[len(s)-len(checker):] == checker
}

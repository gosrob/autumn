package pkginfo_test

import (
	"testing"

	"github.com/gosrob/autumn/internal/util/pkginfo"
)

func TestGetPackageFromFilePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Typical case with multiple directories",
			path:     "/Users/zhangruobin/Episodes/travel/coding/golang/autumn/internal/util/pkginfo/info.go",
			expected: "pkginfo",
		},
		{
			name:     "Root directory",
			path:     "/info.go",
			expected: "",
		},
		{
			name:     "Single directory",
			path:     "pkginfo/info.go",
			expected: "pkginfo",
		},
		{
			name:     "No directory",
			path:     "info.go",
			expected: "",
		},
		{
			name:     "Empty string",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pkginfo.GetPackageFromFilePath(tt.path)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

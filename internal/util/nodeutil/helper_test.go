package nodeutil

import (
	"testing"
)

func TestGetType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Type
	}{
		{
			name:  "simple type",
			input: "int",
			expected: Type{
				PureType:       "int",
				IsArray:        false,
				IsPointer:      false,
				IsArrayPointer: false,
			},
		},
		{
			name:  "pointer type",
			input: "*int",
			expected: Type{
				PureType:       "int",
				IsArray:        false,
				IsPointer:      true,
				IsArrayPointer: false,
			},
		},
		{
			name:  "array type",
			input: "[]int",
			expected: Type{
				PureType:       "int",
				IsArray:        true,
				IsPointer:      false,
				IsArrayPointer: false,
			},
		},
		{
			name:  "array pointer type",
			input: "[]*int",
			expected: Type{
				PureType:       "int",
				IsArray:        true,
				IsPointer:      false,
				IsArrayPointer: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetType(tt.input)
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

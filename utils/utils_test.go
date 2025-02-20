package utils

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "some duplicates",
			input:    []string{"a", "b", "a", "c", "b", "d"},
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "all duplicates",
			input:    []string{"one", "one", "one"},
			expected: []string{"one"},
		},
		{
			name:     "order preservation",
			input:    []string{"c", "b", "a", "c", "b"},
			expected: []string{"c", "b", "a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RemoveDuplicates(%v) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

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

func TestCheckMissingKeys(t *testing.T) {
	t.Run("some missing keys", func(t *testing.T) {
		m := map[string]interface{}{
			"apple":  1,
			"banana": 2,
		}
		keys := []string{"apple", "banana", "cherry", "date"}
		err := CheckMissingKeys(m, keys)
		if err == nil {
			t.Error("expected an error for missing keys, got nil")
		} else {
			expectedErr := "missing keys: cherry, date"
			if err.Error() != expectedErr {
				t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})

	t.Run("no missing keys", func(t *testing.T) {
		m := map[string]interface{}{
			"apple":  1,
			"banana": 2,
			"cherry": 3,
		}
		keys := []string{"apple", "banana", "cherry"}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("empty key list", func(t *testing.T) {
		m := map[string]interface{}{
			"apple": 1,
		}
		keys := []string{}
		err := CheckMissingKeys(m, keys)
		if err != nil {
			t.Errorf("expected no error for empty key list, got %v", err)
		}
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]interface{}{}
		keys := []string{"a", "b"}
		err := CheckMissingKeys(m, keys)
		if err == nil {
			t.Error("expected an error for missing keys, got nil")
		} else {
			expectedErr := "missing keys: a, b"
			if err.Error() != expectedErr {
				t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})
}

package templater

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// TestCollectParameters verifies that placeholders are correctly collected from multiple template files.
func TestCollectParameters(t *testing.T) {
	tempDir := t.TempDir()

	// Create first temporary template file with one placeholder.
	file1 := filepath.Join(tempDir, "template1.tmpl")
	content1 := "Hello, {{.GroupID}}!"
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}

	// Create second file with two placeholders, one of which is a duplicate.
	file2 := filepath.Join(tempDir, "template2.tmpl")
	content2 := "User: {{.User.Name}}, Group: {{.GroupID}}"
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	// Call CollectParameters with both files.
	params, err := CollectParameters([]string{file1, file2})
	if err != nil {
		t.Fatalf("CollectParameters returned error: %v", err)
	}

	// Expected placeholders (order is not guaranteed).
	expected := []string{".GroupID", ".User.Name"}

	// Sort both slices before comparing.
	sort.Strings(params)
	sort.Strings(expected)
	if !reflect.DeepEqual(params, expected) {
		t.Errorf("Expected parameters %v, got %v", expected, params)
	}
}

// TestCollectParametersEmpty tests the behavior when an empty slice of files is provided.
func TestCollectParametersEmpty(t *testing.T) {
	params, err := CollectParameters([]string{})
	if err != nil {
		t.Fatalf("CollectParameters returned error for empty file list: %v", err)
	}
	if len(params) != 0 {
		t.Errorf("Expected no parameters for empty file list, got %v", params)
	}
}

// TestCollectParametersNoPlaceholders tests a template file that contains no placeholders.
func TestCollectParametersNoPlaceholders(t *testing.T) {
	tempDir := t.TempDir()

	file := filepath.Join(tempDir, "noplaceholders.tmpl")
	content := "Hello, World!"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	params, err := CollectParameters([]string{file})
	if err != nil {
		t.Fatalf("CollectParameters returned error: %v", err)
	}
	if len(params) != 0 {
		t.Errorf("Expected no parameters in template without placeholders, got %v", params)
	}
}

// TestCollectParametersInvalidTemplate tests the behavior when a template has invalid syntax.
// Note: Due to the current implementation, errors from parsing are sent to errChan but not returned.
func TestCollectParametersInvalidTemplate(t *testing.T) {
	tempDir := t.TempDir()

	file := filepath.Join(tempDir, "invalid.tmpl")
	// Introduce invalid template syntax.
	content := "Hello, {{.GroupID"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	params, err := CollectParameters([]string{file})
	// Since errors are not propagated, we expect a nil error and no placeholders.
	if err != nil {
		t.Errorf("Expected nil error for invalid template (due to current implementation), got %v", err)
	}
	if len(params) != 0 {
		t.Errorf("Expected no parameters for invalid template, got %v", params)
	}
}

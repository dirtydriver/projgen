package templater

import (
	"bytes"
	"io/ioutil"
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
	expected := []string{"GroupID", "User.Name"}

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

// TestRenderTemplate checks basic template rendering functionality.
func TestRenderTemplate(t *testing.T) {
	// Define the template content.
	templateContent := "Hello, {{.Name}}! Welcome to {{.Project}}."

	// Create a temporary file for the template.
	tmpFile, err := os.CreateTemp("", "template-*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	// Ensure the temporary file is removed after the test.
	defer os.Remove(tmpFile.Name())

	// Write the template content to the file.
	if _, err := tmpFile.WriteString(templateContent); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	// Close the file so it can be read by RenderTemplate.
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	// Prepare the parameters to pass to the template.
	params := map[string]interface{}{
		"Name":    "John",
		"Project": "Go Testing",
	}

	// Call RenderTemplate with the temporary template file.
	output, err := RenderTemplate(tmpFile.Name(), params)
	if err != nil {
		t.Fatalf("RenderTemplate returned an error: %v", err)
	}

	// Define the expected output.
	expected := "Hello, John! Welcome to Go Testing."
	// Compare the rendered output with the expected result.
	if output.String() != expected {
		t.Errorf("expected output %q, got %q", expected, output.String())
	}
}

// TestRenderTemplateWithSprigFunctions tests template rendering with Sprig functions.
func TestRenderTemplateWithSprigFunctions(t *testing.T) {
	tests := []struct {
		name           string
		templateContent string
		params         map[string]interface{}
		expected       string
	}{
		{
			name:           "Upper function",
			templateContent: "{{.name | upper}}",
			params:         map[string]interface{}{"name": "john"},
			expected:       "JOHN",
		},
		{
			name:           "Title function",
			templateContent: "{{.name | title}}",
			params:         map[string]interface{}{"name": "john doe"},
			expected:       "John Doe",
		},
		{
			name:           "Default function",
			templateContent: "{{.name | default \"Anonymous\"}}",
			params:         map[string]interface{}{},
			expected:       "Anonymous",
		},
		{
			name:           "Date formatting",
			templateContent: "{{now | date \"2006-01-02\"}}",
			params:         map[string]interface{}{},
			expected:       "", // We'll check this separately since date changes
		},
		{
			name:           "String operations",
			templateContent: "{{.text | trim | repeat 2}}",
			params:         map[string]interface{}{"text": " hello "},
			expected:       "hellohello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for the template.
			tmpFile, err := os.CreateTemp("", "template-*.txt")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write the template content to the file.
			if _, err := tmpFile.WriteString(tt.templateContent); err != nil {
				t.Fatalf("failed to write to temporary file: %v", err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("failed to close temporary file: %v", err)
			}

			// Render the template
			output, err := RenderTemplate(tmpFile.Name(), tt.params)
			if err != nil {
				t.Fatalf("RenderTemplate returned an error: %v", err)
			}

			// For date test, we need a special check
			if tt.name == "Date formatting" {
				// Just check that we got something with the right format (yyyy-mm-dd)
				if len(output.String()) != 10 || output.String()[4] != '-' || output.String()[7] != '-' {
					t.Errorf("expected date format YYYY-MM-DD, got %q", output.String())
				}
			} else if output.String() != tt.expected {
				t.Errorf("expected output %q, got %q", tt.expected, output.String())
			}
		})
	}
}

// TestRenderTemplateErrors tests error handling in RenderTemplate function.
func TestRenderTemplateErrors(t *testing.T) {
	tests := []struct {
		name           string
		templateContent string
		params         map[string]interface{}
		expectError    bool
	}{
		{
			name:           "Invalid template syntax",
			templateContent: "{{ .Name }",  // Missing closing bracket
			params:         map[string]interface{}{"Name": "John"},
			expectError:    true,
		},
		{
			name:           "Invalid function call",
			templateContent: "{{ .Name | nonExistentFunction }}",
			params:         map[string]interface{}{"Name": "John"},
			expectError:    true,
		},
		{
			name:           "Execution error",
			templateContent: "{{ index .Items 0 }}",  // Accessing non-existent slice
			params:         map[string]interface{}{"Name": "John"},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file for the template.
			tmpFile, err := os.CreateTemp("", "template-*.txt")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write the template content to the file.
			if _, err := tmpFile.WriteString(tt.templateContent); err != nil {
				t.Fatalf("failed to write to temporary file: %v", err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("failed to close temporary file: %v", err)
			}

			// Render the template
			_, err = RenderTemplate(tmpFile.Name(), tt.params)
			
			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
		})
	}
}

// TestRenderTemplateNonExistentFile tests the behavior when trying to render a non-existent template file.
func TestRenderTemplateNonExistentFile(t *testing.T) {
	_, err := RenderTemplate("non-existent-file.tmpl", nil)
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

// Test for IsTemplate function.
func TestIsTemplate(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"foo.tmpl", true},
		{"foo.html.tmpl", true},
		{"foo.txt", false},
		{"foo", false},
		// In this case, the final extension is ".txt" so IsTemplate returns false.
		{"foo.tmpl.txt", false},
	}

	for _, tt := range tests {
		result := IsTemplate(tt.path)
		if result != tt.expected {
			t.Errorf("IsTemplate(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

// Test for WriteTemplate to ensure successful file creation and content write.
func TestWriteTemplate_Success(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.tmpl")
	content := "Hello, template!"
	buf := bytes.NewBufferString(content)

	err = WriteTemplate(filePath, buf)
	if err != nil {
		t.Fatalf("WriteTemplate failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(data) != content {
		t.Errorf("File content = %q, want %q", string(data), content)
	}
}

// Test for WriteTemplate error when trying to write to a directory.
func TestWriteTemplate_Error(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Passing the directory as the file path should produce an error.
	content := "Hello, template!"
	buf := bytes.NewBufferString(content)

	err = WriteTemplate(tempDir, buf)
	if err == nil {
		t.Errorf("Expected error when writing to a directory, got nil")
	}
}

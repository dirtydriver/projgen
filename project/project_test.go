package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerate verifies that Generate renders template files and copies static files correctly.
func TestGenerate(t *testing.T) {
	// Create temporary directories for templates and output.
	templateDir, err := os.MkdirTemp("", "template")
	if err != nil {
		t.Fatalf("failed to create temp template directory: %v", err)
	}
	defer os.RemoveAll(templateDir)

	outputDir, err := os.MkdirTemp("", "output")
	if err != nil {
		t.Fatalf("failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Define project type and create a corresponding subdirectory.
	projectType := "test"
	projectTemplateDir := filepath.Join(templateDir, projectType)
	if err := os.MkdirAll(projectTemplateDir, 0755); err != nil {
		t.Fatalf("failed to create project template directory: %v", err)
	}

	// Create a template file (ending in .tmpl).
	templateFilePath := filepath.Join(projectTemplateDir, "greeting.txt.tmpl")
	templateContent := "Hello, {{.Name}}!"
	if err := os.WriteFile(templateFilePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	// Create a static file.
	staticFilePath := filepath.Join(projectTemplateDir, "readme.md")
	staticContent := "This is a static file."
	if err := os.WriteFile(staticFilePath, []byte(staticContent), 0644); err != nil {
		t.Fatalf("failed to write static file: %v", err)
	}

	// Prepare parameters for the Generate function.
	params := map[string]interface{}{
		"Name": "World",
	}

	// Call the Generate function.
	if err := Generate(projectTemplateDir, outputDir, params); err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	// Verify that the rendered template file exists with the expected content.
	// The output file should have the same relative path, but without the .tmpl extension.
	renderedFilePath := filepath.Join(outputDir, "greeting.txt")
	data, err := os.ReadFile(renderedFilePath)
	if err != nil {
		t.Fatalf("failed to read rendered template file: %v", err)
	}
	expectedRendered := "Hello, World!"
	if strings.TrimSpace(string(data)) != expectedRendered {
		t.Errorf("rendered content mismatch: expected %q, got %q", expectedRendered, string(data))
	}

	// Verify that the static file was copied correctly.
	copiedStaticFilePath := filepath.Join(outputDir, "readme.md")
	staticData, err := os.ReadFile(copiedStaticFilePath)
	if err != nil {
		t.Fatalf("failed to read copied static file: %v", err)
	}
	if string(staticData) != staticContent {
		t.Errorf("static file content mismatch: expected %q, got %q", staticContent, string(staticData))
	}
}

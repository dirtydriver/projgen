package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dirtydriver/projgen/filescheck"
	"github.com/dirtydriver/projgen/templater"
)

// Generate creates a new project from a template directory using the provided parameters.
// It copies all files from the template, rendering any .tmpl files with the given parameters.
func Generate(templateDir, outputDir string, paramsMap map[string]interface{}) error {
	// Validate and extract the project type.
	projectTypeVal, ok := paramsMap["type"]
	if !ok {
		return errors.New("missing 'type' parameter")
	}
	projectType, ok := projectTypeVal.(string)
	if !ok {
		return errors.New("'type' parameter is not a string")
	}

	// Construct the path to the project's templates.
	templatePath := filepath.Join(templateDir, projectType)
	fileList, err := filescheck.FilesInDirectories(templatePath)
	if err != nil {
		return err
	}

	for _, file := range fileList {
		// Compute the relative path from the template directory.
		relPath, err := filepath.Rel(templatePath, file)
		if err != nil {
			return fmt.Errorf("failed to determine relative path for %s: %w", file, err)
		}
		targetPath := filepath.Join(outputDir, relPath)

		if templater.IsTemplate(file) {
			// Render the template with the provided parameters.
			rendered, err := templater.RenderTemplate(file, paramsMap)
			if err != nil {
				return err
			}

			// Remove the .tmpl extension from the target path.
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")

			// Ensure the target directory exists.
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
			}

			// Write the rendered content to the target path.
			if err := templater.WriteTemplate(targetPath, &rendered); err != nil {
				return err
			}
		} else {
			// For non-template files, ensure the directory exists.
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", targetPath, err)
			}

			// Copy the file.
			if err := filescheck.CopyFile(file, filepath.Dir(targetPath)); err != nil {
				return err
			}
		}
	}
	return nil
}

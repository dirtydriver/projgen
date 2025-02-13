package filescheck

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func filesInDirectories(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err

}

func copyFiles(files []string, target_dir string) error {

	for _, file := range files {
		input, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", file, err)
		}

		if err := os.WriteFile(target_dir, input, 0644); err != nil {
			return err
		}

		return err
	}
}

func createDirectory(dir_name string) error {

	isDirExists := isDirectoryExists(dir_name)
	if isDirExists {
		return error.Error("The directory already exists")
	}

	if err := os.MkdirAll(dir_name, os.ModePerm); err != nil {
		return err
	}

}
func isDirectoryExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return false

}

package filescheck

import (
	"fmt"
	"os"
	"path/filepath"
)

func FilesInDirectories(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err

}

func CopyFiles(files []string, targetDir string) error {

	for _, file := range files {
		input, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", file, err)
		}
		destPath := filepath.Join(targetDir, filepath.Base(file))
		if err := os.WriteFile(destPath, input, 0644); err != nil {
			return err
		}

	}
	return nil
}

func CreateDirectory(dirName string) error {
	if info, err := os.Stat(dirName); err == nil {
		if info.IsDir() {
			return fmt.Errorf("the directory already exists")
		}
		// If the path exists but isn't a directory, you might want to handle that differently.
		return fmt.Errorf("a non-directory file exists at %s", dirName)
	} else if !os.IsNotExist(err) {
		// An unexpected error occurred
		return err
	}

	// The directory does not exist, so create it.
	return os.MkdirAll(dirName, os.ModePerm)
}

func IsDirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

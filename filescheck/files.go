package filescheck

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ReadParamsFromFile(paramFilePath string, paramsMap *map[string]interface{}) error {
	file, err := os.Open(paramFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines or lines starting with a comment marker
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid parameter format on line %d: %s. Expected key=value", lineNumber, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := (*paramsMap)[key]; !exists {
			(*paramsMap)[key] = value
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

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

func FindTemplateFiles(path string, pattern string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), pattern) {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil

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

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
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

func main() {

	entries, err := filesInDirectories(".")
	fmt.Println(isDirectoryExists("cmd"))
	fmt.Println(entries)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err)
	}

	for _, e := range entries {
		fmt.Println(e)
	}
}

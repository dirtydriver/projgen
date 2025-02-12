package main

import (
	"fmt"
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

}

func main() {

	entries, err := filesInDirectories("/mnt/f/Progproject/GO")
	fmt.Println(entries)
	if err != nil {
		log.Fatalf("Something went wrong: %s", err)
	}

	for _, e := range entries {
		fmt.Println(e)
	}
}

package filescheck

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFilesInDirectories creates a temporary directory with files and subdirectories,
// then verifies that FilesInDirectories returns the correct list of file paths.
func TestFilesInDirectories(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "testFilesInDirectories")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a subdirectory.
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// Create files in both tempDir and subDir.
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := FilesInDirectories(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// We expect to find exactly two files.
	if len(files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(files))
	}

	// Optionally, you can check that both file1 and file2 are in the list.
	found := make(map[string]bool)
	for _, f := range files {
		found[f] = true
	}
	if !found[file1] || !found[file2] {
		t.Fatalf("Files returned do not match expected files: %v", files)
	}
}

// TestCopyFiles creates a temporary source file and destination directory,
// then uses CopyFiles to copy the file and verifies that the content is identical.
func TestCopyFiles(t *testing.T) {
	// Create temporary source directory.
	srcDir, err := os.MkdirTemp("", "testCopyFiles_src")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(srcDir)

	// Create temporary destination directory.
	destDir, err := os.MkdirTemp("", "testCopyFiles_dest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(destDir)

	// Create a source file.
	srcFile := filepath.Join(srcDir, "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Call CopyFiles to copy the file to destDir.
	if err := CopyFiles([]string{srcFile}, destDir); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify that the copied file exists and its content matches.
	copiedFile := filepath.Join(destDir, "test.txt")
	copiedContent, err := os.ReadFile(copiedFile)
	if err != nil {
		t.Fatalf("Expected file %q to be copied, got error: %v", copiedFile, err)
	}
	if string(copiedContent) != string(content) {
		t.Fatalf("Expected file content %q, got %q", content, copiedContent)
	}
}

// TestCreateDirectory tests that CreateDirectory creates a directory if it does not exist
// and returns an error if the directory already exists.
func TestCreateDirectory(t *testing.T) {
	// Create a temporary parent directory.
	parentDir, err := os.MkdirTemp("", "testCreateDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(parentDir)

	newDir := filepath.Join(parentDir, "newDir")
	// Ensure the directory does not exist.
	if IsDirectoryExists(newDir) {
		t.Fatalf("Expected directory %s to not exist", newDir)
	}

	// Create the directory.
	if err := CreateDirectory(newDir); err != nil {
		t.Fatalf("Expected no error when creating directory, got %v", err)
	}

	// Now the directory should exist.
	if !IsDirectoryExists(newDir) {
		t.Fatalf("Expected directory %s to exist", newDir)
	}

	// Try to create the same directory again and expect an error.
	if err := CreateDirectory(newDir); err == nil {
		t.Fatalf("Expected error when creating an existing directory, got nil")
	}
}

// TestIsDirectoryExists checks that IsDirectoryExists returns true for directories,
// and false for files or non-existent paths.

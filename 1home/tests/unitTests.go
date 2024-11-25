package tests

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"1home/config"
	"1home/filesystem"
)

func setupTestVFS(t *testing.T) *filesystem.VirtualFileSystem {

	zipPath := createTestZIP(t)
	vfs, err := filesystem.NewVFS(zipPath)
	if err != nil {
		t.Fatalf("Failed to create VFS: %v", err)
	}
	return vfs
}

func createTestZIP(t *testing.T) string {

	tempDir, err := os.MkdirTemp("", "test-vfs-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFiles := []struct {
		path     string
		contents string
	}{
		{filepath.Join(tempDir, "test1.txt"), "Hello"},
		{filepath.Join(tempDir, "test2.txt"), "World"},
	}

	for _, file := range testFiles {
		err := os.WriteFile(file.path, []byte(file.contents), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	zipPath := filepath.Join(os.TempDir(), "test-vfs.zip")
	err = createZIP(tempDir, zipPath)
	if err != nil {
		t.Fatalf("Failed to create ZIP: %v", err)
	}
	return zipPath
}

func createZIP(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = writer.Write([]byte{})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func TestCD(t *testing.T) {
	vfs := setupTestVFS(t)

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{"Valid directory", ".", false},
		{"Invalid directory", "/nonexistent", true},
		{"Prevent traversal", "../../../", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vfs.ChangeDirectory(tt.path)
			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestLS(t *testing.T) {
	vfs := setupTestVFS(t)

	files, err := vfs.ListDirectory()
	if err != nil {
		t.Fatalf("Failed to list directory: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}
}

func TestCP(t *testing.T) {
	vfs := setupTestVFS(t)

	tests := []struct {
		name        string
		src         string
		dst         string
		expectError bool
	}{
		{"Valid copy", "test1.txt", "test3.txt", false},
		{"Nonexistent source", "nonexistent.txt", "test4.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vfs.CopyFile(tt.src, tt.dst)
			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTouch(t *testing.T) {
	vfs := setupTestVFS(t)

	tests := []struct {
		name        string
		filename    string
		expectError bool
	}{
		{"Valid touch", "newfile.txt", false},
		{"Empty filename", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vfs.Touch(tt.filename)
			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestWhoami(t *testing.T) {

	cfg, err := config.LoadConfig("testdata/config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	tests := []struct {
		name     string
		expected string
	}{
		{"Check username", "testuser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cfg.Username != tt.expected {
				t.Errorf("Expected username %s, got %s", tt.expected, cfg.Username)
			}
		})
	}
}

package filesystem

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractZip(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	zipPath := filepath.Join(srcDir, "test.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}

	writer := zip.NewWriter(zipFile)

	fw, err := writer.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = fw.Write([]byte("hello"))

	subFw, err := writer.Create("subdir/nested.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = subFw.Write([]byte("nested"))

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	zipFile.Close()

	extractDir := filepath.Join(destDir, "extracted")
	if err := ExtractZip(zipPath, extractDir); err != nil {
		t.Fatalf("ExtractZip() failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(extractDir, "test.txt"))
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("extracted content = %q, want \"hello\"", string(data))
	}

	data, err = os.ReadFile(filepath.Join(extractDir, "subdir", "nested.txt"))
	if err != nil {
		t.Fatalf("failed to read nested file: %v", err)
	}
	if string(data) != "nested" {
		t.Errorf("nested content = %q, want \"nested\"", string(data))
	}
}

func TestValidateZipPath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"test.txt", false},
		{"subdir/test.txt", false},
		{"../evil.txt", true},
		{"/absolute/path", true},
		{"C:\\Windows\\test", true},
	}

	for _, tt := range tests {
		err := validateZipPath(tt.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateZipPath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
		}
	}
}

func TestEnsureDir(t *testing.T) {
	dir := t.TempDir()
	testDir := filepath.Join(dir, "new", "nested", "dir")

	if err := ensureDir(testDir); err != nil {
		t.Fatalf("ensureDir() failed: %v", err)
	}

	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("path is not a directory")
	}

	if err := ensureDir(testDir); err != nil {
		t.Fatalf("ensureDir() on existing dir failed: %v", err)
	}

	testFile := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatal(err)
	}

	err = ensureDir(testFile)
	if err == nil {
		t.Error("ensureDir() on file should fail")
	}
}

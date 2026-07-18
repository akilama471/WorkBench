package pkg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akilama471/WorkBench/internal/filesystem"
)

func TestManifestValidate(t *testing.T) {
	valid := &Manifest{
		ID:           "php",
		Version:      "8.3.30",
		Platform:     "windows",
		Architecture: "x64",
		Download: ManifestDownload{
			URL:    "https://example.com/php.zip",
			SHA256: "abc123",
		},
	}

	if err := valid.Validate(); err != nil {
		t.Errorf("valid manifest should pass validation: %v", err)
	}

	tests := []struct {
		name   string
		modify func(*Manifest)
	}{
		{"missing ID", func(m *Manifest) { m.ID = "" }},
		{"missing version", func(m *Manifest) { m.Version = "" }},
		{"missing URL", func(m *Manifest) { m.Download.URL = "" }},
		{"missing platform", func(m *Manifest) { m.Platform = "" }},
		{"missing architecture", func(m *Manifest) { m.Architecture = "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := *valid
			tt.modify(&m)
			if err := m.Validate(); err == nil {
				t.Error("expected validation error")
			}
		})
	}
}

func TestLoadManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "test.json")

	content := `{
		"id": "php",
		"name": "PHP",
		"version": "8.3.30",
		"platform": "windows",
		"architecture": "x64",
		"download": {
			"url": "https://example.com/php.zip",
			"sha256": "abc123"
		},
		"archive": "zip",
		"install": {
			"directory": "php/8.3.30"
		}
	}`

	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test manifest: %v", err)
	}

	m, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest() failed: %v", err)
	}

	if m.ID != "php" {
		t.Errorf("ID = %q, want \"php\"", m.ID)
	}
	if m.Version != "8.3.30" {
		t.Errorf("Version = %q, want \"8.3.30\"", m.Version)
	}
}

func TestRepository(t *testing.T) {
	repo := NewRepository()

	m := &Manifest{
		ID:           "php",
		Version:      "8.3.30",
		Platform:     "windows",
		Architecture: "x64",
		Download: ManifestDownload{
			URL: "https://example.com/php.zip",
		},
	}

	repo.Register(m)

	found, exists := repo.Find("php", "8.3.30", "windows", "x64")
	if !exists {
		t.Fatal("Find() should find registered manifest")
	}
	if found.ID != "php" {
		t.Errorf("found.ID = %q, want \"php\"", found.ID)
	}

	_, exists = repo.Find("php", "8.3.30", "linux", "x64")
	if exists {
		t.Error("Find() should not find mismatched platform")
	}

	all := repo.All()
	if len(all) != 1 {
		t.Errorf("All() returned %d, want 1", len(all))
	}
}

func TestVerifyChecksum(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)
	mgr := NewManager(paths, nil)

	testFile := filepath.Join(dir, "test.txt")
	content := []byte("hello world")
	if err := os.WriteFile(testFile, content, 0o644); err != nil {
		t.Fatal(err)
	}

	valid, err := mgr.VerifyChecksum(testFile, "")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Error("empty checksum should return true")
	}

	valid, err = mgr.VerifyChecksum(testFile, "wrong_checksum")
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Error("wrong checksum should return false")
	}

	valid, err = mgr.VerifyChecksum("/nonexistent", "abc")
	if err != nil {
		return
	}
	_ = valid
}

func TestNewPaths(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)
	if paths.Root() != dir {
		t.Errorf("Root() = %q, want %q", paths.Root(), dir)
	}
}

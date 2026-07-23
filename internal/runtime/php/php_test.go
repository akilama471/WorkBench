package php

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)

func TestRuntimeInstalledVersions(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)

	phpBinDir := filepath.Join(paths.Bin(), "php")
	if err := os.MkdirAll(filepath.Join(phpBinDir, "8.1.29"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(phpBinDir, "8.1.29", "php.exe"), []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(phpBinDir, "8.3.30"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(phpBinDir, "8.3.30", "php.exe"), []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(phpBinDir, "notaphp"), 0o755); err != nil {
		t.Fatal(err)
	}

	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	versions, err := rt.InstalledVersions()
	if err != nil {
		t.Fatalf("InstalledVersions() failed: %v", err)
	}

	if len(versions) != 2 {
		t.Fatalf("InstalledVersions() returned %d, want 2", len(versions))
	}

	found81 := false
	found83 := false
	for _, v := range versions {
		if v.Version == "8.1.29" {
			found81 = true
		}
		if v.Version == "8.3.30" {
			found83 = true
		}
	}

	if !found81 {
		t.Error("version 8.1.29 not found")
	}
	if !found83 {
		t.Error("version 8.3.30 not found")
	}
}

func TestRuntimeInstalledVersionsFromPackageName(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)

	phpBinDir := filepath.Join(paths.Bin(), "php")

	pkgDir := filepath.Join(phpBinDir, "php-8.3.30-Win32-vs16-x64")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "php.exe"), []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	versions, err := rt.InstalledVersions()
	if err != nil {
		t.Fatalf("InstalledVersions() failed: %v", err)
	}

	if len(versions) != 1 {
		t.Fatalf("InstalledVersions() returned %d, want 1", len(versions))
	}

	if versions[0].Version != "8.3.30" {
		t.Errorf("Version = %q, want \"8.3.30\"", versions[0].Version)
	}
}

func TestExtractVersion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"8.3.30", "8.3.30"},
		{"php-8.3.30-Win32-vs16-x64", "8.3.30"},
		{"php-8.1.29-nts-Win32-vs16-x64", "8.1.29"},
		{"php-8.2.20-Win32-vs16-x64", "8.2.20"},
		{"notaphp", ""},
		{"php", ""},
		{"php-8.3.30", "8.3.30"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractVersion(tt.input)
			if got != tt.want {
				t.Errorf("extractVersion(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRuntimeActiveVersion(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)
	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	active, err := rt.ActiveVersion()
	if err != nil {
		t.Fatalf("ActiveVersion() failed: %v", err)
	}
	if active != nil {
		t.Error("ActiveVersion() should return nil when no active version")
	}
}

func TestRuntimeUseAndActive(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)

	phpBinDir := filepath.Join(paths.Bin(), "php", "8.3.30")
	if err := os.MkdirAll(phpBinDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(phpBinDir, "php.exe"), []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	if err := rt.Use("8.3.30"); err != nil {
		t.Fatalf("Use(\"8.3.30\") failed: %v", err)
	}

	active, err := rt.ActiveVersion()
	if err != nil {
		t.Fatalf("ActiveVersion() failed: %v", err)
	}
	if active == nil {
		t.Fatal("ActiveVersion() returned nil after Use()")
	}
	if active.Version != "8.3.30" {
		t.Errorf("active version = %q, want \"8.3.30\"", active.Version)
	}
}

func TestRuntimeUseWithPackageName(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)

	pkgDir := filepath.Join(paths.Bin(), "php", "php-8.3.30-Win32-vs16-x64")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pkgDir, "php.exe"), []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	if err := rt.Use("8.3.30"); err != nil {
		t.Fatalf("Use(\"8.3.30\") failed: %v", err)
	}

	active, err := rt.ActiveVersion()
	if err != nil {
		t.Fatalf("ActiveVersion() failed: %v", err)
	}
	if active == nil {
		t.Fatal("ActiveVersion() returned nil after Use()")
	}
	if active.Version != "8.3.30" {
		t.Errorf("active version = %q, want \"8.3.30\"", active.Version)
	}
}

func TestRuntimeUseVersionNotFound(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)
	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	err := rt.Use("9.9.99")
	if err == nil {
		t.Error("Use(\"9.9.99\") should fail for nonexistent version")
	}
}

func TestRuntimeUseInvalidInstallation(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)

	phpBinDir := filepath.Join(paths.Bin(), "php", "8.3.0")
	if err := os.MkdirAll(phpBinDir, 0o755); err != nil {
		t.Fatal(err)
	}

	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	err := rt.Use("8.3.0")
	if err == nil {
		t.Error("Use(\"8.3.0\") should fail for invalid installation")
	}
}

func TestRuntimeIDs(t *testing.T) {
	dir := t.TempDir()
	paths := filesystem.NewPaths(dir)
	log := logger.New(logger.LevelDebug, nil)
	rt := NewRuntime(paths, log)

	if rt.ID() != "php" {
		t.Errorf("ID() = %q, want \"php\"", rt.ID())
	}
	if rt.Name() != "PHP" {
		t.Errorf("Name() = %q, want \"PHP\"", rt.Name())
	}
}

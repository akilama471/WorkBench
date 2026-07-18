package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectorDetectType(t *testing.T) {
	d := NewDetector()

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "artisan"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	typ, err := d.DetectType(dir)
	if err != nil {
		t.Fatal(err)
	}
	if typ != "laravel" {
		t.Errorf("DetectType() = %q, want \"laravel\"", typ)
	}

	dir2 := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir2, "go.mod"), []byte("module test"), 0o644); err != nil {
		t.Fatal(err)
	}
	typ, err = d.DetectType(dir2)
	if err != nil {
		t.Fatal(err)
	}
	if typ != "go" {
		t.Errorf("DetectType() = %q, want \"go\"", typ)
	}

	dir3 := t.TempDir()
	typ, err = d.DetectType(dir3)
	if err != nil {
		t.Fatal(err)
	}
	if typ != "generic" {
		t.Errorf("DetectType() = %q, want \"generic\"", typ)
	}
}

func TestManagerAddAndGet(t *testing.T) {
	mgr := NewManager()
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0o644); err != nil {
		t.Fatal(err)
	}

	p, err := mgr.Add(dir)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}
	if p.Type != "go" {
		t.Errorf("project.Type = %q, want \"go\"", p.Type)
	}

	got, exists := mgr.Get(dir)
	if !exists {
		t.Fatal("Get() should find added project")
	}
	if got.Name != filepath.Base(dir) {
		t.Errorf("project.Name = %q", got.Name)
	}
}

func TestManagerList(t *testing.T) {
	mgr := NewManager()
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	_, _ = mgr.Add(dir1)
	_, _ = mgr.Add(dir2)

	if mgr.Count() != 2 {
		t.Errorf("Count() = %d, want 2", mgr.Count())
	}

	list := mgr.List()
	if len(list) != 2 {
		t.Errorf("List() returned %d, want 2", len(list))
	}
}

func TestManagerRemove(t *testing.T) {
	mgr := NewManager()
	dir := t.TempDir()

	_, _ = mgr.Add(dir)
	mgr.Remove(dir)

	_, exists := mgr.Get(dir)
	if exists {
		t.Error("Get() should not find removed project")
	}
}

package project

import (
	"fmt"
	"os"
	"path/filepath"
)

type Project struct {
	ID   int
	Name string
	Path string
	Type string
}

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) DetectType(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to access path %s: %w", path, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", path)
	}

	if d.hasFile(path, "artisan") {
		return "laravel", nil
	}
	if d.hasFile(path, "composer.json") {
		return "php", nil
	}
	if d.hasFile(path, "package.json") {
		return "node", nil
	}
	if d.hasFile(path, "manage.py") {
		return "python", nil
	}
	if d.hasFile(path, "go.mod") {
		return "go", nil
	}

	return "generic", nil
}

func (d *Detector) hasFile(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func (d *Detector) DetectName(path string) string {
	return filepath.Base(path)
}

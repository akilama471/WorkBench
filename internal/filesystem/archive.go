package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip archive %s: %w", src, err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory %s: %w", dest, err)
	}

	for _, f := range r.File {
		if err := validateZipPath(f.Name); err != nil {
			return fmt.Errorf("invalid path in archive: %w", err)
		}

		target := filepath.Join(dest, filepath.Clean(f.Name))
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(dest)) {
			return fmt.Errorf("archive path traversal detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", target, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", target, err)
		}

		if err := extractZipFile(f, target); err != nil {
			return fmt.Errorf("failed to extract %s: %w", f.Name, err)
		}
	}

	return nil
}

func extractZipFile(f *zip.File, target string) error {
	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func validateZipPath(name string) error {
	if strings.Contains(name, "..") {
		return fmt.Errorf("path contains .. : %s", name)
	}
	if strings.HasPrefix(name, "/") || strings.HasPrefix(name, "\\") {
		return fmt.Errorf("path is absolute: %s", name)
	}
	if len(name) >= 2 && name[1] == ':' {
		return fmt.Errorf("path is absolute: %s", name)
	}
	return nil
}

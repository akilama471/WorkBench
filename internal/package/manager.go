package pkg

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
)

type Manager struct {
	paths *filesystem.Paths
	log   *logger.Logger
	repo  *Repository
}

func NewManager(paths *filesystem.Paths, log *logger.Logger) *Manager {
	return &Manager{
		paths: paths,
		log:   log,
		repo:  NewRepository(),
	}
}

func (m *Manager) Repository() *Repository {
	return m.repo
}

func (m *Manager) Download(manifest *Manifest) (string, error) {
	downloadDir := m.paths.CacheDownloads()
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create download directory: %w", err)
	}

	filename := fmt.Sprintf("%s-%s.zip", manifest.ID, manifest.Version)
	destPath := filepath.Join(downloadDir, filename)

	m.log.Info(logger.CategoryPackage, "downloading package",
		"id", manifest.ID, "version", manifest.Version, "url", manifest.Download.URL)

	resp, err := http.Get(manifest.Download.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download package %s: %w", manifest.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status %d for package %s", resp.StatusCode, manifest.ID)
	}

	outFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create download file %s: %w", destPath, err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write downloaded file %s: %w", destPath, err)
	}

	return destPath, nil
}

func (m *Manager) VerifyChecksum(filePath, expectedSHA256 string) (bool, error) {
	if expectedSHA256 == "" {
		return true, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file for checksum verification: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false, fmt.Errorf("failed to compute checksum: %w", err)
	}

	actual := fmt.Sprintf("%x", hasher.Sum(nil))
	return actual == expectedSHA256, nil
}

func (m *Manager) Extract(src, dest string) error {
	m.log.Info(logger.CategoryPackage, "extracting package", "src", src, "dest", dest)

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("failed to create extraction directory: %w", err)
	}

	ext := filepath.Ext(src)
	switch ext {
	case ".zip":
		return filesystem.ExtractZip(src, dest)
	default:
		return fmt.Errorf("unsupported archive type: %s", ext)
	}
}

func (m *Manager) Install(manifest *Manifest) error {
	downloadPath, err := m.Download(manifest)
	if err != nil {
		return err
	}

	if manifest.Download.SHA256 != "" {
		valid, err := m.VerifyChecksum(downloadPath, manifest.Download.SHA256)
		if err != nil {
			return fmt.Errorf("checksum verification error: %w", err)
		}
		if !valid {
			return fmt.Errorf("checksum mismatch for package %s", manifest.ID)
		}
	}

	extractDir := m.paths.CacheExtract()
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("failed to create extract directory: %w", err)
	}

	if err := m.Extract(downloadPath, extractDir); err != nil {
		return fmt.Errorf("failed to extract package %s: %w", manifest.ID, err)
	}

	finalDir := filepath.Join(m.paths.Bin(), manifest.Install.Directory)
	if err := os.MkdirAll(filepath.Dir(finalDir), 0o755); err != nil {
		return fmt.Errorf("failed to create final directory: %w", err)
	}

	if err := os.Rename(extractDir, finalDir); err != nil {
		return fmt.Errorf("failed to move package to final directory: %w", err)
	}

	if err := os.Remove(downloadPath); err != nil {
		m.log.Warn(logger.CategoryPackage, "failed to clean up download", "path", downloadPath)
	}

	m.log.Info(logger.CategoryPackage, "package installed", "id", manifest.ID, "version", manifest.Version)
	return nil
}

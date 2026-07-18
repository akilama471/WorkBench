package php

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	"github.com/akilama471/WorkBench/internal/runtime"
)

type Runtime struct {
	paths *filesystem.Paths
	log   *logger.Logger
}

func NewRuntime(paths *filesystem.Paths, log *logger.Logger) *Runtime {
	return &Runtime{
		paths: paths,
		log:   log,
	}
}

func (r *Runtime) ID() string   { return "php" }
func (r *Runtime) Name() string { return "PHP" }

func (r *Runtime) InstalledVersions() ([]runtime.Version, error) {
	phpBinDir := filepath.Join(r.paths.Bin(), "php")
	entries, err := os.ReadDir(phpBinDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read PHP bin directory: %w", err)
	}

	active, _ := r.activeVersionString()
	var versions []runtime.Version

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if r.isValidPHPInstallation(entry.Name()) {
			versions = append(versions, runtime.Version{
				Version:     entry.Name(),
				Path:        filepath.Join(phpBinDir, entry.Name()),
				IsActive:    entry.Name() == active,
				IsInstalled: true,
			})
		}
	}

	return versions, nil
}

func (r *Runtime) ActiveVersion() (*runtime.Version, error) {
	active, err := r.activeVersionString()
	if err != nil {
		return nil, err
	}

	if active == "" {
		return nil, nil
	}

	versions, err := r.InstalledVersions()
	if err != nil {
		return nil, err
	}

	for _, v := range versions {
		if v.Version == active {
			v.IsActive = true
			return &v, nil
		}
	}

	return nil, fmt.Errorf("active PHP version %s is not installed", active)
}

func (r *Runtime) Use(version string) error {
	r.log.Info(logger.CategoryRuntime, "switching PHP version", "target", version)

	versions, err := r.InstalledVersions()
	if err != nil {
		return fmt.Errorf("failed to list installed PHP versions: %w", err)
	}

	var target *runtime.Version
	for i := range versions {
		if versions[i].Version == version {
			target = &versions[i]
			break
		}
	}

	if target == nil {
		return fmt.Errorf("PHP version %s: %w", version, runtime.ErrVersionNotFound)
	}

	if !r.isValidPHPInstallation(version) {
		return fmt.Errorf("PHP version %s: %w", version, runtime.ErrInvalidVersion)
	}

	previousActive, _ := r.activeVersionString()

	activeFile := r.activePHPFile()
	if err := os.MkdirAll(filepath.Dir(activeFile), 0o755); err != nil {
		return fmt.Errorf("failed to create active directory: %w", err)
	}

	if err := os.WriteFile(activeFile, []byte(version), 0o644); err != nil {
		return fmt.Errorf("failed to write active PHP version: %w", err)
	}

	if err := r.verifyActiveVersion(version); err != nil {
		if previousActive != "" {
			_ = os.WriteFile(activeFile, []byte(previousActive), 0o644)
		} else {
			_ = os.Remove(activeFile)
		}
		return fmt.Errorf("PHP switch failed, rolled back: %w", err)
	}

	r.log.Info(logger.CategoryRuntime, "PHP version switched", "from", previousActive, "to", version)
	return nil
}

func (r *Runtime) activeVersionString() (string, error) {
	data, err := os.ReadFile(r.activePHPFile())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read active PHP version: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func (r *Runtime) activePHPFile() string {
	return filepath.Join(r.paths.Active(), "php")
}

func (r *Runtime) isValidPHPInstallation(version string) bool {
	binDir := filepath.Join(r.paths.Bin(), "php", version)
	entries, err := os.ReadDir(binDir)
	if err != nil {
		return false
	}

	hasPHP := false
	for _, entry := range entries {
		name := strings.ToLower(entry.Name())
		if name == "php.exe" || name == "php" {
			hasPHP = true
			break
		}
	}

	return hasPHP
}

func (r *Runtime) verifyActiveVersion(version string) error {
	data, err := os.ReadFile(r.activePHPFile())
	if err != nil {
		return fmt.Errorf("failed to read active PHP file: %w", err)
	}

	active := strings.TrimSpace(string(data))
	if active != version {
		return fmt.Errorf("active PHP version mismatch: expected %s, got %s", version, active)
	}

	if !r.isValidPHPInstallation(active) {
		return fmt.Errorf("active PHP version %s is not a valid installation", active)
	}

	return nil
}

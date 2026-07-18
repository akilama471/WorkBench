package php

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	"github.com/akilama471/WorkBench/internal/runtime"
)

type Manager struct {
	paths    *filesystem.Paths
	log      *logger.Logger
	php      *Runtime
	versions []runtime.Version
}

func NewManager(paths *filesystem.Paths, log *logger.Logger) *Manager {
	return &Manager{
		paths: paths,
		log:   log,
		php:   NewRuntime(paths, log),
	}
}

func (m *Manager) Runtime() *Runtime {
	return m.php
}

func (m *Manager) ListVersions() ([]runtime.Version, error) {
	versions, err := m.php.InstalledVersions()
	if err != nil {
		return nil, err
	}

	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i].Version, versions[j].Version) < 0
	})

	m.versions = versions
	return versions, nil
}

func (m *Manager) CurrentVersion() (string, error) {
	v, err := m.php.ActiveVersion()
	if err != nil {
		return "", err
	}
	if v == nil {
		return "none", nil
	}
	return v.Version, nil
}

func (m *Manager) SwitchVersion(version string) error {
	return m.php.Use(version)
}

func (m *Manager) PHPBinaryPath(version string) string {
	return filepath.Join(m.paths.PHPBin(version), "php.exe")
}

func (m *Manager) ActivePHPBinaryPath() (string, error) {
	v, err := m.php.ActiveVersion()
	if err != nil {
		return "", err
	}
	if v == nil {
		return "", fmt.Errorf("no active PHP version")
	}
	return m.PHPBinaryPath(v.Version), nil
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			fmt.Sscanf(aParts[i], "%d", &aNum)
		}
		if i < len(bParts) {
			fmt.Sscanf(bParts[i], "%d", &bNum)
		}
		if aNum != bNum {
			return aNum - bNum
		}
	}
	return 0
}

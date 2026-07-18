package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Load(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", path, err)
	}
	return data, nil
}

func (m *Manager) Save(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to save config to %s: %w", path, err)
	}

	return nil
}

func (m *Manager) Backup(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read file for backup %s: %w", path, err)
	}

	backupDir := filepath.Join(filepath.Dir(path), ".backups")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, filepath.Base(path)+"."+timestamp+".bak")

	if err := os.WriteFile(backupPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write backup to %s: %w", backupPath, err)
	}

	return nil
}

func (m *Manager) Replace(path string, data []byte) error {
	if err := m.Backup(path); err != nil {
		return fmt.Errorf("backup failed before replacing %s: %w", path, err)
	}

	if err := m.Save(path, data); err != nil {
		return fmt.Errorf("failed to replace config %s: %w", path, err)
	}

	return nil
}

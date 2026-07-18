package config

import (
	"fmt"
	"strings"
)

type Writer struct {
	manager *Manager
}

func NewWriter(manager *Manager) *Writer {
	return &Writer{manager: manager}
}

func (w *Writer) ReplaceLine(path string, oldDirective string, newDirective string) error {
	data, err := w.manager.Load(path)
	if err != nil {
		return err
	}

	if err := w.manager.Backup(path); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	found := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, oldDirective) {
			lines[i] = newDirective
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, newDirective)
	}

	newData := strings.Join(lines, "\n")
	if err := w.manager.Save(path, []byte(newData)); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

func (w *Writer) AddDirective(path string, directive string) error {
	data, err := w.manager.Load(path)
	if err != nil {
		return err
	}

	if err := w.manager.Backup(path); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	content := string(data)
	if strings.Contains(content, directive) {
		return nil
	}

	content = content + "\n" + directive + "\n"
	if err := w.manager.Save(path, []byte(content)); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

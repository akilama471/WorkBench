package config

import (
	"fmt"
	"strings"
)

type Loader struct {
	manager *Manager
}

func NewLoader(manager *Manager) *Loader {
	return &Loader{manager: manager}
}

func (l *Loader) LoadConfig(path string) (map[string]string, error) {
	data, err := l.manager.Load(path)
	if err != nil {
		return nil, err
	}

	return parseConfig(data), nil
}

func (l *Loader) LoadRaw(path string) (string, error) {
	data, err := l.manager.Load(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseConfig(data []byte) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

func (l *Loader) ValidateApacheConfig(path string) error {
	data, err := l.manager.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load Apache config: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, "Listen") && !strings.Contains(content, "ServerRoot") {
		return fmt.Errorf("Apache config appears invalid: missing common directives")
	}

	return nil
}

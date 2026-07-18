package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
)

type Manifest struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Version      string           `json:"version"`
	Platform     string           `json:"platform"`
	Architecture string           `json:"architecture"`
	Type         string           `json:"type"`
	Download     ManifestDownload `json:"download"`
	Install      ManifestInstall  `json:"install"`
}

type ManifestDownload struct {
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

type ManifestInstall struct {
	Directory string `json:"directory"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %s: %w", path, err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %s: %w", path, err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("manifest validation failed: %w", err)
	}

	return &manifest, nil
}

func (m *Manifest) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("manifest ID is required")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest version is required")
	}
	if m.Download.URL == "" {
		return fmt.Errorf("download URL is required")
	}
	if m.Platform == "" {
		return fmt.Errorf("platform is required")
	}
	if m.Architecture == "" {
		return fmt.Errorf("architecture is required")
	}
	return nil
}

func (m *Manifest) IsPlatformSupported() bool {
	currentPlatform := runtime.GOOS
	return m.Platform == currentPlatform
}

func (m *Manifest) IsArchitectureSupported() bool {
	currentArch := runtime.GOARCH
	return m.Architecture == currentArch
}

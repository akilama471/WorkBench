package models

type PackageInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	IsInstalled bool   `json:"is_installed"`
	InstallPath string `json:"install_path"`
}

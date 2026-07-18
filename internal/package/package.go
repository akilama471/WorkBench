package pkg

type Package struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Architecture string `json:"architecture"`
	Type         string `json:"type"`

	DownloadURL string `json:"download_url"`
	Checksum    string `json:"checksum"`
	ArchiveType string `json:"archive_type"`

	InstallDir string `json:"install_dir"`
}

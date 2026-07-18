package models

type RuntimeInfo struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ActiveVersion string   `json:"active_version"`
	Versions      []string `json:"versions"`
}

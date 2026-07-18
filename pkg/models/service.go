package models

type ServiceInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	IsInstalled bool   `json:"is_installed"`
}

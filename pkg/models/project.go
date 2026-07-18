package models

type ProjectInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

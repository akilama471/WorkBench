//go:build cgo

package state

import "github.com/akilama471/WorkBench/internal/service"

type ServiceState struct {
	ID          string
	Name        string
	Status      service.Status
	IsInstalled bool
}

type AppState struct {
	Services   map[string]*ServiceState
	PHPVersion string
}

func NewAppState() *AppState {
	return &AppState{
		Services: make(map[string]*ServiceState),
	}
}

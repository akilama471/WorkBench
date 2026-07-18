package service

import "errors"

var (
	ErrServiceNotInstalled      = errors.New("service is not installed")
	ErrServiceAlreadyRunning    = errors.New("service is already running")
	ErrServiceNotRunning        = errors.New("service is not running")
	ErrServiceNotFound          = errors.New("service not found")
	ErrServiceAlreadyRegistered = errors.New("service is already registered")
)

type Status int

const (
	StatusUnknown Status = iota
	StatusStopped
	StatusStarting
	StatusRunning
	StatusStopping
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "Unknown"
	case StatusStopped:
		return "Stopped"
	case StatusStarting:
		return "Starting"
	case StatusRunning:
		return "Running"
	case StatusStopping:
		return "Stopping"
	case StatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

type Service interface {
	ID() string
	Name() string

	Start() error
	Stop() error
	Restart() error

	Status() Status
	IsInstalled() bool
}

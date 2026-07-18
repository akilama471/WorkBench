package runtime

import "errors"

var (
	ErrRuntimeNotInstalled      = errors.New("runtime is not installed")
	ErrVersionNotFound          = errors.New("version not found")
	ErrInvalidVersion           = errors.New("invalid version")
	ErrSwitchFailed             = errors.New("failed to switch runtime version")
	ErrRuntimeAlreadyRegistered = errors.New("runtime is already registered")
	ErrRuntimeNotFound          = errors.New("runtime not found")
)

type Version struct {
	Version     string
	Path        string
	IsActive    bool
	IsInstalled bool
}

type Runtime interface {
	ID() string
	Name() string

	InstalledVersions() ([]Version, error)
	ActiveVersion() (*Version, error)
	Use(version string) error
}

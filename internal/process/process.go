package process

type Status int

const (
	StatusUnknown Status = iota
	StatusStopped
	StatusRunning
	StatusError
)

type StartConfig struct {
	Executable  string
	Args        []string
	Directory   string
	Environment []string
}

type Process struct {
	PID        int
	Executable string
	Status     Status
}

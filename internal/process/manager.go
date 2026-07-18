package process

type Manager interface {
	Start(config StartConfig) (*Process, error)
	Stop(pid int) error
	IsRunning(pid int) bool
}

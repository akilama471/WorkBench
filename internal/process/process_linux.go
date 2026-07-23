//go:build linux

package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type manager struct{}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) Start(config StartConfig) (*Process, error) {
	if config.Executable == "" {
		return nil, fmt.Errorf("executable path is required")
	}

	cmd := exec.Command(config.Executable, config.Args...)
	cmd.Dir = config.Directory
	cmd.Env = config.Environment

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process %s: %w", config.Executable, err)
	}

	proc := &Process{
		PID:        cmd.Process.Pid,
		Executable: config.Executable,
		Status:     StatusRunning,
	}

	return proc, nil
}

func (m *manager) Stop(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}

	if err := p.Signal(syscall.SIGTERM); err != nil {
		if err := p.Kill(); err != nil {
			return fmt.Errorf("failed to stop process %d: %w", pid, err)
		}
	}

	time.Sleep(500 * time.Millisecond)

	if m.IsRunning(pid) {
		if err := p.Kill(); err != nil {
			return fmt.Errorf("failed to force kill process %d: %w", pid, err)
		}
	}

	return nil
}

func (m *manager) IsRunning(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = p.Signal(syscall.Signal(0))
	return err == nil
}

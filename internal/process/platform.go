//go:build windows

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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x00000008, // DETACHED_PROCESS
	}

	// Stderr pipe removed because we don't read it and it causes deadlocks if unread

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
	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	handle, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		fmt.Printf("IsRunning OpenProcess failed for pid %d: %v\n", pid, err)
		return false
	}
	defer syscall.CloseHandle(handle)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		fmt.Printf("IsRunning GetExitCodeProcess failed for pid %d: %v\n", pid, err)
		return false
	}

	const STILL_ACTIVE = 259
	isRunning := exitCode == STILL_ACTIVE
	if !isRunning {
		fmt.Printf("IsRunning returning false for pid %d because exitCode is %d\n", pid, exitCode)
	}
	return isRunning
}

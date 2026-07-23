package process

import (
	"testing"
)

func TestStartConfig(t *testing.T) {
	config := StartConfig{
		Executable: "/usr/bin/test",
		Args:       []string{"-flag", "value"},
		Directory:  "/tmp",
		Environment: []string{
			"KEY=value",
		},
	}

	if config.Executable != "/usr/bin/test" {
		t.Errorf("Executable = %q", config.Executable)
	}
	if len(config.Args) != 2 {
		t.Errorf("Args length = %d, want 2", len(config.Args))
	}
	if config.Directory != "/tmp" {
		t.Errorf("Directory = %q", config.Directory)
	}
	if len(config.Environment) != 1 {
		t.Errorf("Environment length = %d, want 1", len(config.Environment))
	}
}

func TestProcessModel(t *testing.T) {
	proc := &Process{
		PID:        1234,
		Executable: "/usr/bin/test",
		Status:     StatusRunning,
	}

	if proc.PID != 1234 {
		t.Errorf("PID = %d, want 1234", proc.PID)
	}
	if proc.Executable != "/usr/bin/test" {
		t.Errorf("Executable = %q", proc.Executable)
	}
	if proc.Status != StatusRunning {
		t.Errorf("Status = %v, want StatusRunning", proc.Status)
	}
}

func TestProcessStatus(t *testing.T) {
	tests := []struct {
		status Status
		name   string
	}{
		{StatusUnknown, "Unknown"},
		{StatusStopped, "Stopped"},
		{StatusRunning, "Running"},
		{StatusError, "Error"},
	}

	for _, tt := range tests {
		if int(tt.status) < 0 || int(tt.status) > 3 {
			t.Errorf("Status %d out of expected range", tt.status)
		}
		_ = tt.name
	}
}

func TestNewManager(t *testing.T) {
	mgr := NewManager()
	if mgr == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestStartWithEmptyExecutable(t *testing.T) {
	mgr := NewManager()
	_, err := mgr.Start(StartConfig{})
	if err == nil {
		t.Error("expected error for empty executable")
	}
}

func TestIsRunningWithInvalidPID(t *testing.T) {
	mgr := NewManager()
	if mgr.IsRunning(99999999) {
		t.Error("IsRunning should return false for nonexistent PID")
	}
}

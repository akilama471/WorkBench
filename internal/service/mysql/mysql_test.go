package mysql

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/akilama471/WorkBench/internal/filesystem"
	"github.com/akilama471/WorkBench/internal/logger"
	"github.com/akilama471/WorkBench/internal/process"
	"github.com/akilama471/WorkBench/internal/service"
)

type mockProcessManager struct {
	started  bool
	stopped  bool
	running  bool
	startErr error
	stopErr  error
}

func (m *mockProcessManager) Start(config process.StartConfig) (*process.Process, error) {
	if m.startErr != nil {
		return nil, m.startErr
	}
	m.started = true
	return &process.Process{
		PID:        12345,
		Executable: config.Executable,
		Status:     process.StatusRunning,
	}, nil
}

func (m *mockProcessManager) Stop(pid int) error {
	if m.stopErr != nil {
		return m.stopErr
	}
	m.stopped = true
	return nil
}

func (m *mockProcessManager) IsRunning(pid int) bool {
	return m.running
}

func newTestService(dir string, proc process.Manager) *Service {
	paths := filesystem.NewPaths(dir)
	log := logger.New(logger.LevelDebug, nil)
	return NewService(paths, proc, log)
}

func TestServiceIDAndName(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	if svc.ID() != "mysql" {
		t.Errorf("ID() = %q, want \"mysql\"", svc.ID())
	}
	if svc.Name() != "MySQL" {
		t.Errorf("Name() = %q, want \"MySQL\"", svc.Name())
	}
}

func TestServiceNotInstalled(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	if svc.IsInstalled() {
		t.Error("IsInstalled() should return false when no MySQL directory exists")
	}
}

func TestServiceInstalled(t *testing.T) {
	dir := t.TempDir()
	mysqlBin := filepath.Join(dir, "bin", "mysql", "8.0.36")
	if err := os.MkdirAll(mysqlBin, 0o755); err != nil {
		t.Fatal(err)
	}

	svc := newTestService(dir, &mockProcessManager{})
	if !svc.IsInstalled() {
		t.Error("IsInstalled() should return true when MySQL directory exists")
	}
}

func TestServiceStatusInitiallyUnknown(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	if svc.Status() != service.StatusUnknown {
		t.Errorf("Status() = %v, want StatusUnknown", svc.Status())
	}
}

func TestStartNotInstalled(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	err := svc.Start()
	if err == nil {
		t.Fatal("expected error starting MySQL when not installed")
	}
	if !errors.Is(err, service.ErrServiceNotInstalled) {
		t.Errorf("expected ErrServiceNotInstalled, got %v", err)
	}
}

func TestStartSuccess(t *testing.T) {
	dir := t.TempDir()
	mysqlBin := filepath.Join(dir, "bin", "mysql", "8.0.36")
	if err := os.MkdirAll(mysqlBin, 0o755); err != nil {
		t.Fatal(err)
	}

	executable := filepath.Join(mysqlBin, "mysqld.exe")
	if err := os.WriteFile(executable, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}

	mockProc := &mockProcessManager{running: true}
	svc := newTestService(dir, mockProc)

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	if !mockProc.started {
		t.Error("process manager Start() was not called")
	}
	if svc.Status() != service.StatusRunning {
		t.Errorf("Status() = %v, want StatusRunning", svc.Status())
	}
}

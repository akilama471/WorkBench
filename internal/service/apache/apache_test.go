package apache

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
	if svc.ID() != "apache" {
		t.Errorf("ID() = %q, want \"apache\"", svc.ID())
	}
	if svc.Name() != "Apache" {
		t.Errorf("Name() = %q, want \"Apache\"", svc.Name())
	}
}

func TestServiceNotInstalled(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	if svc.IsInstalled() {
		t.Error("IsInstalled() should return false when no Apache directory exists")
	}
}

func TestServiceInstalled(t *testing.T) {
	dir := t.TempDir()
	apacheBin := filepath.Join(dir, "bin", "apache", "2.4.57")
	if err := os.MkdirAll(apacheBin, 0o755); err != nil {
		t.Fatal(err)
	}

	svc := newTestService(dir, &mockProcessManager{})
	if !svc.IsInstalled() {
		t.Error("IsInstalled() should return true when Apache directory exists")
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
		t.Fatal("expected error starting Apache when not installed")
	}
	if !errors.Is(err, service.ErrServiceNotInstalled) {
		t.Errorf("expected ErrServiceNotInstalled, got %v", err)
	}
}

func TestStartAlreadyRunning(t *testing.T) {
	dir := t.TempDir()
	apacheBin := filepath.Join(dir, "bin", "apache", "2.4.57")
	httpd := filepath.Join(apacheBin, "httpd.exe")
	if err := os.MkdirAll(apacheBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(httpd, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}
	apacheConf := filepath.Join(dir, "etc", "apache", "httpd.conf")
	if err := os.MkdirAll(filepath.Dir(apacheConf), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(apacheConf, []byte("Listen 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proc := &mockProcessManager{running: true}
	svc := newTestService(dir, proc)

	if err := svc.Start(); err != nil {
		t.Fatalf("first Start() failed: %v", err)
	}

	err := svc.Start()
	if err == nil {
		t.Fatal("expected error starting Apache when already running")
	}
	if !errors.Is(err, service.ErrServiceAlreadyRunning) {
		t.Errorf("expected ErrServiceAlreadyRunning, got %v", err)
	}
}

func TestStopNotRunning(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	err := svc.Stop()
	if err == nil {
		t.Fatal("expected error stopping Apache when not running")
	}
	if !errors.Is(err, service.ErrServiceNotRunning) {
		t.Errorf("expected ErrServiceNotRunning, got %v", err)
	}
}

func TestStartAndStop(t *testing.T) {
	dir := t.TempDir()
	apacheBin := filepath.Join(dir, "bin", "apache", "2.4.57")
	httpd := filepath.Join(apacheBin, "httpd.exe")
	if err := os.MkdirAll(apacheBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(httpd, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}
	apacheConf := filepath.Join(dir, "etc", "apache", "httpd.conf")
	if err := os.MkdirAll(filepath.Dir(apacheConf), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(apacheConf, []byte("Listen 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proc := &mockProcessManager{running: true}
	svc := newTestService(dir, proc)

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	if svc.Status() != service.StatusRunning {
		t.Errorf("after Start, Status() = %v, want StatusRunning", svc.Status())
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("Stop() failed: %v", err)
	}
	if svc.Status() != service.StatusStopped {
		t.Errorf("after Stop, Status() = %v, want StatusStopped", svc.Status())
	}
}

func TestRestart(t *testing.T) {
	dir := t.TempDir()
	apacheBin := filepath.Join(dir, "bin", "apache", "2.4.57")
	httpd := filepath.Join(apacheBin, "httpd.exe")
	if err := os.MkdirAll(apacheBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(httpd, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}
	apacheConf := filepath.Join(dir, "etc", "apache", "httpd.conf")
	if err := os.MkdirAll(filepath.Dir(apacheConf), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(apacheConf, []byte("Listen 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proc := &mockProcessManager{running: true}
	svc := newTestService(dir, proc)

	if err := svc.Restart(); err != nil {
		t.Fatalf("Restart() failed: %v", err)
	}
	if svc.Status() != service.StatusRunning {
		t.Errorf("after Restart, Status() = %v, want StatusRunning", svc.Status())
	}
}

func TestLogPaths(t *testing.T) {
	svc := newTestService(t.TempDir(), &mockProcessManager{})
	if svc.ErrorLogPath() == "" {
		t.Error("ErrorLogPath() should not be empty")
	}
	if svc.AccessLogPath() == "" {
		t.Error("AccessLogPath() should not be empty")
	}
}

func TestStatusDetectsProcessDeath(t *testing.T) {
	dir := t.TempDir()
	apacheBin := filepath.Join(dir, "bin", "apache", "2.4.57")
	httpd := filepath.Join(apacheBin, "httpd.exe")
	if err := os.MkdirAll(apacheBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(httpd, []byte(""), 0o755); err != nil {
		t.Fatal(err)
	}
	apacheConf := filepath.Join(dir, "etc", "apache", "httpd.conf")
	if err := os.MkdirAll(filepath.Dir(apacheConf), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(apacheConf, []byte("Listen 80\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	proc := &mockProcessManager{running: true}
	svc := newTestService(dir, proc)

	if err := svc.Start(); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	proc.running = false

	if svc.Status() != service.StatusStopped {
		t.Errorf("Status() = %v, want StatusStopped after process death", svc.Status())
	}
}

package service

import (
	"errors"
	"testing"
)

type mockService struct {
	id        string
	name      string
	status    Status
	installed bool
	startErr  error
	stopErr   error
}

func (m *mockService) ID() string        { return m.id }
func (m *mockService) Name() string      { return m.name }
func (m *mockService) Status() Status    { return m.status }
func (m *mockService) IsInstalled() bool { return m.installed }
func (m *mockService) Start() error {
	if m.startErr != nil {
		return m.startErr
	}
	m.status = StatusRunning
	return nil
}
func (m *mockService) Stop() error {
	if m.stopErr != nil {
		return m.stopErr
	}
	m.status = StatusStopped
	return nil
}
func (m *mockService) Restart() error {
	if err := m.Stop(); err != nil {
		return err
	}
	return m.Start()
}

func TestServiceManagerRegisterAndGet(t *testing.T) {
	mgr := NewManager()
	svc := &mockService{id: "apache", name: "Apache", status: StatusStopped, installed: true}

	if err := mgr.Register(svc); err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	got, err := mgr.Get("apache")
	if err != nil {
		t.Fatalf("Get(\"apache\") failed: %v", err)
	}
	if got.ID() != "apache" {
		t.Errorf("Get(\"apache\").ID() = %q, want \"apache\"", got.ID())
	}
}

func TestServiceManagerDuplicateRegister(t *testing.T) {
	mgr := NewManager()
	svc1 := &mockService{id: "apache", name: "Apache", status: StatusStopped}
	svc2 := &mockService{id: "apache", name: "Apache2", status: StatusStopped}

	if err := mgr.Register(svc1); err != nil {
		t.Fatalf("first Register() failed: %v", err)
	}

	err := mgr.Register(svc2)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
	if !errors.Is(err, ErrServiceAlreadyRegistered) {
		t.Errorf("expected ErrServiceAlreadyRegistered, got %v", err)
	}
}

func TestServiceManagerNotFound(t *testing.T) {
	mgr := NewManager()

	_, err := mgr.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}
}

func TestServiceManagerAll(t *testing.T) {
	mgr := NewManager()
	svc1 := &mockService{id: "apache", name: "Apache", status: StatusStopped}
	svc2 := &mockService{id: "mariadb", name: "MariaDB", status: StatusStopped}

	_ = mgr.Register(svc1)
	_ = mgr.Register(svc2)

	all := mgr.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d services, want 2", len(all))
	}
}

func TestServiceManagerStartStop(t *testing.T) {
	mgr := NewManager()
	svc := &mockService{id: "apache", name: "Apache", status: StatusStopped, installed: true}
	_ = mgr.Register(svc)

	if err := mgr.Start("apache"); err != nil {
		t.Fatalf("Start(\"apache\") failed: %v", err)
	}
	if svc.status != StatusRunning {
		t.Errorf("after Start, status = %v, want StatusRunning", svc.status)
	}

	if err := mgr.Stop("apache"); err != nil {
		t.Fatalf("Stop(\"apache\") failed: %v", err)
	}
	if svc.status != StatusStopped {
		t.Errorf("after Stop, status = %v, want StatusStopped", svc.status)
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusUnknown, "Unknown"},
		{StatusStopped, "Stopped"},
		{StatusStarting, "Starting"},
		{StatusRunning, "Running"},
		{StatusStopping, "Stopping"},
		{StatusError, "Error"},
		{Status(99), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.status.String(); got != tt.expected {
			t.Errorf("Status(%d).String() = %q, want %q", tt.status, got, tt.expected)
		}
	}
}

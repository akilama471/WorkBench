package service

import "fmt"

type Manager struct {
	services map[string]Service
}

func NewManager() *Manager {
	return &Manager{
		services: make(map[string]Service),
	}
}

func (m *Manager) Register(svc Service) error {
	id := svc.ID()
	if _, exists := m.services[id]; exists {
		return fmt.Errorf("%s: %w", id, ErrServiceAlreadyRegistered)
	}
	m.services[id] = svc
	return nil
}

func (m *Manager) Get(id string) (Service, error) {
	svc, exists := m.services[id]
	if !exists {
		return nil, fmt.Errorf("%s: %w", id, ErrServiceNotFound)
	}
	return svc, nil
}

func (m *Manager) All() []Service {
	result := make([]Service, 0, len(m.services))
	for _, svc := range m.services {
		result = append(result, svc)
	}
	return result
}

func (m *Manager) Start(id string) error {
	svc, err := m.Get(id)
	if err != nil {
		return err
	}
	return svc.Start()
}

func (m *Manager) Stop(id string) error {
	svc, err := m.Get(id)
	if err != nil {
		return err
	}
	return svc.Stop()
}

func (m *Manager) Restart(id string) error {
	svc, err := m.Get(id)
	if err != nil {
		return err
	}
	return svc.Restart()
}

func (m *Manager) Status(id string) (Status, error) {
	svc, err := m.Get(id)
	if err != nil {
		return StatusUnknown, err
	}
	return svc.Status(), nil
}

package runtime

import "fmt"

type Manager struct {
	runtimes map[string]Runtime
}

func NewManager() *Manager {
	return &Manager{
		runtimes: make(map[string]Runtime),
	}
}

func (m *Manager) Register(r Runtime) error {
	id := r.ID()
	if _, exists := m.runtimes[id]; exists {
		return fmt.Errorf("%s: %w", id, ErrRuntimeAlreadyRegistered)
	}
	m.runtimes[id] = r
	return nil
}

func (m *Manager) Get(id string) (Runtime, error) {
	r, exists := m.runtimes[id]
	if !exists {
		return nil, fmt.Errorf("%s: %w", id, ErrRuntimeNotFound)
	}
	return r, nil
}

func (m *Manager) All() []Runtime {
	result := make([]Runtime, 0, len(m.runtimes))
	for _, r := range m.runtimes {
		result = append(result, r)
	}
	return result
}

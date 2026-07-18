package project

type Manager struct {
	projects map[string]*Project
	detector *Detector
}

func NewManager() *Manager {
	return &Manager{
		projects: make(map[string]*Project),
		detector: NewDetector(),
	}
}

func (m *Manager) Add(path string) (*Project, error) {
	projectType, err := m.detector.DetectType(path)
	if err != nil {
		return nil, err
	}

	name := m.detector.DetectName(path)

	if _, exists := m.projects[path]; exists {
		return m.projects[path], nil
	}

	project := &Project{
		Name: name,
		Path: path,
		Type: projectType,
	}

	m.projects[path] = project
	return project, nil
}

func (m *Manager) Get(path string) (*Project, bool) {
	p, exists := m.projects[path]
	return p, exists
}

func (m *Manager) List() []*Project {
	result := make([]*Project, 0, len(m.projects))
	for _, p := range m.projects {
		result = append(result, p)
	}
	return result
}

func (m *Manager) Remove(path string) {
	delete(m.projects, path)
}

func (m *Manager) Count() int {
	return len(m.projects)
}

package pkg

type Repository struct {
	manifests map[string]*Manifest
}

func NewRepository() *Repository {
	return &Repository{
		manifests: make(map[string]*Manifest),
	}
}

func (r *Repository) Register(m *Manifest) {
	key := r.key(m.ID, m.Version, m.Platform, m.Architecture)
	r.manifests[key] = m
}

func (r *Repository) Find(id, version, platform, architecture string) (*Manifest, bool) {
	key := r.key(id, version, platform, architecture)
	m, exists := r.manifests[key]
	return m, exists
}

func (r *Repository) FindByPlatform(id, version, platform, architecture string) (*Manifest, bool) {
	for _, m := range r.manifests {
		if m.ID == id && m.Version == version && m.Platform == platform && m.Architecture == architecture {
			return m, true
		}
	}
	return nil, false
}

func (r *Repository) key(id, version, platform, architecture string) string {
	return id + ":" + version + ":" + platform + ":" + architecture
}

func (r *Repository) All() []*Manifest {
	result := make([]*Manifest, 0, len(r.manifests))
	for _, m := range r.manifests {
		result = append(result, m)
	}
	return result
}

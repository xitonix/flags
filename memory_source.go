package flags

import "go.xitonix.io/flags/internal"

type MemorySource struct {
	cache map[string]string
}

func NewMemorySource() *MemorySource {
	return &MemorySource{
		cache: make(map[string]string),
	}
}

func (m *MemorySource) Read(key string) (string, bool) {
	value, ok := m.cache[key]
	return value, ok
}

func (m *MemorySource) Add(key, value string) {
	if internal.IsEmpty(key) {
		return
	}
	m.cache[key] = value
}

func (m *MemorySource) AddRange(items map[string]string) {
	for key, value := range items {
		m.Add(key, value)
	}
}

package flags

import "go.xitonix.io/flags/internal"

// MemorySource represents an in-memory implementation of `core.Source` interface.
type MemorySource struct {
	cache map[string]string
}

// NewMemorySource creates a new instance of MemorySource type.
func NewMemorySource() *MemorySource {
	return &MemorySource{
		cache: make(map[string]string),
	}
}

// Read reads the in-memory value associated with the specified key.
func (m *MemorySource) Read(key string) (string, bool) {
	value, ok := m.cache[key]
	return value, ok
}

// Add adds a new key-value pair to the source.
//
// If the key already exists, it will override the current value.
func (m *MemorySource) Add(key, value string) {
	if internal.IsEmpty(key) {
		return
	}
	m.cache[key] = value
}

// AddRange adds a list of key-value pairs to the source.
//
// If any of the keys already exists, its value will override the current value.
func (m *MemorySource) AddRange(items map[string]string) {
	for key, value := range items {
		m.Add(key, value)
	}
}

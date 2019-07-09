package mocks

// EnvReader represents a mocked environment variable reader.
//
// This will read from a simulated in-memory variable registry.
type EnvReader struct {
	register map[string]string
}

// NewEnvReader creates a new environment variable reader mock.
func NewEnvReader() *EnvReader {
	return &EnvReader{
		register: make(map[string]string),
	}
}

// Get reads the environment variable from the simulated in-memory variable registry.
func (e *EnvReader) Get(key string) (string, bool) {
	v, ok := e.register[key]
	return v, ok
}

// Set sets a new key-value in the simulated in-memory variable registry
func (e *EnvReader) Set(key, value string) {
	e.register[key] = value
}

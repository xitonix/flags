package mocks

type EnvReader struct {
	register map[string]string
}

func NewEnvReader() *EnvReader {
	return &EnvReader{
		register: make(map[string]string),
	}
}

func (e *EnvReader) Get(key string) (string, bool) {
	v, ok := e.register[key]
	return v, ok
}

func (e *EnvReader) Set(key, value string) {
	e.register[key] = value
}

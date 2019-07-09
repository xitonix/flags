package internal

import "os"

// OSEnvReader environment variable reader.
type OSEnvReader struct{}

// Get reads the value of an environment variable.
func (OSEnvReader) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}

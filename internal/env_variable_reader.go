package internal

// EnvironmentVariableReader environment variable reader interface.
type EnvironmentVariableReader interface {
	Get(key string) (string, bool)
}

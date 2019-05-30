package internal

type EnvironmentVariableReader interface {
	Get(key string) (string, bool)
}

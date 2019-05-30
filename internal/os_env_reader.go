package internal

import "os"

type OSEnvReader struct{}

func (OSEnvReader) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}

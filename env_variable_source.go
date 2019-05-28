package flags

import "os"

type envVariableSource struct {
}

func (*envVariableSource) Read(key string) (string, bool) {
	return os.LookupEnv(key)
}

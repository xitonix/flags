package flags

import (
	"go.xitonix.io/flags/internal"
)

type envVariableSource struct {
	set internal.EnvironmentVariableReader
}

func newEnvironmentVarSource(set internal.EnvironmentVariableReader) *envVariableSource {
	return &envVariableSource{
		set: set,
	}
}

func (e *envVariableSource) Read(key string) (string, bool) {
	return e.set.Get(key)
}

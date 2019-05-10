package core

import (
	"go.xitonix.io/flags/internal"
)

type EnvironmentVariable struct {
	prefix, key string
}

func (v *EnvironmentVariable) Name() string {
	if internal.IsEmpty(v.key) {
		return ""
	}
	if internal.IsEmpty(v.prefix) {
		return v.key
	}
	return v.prefix + "_" + v.key
}

func (v *EnvironmentVariable) Auto(longName string) {
	if internal.IsEmpty(v.key) {
		v.Set(longName)
	}
}

func (v *EnvironmentVariable) SetPrefix(prefix string) {
	v.prefix = internal.SanitiseEnvVarName(prefix)
}

func (v *EnvironmentVariable) Set(key string) {
	v.key = internal.SanitiseEnvVarName(key)
}

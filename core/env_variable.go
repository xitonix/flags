package core

import (
	"go.xitonix.io/flags/internal"
)

// EnvironmentVariable represents the definition of an environment variable.
//
// The name of an environment variable is case sensitive and all uppercase.
// Different segments of the name is concatenated using an "_" character.
// For example PREFIX_PORT_NUMBER
type EnvironmentVariable struct {
	prefix, key string
	isSet       bool
}

// Prefix returns the prefix of the environment variable if applicable.
func (v *EnvironmentVariable) Prefix() string {
	return v.prefix
}

// Name returns the name of the environment variable.
//
// If a prefix is set, it will be included in the return value.
func (v *EnvironmentVariable) Name() string {
	if internal.IsEmpty(v.key) {
		return ""
	}
	if internal.IsEmpty(v.prefix) {
		return v.key
	}
	return v.prefix + "_" + v.key
}

// SetPrefix sets the prefix of the environment variable's name.
func (v *EnvironmentVariable) SetPrefix(prefix string) {
	v.prefix = internal.SanitiseEnvVarName(prefix)
}

// Set sets the name of the environment variable.
//
// An automatically generated name ('auto' == true), will be overridden with
// explicit values (where the method is called with 'auto' parameter set to false).
func (v *EnvironmentVariable) Set(name string, auto bool) {
	if v.isSet && auto {
		return
	}
	v.key = internal.SanitiseEnvVarName(name)
	v.isSet = !auto
}

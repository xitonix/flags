package flags

import (
	"strings"
)

type EnvironmentVariable interface {
	Prefix() string
	Name() string
	Key() string
	auto()
	setPrefix(string)
	set(string)
}

type Variable struct {
	prefix, key string
}

func (v *Variable) Prefix() string {
	return v.prefix
}

func (v *Variable) Name() string {
	if isEmpty(v.key) {
		return ""
	}
	if isEmpty(v.prefix) {
		return v.key
	}
	return v.prefix + "_" + v.key
}

func (v *Variable) Key() string {
	return v.key
}

func (v *Variable) auto() {
	if isEmpty(v.key) {
		v.set(v.key)
	}
}

func (v *Variable) setPrefix(prefix string) {
	v.prefix = sanitise(prefix)
}

func (v *Variable) set(key string) {
	v.key = sanitise(key)
}

func sanitise(value string) string {
	return strings.ToUpper(strings.ReplaceAll(value, "-", "_"))
}

package flags

import (
	"strings"
)

type EnvVariable struct {
	prefix, key string
}

func (v *EnvVariable) Name() string {
	if isEmpty(v.key) {
		return ""
	}
	if isEmpty(v.prefix) {
		return v.key
	}
	return v.prefix + "_" + v.key
}

func (v *EnvVariable) auto(longName string) {
	if isEmpty(v.key) {
		v.set(longName)
	}
}

func (v *EnvVariable) setPrefix(prefix string) {
	v.prefix = sanitise(prefix)
}

func (v *EnvVariable) set(key string) {
	v.key = sanitise(key)
}

func sanitise(value string) string {
	return strings.ToUpper(strings.ReplaceAll(value, "-", "_"))
}

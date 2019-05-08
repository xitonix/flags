package flags

import (
	"strings"
)

type argSource struct {
	arguments map[string]string
}

func newArgSource(args []string) *argSource {
	src := &argSource{
		arguments: make(map[string]string),
	}
	if len(args) == 0 {
		return src
	}
	var prevKey string
	for _, arg := range args {
		isKey := strings.HasPrefix(arg, "-")
		trimmed := strings.TrimLeft(arg, "-")
		parts := strings.Split(trimmed, "=")
		if len(parts) >= 2 {
			// This is to support key=val as well as special cases like
			// key="-a=10 -b=20" OR key="--a=10 --b=20" to cover nested arguments
			src.arguments[parts[0]] = strings.Join(parts[1:], "=")
			prevKey = ""
			continue
		}

		if isKey {
			src.arguments[trimmed] = ""
			prevKey = trimmed
		} else {
			src.arguments[prevKey] = arg
			prevKey = ""
		}
	}
	return src
}

func (a *argSource) Read(key string) (string, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

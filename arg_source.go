package flags

import (
	"strings"
)

type argSource struct {
	arguments map[string]string
}

// creates a new command line argument parser and returns true if one of the arguments is
// --help or -h
func newArgSource(args []string) (*argSource, bool) {
	src := &argSource{
		arguments: make(map[string]string),
	}
	if len(args) == 0 {
		return src, false
	}
	var prevKey string
	var isHelpRequested bool
	for _, arg := range args {
		isKey := strings.HasPrefix(arg, "-")
		parts := strings.Split(arg, "=")
		if !isHelpRequested && isKey {
			ag := strings.TrimSpace(strings.ToLower(arg))
			if ag == "--help" || ag == "-h" || strings.HasPrefix(ag, "-h=") || strings.HasPrefix(ag, "--help=") {
				isHelpRequested = true
			}
		}
		if len(parts) >= 2 {
			// This is to support key=val as well as special cases like
			// key="-a=10 -b=20" OR key="--a=10 --b=20" to cover nested arguments
			src.arguments[parts[0]] = strings.Join(parts[1:], "=")
			prevKey = ""
			continue
		}

		if isKey {
			src.arguments[arg] = ""
			prevKey = arg
		} else {
			src.arguments[prevKey] = arg
			prevKey = ""
		}
	}
	return src, isHelpRequested
}

func (a *argSource) Read(key string) (string, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

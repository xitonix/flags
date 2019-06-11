package flags

import (
	"fmt"
	"strings"
)

type argSource struct {
	arguments map[string]string
	counters  map[string]uint64
}

// creates a new command line argument parser and returns true if one of the arguments is
// --help or -h
func newArgSource(args []string) (*argSource, bool) {
	src := &argSource{
		arguments: make(map[string]string),
		counters:  make(map[string]uint64),
	}
	if len(args) == 0 {
		return src, false
	}
	var prevKey string
	var isHelpRequested bool
	for _, arg := range args {
		isKey := strings.HasPrefix(arg, "-")
		if !isHelpRequested && isKey {
			ag := strings.TrimSpace(strings.ToLower(arg))
			if ag == "--help" || ag == "-h" || strings.HasPrefix(ag, "-h=") || strings.HasPrefix(ag, "--help=") {
				isHelpRequested = true
				continue
			}
		}
		parts := strings.Split(arg, "=")
		if len(parts) >= 2 {
			// This is to support key=val as well as special cases like
			// key="-a=10 -b=20" OR key="--a=10 --b=20" to cover nested arguments
			keys := processKey(parts[0])
			for i, k := range keys {
				src.arguments[k] = ""
				if isShort(k) {
					src.counters[k]++
				}
				if i == len(keys)-1 {
					src.arguments[k] = strings.Join(parts[1:], "=")
				}
			}
			prevKey = ""
			continue
		}

		if isKey {
			keys := processKey(arg)
			for i, k := range keys {
				src.arguments[k] = ""
				if isShort(k) {
					src.counters[k]++
				}
				if i == len(keys)-1 {
					prevKey = k
				}
			}
			continue
		}

		// The argument is not a key. It is a value
		// So the previous key (prevKey) for the next
		// argument won't be a key anymore
		src.arguments[prevKey] = arg
		prevKey = ""

	}
	return src, isHelpRequested
}

func processKey(arg string) []string {
	isShort := isShort(arg)
	if !isShort {
		return []string{arg}
	}
	arg = strings.TrimLeft(arg, "-")
	args := make([]string, len(arg))
	// expand attached short flags
	// for example ‘-abc’ is equivalent to ‘-a -b -c’
	for i, short := range arg {
		args[i] = fmt.Sprintf("-%c", short)
	}
	return args
}

func isShort(arg string) bool {
	return !strings.HasPrefix(arg, "--")
}

func (a *argSource) getCounter(short string) (uint64, bool) {
	c, ok := a.counters["-"+short]
	return c, ok
}

func (a *argSource) Read(key string) (string, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

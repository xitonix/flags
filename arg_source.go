package flags

import (
	"regexp"
	"strings"

	"github.com/xitonix/flags/core"
)

type argSource struct {
	arguments map[string]string
	repeats   map[string]int
}

type argSection struct {
	isKey bool
	value string
}

// creates a new command line argument parser and returns true if one of the arguments is
// --help or -h
func newArgSource(args []string) (*argSource, bool) {
	src := &argSource{
		arguments: make(map[string]string),
		repeats:   make(map[string]int),
	}
	if len(args) == 0 {
		return src, false
	}
	var prevKey string
	var isHelpRequested bool
	for _, arg := range args {
		number := regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+$`)
		isKey := strings.HasPrefix(arg, "-") && !number.Match([]byte(strings.TrimLeft(arg, "-")))
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
			sections := processKey(parts[0])
			for i, section := range sections {
				src.arguments[section.value] = ""
				src.repeats[section.value]++
				if i == len(sections)-1 {
					src.arguments[section.value] = strings.Join(parts[1:], "=")
				}
			}
			prevKey = ""
			continue
		}

		if isKey {
			sections := processKey(arg)
			for i, section := range sections {
				if section.isKey {
					// -short or --long key
					src.arguments[section.value] = ""
					src.repeats[section.value]++
					if i == len(sections)-1 {
						prevKey = section.value
					}
				} else {
					// short form mixed with value (i.e. -A10B2)
					// this section is a value section, not key (i.e. 10 or 2 in -A10B2)
					src.arguments[sections[i-1].value] = section.value
					prevKey = ""
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

func processKey(arg string) []argSection {
	isShort := !strings.HasPrefix(arg, "--")
	if !isShort {
		return []argSection{{
			isKey: true,
			value: arg,
		}}
	}
	arg = strings.TrimLeft(arg, "-")
	args := make([]argSection, 0)
	// expand attached short flags
	// for example ‘-abc’ is equivalent to ‘-a -b -c’
	// or mixed values such as -a100, -a10b20 or -a10b3.4
	mixedKeyValue := regexp.MustCompile(`[a-zA-Z]|[+-]?([0-9]*[.])?[0-9]+`)
	var prevArgWasKey bool
	for _, match := range mixedKeyValue.FindAllString(arg, -1) {
		var isKey bool
		val := match
		// !prevArgWasKey: if the previous item was not a key we will add the current item as a key
		// (even though it's not) in order for the invalid input to be detected as an unknown flag at the time of parsing.
		// An example is -10 is invalid because a value entry must always be preceded by a key.
		if isLetter(match) || !prevArgWasKey {
			isKey = true
			val = "-" + match
		}
		prevArgWasKey = isKey
		args = append(args, argSection{
			isKey: isKey,
			value: val,
		})
	}
	return args
}

func isLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func (a *argSource) getNumberOfRepeats(f core.Flag) int {
	return a.repeats["-"+f.ShortName()] + a.repeats["--"+f.LongName()]
}

func (a *argSource) Read(key string) (string, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

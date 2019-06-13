package test

import "strings"

func ErrorContains(err error, desired string) bool {
	if err == nil {
		return len(desired) == 0
	}
	if desired == "" {
		return false
	}
	return strings.Contains(err.Error(), desired)
}

func StringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

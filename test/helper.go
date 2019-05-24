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

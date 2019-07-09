package test

import (
	"strings"
)

// ErrorContains returns true if the specified error is not nil and contains the `desired` string.
func ErrorContains(err error, desired string) bool {
	if err == nil {
		return len(desired) == 0
	}
	if desired == "" {
		return false
	}
	return strings.Contains(err.Error(), desired)
}

// ErrorContainsExact returns true if the specified error is not nil and is equal to the `desired` string.
func ErrorContainsExact(err error, desired string) bool {
	if err == nil {
		return len(desired) == 0
	}
	if desired == "" {
		return false
	}
	return err.Error() == desired
}

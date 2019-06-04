package core

import (
	"errors"
)

var (
	// ErrEmptyFlagName occurs when a flag with an empty long name is tried to be added to a bucket.
	ErrEmptyFlagName = errors.New("the flag name cannot be empty")
)

package flags

import (
	"errors"
)

var (
	ErrEmptyFlagName = errors.New("the flag name cannot be empty")
)

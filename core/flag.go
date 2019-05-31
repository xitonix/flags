package core

import "go.xitonix.io/flags/data"

// Flag is the interface for defining a CLI flag.
type Flag interface {
	LongName() string
	ShortName() string
	Usage() string
	IsSet() bool
	Type() string
	Key() *data.Key
	Set(value string) error
	ResetToDefault()
	IsHidden() bool
	IsDeprecated() bool
	Default() interface{}
}

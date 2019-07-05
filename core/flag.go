package core

import "github.com/xitonix/flags/data"

// Flag is the interface for defining a CLI flag.
type Flag interface {
	LongName() string
	ShortName() string
	Usage() string
	IsSet() bool
	IsRequired() bool
	Type() string
	Key() *data.Key
	Set(value string) error
	ResetToDefault()
	IsHidden() bool
	IsDeprecated() bool
	Default() interface{}
}

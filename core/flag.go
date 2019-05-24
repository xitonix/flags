package core

// Flag is the interface for defining a CLI flag.
type Flag interface {
	LongName() string
	ShortName() string
	Usage() string
	IsSet() bool
	Type() string
	Env() *EnvironmentVariable
	Set(value string) error
	ResetToDefault()
	IsHidden() bool
	Default() interface{}
}

package core

type Flag interface {
	LongName() string
	ShortName() string
	Usage() string
	IsSet() bool
	Type() string
	Env() *EnvironmentVariable
	Set(value string) error
	ResetToDefault()
	Default() interface{}
}

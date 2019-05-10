package flags

type Flag interface {
	LongName() string
	ShortName() string
	Usage() string
	IsSet() bool
	Type() string
	Env() *EnvVariable
	Set(value string) error
	ResetToDefault()
	FormatHelp() string
}

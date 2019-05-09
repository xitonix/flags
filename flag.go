package flags

type Flag interface {
	Env() EnvironmentVariable
	Name() string
	Short() string
	Usage() string
	IsSet() bool
	ResetToDefault()
	Type() string
}

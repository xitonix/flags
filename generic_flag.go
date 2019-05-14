package flags

import (
	"go.xitonix.io/flags/core"
)

const defaultFlagValue = "default"

type genericFlag struct {
	long, short string
	value       string
	isSet       bool
	env         *core.EnvironmentVariable
}

func newGeneric(long, short string) *genericFlag {
	return &genericFlag{
		long:  long,
		short: short,
		env:   &core.EnvironmentVariable{},
	}
}

func (f *genericFlag) LongName() string {
	return f.long
}

func (f *genericFlag) ShortName() string {
	return f.short
}

func (f *genericFlag) Usage() string {
	return "This is a very useful flag"
}

func (f *genericFlag) IsSet() bool {
	return f.isSet
}

func (f *genericFlag) Type() string {
	return "generic"
}

func (f *genericFlag) Env() *core.EnvironmentVariable {
	return f.env
}

func (f *genericFlag) Set(value string) error {
	f.isSet = true
	f.value = value
	return nil
}

func (f *genericFlag) ResetToDefault() {
	f.value = defaultFlagValue
}

func (f *genericFlag) Default() interface{} {
	return defaultFlagValue
}

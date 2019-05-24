package flags

import (
	"go.xitonix.io/flags/core"
)

const defaultFlagValue = "default"

type flagMock struct {
	long, short string
	value       string
	isSet       bool
	env         *core.EnvironmentVariable
}

func newMockedFlag(long, short string) *flagMock {
	return &flagMock{
		long:  long,
		short: short,
		env:   &core.EnvironmentVariable{},
	}
}

func (f *flagMock) LongName() string {
	return f.long
}

func (f *flagMock) ShortName() string {
	return f.short
}

func (f *flagMock) Usage() string {
	return "This is a very useful flag"
}

func (f *flagMock) IsSet() bool {
	return f.isSet
}

func (f *flagMock) Type() string {
	return "generic"
}

func (f *flagMock) Env() *core.EnvironmentVariable {
	return f.env
}

func (f *flagMock) Set(value string) error {
	f.isSet = true
	f.value = value
	return nil
}

func (f *flagMock) ResetToDefault() {
	f.value = defaultFlagValue
}

func (f *flagMock) Default() interface{} {
	return defaultFlagValue
}

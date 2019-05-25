package flags

import (
	"go.xitonix.io/flags/core"
)

const defaultFlagValue = "default"

type flagMock struct {
	long, short  string
	value        string
	isSet        bool
	isDeprecated bool
	isHidden     bool
	key          *core.Key
}

func newMockedFlag(long, short string) *flagMock {
	return &flagMock{
		long:  long,
		short: short,
		key:   &core.Key{},
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

func (f *flagMock) IsHidden() bool {
	return f.isHidden
}

func (f *flagMock) IsDeprecated() bool {
	return f.isDeprecated
}

func (f *flagMock) Type() string {
	return "generic"
}

func (f *flagMock) Key() *core.Key {
	return f.key
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

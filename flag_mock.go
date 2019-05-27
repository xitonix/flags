package flags

import (
	"errors"
	"go.xitonix.io/flags/core"
)

type flagMock struct {
	long, short   string
	value         string
	isSet         bool
	isDeprecated  bool
	isHidden      bool
	defaultValue  string
	makeSetToFail bool
	key           *core.Key
	usage         string
}

func newMockedFlag(long, short string) *flagMock {
	return newMockedFlagWithUsage(long, short, "this is a mocked flag")
}

func newMockedFlagWithKey(long, short, key string) *flagMock {
	k := &core.Key{}
	k.Set(key)
	return &flagMock{
		long:  long,
		short: short,
		key:   k,
	}
}

func newMockedFlagWithUsage(long, short, usage string) *flagMock {
	return &flagMock{
		long:  long,
		short: short,
		key:   &core.Key{},
		usage: usage,
	}
}

func (f *flagMock) LongName() string {
	return f.long
}

func (f *flagMock) ShortName() string {
	return f.short
}

func (f *flagMock) Usage() string {
	return f.usage
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
	if f.makeSetToFail {
		return errors.New("you asked for it")
	}
	f.isSet = true
	f.value = value
	return nil
}

func (f *flagMock) ResetToDefault() {
	f.value = f.defaultValue
}

func (f *flagMock) SetDefaultValue(defaultValue string) {
	f.defaultValue = defaultValue
}

func (f *flagMock) Default() interface{} {
	return f.defaultValue
}

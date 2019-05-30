package mocks

import (
	"go.xitonix.io/flags/core"
)

type Flag struct {
	long, short   string
	value         interface{}
	isSet         bool
	isDeprecated  bool
	isHidden      bool
	defaultValue  string
	MakeSetToFail bool
	key           *core.Key
	usage         string
}

func NewFlag(long, short string) *Flag {
	return NewFlagWithUsage(long, short, "this is a mocked flag")
}

func NewFlagWithKey(long, short, key string) *Flag {
	k := &core.Key{}
	k.Set(key)
	return &Flag{
		long:  long,
		short: short,
		key:   k,
	}
}

func NewFlagWithUsage(long, short, usage string) *Flag {
	return &Flag{
		long:  long,
		short: short,
		key:   &core.Key{},
		usage: usage,
	}
}

func (f *Flag) LongName() string {
	return f.long
}

func (f *Flag) ShortName() string {
	return f.short
}

func (f *Flag) Usage() string {
	return f.usage
}

func (f *Flag) IsSet() bool {
	return f.isSet
}

func (f *Flag) IsHidden() bool {
	return f.isHidden
}

func (f *Flag) IsDeprecated() bool {
	return f.isDeprecated
}

func (f *Flag) Type() string {
	return "generic"
}

func (f *Flag) Key() *core.Key {
	return f.key
}

func (f *Flag) Set(value string) error {
	if f.MakeSetToFail {
		return ErrExpected
	}
	f.isSet = true
	f.value = value
	return nil
}

func (f *Flag) ResetToDefault() {
	f.value = f.defaultValue
}

func (f *Flag) SetDefaultValue(defaultValue string) {
	f.defaultValue = defaultValue
}

func (f *Flag) Default() interface{} {
	return f.defaultValue
}

func (f *Flag) Get() interface{} {
	return f.value
}

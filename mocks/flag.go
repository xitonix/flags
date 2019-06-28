package mocks

import (
	"go.xitonix.io/flags/data"
)

type Flag struct {
	long, short   string
	value         interface{}
	isSet         bool
	isDeprecated  bool
	isRequired    bool
	isHidden      bool
	hasDefault    bool
	defaultValue  string
	MakeSetToFail bool
	key           *data.Key
	usage         string
}

func NewFlag(long, short string) *Flag {
	return NewFlagWithUsage(long, short, "this is a mocked flag")
}

func NewFlagWithKey(long, short, key string) *Flag {
	k := &data.Key{}
	k.SetID(key)
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
		key:   &data.Key{},
		usage: usage,
	}
}

func (f *Flag) WithKey(keyID string) *Flag {
	f.key.SetID(keyID)
	return f
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

// IsRequired returns true if the flag value must be provided.
func (f *Flag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *Flag) Required() *Flag {
	f.isRequired = true
	return f
}

func (f *Flag) Type() string {
	return "generic"
}

func (f *Flag) Key() *data.Key {
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
	f.hasDefault = true
}

func (f *Flag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	return f.defaultValue
}

func (f *Flag) SetDeprecated(deprecated bool) {
	f.isDeprecated = deprecated
}

func (f *Flag) SetHidden(hidden bool) {
	f.isHidden = hidden
}

func (f *Flag) Get() interface{} {
	return f.value
}

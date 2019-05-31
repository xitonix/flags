package flags

import (
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

type StringFlag struct {
	key                 *data.Key
	defaultValue, value string
	hasDefault          bool
	ptr                 *string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
}

func newString(name, usage, short string) *StringFlag {
	ptr := new(string)
	return &StringFlag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   ptr,
	}
}

func (f *StringFlag) LongName() string {
	return f.long
}

func (f *StringFlag) IsHidden() bool {
	return f.isHidden
}

func (f *StringFlag) IsDeprecated() bool {
	return f.isDeprecated
}

func (f *StringFlag) Type() string {
	return "string"
}

func (f *StringFlag) ShortName() string {
	return f.short
}

func (f *StringFlag) Usage() string {
	return f.usage
}

func (f *StringFlag) IsSet() bool {
	return f.isSet
}

func (f *StringFlag) Var() *string {
	return f.ptr
}

func (f *StringFlag) Get() string {
	return f.value
}

func (f *StringFlag) WithKey(keyID string) *StringFlag {
	f.key.Set(keyID)
	return f
}

func (f *StringFlag) WithDefault(defaultValue string) *StringFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

func (f *StringFlag) Hide() *StringFlag {
	f.isHidden = true
	return f
}

func (f *StringFlag) MarkAsDeprecated() *StringFlag {
	f.isDeprecated = true
	return f
}

func (f *StringFlag) Set(value string) error {
	f.isSet = true
	f.set(value)
	return nil
}

func (f *StringFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

func (f *StringFlag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	if f.defaultValue == "" {
		return "''"
	}
	return f.defaultValue
}

func (f *StringFlag) Key() *data.Key {
	return f.key
}

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

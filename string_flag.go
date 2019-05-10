package flags

import "go.xitonix.io/flags/core"

type StringFlag struct {
	env                 *core.EnvironmentVariable
	defaultValue, value string
	hasDefault          bool
	ptr                 *string
	long, short         string
	usage               string
	isSet               bool
}

func newString(name string, usage string) *StringFlag {
	return newStringInternal(name, "", "", usage, false)
}

func newStringD(name string, defaultValue string, usage string) *StringFlag {
	return newStringInternal(name, "", defaultValue, usage, true)
}

func newStringP(name string, short string, usage string) *StringFlag {
	return newStringInternal(name, short, "", usage, false)
}

func newStringPD(name string, short string, defaultValue string, usage string) *StringFlag {
	return newStringInternal(name, short, defaultValue, usage, true)
}

func newStringInternal(name string, short string, defaultValue string, usage string, hasDefault bool) *StringFlag {
	ptr := new(string)
	return &StringFlag{env: &core.EnvironmentVariable{}, defaultValue: defaultValue, long: name, short: short, usage: usage, ptr: ptr, hasDefault: hasDefault}
}

func (f *StringFlag) Default() interface{} {
	if f.hasDefault {
		return f.defaultValue
	}
	return nil
}

func (f *StringFlag) LongName() string {
	return f.long
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

func (f *StringFlag) WithEnv(variable string) *StringFlag {
	f.env.Set(variable)
	return f
}

func (f *StringFlag) Set(value string) error {
	f.isSet = true
	f.set(value)
	return nil
}

func (f *StringFlag) ResetToDefault() {
	f.isSet = false
	f.set(f.defaultValue)
}

func (f *StringFlag) Env() *core.EnvironmentVariable {
	return f.env
}

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

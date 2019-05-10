package flags

import "fmt"

type StringFlag struct {
	env                 *EnvVariable
	defaultValue, value string
	ptr                 *string
	long, short         string
	usage               string
	isSet               bool
}

func newString(name string, defaultValue string, usage string) *StringFlag {
	return newStringP(name, "", defaultValue, usage)
}

func newStringP(name string, short string, defaultValue string, usage string) *StringFlag {
	ptr := new(string)
	return &StringFlag{env: &EnvVariable{}, defaultValue: defaultValue, long: name, short: short, usage: usage, ptr: ptr}
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
	f.env.set(variable)
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

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

func (f *StringFlag) FormatHelp() string {
	var short string
	if len(f.short) != 0 {
		short = "-" + f.short + ", "
	}
	return fmt.Sprintf("%s --%s\t%s", short, f.long, f.Type())
}

func (f *StringFlag) Env() *EnvVariable {
	return f.env
}

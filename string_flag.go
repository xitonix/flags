package flags

type StringFlag struct {
	envVariable         EnvironmentVariable
	defaultValue, value string
	ptr                 *string
	name, short         string
	usage               string
	isSet               bool
}

func newString(name string, defaultValue string, usage string) *StringFlag {
	return newStringP(name, "", defaultValue, usage)
}

func newStringP(name string, short string, defaultValue string, usage string) *StringFlag {
	ptr := new(string)
	return &StringFlag{envVariable: &Variable{}, defaultValue: defaultValue, name: name, short: short, usage: usage, ptr: ptr}
}

func (f *StringFlag) Env() EnvironmentVariable {
	return f.envVariable
}

func (f *StringFlag) Name() string {
	return f.name
}

func (f *StringFlag) Type() string {
	return "string"
}

func (f *StringFlag) Short() string {
	return f.short
}

func (f *StringFlag) Usage() string {
	return f.usage
}

func (f *StringFlag) IsSet() bool {
	return f.isSet
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

func (f *StringFlag) Var() *string {
	return f.ptr
}

func (f *StringFlag) Get() string {
	return f.value
}

func (f *StringFlag) WithEnv(variable string) *StringFlag {
	f.envVariable.set(variable)
	return f
}

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

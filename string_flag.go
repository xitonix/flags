package flags

type StringFlag struct {
	envVariable         EnvironmentVariable
	defaultValue, value string
	name, short         string
	usage               string
}

func newString(name string, defaultValue string, usage string) *StringFlag {
	return newStringP(name, "", defaultValue, usage)
}

func newStringP(name string, short string, defaultValue string, usage string) *StringFlag {
	return &StringFlag{envVariable: &Variable{}, defaultValue: defaultValue, name: name, short: short, usage: usage}
}

func (s *StringFlag) Env() EnvironmentVariable {
	return s.envVariable
}

func (s *StringFlag) Name() string {
	return s.name
}

func (s *StringFlag) Short() string {
	return s.short
}

func (s *StringFlag) Usage() string {
	return s.usage
}

func (s *StringFlag) Set(value string) {
	if isEmpty(value) {
		value = s.defaultValue
	}
	s.value = value
}

func (s *StringFlag) Get() string {
	return s.value
}

func (s *StringFlag) WithEnv(variable string) *StringFlag {
	s.envVariable.set(variable)
	return s
}

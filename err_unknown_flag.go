package flags

type ErrUnknownFlag struct {
	name string
}

func errUnknownFlag(name string) *ErrUnknownFlag {
	return &ErrUnknownFlag{
		name: name,
	}
}

func (e *ErrUnknownFlag) FieldName() string {
	return e.name
}

func (e *ErrUnknownFlag) Error() string {
	return e.name + " is an unknown flag"
}

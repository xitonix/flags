package core

// ErrUnknownFlag occurs when an undefined flag has been passed to the tool as a command line argument.
type ErrUnknownFlag struct {
	name string
}

// NewUnknownFlagErr creates a new instance of ErrUnknownFlag
func NewUnknownFlagErr(name string) *ErrUnknownFlag {
	return &ErrUnknownFlag{
		name: name,
	}
}

// Error returns the string representation of an ErrUnknownFlag.
func (e *ErrUnknownFlag) Error() string {
	return e.name + " is an unknown flag"
}

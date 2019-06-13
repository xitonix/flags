package core

// Terminator defines the terminator interface.
//
// A terminator is responsible for terminating the execution of the running tool.
// i.e. after printing help or when an error occurred.
type Terminator interface {
	Terminate(code int)
}

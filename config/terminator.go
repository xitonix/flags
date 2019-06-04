package config

import "os"

// Terminator represents an implementation of the interface `core.Terminator`
type Terminator struct{}

// Terminate calls os.Exit() to terminate the execution of the running tool.
func (Terminator) Terminate(code int) {
	os.Exit(code)
}

package core

import "os"

type Terminator interface {
	Terminate(code int)
}

type OSTerminator struct{}

func (OSTerminator) Terminate(code int) {
	os.Exit(code)
}

package defaults

import "os"

type Terminator struct{}

func (Terminator) Terminate(code int) {
	os.Exit(code)
}

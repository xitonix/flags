package mocks

// Terminator represents a mocked terminator.
type Terminator struct {
	IsTerminated bool
	Code         int
}

// Terminate simulates the termination.
func (t *Terminator) Terminate(code int) {
	t.IsTerminated = true
	t.Code = code
}

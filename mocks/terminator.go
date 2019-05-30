package mocks

type Terminator struct {
	IsTerminated bool
	Code         int
}

func (t *Terminator) Terminate(code int) {
	t.IsTerminated = true
	t.Code = code
}

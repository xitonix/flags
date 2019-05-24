package test

type TerminatorMock struct {
	IsTerminated bool
	Code         int
}

func (t *TerminatorMock) Terminate(code int) {
	t.IsTerminated = true
	t.Code = code
}

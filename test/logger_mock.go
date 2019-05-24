package test

type LoggerMock struct {
	IsPrintCalled bool
}

func (l LoggerMock) Print(string) {
	l.IsPrintCalled = true
}

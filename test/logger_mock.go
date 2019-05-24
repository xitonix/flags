package test

type LoggerMock struct {
	Error error
}

func (l *LoggerMock) Print(err error) {
	l.Error = err
}

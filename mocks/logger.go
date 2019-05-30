package mocks

type Logger struct {
	Error error
}

func (l *Logger) Print(err error) {
	l.Error = err
}

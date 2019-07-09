package mocks

// Logger represents a mocked logger.
type Logger struct {
	Error error
}

// Print simulate printing without writing anything to standard output.
func (l *Logger) Print(err error) {
	l.Error = err
}

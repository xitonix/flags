package mocks

// InMemoryWriter represents a mocked writer.
//
// This will write into an in-memory buffer.
type InMemoryWriter struct {
	IsClosed          bool
	WriteCounter      int
	Lines             []string
	ForceWriteToBreak bool
	ForceCloseToBreak bool
}

// NewInMemoryWriter creates a new in-memory writer object
func NewInMemoryWriter() *InMemoryWriter {
	return &InMemoryWriter{
		Lines: make([]string, 0),
	}
}

// Write writes the input into the in-memory buffer.
func (w *InMemoryWriter) Write(p []byte) (n int, err error) {
	if w.ForceWriteToBreak {
		return 0, ErrExpected
	}
	w.WriteCounter++
	w.Lines = append(w.Lines, string(p))
	return len(p), nil
}

// Close closes the writer
func (w *InMemoryWriter) Close() error {
	if w.ForceCloseToBreak {
		return ErrExpected
	}
	w.IsClosed = true
	return nil
}

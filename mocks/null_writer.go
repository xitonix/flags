package mocks

type InMemoryWriter struct {
	IsClosed          bool
	WriteCounter      int
	Lines             []string
	ForceWriteToBreak bool
	ForceCloseToBreak bool
}

func NewInMemoryWriter() *InMemoryWriter {
	return &InMemoryWriter{
		Lines: make([]string, 0),
	}
}

func (w *InMemoryWriter) Write(p []byte) (n int, err error) {
	if w.ForceWriteToBreak {
		return 0, ErrExpected
	}
	w.WriteCounter++
	w.Lines = append(w.Lines, string(p))
	return len(p), nil
}

func (w *InMemoryWriter) Close() error {
	if w.ForceCloseToBreak {
		return ErrExpected
	}
	w.IsClosed = true
	return nil
}

package mocks

type NullWriter struct {
	IsClosed          bool
	WriteCounter      int
	Lines             []string
	ForceWriteToBreak bool
	ForceCloseToBreak bool
}

func NewNullWriter() *NullWriter {
	return &NullWriter{
		Lines: make([]string, 0),
	}
}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	if w.ForceWriteToBreak {
		return 0, ErrExpected
	}
	w.WriteCounter++
	w.Lines = append(w.Lines, string(p))
	return len(p), nil
}

func (w *NullWriter) Close() error {
	if w.ForceCloseToBreak {
		return ErrExpected
	}
	w.IsClosed = true
	return nil
}

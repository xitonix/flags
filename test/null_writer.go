package test

type NullWriter struct {
	IsClosed     bool
	WriteCounter int
	Lines        []string
}

func NewNullWriter() *NullWriter {
	return &NullWriter{
		Lines: make([]string, 0),
	}
}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	w.WriteCounter++
	w.Lines = append(w.Lines, string(p))
	return len(p), nil
}

func (w *NullWriter) Close() error {
	w.IsClosed = true
	return nil
}

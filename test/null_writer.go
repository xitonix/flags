package test

type NullWriter struct {
	IsClosed     bool
	WriteCounter int
}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	w.WriteCounter++
	return len(p), nil
}

func (w *NullWriter) Close() error {
	w.IsClosed = true
	return nil
}

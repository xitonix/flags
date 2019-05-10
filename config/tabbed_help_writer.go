package config

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type TabbedHelpWriter struct {
	w *tabwriter.Writer
}

func NewTabbedHelpWriter() *TabbedHelpWriter {
	return &TabbedHelpWriter{
		w: tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns),
	}
}

func (h *TabbedHelpWriter) Write(p []byte) (int, error) {
	return fmt.Fprintf(h.w, string(p))
}

func (h *TabbedHelpWriter) Close() error {
	return h.w.Flush()
}

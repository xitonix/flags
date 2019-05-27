package core

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
	// Hidden flag
	if len(p) == 0 {
		return 0, nil
	}
	return fmt.Fprintf(h.w, string(p))
}

func (h *TabbedHelpWriter) Close() error {
	return h.w.Flush()
}

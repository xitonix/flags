package core

import (
	"fmt"
	"io"

	"text/tabwriter"
)

// TabbedHelpWriter represents a tabbed help writer.
type TabbedHelpWriter struct {
	w *tabwriter.Writer
}

// NewTabbedHelpWriter creates a new instance of a TabbedHelpWriter.
func NewTabbedHelpWriter(output io.Writer) *TabbedHelpWriter {
	return &TabbedHelpWriter{
		w: tabwriter.NewWriter(output, 0, 0, 2, ' ', tabwriter.DiscardEmptyColumns),
	}
}

// Write writes the formatted help lines to the specified output.
func (h *TabbedHelpWriter) Write(p []byte) (int, error) {
	// Hidden flag
	if len(p) == 0 {
		return 0, nil
	}
	return fmt.Fprintf(h.w, string(p))
}

// Close flushes the buffer.
func (h *TabbedHelpWriter) Close() error {
	return h.w.Flush()
}

package core

import (
	"io"
)

type HelpProvider struct {
	Writer    io.WriteCloser
	Formatter HelpFormatter
}

func NewHelpProvider(w io.WriteCloser, f HelpFormatter) HelpProvider {
	return HelpProvider{
		Writer:    w,
		Formatter: f,
	}
}

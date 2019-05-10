package config

import (
	"io"

	"go.xitonix.io/flags/core"
)

type HelpProvider struct {
	Writer    io.WriteCloser
	Formatter core.HelpFormatter
}

func NewHelpProvider(w io.WriteCloser, f core.HelpFormatter) HelpProvider {
	return HelpProvider{
		Writer:    w,
		Formatter: f,
	}
}

func DefaultHelpProvider() HelpProvider {
	return NewHelpProvider(NewTabbedHelpWriter(), &TabbedHelpFormatter{})
}

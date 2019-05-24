package config

import (
	"go.xitonix.io/flags/core"
)

type Options struct {
	EnvPrefix    string
	AutoEnv      bool
	Terminator   core.Terminator
	Logger       core.Logger
	HelpProvider HelpProvider
}

func NewOptions() *Options {
	return &Options{
		// Set default values here
		EnvPrefix:    "",
		AutoEnv:      false,
		Terminator:   &core.OSTerminator{},
		Logger:       &DefaultLogger{},
		HelpProvider: NewHelpProvider(NewTabbedHelpWriter(), NewTabbedHelpFormatter("(default: %v)", "[DEPRECATED]")),
	}
}

// Option represents an option function
type Option func(options *Options)

func WithHelpProvider(p HelpProvider) Option {
	return func(options *Options) {
		options.HelpProvider = p
	}
}

func WithAutoEnv() Option {
	return func(options *Options) {
		options.AutoEnv = true
	}
}

func WithTerminator(terminator core.Terminator) Option {
	return func(options *Options) {
		options.Terminator = terminator
	}
}

func WithLogger(logger core.Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

func WithEnvPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = prefix
	}
}

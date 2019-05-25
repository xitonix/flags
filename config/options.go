package config

import (
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/defaults"
)

type Options struct {
	KeyPrefix    string
	AutoKeys     bool
	Terminator   core.Terminator
	Logger       core.Logger
	HelpProvider HelpProvider
}

func NewOptions() *Options {
	return &Options{
		// SetID default values here
		KeyPrefix:  "",
		AutoKeys:   false,
		Terminator: &defaults.Terminator{},
		Logger:     &defaults.Logger{},
		HelpProvider: NewHelpProvider(NewTabbedHelpWriter(),
			NewTabbedHelpFormatter(defaults.DefaultValueFormatString, defaults.DeprecatedFlagIndicator)),
	}
}

// Option represents an option function
type Option func(options *Options)

func WithHelpProvider(p HelpProvider) Option {
	return func(options *Options) {
		options.HelpProvider = p
	}
}

func WithAutoKeys() Option {
	return func(options *Options) {
		options.AutoKeys = true
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

func WithKeyPrefix(prefix string) Option {
	return func(options *Options) {
		options.KeyPrefix = prefix
	}
}

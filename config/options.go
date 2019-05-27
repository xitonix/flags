package config

import (
	"go.xitonix.io/flags/by"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type Options struct {
	KeyPrefix    string
	AutoKeys     bool
	Terminator   core.Terminator
	Logger       core.Logger
	HelpProvider core.HelpProvider
	Comparer     by.Comparer
}

func NewOptions() *Options {
	return &Options{
		// SetID default values here
		KeyPrefix:  "",
		AutoKeys:   false,
		Comparer:   by.Declared,
		Terminator: &Terminator{},
		Logger:     &Logger{},
		HelpProvider: core.NewHelpProvider(core.NewTabbedHelpWriter(),
			core.NewTabbedHelpFormatter(DefaultValueFormatString, DeprecatedFlagIndicator)),
	}
}

// Option represents an option function
type Option func(options *Options)

func WithSort(c by.Comparer) Option {
	return func(options *Options) {
		options.Comparer = c
	}
}

func WithHelpProvider(p core.HelpProvider) Option {
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
		options.KeyPrefix = internal.SanitiseFlagID(prefix)
	}
}

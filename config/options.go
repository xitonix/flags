package config

import (
	"go.xitonix.io/flags/core"
)

type Options struct {
	EnvPrefix    string
	AutoEnv      bool
	Log          core.Logger
	HelpProvider HelpProvider
}

func NewOptions() *Options {
	return &Options{
		// Set default values here
		EnvPrefix:    "",
		AutoEnv:      false,
		Log:          &DefaultLogger{},
		HelpProvider: DefaultHelpProvider(),
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

func WithLogger(logger core.Logger) Option {
	return func(options *Options) {
		options.Log = logger
	}
}

func WithEnvPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = prefix
	}
}

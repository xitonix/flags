package config

import "io"

type Options struct {
	EnvPrefix  string
	AutoEnv    bool
	Log        Logger
	HelpWriter io.Writer
}

func NewOptions() *Options {
	return &Options{
		// Set default values here
		EnvPrefix: "",
		AutoEnv:   false,
		Log:       &DefaultLogger{},
	}
}

// Option represents an option function
type Option func(options *Options)

func WithAutoEnv() Option {
	return func(options *Options) {
		options.AutoEnv = true
	}
}

func WithLogger(logger Logger) Option {
	return func(options *Options) {
		options.Log = logger
	}
}

func WithEnvPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = prefix
	}
}

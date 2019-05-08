package flags

type Options struct {
	EnvPrefix string
	AutoEnv   bool
}

func NewOptions() *Options {
	return &Options{
		// Set default values here
		EnvPrefix: "",
		AutoEnv:   false,
	}
}

// Option represents an option function
type Option func(options *Options)

func EnableAutoEnv() Option {
	return func(options *Options) {
		options.AutoEnv = true
	}
}

func EnvPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = prefix
	}
}

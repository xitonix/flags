package flags

type Options struct {
	EnvPrefix string
	AutoEnv   bool
	Log       Logger
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

func EnableAutoEnv() Option {
	return func(options *Options) {
		options.AutoEnv = true
	}
}

func WithLogger(logger Logger) Option {
	return func(options *Options) {
		options.Log = logger
	}
}

func EnvPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvPrefix = prefix
	}
}

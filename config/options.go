package config

import (
	"io"
	"os"

	"go.xitonix.io/flags/by"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

// Options holds the configuration settings of a bucket
type Options struct {
	// KeyPrefix the prefix for the flag keys (default: "").
	KeyPrefix string
	// AutoKey enables automatic key generation for the bucket (default: false)
	//
	// This will generate a unique key for each flag within the bucket. Automatically generated keys are based on the flags'
	// long name. For example 'file-path' will generate the string 'FILE_PATH' as the key.
	AutoKeys bool
	// Terminator is the internal terminator for the bucket.
	//
	// The terminator is responsible for terminating the execution of the running tool.
	// For example, the execution will be terminated after printing help.
	// The default terminator will call os.Exit() internally.
	Terminator core.Terminator
	// Logger is the internal logger of the bucket.
	Logger core.Logger
	// HelpFormatter is the help formatter for the bucket.
	//
	// It is responsible for formatting the help output.
	// The default formatter generates tabbed output.
	HelpFormatter core.HelpFormatter
	// HelpWriter is the bucket's help writer.
	//
	// The help writer is responsible to print the formatted help.
	HelpWriter io.WriteCloser
	// Comparer is the sort order of the bucket (default: by.DeclarationOrder).
	//
	// The comparer decides the order in which the flags will be displayed in the help output.
	// By default the flags will be printed in the same order as they have been defined.
	//
	// You can use the built-in sort orders such as by.KeyAscending, by.LongNameDescending, etc to override the defaults.
	// Alternatively you can implement `by.Comparer` interface and use your own comparer to sort the help output.
	Comparer by.Comparer
	// DeprecationMark is the deprecation mark for the flags within the bucket (default: config.DeprecatedFlagIndicatorDefault).
	//
	// The deprecation mark is used in the help output to draw the users' attention.
	DeprecationMark string
	// DefaultValueFormatString is the flags' Default value format string within the bucket (default: config.DefaultValueFormatStringDefault).
	//
	// The string is used to format the default value in help output (i.e. [Default: %v])
	DefaultValueFormatString string
	// RequiredFlagMark is used to mark a required flag in the help output.
	RequiredFlagMark string
	// PreSetCallback is a callback which will be called before the flag value has been set by a source.
	PreSetCallback core.Callback
	// PostSetCallback is a callback which will be called after the flag value has been set by a source.
	PostSetCallback core.Callback
}

func NewOptions() *Options {
	return &Options{
		// SetID default values here
		KeyPrefix:                "",
		AutoKeys:                 false,
		Comparer:                 by.DeclarationOrder,
		Terminator:               &Terminator{},
		Logger:                   &Logger{},
		HelpFormatter:            &core.TabbedHelpFormatter{},
		HelpWriter:               core.NewTabbedHelpWriter(os.Stdout),
		DeprecationMark:          DeprecatedFlagIndicatorDefault,
		DefaultValueFormatString: DefaultValueFormatStringDefault,
		RequiredFlagMark:         RequiredFlagMarkDefault,
		PreSetCallback:           nil,
		PostSetCallback:          nil,
	}
}

// Option represents an option function
type Option func(options *Options)

// WithPreSetCallback sets the callback function which will be called before the flag value is being set by a source.
func WithPreSetCallback(callback core.Callback) Option {
	return func(options *Options) {
		options.PreSetCallback = callback
	}
}

// WithPostSetCallback sets the callback function which will be called after the flag value has been set by a source.
//
// The post Set callback will not get called if the Set operation fails.
func WithPostSetCallback(callback core.Callback) Option {
	return func(options *Options) {
		options.PostSetCallback = callback
	}
}

// WithSortOrder sets the sort order of the bucket.
//
// The comparer decides the order in which the flags will be displayed in the help output.
// By default the flags will be printed in the same order as they have been defined.
//
// You can use the built-in sort orders such as by.KeyAscending, by.LongNameDescending, etc to override the defaults.
// Alternatively you can implement `by.Comparer` interface and use your own comparer to sort the help output.
func WithSortOrder(c by.Comparer) Option {
	return func(options *Options) {
		options.Comparer = c
	}
}

// WithDeprecationMark sets the deprecation mark for the flags within the bucket.
//
// The deprecation mark is used in the help output to draw the users' attention.
func WithDeprecationMark(deprecationMark string) Option {
	return func(options *Options) {
		options.DeprecationMark = deprecationMark
	}
}

// WithRequiredMark sets the indicator for the required flags within the bucket.
//
// The required mark is used in the help output to draw the users' attention.
func WithRequiredMark(requiredMark string) Option {
	return func(options *Options) {
		options.RequiredFlagMark = requiredMark
	}
}

// WithDefaultValueFormatString sets the bucket's Default value format string.
//
// The string is used to format the default value in the help output (i.e. [Default: %v])
func WithDefaultValueFormatString(defaultValueFormatString string) Option {
	return func(options *Options) {
		options.DefaultValueFormatString = defaultValueFormatString
	}
}

// WithHelpFormatter sets the help formatter of the bucket.
//
// The help formatter is responsible for formatting the help output.
// The default help formatter generates tabbed output.
func WithHelpFormatter(hf core.HelpFormatter) Option {
	return func(options *Options) {
		options.HelpFormatter = hf
	}
}

// WithHelpWriter sets the help writer of the bucket.
//
// The help writer is responsible for printing the formatted help output.
// The default help writer writes tabbed output to os.Stdout.
func WithHelpWriter(hw io.WriteCloser) Option {
	return func(options *Options) {
		options.HelpWriter = hw
	}
}

// WithAutoKeys enables automatic key generation for the bucket.
//
// This will generate a unique key for each flag within the bucket. Automatically generated keys are based on the flags'
// long name. For example 'file-path' will result in 'FILE_PATH' as the key.
//
// All the keys are uppercase strings concatenated by underscore character.
func WithAutoKeys() Option {
	return func(options *Options) {
		options.AutoKeys = true
	}
}

// WithTerminator sets the internal terminator for the bucket.
//
// The terminator is responsible for terminating the execution of the running tool.
// For example, the execution will be terminated after printing help.
// The default terminator will call os.Exit() internally.
func WithTerminator(terminator core.Terminator) Option {
	return func(options *Options) {
		options.Terminator = terminator
	}
}

// WithLogger sets the internal logger of the bucket.
func WithLogger(logger core.Logger) Option {
	return func(options *Options) {
		options.Logger = logger
	}
}

// WithKeyPrefix sets the prefix for all the automatically generated (or explicitly defined) keys.
//
// For example 'file-path' with 'Prefix' will result in 'PREFIX_FILE_PATH' as the key.
func WithKeyPrefix(prefix string) Option {
	return func(options *Options) {
		options.KeyPrefix = internal.SanitiseFlagID(prefix)
	}
}

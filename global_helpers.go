package flags

import (
	"go.xitonix.io/flags/by"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

var (
	// DefaultBucket holds the default bucket instance.
	DefaultBucket = NewBucket()
)

// EnableAutoKeyGeneration enables automatic key generation for the default bucket.
//
// This will generate a unique key for each flag within the bucket. Automatically generated keys are based on the flags'
// long name. For example 'file-path' will result in 'FILE_PATH' as the key.
//
// All the keys are uppercase strings concatenated by underscore character.
func EnableAutoKeyGeneration() {
	DefaultBucket.opts.AutoKeys = true
}

// SetKeyPrefix sets the prefix for all the automatically generated or explicitly defined keys.
//
// For example 'file-path' with 'Prefix' will result in 'PREFIX_FILE_PATH' as the key.
func SetKeyPrefix(prefix string) {
	DefaultBucket.opts.KeyPrefix = internal.SanitiseFlagID(prefix)
}

// SetLogger sets the internal logger for the default bucket.
func SetLogger(logger core.Logger) {
	DefaultBucket.opts.Logger = logger
}

// SetSortOrder sets the sort order of the default bucket.
//
// It decides the order in which the flags will be displayed in the help output.
// By default the flags will be printed in the same order as they have been defined.
//
// You can use the built-in sort orders such as by.KeyAscending, by.LongNameDescending, etc to override the defaults.
// Alternatively you can implement `by.Comparer` interface and use your own comparer to sort the help output.
func SetSortOrder(comparer by.Comparer) {
	DefaultBucket.opts.Comparer = comparer
}

// SetTerminator sets the internal terminator for the default bucket.
//
// The terminator is responsible for terminating the execution of the running tool.
// For example after printing help, the execution will be terminated. The default terminator
// will call os.Exit() internally.
func SetTerminator(terminator core.Terminator) {
	DefaultBucket.opts.Terminator = terminator
}

// SetHelpProvider sets the help provider for the default bucket.
//
// The help provider is responsible for formatting and printing the help output.
// The default help provider generates tabbed output.
func SetHelpProvider(p core.HelpProvider) {
	DefaultBucket.opts.HelpProvider = p
}

// SetDeprecationMark sets the default bucket's deprecation mark.
//
// The deprecation mark is used in the help output to draw the users' attention.
func SetDeprecationMark(m string) {
	DefaultBucket.opts.DeprecationMark = m
}

// SetDefaultValueFormatString sets the default bucket's Default value format string.
//
// The string is used to format the default value in help output (ie. [Default: %v])
func SetDefaultValueFormatString(f string) {
	DefaultBucket.opts.DefaultValueFormatString = f
}

// Parse this is a shortcut for calling the default bucket's Parse method.
//
// It parses the flags and queries all the available sources in order, to fill the value of each flag.
//
// If none of the sources offers any value, the flag will be set to the specified Default value (if any).
// In case no Default value is defined, the flag will be set to the zero value of its type. For example an
// Int flag will be set to zero.
//
// The order of the default sources is Command Line Arguments > Environment Variables > [Default Value]
func Parse() {
	DefaultBucket.Parse()
}

// AppendSource appends a new source to the default bucket.
//
// With the default configuration, the order will be:
// Command Line Arguments > Environment Variables > src > [Default Value]
//
// Note that the Parse method will query the sources in order.
func AppendSource(src core.Source) {
	DefaultBucket.AppendSource(src)
}

// PrependSource prepends a new source to the default bucket.
// This is an alias for AddSource(src, 0)
//
// With the default configuration, the order will be:
// src > Command Line Arguments > Environment Variables > [Default Value]
//
// Note that the Parse method will query the sources in order.
func PrependSource(src core.Source) {
	DefaultBucket.PrependSource(src)
}

// AddSource inserts the new source into the default bucket at the specified index.
//
// If the index is <= 0 the new source will get added to the beginning of the chain. If the index is greater than the
// current number of sources, it will get be appended the end.
//
// Note that the Parse method will query the sources in order.
func AddSource(src core.Source, index int) {
	DefaultBucket.AddSource(src, index)
}

// String adds a new string flag to the default bucket.
//
// The long names will be automatically converted to lowercase by the library.
func String(longName, usage string) *StringFlag {
	return DefaultBucket.String(longName, usage)
}

// StringP adds a new string flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

// Int adds a new Int flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func Int(longName, usage string) *IntFlag {
	return DefaultBucket.Int(longName, usage)
}

// IntP adds a new Int flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func IntP(longName, usage, shortName string) *IntFlag {
	return DefaultBucket.IntP(longName, usage, shortName)
}

// Int64 adds a new Int64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func Int64(longName, usage string) *Int64Flag {
	return DefaultBucket.Int64(longName, usage)
}

// Int64P adds a new Int64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func Int64P(longName, usage, shortName string) *Int64Flag {
	return DefaultBucket.Int64P(longName, usage, shortName)
}

// Int32 adds a new Int32 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func Int32(longName, usage string) *Int32Flag {
	return DefaultBucket.Int32(longName, usage)
}

// Int32P adds a new Int32 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func Int32P(longName, usage, shortName string) *Int32Flag {
	return DefaultBucket.Int32P(longName, usage, shortName)
}

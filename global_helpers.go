package flags

import (
	"go.xitonix.io/flags/by"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
	"io"
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
//
// Note: In order for the flag values to be extractable from the environment variables (or all the other custom sources),
// each flag MUST have a key associated with it.
func EnableAutoKeyGeneration() {
	DefaultBucket.opts.AutoKeys = true
}

// SetKeyPrefix sets the prefix for all the automatically generated (or explicitly defined) keys.
//
// For example 'file-path' with 'Prefix' will result in 'PREFIX_FILE_PATH' as the key.
func SetKeyPrefix(prefix string) {
	DefaultBucket.opts.KeyPrefix = internal.SanitiseFlagID(prefix)
}

// SetLogger sets the internal logger of the default bucket.
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
// For example, the execution will be terminated after printing help.
// The default terminator will call os.Exit() internally.
func SetTerminator(terminator core.Terminator) {
	DefaultBucket.opts.Terminator = terminator
}

// SetHelpFormatter sets the help formatter of the default bucket.
//
// The help formatter is responsible for formatting the help output.
// The default help formatter generates tabbed output.
func SetHelpFormatter(hf core.HelpFormatter) {
	DefaultBucket.opts.HelpFormatter = hf
}

// SetHelpWriter sets the help writer of the default bucket.
//
// The help writer is responsible for printing the formatted help output.
// The default help writer writes tabbed output to os.Stdout.
func SetHelpWriter(hw io.WriteCloser) {
	DefaultBucket.opts.HelpWriter = hw
}

// SetDeprecationMark sets the default bucket's deprecation mark.
//
// The deprecation mark is used in the help output to draw the users' attention.
func SetDeprecationMark(m string) {
	DefaultBucket.opts.DeprecationMark = m
}

// SetDefaultValueFormatString sets the default bucket's Default value format string.
//
// The string is used to format the default value in the help output (ie. [Default: %v])
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
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func Int(longName, usage string) *IntFlag {
	return DefaultBucket.Int(longName, usage)
}

// IntP adds a new Int flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func IntP(longName, usage, shortName string) *IntFlag {
	return DefaultBucket.IntP(longName, usage, shortName)
}

// Int8 adds a new Int8 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func Int8(longName, usage string) *Int8Flag {
	return DefaultBucket.Int8(longName, usage)
}

// Int8P adds a new Int8 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func Int8P(longName, usage, shortName string) *Int8Flag {
	return DefaultBucket.Int8P(longName, usage, shortName)
}

// Int16 adds a new Int16 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func Int16(longName, usage string) *Int16Flag {
	return DefaultBucket.Int16(longName, usage)
}

// Int16P adds a new Int16 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func Int16P(longName, usage, shortName string) *Int16Flag {
	return DefaultBucket.Int16P(longName, usage, shortName)
}

// Int32 adds a new Int32 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func Int32(longName, usage string) *Int32Flag {
	return DefaultBucket.Int32(longName, usage)
}

// Int32P adds a new Int32 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func Int32P(longName, usage, shortName string) *Int32Flag {
	return DefaultBucket.Int32P(longName, usage, shortName)
}

// Int64 adds a new Int64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func Int64(longName, usage string) *Int64Flag {
	return DefaultBucket.Int64(longName, usage)
}

// Int64P adds a new Int64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func Int64P(longName, usage, shortName string) *Int64Flag {
	return DefaultBucket.Int64P(longName, usage, shortName)
}

// UInt adds a new UInt flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func UInt(longName, usage string) *UIntFlag {
	return DefaultBucket.UInt(longName, usage)
}

// UIntP adds a new UInt flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func UIntP(longName, usage, shortName string) *UIntFlag {
	return DefaultBucket.UIntP(longName, usage, shortName)
}

// UInt64 adds a new UInt64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func UInt64(longName, usage string) *UInt64Flag {
	return DefaultBucket.UInt64(longName, usage)
}

// UInt64P adds a new UInt64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func UInt64P(longName, usage, shortName string) *UInt64Flag {
	return DefaultBucket.UInt64P(longName, usage, shortName)
}

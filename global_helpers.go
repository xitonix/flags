package flags

import (
	"io"

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

// SetPreSetCallback sets the pre Set callback function for the default bucket.
//
// The function will be called before the flag value is being set by a source.
func SetPreSetCallback(callback core.Callback) {
	DefaultBucket.opts.PreSetCallback = callback
}

// SetPostSetCallback sets the post Set callback function for the default bucket.
//
// The function will be called after the flag value has been set by a source.
// The post Set callback will not get called if the Set operation fails.
func SetPostSetCallback(callback core.Callback) {
	DefaultBucket.opts.PostSetCallback = callback
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
// The string is used to format the default value in the help output (i.e. [Default: %v])
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
// Long names will be automatically converted to lowercase by the library (i.e. file-path).
// A valid short name is a case sensitive single character string (i.e. f or F).
func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

// Int adds a new Int flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Int(longName, usage string) *IntFlag {
	return DefaultBucket.Int(longName, usage)
}

// IntP adds a new Int flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func IntP(longName, usage, shortName string) *IntFlag {
	return DefaultBucket.IntP(longName, usage, shortName)
}

// Int8 adds a new Int8 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Int8(longName, usage string) *Int8Flag {
	return DefaultBucket.Int8(longName, usage)
}

// Int8P adds a new Int8 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int8P(longName, usage, shortName string) *Int8Flag {
	return DefaultBucket.Int8P(longName, usage, shortName)
}

// Int16 adds a new Int16 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Int16(longName, usage string) *Int16Flag {
	return DefaultBucket.Int16(longName, usage)
}

// Int16P adds a new Int16 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int16P(longName, usage, shortName string) *Int16Flag {
	return DefaultBucket.Int16P(longName, usage, shortName)
}

// Int32 adds a new Int32 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Int32(longName, usage string) *Int32Flag {
	return DefaultBucket.Int32(longName, usage)
}

// Int32P adds a new Int32 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int32P(longName, usage, shortName string) *Int32Flag {
	return DefaultBucket.Int32P(longName, usage, shortName)
}

// Int64 adds a new Int64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Int64(longName, usage string) *Int64Flag {
	return DefaultBucket.Int64(longName, usage)
}

// Int64P adds a new Int64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int64P(longName, usage, shortName string) *Int64Flag {
	return DefaultBucket.Int64P(longName, usage, shortName)
}

// UInt adds a new UInt flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func UInt(longName, usage string) *UIntFlag {
	return DefaultBucket.UInt(longName, usage)
}

// UIntP adds a new UInt flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UIntP(longName, usage, shortName string) *UIntFlag {
	return DefaultBucket.UIntP(longName, usage, shortName)
}

// UInt64 adds a new UInt64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func UInt64(longName, usage string) *UInt64Flag {
	return DefaultBucket.UInt64(longName, usage)
}

// UInt64P adds a new UInt64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt64P(longName, usage, shortName string) *UInt64Flag {
	return DefaultBucket.UInt64P(longName, usage, shortName)
}

// UInt32 adds a new UInt32 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func UInt32(longName, usage string) *UInt32Flag {
	return DefaultBucket.UInt32(longName, usage)
}

// UInt32P adds a new UInt32 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt32P(longName, usage, shortName string) *UInt32Flag {
	return DefaultBucket.UInt32P(longName, usage, shortName)
}

// UInt16 adds a new UInt16 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func UInt16(longName, usage string) *UInt16Flag {
	return DefaultBucket.UInt16(longName, usage)
}

// UInt16P adds a new UInt16 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt16P(longName, usage, shortName string) *UInt16Flag {
	return DefaultBucket.UInt16P(longName, usage, shortName)
}

// UInt8 adds a new UInt8 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func UInt8(longName, usage string) *UInt8Flag {
	return DefaultBucket.UInt8(longName, usage)
}

// UInt8P adds a new UInt8 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt8P(longName, usage, shortName string) *UInt8Flag {
	return DefaultBucket.UInt8P(longName, usage, shortName)
}

// Bool adds a new Bool flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func Bool(longName, usage string) *BoolFlag {
	return DefaultBucket.Bool(longName, usage)
}

// BoolP adds a new Bool flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func BoolP(longName, usage, shortName string) *BoolFlag {
	return DefaultBucket.BoolP(longName, usage, shortName)
}

// Float64 adds a new Float64 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
func Float64(longName, usage string) *Float64Flag {
	return DefaultBucket.Float64(longName, usage)
}

// Float64P adds a new Float64 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
// A valid short name is a case sensitive single character string (i.e. c or C).
func Float64P(longName, usage, shortName string) *Float64Flag {
	return DefaultBucket.Float64P(longName, usage, shortName)
}

// Float32 adds a new Float32 flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
func Float32(longName, usage string) *Float32Flag {
	return DefaultBucket.Float32(longName, usage)
}

// Float32P adds a new Float32 flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
// A valid short name is a case sensitive single character string (i.e. c or C).
func Float32P(longName, usage, shortName string) *Float32Flag {
	return DefaultBucket.Float32P(longName, usage, shortName)
}

// CounterP adds a new counter flag with a short name to the default bucket.
//
// The value of a counter flag can be increased by repeating the short or the long form of the flag.
// For example the presence of -vv command line argument will set the value of the counter to 2.
//
// Long names will be automatically converted to lowercase by the library (i.e. verbosity).
// A valid short name is a case sensitive single character string (i.e. v or V).
func CounterP(longName, usage, shortName string) *CounterFlag {
	return DefaultBucket.CounterP(longName, usage, shortName)
}

// VerbosityP is an alias for CounterP("verbose", usage, "v").
//
// The value of the verbosity flag can be increased by repeating the short or the long form.
// For example the presence of -vv command line argument will set the verbosity level to 2.
// Having '--verbose -v', '--verbose --verbose' or '-v -v' would have the same effect.
func VerbosityP(usage string) *CounterFlag {
	return DefaultBucket.CounterP("verbose", usage, "v")
}

// Duration adds a new Duration flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ttl).
//
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func Duration(longName, usage string) *DurationFlag {
	return DefaultBucket.Duration(longName, usage)
}

// DurationP adds a new Duration flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ttl).
// A valid short name is a case sensitive single character string (i.e. t or T).
//
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func DurationP(longName, usage, shortName string) *DurationFlag {
	return DefaultBucket.DurationP(longName, usage, shortName)
}

// Time adds a new Time flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. birthday).
//
// Supported layouts are:
//
// Full Date and Time
//
//  dd-MM-yyyyThh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980T14:22:20)
//  dd-MM-yyyy hh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980 14:22:20)
//  dd-MM-yyyyThh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980T02:22:20 PM)
//  dd-MM-yyyy hh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980 02:22:20 PM)
//
//  dd/MM/yyyyThh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyy hh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyyThh:mm:SS[.999999999] AM/PM
//  dd/MM/yyyy hh:mm:SS[.999999999] AM/PM
//
// Date
//
//  dd-MM-yyyy
//  dd/MM/yyyy
//
// Timestamp
//
//  MMM dd hh:mm:ss[.999999999] (24 hrs, i.e. Aug 27 14:22:20)
//  MMM dd hh:mm:ss[.999999999] AM/PM (i.e. Aug 27 02:22:20 PM)
//
// Time
//
//  hh:mm:ss[.999999999] (24 hrs, i.e. 14:22:20)
//  hh:mm:ss[.999999999] AM/PM (i.e. 02:22:20 PM)
//
// [.999999999] is the optional nano second component for time.
func Time(longName, usage string) *TimeFlag {
	return DefaultBucket.Time(longName, usage)
}

// TimeP adds a new Time flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. birthday).
// A valid short name is a case sensitive single character string (i.e. b or B).
//
// Supported layouts are:
//
// Full Date and Time
//
//  dd-MM-yyyyThh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980T14:22:20)
//  dd-MM-yyyy hh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980 14:22:20)
//  dd-MM-yyyyThh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980T02:22:20 PM)
//  dd-MM-yyyy hh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980 02:22:20 PM)
//
//  dd/MM/yyyyThh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyy hh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyyThh:mm:SS[.999999999] AM/PM
//  dd/MM/yyyy hh:mm:SS[.999999999] AM/PM
//
// Date
//
//  dd-MM-yyyy
//  dd/MM/yyyy
//
// Timestamp
//
//  MMM dd hh:mm:ss[.999999999] (24 hrs, i.e. Aug 27 14:22:20)
//  MMM dd hh:mm:ss[.999999999] AM/PM (i.e. Aug 27 02:22:20 PM)
//
// Time
//
//  hh:mm:ss[.999999999] (24 hrs, i.e. 14:22:20)
//  hh:mm:ss[.999999999] AM/PM (i.e. 02:22:20 PM)
//
// [.999999999] is the optional nano second component for time.
func TimeP(longName, usage, shortName string) *TimeFlag {
	return DefaultBucket.TimeP(longName, usage, shortName)
}

// StringSlice adds a new string slice flag to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. week-days).
//
// The value of a StringSlice flag can be set using comma (or any custom delimiter) separated strings.
// For example --week-days "Sat,Sun,Mon,Tue,Wed,Thu,Fri"
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, --weekends "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
func StringSlice(longName, usage string) *StringSliceFlag {
	return DefaultBucket.StringSlice(longName, usage)
}

// StringSliceP adds a new string slice flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. week-days).
// A valid short name is a case sensitive single character string (i.e. w or W).
//
// The value of a StringSlice flag can be set using comma (or any custom delimiter) separated strings.
// For example --week-days "Sat,Sun,Mon,Tue,Wed,Thu,Fri"
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, --weekends "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
func StringSliceP(longName, usage, shortName string) *StringSliceFlag {
	return DefaultBucket.StringSliceP(longName, usage, shortName)
}

// IntSlice adds a new string slice flag to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a IntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IntSlice(longName, usage string) *IntSliceFlag {
	return DefaultBucket.IntSlice(longName, usage)
}

// IntSliceP adds a new string slice flag with a short name to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a IntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IntSliceP(longName, usage, shortName string) *IntSliceFlag {
	return DefaultBucket.IntSliceP(longName, usage, shortName)
}

// UIntSlice adds a new string slice flag to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a UIntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func UIntSlice(longName, usage string) *UIntSliceFlag {
	return DefaultBucket.UIntSlice(longName, usage)
}

// UIntSliceP adds a new string slice flag with a short name to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a UIntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func UIntSliceP(longName, usage, shortName string) *UIntSliceFlag {
	return DefaultBucket.UIntSliceP(longName, usage, shortName)
}

// Float64Slice adds a new string slice flag to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func Float64Slice(longName, usage string) *Float64SliceFlag {
	return DefaultBucket.Float64Slice(longName, usage)
}

// Float64SliceP adds a new string slice flag with a short name to the default bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func Float64SliceP(longName, usage, shortName string) *Float64SliceFlag {
	return DefaultBucket.Float64SliceP(longName, usage, shortName)
}

// IPAddress adds a new IPAddress flag to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ip-address).
//
// The value of an IP address flag can be specified using a dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func IPAddress(longName, usage string) *IPAddressFlag {
	return DefaultBucket.IPAddress(longName, usage)
}

// IPAddressP adds a new IPAddress flag with short name to the default bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ip-address).
// A valid short name is a case sensitive single character string (i.e. i or I).
//
// The value of an IP address flag can be specified using a dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func IPAddressP(longName, usage, shortName string) *IPAddressFlag {
	return DefaultBucket.IPAddressP(longName, usage, shortName)
}

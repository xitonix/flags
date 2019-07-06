package flags

import (
	"io"

	"github.com/xitonix/flags/by"
	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/internal"
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

// SetRequiredFlagMark sets the indicator for the required flags within the default bucket.
//
// The required mark is used in the help output to draw the users' attention.
func SetRequiredFlagMark(m string) {
	DefaultBucket.opts.RequiredFlagMark = m
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

// Add adds a new custom flag type to the default bucket.
//
// This method must be called before calling Parse().
func Add(f core.Flag) {
	DefaultBucket.Add(f)
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
// The long name will be automatically converted to lowercase by the library.
func String(longName, usage string) *StringFlag {
	return DefaultBucket.String(longName, usage)
}

// StringP adds a new string flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. file-path).
// A valid short name is a case sensitive single character string (i.e. f or F).
func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

// Int adds a new int flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func Int(longName, usage string) *IntFlag {
	return DefaultBucket.Int(longName, usage)
}

// IntP adds a new int flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func IntP(longName, usage, shortName string) *IntFlag {
	return DefaultBucket.IntP(longName, usage, shortName)
}

// Int8 adds a new int8 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func Int8(longName, usage string) *Int8Flag {
	return DefaultBucket.Int8(longName, usage)
}

// Int16 adds a new int16 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func Int16(longName, usage string) *Int16Flag {
	return DefaultBucket.Int16(longName, usage)
}

// Int32 adds a new int32 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func Int32(longName, usage string) *Int32Flag {
	return DefaultBucket.Int32(longName, usage)
}

// Int32P adds a new int32 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int32P(longName, usage, shortName string) *Int32Flag {
	return DefaultBucket.Int32P(longName, usage, shortName)
}

// Int64 adds a new int64 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func Int64(longName, usage string) *Int64Flag {
	return DefaultBucket.Int64(longName, usage)
}

// Int64P adds a new int64 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func Int64P(longName, usage, shortName string) *Int64Flag {
	return DefaultBucket.Int64P(longName, usage, shortName)
}

// UInt adds a new uint flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func UInt(longName, usage string) *UIntFlag {
	return DefaultBucket.UInt(longName, usage)
}

// UIntP adds a new uint flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UIntP(longName, usage, shortName string) *UIntFlag {
	return DefaultBucket.UIntP(longName, usage, shortName)
}

// UInt64 adds a new uint64 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func UInt64(longName, usage string) *UInt64Flag {
	return DefaultBucket.UInt64(longName, usage)
}

// UInt64P adds a new uint64 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt64P(longName, usage, shortName string) *UInt64Flag {
	return DefaultBucket.UInt64P(longName, usage, shortName)
}

// UInt32 adds a new uint32 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func UInt32(longName, usage string) *UInt32Flag {
	return DefaultBucket.UInt32(longName, usage)
}

// UInt32P adds a new uint32 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt32P(longName, usage, shortName string) *UInt32Flag {
	return DefaultBucket.UInt32P(longName, usage, shortName)
}

// UInt16 adds a new uint16 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func UInt16(longName, usage string) *UInt16Flag {
	return DefaultBucket.UInt16(longName, usage)
}

// UInt16P adds a new uint16 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt16P(longName, usage, shortName string) *UInt16Flag {
	return DefaultBucket.UInt16P(longName, usage, shortName)
}

// UInt8 adds a new uint8 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
func UInt8(longName, usage string) *UInt8Flag {
	return DefaultBucket.UInt8(longName, usage)
}

// UInt8P adds a new uint8 flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func UInt8P(longName, usage, shortName string) *UInt8Flag {
	return DefaultBucket.UInt8P(longName, usage, shortName)
}

// Byte adds a new Byte flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. byte).
func Byte(longName, usage string) *ByteFlag {
	return DefaultBucket.Byte(longName, usage)
}

// Bool adds a new boolean flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library.
//
// The value of a boolean flag can be explicitly set using true, false, 1 and 0 (i.e. --enabled true OR --enabled=1).
// The presence of the flag as a CLI argument will also set the flag to true (i.e. --enabled).
func Bool(longName, usage string) *BoolFlag {
	return DefaultBucket.Bool(longName, usage)
}

// BoolSlice adds a new boolean slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. bits).
//
// The value of a BoolSlice flag can be set using a comma (or any custom delimiter) separated string of booleans.
// For example --bits "0, 1, true, false"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func BoolSlice(longName, usage string) *BoolSliceFlag {
	return DefaultBucket.BoolSlice(longName, usage)
}

// Float64 adds a new float64 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. conversion-rate).
func Float64(longName, usage string) *Float64Flag {
	return DefaultBucket.Float64(longName, usage)
}

// Float32 adds a new float32 flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. conversion-rate).
func Float32(longName, usage string) *Float32Flag {
	return DefaultBucket.Float32(longName, usage)
}

// Counter adds a new counter flag to the default bucket.
//
// The value of a counter flag can be increased by repeating the short or the long form of the flag.
// For example, if the short name is 'c', the presence of -cc command line argument will set the value of the counter to 2.
//
// The long name will be automatically converted to lowercase by the library (i.e. count).
func Counter(longName, usage string) *CounterFlag {
	return DefaultBucket.Counter(longName, usage)
}

// Verbosity is an alias for CounterP("verbose", usage, "v").
//
// The value of the verbosity flag can be increased by repeating the short or the long form.
// For example the presence of -vv command line argument will set the verbosity level to 2.
// Having '--verbose -v', '--verbose --verbose' or '-v -v' would have the same effect.
func Verbosity(usage string) *CounterFlag {
	return DefaultBucket.Verbosity(usage)
}

// Duration adds a new Duration flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. ttl).
//
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func Duration(longName, usage string) *DurationFlag {
	return DefaultBucket.Duration(longName, usage)
}

// DurationSlice adds a new duration slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. durations).
//
// The value of a Duration slice flag can be set using a comma (or any custom delimiter) separated string of durations.
//
// Each duration string is a possibly signed sequence of decimal numbers, with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// For example --durations "2s, 2.5s, 5s".
//
// A custom delimiter string can be defined using WithDelimiter() method.
func DurationSlice(longName, usage string) *DurationSliceFlag {
	return DefaultBucket.DurationSlice(longName, usage)
}

// Time adds a new Time flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. birthday).
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
// The long name will be automatically converted to lowercase by the library (i.e. birthday).
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
// The long name will be automatically converted to lowercase by the library (i.e. week-days).
//
// The value of a string slice flag can be set using comma (or any custom delimiter) separated strings.
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
// The long name will be automatically converted to lowercase by the library (i.e. week-days).
// A valid short name is a case sensitive single character string (i.e. w or W).
//
// The value of a string slice flag can be set using comma (or any custom delimiter) separated strings.
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

// IntSlice adds a new int slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of an int slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IntSlice(longName, usage string) *IntSliceFlag {
	return DefaultBucket.IntSlice(longName, usage)
}

// IntSliceP adds a new int slice flag with a short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of an int slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IntSliceP(longName, usage, shortName string) *IntSliceFlag {
	return DefaultBucket.IntSliceP(longName, usage, shortName)
}

// UIntSlice adds a new uint slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a uint slice flag can be set using a comma (or any custom delimiter) separated string of unsigned integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func UIntSlice(longName, usage string) *UIntSliceFlag {
	return DefaultBucket.UIntSlice(longName, usage)
}

// UIntSliceP adds a new uint slice flag with a short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a uint slice flag can be set using a comma (or any custom delimiter) separated string of unsigned integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func UIntSliceP(longName, usage, shortName string) *UIntSliceFlag {
	return DefaultBucket.UIntSliceP(longName, usage, shortName)
}

// Float64Slice adds a new float64 slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a float64 slice flag can be set using a comma (or any custom delimiter) separated string of floating point numbers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func Float64Slice(longName, usage string) *Float64SliceFlag {
	return DefaultBucket.Float64Slice(longName, usage)
}

// IPAddress adds a new IP address flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. endpoint).
//
// The value of an IP address flag can be specified using an IPv4 dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func IPAddress(longName, usage string) *IPAddressFlag {
	return DefaultBucket.IPAddress(longName, usage)
}

// IPAddressP adds a new IP address flag with short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. endpoint).
// A valid short name is a case sensitive single character string (i.e. e or E).
//
// The value of an IP address flag can be specified using an IPv4 dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func IPAddressP(longName, usage, shortName string) *IPAddressFlag {
	return DefaultBucket.IPAddressP(longName, usage, shortName)
}

// IPAddressSlice adds a new IP address slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. endpoints)
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IPAddressSlice(longName, usage string) *IPAddressSliceFlag {
	return DefaultBucket.IPAddressSlice(longName, usage)
}

// IPAddressSliceP adds a new IP address slice flag with a short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. endpoints)
// A valid short name is a case sensitive single character string (i.e. e or E).
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
//
// A custom delimiter string can be defined using WithDelimiter() method.
func IPAddressSliceP(longName, usage, shortName string) *IPAddressSliceFlag {
	return DefaultBucket.IPAddressSliceP(longName, usage, shortName)
}

// CIDR adds a new CIDR flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. network).
//
// The value of a CIDR flag can be defined using a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291. The input will be
// parsed to the IP address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
func CIDR(longName, usage string) *CIDRFlag {
	return DefaultBucket.CIDR(longName, usage)
}

// CIDRSlice adds a new boolean slice flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. networks).
//
// The value of a CIDR slice flag can be defined using a list of CIDR notation IP addresses and prefix length,
// like "192.0.2.0/24, 2001:db8::/32", as defined in RFC 4632 and RFC 4291. Each item will be parsed to the
// address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
//
// A custom delimiter string can be defined using WithDelimiter() method.
func CIDRSlice(longName, usage string) *CIDRSliceFlag {
	return DefaultBucket.CIDRSlice(longName, usage)
}

// StringMap adds a new string map flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. mappings)
//
// The value of a string map flag can be set using standard map initialisation strings.
// For example --mappings '{"key1":"value1", "key2":"value2"}'
func StringMap(longName, usage string) *StringMapFlag {
	return DefaultBucket.StringMapP(longName, usage, "")
}

// StringMapP adds a new string map flag with a short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. mappings).
// A valid short name is a case sensitive single character string (i.e. m or M).
//
// The value of a string map flag can be set using standard map initialisation strings.
// For example --mappings '{"key1":"value1", "key2":"value2"}'
func StringMapP(longName, usage, shortName string) *StringMapFlag {
	return DefaultBucket.StringMapP(longName, usage, shortName)
}

// StringSliceMap adds a new string slice map flag to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. days)
//
// The value of a string slice map flag can be set using standard map initialisation strings.
// Keys are strings and each value is a set of comma (or any custom delimiter) separated strings.
// For example --days '{"Week Days":"Mon,Tue,Wed,Thu,Fri", "Weekend":"Sat,Sun"}'
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
func StringSliceMap(longName, usage string) *StringSliceMapFlag {
	return DefaultBucket.StringSliceMapP(longName, usage, "")
}

// StringSliceMapP adds a new string slice map flag with a short name to the default bucket.
//
// The long name will be automatically converted to lowercase by the library (i.e. days).
// A valid short name is a case sensitive single character string (i.e. d or D).
//
// The value of a string slice map flag can be set using standard map initialisation strings.
// Keys are strings and each value is a set of comma (or any custom delimiter) separated strings.
// For example --days '{"Week Days":"Mon,Tue,Wed,Thu,Fri", "Weekend":"Sat,Sun"}'
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
func StringSliceMapP(longName, usage, shortName string) *StringSliceMapFlag {
	return DefaultBucket.StringSliceMapP(longName, usage, shortName)
}

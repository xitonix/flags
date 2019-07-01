package flags

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

// Bucket represents a container that holds a group of unique flags.
//
// The value of the registered flags will be provided by one of the Sources in the bucket. Each bucket comes with two
// preconfigured sources by default. A command line argument source which is responsible to parse the provided command
// line arguments and an Environment Variable source which queries the system's environment variable registry to extract
// the flag value. By default, the command line argument source has a higher priority over the environment variable source.
// That means the values provided with command line will override their environment variable counterpart.
//
// Apart from the predefined sources, any custom implementation of the `core.Source` interface can be added to the bucket's
// chain of sources (See `flags.MemorySource` for an example). Custom sources can be added using AddSource(), AppendSource()
// and PrependSource() methods.
//
// The Parse method will query all the available sources for a specified key in order.
// The querying process will be stopped as soon as a source has provided a value. If none of the sources has a value to offer,
// the flag will be set to the Default value. In cases where the flag does not have a default value, it will be set to
// the flag type's zero value (for example 0, for an Int flag).
type Bucket struct {
	opts          *config.Options
	reg           *registry
	flags         []core.Flag
	sources       []core.Source
	argSource     *argSource
	helpRequested bool
}

// NewBucket creates a new bucket.
func NewBucket(opts ...config.Option) *Bucket {
	return newBucket(os.Args[1:], internal.OSEnvReader{}, opts...)
}

func newBucket(args []string, envReader internal.EnvironmentVariableReader, opts ...config.Option) *Bucket {
	ops := config.NewOptions()
	for _, option := range opts {
		option(ops)
	}

	argSource, helpRequested := newArgSource(args)
	return &Bucket{
		reg:   newRegistry(),
		flags: make([]core.Flag, 0),
		sources: []core.Source{
			argSource,
			newEnvironmentVarSource(envReader),
		},
		argSource:     argSource,
		helpRequested: helpRequested,
		opts:          ops,
	}
}

// Options returns the current configuration values of the bucket.
func (b *Bucket) Options() *config.Options {
	return b.opts
}

// Flags returns a list of all the registered flags within the bucket.
func (b *Bucket) Flags() []core.Flag {
	return b.flags
}

// Help prints the documentation of the currently registered flag.
//
// You can change the default format by overriding the default HelpFormatter and HelpWriter.
func (b *Bucket) Help() {
	err := b.help()
	if err != nil {
		b.opts.Logger.Print(err)
		b.opts.Terminator.Terminate(core.FailureExitCode)
	}
}

// Parse parses the flags and queries all the available sources in order, to fill the value of each flag.
//
// If none of the sources offers any value, the flag will be set to the specified Default value (if any).
// In case no Default value is defined, the flag will be set to the zero value of its type. For example an
// Int flag will be set to zero.
//
// The order of the default sources is Command Line Arguments > Environment Variables > [Default Value]
//
// Remember that in order for the flag values to be extractable from the environment variables
// (or all the other custom sources) it MUST have a key associated with it.
//
// See flags.EnableAutoKeyGeneration(), flags.SetKeyPrefix() and different flags' WithKey() method for more details.
func (b *Bucket) Parse() {
	b.init()

	if b.helpRequested {
		b.Help()
		b.opts.Terminator.Terminate(core.SuccessExitCode)
		return
	}
	if err := b.checkForUnknownFlags(); err != nil {
		b.Help()
		b.terminateWithError(err)
		return
	}

	for _, f := range b.flags {
		if f.IsRequired() && f.IsDeprecated() {
			pn := internal.GetPrintName(f.LongName(), f.ShortName())
			b.terminateWithError(fmt.Errorf("%s is marked as deprecated. An obsolete flag cannot be mandatory", pn))
			return
		}
		for _, src := range b.sources {
			var (
				found bool
				value string
			)

			argSrc, isArgs := src.(*argSource)

			if isArgs {
				value, found = src.Read("--" + f.LongName())
				if !found {
					value, found = src.Read("-" + f.ShortName())
				}
				if !found || internal.IsEmpty(value) {
					if repeatable, isRepeatable := f.(core.Repeatable); isRepeatable {
						count := argSrc.getNumberOfRepeats(f)
						if count > 0 {
							// Either the short form or the long form has been
							// provided at least once
							value = strconv.Itoa(count * repeatable.Once())
							found = true
						}
					}
				}
			}

			if !found && !isArgs && f.Key().IsSet() {
				value, found = src.Read(f.Key().String())
			}

			if !found {
				f.ResetToDefault()
				continue
			}

			if p, ok := f.(core.EmptyValueProvider); ok && found && internal.IsEmpty(value) {
				value = p.EmptyValue()
			}
			if !b.executeCallback(f, value, false) {
				return
			}
			err := f.Set(value)
			if err != nil {
				b.terminateWithError(err)
				return
			}

			if !b.executeCallback(f, value, true) {
				return
			}
			break
		}
		if f.IsRequired() && !f.IsSet() {
			pn := internal.GetPrintName(f.LongName(), f.ShortName())
			b.terminateWithError(fmt.Errorf("%s flag is required.", pn))
			return
		}
	}
}

// AppendSource appends a new source to the end of the chain.
//
// With the default configuration, the order will be:
// Command Line Arguments > Environment Variables > src > [Default Value]
//
// Note that the Parse method will query the sources in order.
func (b *Bucket) AppendSource(src core.Source) {
	if src == nil {
		return
	}
	b.sources = append(b.sources, src)
}

// PrependSource prepends a new source to the beginning of the chain.
// This is an alias for AddSource(src, 0)
//
// With the default configuration, the order will be:
// src > Command Line Arguments > Environment Variables > [Default Value]
//
// Note that the Parse method will query the sources in order.
func (b *Bucket) PrependSource(src core.Source) {
	if src == nil {
		return
	}
	b.sources = append([]core.Source{src}, b.sources...)
}

// AddSource inserts the new source at the specified index.
//
// If the index is <= 0 the new source will get added to the beginning of the chain. If the index is greater than the
// current number of sources, it will get be appended the end.
//
// Note that the Parse method will query the sources in order.
func (b *Bucket) AddSource(src core.Source, index int) {
	if src == nil {
		return
	}

	if index < 0 {
		index = 0
	}
	if index > len(b.sources) {
		index = len(b.sources)
	}
	b.sources = append(b.sources[:index], append([]core.Source{src}, b.sources[index:]...)...)
}

// FullString adds a new string flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library.
func (b *Bucket) String(longName, usage string) *StringFlag {
	return b.StringP(longName, usage, "")
}

// StringP adds a new string flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. file-path).
// A valid short name is a case sensitive single character string (i.e. f or F).
func (b *Bucket) StringP(longName, usage, shortName string) *StringFlag {
	f := newString(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int adds a new Int flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Int(longName, usage string) *IntFlag {
	return b.IntP(longName, usage, "")
}

// IntP adds a new Int flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) IntP(longName, usage, shortName string) *IntFlag {
	f := newInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int8 adds a new Int8 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Int8(longName, usage string) *Int8Flag {
	return b.Int8P(longName, usage, "")
}

// Int8P adds a new Int8 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) Int8P(longName, usage, shortName string) *Int8Flag {
	f := newInt8(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int16 adds a new Int16 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Int16(longName, usage string) *Int16Flag {
	return b.Int16P(longName, usage, "")
}

// Int16P adds a new Int16 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) Int16P(longName, usage, shortName string) *Int16Flag {
	f := newInt16(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int32 adds a new Int32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Int32(longName, usage string) *Int32Flag {
	return b.Int32P(longName, usage, "")
}

// Int32P adds a new Int32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) Int32P(longName, usage, shortName string) *Int32Flag {
	f := newInt32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int64 adds a new Int64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Int64(longName, usage string) *Int64Flag {
	return b.Int64P(longName, usage, "")
}

// Int64P adds a new Int64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) Int64P(longName, usage, shortName string) *Int64Flag {
	f := newInt64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt adds a new UInt flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) UInt(longName, usage string) *UIntFlag {
	return b.UIntP(longName, usage, "")
}

// UIntP adds a new UInt flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) UIntP(longName, usage, shortName string) *UIntFlag {
	f := newUInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt64 adds a new UInt64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) UInt64(longName, usage string) *UInt64Flag {
	return b.UInt64P(longName, usage, "")
}

// UInt64P adds a new UInt64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) UInt64P(longName, usage, shortName string) *UInt64Flag {
	f := newUInt64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt32 adds a new UInt32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) UInt32(longName, usage string) *UInt32Flag {
	return b.UInt32P(longName, usage, "")
}

// UInt32P adds a new UInt32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) UInt32P(longName, usage, shortName string) *UInt32Flag {
	f := newUInt32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt16 adds a new UInt16 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) UInt16(longName, usage string) *UInt16Flag {
	return b.UInt16P(longName, usage, "")
}

// UInt16P adds a new UInt16 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) UInt16P(longName, usage, shortName string) *UInt16Flag {
	f := newUInt16(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt8 adds a new uint8 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) UInt8(longName, usage string) *UInt8Flag {
	return b.UInt8P(longName, usage, "")
}

// UInt8P adds a new uint8 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) UInt8P(longName, usage, shortName string) *UInt8Flag {
	f := newUInt8(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Byte adds a new byte flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
func (b *Bucket) Byte(longName, usage string) *ByteFlag {
	return b.ByteP(longName, usage, "")
}

// ByteP adds a new byte flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. port-number).
// A valid short name is a case sensitive single character string (i.e. p or P).
func (b *Bucket) ByteP(longName, usage, shortName string) *ByteFlag {
	f := newByte(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Bool adds a new bool flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. enable-write-access).
func (b *Bucket) Bool(longName, usage string) *BoolFlag {
	return b.BoolP(longName, usage, "")
}

// BoolP adds a new bool flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. enable-write-access).
// A valid short name is a case sensitive single character string (i.e. e or E).
func (b *Bucket) BoolP(longName, usage, shortName string) *BoolFlag {
	f := newBool(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// BoolSlice adds a new int slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. bits)
//
// The value of a BoolSlice flag can be set using a comma (or any custom delimiter) separated string of booleans.
// For example --bits "0, 1, true, false"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) BoolSlice(longName, usage string) *BoolSliceFlag {
	return b.BoolSliceP(longName, usage, "")
}

// BoolSliceP adds a new int slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. bits)
// A valid short name is a case sensitive single character string (i.e. b or B).
//
// The value of a BoolSlice flag can be set using a comma (or any custom delimiter) separated string of booleans.
// For example --bits "0, 1, true, false"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) BoolSliceP(longName, usage, shortName string) *BoolSliceFlag {
	f := newBoolSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Float64 adds a new float64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
func (b *Bucket) Float64(longName, usage string) *Float64Flag {
	return b.Float64P(longName, usage, "")
}

// Float64P adds a new float64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
// A valid short name is a case sensitive single character string (i.e. c or C).
func (b *Bucket) Float64P(longName, usage, shortName string) *Float64Flag {
	f := newFloat64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Float32 adds a new float32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
func (b *Bucket) Float32(longName, usage string) *Float32Flag {
	return b.Float32P(longName, usage, "")
}

// Float32P adds a new float32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. conversion-rate).
// A valid short name is a case sensitive single character string (i.e. c or C).
func (b *Bucket) Float32P(longName, usage, shortName string) *Float32Flag {
	f := newFloat32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// CounterP adds a new counter flag with a short name to the bucket.
//
// The value of a counter flag can be increased by repeating the short or the long form.
// For example the presence of -vv command line argument will set the value of the counter to 2.
//
// Long names will be automatically converted to lowercase by the library (i.e. verbosity).
// A valid short name is a case sensitive single character string (i.e. v or V).
func (b *Bucket) CounterP(longName, usage, shortName string) *CounterFlag {
	f := newCounter(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// VerbosityP is an alias for CounterP("verbose", usage, "v").
//
// The value of the verbosity flag can be increased by repeating the short or the long form.
// For example the presence of -vv command line argument will set the verbosity level to 2.
// Having '--verbose -v', '--verbose --verbose' or '-v -v' would have the same effect.
func (b *Bucket) VerbosityP(usage string) *CounterFlag {
	return b.CounterP("verbose", usage, "v")
}

// Duration adds a new Duration flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ttl).
//
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (b *Bucket) Duration(longName, usage string) *DurationFlag {
	return b.DurationP(longName, usage, "")
}

// DurationP adds a new Duration flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ttl).
// A valid short name is a case sensitive single character string (i.e. t or T).
//
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (b *Bucket) DurationP(longName, usage, shortName string) *DurationFlag {
	f := newDuration(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// DurationSlice adds a new duration slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. durations)
//
// The value of a Duration slice flag can be set using a comma (or any custom delimiter) separated string of durations.
//
// Each duration string is a possibly signed sequence of decimal numbers, with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// For example --durations "2s, 2.5s, 5s".
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) DurationSlice(longName, usage string) *DurationSliceFlag {
	return b.DurationSliceP(longName, usage, "")
}

// DurationSliceP adds a new duration slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. durations)
// A valid short name is a case sensitive single character string (i.e. d or D).
//
// The value of a Duration slice flag can be set using a comma (or any custom delimiter) separated string of durations.
//
// Each duration string is a possibly signed sequence of decimal numbers, with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// For example --durations "2s, 2.5s, 5s".
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) DurationSliceP(longName, usage, shortName string) *DurationSliceFlag {
	f := newDurationSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Time adds a new Time flag to the bucket.
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
func (b *Bucket) Time(longName, usage string) *TimeFlag {
	return b.TimeP(longName, usage, "")
}

// TimeP adds a new Time flag with a short name to the bucket.
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
func (b *Bucket) TimeP(longName, usage, shortName string) *TimeFlag {
	f := newTime(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// StringSlice adds a new string slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. week-days)
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
func (b *Bucket) StringSlice(longName, usage string) *StringSliceFlag {
	return b.StringSliceP(longName, usage, "")
}

// StringSliceP adds a new string slice flag with a short name to the bucket.
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
func (b *Bucket) StringSliceP(longName, usage, shortName string) *StringSliceFlag {
	f := newStringSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// IntSlice adds a new int slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a IntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) IntSlice(longName, usage string) *IntSliceFlag {
	return b.IntSliceP(longName, usage, "")
}

// IntSliceP adds a new int slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a IntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) IntSliceP(longName, usage, shortName string) *IntSliceFlag {
	f := newIntSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UIntSlice adds a new uint slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a UIntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) UIntSlice(longName, usage string) *UIntSliceFlag {
	return b.UIntSliceP(longName, usage, "")
}

// UIntSliceP adds a new uint slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a UIntSlice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) UIntSliceP(longName, usage, shortName string) *UIntSliceFlag {
	f := newUIntSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Float64Slice adds a new float64 slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) Float64Slice(longName, usage string) *Float64SliceFlag {
	return b.Float64SliceP(longName, usage, "")
}

// Float64SliceP adds a new float64 slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. numbers)
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) Float64SliceP(longName, usage, shortName string) *Float64SliceFlag {
	f := newFloat64Slice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// IPAddress adds a new IPAddress flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ip-address).
//
// The value of an IP address flag can be specified using a dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func (b *Bucket) IPAddress(longName, usage string) *IPAddressFlag {
	return b.IPAddressP(longName, usage, "")
}

// IPAddressP adds a new IPAddress flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. ip-address).
// A valid short name is a case sensitive single character string (i.e. i or I).
//
// The value of an IP address flag can be specified using a dotted decimal (i.e. "192.0.2.1")
// or an IPv6 ("2001:db8::68") formatted string.
func (b *Bucket) IPAddressP(longName, usage, shortName string) *IPAddressFlag {
	f := newIPAddress(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// IPAddressSlice adds a new IP Address slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. ip-addresses)
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) IPAddressSlice(longName, usage string) *IPAddressSliceFlag {
	return b.IPAddressSliceP(longName, usage, "")
}

// IPAddressSliceP adds a new IP Address slice flag with a short name to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. ip-addresses)
// A valid short name is a case sensitive single character string (i.e. i or I).
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) IPAddressSliceP(longName, usage, shortName string) *IPAddressSliceFlag {
	f := newIPAddressSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// CIDR adds a new CIDR flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. network).
//
// The value of a CIDR flag can be defined using a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291. The input will be
// parsed to the IP address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
func (b *Bucket) CIDR(longName, usage string) *CIDRFlag {
	return b.CIDRP(longName, usage, "")
}

// CIDRP adds a new CIDR flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (i.e. network).
// A valid short name is a case sensitive single character string (i.e. n or N).
//
// The value of a CIDR flag can be defined using a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291. The input will be
// parsed to the IP address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
func (b *Bucket) CIDRP(longName, usage, shortName string) *CIDRFlag {
	f := newCIDR(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// CIDRSlice adds a new CIDR slice flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library (i.e. durations)
//
// The value of a CIDR slice flag can be defined using a list of CIDR notation IP addresses and prefix length,
// like "192.0.2.0/24, 2001:db8::/32", as defined in RFC 4632 and RFC 4291. Each item will be parsed to the
// address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) CIDRSlice(longName, usage string) *CIDRSliceFlag {
	return b.CIDRSliceP(longName, usage, "")
}

// CIDRSliceP adds a new CIDR slice flag with a short name to the bucket.
//
// The value of a CIDR slice flag can be defined using a list of CIDR notation IP addresses and prefix length,
// like "192.0.2.0/24, 2001:db8::/32", as defined in RFC 4632 and RFC 4291. Each item will be parsed to the
// address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (b *Bucket) CIDRSliceP(longName, usage, shortName string) *CIDRSliceFlag {
	f := newCIDRSlice(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

func (b *Bucket) help() error {
	flags := b.sortFlags()
	for _, flag := range flags {
		_, err := b.opts.HelpWriter.Write([]byte(b.opts.HelpFormatter.Format(flag, b.opts.DeprecationMark, b.opts.DefaultValueFormatString)))
		if err != nil {
			return err
		}
	}
	return b.opts.HelpWriter.Close()
}

func (b *Bucket) checkForUnknownFlags() error {
	for arg := range b.argSource.arguments {
		if b.reg.isRegistered(arg) || b.reg.isReserved(arg) {
			continue
		}
		return core.NewUnknownFlagErr(arg)
	}
	return nil
}

func (b *Bucket) init() {
	for _, f := range b.flags {
		if !internal.IsEmpty(b.opts.KeyPrefix) {
			f.Key().SetPrefix(b.opts.KeyPrefix)
		}

		if b.opts.AutoKeys && !f.Key().IsSet() {
			f.Key().SetID(f.LongName())
		}
		err := b.reg.add(f)
		if err != nil {
			b.opts.Logger.Print(err)
			b.opts.Terminator.Terminate(core.FailureExitCode)
		}
	}
}

func (b *Bucket) sortFlags() []core.Flag {
	if b.opts.Comparer == nil {
		return b.flags
	}

	clone := make([]core.Flag, len(b.flags))
	copy(clone, b.flags)
	sort.Slice(clone, func(i, j int) bool {
		return b.opts.Comparer.LessThan(clone[i], clone[j])
	})
	return clone
}

func (b *Bucket) executeCallback(f core.Flag, value string, post bool) bool {
	cb := b.opts.PreSetCallback
	if post {
		cb = b.opts.PostSetCallback
	}
	if cb == nil {
		return true
	}
	if err := cb(f, value); err != nil {
		b.terminateWithError(err)
		return false
	}

	return true
}

func (b *Bucket) terminateWithError(err error) {
	b.opts.Logger.Print(err)
	b.opts.Terminator.Terminate(core.FailureExitCode)
}

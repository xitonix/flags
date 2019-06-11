package flags

import (
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
		b.opts.Logger.Print(err)
		b.opts.Terminator.Terminate(core.FailureExitCode)
		return
	}

	for _, f := range b.flags {
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

			err := f.Set(value)
			if err != nil {
				b.opts.Logger.Print(err)
				b.opts.Terminator.Terminate(core.FailureExitCode)
			}
			break
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

// String adds a new string flag to the bucket.
//
// The long names will be automatically converted to lowercase by the library.
func (b *Bucket) String(longName, usage string) *StringFlag {
	return b.StringP(longName, usage, "")
}

// StringP adds a new string flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func (b *Bucket) StringP(longName, usage, shortName string) *StringFlag {
	f := newString(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int adds a new Int flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) Int(longName, usage string) *IntFlag {
	return b.IntP(longName, usage, "")
}

// IntP adds a new Int flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) IntP(longName, usage, shortName string) *IntFlag {
	f := newInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int8 adds a new Int8 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) Int8(longName, usage string) *Int8Flag {
	return b.Int8P(longName, usage, "")
}

// Int8P adds a new Int8 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) Int8P(longName, usage, shortName string) *Int8Flag {
	f := newInt8(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int16 adds a new Int16 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) Int16(longName, usage string) *Int16Flag {
	return b.Int16P(longName, usage, "")
}

// Int16P adds a new Int16 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) Int16P(longName, usage, shortName string) *Int16Flag {
	f := newInt16(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int32 adds a new Int32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) Int32(longName, usage string) *Int32Flag {
	return b.Int32P(longName, usage, "")
}

// Int32P adds a new Int32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) Int32P(longName, usage, shortName string) *Int32Flag {
	f := newInt32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int64 adds a new Int64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) Int64(longName, usage string) *Int64Flag {
	return b.Int64P(longName, usage, "")
}

// Int64P adds a new Int64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) Int64P(longName, usage, shortName string) *Int64Flag {
	f := newInt64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt adds a new UInt flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) UInt(longName, usage string) *UIntFlag {
	return b.UIntP(longName, usage, "")
}

// UIntP adds a new UInt flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) UIntP(longName, usage, shortName string) *UIntFlag {
	f := newUInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt64 adds a new UInt64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) UInt64(longName, usage string) *UInt64Flag {
	return b.UInt64P(longName, usage, "")
}

// UInt64P adds a new UInt64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) UInt64P(longName, usage, shortName string) *UInt64Flag {
	f := newUInt64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt32 adds a new UInt32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) UInt32(longName, usage string) *UInt32Flag {
	return b.UInt32P(longName, usage, "")
}

// UInt32P adds a new UInt32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) UInt32P(longName, usage, shortName string) *UInt32Flag {
	f := newUInt32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt16 adds a new UInt16 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) UInt16(longName, usage string) *UInt16Flag {
	return b.UInt16P(longName, usage, "")
}

// UInt16P adds a new UInt16 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) UInt16P(longName, usage, shortName string) *UInt16Flag {
	f := newUInt16(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// UInt8 adds a new UInt8 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
func (b *Bucket) UInt8(longName, usage string) *UInt8Flag {
	return b.UInt8P(longName, usage, "")
}

// UInt8P adds a new UInt8 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. port-number).
// A valid short name is a case sensitive single character string (ie. p or P).
func (b *Bucket) UInt8P(longName, usage, shortName string) *UInt8Flag {
	f := newUInt8(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Bool adds a new Bool flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. enable-write-access).
func (b *Bucket) Bool(longName, usage string) *BoolFlag {
	return b.BoolP(longName, usage, "")
}

// BoolP adds a new Bool flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. enable-write-access).
// A valid short name is a case sensitive single character string (ie. e or E).
func (b *Bucket) BoolP(longName, usage, shortName string) *BoolFlag {
	f := newBool(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Float64 adds a new Float64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. conversion-rate).
func (b *Bucket) Float64(longName, usage string) *Float64Flag {
	return b.Float64P(longName, usage, "")
}

// Float64P adds a new Float64 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. conversion-rate).
// A valid short name is a case sensitive single character string (ie. c or C).
func (b *Bucket) Float64P(longName, usage, shortName string) *Float64Flag {
	f := newFloat64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Float32 adds a new Float32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. conversion-rate).
func (b *Bucket) Float32(longName, usage string) *Float32Flag {
	return b.Float32P(longName, usage, "")
}

// Float32P adds a new Float32 flag with a short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. conversion-rate).
// A valid short name is a case sensitive single character string (ie. c or C).
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
// Long names will be automatically converted to lowercase by the library (ie. verbosity).
// A valid short name is a case sensitive single character string (ie. v or V).
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
// Long names will be automatically converted to lowercase by the library (ie. ttl).
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
// Long names will be automatically converted to lowercase by the library (ie. ttl).
// A valid short name is a case sensitive single character string (ie. t or T).
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

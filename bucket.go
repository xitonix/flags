package flags

import (
	"os"
	"sort"

	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

// Bucket represents a container that holds a group of flags.
//
// Each bucket may contain a set of unique flags.
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
// You can customise the default format by overriding the help provider.
// Check config.WithHelpProvider(...) method for more details.
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

			_, isArgs := src.(*argSource)

			if isArgs {
				value, found = src.Read("--" + f.LongName())
				if !found {
					value, found = src.Read("-" + f.ShortName())
				}
			}

			if !found && !isArgs && f.Key().IsSet() {
				value, found = src.Read(f.Key().Get())
			}

			if !found {
				f.ResetToDefault()
				continue
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

// StringP adds a new string flag with short name to the bucket.
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
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func (b *Bucket) Int(longName, usage string) *IntFlag {
	return b.IntP(longName, usage, "")
}

// IntP adds a new Int flag with short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func (b *Bucket) IntP(longName, usage, shortName string) *IntFlag {
	f := newInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int64 adds a new Int64 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func (b *Bucket) Int64(longName, usage string) *Int64Flag {
	return b.Int64P(longName, usage, "")
}

// Int64P adds a new Int64 flag with short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func (b *Bucket) Int64P(longName, usage, shortName string) *Int64Flag {
	f := newInt64(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

// Int32 adds a new Int32 flag to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
func (b *Bucket) Int32(longName, usage string) *Int32Flag {
	return b.Int32P(longName, usage, "")
}

// Int32P adds a new Int32 flag with short name to the bucket.
//
// Long names will be automatically converted to lowercase by the library (ie. file-path).
// A valid short name is a case sensitive single character string (ie. f or F).
func (b *Bucket) Int32P(longName, usage, shortName string) *Int32Flag {
	f := newInt32(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

func (b *Bucket) help() error {
	flags := b.sortFlags()
	for _, flag := range flags {
		_, err := b.opts.HelpProvider.Writer.Write([]byte(b.opts.HelpProvider.Formatter.Format(flag, b.opts.DeprecationMark, b.opts.DefaultValueFormatString)))
		if err != nil {
			return err
		}
	}
	return b.opts.HelpProvider.Writer.Close()
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
			f.Key().Set(f.LongName())
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

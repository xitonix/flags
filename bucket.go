package flags

import (
	"os"
	"sort"

	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type Bucket struct {
	opts          *config.Options
	reg           *registry
	flags         []core.Flag
	sources       []core.Source
	argSource     *argSource
	helpRequested bool
}

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

func (b *Bucket) Options() *config.Options {
	return b.opts
}

func (b *Bucket) Flags() []core.Flag {
	return b.flags
}

func (b *Bucket) Help() {
	err := b.help()
	if err != nil {
		b.opts.Logger.Print(err)
		b.opts.Terminator.Terminate(core.FailureExitCode)
	}
}

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

func (b *Bucket) AppendSource(src core.Source) {
	if src == nil {
		return
	}
	b.sources = append(b.sources, src)
}

func (b *Bucket) PrependSource(src core.Source) {
	if src == nil {
		return
	}
	b.sources = append([]core.Source{src}, b.sources...)
}

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

func (b *Bucket) String(longName, usage string) *StringFlag {
	return b.StringP(longName, usage, "")
}

func (b *Bucket) StringP(longName, usage, shortName string) *StringFlag {
	f := newString(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

func (b *Bucket) Int(longName, usage string) *IntFlag {
	return b.IntP(longName, usage, "")
}

func (b *Bucket) IntP(longName, usage, shortName string) *IntFlag {
	f := newInt(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

func (b *Bucket) Int64(longName, usage string) *Int64Flag {
	return b.Int64P(longName, usage, "")
}

func (b *Bucket) Int64P(longName, usage, shortName string) *Int64Flag {
	f := newInt64(longName, usage, shortName)
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

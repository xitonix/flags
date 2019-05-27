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
	sources       []Source
	argSource     *argSource
	helpRequested bool
}

func NewBucket(opts ...config.Option) *Bucket {
	return newBucket(os.Args[1:], opts...)
}

func newBucket(args []string, opts ...config.Option) *Bucket {
	ops := config.NewOptions()
	for _, option := range opts {
		option(ops)
	}

	argSource, helpRequested := newArgSource(args)
	return &Bucket{
		reg:   newRegistry(),
		flags: make([]core.Flag, 0),
		sources: []Source{
			argSource,
		},
		argSource:     argSource,
		helpRequested: helpRequested,
		opts:          ops,
	}
}

func (b *Bucket) Options() *config.Options {
	return b.opts
}

func (b *Bucket) String(longName, usage string) *StringFlag {
	return StringP(longName, usage, "")
}

func (b *Bucket) StringP(longName, usage, shortName string) *StringFlag {
	f := newString(longName, usage, shortName)
	b.flags = append(b.flags, f)
	return f
}

func (b *Bucket) Flags() []core.Flag {
	return b.flags
}

func (b *Bucket) Help() {
	flags := b.sortFlags()
	for _, flag := range flags {
		_, err := b.opts.HelpProvider.Writer.Write([]byte(b.opts.HelpProvider.Formatter.Format(flag)))
		if err != nil {
			b.opts.Logger.Print(err)
			b.opts.Terminator.Terminate(core.FailureExitCode)
		}
	}
	err := b.opts.HelpProvider.Writer.Close()
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
	}

	if err := b.checkForUnknownFlags(); err != nil {
		b.Help()
		b.opts.Logger.Print(err)
		b.opts.Terminator.Terminate(core.FailureExitCode)
	}
	for _, f := range b.flags {
		for _, src := range b.sources {
			value, found := src.Read(f.LongName())
			if !found {
				value, found = src.Read(f.ShortName())
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

func (b *Bucket) enableAutoKeyGen() {
	b.opts.AutoKeys = true
}

func (b *Bucket) setKeyPrefix(prefix string) {
	b.opts.KeyPrefix = internal.SanitiseFlagID(prefix)
}

func (b *Bucket) setLogger(logger core.Logger) {
	b.opts.Logger = logger
}

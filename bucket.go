package flags

import (
	"os"

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
	ops := config.NewOptions()
	for _, option := range opts {
		option(ops)
	}

	argSource, helpRequested := newArgSource(os.Args[1:])
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

func (b *Bucket) String(longName, usage string) *StringFlag {
	f := newString(longName, usage, "")
	b.addFlag(f)
	return f
}

func (b *Bucket) StringP(longName, usage, shortName string) *StringFlag {
	f := newString(longName, usage, shortName)
	b.addFlag(f)
	return f
}

func (b *Bucket) Flags() []core.Flag {
	return b.flags
}

func (b *Bucket) Help() {
	for _, flag := range b.flags {
		_, _ = b.opts.HelpProvider.Writer.Write([]byte(b.opts.HelpProvider.Formatter.Format(flag)))
	}
	_ = b.opts.HelpProvider.Writer.Close()
}

func (b *Bucket) Parse() {
	b.init()

	if b.helpRequested {
		b.Help()
		os.Exit(0)
	}

	if err := b.checkForUnknownFlags(); err != nil {
		b.Help()
		b.opts.Log.Fatal(err)
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
				b.opts.Log.Fatal(err)
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
	if !internal.IsEmpty(b.opts.EnvPrefix) {
		for _, f := range b.flags {
			f.Env().SetPrefix(b.opts.EnvPrefix)
		}
	}

	if b.opts.AutoEnv {
		for _, f := range b.flags {
			f.Env().Set(f.LongName(), true)
		}
	}
}

func (b *Bucket) addFlag(flag core.Flag) {
	err := b.reg.add(flag)
	if err != nil {
		b.opts.Log.Fatal(err)
	}
	b.flags = append(b.flags, flag)
}

func (b *Bucket) enableAutoEnv() {
	b.opts.AutoEnv = true
}

func (b *Bucket) setEnvPrefix(prefix string) {
	b.opts.EnvPrefix = internal.SanitiseEnvVarName(prefix)
}

func (b *Bucket) setLogger(logger core.Logger) {
	b.opts.Log = logger
}

package flags

import (
	"fmt"
	"os"
	"strings"

	"go.xitonix.io/flags/config"
	"text/tabwriter"
)

type Bucket struct {
	reg           *registry
	flags         []Flag
	sources       []Source
	ops           *config.Options
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
		flags: make([]Flag, 0),
		sources: []Source{
			argSource,
		},
		argSource:     argSource,
		helpRequested: helpRequested,
		ops:           ops,
	}
}

func (b *Bucket) String(name string, defaultValue string, usage string) *StringFlag {
	return b.StringP(name, "", defaultValue, usage)
}

func (b *Bucket) StringP(name string, short string, defaultValue string, usage string) *StringFlag {
	f := newStringP(strings.TrimSpace(strings.ToLower(name)), strings.TrimSpace(strings.ToLower(short)), defaultValue, usage)
	b.addFlag(f)
	return f
}

func (b *Bucket) Flags() []Flag {
	return b.flags
}

func (b *Bucket) Help() {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	for _, flag := range b.flags {
		_, _ = fmt.Fprintln(tw, flag.FormatHelp())
	}
	_ = tw.Flush()
}

func (b *Bucket) Parse() {
	b.init()

	if b.helpRequested {
		b.Help()
		os.Exit(0)
	}

	if err := b.checkForUnknownFlags(); err != nil {
		b.Help()
		b.ops.Log.Fatal(err)
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

			var err error
			switch flag := f.(type) {
			case *StringFlag:
				err = flag.Set(value)
			}
			if err != nil {
				b.ops.Log.Fatal(err)
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
		return errUnknownFlag(arg)
	}
	return nil
}

func (b *Bucket) init() {
	if !isEmpty(b.ops.EnvPrefix) {
		for _, f := range b.flags {
			f.Env().setPrefix(b.ops.EnvPrefix)
		}
	}

	if b.ops.AutoEnv {
		for _, f := range b.flags {
			f.Env().auto(f.LongName())
		}
	}
}

func (b *Bucket) addFlag(flag Flag) {
	err := b.reg.add(flag)
	if err != nil {
		b.ops.Log.Fatal(err)
	}
	b.flags = append(b.flags, flag)
}

func (b *Bucket) enableAutoEnv() {
	b.ops.AutoEnv = true
}

func (b *Bucket) setEnvPrefix(prefix string) {
	b.ops.EnvPrefix = sanitise(prefix)
}

func (b *Bucket) setLogger(logger config.Logger) {
	b.ops.Log = logger
}

package flags

import "os"

type Bucket struct {
	flags   map[string]Flag
	sources []Source
	ops     *Options
}

var DefaultBucket = NewBucket()

func NewBucket(opts ...Option) *Bucket {
	ops := NewOptions()
	for _, option := range opts {
		option(ops)
	}

	return &Bucket{
		flags: make(map[string]Flag),
		sources: []Source{
			newArgSource(os.Args[1:]),
		},
		ops: ops,
	}
}

func (b *Bucket) String(name string, defaultValue string, usage string) *StringFlag {
	return b.StringP(name, "", defaultValue, usage)
}

func (b *Bucket) StringP(name string, short string, defaultValue string, usage string) *StringFlag {
	id := name + ":" + short
	f := newStringP(name, short, defaultValue, usage)
	//TODO: Fail to duplicate keys
	b.flags[id] = f
	return f
}

func (b *Bucket) Parse() {
	if !isEmpty(b.ops.EnvPrefix) {
		for _, f := range b.flags {
			f.Env().setPrefix(b.ops.EnvPrefix)
		}
	}

	if b.ops.AutoEnv {
		for _, f := range b.flags {
			f.Env().auto()
		}
	}

	for _, f := range b.flags {
		for _, src := range b.sources {
			value, found := src.Read(f.Name())
			if !found {
				value, found = src.Read(f.Short())
			}
			if found {
				switch flag := f.(type) {
				case *StringFlag:
					flag.Set(value)
				}
				continue
			}
		}
	}
}

func Parse() {
	DefaultBucket.Parse()
}

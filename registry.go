package flags

import (
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type registry struct {
	shorts map[string]interface{}
	longs  map[string]interface{}
}

var (
	reserved = map[string]interface{}{
		"h":    nil,
		"help": nil,
	}
)

func newRegistry() *registry {
	return &registry{
		shorts: make(map[string]interface{}),
		longs:  make(map[string]interface{}),
	}
}

func (r *registry) add(flag core.Flag) error {
	long := strings.ToLower(strings.TrimSpace(flag.LongName()))

	if internal.IsEmpty(long) {
		return core.ErrEmptyFlagName
	}

	if r.isReserved(long) {
		return core.NewInvalidFlagErr(long, "", "is a reserved flag")
	}

	short := flag.ShortName()
	if r.isReserved(short) {
		return core.NewInvalidFlagErr("", short, "is a reserved flag")
	}

	if r.isLongNameRegistered(long) {
		return core.NewInvalidFlagErr(long, "", "flag already exists")
	}

	if r.isShortNameRegistered(short) {
		return core.NewInvalidFlagErr("", short, "flag already exists")
	}

	r.longs[long] = nil
	if len(short) > 0 {
		r.shorts[short] = nil
	}
	return nil
}

func (r *registry) isRegistered(name string) bool {
	return r.isLongNameRegistered(name) || r.isShortNameRegistered(name)
}

func (r *registry) isReserved(name string) bool {
	_, ok := reserved[name]
	return ok
}

func (r *registry) isLongNameRegistered(name string) bool {
	_, ok := r.longs[name]
	return ok
}

func (r *registry) isShortNameRegistered(name string) bool {
	_, ok := r.shorts[name]
	return ok
}

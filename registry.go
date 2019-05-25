package flags

import (
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type registry struct {
	book map[string]interface{}
}

var (
	reserved = map[string]interface{}{
		"h":    nil,
		"help": nil,
	}
)

func newRegistry() *registry {
	return &registry{
		book: make(map[string]interface{}),
	}
}

func (r *registry) add(flag core.Flag) error {
	long := strings.ToLower(strings.TrimSpace(flag.LongName()))

	if internal.IsEmpty(long) {
		return core.ErrEmptyFlagName
	}

	if r.isReserved(long) {
		return core.NewInvalidFlagErr(long, "", "", "is a reserved flag")
	}

	short := flag.ShortName()
	if r.isReserved(short) {
		return core.NewInvalidFlagErr("", short, "", "is a reserved flag")
	}

	if r.isRegistered(long) {
		return core.NewInvalidFlagErr(long, "", "", "flag already exists")
	}

	if r.isRegistered(short) {
		return core.NewInvalidFlagErr("", short, "", "flag already exists")
	}

	key := flag.Key().Value()
	if r.isRegistered(key) {
		return core.NewInvalidFlagErr("", "", key, "flag key already exists")
	}

	r.book[long] = nil
	if len(short) > 0 {
		r.book[short] = nil
	}

	if len(key) > 0 {
		r.book[key] = key
	}

	return nil
}

func (r *registry) isRegistered(name string) bool {
	_, ok := r.book[name]
	return ok
}

func (r *registry) isReserved(name string) bool {
	_, ok := reserved[name]
	return ok
}

func (r *registry) isShortNameRegistered(name string) bool {
	_, ok := r.book[name]
	return ok
}

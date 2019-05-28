package flags

import (
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

type registry struct {
	catalogue map[string]interface{}
}

var (
	reserved = map[string]interface{}{
		"-h":     nil,
		"--help": nil,
	}
)

func newRegistry() *registry {
	return &registry{
		catalogue: make(map[string]interface{}),
	}
}

func (r *registry) addLongNameIfValid(longName string) error {
	long := strings.ToLower(strings.TrimSpace(longName))

	if internal.IsEmpty(long) {
		return core.ErrEmptyFlagName
	}

	long = "--" + long

	if r.isReserved(long) {
		return core.NewInvalidFlagErr(long, "", "", "is a reserved flag")
	}

	if _, ok := r.catalogue[long]; ok {
		return core.NewInvalidFlagErr(long, "", "", "flag already exists")
	}

	r.catalogue[long] = nil
	return nil
}

func (r *registry) addShortNameIfValid(shortName string) error {
	if len(shortName) == 0 {
		return nil
	}
	shortName = "-" + shortName
	if r.isReserved(shortName) {
		return core.NewInvalidFlagErr("", shortName, "", "is a reserved flag")
	}

	if _, ok := r.catalogue[shortName]; ok {
		return core.NewInvalidFlagErr("", shortName, "", "flag already exists")
	}
	r.catalogue[shortName] = nil
	return nil
}

func (r *registry) addKeyIfValid(key string) error {
	if len(key) == 0 {
		return nil
	}
	if _, ok := r.catalogue[key]; ok {
		return core.NewInvalidFlagErr("", "", key, "flag key already exists")
	}

	if len(key) > 0 {
		r.catalogue[key] = nil
	}
	return nil
}

func (r *registry) add(flag core.Flag) error {
	if err := r.addLongNameIfValid(flag.LongName()); err != nil {
		return err
	}

	if err := r.addShortNameIfValid(flag.ShortName()); err != nil {
		return err
	}

	return r.addKeyIfValid(flag.Key().Get())
}

func (r *registry) isRegistered(arg string) bool {
	_, ok := r.catalogue[arg]
	return ok
}

func (r *registry) isReserved(name string) bool {
	_, ok := reserved[name]
	return ok
}

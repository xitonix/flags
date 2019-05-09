package flags

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

func (r *registry) add(flag Flag) error {
	long := flag.Name()

	if isEmpty(long) {
		return ErrEmptyFlagName
	}

	if r.isReserved(long) {
		return errInvalidFlag(long, "", "is a reserved flag")
	}

	short := flag.Short()
	if r.isReserved(short) {
		return errInvalidFlag("", short, "is a reserved flag")
	}

	if r.isLongNameRegistered(long) {
		return errInvalidFlag(long, "", "flag already exists")
	}

	if r.isShortNameRegistered(short) {
		return errInvalidFlag("", short, "flag already exists")
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

func (r *registry) isLongNameRegistered(name string) bool {
	_, ok := r.longs[name]
	return ok
}

func (r *registry) isShortNameRegistered(name string) bool {
	_, ok := r.shorts[name]
	return ok
}

func (r *registry) isReserved(name string) bool {
	_, ok := reserved[name]
	return ok
}

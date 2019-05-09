package flags

type ErrInvalidFlag struct {
	long, short string
	msg         string
}

func errInvalidFlag(long, short, msg string) *ErrInvalidFlag {
	return &ErrInvalidFlag{
		long:  long,
		short: short,
		msg:   msg,
	}
}

func (e *ErrInvalidFlag) FieldName() string {
	return e.long
}

func (e *ErrInvalidFlag) Error() string {
	var str string
	if !isEmpty(e.long) {
		str += "--" + e.long
	}
	if !isEmpty(e.short) {
		var comma string
		if !isEmpty(str) {
			comma = ", "
		}
		str += comma + "-" + e.short
	}
	return str + " " + e.msg
}

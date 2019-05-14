package core

import "go.xitonix.io/flags/internal"

type ErrInvalidFlag struct {
	long, short string
	msg         string
}

func NewInvalidFlagErr(long, short, msg string) *ErrInvalidFlag {
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
	if !internal.IsEmpty(e.long) {
		str += "--" + e.long
	}
	if !internal.IsEmpty(e.short) {
		var comma string
		if !internal.IsEmpty(str) {
			comma = ", "
		}
		str += comma + "-" + e.short
	}
	return str + " " + e.msg
}

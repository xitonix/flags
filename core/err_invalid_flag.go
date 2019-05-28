package core

import "go.xitonix.io/flags/internal"

type ErrInvalidFlag struct {
	long, short, key string
	msg              string
}

func NewInvalidFlagErr(long, short, key, msg string) *ErrInvalidFlag {
	return &ErrInvalidFlag{
		long:  long,
		short: short,
		key:   key,
		msg:   msg,
	}
}

func (e *ErrInvalidFlag) FieldName() string {
	return e.long
}

func (e *ErrInvalidFlag) Error() string {
	var str string
	if !internal.IsEmpty(e.long) {
		str += e.long
	}
	if !internal.IsEmpty(e.short) {
		var comma string
		if !internal.IsEmpty(str) {
			comma = ", "
		}
		str += comma + e.short
	}
	if !internal.IsEmpty(e.key) {
		str = e.key
	}
	return str + " " + e.msg
}

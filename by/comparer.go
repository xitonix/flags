package by

import "github.com/xitonix/flags/core"

// Comparer is the interface for comparing two flags.
//
// An implementation of this interface can be used to sort the flags in the help output.
//
// Checkout pre-built implementations such as by.LongNameAscending, by.KeyAscending etc
type Comparer interface {
	LessThan(f1, f2 core.Flag) bool
}

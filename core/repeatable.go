package core

// Repeatable is the interface for a special type of numeric flags when the flags' short or long names can be provided
// multiple times by the command line arguments.
//
// Each repeat of the associated names, will add 'Once()' to the final value of the flag.
//
// A good example would be the CounterFlag type which implements the Repeatable interface with the return value of 1 for
// its Once() method. If a flag of the Counter type is defined with -c, --count names, providing '-ccc' as command line
// argument is expected to set the final value of the flag to 3. '--counter -cc', '--counter --counter --counter' or
// '-c --counter --counter' would result in the same final value.
type Repeatable interface {
	Once() int
}

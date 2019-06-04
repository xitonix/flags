package core

// Logger is the interface that adds logging functionality to the package.
//
// A logger is only used to print out the runtime errors.
type Logger interface {
	Print(error)
}

package config

import (
	"fmt"
)

// Logger represent the default logger which writes to standard output.
type Logger struct{}

// Print prints the error to standard output.
func (*Logger) Print(err error) {
	fmt.Println(err.Error())
}

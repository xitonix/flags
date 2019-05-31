package core

// HelpFormatter is an interface that contains a method for formatting a flag in
// order to provide a string to --help or -h requests.
type HelpFormatter interface {
	Format(f Flag, deprecationMark, defaultValueFormatString string) string
}

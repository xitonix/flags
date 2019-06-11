package core

// EmptyValueProvider is an interface which provides a fallback value when none of the sources has provided a
// none-empty value for a flag.
//
// This is different to Default values in which none of the sources provides ANY value.
type EmptyValueProvider interface {
	EmptyValue() string
}

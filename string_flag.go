package flags

import (
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// StringFlag represents a string flag
type StringFlag struct {
	key                 *data.Key
	defaultValue, value string
	hasDefault          bool
	ptr                 *string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
}

func newString(name, usage, short string) *StringFlag {
	ptr := new(string)
	return &StringFlag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   ptr,
	}
}

// LongName returns the long name of the flag (ie. --port).
//
// Long name is lower case and separated by hyphens (ie. --port-number)
func (f *StringFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed when the documentation is requested by the user (--help/-h)
func (f *StringFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *StringFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the the flag's type.
//
// This will be printed in the help output.
func (f *StringFlag) Type() string {
	return "string"
}

// ShortName returns the short name of the flag (ie. -p)
//
// Short name is case sensitive
func (f *StringFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *StringFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by any specified source.
func (f *StringFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying value.
//
// You can also use the Get() method instead.
func (f *StringFlag) Var() *string {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *StringFlag) Get() string {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values,
// defined at bucket level (if enabled)
func (f *StringFlag) WithKey(keyID string) *StringFlag {
	f.key.Set(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value for this flag, the default value will be returned
// by Get(). The same value will also be assigned to the underlying pointer (accessible through Var())
func (f *StringFlag) WithDefault(defaultValue string) *StringFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *StringFlag) Hide() *StringFlag {
	f.isHidden = true
	return f
}

// MarkAsDeprecated marks the flag as deprecated.
//
// A deprecated flag will be marked in the help output to draw users' attention.
// The default deprecation mark (config.DeprecatedFlagIndicatorDefault) can be overridden by calling
// the bucket's WithDeprecationMark(...) method
func (f *StringFlag) MarkAsDeprecated() *StringFlag {
	f.isDeprecated = true
	return f
}

// Set sets the value of this flag
func (f *StringFlag) Set(value string) error {
	f.isSet = true
	f.set(value)
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method
func (f *StringFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *StringFlag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	if f.defaultValue == "" {
		return "''"
	}
	return f.defaultValue
}

// Key returns the current key of the flag.
//
// Each flag within a bucket may have an optional unique key which can be used to retrieve its value
// from different sources.
func (f *StringFlag) Key() *data.Key {
	return f.key
}

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

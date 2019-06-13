package flags

import (
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// StringSliceFlag represents a StringSlice flag
type StringSliceFlag struct {
	key                 *data.Key
	defaultValue, value []string
	hasDefault          bool
	ptr                 *[]string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	delimiter           string
}

func newStringSlice(name, usage, short string) *StringSliceFlag {
	ptr := new([]string)
	return &StringSliceFlag{
		key:       &data.Key{},
		short:     internal.SanitiseShortName(short),
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       ptr,
		value:     make([]string, 0),
		delimiter: core.DefaultSliceDelimiter,
	}
}

// LongName returns the long name of the flag (i.e. --file).
//
// Long name is case insensitive and always lower case (i.e. --file-path).
func (f *StringSliceFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *StringSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *StringSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *StringSliceFlag) Type() string {
	return "[]string"
}

// ShortName returns the flag's short name (i.e. -p).
//
// Short name is a single case sensitive character.
func (f *StringSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *StringSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *StringSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *StringSliceFlag) Var() *[]string {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *StringSliceFlag) Get() []string {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *StringSliceFlag) WithKey(keyID string) *StringSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *StringSliceFlag) WithDefault(defaultValue []string) *StringSliceFlag {
	if defaultValue != nil {
		f.defaultValue = defaultValue
		f.hasDefault = true
	}
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *StringSliceFlag) Hide() *StringSliceFlag {
	f.isHidden = true
	return f
}

// MarkAsDeprecated marks the flag as deprecated.
//
// A deprecated flag will be marked in the help output to draw the users' attention.
// The default deprecation mark (config.DeprecatedFlagIndicatorDefault) can be overridden by configuring the bucket.
//
// Example:
//
// 	flags.SetDeprecationMark("**DEPRECATED**")
//  OR
//	bucket := flags.NewBucket(config.WithDeprecationMark("**DEPRECATED**"))
func (f *StringSliceFlag) MarkAsDeprecated() *StringSliceFlag {
	f.isDeprecated = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *StringSliceFlag) WithDelimiter(delimiter string) *StringSliceFlag {
	if len(delimiter) == 0 {
		delimiter = core.DefaultSliceDelimiter
	}
	f.delimiter = delimiter
	return f
}

// Set sets the flag value.
func (f *StringSliceFlag) Set(value string) error {
	var parts []string
	if len(value) == 0 {
		parts = make([]string, 0)
	} else {
		parts = strings.Split(value, f.delimiter)
	}
	f.set(parts)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *StringSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *StringSliceFlag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	return f.defaultValue
}

// Key returns the current key of the flag.
//
// Each flag within a bucket may have an optional UNIQUE key which will be used to retrieve its value
// from different sources. This is the key which will be used internally to retrieve the flag's value
// from the environment variables.
func (f *StringSliceFlag) Key() *data.Key {
	return f.key
}

func (f *StringSliceFlag) set(value []string) {
	f.value = value
	*f.ptr = value
}

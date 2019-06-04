package flags

import (
	"fmt"
	"strconv"
	"strings"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// Int64Flag represents an int64 flag
type Int64Flag struct {
	key                 *data.Key
	defaultValue, value int64
	hasDefault          bool
	ptr                 *int64
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
}

func newInt64(name, usage, short string) *Int64Flag {
	ptr := new(int64)
	return &Int64Flag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   ptr,
	}
}

// LongName returns the long name of the flag (ie. --port).
//
// Long name is case insensitive and always lower case (ie. --port-number).
func (f *Int64Flag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be print64ed in the help output.
func (f *Int64Flag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *Int64Flag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be print64ed in the help output.
func (f *Int64Flag) Type() string {
	return "int64"
}

// ShortName returns the flag's short name (ie. -p).
//
// Short name is a single case sensitive character.
func (f *Int64Flag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be print64ed in the help output.
func (f *Int64Flag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *Int64Flag) IsSet() bool {
	return f.isSet
}

// Var returns a point64er to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *Int64Flag) Var() *int64 {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *Int64Flag) Get() int64 {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values,
// defined at bucket level (if enabled).
func (f *Int64Flag) WithKey(keyID string) *Int64Flag {
	f.key.Set(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *Int64Flag) WithDefault(defaultValue int64) *Int64Flag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *Int64Flag) Hide() *Int64Flag {
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
func (f *Int64Flag) MarkAsDeprecated() *Int64Flag {
	f.isDeprecated = true
	return f
}

// Set sets the value of this flag.
func (f *Int64Flag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = "0"
	}
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("%s is not a valid %s value for --%s", value, f.Type(), f.long)
	}
	f.set(int64(v))
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *Int64Flag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *Int64Flag) Default() interface{} {
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
func (f *Int64Flag) Key() *data.Key {
	return f.key
}

func (f *Int64Flag) set(value int64) {
	f.value = value
	*f.ptr = value
}
package flags

import (
	"fmt"
	"strings"
	"time"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// TimeFlag represents a Time flag
type TimeFlag struct {
	key                 *data.Key
	defaultValue, value time.Time
	hasDefault          bool
	ptr                 *time.Time
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
}

func newTime(name, usage, short string) *TimeFlag {
	ptr := new(time.Time)
	return &TimeFlag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   ptr,
	}
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (ie. --port-number).
func (f *TimeFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *TimeFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *TimeFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *TimeFlag) Type() string {
	return "time"
}

// ShortName returns the flag's short name (ie. -p).
//
// Short name is a single case sensitive character.
func (f *TimeFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *TimeFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *TimeFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *TimeFlag) Var() *time.Time {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *TimeFlag) Get() time.Time {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *TimeFlag) WithKey(keyID string) *TimeFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *TimeFlag) WithDefault(defaultValue time.Time) *TimeFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *TimeFlag) Hide() *TimeFlag {
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
func (f *TimeFlag) MarkAsDeprecated() *TimeFlag {
	f.isDeprecated = true
	return f
}

// Set sets the value of this flag.
//
// A time string must be in `2006-01-02T15:04:05.999999999Z07:00` format
func (f *TimeFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = time.Time{}.Format(time.RFC3339Nano)
	}
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid %s value for --%s", value, f.Type(), f.long)
	}
	f.set(t)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *TimeFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *TimeFlag) Default() interface{} {
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
func (f *TimeFlag) Key() *data.Key {
	return f.key
}

func (f *TimeFlag) set(value time.Time) {
	f.value = value
	*f.ptr = value
}

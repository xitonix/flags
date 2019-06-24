package flags

import (
	"strings"
	"time"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// DurationSliceFlag represents a Duration slice flag.
//
// The value of a Duration slice flag can be set using a comma (or any custom delimiter) separated string of durations.
//
// Each duration string is a possibly signed sequence of decimal numbers, with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// For example --durations "2s, 2.5s, 5s".
//
// A custom delimiter string can be defined using WithDelimiter() method.
type DurationSliceFlag struct {
	key                 *data.Key
	defaultValue, value []time.Duration
	hasDefault          bool
	ptr                 *[]time.Duration
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	delimiter           string
	validate            func(in time.Duration) error
	validM              map[time.Duration]interface{}
	valid               string
}

func newDurationSlice(name, usage, short string) *DurationSliceFlag {
	f := &DurationSliceFlag{
		key:       &data.Key{},
		short:     internal.SanitiseShortName(short),
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]time.Duration),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(make([]time.Duration, 0))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --durations).
func (f *DurationSliceFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *DurationSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *DurationSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *DurationSliceFlag) Type() string {
	return "[]duration"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *DurationSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *DurationSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *DurationSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *DurationSliceFlag) Var() *[]time.Duration {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *DurationSliceFlag) Get() []time.Duration {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *DurationSliceFlag) WithKey(keyID string) *DurationSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *DurationSliceFlag) WithDefault(defaultValue []time.Duration) *DurationSliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *DurationSliceFlag) Hide() *DurationSliceFlag {
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
// 	bucket := flags.NewBucket(config.WithDeprecationMark("**DEPRECATED**"))
func (f *DurationSliceFlag) MarkAsDeprecated() *DurationSliceFlag {
	f.isDeprecated = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *DurationSliceFlag) WithDelimiter(delimiter string) *DurationSliceFlag {
	if len(delimiter) == 0 {
		delimiter = core.DefaultSliceDelimiter
	}
	f.delimiter = delimiter
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *DurationSliceFlag) WithValidationCallback(validate func(in time.Duration) error) *DurationSliceFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *DurationSliceFlag) WithValidRange(valid ...time.Duration) *DurationSliceFlag {
	l := len(valid)
	if len(valid) == 0 {
		return f
	}
	f.validM = make(map[time.Duration]interface{})
	for i, v := range valid {
		f.valid += internal.GetExpectedValueString(v, i, l)
		f.validM[v] = nil
	}
	return f
}

// Set sets the flag value.
//
// The value of a Duration slice flag can be set using a comma (or any custom delimiter) separated string of durations.
//
// Each duration string is a possibly signed sequence of decimal numbers, with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// For example --durations "2s, 2.5s, 5s".
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (f *DurationSliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]time.Duration, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}
		item, err := time.ParseDuration(value)
		if err != nil {
			return internal.InvalidValueErr(value, f.long, f.Type())
		}
		list = append(list, item)
	}

	if f.validate != nil {
		for _, item := range list {
			err := f.validate(item)
			if err != nil {
				return err
			}
		}
	}

	// Validation callback takes priority over validation list
	if f.validate == nil && f.validM != nil {
		for _, item := range list {
			if _, ok := f.validM[item]; !ok {
				return internal.OutOfRangeErr(value, f.long, f.valid, len(f.validM))
			}
		}
	}

	f.set(list)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *DurationSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *DurationSliceFlag) Default() interface{} {
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
func (f *DurationSliceFlag) Key() *data.Key {
	return f.key
}

func (f *DurationSliceFlag) set(value []time.Duration) {
	f.value = value
	*f.ptr = value
}

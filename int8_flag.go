package flags

import (
	"strconv"
	"strings"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// Int8Flag represents an int8 flag
type Int8Flag struct {
	key                 *data.Key
	defaultValue, value int8
	hasDefault          bool
	ptr                 *int8
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	validate            func(in int8) error
	validM              map[int8]interface{}
	valid               string
}

func newInt8(name, usage, short string) *Int8Flag {
	f := &Int8Flag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   new(int8),
	}
	f.set(0)
	return f
}

// LongName returns the long name of the flag..
//
// Long name is case insensitive and always lower case (i.e. --port-number).
func (f *Int8Flag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *Int8Flag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *Int8Flag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *Int8Flag) Type() string {
	return "int8"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *Int8Flag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *Int8Flag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *Int8Flag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *Int8Flag) Var() *int8 {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *Int8Flag) Get() int8 {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *Int8Flag) WithKey(keyID string) *Int8Flag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *Int8Flag) WithDefault(defaultValue int8) *Int8Flag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *Int8Flag) Hide() *Int8Flag {
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
func (f *Int8Flag) MarkAsDeprecated() *Int8Flag {
	f.isDeprecated = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *Int8Flag) WithValidationCallback(validate func(in int8) error) *Int8Flag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *Int8Flag) WithValidRange(valid ...int8) *Int8Flag {
	l := len(valid)
	if l == 0 {
		return f
	}
	f.validM = make(map[int8]interface{})
	for i, v := range valid {
		f.valid += internal.GetExpectedValueString(v, i, l)
		f.validM[v] = nil
	}
	return f
}

// Set sets the flag value.
func (f *Int8Flag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = "0"
	}
	v, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return internal.InvalidValueErr(value, f.long, f.Type())
	}

	if f.validate != nil {
		err := f.validate(int8(v))
		if err != nil {
			return err
		}
	}

	// Validation callback takes priority over validation list
	if f.validate == nil && f.validM != nil {
		if _, ok := f.validM[int8(v)]; !ok {
			return internal.OutOfRangeErr(value, f.long, f.valid, len(f.validM))
		}
	}

	f.set(int8(v))
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *Int8Flag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *Int8Flag) Default() interface{} {
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
func (f *Int8Flag) Key() *data.Key {
	return f.key
}

func (f *Int8Flag) set(value int8) {
	f.value = value
	*f.ptr = value
}

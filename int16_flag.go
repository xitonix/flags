package flags

import (
	"strconv"
	"strings"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// Int16Flag represents an int16 flag
type Int16Flag struct {
	key                 *data.Key
	defaultValue, value int16
	hasDefault          bool
	ptr                 *int16
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	validate            func(in int16) error
	validM              map[int16]interface{}
	valid               string
}

func newInt16(name, usage, short string) *Int16Flag {
	f := &Int16Flag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   new(int16),
	}
	f.set(0)
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --port-number).
func (f *Int16Flag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *Int16Flag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *Int16Flag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *Int16Flag) Type() string {
	return "int16"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *Int16Flag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *Int16Flag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *Int16Flag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *Int16Flag) Var() *int16 {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *Int16Flag) Get() int16 {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *Int16Flag) WithKey(keyID string) *Int16Flag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *Int16Flag) WithDefault(defaultValue int16) *Int16Flag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *Int16Flag) Hide() *Int16Flag {
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
func (f *Int16Flag) MarkAsDeprecated() *Int16Flag {
	f.isDeprecated = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *Int16Flag) WithValidationCallback(validate func(in int16) error) *Int16Flag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *Int16Flag) WithValidRange(valid ...int16) *Int16Flag {
	l := len(valid)
	if l == 0 {
		return f
	}
	f.validM = make(map[int16]interface{})
	for i, v := range valid {
		f.valid += internal.GetExpectedValueString(v, i, l)
		f.validM[v] = nil
	}
	return f
}

// Set sets the flag value.
func (f *Int16Flag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = "0"
	}
	v, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return internal.InvalidValueErr(value, f.long, f.Type())
	}

	if f.validate != nil {
		err := f.validate(int16(v))
		if err != nil {
			return err
		}
	}

	// Validation callback takes priority over validation list
	if f.validate == nil && f.validM != nil {
		if _, ok := f.validM[int16(v)]; !ok {
			return internal.OutOfRangeErr(value, f.long, f.valid, len(f.validM))
		}
	}

	f.set(int16(v))
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *Int16Flag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *Int16Flag) Default() interface{} {
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
func (f *Int16Flag) Key() *data.Key {
	return f.key
}

func (f *Int16Flag) set(value int16) {
	f.value = value
	*f.ptr = value
}

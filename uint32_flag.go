package flags

import (
	"strconv"
	"strings"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// UInt32Flag represents an uint32 flag.
type UInt32Flag struct {
	key                 *data.Key
	defaultValue, value uint32
	hasDefault          bool
	ptr                 *uint32
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	validate            func(in uint32) error
	validationList      map[uint32]interface{}
	acceptableItems     []string
}

func newUInt32(name, usage, short string) *UInt32Flag {
	f := &UInt32Flag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   new(uint32),
	}
	f.set(0)
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --port-number).
func (f *UInt32Flag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *UInt32Flag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *UInt32Flag) IsDeprecated() bool {
	return f.isDeprecated
}

// IsRequired returns true if the flag value must be provided.
func (f *UInt32Flag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *UInt32Flag) Required() *UInt32Flag {
	f.isRequired = true
	return f
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *UInt32Flag) Type() string {
	return "uint32"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *UInt32Flag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *UInt32Flag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *UInt32Flag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *UInt32Flag) Var() *uint32 {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *UInt32Flag) Get() uint32 {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *UInt32Flag) WithKey(keyID string) *UInt32Flag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *UInt32Flag) WithDefault(defaultValue uint32) *UInt32Flag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *UInt32Flag) Hide() *UInt32Flag {
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
func (f *UInt32Flag) MarkAsDeprecated() *UInt32Flag {
	f.isDeprecated = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *UInt32Flag) WithValidationCallback(validate func(in uint32) error) *UInt32Flag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *UInt32Flag) WithValidRange(valid ...uint32) *UInt32Flag {
	if len(valid) == 0 {
		return f
	}
	f.validationList = make(map[uint32]interface{})
	f.acceptableItems = make([]string, 0)
	for _, v := range valid {
		if _, ok := f.validationList[v]; !ok {
			f.validationList[v] = nil
			f.acceptableItems = append(f.acceptableItems, strconv.FormatUint(uint64(v), 10))
		}
	}
	return f
}

// Set sets the flag value.
func (f *UInt32Flag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = "0"
	}
	v, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return internal.InvalidValueErr(value, f.long, f.short, f.Type())
	}

	if f.validate != nil {
		err := f.validate(uint32(v))
		if err != nil {
			return err
		}
	}

	// Validation callback takes priority over validation list
	if f.validate == nil && len(f.validationList) > 0 {
		if _, ok := f.validationList[uint32(v)]; !ok {
			return internal.OutOfRangeErr(value, f.long, f.short, f.acceptableItems)
		}
	}

	f.set(uint32(v))
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *UInt32Flag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *UInt32Flag) Default() interface{} {
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
func (f *UInt32Flag) Key() *data.Key {
	return f.key
}

func (f *UInt32Flag) set(value uint32) {
	f.value = value
	*f.ptr = value
}

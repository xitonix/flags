package flags

import (
	"fmt"
	"strconv"
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// Float64SliceFlag represents a Float64Slice flag.
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
type Float64SliceFlag struct {
	key                 *data.Key
	defaultValue, value []float64
	hasDefault          bool
	ptr                 *[]float64
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	delimiter           string
	validate            func(in float64) error
	validationList      map[float64]interface{}
	acceptableItems     []string
}

func newFloat64Slice(name, usage, short string) *Float64SliceFlag {
	f := &Float64SliceFlag{
		key:       &data.Key{},
		short:     internal.SanitiseShortName(short),
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]float64),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(make([]float64, 0))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --rates).
func (f *Float64SliceFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *Float64SliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *Float64SliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *Float64SliceFlag) Type() string {
	return "[]float64"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *Float64SliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *Float64SliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *Float64SliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *Float64SliceFlag) Var() *[]float64 {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *Float64SliceFlag) Get() []float64 {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *Float64SliceFlag) WithKey(keyID string) *Float64SliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *Float64SliceFlag) WithDefault(defaultValue []float64) *Float64SliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *Float64SliceFlag) Hide() *Float64SliceFlag {
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
func (f *Float64SliceFlag) MarkAsDeprecated() *Float64SliceFlag {
	f.isDeprecated = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *Float64SliceFlag) WithDelimiter(delimiter string) *Float64SliceFlag {
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
func (f *Float64SliceFlag) WithValidationCallback(validate func(in float64) error) *Float64SliceFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *Float64SliceFlag) WithValidRange(valid ...float64) *Float64SliceFlag {
	if len(valid) == 0 {
		return f
	}
	f.validationList = make(map[float64]interface{})
	f.acceptableItems = make([]string, 0)
	for _, v := range valid {
		if _, ok := f.validationList[v]; !ok {
			f.validationList[v] = nil
			f.acceptableItems = append(f.acceptableItems, fmt.Sprintf("%v", v))
		}
	}
	return f
}

// Set sets the flag value.
//
// The value of a Float64Slice flag can be set using a comma (or any custom delimiter) separated string of integers.
// For example --rates "1.0, 1.5, 3.0, 3.5, 5.0"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (f *Float64SliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]float64, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}
		item, err := strconv.ParseFloat(value, 64)
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
	if f.validate == nil && len(f.validationList) > 0 {
		for _, item := range list {
			if _, ok := f.validationList[item]; !ok {
				return internal.OutOfRangeErr(value, f.long, f.acceptableItems)
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
func (f *Float64SliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *Float64SliceFlag) Default() interface{} {
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
func (f *Float64SliceFlag) Key() *data.Key {
	return f.key
}

func (f *Float64SliceFlag) set(value []float64) {
	f.value = value
	*f.ptr = value
}

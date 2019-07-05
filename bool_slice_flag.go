package flags

import (
	"strconv"
	"strings"

	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/data"
	"github.com/xitonix/flags/internal"
)

// BoolSliceFlag represents a boolean slice flag.
//
// The value of a boolean slice flag can be set using a comma (or any custom delimiter) separated string of true, false, 0 or 1.
// For example --bits "0, 1, true, false".
//
// A custom delimiter string can be defined using WithDelimiter() method.
type BoolSliceFlag struct {
	key                 *data.Key
	defaultValue, value []bool
	hasDefault          bool
	ptr                 *[]bool
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	delimiter           string
	validate            func(in bool) error
}

func newBoolSlice(name, usage string) *BoolSliceFlag {
	f := &BoolSliceFlag{
		key:       &data.Key{},
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]bool),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(make([]bool, 0))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --bits).
func (f *BoolSliceFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *BoolSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *BoolSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// WithShort sets the short name of the flag.
//
// The short name is a single case sensitive character (i.e. -B).
func (f *BoolSliceFlag) WithShort(short string) *BoolSliceFlag {
	f.short = internal.SanitiseShortName(short)
	return f
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *BoolSliceFlag) Type() string {
	return "[]bool"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -B).
func (f *BoolSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *BoolSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *BoolSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *BoolSliceFlag) Var() *[]bool {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *BoolSliceFlag) Get() []bool {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *BoolSliceFlag) WithKey(keyID string) *BoolSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *BoolSliceFlag) WithDefault(defaultValue []bool) *BoolSliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *BoolSliceFlag) Hide() *BoolSliceFlag {
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
func (f *BoolSliceFlag) MarkAsDeprecated() *BoolSliceFlag {
	f.isDeprecated = true
	return f
}

// IsRequired returns true if the flag value must be provided.
func (f *BoolSliceFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *BoolSliceFlag) Required() *BoolSliceFlag {
	f.isRequired = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *BoolSliceFlag) WithDelimiter(delimiter string) *BoolSliceFlag {
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
func (f *BoolSliceFlag) WithValidationCallback(validate func(in bool) error) *BoolSliceFlag {
	f.validate = validate
	return f
}

// Set sets the flag value.
//
// The value of a boolean slice flag can be set using a comma (or any custom delimiter) separated string of true, false, 0 or 1.
// For example --bits "0, 1, true, false"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (f *BoolSliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]bool, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}

		item, err := strconv.ParseBool(value)
		if err != nil {
			return internal.InvalidValueErr(value, f.long, f.short, f.Type())
		}

		if f.validate != nil {
			err := f.validate(item)
			if err != nil {
				return err
			}
		}
		list = append(list, item)
	}

	f.set(list)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *BoolSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *BoolSliceFlag) Default() interface{} {
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
func (f *BoolSliceFlag) Key() *data.Key {
	return f.key
}

func (f *BoolSliceFlag) set(value []bool) {
	f.value = value
	*f.ptr = value
}

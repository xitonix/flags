package core

import (
	"strconv"
	"strings"

	"github.com/xitonix/flags/internal"
)

// UIntSliceFlag represents an uint flag.
//
// The value of a uint slice flag can be set using a comma (or any custom delimiter) separated string of unsigned integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
type UIntSliceFlag struct {
	key                 *Key
	defaultValue, value []uint
	hasDefault          bool
	ptr                 *[]uint
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	delimiter           string
	validate            func(in uint) error
	validationList      map[uint]interface{}
	acceptableItems     []string
}

// NewUInt creates a new uint slice flag.
func NewUIntSlice(name, usage string) *UIntSliceFlag {
	f := &UIntSliceFlag{
		key:       &Key{},
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]uint),
		delimiter: DefaultDelimiter,
	}
	f.set(make([]uint, 0))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --numbers).
func (f *UIntSliceFlag) LongName() string {
	return f.long
}

// WithShort sets the short name of the flag.
//
// The short name is a single case sensitive character (i.e. -n).
func (f *UIntSliceFlag) WithShort(short string) *UIntSliceFlag {
	f.short = internal.SanitiseShortName(short)
	return f
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *UIntSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *UIntSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *UIntSliceFlag) Type() string {
	return "[]uint"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -N).
func (f *UIntSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *UIntSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *UIntSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *UIntSliceFlag) Var() *[]uint {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *UIntSliceFlag) Get() []uint {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *UIntSliceFlag) WithKey(keyID string) *UIntSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *UIntSliceFlag) WithDefault(defaultValue []uint) *UIntSliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *UIntSliceFlag) Hide() *UIntSliceFlag {
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
func (f *UIntSliceFlag) MarkAsDeprecated() *UIntSliceFlag {
	f.isDeprecated = true
	return f
}

// IsRequired returns true if the flag value must be provided.
func (f *UIntSliceFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *UIntSliceFlag) Required() *UIntSliceFlag {
	f.isRequired = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultDelimiter).
func (f *UIntSliceFlag) WithDelimiter(delimiter string) *UIntSliceFlag {
	if len(delimiter) == 0 {
		delimiter = DefaultDelimiter
	}
	f.delimiter = delimiter
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *UIntSliceFlag) WithValidationCallback(validate func(in uint) error) *UIntSliceFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *UIntSliceFlag) WithValidRange(valid ...uint) *UIntSliceFlag {
	if len(valid) == 0 {
		return f
	}
	f.validationList = make(map[uint]interface{})
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
//
// The value of a uint slice flag can be set using a comma (or any custom delimiter) separated string of unsigned integers.
// For example --numbers "1,8,70,60,100"
//
// A custom delimiter string can be defined using WithDelimiter() method.
func (f *UIntSliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]uint, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}
		item, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return internal.InvalidValueErr(value, f.long, f.short, f.Type())
		}

		if f.validate != nil {
			err := f.validate(uint(item))
			if err != nil {
				return err
			}
		}

		// Validation callback takes priority over validation list
		if f.validate == nil && len(f.validationList) > 0 {
			if _, ok := f.validationList[uint(item)]; !ok {
				return internal.OutOfRangeErr(value, f.long, f.short, f.acceptableItems)
			}
		}

		list = append(list, uint(item))
	}

	f.set(list)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *UIntSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *UIntSliceFlag) Default() interface{} {
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
func (f *UIntSliceFlag) Key() *Key {
	return f.key
}

func (f *UIntSliceFlag) set(value []uint) {
	f.value = value
	*f.ptr = value
}

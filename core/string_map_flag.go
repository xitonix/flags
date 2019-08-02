package core

import (
	"strings"

	"github.com/xitonix/flags/internal"
)

// StringMapFlag represents a string map flag.
//
// The value of a string map flag can be set using map initialisation literals.
// For example --mappings "key1:value1, key2:value2"
//
// By default, the leading and trailing white spaces will be automatically trimmed from each key/value pairs.
// With trimming enabled, "key1 : value1 , key2:  value2  " will be parsed into
// {"key1", "value1", "key2":"value2"} instead of {"key1 ", " value1 ", " key2":"  value2  "}.
// Notice that all the leading/trailing white space characters have been removed from all the keys and the values.
// Trimming can be disabled by calling the DisableKeyTrimming(), DisableValueTrimming() methods.
type StringMapFlag struct {
	key                 *Key
	defaultValue, value map[string]string
	hasDefault          bool
	ptr                 *map[string]string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	trimKey             bool
	trimValue           bool
	delimiter           string
	validate            func(key, value string) error
}

// NewStringMap creates a new string map flag.
func NewStringMap(name, usage string) *StringMapFlag {
	f := &StringMapFlag{
		key:       &Key{},
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new(map[string]string),
		trimKey:   true,
		trimValue: true,
		delimiter: DefaultDelimiter,
	}
	f.set(make(map[string]string))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --mappings).
func (f *StringMapFlag) LongName() string {
	return f.long
}

// WithShort sets the short name of the flag.
//
// The short name is a single case sensitive character (i.e. -m).
func (f *StringMapFlag) WithShort(short string) *StringMapFlag {
	f.short = internal.SanitiseShortName(short)
	return f
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *StringMapFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *StringMapFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *StringMapFlag) Type() string {
	return "[string]string"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -M).
func (f *StringMapFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *StringMapFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *StringMapFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *StringMapFlag) Var() *map[string]string {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *StringMapFlag) Get() map[string]string {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it. You can set the key to "-" to disable the auto generated ID (if there is one) for this flag.
func (f *StringMapFlag) WithKey(keyID string) *StringMapFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *StringMapFlag) WithDefault(defaultValue map[string]string) *StringMapFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *StringMapFlag) Hide() *StringMapFlag {
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
func (f *StringMapFlag) MarkAsDeprecated() *StringMapFlag {
	f.isDeprecated = true
	return f
}

// IsRequired returns true if the flag value must be provided.
func (f *StringMapFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *StringMapFlag) Required() *StringMapFlag {
	f.isRequired = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
func (f *StringMapFlag) WithValidationCallback(validate func(key, value string) error) *StringMapFlag {
	f.validate = validate
	return f
}

// Set sets the flag value.
//
// The value of a string map flag can be set using map initialisation literals.
// For example --mappings "key1:value1, key2:value2"
//
// By default, the leading and trailing white spaces will be automatically trimmed from each key/value pairs.
// With trimming enabled, "key1 : value1 , key2:  value2  " will be parsed into
// {"key1", "value1", "key2":"value2"} instead of {"key1 ", " value1 ", " key2":"  value2  "}.
// Notice that all the leading/trailing white space characters have been removed from all the keys and the values.
// Trimming can be disabled by calling the DisableKeyTrimming(), DisableValueTrimming() methods.
func (f *StringMapFlag) Set(value string) error {
	if internal.IsEmpty(value) {
		f.set(make(map[string]string))
		f.isSet = true
		return nil
	}

	mp := make(map[string]string)
	pairs := strings.Split(value, f.delimiter)
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return internal.InvalidValueErr(value, f.long, f.short, f.Type())
		}
		if f.trimKey {
			parts[0] = strings.TrimSpace(parts[0])
		}
		if f.trimValue {
			parts[1] = strings.TrimSpace(parts[1])
		}
		mp[parts[0]] = parts[1]
	}

	if f.validate != nil {
		for k, v := range mp {
			if err := f.validate(k, v); err != nil {
				return err
			}
		}
	}

	f.set(mp)
	f.isSet = true
	return nil
}

// WithDelimiter sets the delimiter for separating key/value pairs within the input string (Default: core.DefaultDelimiter).
func (f *StringMapFlag) WithDelimiter(delimiter string) *StringMapFlag {
	if len(delimiter) == 0 {
		delimiter = DefaultDelimiter
	}
	f.delimiter = delimiter
	return f
}

// DisableTrimming disables trimming the leading and trailing white space characters from each key in the map.
func (f *StringMapFlag) DisableKeyTrimming() *StringMapFlag {
	f.trimKey = false
	return f
}

// DisableValueTrimming disables trimming the leading and trailing white space characters from each value in the map.
func (f *StringMapFlag) DisableValueTrimming() *StringMapFlag {
	f.trimValue = false
	return f
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *StringMapFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *StringMapFlag) Default() interface{} {
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
func (f *StringMapFlag) Key() *Key {
	return f.key
}

func (f *StringMapFlag) set(value map[string]string) {
	f.value = value
	*f.ptr = value
}

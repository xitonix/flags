package flags

import (
	"encoding/json"
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// StringSliceMapFlag represents a map[string][]string flag.
//
// The value of a StringSliceMap flag can be set using standard map initialisation strings.
// Keys are strings and each value is a set of comma (or any custom delimiter) separated strings.
// For example --days '{"Week Days":"Mon,Tue,Wed,Thu,Fri", "Weekend":"Sat,Sun"}'
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
type StringSliceMapFlag struct {
	key                 *data.Key
	defaultValue, value map[string][]string
	hasDefault          bool
	ptr                 *map[string][]string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	validate            func(key string, value []string) error
	trimSpaces          bool
	delimiter           string
}

func newStringSliceMap(name, usage, short string) *StringSliceMapFlag {
	f := &StringSliceMapFlag{
		key:       &data.Key{},
		short:     internal.SanitiseShortName(short),
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new(map[string][]string),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(make(map[string][]string))
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --numbers).
func (f *StringSliceMapFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *StringSliceMapFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *StringSliceMapFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// WithTrimming enables trimming the leading and trailing white space characters from each list item.
func (f *StringSliceMapFlag) WithTrimming() *StringSliceMapFlag {
	f.trimSpaces = true
	return f
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *StringSliceMapFlag) Type() string {
	return "[string][]string"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *StringSliceMapFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *StringSliceMapFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *StringSliceMapFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *StringSliceMapFlag) Var() *map[string][]string {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *StringSliceMapFlag) Get() map[string][]string {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *StringSliceMapFlag) WithKey(keyID string) *StringSliceMapFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *StringSliceMapFlag) WithDefault(defaultValue map[string][]string) *StringSliceMapFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *StringSliceMapFlag) WithDelimiter(delimiter string) *StringSliceMapFlag {
	if len(delimiter) == 0 {
		delimiter = core.DefaultSliceDelimiter
	}
	f.delimiter = delimiter
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *StringSliceMapFlag) Hide() *StringSliceMapFlag {
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
func (f *StringSliceMapFlag) MarkAsDeprecated() *StringSliceMapFlag {
	f.isDeprecated = true
	return f
}

// IsRequired returns true if the flag value must be provided.
func (f *StringSliceMapFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *StringSliceMapFlag) Required() *StringSliceMapFlag {
	f.isRequired = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
func (f *StringSliceMapFlag) WithValidationCallback(validate func(key string, value []string) error) *StringSliceMapFlag {
	f.validate = validate
	return f
}

// Set sets the flag value.
//
// The value of a StringSliceMap flag can be set using standard map initialisation strings.
// Keys are strings and each value is a set of comma (or any custom delimiter) separated strings.
// For example --days '{"Week Days":"Mon,Tue,Wed,Thu,Fri", "Weekend":"Sat,Sun"}'
//
// A custom delimiter string can be defined using WithDelimiter() method.
//
// You can also trim the leading and trailing white spaces from each list item by enabling the feature
// using WithTrimming() method. With trimming enabled, "Sat, Sun" will be parsed into
// {"Sat", "Sun"} instead of {"Sat", " Sun"}.
// Notice that the leading white space before " Sun" has been removed.
func (f *StringSliceMapFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		value = "{}"
	}
	mp := make(map[string]string)
	err := json.Unmarshal([]byte(value), &mp)
	if err != nil {
		return internal.InvalidValueErr(value, f.long, f.short, f.Type())
	}

	result := make(map[string][]string)
	for k, v := range mp {
		parts := strings.Split(v, f.delimiter)
		if f.validate != nil {
			if err := f.validate(k, parts); err != nil {
				return err
			}
		}
		slice := make([]string, len(parts))
		for i, item := range parts {
			if f.trimSpaces {
				item = strings.TrimSpace(item)
			}
			slice[i] = item
		}
		result[k] = slice
	}

	f.set(result)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *StringSliceMapFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *StringSliceMapFlag) Default() interface{} {
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
func (f *StringSliceMapFlag) Key() *data.Key {
	return f.key
}

func (f *StringSliceMapFlag) set(value map[string][]string) {
	f.value = value
	*f.ptr = value
}

package flags

import (
	"fmt"
	"strings"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// StringFlag represents a string flag
type StringFlag struct {
	key                 *data.Key
	defaultValue, value string
	hasDefault          bool
	ptr                 *string
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	validate            func(in string) error
	validM              map[string]interface{}
	valid               []string
	ignoreCase          bool
}

func newString(name, usage, short string) *StringFlag {
	ptr := new(string)
	return &StringFlag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   ptr,
	}
}

// LongName returns the long name of the flag (i.e. --file).
//
// Long name is case insensitive and always lower case (i.e. --file-path).
func (f *StringFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *StringFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *StringFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *StringFlag) Type() string {
	return "string"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *StringFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *StringFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *StringFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *StringFlag) Var() *string {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *StringFlag) Get() string {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *StringFlag) WithKey(keyID string) *StringFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *StringFlag) WithDefault(defaultValue string) *StringFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *StringFlag) Hide() *StringFlag {
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
func (f *StringFlag) MarkAsDeprecated() *StringFlag {
	f.isDeprecated = true
	return f
}

func (f *StringFlag) WithValidationFunc(validate func(in string) error) *StringFlag {
	f.validate = validate
	return f
}

func (f *StringFlag) WithValidRange(valid []string, ignoreCase bool) *StringFlag {
	if len(valid) > 0 {
		f.ignoreCase = ignoreCase
		f.validM = make(map[string]interface{})
		f.valid = valid
		for _, v := range valid {
			item := v
			if ignoreCase {
				item = strings.ToLower(v)
			}
			f.validM[item] = nil
		}
	}
	return f
}

// Set sets the flag value.
func (f *StringFlag) Set(value string) error {
	if f.validate != nil {
		err := f.validate(value)
		if err != nil {
			return err
		}
	}

	if f.validate == nil && f.validM != nil {
		item := value
		if f.ignoreCase {
			item = strings.ToLower(item)
		}
		if _, ok := f.validM[item]; !ok {
			return fmt.Errorf("'%s' is not an acceptable value for --%s. Valid values are: %s", value, f.long, strings.Join(f.valid, ","))
		}
	}
	f.set(value)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *StringFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *StringFlag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	if f.defaultValue == "" {
		return "''"
	}
	return f.defaultValue
}

// Key returns the current key of the flag.
//
// Each flag within a bucket may have an optional UNIQUE key which will be used to retrieve its value
// from different sources. This is the key which will be used internally to retrieve the flag's value
// from the environment variables.
func (f *StringFlag) Key() *data.Key {
	return f.key
}

func (f *StringFlag) set(value string) {
	f.value = value
	*f.ptr = value
}

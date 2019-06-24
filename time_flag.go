package flags

import (
	"fmt"
	"strings"
	"time"

	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

var (
	layouts = []string{
		"02-01-2006T15:04:05.999999999",
		"02-01-2006T03:04:05.999999999PM",
		"02-01-2006T03:04:05.999999999 PM",
		"02-01-2006T03:04:05.999999999pm",
		"02-01-2006T03:04:05.999999999 pm",

		"02-01-2006 15:04:05.999999999",
		"02-01-2006 03:04:05.999999999PM",
		"02-01-2006 03:04:05.999999999 PM",
		"02-01-2006 03:04:05.999999999pm",
		"02-01-2006 03:04:05.999999999 pm",

		"02/01/2006T15:04:05.999999999",
		"02/01/2006T03:04:05.999999999PM",
		"02/01/2006T03:04:05.999999999 PM",
		"02/01/2006T03:04:05.999999999pm",
		"02/01/2006T03:04:05.999999999 pm",

		"02/01/2006 15:04:05.999999999",
		"02/01/2006 03:04:05.999999999PM",
		"02/01/2006 03:04:05.999999999 PM",
		"02/01/2006 03:04:05.999999999pm",
		"02/01/2006 03:04:05.999999999 pm",

		"02-01-2006",
		"02/01/2006",

		"Jan _2 15:04:05.999999999",
		"Jan _2 03:04:05.999999999PM",
		"Jan _2 03:04:05.999999999 PM",
		"Jan _2 03:04:05.999999999pm",
		"Jan _2 03:04:05.999999999 pm",

		"15:04:05.999999999",
		"03:04:05.999999999 PM",
		"03:04:05.999999999 pm",
		"03:04:05.999999999PM",
		"03:04:05.999999999pm",
	}
)

// TimeFlag represents a time flag.
//
// Supported layouts for time flags are:
//
// Full Date and Time
//
//  dd-MM-yyyyThh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980T14:22:20)
//  dd-MM-yyyy hh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980 14:22:20)
//  dd-MM-yyyyThh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980T02:22:20 PM)
//  dd-MM-yyyy hh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980 02:22:20 PM)
//
//  dd/MM/yyyyThh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyy hh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyyThh:mm:SS[.999999999] AM/PM
//  dd/MM/yyyy hh:mm:SS[.999999999] AM/PM
//
// Date
//
//  dd-MM-yyyy
//  dd/MM/yyyy
//
// Timestamp
//
//  MMM dd hh:mm:ss[.999999999] (24 hrs, i.e. Aug 27 14:22:20)
//  MMM dd hh:mm:ss[.999999999] AM/PM (i.e. Aug 27 02:22:20 PM)
//
// Time
//
//  hh:mm:ss[.999999999] (24 hrs, i.e. 14:22:20)
//  hh:mm:ss[.999999999] AM/PM (i.e. 02:22:20 PM)
//
// [.999999999] is the optional nano second component for time.
type TimeFlag struct {
	key                 *data.Key
	defaultValue, value time.Time
	hasDefault          bool
	ptr                 *time.Time
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isHidden            bool
	validate            func(in time.Time) error
	validM              map[time.Time]interface{}
	valid               []time.Time
}

func newTime(name, usage, short string) *TimeFlag {
	f := &TimeFlag{
		key:   &data.Key{},
		short: internal.SanitiseShortName(short),
		long:  internal.SanitiseLongName(name),
		usage: usage,
		ptr:   new(time.Time),
	}
	f.set(time.Time{})
	return f
}

// LongName returns the long name of the flag..
//
// Long name is case insensitive and always lower case (i.e. --port-number).
func (f *TimeFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *TimeFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *TimeFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *TimeFlag) Type() string {
	return "time"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -P).
func (f *TimeFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *TimeFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *TimeFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *TimeFlag) Var() *time.Time {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *TimeFlag) Get() time.Time {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *TimeFlag) WithKey(keyID string) *TimeFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *TimeFlag) WithDefault(defaultValue time.Time) *TimeFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *TimeFlag) Hide() *TimeFlag {
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
func (f *TimeFlag) MarkAsDeprecated() *TimeFlag {
	f.isDeprecated = true
	return f
}

// WithValidationCallback sets the validation callback function which will be called when the flag value is being set.
//
// The set operation will fail if the callback returns an error.
// You can also define a list of acceptable values using WithValidRange(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *TimeFlag) WithValidationCallback(validate func(in time.Time) error) *TimeFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *TimeFlag) WithValidRange(valid ...time.Time) *TimeFlag {
	l := len(valid)
	if l == 0 {
		return f
	}
	f.validM = make(map[time.Time]interface{})
	for _, v := range valid {
		f.valid = valid
		f.validM[v] = nil
	}
	return f
}

// Set sets the flag value.
//
// Supported layouts are:
//
// Full Date and Time
//
//  dd-MM-yyyyThh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980T14:22:20)
//  dd-MM-yyyy hh:mm:SS[.999999999] (24 hrs, i.e. 27-08-1980 14:22:20)
//  dd-MM-yyyyThh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980T02:22:20 PM)
//  dd-MM-yyyy hh:mm:SS[.999999999] AM/PM (i.e. 27-08-1980 02:22:20 PM)
//
//  dd/MM/yyyyThh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyy hh:mm:SS[.999999999] (24 hrs)
//  dd/MM/yyyyThh:mm:SS[.999999999] AM/PM
//  dd/MM/yyyy hh:mm:SS[.999999999] AM/PM
//
// Date
//
//  dd-MM-yyyy
//  dd/MM/yyyy
//
// Timestamp
//
//  MMM dd hh:mm:ss[.999999999] (24 hrs, i.e. Aug 27 14:22:20)
//  MMM dd hh:mm:ss[.999999999] AM/PM (i.e. Aug 27 02:22:20 PM)
//
// Time
//
//  hh:mm:ss[.999999999] (24 hrs, i.e. 14:22:20)
//  hh:mm:ss[.999999999] AM/PM (i.e. 02:22:20 PM)
//
// [.999999999] is the optional nano second component for time.
func (f *TimeFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		value = time.Time{}.Format("02-01-2006T15:04:05.999999999")
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {

			if f.validate != nil {
				err := f.validate(t)
				if err != nil {
					return err
				}
			}

			// Validation callback takes priority over validation list
			if f.validate == nil && f.validM != nil {
				if _, ok := f.validM[t]; !ok {
					return fmt.Errorf("%v is not an acceptable value for --%s. You must pick a value from %s.", t.Format(layout), f.long, f.getValidRangeString(layout))
				}
			}

			f.set(t)
			f.isSet = true
			return nil
		}
	}
	return internal.InvalidValueErr(value, f.long, f.Type())
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *TimeFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *TimeFlag) Default() interface{} {
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
func (f *TimeFlag) Key() *data.Key {
	return f.key
}

func (f *TimeFlag) set(value time.Time) {
	f.value = value
	*f.ptr = value
}

func (f *TimeFlag) getValidRangeString(layout string) string {
	v := make([]string, 0)
	for _, t := range f.valid {
		v = append(v, t.Format(layout))
	}
	return strings.Join(v, ",")
}

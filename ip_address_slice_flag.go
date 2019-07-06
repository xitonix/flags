package flags

import (
	"net"
	"strings"

	"github.com/xitonix/flags/core"
	"github.com/xitonix/flags/data"
	"github.com/xitonix/flags/internal"
)

// IPAddressSliceFlag represents an IP Address slice flag.
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
type IPAddressSliceFlag struct {
	key                 *data.Key
	defaultValue, value []net.IP
	hasDefault          bool
	ptr                 *[]net.IP
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	delimiter           string
	validate            func(in net.IP) error
	validationList      map[string]interface{}
	acceptableItems     []string
}

func newIPAddressSlice(name, usage string) *IPAddressSliceFlag {
	f := &IPAddressSliceFlag{
		key:       &data.Key{},
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]net.IP),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(nil)
	return f
}

// LongName returns the long name of the flag.
//
// Long name is case insensitive and always lower case (i.e. --endpoints).
func (f *IPAddressSliceFlag) LongName() string {
	return f.long
}

// WithShort sets the short name of the flag.
//
// The short name is a single case sensitive character (i.e. -e).
func (f *IPAddressSliceFlag) WithShort(short string) *IPAddressSliceFlag {
	f.short = internal.SanitiseShortName(short)
	return f
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *IPAddressSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *IPAddressSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// IsRequired returns true if the flag value must be provided.
func (f *IPAddressSliceFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *IPAddressSliceFlag) Required() *IPAddressSliceFlag {
	f.isRequired = true
	return f
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *IPAddressSliceFlag) Type() string {
	return "[]ip"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -E).
func (f *IPAddressSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *IPAddressSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *IPAddressSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *IPAddressSliceFlag) Var() *[]net.IP {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *IPAddressSliceFlag) Get() []net.IP {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *IPAddressSliceFlag) WithKey(keyID string) *IPAddressSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *IPAddressSliceFlag) WithDefault(defaultValue []net.IP) *IPAddressSliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *IPAddressSliceFlag) Hide() *IPAddressSliceFlag {
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
func (f *IPAddressSliceFlag) MarkAsDeprecated() *IPAddressSliceFlag {
	f.isDeprecated = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *IPAddressSliceFlag) WithDelimiter(delimiter string) *IPAddressSliceFlag {
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
func (f *IPAddressSliceFlag) WithValidationCallback(validate func(in net.IP) error) *IPAddressSliceFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *IPAddressSliceFlag) WithValidRange(valid ...net.IP) *IPAddressSliceFlag {
	if len(valid) == 0 {
		return f
	}
	f.validationList = make(map[string]interface{})
	f.acceptableItems = make([]string, 0)
	for _, v := range valid {
		if len(v) == 0 {
			continue
		}
		s := v.String()
		if _, ok := f.validationList[s]; !ok {
			f.acceptableItems = append(f.acceptableItems, s)
			f.validationList[s] = nil
		}
	}
	return f
}

// Set sets the flag value.
//
// The value of an IP address slice flag can be specified using a comma (or any custom delimiter) separated string of
// IPv4 (i.e. "192.0.2.1, 192.0.2.2") or IPv6 ("2001:db8::68, 2001:ab8::69") formatted strings.
// Different IP address versions can also be combined into a single string (i.e. "192.0.2.1, 2001:db8::68").
func (f *IPAddressSliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]net.IP, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}
		ip := net.ParseIP(value)
		if ip == nil {
			return internal.InvalidValueErr(value, f.long, f.short, f.Type())
		}

		if f.validate != nil {
			err := f.validate(ip)
			if err != nil {
				return err
			}
		}

		// Validation callback takes priority over validation list
		if f.validate == nil && len(f.validationList) > 0 {
			if _, ok := f.validationList[ip.String()]; !ok {
				return internal.OutOfRangeErr(value, f.long, f.short, f.acceptableItems)
			}
		}

		list = append(list, ip)
	}

	f.set(list)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *IPAddressSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil.
//
// The default value can be defined using WithDefault(...) method.
func (f *IPAddressSliceFlag) Default() interface{} {
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
func (f *IPAddressSliceFlag) Key() *data.Key {
	return f.key
}

func (f *IPAddressSliceFlag) set(value []net.IP) {
	f.value = value
	*f.ptr = value
}

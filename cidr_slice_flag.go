package flags

import (
	"strings"

	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/data"
	"go.xitonix.io/flags/internal"
)

// CIDRSliceFlag represents a CIDR slice flag
//
// The value of a CIDR slice flag can be defined using a list of CIDR notation IP addresses and prefix length,
// like "192.0.2.0/24, 2001:db8::/32", as defined in RFC 4632 and RFC 4291. Each item will be parsed to the
// address and the network implied by the IP and prefix length.
//
// For example, "192.0.2.1/24" will be translated to the IP address 192.0.2.1 and the network 192.0.2.0/24.
type CIDRSliceFlag struct {
	key                 *data.Key
	defaultValue, value []core.CIDR
	hasDefault          bool
	ptr                 *[]core.CIDR
	long, short         string
	usage               string
	isSet               bool
	isDeprecated        bool
	isRequired          bool
	isHidden            bool
	delimiter           string
	validate            func(in core.CIDR) error
	validationList      map[string]interface{}
	acceptableItems     []string
}

func newCIDRSlice(name, usage, short string) *CIDRSliceFlag {
	f := &CIDRSliceFlag{
		key:       &data.Key{},
		short:     internal.SanitiseShortName(short),
		long:      internal.SanitiseLongName(name),
		usage:     usage,
		ptr:       new([]core.CIDR),
		delimiter: core.DefaultSliceDelimiter,
	}
	f.set(nil)
	return f
}

// LongName returns the long name of the flag..
//
// Long name is case insensitive and always lower case (i.e. --networks).
func (f *CIDRSliceFlag) LongName() string {
	return f.long
}

// IsHidden returns true if the flag is hidden.
//
// A hidden flag won't be printed in the help output.
func (f *CIDRSliceFlag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is deprecated.
func (f *CIDRSliceFlag) IsDeprecated() bool {
	return f.isDeprecated
}

// IsRequired returns true if the flag value must be provided.
func (f *CIDRSliceFlag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *CIDRSliceFlag) Required() *CIDRSliceFlag {
	f.isRequired = true
	return f
}

// Type returns the string representation of the flag's type.
//
// This will be printed in the help output.
func (f *CIDRSliceFlag) Type() string {
	return "[]cidr"
}

// ShortName returns the flag's short name.
//
// Short name is a single case sensitive character (i.e. -N).
func (f *CIDRSliceFlag) ShortName() string {
	return f.short
}

// Usage returns the usage string of the flag.
//
// This will be printed in the help output.
func (f *CIDRSliceFlag) Usage() string {
	return f.usage
}

// IsSet returns true if the value of this flag is set by one of the available sources.
//
// This method returns false if none of the sources has a value to offer, or the value
// has been set to Default (if specified).
func (f *CIDRSliceFlag) IsSet() bool {
	return f.isSet
}

// Var returns a pointer to the underlying variable.
//
// You can also use the Get() method as an alternative.
func (f *CIDRSliceFlag) Var() *[]core.CIDR {
	return f.ptr
}

// Get returns the current value of the flag.
func (f *CIDRSliceFlag) Get() []core.CIDR {
	return f.value
}

// WithKey explicitly defines the key for this flag.
//
// Explicit keys will override the automatically generated values, defined at bucket level (if enabled).
//
// In order for the flag value to be extractable from the environment variables, or all the other custom sources,
// it MUST have a key associated with it.
func (f *CIDRSliceFlag) WithKey(keyID string) *CIDRSliceFlag {
	f.key.SetID(keyID)
	return f
}

// WithDefault sets the default value of the flag.
//
// If none of the available sources offers a value, the default value will be assigned to the flag.
func (f *CIDRSliceFlag) WithDefault(defaultValue []core.CIDR) *CIDRSliceFlag {
	f.defaultValue = defaultValue
	f.hasDefault = true
	return f
}

// Hide marks the flag as hidden.
//
// A hidden flag will not be displayed in the help output.
func (f *CIDRSliceFlag) Hide() *CIDRSliceFlag {
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
func (f *CIDRSliceFlag) MarkAsDeprecated() *CIDRSliceFlag {
	f.isDeprecated = true
	return f
}

// WithDelimiter sets the delimiter for splitting the input string (Default: core.DefaultSliceDelimiter)
func (f *CIDRSliceFlag) WithDelimiter(delimiter string) *CIDRSliceFlag {
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
func (f *CIDRSliceFlag) WithValidationCallback(validate func(in core.CIDR) error) *CIDRSliceFlag {
	f.validate = validate
	return f
}

// WithValidRange defines a list of acceptable values from which the final flag value can be chosen.
//
// The set operation will fail if the flag value is not from the specified list.
// You can also define a custom validation callback function using WithValidationCallback(...) method.
// Remember that setting the valid range will have no effect if a validation callback has been specified.
func (f *CIDRSliceFlag) WithValidRange(valid ...core.CIDR) *CIDRSliceFlag {
	if len(valid) == 0 {
		return f
	}
	f.validationList = make(map[string]interface{})
	f.acceptableItems = make([]string, 0)
	for _, cidr := range valid {
		s := cidr.String()
		if len(s) == 0 {
			continue
		}
		if _, ok := f.validationList[s]; !ok {
			f.acceptableItems = append(f.acceptableItems, s)
			f.validationList[s] = nil
		}
	}
	return f
}

// Set sets the flag value.
//
// The value of a CIDR slice flag can be defined using a list of CIDR notation IP addresses and prefix length,
// like "192.0.2.0/24, 2001:db8::/32", as defined in RFC 4632 and RFC 4291. Each item will be parsed to the
// address and the network implied by the IP and prefix length.
func (f *CIDRSliceFlag) Set(value string) error {
	parts := strings.Split(strings.TrimSpace(value), f.delimiter)
	list := make([]core.CIDR, 0)
	for _, v := range parts {
		value = strings.TrimSpace(v)
		if internal.IsEmpty(v) {
			continue
		}
		cidr, err := core.ParseCIDR(value)
		if err != nil {
			return internal.InvalidValueErr(value, f.long, f.short, f.Type())
		}

		if f.validate != nil {
			err := f.validate(*cidr)
			if err != nil {
				return err
			}
		}

		// Validation callback takes priority over validation list
		if f.validate == nil && len(f.validationList) > 0 {
			if _, ok := f.validationList[cidr.String()]; !ok {
				return internal.OutOfRangeErr(value, f.long, f.short, f.acceptableItems)
			}
		}

		list = append(list, *cidr)
	}

	f.set(list)
	f.isSet = true
	return nil
}

// ResetToDefault resets the value of this flag to default if a default value is specified.
//
// Calling this method on a flag without a default value will have no effect.
// The default value can be defined using WithDefault(...) method.
func (f *CIDRSliceFlag) ResetToDefault() {
	if !f.hasDefault {
		return
	}
	f.isSet = false
	f.set(f.defaultValue)
}

// Default returns the default value if specified, otherwise returns nil
//
// The default value can be defined using WithDefault(...) method
func (f *CIDRSliceFlag) Default() interface{} {
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
func (f *CIDRSliceFlag) Key() *data.Key {
	return f.key
}

func (f *CIDRSliceFlag) set(value []core.CIDR) {
	f.value = value
	*f.ptr = value
}

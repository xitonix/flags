package mocks

import (
	"github.com/xitonix/flags/core"
)

// Flag represents a mocked flag object
type Flag struct {
	long, short   string
	value         interface{}
	isSet         bool
	isDeprecated  bool
	isRequired    bool
	isHidden      bool
	hasDefault    bool
	defaultValue  string
	MakeSetToFail bool
	key           *core.Key
	usage         string
}

// NewFlag creates a new flag mock.
func NewFlag(long, short string) *Flag {
	return NewFlagWithUsage(long, short, "this is a mocked flag")
}

// NewFlagWithUsage creates a new flag mock along with usage string.
func NewFlagWithUsage(long, short, usage string) *Flag {
	return &Flag{
		long:  long,
		short: short,
		key:   &core.Key{},
		usage: usage,
	}
}

// WithKey sets the key
func (f *Flag) WithKey(keyID string) *Flag {
	f.key.SetID(keyID)
	return f
}

// LongName returns the flag's long name
func (f *Flag) LongName() string {
	return f.long
}

// ShortName returns the flag's short name
func (f *Flag) ShortName() string {
	return f.short
}

// Usage returns the flag's usage string
func (f *Flag) Usage() string {
	return f.usage
}

// IsSet returns true if the flag value has been set by one of the available sources.
func (f *Flag) IsSet() bool {
	return f.isSet
}

// IsHidden returns true if the flag is marked as hidden
func (f *Flag) IsHidden() bool {
	return f.isHidden
}

// IsDeprecated returns true if the flag is marked as deprecated
func (f *Flag) IsDeprecated() bool {
	return f.isDeprecated
}

// IsRequired returns true if the flag value must be provided.
func (f *Flag) IsRequired() bool {
	return f.isRequired
}

// Required makes the flag mandatory.
//
// Setting the default value of a required flag will have no effect.
func (f *Flag) Required() *Flag {
	f.isRequired = true
	return f
}

// Type returns "generic" as the flag type.
func (f *Flag) Type() string {
	return "generic"
}

// Key returns the key
func (f *Flag) Key() *core.Key {
	return f.key
}

// Set sets the flag value
func (f *Flag) Set(value string) error {
	if f.MakeSetToFail {
		return ErrExpected
	}
	f.isSet = true
	f.value = value
	return nil
}

// ResetToDefault resets the flag value to default.
func (f *Flag) ResetToDefault() {
	f.value = f.defaultValue
}

// SetDefaultValue sets the default value of the flag.
func (f *Flag) SetDefaultValue(defaultValue string) {
	f.defaultValue = defaultValue
	f.hasDefault = true
}

// Default returns the default value.
func (f *Flag) Default() interface{} {
	if !f.hasDefault {
		return nil
	}
	return f.defaultValue
}

// MarkAsDeprecated marks the flag as deprecated.
func (f *Flag) MarkAsDeprecated() *Flag {
	f.isDeprecated = true
	return f
}

// SetHidden marks the flag as hidden.
func (f *Flag) SetHidden(hidden bool) {
	f.isHidden = hidden
}

// Get returns the flag value.
func (f *Flag) Get() interface{} {
	return f.value
}

package core

import (
	"go.xitonix.io/flags/internal"
)

// Key represents a unique flag ID
//
// Keys are all uppercase and consist of an optional prefix concatenated with an ID.
// The combination of the prefix and the ID must be unique within a bucket.
//
// Different segments of a key are concatenated using underscore characters.
// For example PREFIX_PORT_NUMBER
type Key struct {
	prefix, id string
	isSet      bool
}

// Prefix returns the prefix of the key if it is set.
func (k *Key) Prefix() string {
	return k.prefix
}

// Value returns [PREFIX]_ID
func (k *Key) Value() string {
	if k == nil {
		return ""
	}
	if internal.IsEmpty(k.id) {
		return ""
	}
	if internal.IsEmpty(k.prefix) {
		return k.id
	}
	return k.prefix + "_" + k.id
}

// SetPrefix sets the prefix of the flag keys.
func (k *Key) SetPrefix(prefix string) {
	k.prefix = internal.SanitiseFlagID(prefix)
}

// SetID sets the key ID.
//
// An automatically generated ID ('auto' == true), will be overridden with
// explicit values (where the method is called with 'auto' parameter set to false).
func (k *Key) SetID(id string, auto bool) {
	if k.isSet && auto {
		return
	}
	k.id = internal.SanitiseFlagID(id)
	k.isSet = !auto
}

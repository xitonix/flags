package data

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
}

// Prefix returns the prefix of the key if it is set.
func (k *Key) Prefix() string {
	return k.prefix
}

// Get returns [PREFIX]_ID
func (k *Key) Get() string {
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

func (k *Key) Set(id string) {
	k.id = internal.SanitiseFlagID(id)
}

func (k *Key) IsSet() bool {
	return !internal.IsEmpty(k.id)
}

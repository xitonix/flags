package core

import (
	"github.com/xitonix/flags/internal"
)

// Key represents a unique flag identifier.
//
// Keys are used to extract the flag value from the environment variable source and all the other registered
// custom sources within a bucket.
//
// Each key consists of two parts, an optional prefix and an ID. The combination of these two
// sections must be unique within a bucket.
type Key struct {
	prefix, id string
}

// Prefix returns the prefix of the key if it has been set.
func (k *Key) Prefix() string {
	return k.prefix
}

// FullString returns the string representation of the key.
//
// The return value is uppercase and consists of the optional prefix, concatenated with the ID in PREFIX_ID format.
// If the ID has not been set, this method will return an empty string, even if the prefix is not empty.
func (k *Key) String() string {
	if internal.IsEmpty(k.id) {
		return ""
	}
	if internal.IsEmpty(k.prefix) {
		return k.id
	}
	return k.prefix + "_" + k.id
}

// SetPrefix sets the prefix of the key.
func (k *Key) SetPrefix(prefix string) {
	k.prefix = internal.SanitiseFlagID(prefix)
}

// SetID sets the ID of the key.
func (k *Key) SetID(id string) {
	k.id = internal.SanitiseFlagID(id)
}

// IsSet returns true if the key ID is not empty.
func (k *Key) IsSet() bool {
	return !internal.IsEmpty(k.id)
}

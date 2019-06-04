package core

// Source is an interface to define data providers for the flags.
//
// See `Bucket` for more details.
type Source interface {
	Read(key string) (string, bool)
}

package flags

import (
	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
	"go.xitonix.io/flags/internal"
)

var (
	DefaultBucket = NewBucket()
)

func EnableAutoKeyGeneration() {
	DefaultBucket.opts.AutoKeys = true
}

func SetKeyPrefix(prefix string) {
	DefaultBucket.opts.KeyPrefix = internal.SanitiseFlagID(prefix)
}

func SetLogger(logger core.Logger) {
	DefaultBucket.opts.Logger = logger
}

func Parse() {
	DefaultBucket.Parse()
}

func Options() *config.Options {
	return DefaultBucket.opts
}

func AppendSource(src core.Source) {
	DefaultBucket.AppendSource(src)
}

func PrependSource(src core.Source) {
	DefaultBucket.PrependSource(src)
}

func AddSource(src core.Source, index int) {
	DefaultBucket.AddSource(src, index)
}

func String(longName, usage string) *StringFlag {
	return DefaultBucket.String(longName, usage)
}

func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

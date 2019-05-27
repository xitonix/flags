package flags

import (
	"go.xitonix.io/flags/config"
	"go.xitonix.io/flags/core"
)

var (
	DefaultBucket = NewBucket()
)

func EnableAutoKeyGeneration() {
	DefaultBucket.enableAutoKeyGen()
}

func SetKeyPrefix(prefix string) {
	DefaultBucket.setKeyPrefix(prefix)
}

func SetLogger(logger core.Logger) {
	DefaultBucket.setLogger(logger)
}

func Parse() {
	DefaultBucket.Parse()
}

func Options() *config.Options {
	return DefaultBucket.opts
}

func String(longName, usage string) *StringFlag {
	return DefaultBucket.String(longName, usage)
}

func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

package flags

import (
	"go.xitonix.io/flags/core"
)

var (
	DefaultBucket = NewBucket()
)

func EnableAutoEnv() {
	DefaultBucket.enableAutoEnv()
}

func SetEnvPrefix(prefix string) {
	DefaultBucket.setEnvPrefix(prefix)
}

func SetLogger(logger core.Logger) {
	DefaultBucket.setLogger(logger)
}

func Parse() {
	DefaultBucket.Parse()
}

func String(longName, usage string) *StringFlag {
	return DefaultBucket.String(longName, usage)
}

func StringP(longName, usage, shortName string) *StringFlag {
	return DefaultBucket.StringP(longName, usage, shortName)
}

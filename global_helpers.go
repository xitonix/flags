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

func String(name, usage string) *StringFlag {
	return DefaultBucket.String(name, usage)
}

func StringD(name, usage, defaultValue string) *StringFlag {
	return DefaultBucket.StringD(name, usage, defaultValue)
}

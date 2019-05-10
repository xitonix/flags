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

func String(name string, usage string) *StringFlag {
	return DefaultBucket.String(name, usage)
}

func StringP(name string, short string, usage string) *StringFlag {
	return DefaultBucket.StringP(name, short, usage)
}

func StringPD(name string, short string, defaultValue string, usage string) *StringFlag {
	return DefaultBucket.StringPD(name, short, defaultValue, usage)
}

func StringD(name string, defaultValue string, usage string) *StringFlag {
	return DefaultBucket.StringD(name, defaultValue, usage)
}

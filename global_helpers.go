package flags

import "go.xitonix.io/flags/config"

var (
	DefaultBucket = NewBucket()
)

func EnableAutoEnv() {
	DefaultBucket.enableAutoEnv()
}

func SetEnvPrefix(prefix string) {
	DefaultBucket.setEnvPrefix(prefix)
}

func SetLogger(logger config.Logger) {
	DefaultBucket.setLogger(logger)
}

func Parse() {
	DefaultBucket.Parse()
}

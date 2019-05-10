package config

type Logger interface {
	Fatal(err error)
}

package core

type Source interface {
	Read(key string) (string, bool)
}

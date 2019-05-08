package flags

type Source interface {
	Read(key string) (string, bool)
}

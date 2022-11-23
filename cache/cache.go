package cache

type Cache interface {
	Set(key any, value any)
	Get(key any) (any, bool)
	Del(key any)
	Flush()
}

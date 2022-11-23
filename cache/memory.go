package cache

import `sync`

type memoryCache struct {
	sync.Map
}

func NewDefaultMemoryCache() Cache {
	return new(memoryCache)
}

func (mc *memoryCache) Set(key any, value any) {
	mc.Map.Store(key, value)
}

func (mc *memoryCache) Get(key any) (any, bool) {
	return mc.Map.Load(key)
}

func (mc *memoryCache) Del(key any) {
	mc.Map.Delete(key)
}

func (mc *memoryCache) Flush() {
	mc.Map.Range(func(key, value any) bool {
		mc.Map.Delete(key)

		return true
	})
}

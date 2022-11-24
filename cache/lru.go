package cache

import (
	`container/list`
	`math`
	`sync`
)

const defaultLRUCap = math.MaxInt16

type lruCache struct {
	cap   int
	lst   *list.List
	cache sync.Map
	mutex sync.Mutex
}

type node struct {
	key   any
	value any
}

func NewLRUCache(cap ...int) Cache {
	c := defaultLRUCap
	if len(cap) != 0 {
		c = cap[0]
	}

	if c <= 0 {
		panic("a positive cap is required")
	}

	return &lruCache{
		cap: c,
		lst: list.New(),
	}
}

func (lc *lruCache) Set(key any, value any) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	// key存在，移到最近的节点，更新value值
	if element, ok := lc.cache.Load(key); ok {
		e := element.(*list.Element)

		e.Value.(*node).value = value
		lc.lst.MoveToFront(e)

		lc.cache.Store(key, e)

		return
	}

	// 已达到最大数量，移除最远节点
	if lc.lst.Len() == lc.cap {
		back := lc.lst.Back()
		lc.lst.Remove(back)

		lc.cache.Delete(back.Value.(*node).key)
	}

	// 添加最近的节点
	lc.cache.Store(key, lc.lst.PushFront(&node{key: key, value: value}))

	return
}

func (lc *lruCache) Get(key any) (any, bool) {
	val, ok := lc.cache.Load(key)
	if !ok {
		return nil, false
	}

	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	value := val.(*list.Element)
	lc.lst.MoveToFront(value)

	return value.Value.(*node).value, true
}

func (lc *lruCache) Del(key any) {
	val, ok := lc.cache.Load(key)
	if !ok {
		return
	}

	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.lst.Remove(val.(*list.Element))
	lc.cache.Delete(key)
}

func (lc *lruCache) Flush() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.lst = list.New()
	lc.cache.Range(func(key, value any) bool {
		lc.cache.Delete(key)

		return true
	})
}

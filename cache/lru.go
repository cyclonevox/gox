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
	cache map[any]*list.Element
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
		cap:   c,
		lst:   list.New(),
		cache: make(map[any]*list.Element, c),
	}
}

func (lc *lruCache) Set(key any, value any) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	// key存在，移到最近的节点，更新value值
	if element, ok := lc.cache[key]; ok {
		element.Value.(*node).value = value
		lc.lst.MoveToFront(element)

		lc.cache[key] = element

		return
	}

	// 已达到最大数量，移除最远节点
	if len(lc.cache) == lc.cap {
		back := lc.lst.Back()
		lc.lst.Remove(back)

		delete(lc.cache, back.Value.(*node).key)
	}

	// 添加最近的节点
	lc.cache[key] = lc.lst.PushFront(&node{key: key, value: value})

	return
}

func (lc *lruCache) Get(key any) (any, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	val, ok := lc.cache[key]
	if !ok {
		return nil, false
	}

	lc.lst.MoveToFront(val)

	return val.Value.(*node).value, true
}

func (lc *lruCache) Del(key any) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	val, ok := lc.cache[key]
	if !ok {
		return
	}

	lc.lst.Remove(val)
	delete(lc.cache, key)
}

func (lc *lruCache) Flush() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.lst = list.New()
	lc.cache = make(map[any]*list.Element, lc.cap)
}

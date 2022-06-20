package util

import "sync"

type LruCache[k comparable, v any] struct {
	cap   int
	store *OrderedMap[k, v]
	clear int

	lock sync.RWMutex
}

func NewLruCache[k comparable, v any](cap int, clear int) *LruCache[k, v] {
	if clear == 0 || clear > cap {
		clear = cap
	}
	return &LruCache[k, v]{
		cap:   cap,
		store: NewOrderedMap[k, v](),
		clear: clear,

		lock: sync.RWMutex{},
	}
}

func (c *LruCache[k, v]) HasKey(key k) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.store.HasKey(key)
}

func (c *LruCache[k, v]) Get(key k) v {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.store.Get(key)
}

func (c *LruCache[k, v]) Set(key k, val v) {
	if c.HasKey(key) {
		c.lock.Lock()
		defer c.lock.Unlock()

		elem := c.store.Set(key, val)
		c.store.list.MoveToBack(elem)
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store.list.Len()+1 > c.cap {
		for i := 0; i < c.clear; i++ {
			elem := c.store.list.Front()
			if elem != nil {
				c.store.list.Remove(elem)
			}
		}
	}
	c.store.Set(key, val)
}

package util

import "sync"

type LruCache[K comparable, V any] struct {
	cap   int
	store *OrderedMap[K, V]
	clear int

	lock sync.RWMutex
}

// - cap, the capacity of the store, normally the max number elements can be stored
// - clear, the number of elements to be cleared when the store is full, default is 1/3 of capacity
func NewLruCache[K comparable, V any](cap int, clear int) *LruCache[K, V] {
	if clear == 0 || clear > cap {
		clear = cap / 3
	}
	return &LruCache[K, V]{
		cap:   cap,
		store: NewOrderedMap[K, V](),
		clear: clear,

		lock: sync.RWMutex{},
	}
}

func (c *LruCache[K, V]) HasKey(key K) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.store.HasKey(key)
}

func (c *LruCache[K, V]) Get(key K) V {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.store.Get(key)
}

func (c *LruCache[K, V]) Set(key K, val V) {
	if c.HasKey(key) {
		c.lock.Lock()
		defer c.lock.Unlock()

		elem := c.store.Set(key, val)
		c.store.list.MoveToBack(elem)
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// clear the store if it's full
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

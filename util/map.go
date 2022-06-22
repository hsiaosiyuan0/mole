package util

import "container/list"

type KvPair[K comparable, V any] struct {
	Key K
	Val V
}

type OrderedMap[K comparable, V any] struct {
	dict map[K]*list.Element
	list *list.List
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		dict: map[K]*list.Element{},
		list: list.New(),
	}
}

func (o *OrderedMap[K, V]) HasKey(key K) bool {
	_, ok := o.dict[key]
	return ok
}

func (o *OrderedMap[K, V]) Get(key K) (v V) {
	kv, ok := o.dict[key]
	if ok {
		v = kv.Value.(*KvPair[K, V]).Val
		return
	}
	return
}

func (o *OrderedMap[K, V]) Set(key K, val V) *list.Element {
	pair := &KvPair[K, V]{key, val}

	if old, ok := o.dict[key]; ok {
		elem := o.list.PushBack(pair)
		o.dict[key] = elem
		o.list.Remove(old)
		return elem
	}

	elem := o.list.PushBack(pair)
	o.dict[key] = elem
	return elem
}

func (o *OrderedMap[k, v]) Remove(key k) {
	if old, ok := o.dict[key]; ok {
		delete(o.dict, key)
		o.list.Remove(old)
	}
}

func (o *OrderedMap[k, v]) Front() *list.Element {
	return o.list.Front()
}

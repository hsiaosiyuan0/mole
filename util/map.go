package util

import "container/list"

type KvPair[k comparable, v any] struct {
	Key k
	Val v
}

type OrderedMap[k comparable, v any] struct {
	dict map[k]*list.Element
	list *list.List
}

func NewOrderedMap[k comparable, v any]() *OrderedMap[k, v] {
	return &OrderedMap[k, v]{
		dict: map[k]*list.Element{},
		list: list.New(),
	}
}

func (o *OrderedMap[k, v]) HasKey(key k) bool {
	_, ok := o.dict[key]
	return ok
}

func (o *OrderedMap[k, v]) Set(key k, val v) {
	pair := &KvPair[k, v]{key, val}

	if old, ok := o.dict[key]; ok {
		elem := o.list.InsertBefore(pair, old)
		o.dict[key] = elem
		o.list.Remove(old)
	} else {
		elem := o.list.PushBack(pair)
		o.dict[key] = elem
	}
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

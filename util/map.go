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

func (o *OrderedMap[k, v]) Get(key k) (vv v) {
	kv, ok := o.dict[key]
	if ok {
		vv = kv.Value.(*KvPair[k, v]).Val
		return
	}
	return
}

func (o *OrderedMap[k, v]) Set(key k, val v) *list.Element {
	pair := &KvPair[k, v]{key, val}

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

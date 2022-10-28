package parseutil

import (
	"sync"
)

type ParserType = uint8

const (
	PT_JS ParserType = iota
)

type FutureParsed = chan interface{}
type Parsed = interface{}

type Parser interface {
	Type() ParserType
	Parse(file, code string, readFile bool, biz interface{}) Parsed
}

type ParseCache[T comparable] struct {
	ps map[ParserType]Parser

	store     map[T]interface{}
	storeLock sync.Mutex
}

func NewParseCache[T comparable]() *ParseCache[T] {
	return &ParseCache[T]{
		ps:        map[ParserType]Parser{},
		store:     map[T]interface{}{},
		storeLock: sync.Mutex{},
	}
}

func (c *ParseCache[T]) SetParser(p Parser) {
	c.ps[p.Type()] = p
}

func (c *ParseCache[T]) Parse(typ ParserType, key T, file, code string, readFile bool, biz interface{}) FutureParsed {
	var r interface{}
	c.storeLock.Lock()
	r = c.store[key]
	c.storeLock.Unlock()

	switch v := r.(type) {
	case FutureParsed: // already has a pending process
		return v

	case *Parsed: // process has been done before
		ret := make(FutureParsed, 1)
		ret <- v
		return ret

	case nil: // first time
		ret := make(FutureParsed, 1)
		c.storeLock.Lock()
		c.store[key] = ret
		c.storeLock.Unlock()

		go func() {
			r := c.ps[typ].Parse(file, code, readFile, biz)

			c.storeLock.Lock()
			c.store[key] = r
			c.storeLock.Unlock()

			ret <- r
		}()

		return ret
	}
	panic("unreachable")
}

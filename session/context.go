package session

import (
	"sync"
)

type Context map[interface{}]interface{}

func newContext() Context {
	return make(map[interface{}]interface{})
}

var lock sync.Mutex

func (c Context) Set(key, value interface{}) {
	lock.Lock()
	c[key] = value
	lock.Unlock()
}

func (c Context) Get(key interface{}) interface{} {
	lock.Lock()
	if v, ok := c[key]; ok {
		lock.Unlock()
		return v
	}
	lock.Unlock()
	return nil
}

func (c Context) Delete(key interface{}) {
	lock.Lock()
	delete(c, key)
	lock.Unlock()
}

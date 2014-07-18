package context

import (
	"sync"
)

type Context map[interface{}]interface{}

func NewContext() Context {
	return make(map[interface{}]interface{})
}

func (c Context) Set(key, value interface{}) {
	c[key] = value
}

func (c Context) Get(key interface{}) interface{} {
	if v, ok := c[key]; ok {
		return v
	}
	return nil
}

func (c Context) Delete(key interface{}) {
	delete(c, key)
}

var globalContext Context
var globalLock sync.Mutex

func Set(key, value interface{}) {
	globalLock.Lock()
	globalContext.Set(key, value)
	globalLock.Unlock()
}

func Get(key interface{}) interface{} {
	globalLock.Lock()
	v := globalContext.Get(key)
	globalLock.Unlock()
	return v
}

func Delete(key interface{}) {
	globalLock.Lock()
	globalContext.Delete(key)
	globalLock.Unlock()
}

func init() {
	globalContext = NewContext()
}

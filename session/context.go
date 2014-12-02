package session

//single context
type Context map[interface{}]interface{}

type entireContext struct {
	Input  Context
	Output Context
}

func newContext() *entireContext {
	return &entireContext{
		make(map[interface{}]interface{}),
		make(map[interface{}]interface{}),
	}
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

func (c Context) Purge() {
	for k, _ := range c {
		c.Delete(k)
	}
}

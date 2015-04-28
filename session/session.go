package session

type Session struct {
	Ctx     *entireContext
	Cache   *Cache
	IsLogin bool
}

func NewSession(store Store) *Session {
	return &Session{
		Ctx:     newContext(),
		Cache:   newCache(store),
		IsLogin: false,
	}
}

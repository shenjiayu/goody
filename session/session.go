package session

//admin account
var (
	admin_emails = []string{"sjy19930312@gmail.com"}
)

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

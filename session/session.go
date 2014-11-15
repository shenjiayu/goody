package session

//admin account
var (
	admin_emails = []string{"sjy19930312@gmail.com"}
)

type Session struct {
	Ctx         Context
	Cache       *Cache
	IsLogin     bool
	IsSuperUser bool
}

func NewSession(store Store) *Session {
	return &Session{
		Ctx:         newContext(),
		Cache:       NewCache(store),
		IsLogin:     false,
		IsSuperUser: false,
	}
}

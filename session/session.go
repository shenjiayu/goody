package session

//admin account
var (
	admin_accounts = []string{"shenjiayu"}
)

type Session struct {
	Ctx         Context
	Cache       *Cache
	IsLogin     bool
	IsSuperUser bool
	ValidCsrf   bool
}

func NewSession(store Store) *Session {
	return &Session{
		Ctx:         newContext(),
		Cache:       NewCache(store),
		IsLogin:     false,
		IsSuperUser: false,
		ValidCsrf:   false,
	}
}

func (s *Session) SetCsrf(isValid bool) {
	s.ValidCsrf = isValid
}

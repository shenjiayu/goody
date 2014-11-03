package session

import (
	"fmt"
)

//admin account
var (
	admin_accounts = []string{"shenjiayu"}
)

type Session struct {
	Ctx         Context
	Cache       *Cache
	IsLogin     bool
	IsSuperUser bool
}

func NewSession(store Store, name string) *Session {
	return &Session{
		Ctx:         newContext(),
		Cache:       NewCache(store, name),
		IsLogin:     false,
		IsSuperUser: false,
	}
}

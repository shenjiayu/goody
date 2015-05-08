package session

import (
	"net/http"
)

type Session struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Ctx            *entireContext
	Cache          *Cache
	IsLogin        bool
}

func NewSession(store Store, w http.ResponseWriter, r *http.Request) *Session {
	return &Session{
		ResponseWriter: w,
		Request:        r,
		Ctx:            newContext(),
		Cache:          newCache(store),
		IsLogin:        false,
	}
}

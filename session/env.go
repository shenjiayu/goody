package session

import (
	"net/http"
)

type Env struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
	Status         int
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	env := new(Env)
	env.ResponseWriter = w
	env.Request = r
	store := &RedisStore{}
	env.Session = NewSession(store, "Session_ID")
	return env
}

func (e *Env) SetStatus(status int) {
	e.Status = status
}

func (e *Env) Redirect(url string) {
	http.Redirect(e.ResponseWriter, e.Request, url, http.StatusFound)
}

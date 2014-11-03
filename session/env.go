package session

import (
	"net/http"
)

type Env struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
	Status         int
	finished       bool
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	env := new(Env)
	env.ResponseWriter = w
	env.Request = r
	store := &RedisStore{}
	env.Session = NewSession(store, "Session_ID")
	return env
}

func (e *Env) ProcessRequest() error {
	if session, err := e.Session.Cache.store.New(e.Request, e.ResponseWriter, "Session_ID"); err != nil {
		e.Session = session
		return err
	}
	return nil
}

func (e *Env) ProcessResponse() error {
	//none
	return nil
}

func (e *Env) SetStatus(status int) {
	e.Status = status
}

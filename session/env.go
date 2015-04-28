package session

import (
	"net/http"
)

type Env struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Session        *Session
}

func NewEnv(w http.ResponseWriter, r *http.Request) *Env {
	env := new(Env)
	env.ResponseWriter = w
	env.Request = r
	return env
}

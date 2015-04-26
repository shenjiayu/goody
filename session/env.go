package session

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
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

func Encrypt(data interface{}) string {
	h := sha1.New()
	buf := make([]byte, 5)
	switch data.(type) {
	case int64:
		binary.PutVarint(buf, data.(int64))
	case string:
		buf = []byte(data.(string))
	}
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

type Session struct {
	ID   string
	Ctx  Context
	name string
	//Values  map[interface{}]interface{}
	Username string
	Options  *Options
	store    Store
	IsNew    bool
}

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func NewSession(store Store, name string) *Session {
	return &Session{
		Ctx:  newContext(),
		name: name,
		//Values: make(map[interface{}]interface{}),
		Options: DefaultOptions(),
		store:   store,
	}
}

func DefaultOptions() *Options {
	return &Options{
		Path:     "/",
		Domain:   "coddict.co",
		MaxAge:   720000,
		HttpOnly: true,
	}
}

func (s *Session) NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	if options.MaxAge > 0 {
		d := time.Duration(options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if options.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) Store() Store {
	return s.store
}

func (s *Session) NewID() string {
	return generateID()
}

func generateID() string {
	buf := make([]byte, 40)
	if n, err := rand.Read(buf); err == nil {
		return base64.URLEncoding.EncodeToString(buf[:n])
	}
	return ""
}

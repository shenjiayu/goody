package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Session struct {
	ID      string
	Ctx     Context
	name    string
	Values  Values
	Options *Options
	store   Store
	IsNew   bool
}

func NewSession(store Store, name string) *Session {
	return &Session{
		Ctx:     newContext(),
		name:    name,
		Values:  NewValues(),
		Options: DefaultOptions(),
		store:   store,
	}
}

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func DefaultOptions() *Options {
	return &Options{
		Path:     "/",
		MaxAge:   720000,
		HttpOnly: true,
	}
}

type Values struct {
	Username string `json:"username"`
}

func NewValues() Values {
	return Values{}
}

func (s *Session) NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
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

func (s *Session) EncodingToJson() string {
	data, _ := json.Marshal(s.Values)
	return fmt.Sprintf("%s", data)
}

func (s *Session) DecodingFromJson(data string) string {
	json.Unmarshal([]byte(data), &s.Values)
	return s.Values.Username
}

func generateID() string {
	buf := make([]byte, 40)
	if n, err := rand.Read(buf); err == nil {
		return base64.URLEncoding.EncodeToString(buf[:n])
	}
	return ""
}

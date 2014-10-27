package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
		Ctx:     newContext(),
		Cache:   NewCache(store, name),
		IsLogin: false,
	}
}

type Cache struct {
	ID      string
	name    string
	Values  Values
	Options *Options
	store   Store
}

func NewCache(store Store, name string) *Cache {
	return &Cache{
		name:    name,
		Values:  NewValues(),
		Options: DefaultOptions(),
		store:   store,
	}
}

type Values struct {
	Username string `json:"username"`
	Level    int    `json:"level"`
}

func NewValues() Values {
	return Values{}
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

func (c *Cache) NewCookie(name, value string, options *Options) *http.Cookie {
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

func (c *Cache) Name() string {
	return c.name
}

func (c *Cache) Store() Store {
	return c.store
}

func (c *Cache) NewID() string {
	return generateID()
}

func (c *Cache) EncodingToJson() string {
	data, _ := json.Marshal(c.Values)
	return fmt.Sprintf("%s", data)
}

func (c *Cache) DecodingFromJson(data string) {
	json.Unmarshal([]byte(data), &c.Values)
}

func generateID() string {
	buf := make([]byte, 40)
	if n, err := rand.Read(buf); err == nil {
		return base64.URLEncoding.EncodeToString(buf[:n])
	}
	return ""
}

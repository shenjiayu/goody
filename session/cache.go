package session

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Cache struct {
	ID      string
	Values  *Values
	Options *Options
	store   Store
}

func newCache(store Store) *Cache {
	return &Cache{
		Values:  NewValues(),
		Options: DefaultOptions(),
		store:   store,
	}
}

//init logined users
func NewUser(user_id int, username string, avatar string, status int, email string) *Values {
	return &Values{
		User_id:  user_id,
		Username: username,
		Avatar:   avatar,
		Status:   status,
		Email:    email,
	}
}

//init anonymous users
func AnonymousUser(store Store) *Cache {
	c := new(Cache)
	c.ID = c.NewID()
	c.Values = NewValues()
	c.Values.Csrf = c.NewCsrf()
	c.Values.Level = -1
	c.Options = DefaultOptions()
	c.store = store
	return c
}

type Values struct {
	User_id  int    `json:"user_id"`
	Status   int    `json:"status"`
	Username string `json:"username"`
	Avatar   string `json:avatar`
	Email    string `json:"email"`
	Level    int    `json:"level"` //-1 is anonymous user, 0 is normal user, 1 is admin
	Csrf     string `json:"Csrf"`
}

func NewValues() *Values {
	return &Values{}
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

func (c *Cache) Store() Store {
	return c.store
}

func (c *Cache) NewID() string {
	return generateID()
}

func (c *Cache) NewCsrf() string {
	h := sha1.New()
	buf := make([]byte, 0)
	buf, _ = time.Now().MarshalBinary()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c *Cache) EncodingToJson() string {
	data, _ := json.Marshal(c.Values)
	return fmt.Sprintf("%s", data)
}

func (c *Cache) DecodingFromJson(data string) {
	json.Unmarshal([]byte(data), &c.Values)
}

func generateID() string {
	buf := make([]byte, 20)
	if n, err := rand.Read(buf); err == nil {
		return base64.URLEncoding.EncodeToString(buf[:n])
	}
	return ""
}

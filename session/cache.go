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
func NewUser(user_id int, student_id int, username string, name string, grade int, year int, class int, status int) *Values {
	return &Values{
		User_id:    user_id,
		Student_id: student_id,
		Username:   username,
		Name:       name,
		Grade:      grade,
		Year:       year,
		Class:      class,
		Status:     status,
	}
}

//init anonymous users
func AnonymousUser(store Store) *Cache {
	c := newCache(store)
	c.ID = c.NewID()
	c.Values.Level = -1
	return c
}

type Values struct {
	User_id    int    `json:"user_id"`
	Student_id int    `json:"student_id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Year       int    `json:"year"`
	Class      int    `json:"class"`
	Level      int    `json:"level"` //-1 is anonymous user, 0 is normal user, 1 is admin
	Status     int    `json:"status"`
	Csrf       string `json:"Csrf"`
}

func NewValues() *Values {
	return &Values{Csrf: NewCsrf()}
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
		MaxAge:   2592000,
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

func NewCsrf() string {
	h := sha1.New()
	buf := make([]byte, 0)
	buf, _ = time.Now().MarshalBinary()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

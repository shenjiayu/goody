package session

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

type RedisStore struct {
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

//initiate a pool for incoming connections
var (
	pool = newPool(":6379")
)

func (r *RedisStore) New(req *http.Request, w http.ResponseWriter) (*Session, error) {
	session := NewSession(r)
	if cookie, err := req.Cookie("Session_ID"); err == nil {
		session.Cache.ID = cookie.Value
		if ok, err2 := r.load(session.Cache); err2 == nil && ok {
			if session.Cache.Values.Email != "" {
				session.IsLogin = true
				if session.Cache.Values.Level == 1 {
					session.IsSuperUser = true
				}
			}
		} else {
			session.Cache = AnonymousUser(r)
			if err := r.Save(req, w, session.Cache); err != nil {
				return nil, err
			}
			http.SetCookie(w, session.Cache.NewCookie("Session_ID", session.Cache.ID, session.Cache.Options))
		}
	} else if err == http.ErrNoCookie {
		session.Cache = AnonymousUser(r)
		if err := r.Save(req, w, session.Cache); err != nil {
			return nil, err
		}
		http.SetCookie(w, session.Cache.NewCookie("Session_ID", session.Cache.ID, session.Cache.Options))
	}
	return session, nil
}

func (r *RedisStore) Save(req *http.Request, w http.ResponseWriter, c *Cache) error {
	if c.Options.MaxAge < 0 {
		if err := r.delete(c); err != nil {
			return err
		}
		http.SetCookie(w, c.NewCookie("Session_ID", "", c.Options))
	} else {
		if c.ID == "" {
			c.ID = c.NewID()
		}
		if c.Values.Csrf == "" {
			c.Values.Csrf = c.NewCsrf()
		}
		for _, v := range admin_emails {
			if c.Values.Email == v {
				c.Values.Level = 1
				break
			}
		}
		if err := r.save(c); err != nil {
			fmt.Println(err)
			return err
		}
		http.SetCookie(w, c.NewCookie("Session_ID", c.ID, c.Options))
	}
	return nil
}

func (r *RedisStore) Get(req *http.Request) (*Cache, error) {
	conn := pool.Get()
	defer conn.Close()
	c := &Cache{}
	if cookie, err := req.Cookie("Session_ID"); err == nil {
		_, err := redis.String(conn.Do("GET", "session_"+cookie.Value))
		if err != nil {
			return nil, err
		}
		c.ID = cookie.Value
		c.store = r
	}
	return c, nil
}

func (r *RedisStore) Delete(w http.ResponseWriter, c *Cache) error {
	if c == nil {
		return errors.New("cache cannot be nil")
	}
	if err := r.delete(c); err != nil {
		return err
	}
	c.Options.MaxAge = -1
	http.SetCookie(w, c.NewCookie("Session_ID", "", c.Options))
	return nil
}

func (r *RedisStore) load(c *Cache) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}
	data, err := redis.String(conn.Do("GET", "session_"+c.ID))
	if err != nil {
		return false, err
	}
	//no asociated value for such key
	if data == "" {
		return false, nil
	}
	c.DecodingFromJson(data)
	return true, nil
}

func (r *RedisStore) save(c *Cache) error {
	data := c.EncodingToJson()
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	if _, err := conn.Do("SETEX", "session_"+c.ID, c.Options.MaxAge, data); err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) delete(c *Cache) error {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", "session_"+c.ID); err != nil {
		return err
	}
	return nil
}

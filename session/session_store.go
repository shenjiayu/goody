package session

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

type Store interface {
	//initiate a new session (normally saved it into the cookie of users)
	New(*http.Request, http.ResponseWriter, string) (*Session, error)
	//save the session into backend (redis)
	Save(*http.Request, http.ResponseWriter, *Cache) error
	//get the session out of the backend (redis)
	Get(*http.Request, string) (*Cache, error)
}

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

func (r *RedisStore) New(req *http.Request, w http.ResponseWriter, name string) (*Session, error) {
	session := NewSession(r, name)
	if cookie, err := req.Cookie(name); err == nil {
		session.Cache.ID = cookie.Value
		if ok, username, err2 := r.load(session.Cache); err2 == nil && ok {
			session.Cache.Values.Username = username
			session.IsLogin = true
		} else {
			session.Cache.Options.MaxAge = -1
			http.SetCookie(w, session.Cache.NewCookie(session.Cache.Name(), "", session.Cache.Options))
		}
	}
	return session, nil
}

func (r *RedisStore) Save(req *http.Request, w http.ResponseWriter, c *Cache) error {
	if c.Options.MaxAge < 0 {
		if err := r.delete(c); err != nil {
			return err
		}
		http.SetCookie(w, c.NewCookie(c.Name(), "", c.Options))
	} else {
		if c.ID == "" {
			c.ID = c.NewID()
		}
		if err := r.save(c); err != nil {
			fmt.Println(err)
			return err
		}
		http.SetCookie(w, c.NewCookie(c.Name(), c.ID, c.Options))
	}
	return nil
}

func (r *RedisStore) Get(req *http.Request, name string) (*Cache, error) {
	conn := pool.Get()
	defer conn.Close()
	c := &Cache{}
	if cookie, err := req.Cookie(name); err == nil {
		_, err := redis.String(conn.Do("GET", "session_"+cookie.Value))
		if err != nil {
			return nil, err
		}
		c.ID = cookie.Value
		c.name = name
		c.store = r
	}
	return c, nil
}

func (r *RedisStore) load(c *Cache) (bool, string, error) {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, "", err
	}
	data, err := redis.String(conn.Do("GET", "session_"+c.ID))
	if err != nil {
		return false, "", err
	}
	//no asociated value for such key
	if data == "" {
		return false, "", nil
	}
	username := c.DecodingFromJson(data)
	return true, username, nil
}

func (r *RedisStore) save(c *Cache) error {
	data := c.EncodingToJson()
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", "session_"+c.ID, c.Options.MaxAge, data)
	return err
}

func (r *RedisStore) delete(c *Cache) error {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", "session_"+c.ID); err != nil {
		return err
	}
	return nil
}

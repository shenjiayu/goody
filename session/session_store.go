package session

import (
	"github.com/garyburd/redigo/redis"
	"net/http"
	"time"
)

type Store interface {
	//initiate a new session (normally saved it into the cookie of users)
	New(*http.Request, http.ResponseWriter, string) (*Session, error)
	//save the session into backend (redis)
	Save(*http.Request, http.ResponseWriter, *Session) error
	//get the session out of the backend (redis)
	Get(*http.Request, string) (*Session, error)
}

type RedisStore struct {
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
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
	session.IsNew = true
	if cookie, err := req.Cookie(name); err == nil {
		session.ID = cookie.Value
		if ok, username, err2 := r.load(session); err2 == nil && ok {
			session.Values.Username = username
			session.IsNew = false
		} else {
			session.Options.MaxAge = -1
			http.SetCookie(w, session.NewCookie(session.Name(), "", session.Options))
		}
	}
	return session, nil
}

func (r *RedisStore) Save(req *http.Request, w http.ResponseWriter, s *Session) error {
	if s.Options.MaxAge < 0 {
		if err := r.delete(s); err != nil {
			return err
		}
		http.SetCookie(w, s.NewCookie(s.Name(), "", s.Options))
	} else {
		if s.ID == "" {
			s.ID = s.NewID()
		}
		if err := r.save(s); err != nil {
			return err
		}
		http.SetCookie(w, s.NewCookie(s.Name(), s.ID, s.Options))
	}
	return nil
}

func (r *RedisStore) Get(req *http.Request, name string) (*Session, error) {
	conn := pool.Get()
	defer conn.Close()
	s := &Session{}
	if cookie, err := req.Cookie(name); err == nil {
		_, err := redis.String(conn.Do("GET", "session_"+cookie.Value))
		if err != nil {
			return nil, err
		}
		s.ID = cookie.Value
		s.name = name
		s.store = r
	}
	return s, nil
}

func (r *RedisStore) load(s *Session) (bool, string, error) {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, "", err
	}
	data, err := redis.String(conn.Do("GET", "session_"+s.ID))
	if err != nil {
		return false, "", err
	}
	//no asociated value for such key
	if data == "" {
		return false, "", nil
	}
	username := s.DecodingFromJson(data)
	return true, username, nil
}

func (r *RedisStore) save(s *Session) error {
	data := s.EncodingToJson()
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	_, err := conn.Do("SETEX", "session_"+s.ID, s.Options.MaxAge, data)
	return err
}

func (r *RedisStore) delete(s *Session) error {
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", "session_"+s.ID); err != nil {
		return err
	}
	return nil
}

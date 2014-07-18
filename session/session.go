package session

type Session struct {
	Values map[interface{}]interface{}
	Id     string
	MaxAge int
}

func NewSession(id string, maxAge int) *Session {
	return &Session{Id: id, MaxAge: maxAge, Values: make(map[interface{}]interface{})}
}

func (s *Session) Set(key, value interface{}) {
	s.Values[key] = value
}

func (s *Session) Get(key interface{}) interface{} {
	if v, ok := s.Values[key]; ok {
		return v
	}
	return nil
}

func (s *Session) Delete(key interface{}) {
	delete(s.Values, key)
}

func (s *Session) Expire(seconds int) {
	if seconds < 0 {
		seconds = 0
	}
	s.MaxAge = seconds
}

func (s *Session) Flush() {
	s.Values = make(map[interface{}]interface{})
	s.Id = NewUUID().HexString()
	s.MaxAge = 0
}

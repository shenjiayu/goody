package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type Session struct {
	Ctx      Context
	Username string
	Cookies  map[string]*http.Cookie
}

func NewSession(r *http.Request) *Session {
	s := new(Session)
	s.Ctx = newContext()
	s.Username = ""
	s.Cookies = make(map[string]*http.Cookie)
	if cookie, err := r.Cookie("Session_ID"); err != http.ErrNoCookie {
		s.Cookies["Session_ID"] = cookie
	}
	if cookie, err := r.Cookie("admin"); err != http.ErrNoCookie {
		s.Cookies["admin"] = cookie
	}
	if cookie, err := r.Cookie("token"); err != http.ErrNoCookie {
		s.Cookies["token"] = cookie
	}
	if cookie, err := r.Cookie("access_token"); err != http.ErrNoCookie {
		s.Cookies["access_token"] = cookie
	}
	return s
}

func (s *Session) NewCookie(name string) {
	switch name {
	case "Session_ID":
		s.Cookies["Session_ID"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "token":
		s.Cookies["token"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "admin":
		s.Cookies["admin"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "access_token":
		s.Cookies["access_token"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	}
}

func (s *Session) DestroyCookie(w http.ResponseWriter) {
	for _, v := range s.Cookies {
		v = &http.Cookie{Name: v.Name, MaxAge: -1}
		s.setCookies(w, v)
	}
}

func (s *Session) NotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Session) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (s *Session) setCookies(w http.ResponseWriter, cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	http.SetCookie(w, cookie)
}

func (s *Session) SetCookies(w http.ResponseWriter, cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	s.setCookies(w, cookie)
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

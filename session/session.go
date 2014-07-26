package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

type Session struct {
	Ctx                 Context
	Username            string
	Cookie_Session      *http.Cookie
	Cookie_Admin        *http.Cookie
	Cookie_Token        *http.Cookie
	Cookie_Access_Token *http.Cookie
}

func NewSession(r *http.Request) *Session {
	s := new(Session)
	s.Ctx = newContext()
	s.Username = ""
	if cookie, err := r.Cookie("Session_ID"); err != http.ErrNoCookie {
		s.Cookie_Session = cookie
	}
	if cookie, err := r.Cookie("admin"); err != http.ErrNoCookie {
		s.Cookie_Admin = cookie
	}
	if cookie, err := r.Cookie("token"); err != http.ErrNoCookie {
		s.Cookie_Token = cookie
	}
	if cookie, err := r.Cookie("access_token"); err != http.ErrNoCookie {
		s.Cookie_Access_Token = cookie
	}
	return s
}

func (s *Session) NewCookie(name string) {
	switch name {
	case "Session_ID":
		s.Cookie_Session = &http.Cookie{Name: name, Path: "/", MaxAge: 0, HttpOnly: true}
	case "token":
		s.Cookie_Token = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "admin":
		s.Cookie_Admin = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "access_token":
		s.Cookie_Access_Token = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	}
}

func (s *Session) DestroyCookie(w http.ResponseWriter) {
	cookie := http.Cookie{Name: "Session_ID", MaxAge: -1}
	s.setCookies(w, &cookie)
	cookie = http.Cookie{Name: "token", MaxAge: -1}
	s.setCookies(w, &cookie)
	cookie = http.Cookie{Name: "admin", MaxAge: -1}
	s.setCookies(w, &cookie)
	cookie = http.Cookie{Name: "access_token", MaxAge: -1}
	s.setCookies(w, &cookie)
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

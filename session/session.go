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
	return s
}

func (s *Session) NewCookie(name string) {
	switch name {
	case "Session_ID":
		s.Cookies["Session_ID"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	case "admin":
		s.Cookies["admin"] = &http.Cookie{Name: name, Path: "/", MaxAge: 72000, HttpOnly: true}
	}
}

func (s *Session) DestroyCookie(w http.ResponseWriter) {
	for _, v := range s.Cookies {
		v = &http.Cookie{Name: v.Name, MaxAge: -1}
		s.SetCookies(w, v)
	}
}

func (s *Session) NotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404", http.StatusMovedPermanently)
}

func (s *Session) Redirect(w http.ResponseWriter, r *http.Request, redirect_url string) {
	http.Redirect(w, r, redirect_url, http.StatusMovedPermanently)
}

func (s *Session) SetCookies(w http.ResponseWriter, cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	http.SetCookie(w, cookie)
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

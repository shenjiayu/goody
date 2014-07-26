package session

import (
	"errors"
	"github.com/shenjiayu/coddict/models"
	"net/http"
)

func (s *Session) Get() error {
	db := models.OpenDB()
	defer db.Close()
	err := db.QueryRow("SELECT username FROM session_store WHERE session_id = $1 AND token = $2", s.Cookie_Session.Value, s.Cookie_Token.Value).Scan(&s.Username)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) New(w http.ResponseWriter) error {
	db := models.OpenDB()
	defer db.Close()
	var str string
	db.QueryRow("SELECT username FROM session_store WHERE username = $1", s.Username).Scan(&str)
	if str == "" {
		stmt, err := db.Prepare("INSERT INTO session_store (username, session_id, token, access_token)VALUES($1, $2, $3, $4)")
		if err != nil {
			return errors.New("error on preparing")
		}
		_, err = stmt.Exec(s.Username, s.Cookie_Session.Value, s.Cookie_Token.Value, s.Cookie_Access_Token.Value)
		if err != nil {
			return errors.New("error on inserting")
		}
		s.setCookies(w, s.Cookie_Session)
		s.setCookies(w, s.Cookie_Token)
		s.setCookies(w, s.Cookie_Access_Token)
		return nil
	}
	if err := s.Save(w); err != nil {
		return err
	}
	return nil
}

func (s *Session) Save(w http.ResponseWriter) error {
	db := models.OpenDB()
	defer db.Close()
	if s.Cookie_Session.Value == "" {
		return errors.New("session_id should not be empty")
	}
	if s.Cookie_Session.MaxAge < 0 {
		return errors.New("session expired")
	}
	stmt, err := db.Prepare("UPDATE session_store SET session_id = $1, token = $2 WHERE username = $3")
	if err != nil {
		return errors.New("error on statement")
	}
	_, err = stmt.Exec(s.Cookie_Session.Value, s.Cookie_Token.Value, s.Username)
	if err != nil {
		return errors.New("error on saving")
	}
	s.setCookies(w, s.Cookie_Session)
	s.setCookies(w, s.Cookie_Token)
	s.setCookies(w, s.Cookie_Access_Token)
	return nil
}

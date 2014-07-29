package session

import (
	"errors"
	"github.com/shenjiayu/coddict/models"
	"net/http"
)

func (s *Session) Get() error {
	db := models.OpenDB()
	defer db.Close()
	db.QueryRow("SELECT username FROM session_store WHERE session_id = $1 AND token = $2", s.Cookies["Session_ID"].Value, s.Cookies["token"].Value).Scan(&s.Username)
	if s.Username == "" {
		return errors.New("no record.")
	}
	return nil
}

func (s *Session) New(w http.ResponseWriter) error {
	db := models.OpenDB()
	defer db.Close()
	var str string
	db.QueryRow("SELECT username FROM session_store WHERE username = $1", s.Username).Scan(&str)
	if str == "" {
		if s.Cookies["access_token"] == nil {
			stmt, err := db.Prepare("INSERT INTO session_store (username, session_id, token)VALUES($1, $2, $3)")
			if err != nil {
				return errors.New("error on preparing")
			}
			_, err = stmt.Exec(s.Username, s.Cookies["Session_ID"].Value, s.Cookies["token"].Value)
			if err != nil {
				return errors.New("error on inserting")
			}
		} else {
			stmt, err := db.Prepare("INSERT INTO session_store (username, session_id, token, access_token)VALUES($1, $2, $3, $4)")
			if err != nil {
				return errors.New("error on preparing")
			}
			_, err = stmt.Exec(s.Username, s.Cookies["Session_ID"].Value, s.Cookies["token"].Value, s.Cookies["access_token"].Value)
			if err != nil {
				return errors.New("error on inserting")
			}
			db.QueryRow("SELECT access_token FROM session_store WHERE username = $1", s.Username).Scan(&s.Cookies["access_token"].Value)
		}
		s.setCookies(w, s.Cookies["Session_ID"])
		s.setCookies(w, s.Cookies["token"])
		s.setCookies(w, s.Cookies["access_token"])
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
	if s.Cookies["Session_ID"].Value == "" {
		return errors.New("session_id should not be empty")
	}
	if s.Cookies["Session_ID"].MaxAge < 0 {
		return errors.New("session expired")
	}
	stmt, err := db.Prepare("UPDATE session_store SET session_id = $1, token = $2 WHERE username = $3")
	if err != nil {
		return errors.New("error on statement")
	}
	_, err = stmt.Exec(s.Cookies["Session_ID"].Value, s.Cookies["token"].Value, s.Username)
	if err != nil {
		return errors.New("error on saving")
	}
	db.QueryRow("SELECT access_token FROM session_store WHERE username = $1", s.Username).Scan(&s.Cookies["access_token"].Value)
	s.setCookies(w, s.Cookies["Session_ID"])
	s.setCookies(w, s.Cookies["token"])
	s.setCookies(w, s.Cookies["access_token"])
	return nil
}

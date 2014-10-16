package session

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"net/http"
)

var (
	db, _ = sql.Open("postgres", "user=postgres password=postgres dbname=xinjiapoyan sslmode=disable")
)

func (s *Session) GetUser() (string, error) {
	db.QueryRow("SELECT username FROM session_store WHERE session_id = $1", s.Cookies["Session_ID"].Value).Scan(&s.Username)
	if s.Username == "" {
		return "", errors.New("no record")
	}
	return s.Username, nil
}

func (s *Session) GetToken() (string, error) {
	exist := ""
	db.QueryRow("SELECT token FROM session_store WHERE session_id = $1", s.Cookies["Session_ID"].Value).Scan(&exist)
	if exist == "" {
		return "", errors.New("not valid")
	}
	return exist, nil
}

func (s *Session) New(w http.ResponseWriter) error {
	exist := ""
	db.QueryRow("SELECT username FROM session_store WHERE username = $1", s.Username).Scan(&exist)
	if exist == "" {
		stmt, err := db.Prepare("INSERT INTO session_store (username, session_id, token)VALUES($1, $2, $3)")
		if err != nil {
			return errors.New("error on preparing")
		}
		_, err = stmt.Exec(s.Username, s.Cookies["Session_ID"].Value, "")
		if err != nil {
			return errors.New("error on inserting")
		}
		s.SetCookies(w, s.Cookies["Session_ID"])
		return nil
	} else if err := s.Save(w); err != nil {
		return err
	}
	return nil
}

func (s *Session) Save(w http.ResponseWriter) error {
	if s.Cookies["Session_ID"].Value == "" {
		return errors.New("session_id should not be empty")
	}
	if s.Cookies["Session_ID"].MaxAge < 0 {
		return errors.New("session expired")
	}
	stmt, err := db.Prepare("UPDATE session_store SET session_id = $1 WHERE username = $2")
	defer stmt.Close()
	if err != nil {
		return errors.New("error on statement")
	}
	_, err = stmt.Exec(s.Cookies["Session_ID"].Value, s.Username)
	if err != nil {
		return errors.New("error on saving")
	}
	s.SetCookies(w, s.Cookies["Session_ID"])
	return nil
}

func (s *Session) RefreshToken(token string) (string, error) {
	stmt, err := db.Prepare("UPDATE session_store SET token = $1 WHERE session_id = $2")
	defer stmt.Close()
	if err != nil {
		return "", err
	}
	_, err = stmt.Exec(token, s.Cookies["Session_ID"].Value)
	if err != nil {
		return "", err
	}
	return token, nil
}

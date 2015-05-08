package middleware

import (
	"fmt"
	"github.com/shenjiayu/goody/session"
	"net/http"
)

func ProcessRequest(req *http.Request, w http.ResponseWriter) (*session.Session, error) {
	store := session.RedisStore{}
	s, err := store.New(req, w)
	if err != nil {
		return nil, err
	}
	if s.Request.Method != "GET" {
		s.Request.ParseForm()
		csrf := s.Request.FormValue("csrf")
		if csrf != s.Cache.Values.Csrf {
			return nil, fmt.Errorf("error:csrf is invalid")
		}
		s.Ctx.Input.Set("form", s.Request.Form)
	} else {
		s.Ctx.Output.Set("Csrf", s.Cache.Values.Csrf)
	}
	if s.IsLogin {
		s.Ctx.Output.Set("IsLogin", true)
		s.Ctx.Output.Set("Cache", s.Cache.Values)
	}
	return s, nil
}

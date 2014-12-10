package middleware

import (
	"fmt"
	"github.com/shenjiayu/goody/session"
	"github.com/shenjiayu/goody/template"
)

func ProcessRequest(env *session.Env) error {
	store := session.RedisStore{}
	if s, err := store.New(env.Request, env.ResponseWriter); err != nil {
		return err
	} else {
		env.Session = s
		if env.Request.Method != "GET" {
			env.Request.ParseForm()
			token := env.Request.FormValue("csrf")
			if token != env.Session.Cache.Values.Csrf {
				return fmt.Errorf("error:csrf")
			} else {
				env.Session.Ctx.Input.Set("form", env.Request.Form)
			}
		} else {
			env.Session.Ctx.Output.Set("Csrf", env.Session.Cache.Values.Csrf)
		}
		if env.Session.IsLogin {
			env.Session.Ctx.Output.Set("IsLogin", true)
			env.Session.Ctx.Output.Set("Cache", env.Session.Cache.Values)
		}
	}
	return nil
}

func ProcessResponse(env *session.Env) error {
	switch env.Output_method {
	case "render":
		template.RenderTemplate(env.ResponseWriter, env.Output_data.(string), env.Session.Ctx.Output)
	case "json":
		env.ServeJson(env.ResponseWriter, env.Output_data)
	case "eventstream":
		env.ServeEventStream(env.ResponseWriter, env.Output_data)
	default:
		return fmt.Errorf("Only supports ['render', 'json', 'eventstream'] methods for responsing")
	}
	return nil
}

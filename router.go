package goody

import (
	"fmt"
	"github.com/shenjiayu/goody/session"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type router struct {
	literalLocs map[string]*location
	regexpLocs  []*location
}

func newRouter() *router {
	return &router{
		literalLocs: make(map[string]*location),
		regexpLocs:  make([]*location, 0),
	}
}

func (router *router) registerLocation(pattern string, l *location) error {
	meta := regexp.QuoteMeta(pattern)
	//if meta is same as pattern which is not a regular expression
	if meta == pattern {
		if _, ok := router.literalLocs[pattern]; ok {
			return fmt.Errorf("literal %s location has been registered", pattern)
		}
		router.literalLocs[pattern] = l
	} else {
		if strings.HasPrefix(pattern, "^") {
			pattern = "^" + pattern
		}
		if strings.HasSuffix(pattern, "$") {
			pattern = pattern + "$"
		}
		for _, l := range router.regexpLocs {
			if l.pattern == pattern {
				return fmt.Errorf("regexp %s location has been registered", pattern)
			}
		}
		var err error
		if l.regexpPattern, err = regexp.Compile(pattern); err != nil {
			return err
		}
		router.regexpLocs = append(router.regexpLocs, l)
	}
	return nil
}

func (router *router) Handle(pattern string, handler interface{}) error {
	if len(pattern) == 0 {
		return fmt.Errorf("pattern cannot be empty")
	}
	l, err := newLocation(pattern, handler)
	if err != nil {
		return err
	}
	if err := router.registerLocation(pattern, l); err != nil {
		return err
	}
	return nil
}

// part to be refactored
func (router *router) processRequest(env *session.Env) error {
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
				env.Session.Ctx.Set("form", env.Request.Form)
			}
		} else {
			env.Session.Ctx.Set("Csrf", env.Session.Cache.Values.Csrf)
		}
		if env.Session.IsLogin {
			env.Session.Ctx.Set("IsLogin", true)
			env.Session.Ctx.Set("User_id", env.Session.Cache.Values.User_id)
			if env.Session.Cache.Values.Username == "" {
				env.Session.Ctx.Set("Display_info", env.Session.Cache.Values.Email)
			} else {
				env.Session.Ctx.Set("Display_info", env.Session.Cache.Values.Username)
			}
		}
	}
	return nil
}

func (router *router) processResponse(env *session.Env) error {
	switch env.Output_method {
	case "render":
		env.RenderTemplate(env.ResponseWriter, env.Output_data.(string), env.Session.Ctx)
	case "json":
		env.ServeJson(env.Output_data, env.ResponseWriter)
	case "eventstream":
		env.ServeEventStream(env.Output_data, env.ResponseWriter)
	default:
		return fmt.Errorf("Only supports ['render', 'json', 'eventstream'] methods for responsing")
	}
	return nil
}

func (router *router) CallMethod(w http.ResponseWriter, r *http.Request, l *location, args ...string) {
	env := session.NewEnv(w, r)
	if err := router.processRequest(env); err != nil {
		return
	}
	//fmt.Println(env.Session.ValidToken)
	envValue := reflect.ValueOf(env)
	m, _ := l.methods[r.Method]
	//init the arguments
	in := make([]reflect.Value, m.Type().NumIn())
	//the first argument is '*session.Env'.
	in[0] = envValue
	//iterate over the passed arguments 'args' to in variables.
	for i, v := range args {
		in[i+1] = reflect.ValueOf(v)
	}
	//call corresponding method and pass in the 'in' variable.
	m.Call(in)
	if err := router.processResponse(env); err != nil {
		return
	}
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if l, ok := router.literalLocs[path]; ok {
		router.CallMethod(w, r, l)
	} else {
		//iterate over all regular expression locations
		for _, l := range router.regexpLocs {
			//args will be nil, if the regular expression cannot find submatch of this path
			arg := l.regexpPattern.FindStringSubmatch(path)
			if arg != nil {
				router.CallMethod(w, r, l, arg[1])
				return
			}
		}
		http.NotFound(w, r)
		return
	}
}

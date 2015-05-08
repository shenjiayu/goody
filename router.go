package goody

import (
	"fmt"
	"github.com/shenjiayu/goody/middleware"
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

func (router *router) CallMethod(req *http.Request, w http.ResponseWriter, l *location, args ...string) {
	s, err := middleware.ProcessRequest(req, w)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	sessionValue := reflect.ValueOf(s)
	m, _ := l.methods[req.Method]
	if m.Kind() == reflect.Invalid {
		return
	}
	//init the arguments
	in := make([]reflect.Value, m.Type().NumIn())
	//the first argument is '*session.Env'.
	in[0] = sessionValue
	//iterate over the passed arguments 'args' to in variables.
	for i, v := range args {
		in[i+1] = reflect.ValueOf(v)
	}
	//call corresponding method and pass in the 'in' variable.
	m.Call(in)
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if l, ok := router.literalLocs[path]; ok {
		router.CallMethod(r, w, l)
	} else {
		//iterate over all regular expression locations
		for _, l := range router.regexpLocs {
			//args will be nil, if the regular expression cannot find submatch of this path
			arg := l.regexpPattern.FindStringSubmatch(path)
			if arg != nil {
				router.CallMethod(r, w, l, arg[1])
				return
			}
		}
		http.NotFound(w, r)
		return
	}
}

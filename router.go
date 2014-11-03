package goody

import (
	"fmt"
	"net/http"
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

func (router *router) newRegLocation(pattern string, l *location) error {
	meta := regexp.QuoteMeta(pattern)
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
	if err := router.newRegLocation(pattern, l); err != nil {
		return err
	}
	return nil
}

func (router *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if l, ok := router.literalLocs[path]; ok {
		l.invoke(w, r)
	} else {
		for _, l := range router.regexpLocs {
			args := l.regexpPattern.FindStringSubmatch(path)
			if args != nil {
				l.invoke(w, r, args[1:]...)
				return
			}
		}
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

package goody

import (
	"fmt"
	"github.com/shenjiayu/goody/session"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type location struct {
	pattern       string
	regexpPattern *regexp.Regexp
	methods       map[string]reflect.Value
}

func (l *location) invoke(w http.ResponseWriter, r *http.Request, args ...string) {
	env := session.NewEnv(w, r)
	/*
		if err := env.Session.ProcessRequest(env); err != nil {
			return
		}*/
	envValue := reflect.ValueOf(env)
	m, ok := l.methods[r.Method]
	if !ok {
		env.SetStatus(http.StatusMethodNotAllowed)
	}
	in := make([]reflect.Value, m.Type().NumIn())
	in[0] = envValue
	for i, v := range args {
		in[i+1] = reflect.ValueOf(v)
	}
	m.Call(in)
}

var supportMethods = []string{"Get", "Post", "Put", "Head", "Delete"}

//if the arguments that passed in is valid
func (l *location) checkMethod(handler interface{}, m reflect.Type, name string) error {
	nIn := m.NumIn()
	if nIn == 0 || m.In(0).Kind() != reflect.Ptr {
		return fmt.Errorf("%T:function %s first input argument must *session.Env", handler, name)
	}
	if m.In(0).String() != "*session.Env" {
		return fmt.Errorf("%T:function %s first input argument must *session.Env", handler, name)
	}
	if name == "Prepare" && nIn > 1 {
		return fmt.Errorf("%T:function %s must have one input argument", handler, name)
	}
	for i := 1; i < nIn; i++ {
		//right arguments must be string
		if m.In(i).Kind() != reflect.String {
			return fmt.Errorf("%T:function %s %d input arguments must be string", handler, name, i)
		}
	}
	return nil
}

func newLocation(pattern string, handler interface{}) (*location, error) {
	v := reflect.ValueOf(handler)
	l := new(location)
	l.methods = make(map[string]reflect.Value)
	l.pattern = pattern
	hasMethod := false
	for _, method := range supportMethods {
		m := v.MethodByName(method)
		if m.Kind() == reflect.Func {
			if err := l.checkMethod(handler, m.Type(), method); err != nil {
				return nil, err
			}
			hasMethod = true
			l.methods[strings.ToUpper(method)] = m
		}
	}
	if !hasMethod {
		return nil, fmt.Errorf("handler has no any one method in ['GET', 'POST', 'PUT', 'HEAD', 'DELETE']")
	}
	return l, nil
}

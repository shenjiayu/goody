package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"html"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//it is an encapsulation on r.ParseForm() method
//to auto-detect the struct that should be initiated to.
func Form2Struct(r *http.Request, s interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if s == nil {
		return errors.New("struct should not be nil")
	}
	value := reflect.ValueOf(s).Elem()
	for k, v := range r.Form {
		escaped_data := html.EscapeString(v[0])
		//the field of struct is Capitalized to be read from other packages
		k = strings.Title(k)
		//retrieve tag from struct to do validation
		tmp := reflect.TypeOf(s).Elem()
		if err := processTag(tmp, k, &escaped_data); err != nil {
			return err
		}
		field := value.FieldByName(k)
		//validate if this field belongs to this struct
		if field.Kind() != reflect.Invalid {
			//A Kind represents the specific kind of type that a Type represents. The zero Kind is not a valid kind.
			field_kind := field.Kind()
			switch field_kind {
			//here normally, we just had two types of fields, which are 'string', 'int64'.
			case reflect.String:
				field.SetString(escaped_data)
			case reflect.Int:
				tmp, _ := strconv.Atoi(escaped_data)
				field.SetInt(int64(tmp))
			default:
				return errors.New("invalid field")
			}
		}
	}
	return nil
}

//TODO list
//encryption over certain field like 'token', 'password'
func processTag(s reflect.Type, k string, v *string) error {
	if field, ok := s.FieldByName(k); ok {
		if tag := field.Tag.Get("required"); tag == "true" && len(*v) == 0 {
			return errors.New(k + "不能为空")
		}
		if tag := field.Tag.Get("reg"); tag != "" {
			if err := processReg(tag, *v); err != nil {
				return errors.New(k + "不符合标准")
			}
		}
		if tag := field.Tag.Get("encrypt"); tag == "true" {
			*v = encrypt(*v)
		}
		return nil
	}
	return nil
}

//process regular expression
func processReg(pattern, v string) error {
	reg := regexp.MustCompile(pattern)
	if ok := reg.MatchString(v); !ok {
		return errors.New("error")
	}
	return nil
}

//encrpyt certain field
func encrypt(data interface{}) string {
	h := sha1.New()
	buf := make([]byte, 5)
	switch data.(type) {
	case int64:
		binary.PutVarint(buf, data.(int64))
	case string:
		buf = []byte(data.(string))
	}
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//TODO LIST
//This utils is in progress

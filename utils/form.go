package utils

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"html"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	InValidField = errors.New("Invalid field.")
)

//it is an encapsulation on r.ParseForm() method
//to auto-detect the struct that should be parsed to.
func Form2Struct(form url.Values, s interface{}) error {
	if form == nil {
		return errors.New("form should not be nill")
	}
	if s == nil {
		return errors.New("struct should not be nil")
	}
	value := reflect.ValueOf(s).Elem()
	for k, v := range form {
		data := v[0]
		//the field of struct is Capitalized to be read from outside
		k = strings.Title(k)
		//retrieve tag from struct to do validation
		tmp := reflect.TypeOf(s).Elem()
		if err := processTag(tmp, k, &data); err != nil {
			return err
		}
		field := value.FieldByName(k)
		//validate if this field belongs to this struct
		if field.Kind() != reflect.Invalid {
			//A Kind represents the specific kind of type that a Type represents. The zero Kind is not a valid kind.
			field_kind := field.Kind()
			data = html.EscapeString(data)
			switch field_kind {
			//here normally, we just had two types of fields, which are 'string', 'int64'.
			case reflect.String:
				field.SetString(data)
			case reflect.Int:
				tmp, _ := strconv.Atoi(data)
				field.SetInt(int64(tmp))
			case reflect.Int64:
				tmp, _ := strconv.ParseInt(data, 10, 64)
				field.SetInt(tmp)
			case reflect.Float32:
				tmp, _ := strconv.ParseFloat(data, 32)
				field.SetFloat(tmp)
			case reflect.Float64:
				tmp, _ := strconv.ParseFloat(data, 64)
				field.SetFloat(tmp)
			default:
				return errors.New("invalid field")
			}
		} else {
			return InValidField
		}
	}
	return nil
}

//TODO list
//encryption over certain field like 'token', 'password'
func processTag(s reflect.Type, k string, v *string) error {
	if field, ok := s.FieldByName(k); ok {
		if tag := field.Tag.Get("reg"); tag != "" {
			if err := processReg(tag, *v); err != nil {
				return errors.New("format is wrong.")
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
//This utils is under developing

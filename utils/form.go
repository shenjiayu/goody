package utils

import (
	"errors"
	"html"
	"net/http"
	"reflect"
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
	value := reflect.ValueOf(t).Elem()
	for k, v := range r.Form {
		escaped_data := html.EscapeString(v[0])
		//the field of struct is Capitalized to be read from other packages
		k = strings.Title(k)
		field := value.FieldByName(k)
		//validate if this field belongs to this struct
		if field.Kind() != reflect.Invalid {
			//A Kind represents the specific kind of type that a Type represents. The zero Kind is not a valid kind.
			field_kind := field.Kind().String()
			switch field_kind {
			//here normally, we just had two types of fields, which are 'string', 'int64'.
			case "string":
				field.SetString(escaped_data)
			case "int":
				tmp, _ := strconv.Atoi(escaped_data)
				field.SetInt(int64(tmp))
			default:
				return errors.New("invalid field")
			}
		} else {
			return errors.New("invalid field")
		}
	}
	return nil
}

//TODO list
//encryption over certain field like 'token', 'password'

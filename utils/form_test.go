package utils

import (
	"fmt"
	"net/url"
	"testing"
)

type User struct {
	Id       int64  `-`
	Username string `-`
	Password string `encrypt:"true"`
}

func Test_Form2Struct(t *testing.T) {
	form := url.Values{}
	form.Add("id", "1")
	form.Add("username", "shenjiayu")
	form.Add("password", "123456")
	form.Add("email", "xxx@gmail.com")

	user := new(User)

	if err := Form2Struct(form, user); err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Id: %d\nUsername: %s\nPassword: %s\n", user.Id, user.Username, user.Password)
	}
}

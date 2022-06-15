package plugin

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestValidateOpts(t *testing.T) {
	type Address struct {
		City string `json:"city" validate:"required"`
	}

	type User struct {
		Email     string     `json:"email" validate:"required,email"`
		Addresses []*Address `json:"addresses" validate:"required,dive,required"`
	}

	typ := reflect.TypeOf(User{})
	user := reflect.New(typ).Interface().(*User)

	if err := json.Unmarshal([]byte(`
  {
    "email": "deformed email",
    "addresses": [
      {
        "city": "hangzhou"
      },
      {
        "city": "nanjing"
      }
    ]
  }
    `), &user); err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 2, len(user.Addresses), "should be ok")
	util.AssertEqual(t, "hangzhou", user.Addresses[0].City, "should be ok")

	validate := validator.New()
	err := validate.Struct(user)

	util.AssertEqual(t, true, err != nil, "should be ok")
}

func TestOptsParse(t *testing.T) {
	type Opts struct {
		AvoidEscape bool `json:"avoidEscape"`
	}

	optsSchema := &Options{Typ: []reflect.Type{
		reflect.TypeOf(""),
		reflect.TypeOf(Opts{}),
	}}

	_, err := optsSchema.ParseOpts(`["backtick", { "avoidEscape": true }]`, validator.New(), map[int]Validate{
		0: func(in interface{}) error {
			if *in.(*string) != "backtick" {
				return errors.New("not backtick")
			}
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

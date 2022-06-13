package plugin

import (
	"encoding/json"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Options struct {
	// for compatible with the eslint's behavior which encapsulates the options in a list, here
	// we should define the types for each option
	Typ []reflect.Type
}

type Validate = func(in interface{}) error

// only the options whose type is struct can apply the second parameter `validate` as its validator,
// the other options which have primitive type such as string cannot apply any validator by default.
//
// for resolving this problem, the third parameter `validates` is introduced, its key is the index
// of the target option in the options list, for any target option has its index as key in `validates`,
// the value of that key will be used as the option's validator
func (o *Options) ParseOpts(jsonStr string, validate *validator.Validate, validates map[int]Validate) ([]interface{}, error) {
	opts := []interface{}{}
	for _, typ := range o.Typ {
		opt := reflect.New(typ)
		opts = append(opts, opt)
	}

	if err := json.Unmarshal([]byte(jsonStr), &opts); err != nil {
		return nil, err
	}

	for i, opt := range opts {
		if validates != nil {
			if v, ok := validates[i]; ok {
				if err := v(opt); err != nil {
					return nil, err
				}
				continue
			}
		}

		if reflect.ValueOf(opt).Kind() == reflect.Struct {
			if err := validate.Struct(opt); err != nil {
				return nil, err
			}
		}
	}

	return opts, nil
}

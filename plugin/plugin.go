package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
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

	rawOpts := []interface{}{}
	if err := json.Unmarshal([]byte(jsonStr), &rawOpts); err != nil {
		return nil, err
	}

	opts := make([]interface{}, len(o.Typ))
	for i := range opts {
		if rawOpts[i] == nil {
			return nil, errors.New(fmt.Sprintf("missing option for %v\n", o.Typ[i]))
		}

		raw, err := json.Marshal(rawOpts[i])
		if err != nil {
			return nil, err
		}

		opt := reflect.New(o.Typ[i]).Interface()
		if err := json.Unmarshal(raw, &opt); err != nil {
			return nil, err
		}

		if validates != nil {
			if v, ok := validates[i]; ok {
				if err := v(opt); err != nil {
					return nil, err
				}
				continue
			}
		}

		if validate != nil && reflect.ValueOf(opt).Kind() == reflect.Struct {
			if err := validate.Struct(opt); err != nil {
				return nil, err
			}
		}

		opts[i] = opt
	}

	return opts, nil
}

func DefineOptions(opts ...interface{}) *Options {
	o := &Options{
		Typ: []reflect.Type{},
	}
	for _, opt := range opts {
		o.Typ = append(o.Typ, reflect.TypeOf(opt))
	}
	return o
}

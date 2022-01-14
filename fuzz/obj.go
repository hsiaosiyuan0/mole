package fuzz

import (
	"encoding/json"
	"strconv"
	"strings"
)

type WalkObjFn func(path []string, val interface{}) bool

func WalkObj(obj interface{}, path []string, cb WalkObjFn) bool {
	if mv, ok := obj.(map[string]interface{}); ok {
		for k, v := range mv {
			path = append(path, k)
			goon := cb(path, v)
			if !goon {
				return false
			}
			goon = WalkObj(v, path, cb)
			if !goon {
				return false
			}
			path = path[:len(path)-1]
		}
		return true
	}

	if av, ok := obj.([]interface{}); ok {
		for i, v := range av {
			path = append(path, strconv.Itoa(i))
			goon := cb(path, v)
			if !goon {
				return false
			}
			goon = WalkObj(v, path, cb)
			if !goon {
				return false
			}
			path = path[:len(path)-1]
		}
		return true
	}

	return cb(path, obj)
}

func FlattenMap(m map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	cb := func(path []string, val interface{}) bool {
		_, isMap := val.(map[string]interface{})
		_, isArr := val.([]interface{})
		if !isMap && !isArr {
			ret[strings.Join(path, ".")] = val
		}
		return true
	}
	WalkObj(m, make([]string, 0), cb)
	return ret
}

func JsonToKeyPathAndVal(jsonStr string) (map[string]interface{}, error) {
	jsonObj := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &jsonObj)
	if err != nil {
		return nil, err
	}
	return FlattenMap(jsonObj), nil
}

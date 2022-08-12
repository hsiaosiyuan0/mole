package util

import (
	"encoding/json"
	"strconv"
	"strings"
)

type WalkObjFn func(path []string, val interface{}, key interface{}, parent interface{}, arr bool) bool

func WalkObj(obj interface{}, path []string, cb WalkObjFn, key interface{}, parent interface{}, arr bool) bool {
	if mv, ok := obj.(map[string]interface{}); ok {
		for k, v := range mv {
			path = append(path, k)
			goon := cb(path, v, k, mv, false)
			if !goon {
				return false
			}
			goon = WalkObj(v, path, cb, k, mv, false)
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
			goon := cb(path, v, i, av, true)
			if !goon {
				return false
			}
			goon = WalkObj(v, path, cb, i, av, true)
			if !goon {
				return false
			}
			path = path[:len(path)-1]
		}
		return true
	}

	return cb(path, obj, key, parent, arr)
}

func GetByPath(obj map[string]interface{}, path []string) interface{} {
	for len(path) > 0 {
		v := obj[path[0]]
		if v == nil {
			return nil
		}

		path = path[1:]
		if len(path) == 0 {
			return v
		}

		if vo, ok := v.(map[string]interface{}); ok {
			obj = vo
		} else {
			return nil
		}
	}
	return obj
}

func FlattenMap(m map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	cb := func(path []string, val interface{}, key interface{}, parent interface{}, arr bool) bool {
		_, isMap := val.(map[string]interface{})
		_, isArr := val.([]interface{})
		if !isMap && !isArr {
			ret[strings.Join(path, ".")] = val
		}
		return true
	}
	WalkObj(m, make([]string, 0), cb, nil, nil, false)
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

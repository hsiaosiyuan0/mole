package fuzz

import (
	"reflect"
	"testing"
)

func IsNilPtr(v interface{}) bool {
	if v == nil {
		return true
	}
	vv := reflect.ValueOf(v)
	return vv.Kind() == reflect.Ptr && vv.IsNil()
}

func AssertEqual(t *testing.T, except, actual interface{}, msg string) {
	if except == nil && IsNilPtr(actual) {
		return
	}
	if !reflect.DeepEqual(except, actual) {
		t.Fatalf("%s Except: \n%v\nActual: \n%v", msg, except, actual)
	}
}

func AssertEqualJson(t *testing.T, except, actual string) {
	exceptKV, err := JsonToKeyPathAndVal(except)
	if err != nil {
		t.Fatalf("deformed except %v\n", err)
	}
	actualKV, err := JsonToKeyPathAndVal(actual)
	if err != nil {
		t.Fatalf("deformed actual %v\n", err)
	}
	for k, v := range exceptKV {
		exceptV := v
		actualV := actualKV[k]
		if !reflect.DeepEqual(exceptV, actualV) {
			t.Fatalf("\nUnexpected at:\n  %s\nExcept:\n  %v\nActual:\n  %v\n===========================AST===========================\n%v",
				k, exceptV, actualV, actual)
		}
	}
}

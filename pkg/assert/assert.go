package assert

import (
	"reflect"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/utils"
)

func Equal(t *testing.T, except, actual interface{}, msg string) {
	if !reflect.DeepEqual(except, actual) {
		t.Fatalf("%s Except: \n%v\nActual: \n%v", msg, except, actual)
	}
}

func EqualJson(t *testing.T, except, actual string) {
	exceptKV, err := utils.JsonToKeyPathAndVal(except)
	if err != nil {
		t.Fatalf("deformed except %v\n", err)
	}
	actualKV, err := utils.JsonToKeyPathAndVal(actual)
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

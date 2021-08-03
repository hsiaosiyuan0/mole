package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, a, b interface{}, msg string) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%s except: %v actual: %v", msg, a, b)
	}
}

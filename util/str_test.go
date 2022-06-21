package util

import (
	"testing"
)

func TestRemoveJsonComments1(t *testing.T) {
	s, err := RemoveJsonComments(`{
  "a": "" // comment
}`)
	if err != nil {
		t.Fatal(err)
	}
	AssertEqualString(t, "{\n  \"a\": \"\"           \n}", string(s), "should be ok")
}

func TestRemoveJsonComments2(t *testing.T) {
	s, err := RemoveJsonComments(`{
  "a": "// not comment" // comment
}`)
	if err != nil {
		t.Fatal(err)
	}
	AssertEqualString(t, "{\n  \"a\": \"// not comment\"           \n}", string(s), "should be ok")
}

func TestRemoveJsonComments3(t *testing.T) {
	s, err := RemoveJsonComments(`{
  "a": /* comment */ "// not comment"
}`)
	if err != nil {
		t.Fatal(err)
	}
	AssertEqualString(t, "{\n  \"a\":               \"// not comment\"\n}", string(s), "should be ok")
}

func TestRemoveJsonComments4(t *testing.T) {
	s, err := RemoveJsonComments(`{
  "a": /* comment
  */ "// not comment"
}`)
	if err != nil {
		t.Fatal(err)
	}
	AssertEqualString(t, "{\n  \"a\":           \n      \"// not comment\"\n}", string(s), "should be ok")
}

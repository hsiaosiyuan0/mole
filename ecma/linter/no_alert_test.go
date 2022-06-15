package linter

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestNoAlert(t *testing.T) {
	r := lint(t, "alert('no-alert')", &NoAlert{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow the use of `alert`, `confirm`, and `prompt`", r.Diagnoses[0].Msg, "should be ok")
}

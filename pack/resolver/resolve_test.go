package resolver

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestSubpath(t *testing.T) {
	name, subpath := subpathOf("@foo/bar")
	util.AssertEqual(t, "@foo/bar", name, "should be ok")
	util.AssertEqual(t, ".", subpath, "should be ok")

	name, subpath = subpathOf("@foo/bar/baz")
	util.AssertEqual(t, name, "@foo/bar", "should be ok")
	util.AssertEqual(t, "./baz", subpath, "should be ok")
}

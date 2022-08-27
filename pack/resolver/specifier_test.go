package resolver

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestSpecifierImplicitFile(t *testing.T) {
	s, err := NewSpecifier("./foo.mjs?query=1", "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_FILE, s.kind, "should be ok")
	util.AssertEqual(t, "/foo.mjs", s.s, "should be ok")

	s, err = NewSpecifier("../foo.mjs?query=1", "/a")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_FILE, s.kind, "should be ok")
	util.AssertEqual(t, "/foo.mjs", s.s, "should be ok")
}

func TestSpecifierExplicitFile(t *testing.T) {
	s, err := NewSpecifier("file:///opt/nodejs/config.js", "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_FILE, s.kind, "should be ok")
	util.AssertEqual(t, "/opt/nodejs/config.js", s.s, "should be ok")
}

func TestSpecifierDataJs(t *testing.T) {
	s, err := NewSpecifier(`data:text/javascript,console.log("hello!");`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_DATA, s.kind, "should be ok")
	util.AssertEqual(t, SPDK_JS, s.dk, "should be ok")
	util.AssertEqual(t, `console.log("hello!");`, s.d, "should be ok")
}

func TestSpecifierDataJson(t *testing.T) {
	s, err := NewSpecifier(`data:application/json,"world!"`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_DATA, s.kind, "should be ok")
	util.AssertEqual(t, SPDK_JSON, s.dk, "should be ok")
	util.AssertEqual(t, `"world!"`, s.d, "should be ok")
}

func TestSpecifierNode(t *testing.T) {
	s, err := NewSpecifier(`node:fs/promises`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_NODE, s.kind, "should be ok")
	util.AssertEqual(t, `fs/promises`, s.s, "should be ok")
}

func TestSpecifierBare(t *testing.T) {
	s, err := NewSpecifier(`foo`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_BARE, s.kind, "should be ok")
	util.AssertEqual(t, `foo`, s.s, "should be ok")

	s, err = NewSpecifier(`foo/bar`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_BARE, s.kind, "should be ok")
	util.AssertEqual(t, `foo`, s.s, "should be ok")
	util.AssertEqual(t, `./bar`, s.ss, "should be ok")

	s, err = NewSpecifier(`@foo/bar`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_BARE, s.kind, "should be ok")
	util.AssertEqual(t, `@foo/bar`, s.s, "should be ok")
	util.AssertEqual(t, `.`, s.ss, "should be ok")
}

func TestSpecifierImport(t *testing.T) {
	s, err := NewSpecifier(`#foo`, "/")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, SPK_IMPT, s.kind, "should be ok")
	util.AssertEqual(t, `#foo`, s.s, "should be ok")
}

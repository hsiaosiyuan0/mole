package macro

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestMacroLike(t *testing.T) {
	m, ok := HasMacroLike("// #[visitor]")
	util.AssertEqual(t, true, ok, "should be ok")
	util.AssertEqual(t, "visitor", m, "should be ok")
}

func TestMacro(t *testing.T) {
	m, ok := HasMacroLike(`// #[visitor(a.b.c)]`)
	util.AssertEqual(t, true, ok, "should be ok")

	ctxs, err := ParseMacro("", m, nil, nil)
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, "a.b.c", ctxs[0].Args[0], "should be ok")
}

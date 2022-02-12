package macro

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/fuzz"
)

func TestMacroLike(t *testing.T) {
	m, ok := HasMacroLike("// #[visitor]")
	fuzz.AssertEqual(t, true, ok, "should be ok")
	fuzz.AssertEqual(t, "visitor", m, "should be ok")
}

func TestMacro(t *testing.T) {
	m, ok := HasMacroLike(`// #[visitor("true")]`)
	fuzz.AssertEqual(t, true, ok, "should be ok")

	ctxs, err := ParseMacro("", m, nil, nil)
	fuzz.AssertEqual(t, nil, err, "should be ok")

	fuzz.AssertEqual(t, "true", ctxs[0].Args[0], "should be ok")
}

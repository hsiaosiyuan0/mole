package pack

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compile(code string, opts *parser.ParserOpts) (parser.Node, *parser.Parser, error) {
	p := newParser(code, opts)
	prog, err := p.Prog()
	return prog, p, err
}

func TestResolveTopmostStmts(t *testing.T) {
	ast, p, err := compile(`
  function f() {}
  let f1 = f
  var b = c + d
  let e = /*#__PURE__*/ f + g
  double(55);
  /*#__PURE__*/ double(55);

  class A {}

  const hooks = /*#__PURE__*/ {
    fetch: fetch,
    fetchBatch: fetchBatch,
  }
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	stmts := ast.(*parser.Prog).Body()
	tds, _, _ := resolveTopmostStmts(p)
	util.AssertEqual(t, false, tds[stmts[0]].SideEffect, "")
	util.AssertEqual(t, false, tds[stmts[1]].SideEffect, "")
	util.AssertEqual(t, true, tds[stmts[2]].SideEffect, "")
	util.AssertEqual(t, false, tds[stmts[3]].SideEffect, "")
	util.AssertEqual(t, true, tds[stmts[4]].SideEffect, "")
	util.AssertEqual(t, false, tds[stmts[5]].SideEffect, "")
	util.AssertEqual(t, false, tds[stmts[7]].SideEffect, "")
}

func TestResolveTopmostStmtsRef(t *testing.T) {
	ast, p, err := compile(`
  function f() {}
  let f1 = f
  function f2() {
    return () => {
      f()
      f3()
    }
  }
  export {
    f2,
    f1
  }
  import { f3 } from "./test.js"
  export { f3 } from "./test.js"
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	stmts := ast.(*parser.Prog).Body()
	tds, _, _ := resolveTopmostStmts(p)

	f := tds[stmts[0]]
	f1 := tds[stmts[1]]
	f2 := tds[stmts[2]]
	exp := tds[stmts[3]]
	exp1 := tds[stmts[4]]
	util.AssertEqual(t, true, f.Owners[f1] != nil, "")
	util.AssertEqual(t, true, f.Owners[f2] != nil, "")
	util.AssertEqual(t, true, f1.Owned[f] != nil, "")
	util.AssertEqual(t, true, f2.Owned[f] != nil, "")
	util.AssertEqual(t, true, exp.Owned[f1] != nil, "")
	util.AssertEqual(t, true, exp.Owned[f2] != nil, "")
	util.AssertEqual(t, true, f2.Owned[exp1] != nil, "")
}

func TestResolveTopmostStmtsJsx(t *testing.T) {
	ast, p, err := compile(`
  import { F3 } from "./test.js"
  const el = <F3 />
  export default el;
  const px = {};
  px.el = el;
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	stmts := ast.(*parser.Prog).Body()
	tds, _, _ := resolveTopmostStmts(p)

	f3 := tds[stmts[0]]
	el := tds[stmts[1]]
	dft := tds[stmts[2]]
	px := tds[stmts[4]]
	util.AssertEqual(t, true, el.Owned[f3] != nil, "")
	util.AssertEqual(t, true, dft.Owned[el] != nil, "")
	util.AssertEqual(t, true, px.Owned[el] != nil, "")
}

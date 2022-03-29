package parser

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/util"
)

func TestJSX(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.On(FEAT_JSX_NS)
	ast, err := compile(`
  <div:a attr0 attr1={true}
  attr2 = <b/> {...b}
  >&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {...e}
  </div:a>
  `, opts)
	AssertEqual(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	AssertEqual(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JsxElem)
	open := elem.open.(*JsxOpen)
	AssertEqual(t, false, open.closed, "should not closed")

	AssertEqual(t, 4, len(open.attrs), "should have 3 attrs")
	attr2 := open.attrs[2].(*JsxAttr)
	_ = attr2.val.(*JsxElem)

	AssertEqual(t, "attr2", attr2.nameStr, "should be name attr2")

	attr3 := open.attrs[3]
	AssertEqual(t, N_JSX_ATTR_SPREAD, attr3.Type(), "should be name attr3")

	children := elem.children
	AssertEqual(t, 5, len(children), "should have 4 children")

	child1 := children[1].(*JsxExprSpan).expr.(*JsxElem)
	child1_0 := child1.children[0]
	AssertEqual(t, N_JSX_EXPR_SPAN, child1_0.Type(), "should be jsx expr box")

	child3 := children[3]
	AssertEqual(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXFragment(t *testing.T) {
	ast, err := compile(`
  <>&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {...e}
  </>
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	AssertEqual(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JsxElem)
	open := elem.open.(*JsxOpen)
	AssertEqual(t, false, open.closed, "should not closed")

	children := elem.children
	AssertEqual(t, 5, len(children), "should have 4 children")

	child3 := children[3]
	AssertEqual(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXEmpty(t *testing.T) {
	ast, err := compile(`
  <>&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {/* empty */}
  {...e}
  </>
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	AssertEqual(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JsxElem)
	open := elem.open.(*JsxOpen)
	AssertEqual(t, false, open.closed, "should not closed")

	children := elem.children
	AssertEqual(t, 7, len(children), "should have 7 children")

	empty := children[3]
	AssertEqual(t, N_JSX_EXPR_SPAN, empty.Type(), "should be span")
	AssertEqual(t, N_JSX_EMPTY, empty.(*JsxExprSpan).expr.Type(), "should be empty")

	child3 := children[5]
	AssertEqual(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXNoChild(t *testing.T) {
	ast, err := compile(`
  <a></a>
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	AssertEqual(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JsxElem)
	children := elem.children
	AssertEqual(t, 0, len(children), "should have 4 children")
}

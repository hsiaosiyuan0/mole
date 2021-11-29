package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestJSX(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.On(FEAT_JSX_NS)
	ast, err := compile(`
  <div:a attr0 attr1={true}
  attr2 = <b/> attr3 = {...b}
  >&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {...e}
  </div:a>
  `, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	assert.Equal(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JSXElem)
	open := elem.open.(*JSXOpen)
	assert.Equal(t, false, open.closed, "should not closed")

	assert.Equal(t, 4, len(open.attrs), "should have 3 attrs")
	attr2 := open.attrs[2].(*JSXAttr)
	_ = attr2.val.(*JSXElem)

	assert.Equal(t, "attr2", attr2.nameStr, "should be name attr2")

	attr3 := open.attrs[3]
	assert.Equal(t, N_JSX_ATTR_SPREAD, attr3.Type(), "should be name attr3")

	children := elem.children
	assert.Equal(t, 4, len(children), "should have 4 children")

	child1 := children[1].(*JSXExprSpan).expr.(*JSXElem)
	child1_0 := child1.children[0]
	assert.Equal(t, N_JSX_EXPR_SPAN, child1_0.Type(), "should be jsx expr box")

	child3 := children[3]
	assert.Equal(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXFragment(t *testing.T) {
	ast, err := compile(`
  <>&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {...e}
  </>
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	assert.Equal(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JSXElem)
	open := elem.open.(*JSXOpen)
	assert.Equal(t, false, open.closed, "should not closed")

	children := elem.children
	assert.Equal(t, 4, len(children), "should have 4 children")

	child3 := children[3]
	assert.Equal(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXEmpty(t *testing.T) {
	ast, err := compile(`
  <>&CounterClockwiseContourIntegral;{<i>{a}</i>}t2
  {/* empty */}
  {...e}
  </>
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	assert.Equal(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JSXElem)
	open := elem.open.(*JSXOpen)
	assert.Equal(t, false, open.closed, "should not closed")

	children := elem.children
	assert.Equal(t, 5, len(children), "should have 4 children")

	empty := children[3]
	assert.Equal(t, N_JSX_EXPR_SPAN, empty.Type(), "should be span")
	assert.Equal(t, N_JSX_EMPTY, empty.(*JSXExprSpan).expr.Type(), "should be empty")

	child3 := children[4]
	assert.Equal(t, N_JSX_CHILD_SPREAD, child3.Type(), "should be spread")
}

func TestJSXNoChild(t *testing.T) {
	ast, err := compile(`
  <a></a>
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	assert.Equal(t, 1, len(prog.stmts), "should have 1 stmt")

	elem := prog.stmts[0].(*ExprStmt).expr.(*JSXElem)
	children := elem.children
	assert.Equal(t, 0, len(children), "should have 4 children")

}

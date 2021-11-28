package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestJSX(t *testing.T) {
	ast, err := compile(`
  <div:a attr1 attr2={true}
  attr3 = <b/>
  >&CounterClockwiseContourIntegral;{<b>{a}</b>}t2</div:a>
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	elem := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*JSXElem)
	open := elem.open.(*JSXOpen)
	assert.Equal(t, false, open.closed, "should not closed")

	assert.Equal(t, 3, len(open.attrs), "should have 3 attrs")
	attr2 := open.attrs[2].(*JSXAttr)

	assert.Equal(t, "attr3", attr2.nameStr, "should be name attr3")
	_ = attr2.val.(*JSXElem)
}

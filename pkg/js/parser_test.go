package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestProg(t *testing.T) {
	s := NewSource("", "a + b * c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	lhs := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "a", lhs.(*Ident).val.Text(), "should be name a")

	rhs := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr).rhs
	assert.Equal(t, N_EXPR_BIN, rhs.Type(), "should be bin *")

	lhs = rhs.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "b", lhs.(*Ident).val.Text(), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	assert.Equal(t, N_NAME, rhs.Type(), "should be name")
	assert.Equal(t, "c", rhs.(*Ident).val.Text(), "should be name c")
}

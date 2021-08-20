package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestExpr(t *testing.T) {
	s := NewSource("", "a + b - c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	assert.Equal(t, "-", expr.op.Kind().Name, "should be -")

	ab := expr.lhs.(*BinExpr)

	a := ab.lhs
	assert.Equal(t, "a", a.(*Ident).val.Text(), "should be name a")
	assert.Equal(t, "+", ab.op.Kind().Name, "should be +")
	b := ab.rhs
	assert.Equal(t, "b", b.(*Ident).val.Text(), "should be name b")

	c := expr.rhs
	assert.Equal(t, "c", c.(*Ident).val.Text(), "should be name c")
}

func TestExprPcdHigherRight(t *testing.T) {
	s := NewSource("", "a + b * c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "a", lhs.(*Ident).val.Text(), "should be name a")

	rhs := expr.rhs
	assert.Equal(t, N_EXPR_BIN, rhs.Type(), "should be bin *")

	lhs = rhs.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "b", lhs.(*Ident).val.Text(), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	assert.Equal(t, N_NAME, rhs.Type(), "should be name")
	assert.Equal(t, "c", rhs.(*Ident).val.Text(), "should be name c")
}

func TestExprPcdHigherLeft(t *testing.T) {
	s := NewSource("", "a * b + c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	assert.Equal(t, "+", expr.op.Kind().Name, "should be +")

	ab := expr.lhs.(*BinExpr)
	assert.Equal(t, "*", ab.op.Kind().Name, "should be *")
	a := ab.lhs
	assert.Equal(t, "a", a.(*Ident).val.Text(), "should be name a")
	b := ab.rhs
	assert.Equal(t, "b", b.(*Ident).val.Text(), "should be name b")

	c := expr.rhs
	assert.Equal(t, "c", c.(*Ident).val.Text(), "should be name c")
}

func TestExprAssoc(t *testing.T) {
	s := NewSource("", "a ** b ** c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "a", lhs.(*Ident).val.Text(), "should be name a")

	rhs := expr.rhs
	assert.Equal(t, N_EXPR_BIN, rhs.Type(), "should be bin **")

	lhs = rhs.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "b", lhs.(*Ident).val.Text(), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	assert.Equal(t, N_NAME, rhs.Type(), "should be name")
	assert.Equal(t, "c", rhs.(*Ident).val.Text(), "should be name c")
}

func TestCond(t *testing.T) {
	s := NewSource("", "a > 0 ? a : b")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CondExpr)
	test := expr.test.(*BinExpr)
	assert.Equal(t, ">", test.op.Kind().Name, "should be >")
}

func TestAssign(t *testing.T) {
	s := NewSource("", "a = a > 0 ? a : b")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	a := expr.lhs.(*Ident)
	assert.Equal(t, "a", a.val.Text(), "should be a")

	cond := expr.rhs.(*CondExpr)
	test := cond.test.(*BinExpr)
	assert.Equal(t, ">", test.op.Kind().Name, "should be >")
	assert.Equal(t, "a", a.val.Text(), "should be a")
}

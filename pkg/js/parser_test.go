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

func TestMemberExprSubscript(t *testing.T) {
	s := NewSource("", "a[b][c]")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).val.Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).val.Text(), "should be c")
}

func TestMemberExprDot(t *testing.T) {
	s := NewSource("", "a.b.c")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).val.Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).val.Text(), "should be c")
}

func TestUnaryExpr(t *testing.T) {
	s := NewSource("", "a + void 0")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	a := expr.lhs.(*Ident)
	assert.Equal(t, "a", a.val.Text(), "should be a")

	v0 := expr.rhs.(*UnaryExpr)
	assert.Equal(t, "void", v0.op.Text(), "should be void")
	assert.Equal(t, "0", v0.arg.(*NumLit).val.Text(), "should be 0")
}

func TestUpdateExpr(t *testing.T) {
	s := NewSource("", "a + ++b + c++")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	ab := expr.lhs.(*BinExpr)
	assert.Equal(t, "a", ab.lhs.(*Ident).val.Text(), "should be a")

	u1 := ab.rhs.(*UpdateExpr)
	assert.Equal(t, "b", u1.arg.(*Ident).val.Text(), "should be b")
	assert.Equal(t, true, u1.prefix, "should be prefix")

	u2 := expr.rhs.(*UpdateExpr)
	assert.Equal(t, "c", u2.arg.(*Ident).val.Text(), "should be c")
	assert.Equal(t, false, u2.prefix, "should be postfix")
}

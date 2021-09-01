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

func TestNewExpr(t *testing.T) {
	s := NewSource("", "new new a")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*NewExpr).expr.(*NewExpr)
	assert.Equal(t, "a", expr.expr.(*Ident).val.Text(), "should be a")
}

func TestVarDec(t *testing.T) {
	s := NewSource("", "var a = 1")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]
	id := varDec.id.(*Ident)
	init := varDec.init.(*NumLit)
	assert.Equal(t, "a", id.val.Text(), "should be a")
	assert.Equal(t, "1", init.val.Text(), "should be 1")
}

func TestVarDecArrPattern(t *testing.T) {
	s := NewSource("", "var [a, b = 1, [c] = 1, [d = 1]] = e")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.init.(*Ident)
	assert.Equal(t, "e", init.val.Text(), "should be e")

	arr := varDec.id.(*ArrayPattern)
	elem0 := arr.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.val.Text(), "should be a")

	elem1 := arr.elems[1].(*AssignPattern)
	elem1Lhs := elem1.left.(*Ident)
	elem1Rhs := elem1.right.(*NumLit)
	assert.Equal(t, "b", elem1Lhs.val.Text(), "should be b")
	assert.Equal(t, "1", elem1Rhs.val.Text(), "should be 1")

	elem2 := arr.elems[2].(*AssignPattern)
	elem2Lhs := elem2.left.(*ArrayPattern)
	elem2Rhs := elem2.right.(*NumLit)
	assert.Equal(t, "c", elem2Lhs.elems[0].(*Ident).val.Text(), "should be c")
	assert.Equal(t, "1", elem2Rhs.val.Text(), "should be 1")

	elem3 := arr.elems[3].(*ArrayPattern)
	elem31 := elem3.elems[0].(*AssignPattern)
	assert.Equal(t, "d", elem31.left.(*Ident).val.Text(), "should be d")
	assert.Equal(t, "1", elem31.right.(*NumLit).val.Text(), "should be 1")
}

func TestVarDecArrPatternElision(t *testing.T) {
	s := NewSource("", "var [a, , b, , , c, ,] = e")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.init.(*Ident)
	assert.Equal(t, "e", init.val.Text(), "should be e")

	arr := varDec.id.(*ArrayPattern)
	assert.Equal(t, 7, len(arr.elems), "should be len 7")

	elem0 := arr.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.val.Text(), "should be a")

	elem1 := arr.elems[1]
	assert.Equal(t, nil, elem1, "should be nil")

	elem2 := arr.elems[2].(*Ident)
	assert.Equal(t, "b", elem2.val.Text(), "should be b")

	elem3 := arr.elems[3]
	assert.Equal(t, nil, elem3, "should be nil")

	elem4 := arr.elems[4]
	assert.Equal(t, nil, elem4, "should be nil")

	elem5 := arr.elems[5].(*Ident)
	assert.Equal(t, "c", elem5.val.Text(), "should be c")

	elem6 := arr.elems[6]
	assert.Equal(t, nil, elem6, "should be nil")
}

func TestArrLit(t *testing.T) {
	s := NewSource("", "[a, , b, , , c, ,]")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	arrLit := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*ArrLit)
	assert.Equal(t, 7, len(arrLit.elems), "should be len 7")

	elem0 := arrLit.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.val.Text(), "should be a")

	elem1 := arrLit.elems[1]
	assert.Equal(t, nil, elem1, "should be nil")

	elem2 := arrLit.elems[2].(*Ident)
	assert.Equal(t, "b", elem2.val.Text(), "should be b")

	elem3 := arrLit.elems[3]
	assert.Equal(t, nil, elem3, "should be nil")

	elem4 := arrLit.elems[4]
	assert.Equal(t, nil, elem4, "should be nil")

	elem5 := arrLit.elems[5].(*Ident)
	assert.Equal(t, "c", elem5.val.Text(), "should be c")

	elem6 := arrLit.elems[6]
	assert.Equal(t, nil, elem6, "should be nil")
}

func TestObjLit(t *testing.T) {
	s := NewSource("", "var a = {...a, b, ...c, \"d\": 1, [e]: {f: 1}, ...g}")
	p := NewParser(s, make([]string, 0))
	ast, err := p.Prog()
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.id.(*Ident)
	assert.Equal(t, "a", id.val.Text(), "should be a")

	objLit := varDec.init.(*ObjLit)
	assert.Equal(t, 6, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Spread)
	assert.Equal(t, "a", prop0.arg.(*Ident).val.Text(), "should be ...a")

	prop1 := objLit.props[1].(*Ident)
	assert.Equal(t, "b", prop1.val.Text(), "should be b")

	prop2 := objLit.props[2].(*Spread)
	assert.Equal(t, "c", prop2.arg.(*Ident).val.Text(), "should be ...c")

	prop3 := objLit.props[3].(*Prop)
	assert.Equal(t, "d", prop3.key.(*StrLit).val.Text(), "should be d")
	assert.Equal(t, "1", prop3.value.(*NumLit).val.Text(), "should be 1")

	prop4 := objLit.props[4].(*Prop)
	assert.Equal(t, "e", prop4.key.(*Ident).val.Text(), "should be e")
	assert.Equal(t, "f", prop4.value.(*ObjLit).props[0].(*Prop).key.(*Ident).val.Text(), "should be f")
	assert.Equal(t, "1", prop4.value.(*ObjLit).props[0].(*Prop).value.(*NumLit).val.Text(), "should be 1")

	prop5 := objLit.props[5].(*Spread)
	assert.Equal(t, "g", prop5.arg.(*Ident).val.Text(), "should be ...g")
}

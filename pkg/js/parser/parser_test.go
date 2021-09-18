package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func compile(code string) (Node, error) {
	s := NewSource("", code)
	p := NewParser(s, make([]string, 0))
	return p.Prog()
}

func TestExpr(t *testing.T) {
	ast, err := compile("a + b - c")
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
	ast, err := compile("a + b * c")
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
	ast, err := compile("a * b + c")
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
	ast, err := compile("a ** b ** c")
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
	ast, err := compile("a > 0 ? a : b")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CondExpr)
	test := expr.test.(*BinExpr)
	assert.Equal(t, ">", test.op.Kind().Name, "should be >")
}

func TestAssign(t *testing.T) {
	ast, err := compile("a = a > 0 ? a : b")
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
	ast, err := compile("a[b][c]")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).val.Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).val.Text(), "should be c")
}

func TestMemberExprDot(t *testing.T) {
	ast, err := compile("a.b.c")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).val.Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).val.Text(), "should be c")
}

func TestUnaryExpr(t *testing.T) {
	ast, err := compile("a + void 0")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	a := expr.lhs.(*Ident)
	assert.Equal(t, "a", a.val.Text(), "should be a")

	v0 := expr.rhs.(*UnaryExpr)
	assert.Equal(t, "void", v0.op.Text(), "should be void")
	assert.Equal(t, "0", v0.arg.(*NumLit).val.Text(), "should be 0")
}

func TestUpdateExpr(t *testing.T) {
	ast, err := compile("a + ++b + c++")
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
	ast, err := compile("new new a")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*NewExpr).callee.(*NewExpr)
	assert.Equal(t, "a", expr.callee.(*Ident).val.Text(), "should be a")
}

func TestCallExpr(t *testing.T) {
	ast, err := compile("a()(c, ...a, b)")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)
	callee := expr.callee.(*CallExpr)
	assert.Equal(t, "a", callee.callee.(*Ident).val.Text(), "should be a")

	params := expr.args
	assert.Equal(t, "c", params[0].(*Ident).val.Text(), "should be c")
	assert.Equal(t, "a", params[1].(*Spread).arg.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", params[2].(*Ident).val.Text(), "should be b")
}

func TestCallExprMem(t *testing.T) {
	ast, err := compile("a(b).c")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	obj := expr.obj.(*CallExpr)
	callee := obj.callee.(*Ident)
	assert.Equal(t, "a", callee.val.Text(), "should be a")
	assert.Equal(t, "c", expr.prop.(*Ident).val.Text(), "should be c")

	params := obj.args
	assert.Equal(t, "b", params[0].(*Ident).val.Text(), "should be b")
}

func TestCallExprLit(t *testing.T) {
	_, err := compile("a('b')")
	assert.Equal(t, nil, err, "should be prog ok")
}

func TestCallCascadeExpr(t *testing.T) {
	ast, err := compile("a[b][c]()[d][e]()")
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)

	// a[b][c]()[d][e]
	expr0 := expr.callee.(*MemberExpr)

	// a[b][c]()[d]
	expr1 := expr0.obj.(*MemberExpr)
	e := expr0.prop.(*Ident)
	assert.Equal(t, "e", e.val.Text(), "should be e")

	// a[b][c]()
	expr2 := expr1.obj.(*CallExpr)
	d := expr1.prop.(*Ident)
	assert.Equal(t, "d", d.val.Text(), "should be d")

	// a[b][c]
	expr3 := expr2.callee.(*MemberExpr)
	c := expr3.prop.(*Ident)
	assert.Equal(t, "c", c.val.Text(), "should be c")

	// a[b]
	expr4 := expr3.obj.(*MemberExpr)
	b := expr4.prop.(*Ident)
	assert.Equal(t, "b", b.val.Text(), "should be b")
	assert.Equal(t, "a", expr4.obj.(*Ident).val.Text(), "should be a")
}

func TestVarDec(t *testing.T) {
	ast, err := compile("var a = 1")
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]
	id := varDec.id.(*Ident)
	init := varDec.init.(*NumLit)
	assert.Equal(t, "a", id.val.Text(), "should be a")
	assert.Equal(t, "1", init.val.Text(), "should be 1")
}

func TestVarDecArrPattern(t *testing.T) {
	ast, err := compile("var [a, b = 1, [c] = 1, [d = 1]] = e")
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
	ast, err := compile("var [a, , b, , , c, ,] = e")
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
	ast, err := compile("[a, , b, , , c, ,]")
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
	ast, err := compile(`var a = {...a, b, ...c, "d": 1, [e]: {f: 1}, ...g}`)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.id.(*Ident)
	assert.Equal(t, "a", id.val.Text(), "should be a")

	objLit := varDec.init.(*ObjLit)
	assert.Equal(t, 6, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Spread)
	assert.Equal(t, "a", prop0.arg.(*Ident).val.Text(), "should be ...a")

	prop1 := objLit.props[1].(*Prop)
	assert.Equal(t, "b", prop1.key.(*Ident).val.Text(), "should be b")

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

func TestObjLitMethod(t *testing.T) {
	ast, err := compile(`
  var o = {
    a,
    [b] () {},
    c,
    e: () => {},
  }
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.id.(*Ident)
	assert.Equal(t, "o", id.val.Text(), "should be o")

	objLit := varDec.init.(*ObjLit)
	assert.Equal(t, 4, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Prop)
	assert.Equal(t, "a", prop0.key.(*Ident).val.Text(), "should be a")

	prop1 := objLit.props[1].(*Prop)
	assert.Equal(t, "b", prop1.key.(*Ident).val.Text(), "should be b")
	_ = prop1.value.(*FnDec)

	prop2 := objLit.props[2].(*Prop)
	assert.Equal(t, "c", prop2.key.(*Ident).val.Text(), "should be c")

	prop3 := objLit.props[3].(*Prop)
	assert.Equal(t, "e", prop3.key.(*Ident).val.Text(), "should be e")
	_ = prop3.value.(*ArrowFn)
}

func TestFnDec(t *testing.T) {
	ast, err := compile(`
  function a({ b }) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.val.Text(), "should be a")
}

func TestFnExpr(t *testing.T) {
	ast, err := compile(`
  let a = function a({ b }) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*VarDecStmt).decList[0].init.(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.val.Text(), "should be a")
}

func TestAsyncFnDec(t *testing.T) {
	ast, err := compile(`
  async function a({ b }) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.val.Text(), "should be a")
	assert.Equal(t, true, fn.async, "should be true")
}

func TestArrowFn(t *testing.T) {
	ast, err := compile(`
  a = () => {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = expr.rhs.(*ArrowFn)
}

func TestDoWhileStmt(t *testing.T) {
	ast, err := compile(`
  do {} while(1)
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*DoWhileStmt)
}

func TestWhileStmt(t *testing.T) {
	ast, err := compile(`
  while(1) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*WhileStmt)
}

func TestForStmt(t *testing.T) {
	ast, err := compile(`
  for(;;) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForStmt)
}

func TestForInStmt(t *testing.T) {
	ast, err := compile(`
  for (a in b) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)
}

func TestForOfStmt(t *testing.T) {
	ast, err := compile(`
  for (a of b) {}
  for await (a of b) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)

	forAwait := ast.(*Prog).stmts[1].(*ForInOfStmt)
	assert.Equal(t, true, forAwait.await, "should be await")
}

func TestIfStmt(t *testing.T) {
	ast, err := compile(`
  if (a) {} else b
  if (c) {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt := ast.(*Prog).stmts[0].(*IfStmt)
	assert.Equal(t, "a", stmt.test.(*Ident).val.Text(), "should be a")

	stmt = ast.(*Prog).stmts[1].(*IfStmt)
	assert.Equal(t, "c", stmt.test.(*Ident).val.Text(), "should be c")
}

func TestSwitchStmtEmpty(t *testing.T) {
	ast, err := compile(`
	switch (a) {
	}
	`)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*SwitchStmt)
}

func TestSwitchStmt(t *testing.T) {
	ast, err := compile(`
  switch (a) {
    case b in c:
      d
      e
    case f:
    default:
  }
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.test.(*BinExpr)
	assert.Equal(t, "b", test0.lhs.(*Ident).val.Text(), "should be prog b")
	assert.Equal(t, "c", test0.rhs.(*Ident).val.Text(), "should be prog c")

	cons00 := case0.cons[0].(*ExprStmt)
	assert.Equal(t, "d", cons00.expr.(*Ident).val.Text(), "should be prog d")

	cons01 := case0.cons[1].(*ExprStmt)
	assert.Equal(t, "e", cons01.expr.(*Ident).val.Text(), "should be prog e")

	case1 := stmt.cases[1]
	assert.Equal(t, "f", case1.test.(*Ident).val.Text(), "should be prog f")

	case2 := stmt.cases[2]
	assert.Equal(t, nil, case2.test, "should be default")
}

func TestSwitchStmtDefaultMiddle(t *testing.T) {
	ast, err := compile(`
  switch (a) {
    case b in c:
      d
      e
    default:
    case f:
  }
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.test.(*BinExpr)
	assert.Equal(t, "b", test0.lhs.(*Ident).val.Text(), "should be prog b")
	assert.Equal(t, "c", test0.rhs.(*Ident).val.Text(), "should be prog c")

	cons00 := case0.cons[0].(*ExprStmt)
	assert.Equal(t, "d", cons00.expr.(*Ident).val.Text(), "should be prog d")

	cons01 := case0.cons[1].(*ExprStmt)
	assert.Equal(t, "e", cons01.expr.(*Ident).val.Text(), "should be prog e")

	case1 := stmt.cases[1]
	assert.Equal(t, nil, case1.test, "should be default")

	case2 := stmt.cases[2]
	assert.Equal(t, "f", case2.test.(*Ident).val.Text(), "should be prog f")
}

func TestBrkStmt(t *testing.T) {
	ast, err := compile(`
  break
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*BrkStmt)

	ast, err = compile(`
  break a;
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*BrkStmt)
	assert.Equal(t, "a", stmt.label.val.Text(), "should be a")
}

func TestContStmt(t *testing.T) {
	ast, err := compile(`
  continue
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*ContStmt)

	ast, err = compile(`
  continue a;
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*ContStmt)
	assert.Equal(t, "a", stmt.label.val.Text(), "should be a")
}

func TestLabelStmt(t *testing.T) {
	ast, err := compile(`
  a:
  b
  c
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	lbStmt := ast.(*Prog).stmts[0].(*LabelStmt)
	assert.Equal(t, "a", lbStmt.label.val.Text(), "should be a")

	lbBody := lbStmt.body.(*ExprStmt)
	assert.Equal(t, "b", lbBody.expr.(*Ident).val.Text(), "should be b")

	expr := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "c", expr.expr.(*Ident).val.Text(), "should be c")

	ast, err = compile(`
  a: b
  c
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	lbStmt = ast.(*Prog).stmts[0].(*LabelStmt)
	assert.Equal(t, "a", lbStmt.label.val.Text(), "should be a")

	lbBody = lbStmt.body.(*ExprStmt)
	assert.Equal(t, "b", lbBody.expr.(*Ident).val.Text(), "should be b")

	expr = ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "c", expr.expr.(*Ident).val.Text(), "should be c")
}

func TestRetStmt(t *testing.T) {
	ast, err := compile(`
  return a
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	retStmt := ast.(*Prog).stmts[0].(*RetStmt)
	assert.Equal(t, "a", retStmt.arg.(*Ident).val.Text(), "should be a")

	ast, err = compile(`
  return
  a
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*RetStmt)
	assert.Equal(t, nil, stmt0.arg, "should be nil")

	stmt1 := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "a", stmt1.expr.(*Ident).val.Text(), "should be a")
}

func TestThrowStmt(t *testing.T) {
	ast, err := compile(`
  throw a
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*ThrowStmt)
	assert.Equal(t, "a", stmt0.arg.(*Ident).val.Text(), "should be a")

	ast, err = compile(`
  throw
  a
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 = ast.(*Prog).stmts[0].(*ThrowStmt)
	assert.Equal(t, nil, stmt0.arg, "should be nil")

	stmt1 := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "a", stmt1.expr.(*Ident).val.Text(), "should be a")
}

func TestTryStmt(t *testing.T) {
	ast, err := compile(`
  try {} catch(e) {} finally {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*TryStmt)
	catch := stmt0.catch
	assert.Equal(t, "e", catch.param.(*Ident).val.Text(), "should be e")

	assert.Equal(t, true, stmt0.fin != nil, "should have fin")

	ast, err = compile(`
  try {} finally {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 = ast.(*Prog).stmts[0].(*TryStmt)
	assert.Equal(t, true, stmt0.fin != nil, "should have fin")

	_, err = compile(`
  try {}
  `)
	assert.Equal(t, true, err != nil, "should be err")
}

func TestDebugStmt(t *testing.T) {
	ast, err := compile(`
  a
  debugger
  b
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ExprStmt)
	_ = ast.(*Prog).stmts[1].(*DebugStmt)
	_ = ast.(*Prog).stmts[2].(*ExprStmt)
}

func TestEmptyStmt(t *testing.T) {
	ast, err := compile(`
  ;a;;
  ;
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*EmptyStmt)
	_ = ast.(*Prog).stmts[1].(*ExprStmt)
	_ = ast.(*Prog).stmts[2].(*EmptyStmt)
	_ = ast.(*Prog).stmts[3].(*EmptyStmt)
}

func TestClassStmt(t *testing.T) {
	ast, err := compile(`
  class a {}
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*ClassDec)
}

func TestClassField(t *testing.T) {
	ast, err := compile(`
  class a {
    #f1
  }
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.elems[0].(*Field)
	assert.Equal(t, true, elem0.key.(*Ident).pvt, "should be pvt")
	assert.Equal(t, "f1", elem0.key.(*Ident).val.Text(), "should be f1")
}

func TestClassMethod(t *testing.T) {
	ast, err := compile(`
  class a {
    [a] (b) {
      c
    }

    e
    #f () {}
  }
  `)
	assert.Equal(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.elems[0].(*Method)
	assert.Equal(t, "a", elem0.key.(*Ident).val.Text(), "should be a")
	assert.Equal(t, "b", elem0.value.(*FnDec).params[0].(*Ident).val.Text(), "should be b")

	elem1 := cls.body.elems[1].(*Field)
	assert.Equal(t, "e", elem1.key.(*Ident).val.Text(), "should be e")

	elem2 := cls.body.elems[2].(*Method)
	assert.Equal(t, true, elem2.key.(*Ident).pvt, "should be pvt")
	assert.Equal(t, "f", elem2.key.(*Ident).val.Text(), "should be f")
}

func TestSeqExpr(t *testing.T) {
	ast, err := compile(`
  a = (b, c)
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	elem0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	seq := elem0.rhs.(*SeqExpr)
	assert.Equal(t, "b", seq.elems[0].(*Ident).val.Text(), "should be b")
	assert.Equal(t, "c", seq.elems[1].(*Ident).val.Text(), "should be c")
}

func TestClassExpr(t *testing.T) {
	ast, err := compile(`
  a = class {};
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*ClassDec)
}

func TestRegexpExpr(t *testing.T) {
	ast, err := compile(`
  a = /a/
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*RegexpLit)
}

func TestParenExpr(t *testing.T) {
	ast, err := compile(`
  a = (b)
  `)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*Ident)
}

func TestTplExpr(t *testing.T) {
	ast, err := compile("tag`\na${b}c`")
	assert.Equal(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.val.Text(), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	assert.Equal(t, "\na", span0.val.Text(), "should be a")

	span1 := tpl.elems[1].(*Ident)
	assert.Equal(t, "b", span1.val.Text(), "should be b")

	span2 := tpl.elems[2].(*StrLit)
	assert.Equal(t, "c", span2.val.Text(), "should be c")
}

func TestTplExprNest(t *testing.T) {
	ast, err := compile("tag`\na${ f`g\n${d}e` }c`")
	assert.Equal(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.val.Text(), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	assert.Equal(t, "\na", span0.val.Text(), "should be a")

	span2 := tpl.elems[2].(*StrLit)
	assert.Equal(t, "c", span2.val.Text(), "should be c")

	tpl = tpl.elems[1].(*TplExpr)
	tag = tpl.tag.(*Ident)
	assert.Equal(t, "f", tag.val.Text(), "should be f")

	span0 = tpl.elems[0].(*StrLit)
	assert.Equal(t, "g\n", span0.val.Text(), "should be g")

	span2 = tpl.elems[2].(*StrLit)
	assert.Equal(t, "e", span2.val.Text(), "should be e")

	span1 := tpl.elems[1].(*Ident)
	assert.Equal(t, "d", span1.val.Text(), "should be d")
}

func TestTplExprMember(t *testing.T) {
	ast, err := compile("tag`\na${b}c`[d]")
	assert.Equal(t, nil, err, "should be prog ok")
	member := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	tpl := member.obj.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.val.Text(), "should be tag")
}

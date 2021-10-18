package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func compile(code string, opts *ParserOpts) (Node, error) {
	if opts == nil {
		opts = NewParserOpts()
	}
	s := NewSource("", code)
	p := NewParser(s, opts)
	return p.Prog()
}

func testFail(t *testing.T, code, errMs string, opts *ParserOpts) {
	ast, err := compile(code, opts)
	if err == nil {
		t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
	}
	assert.Equal(t, errMs, err.Error(), "")
}

func TestExpr(t *testing.T) {
	ast, err := compile("a + b - c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	assert.Equal(t, "-", expr.OpText(), "should be -")

	ab := expr.lhs.(*BinExpr)

	a := ab.lhs
	assert.Equal(t, "a", a.(*Ident).Text(), "should be name a")
	assert.Equal(t, "+", ab.OpText(), "should be +")
	b := ab.rhs
	assert.Equal(t, "b", b.(*Ident).Text(), "should be name b")

	c := expr.rhs
	assert.Equal(t, "c", c.(*Ident).Text(), "should be name c")
}

func TestExprPcdHigherRight(t *testing.T) {
	ast, err := compile("a + b * c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "a", lhs.(*Ident).Text(), "should be name a")

	rhs := expr.rhs
	assert.Equal(t, N_EXPR_BIN, rhs.Type(), "should be bin *")

	lhs = rhs.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "b", lhs.(*Ident).Text(), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	assert.Equal(t, N_NAME, rhs.Type(), "should be name")
	assert.Equal(t, "c", rhs.(*Ident).Text(), "should be name c")
}

func TestExprPcdHigherLeft(t *testing.T) {
	ast, err := compile("a * b + c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	assert.Equal(t, "+", expr.OpText(), "should be +")

	ab := expr.lhs.(*BinExpr)
	assert.Equal(t, "*", ab.OpText(), "should be *")
	a := ab.lhs
	assert.Equal(t, "a", a.(*Ident).Text(), "should be name a")
	b := ab.rhs
	assert.Equal(t, "b", b.(*Ident).Text(), "should be name b")

	c := expr.rhs
	assert.Equal(t, "c", c.(*Ident).Text(), "should be name c")
}

func TestExprAssoc(t *testing.T) {
	ast, err := compile("a ** b ** c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "a", lhs.(*Ident).Text(), "should be name a")

	rhs := expr.rhs
	assert.Equal(t, N_EXPR_BIN, rhs.Type(), "should be bin **")

	lhs = rhs.(*BinExpr).lhs
	assert.Equal(t, N_NAME, lhs.Type(), "should be name")
	assert.Equal(t, "b", lhs.(*Ident).Text(), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	assert.Equal(t, N_NAME, rhs.Type(), "should be name")
	assert.Equal(t, "c", rhs.(*Ident).Text(), "should be name c")
}

func TestCond(t *testing.T) {
	ast, err := compile("a > 0 ? a : b", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CondExpr)
	test := expr.test.(*BinExpr)
	assert.Equal(t, ">", test.OpText(), "should be >")
}

func TestAssign(t *testing.T) {
	ast, err := compile("a = a > 0 ? a : b", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	a := expr.lhs.(*Ident)
	assert.Equal(t, "a", a.Text(), "should be a")

	cond := expr.rhs.(*CondExpr)
	test := cond.test.(*BinExpr)
	assert.Equal(t, ">", test.OpText(), "should be >")
	assert.Equal(t, "a", a.Text(), "should be a")
}

func TestMemberExprSubscript(t *testing.T) {
	ast, err := compile("a[b][c]", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).Text(), "should be c")
}

func TestMemberExprDot(t *testing.T) {
	ast, err := compile("a.b.c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	assert.Equal(t, "a", ab.obj.(*Ident).Text(), "should be a")
	assert.Equal(t, "b", ab.prop.(*Ident).Text(), "should be b")
	assert.Equal(t, "c", expr.prop.(*Ident).Text(), "should be c")
}

func TestUnaryExpr(t *testing.T) {
	ast, err := compile("a + void 0", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	a := expr.lhs.(*Ident)
	assert.Equal(t, "a", a.Text(), "should be a")

	v0 := expr.rhs.(*UnaryExpr)
	assert.Equal(t, "void", v0.OpText(), "should be void")
	assert.Equal(t, "0", v0.arg.(*NumLit).Text(), "should be 0")
}

func TestUpdateExpr(t *testing.T) {
	ast, err := compile("a + ++b + c++", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	ab := expr.lhs.(*BinExpr)
	assert.Equal(t, "a", ab.lhs.(*Ident).Text(), "should be a")

	u1 := ab.rhs.(*UpdateExpr)
	assert.Equal(t, "b", u1.arg.(*Ident).Text(), "should be b")
	assert.Equal(t, true, u1.prefix, "should be prefix")

	u2 := expr.rhs.(*UpdateExpr)
	assert.Equal(t, "c", u2.arg.(*Ident).Text(), "should be c")
	assert.Equal(t, false, u2.prefix, "should be postfix")
}

func TestNewExpr(t *testing.T) {
	ast, err := compile("new new a", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*NewExpr).callee.(*NewExpr)
	assert.Equal(t, "a", expr.callee.(*Ident).Text(), "should be a")
}

func TestCallExpr(t *testing.T) {
	ast, err := compile("a()(c, ...a, b)", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)
	callee := expr.callee.(*CallExpr)
	assert.Equal(t, "a", callee.callee.(*Ident).Text(), "should be a")

	params := expr.args
	assert.Equal(t, "c", params[0].(*Ident).Text(), "should be c")
	assert.Equal(t, "a", params[1].(*Spread).arg.(*Ident).Text(), "should be a")
	assert.Equal(t, "b", params[2].(*Ident).Text(), "should be b")
}

func TestCallExprMem(t *testing.T) {
	ast, err := compile("a(b).c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	obj := expr.obj.(*CallExpr)
	callee := obj.callee.(*Ident)
	assert.Equal(t, "a", callee.Text(), "should be a")
	assert.Equal(t, "c", expr.prop.(*Ident).Text(), "should be c")

	params := obj.args
	assert.Equal(t, "b", params[0].(*Ident).Text(), "should be b")
}

func TestCallExprLit(t *testing.T) {
	_, err := compile("a('b')", nil)
	assert.Equal(t, nil, err, "should be prog ok")
}

func TestCallCascadeExpr(t *testing.T) {
	ast, err := compile("a[b][c]()[d][e]()", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)

	// a[b][c]()[d][e]
	expr0 := expr.callee.(*MemberExpr)

	// a[b][c]()[d]
	expr1 := expr0.obj.(*MemberExpr)
	e := expr0.prop.(*Ident)
	assert.Equal(t, "e", e.Text(), "should be e")

	// a[b][c]()
	expr2 := expr1.obj.(*CallExpr)
	d := expr1.prop.(*Ident)
	assert.Equal(t, "d", d.Text(), "should be d")

	// a[b][c]
	expr3 := expr2.callee.(*MemberExpr)
	c := expr3.prop.(*Ident)
	assert.Equal(t, "c", c.Text(), "should be c")

	// a[b]
	expr4 := expr3.obj.(*MemberExpr)
	b := expr4.prop.(*Ident)
	assert.Equal(t, "b", b.Text(), "should be b")
	assert.Equal(t, "a", expr4.obj.(*Ident).Text(), "should be a")
}

func TestVarDec(t *testing.T) {
	ast, err := compile("var a = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]
	id := varDec.id.(*Ident)
	init := varDec.init.(*NumLit)
	assert.Equal(t, "a", id.Text(), "should be a")
	assert.Equal(t, "1", init.Text(), "should be 1")
}

func TestVarDecArrPattern(t *testing.T) {
	ast, err := compile("var [a, b = 1, [c] = 1, [d = 1]] = e", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.init.(*Ident)
	assert.Equal(t, "e", init.Text(), "should be e")

	arr := varDec.id.(*ArrayPattern)
	elem0 := arr.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.Text(), "should be a")

	elem1 := arr.elems[1].(*AssignPattern)
	elem1Lhs := elem1.left.(*Ident)
	elem1Rhs := elem1.right.(*NumLit)
	assert.Equal(t, "b", elem1Lhs.Text(), "should be b")
	assert.Equal(t, "1", elem1Rhs.Text(), "should be 1")

	elem2 := arr.elems[2].(*AssignPattern)
	elem2Lhs := elem2.left.(*ArrayPattern)
	elem2Rhs := elem2.right.(*NumLit)
	assert.Equal(t, "c", elem2Lhs.elems[0].(*Ident).Text(), "should be c")
	assert.Equal(t, "1", elem2Rhs.Text(), "should be 1")

	elem3 := arr.elems[3].(*ArrayPattern)
	elem31 := elem3.elems[0].(*AssignPattern)
	assert.Equal(t, "d", elem31.left.(*Ident).Text(), "should be d")
	assert.Equal(t, "1", elem31.right.(*NumLit).Text(), "should be 1")
}

func TestVarDecArrPatternElision(t *testing.T) {
	ast, err := compile("var [a, , b, , , c, ,] = e", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.init.(*Ident)
	assert.Equal(t, "e", init.Text(), "should be e")

	arr := varDec.id.(*ArrayPattern)
	assert.Equal(t, 7, len(arr.elems), "should be len 7")

	elem0 := arr.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.Text(), "should be a")

	elem1 := arr.elems[1]
	assert.Equal(t, nil, elem1, "should be nil")

	elem2 := arr.elems[2].(*Ident)
	assert.Equal(t, "b", elem2.Text(), "should be b")

	elem3 := arr.elems[3]
	assert.Equal(t, nil, elem3, "should be nil")

	elem4 := arr.elems[4]
	assert.Equal(t, nil, elem4, "should be nil")

	elem5 := arr.elems[5].(*Ident)
	assert.Equal(t, "c", elem5.Text(), "should be c")

	elem6 := arr.elems[6]
	assert.Equal(t, nil, elem6, "should be nil")
}

func TestArrLit(t *testing.T) {
	ast, err := compile("[a, , b, , , c, ,]", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	arrLit := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*ArrLit)
	assert.Equal(t, 7, len(arrLit.elems), "should be len 7")

	elem0 := arrLit.elems[0].(*Ident)
	assert.Equal(t, "a", elem0.Text(), "should be a")

	elem1 := arrLit.elems[1]
	assert.Equal(t, nil, elem1, "should be nil")

	elem2 := arrLit.elems[2].(*Ident)
	assert.Equal(t, "b", elem2.Text(), "should be b")

	elem3 := arrLit.elems[3]
	assert.Equal(t, nil, elem3, "should be nil")

	elem4 := arrLit.elems[4]
	assert.Equal(t, nil, elem4, "should be nil")

	elem5 := arrLit.elems[5].(*Ident)
	assert.Equal(t, "c", elem5.Text(), "should be c")

	elem6 := arrLit.elems[6]
	assert.Equal(t, nil, elem6, "should be nil")
}

func TestObjLit(t *testing.T) {
	ast, err := compile(`var a = {...a, b, ...c, "d": 1, [e]: {f: 1}, ...g}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.id.(*Ident)
	assert.Equal(t, "a", id.Text(), "should be a")

	objLit := varDec.init.(*ObjLit)
	assert.Equal(t, 6, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Spread)
	assert.Equal(t, "a", prop0.arg.(*Ident).Text(), "should be ...a")

	prop1 := objLit.props[1].(*Prop)
	assert.Equal(t, "b", prop1.key.(*Ident).Text(), "should be b")

	prop2 := objLit.props[2].(*Spread)
	assert.Equal(t, "c", prop2.arg.(*Ident).Text(), "should be ...c")

	prop3 := objLit.props[3].(*Prop)
	assert.Equal(t, "d", prop3.key.(*StrLit).val, "should be d")
	assert.Equal(t, "1", prop3.value.(*NumLit).Text(), "should be 1")

	prop4 := objLit.props[4].(*Prop)
	assert.Equal(t, "e", prop4.key.(*Ident).Text(), "should be e")
	assert.Equal(t, "f", prop4.value.(*ObjLit).props[0].(*Prop).key.(*Ident).Text(), "should be f")
	assert.Equal(t, "1", prop4.value.(*ObjLit).props[0].(*Prop).value.(*NumLit).Text(), "should be 1")

	prop5 := objLit.props[5].(*Spread)
	assert.Equal(t, "g", prop5.arg.(*Ident).Text(), "should be ...g")
}

func TestObjLitMethod(t *testing.T) {
	ast, err := compile(`
  var o = {
    a,
    [b] () {},
    c,
    e: () => {},
  }
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.id.(*Ident)
	assert.Equal(t, "o", id.Text(), "should be o")

	objLit := varDec.init.(*ObjLit)
	assert.Equal(t, 4, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Prop)
	assert.Equal(t, "a", prop0.key.(*Ident).Text(), "should be a")

	prop1 := objLit.props[1].(*Prop)
	assert.Equal(t, "b", prop1.key.(*Ident).Text(), "should be b")
	_ = prop1.value.(*FnDec)

	prop2 := objLit.props[2].(*Prop)
	assert.Equal(t, "c", prop2.key.(*Ident).Text(), "should be c")

	prop3 := objLit.props[3].(*Prop)
	assert.Equal(t, "e", prop3.key.(*Ident).Text(), "should be e")
	_ = prop3.value.(*ArrowFn)
}

func TestFnDec(t *testing.T) {
	ast, err := compile(`
  function a({ b }) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.Text(), "should be a")
}

func TestFnExpr(t *testing.T) {
	ast, err := compile(`
  let a = function a({ b }) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*VarDecStmt).decList[0].init.(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.Text(), "should be a")
}

func TestAsyncFnDec(t *testing.T) {
	ast, err := compile(`
  async function a({ b }) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.Text(), "should be a")
	assert.Equal(t, true, fn.async, "should be true")
}

func TestArrowFn(t *testing.T) {
	ast, err := compile(`
  a = () => {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = expr.rhs.(*ArrowFn)
}

func TestDoWhileStmt(t *testing.T) {
	ast, err := compile(`
  do {} while(1)
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*DoWhileStmt)
}

func TestWhileStmt(t *testing.T) {
	ast, err := compile(`
  while(1) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*WhileStmt)
}

func TestForStmt(t *testing.T) {
	ast, err := compile(`
  for(;;) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForStmt)
}

func TestForInStmt(t *testing.T) {
	ast, err := compile(`
  for (a in b) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)
}

func TestForOfStmt(t *testing.T) {
	ast, err := compile(`
  for (a of b) {}
  for await (a of b) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)

	forAwait := ast.(*Prog).stmts[1].(*ForInOfStmt)
	assert.Equal(t, true, forAwait.await, "should be await")
}

func TestIfStmt(t *testing.T) {
	ast, err := compile(`
  if (a) {} else b
  if (c) {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt := ast.(*Prog).stmts[0].(*IfStmt)
	assert.Equal(t, "a", stmt.test.(*Ident).Text(), "should be a")

	stmt = ast.(*Prog).stmts[1].(*IfStmt)
	assert.Equal(t, "c", stmt.test.(*Ident).Text(), "should be c")
}

func TestSwitchStmtEmpty(t *testing.T) {
	ast, err := compile(`
	switch (a) {
	}
	`, nil)
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
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.test.(*BinExpr)
	assert.Equal(t, "b", test0.lhs.(*Ident).Text(), "should be prog b")
	assert.Equal(t, "c", test0.rhs.(*Ident).Text(), "should be prog c")

	cons00 := case0.cons[0].(*ExprStmt)
	assert.Equal(t, "d", cons00.expr.(*Ident).Text(), "should be prog d")

	cons01 := case0.cons[1].(*ExprStmt)
	assert.Equal(t, "e", cons01.expr.(*Ident).Text(), "should be prog e")

	case1 := stmt.cases[1]
	assert.Equal(t, "f", case1.test.(*Ident).Text(), "should be prog f")

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
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.test.(*BinExpr)
	assert.Equal(t, "b", test0.lhs.(*Ident).Text(), "should be prog b")
	assert.Equal(t, "c", test0.rhs.(*Ident).Text(), "should be prog c")

	cons00 := case0.cons[0].(*ExprStmt)
	assert.Equal(t, "d", cons00.expr.(*Ident).Text(), "should be prog d")

	cons01 := case0.cons[1].(*ExprStmt)
	assert.Equal(t, "e", cons01.expr.(*Ident).Text(), "should be prog e")

	case1 := stmt.cases[1]
	assert.Equal(t, nil, case1.test, "should be default")

	case2 := stmt.cases[2]
	assert.Equal(t, "f", case2.test.(*Ident).Text(), "should be prog f")
}

func TestBrkStmt(t *testing.T) {
	ast, err := compile(`
  break
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*BrkStmt)

	ast, err = compile(`
  break a;
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*BrkStmt)
	assert.Equal(t, "a", stmt.label.(*Ident).Text(), "should be a")
}

func TestContStmt(t *testing.T) {
	ast, err := compile(`
  continue
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*ContStmt)

	ast, err = compile(`
  continue a;
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*ContStmt)
	assert.Equal(t, "a", stmt.label.(*Ident).Text(), "should be a")
}

func TestLabelStmt(t *testing.T) {
	ast, err := compile(`
  a:
  b
  c
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	lbStmt := ast.(*Prog).stmts[0].(*LabelStmt)
	assert.Equal(t, "a", lbStmt.label.(*Ident).Text(), "should be a")

	lbBody := lbStmt.body.(*ExprStmt)
	assert.Equal(t, "b", lbBody.expr.(*Ident).Text(), "should be b")

	expr := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "c", expr.expr.(*Ident).Text(), "should be c")

	ast, err = compile(`
  a: b
  c
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	lbStmt = ast.(*Prog).stmts[0].(*LabelStmt)
	assert.Equal(t, "a", lbStmt.label.(*Ident).Text(), "should be a")

	lbBody = lbStmt.body.(*ExprStmt)
	assert.Equal(t, "b", lbBody.expr.(*Ident).Text(), "should be b")

	expr = ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "c", expr.expr.(*Ident).Text(), "should be c")
}

func TestRetStmt(t *testing.T) {
	ast, err := compile(`
  return a
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	retStmt := ast.(*Prog).stmts[0].(*RetStmt)
	assert.Equal(t, "a", retStmt.arg.(*Ident).Text(), "should be a")

	ast, err = compile(`
  return
  a
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*RetStmt)
	assert.Equal(t, nil, stmt0.arg, "should be nil")

	stmt1 := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "a", stmt1.expr.(*Ident).Text(), "should be a")
}

func TestThrowStmt(t *testing.T) {
	ast, err := compile(`
  throw a
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*ThrowStmt)
	assert.Equal(t, "a", stmt0.arg.(*Ident).Text(), "should be a")

	ast, err = compile(`
  throw
  a
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 = ast.(*Prog).stmts[0].(*ThrowStmt)
	assert.Equal(t, nil, stmt0.arg, "should be nil")

	stmt1 := ast.(*Prog).stmts[1].(*ExprStmt)
	assert.Equal(t, "a", stmt1.expr.(*Ident).Text(), "should be a")
}

func TestTryStmt(t *testing.T) {
	ast, err := compile(`
  try {} catch(e) {} finally {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*TryStmt)
	catch := stmt0.catch
	assert.Equal(t, "e", catch.(*Catch).param.(*Ident).Text(), "should be e")

	assert.Equal(t, true, stmt0.fin != nil, "should have fin")

	ast, err = compile(`
  try {} finally {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 = ast.(*Prog).stmts[0].(*TryStmt)
	assert.Equal(t, true, stmt0.fin != nil, "should have fin")

	_, err = compile(`
  try {}
  `, nil)
	assert.Equal(t, true, err != nil, "should be err")
}

func TestDebugStmt(t *testing.T) {
	ast, err := compile(`
  a
  debugger
  b
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ExprStmt)
	_ = ast.(*Prog).stmts[1].(*DebugStmt)
	_ = ast.(*Prog).stmts[2].(*ExprStmt)
}

func TestEmptyStmt(t *testing.T) {
	ast, err := compile(`
  ;a;;
  ;
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*EmptyStmt)
	_ = ast.(*Prog).stmts[1].(*ExprStmt)
	_ = ast.(*Prog).stmts[2].(*EmptyStmt)
	_ = ast.(*Prog).stmts[3].(*EmptyStmt)
}

func TestClassStmt(t *testing.T) {
	ast, err := compile(`
  class a {}
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*ClassDec)
}

func TestClassField(t *testing.T) {
	ast, err := compile(`
  class a {
    #f1
  }
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.(*ClassBody).elems[0].(*Field)
	assert.Equal(t, true, elem0.key.(*Ident).pvt, "should be pvt")
	assert.Equal(t, "f1", elem0.key.(*Ident).Text(), "should be f1")
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
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, "a", elem0.key.(*Ident).Text(), "should be a")
	assert.Equal(t, "b", elem0.value.(*FnDec).params[0].(*Ident).Text(), "should be b")

	elem1 := cls.body.(*ClassBody).elems[1].(*Field)
	assert.Equal(t, "e", elem1.key.(*Ident).Text(), "should be e")

	elem2 := cls.body.(*ClassBody).elems[2].(*Method)
	assert.Equal(t, true, elem2.key.(*Ident).pvt, "should be pvt")
	assert.Equal(t, "f", elem2.key.(*Ident).Text(), "should be f")
}

func TestSeqExpr(t *testing.T) {
	ast, err := compile(`
  a = (b, c)
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	elem0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	seq := elem0.rhs.(*SeqExpr)
	assert.Equal(t, "b", seq.elems[0].(*Ident).Text(), "should be b")
	assert.Equal(t, "c", seq.elems[1].(*Ident).Text(), "should be c")
}

func TestClassExpr(t *testing.T) {
	ast, err := compile(`
  a = class {};
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*ClassDec)
}

func TestRegexpExpr(t *testing.T) {
	ast, err := compile(`
  a = /a/
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*RegexpLit)
}

func TestParenExpr(t *testing.T) {
	ast, err := compile(`
  a = (b)
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*Ident)
}

func TestTplExpr(t *testing.T) {
	ast, err := compile("tag`\na${b}c`", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.Text(), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	assert.Equal(t, "\na", span0.val, "should be a")

	span1 := tpl.elems[1].(*Ident)
	assert.Equal(t, "b", span1.Text(), "should be b")

	span2 := tpl.elems[2].(*StrLit)
	assert.Equal(t, "c", span2.val, "should be c")
}

func TestTplExprNest(t *testing.T) {
	ast, err := compile("tag`\na${ f`g\n${d}e` }c`", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.Text(), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	assert.Equal(t, "\na", span0.val, "should be a")

	span2 := tpl.elems[2].(*StrLit)
	assert.Equal(t, "c", span2.val, "should be c")

	tpl = tpl.elems[1].(*TplExpr)
	tag = tpl.tag.(*Ident)
	assert.Equal(t, "f", tag.Text(), "should be f")

	span0 = tpl.elems[0].(*StrLit)
	assert.Equal(t, "g\n", span0.val, "should be g")

	span2 = tpl.elems[2].(*StrLit)
	assert.Equal(t, "e", span2.val, "should be e")

	span1 := tpl.elems[1].(*Ident)
	assert.Equal(t, "d", span1.Text(), "should be d")
}

func TestTplExprMember(t *testing.T) {
	ast, err := compile("tag`\na${b}c`[d]", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	member := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	tpl := member.obj.(*TplExpr)
	tag := tpl.tag.(*Ident)
	assert.Equal(t, "tag", tag.Text(), "should be tag")
}

func TestSuper(t *testing.T) {
	ast, err := compile("class a { constructor() { super() } }", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	ctor := ast.(*Prog).stmts[0].(*ClassDec).body.(*ClassBody).elems[0].(*Method).value.(*FnDec)
	expr := ctor.body.(*BlockStmt).body[0].(*ExprStmt).expr
	call := expr.(*CallExpr)
	assert.Equal(t, N_SUPER, call.callee.Type(), "should be tag")
}

func TestImportCall(t *testing.T) {
	ast, err := compile("a = import(b)", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	assign := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	importCall := assign.rhs.(*ImportCall)
	assert.Equal(t, "b", importCall.src.(*Ident).Text(), "should be b")
}

func TestMetaProp(t *testing.T) {
	ast, err := compile("a = import.meta", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	assign := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	metaProp := assign.rhs.(*MetaProp)
	assert.Equal(t, "meta", metaProp.prop.(*Ident).Text(), "should be meta")
}

func TestFail1(t *testing.T) {
	testFail(t, "{", "Unexpected token `EOF` at (1:1)", nil)
}

func TestFail2(t *testing.T) {
	testFail(t, "}", "Unexpected token `}` at (1:0)", nil)
}

func TestFail3(t *testing.T) {
	testFail(t, "3ea", "Invalid number at (1:0)", nil)
}

func TestFail4(t *testing.T) {
	testFail(t, "3in []", "Identifier directly after number at (1:1)", nil)
}

func TestFail5(t *testing.T) {
	testFail(t, "3e", "Invalid number at (1:0)", nil)
}

func TestFail6(t *testing.T) {
	testFail(t, "3e+", "Invalid number at (1:0)", nil)
}

func TestFail7(t *testing.T) {
	testFail(t, "3e-", "Invalid number at (1:0)", nil)
}

func TestFail8(t *testing.T) {
	testFail(t, "3x", "Identifier directly after number at (1:1)", nil)
}

func TestFail9(t *testing.T) {
	testFail(t, "3x0", "Identifier directly after number at (1:1)", nil)
}

func TestFail10(t *testing.T) {
	testFail(t, "0x", "Expected number in radix 16 at (1:2)", nil)
}

func TestFail11(t *testing.T) {
	testFail(t, "'use strict'; 09", "Invalid number at (1:14)", nil)
}

func TestFail12(t *testing.T) {
	testFail(t, "01a", "Identifier directly after number at (1:2)", nil)
}

func TestFail13(t *testing.T) {
	testFail(t, "3in[]", "Identifier directly after number at (1:1)", nil)
}

func TestFail14(t *testing.T) {
	testFail(t, "0x3in[]", "Identifier directly after number at (1:3)", nil)
}

func TestFail15(t *testing.T) {
	testFail(t, "\"Hello\nWorld\"", "Unterminated string constant at (1:0)", nil)
}

func TestFail16(t *testing.T) {
	testFail(t, "x\\", "Expecting Unicode escape sequence \\uXXXX at (1:1)", nil)
}

func TestFail17(t *testing.T) {
	testFail(t, "x\\u005c", "Invalid Unicode escape at (1:1)", nil)
}

func TestFail18(t *testing.T) {
	testFail(t, "/", "Unterminated regular expression at (1:1)", nil)
}

func TestFail19(t *testing.T) {
	testFail(t, "/test", "Unterminated regular expression at (1:1)", nil)
}

func TestFail20(t *testing.T) {
	testFail(t, "var x = /[a-z]/\\ux", "Bad character escape sequence at (1:17)", nil)
}

func TestFail21(t *testing.T) {
	testFail(t, "3 = 4", "Assigning to rvalue at (1:0)", nil)
}

func TestFail22(t *testing.T) {
	testFail(t, "func() = 4", "Assigning to rvalue at (1:0)", nil)
}

func TestFail23(t *testing.T) {
	testFail(t, "(1 + 1) = 10", "Assigning to rvalue at (1:0)", nil)
}

func TestFail24(t *testing.T) {
	testFail(t, "1++", "Assigning to rvalue at (1:0)", nil)
}

func TestFail25(t *testing.T) {
	testFail(t, "1--", "Assigning to rvalue at (1:0)", nil)
}

func TestFail26(t *testing.T) {
	testFail(t, "++1", "Assigning to rvalue at (1:2)", nil)
}

func TestFail27(t *testing.T) {
	testFail(t, "--1", "Assigning to rvalue at (1:2)", nil)
}

func TestFail28(t *testing.T) {
	testFail(t, "for((1 + 1) in list) process(x);", "Assigning to rvalue at (1:4)", nil)
}

func TestFail29(t *testing.T) {
	testFail(t, "[", "Unexpected token `EOF` at (1:1)", nil)
}

func TestFail30(t *testing.T) {
	testFail(t, "[,", "Unexpected token `EOF` at (1:2)", nil)
}

func TestFail31(t *testing.T) {
	testFail(t, "1 + {", "Unexpected token `EOF` at (1:5)", nil)
}

func TestFail32(t *testing.T) {
	testFail(t, "1 + { t:t ", "Unexpected token `EOF` at (1:10)", nil)
}

func TestFail33(t *testing.T) {
	testFail(t, "1 + { t:t,", "Unexpected token `EOF` at (1:10)", nil)
}

func TestFail34(t *testing.T) {
	testFail(t, "var x = /\n/", "Unterminated regular expression at (1:8)", nil)
}

func TestFail35(t *testing.T) {
	testFail(t, "var x = \"\n", "Unterminated string constant at (1:8)", nil)
}

func TestFail36(t *testing.T) {
	testFail(t, "var if = 42", "Unexpected token `if` at (1:4)", nil)
}

func TestFail37(t *testing.T) {
	testFail(t, "i + 2 = 42", "Assigning to rvalue at (1:0)", nil)
}

func TestFail38(t *testing.T) {
	testFail(t, "+i = 42", "Assigning to rvalue at (1:0)", nil)
}

func TestFail39(t *testing.T) {
	testFail(t, "1 + (", "Unexpected token `EOF` at (1:5)", nil)
}

func TestFail40(t *testing.T) {
	testFail(t, "\n\n\n{", "Unexpected token `EOF` at (4:1)", nil)
}

func TestFail41(t *testing.T) {
	testFail(t, "\n/* Some multiline\ncomment */\n)", "Unexpected token `)` at (4:0)", nil)
}

func TestFail42(t *testing.T) {
	testFail(t, "{ set 1 }", "Unexpected token at (1:6)", nil)
}

func TestFail43(t *testing.T) {
	testFail(t, "{ get 2 }", "Unexpected token at (1:6)", nil)
}

func TestFail44(t *testing.T) {
	testFail(t, "({ set: s(if) { } })", "Unexpected token `if` at (1:10)", nil)
}

func TestFail45(t *testing.T) {
	testFail(t, "({ set s(.) { } })", "Unexpected token `.` at (1:9)", nil)
}

func TestFail46(t *testing.T) {
	testFail(t, "({ set: s() { } })", "Unexpected token `{` at (1:12)", nil)
}

func TestFail47(t *testing.T) {
	testFail(t, "({ set: s(a, b) { } })", "Unexpected token `{` at (1:16)", nil)
}

func TestFail48(t *testing.T) {
	testFail(t, "({ get: g(d) { } })", "Unexpected token `{` at (1:13)", nil)
}

func TestFail49(t *testing.T) {
	testFail(t, "'use strict'; ({ __proto__: 1, __proto__: 2 })", "Redefinition of property at (1:31)", nil)
}

func TestFail50(t *testing.T) {
	testFail(t, "function t(...) { }", "Unexpected token at (1:11)", &ParserOpts{Version: ES5})
}

func TestFail51(t *testing.T) {
	testFail(t, "function t(...) { }", "Unexpected token `)` at (1:14)", nil)
}

func TestFail52(t *testing.T) {
	testFail(t, "function t(...rest,) { }",
		"Unexpected trailing comma after rest element at (1:18)", nil)
}

func TestFail53(t *testing.T) {
	testFail(t, "function t(...rest, b) { }",
		"Rest element must be last element at (1:18)", nil)
}

func TestFail54(t *testing.T) {
	testFail(t, "function t(if) { }",
		"Unexpected token `if` at (1:11)", nil)
}

func TestFail56(t *testing.T) {
	testFail(t, "function t(false) { }",
		"Unexpected token `false` at (1:11)", nil)
}

func TestFail57(t *testing.T) {
	testFail(t, "function t(true) { }",
		"Unexpected token `true` at (1:11)", nil)
}

func TestFail58(t *testing.T) {
	testFail(t, "function t(null) { }",
		"Unexpected token `null` at (1:11)", nil)
}

func TestFail59(t *testing.T) {
	testFail(t, "function true() { }",
		"Unexpected token `true` at (1:9)", nil)
}

func TestFail60(t *testing.T) {
	testFail(t, "function false() { }",
		"Unexpected token `false` at (1:9)", nil)
}

func TestFail61(t *testing.T) {
	testFail(t, "function if() { }",
		"Unexpected token `if` at (1:9)", nil)
}

func TestFail62(t *testing.T) {
	testFail(t, "a b;",
		"Unexpected token at (1:2)", nil)
}

func TestFail63(t *testing.T) {
	testFail(t, "if.a;",
		"Unexpected token `.` at (1:2)", nil)
}

func TestFail64(t *testing.T) {
	testFail(t, "a if;",
		"Unexpected token at (1:2)", nil)
}

func TestFail65(t *testing.T) {
	testFail(t, "a class;",
		"Unexpected token at (1:2)", nil)
}

func TestFail66(t *testing.T) {
	// testFail(t, "break\n",
	// 	"Unsyntactic break (1:0)", nil)
}

func TestFail67(t *testing.T) {}

func TestFail68(t *testing.T) {}

func TestFail69(t *testing.T) {}

func TestFail70(t *testing.T) {}

func TestFail71(t *testing.T) {}

func TestFail72(t *testing.T) {}

func TestFail73(t *testing.T) {}

func TestFail74(t *testing.T) {}

func TestFail75(t *testing.T) {}

func TestFail76(t *testing.T) {}

func TestFail77(t *testing.T) {}

func TestFail78(t *testing.T) {}

func TestFail79(t *testing.T) {}

func TestFail80(t *testing.T) {}

func TestFail81(t *testing.T) {}

func TestFail82(t *testing.T) {}

func TestFail83(t *testing.T) {}

func TestFail84(t *testing.T) {}

func TestFail85(t *testing.T) {}

func TestFail86(t *testing.T) {}

func TestFail87(t *testing.T) {}

func TestFail88(t *testing.T) {}

func TestFail89(t *testing.T) {}

func TestFail90(t *testing.T) {}

func TestFail91(t *testing.T) {}

func TestFail92(t *testing.T) {}

func TestFail93(t *testing.T) {}

func TestFail94(t *testing.T) {}

func TestFail95(t *testing.T) {}

func TestFail96(t *testing.T) {}

func TestFail97(t *testing.T) {}

func TestFail98(t *testing.T) {}

func TestFail99(t *testing.T) {}

func TestFail100(t *testing.T) {}

func TestFail101(t *testing.T) {}

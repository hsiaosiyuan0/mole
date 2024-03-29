package parser

import (
	"testing"

	span "github.com/hsiaosiyuan0/mole/span"
	. "github.com/hsiaosiyuan0/mole/util"
)

func newParser(code string, opts *ParserOpts) *Parser {
	if opts == nil {
		opts = NewParserOpts()
	}
	s := span.NewSource("", code)
	return NewParser(s, opts)
}

func compile(code string, opts *ParserOpts) (Node, *Parser, error) {
	p := newParser(code, opts)
	prog, err := p.Prog()
	return prog, p, err
}

func testFail(t *testing.T, code, errMs string, opts *ParserOpts) {
	ast, _, err := compile(code, opts)
	if err == nil {
		t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
	}
	AssertEqual(t, errMs, err.Error(), "")
}

func testPass(t *testing.T, code string, opts *ParserOpts) {
	_, _, err := compile(code, opts)
	if err != nil {
		t.Fatalf("should pass code:\n%s\nerr:\n%v", code, err)
	}
}

func TestExpr(t *testing.T) {
	ast, p, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	AssertEqual(t, "-", expr.OpText(), "should be -")

	ab := expr.lhs.(*BinExpr)

	a := ab.lhs
	AssertEqual(t, "a", p.NodeText(a), "should be name a")
	AssertEqual(t, "+", ab.OpText(), "should be +")
	b := ab.rhs
	AssertEqual(t, "b", p.NodeText(b), "should be name b")

	c := expr.rhs
	AssertEqual(t, "c", p.NodeText(c), "should be name c")
}

func TestExprPcdHigherRight(t *testing.T) {
	ast, p, err := compile("a + b * c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	AssertEqual(t, N_NAME, lhs.Type(), "should be name")
	AssertEqual(t, "a", p.NodeText(lhs), "should be name a")

	rhs := expr.rhs
	AssertEqual(t, N_EXPR_BIN, rhs.Type(), "should be bin *")

	lhs = rhs.(*BinExpr).lhs
	AssertEqual(t, N_NAME, lhs.Type(), "should be name")
	AssertEqual(t, "b", p.NodeText(lhs), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	AssertEqual(t, N_NAME, rhs.Type(), "should be name")
	AssertEqual(t, "c", p.NodeText(rhs), "should be name c")
}

func TestExprPcdHigherLeft(t *testing.T) {
	ast, p, err := compile("a * b + c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	AssertEqual(t, "+", expr.OpText(), "should be +")

	ab := expr.lhs.(*BinExpr)
	AssertEqual(t, "*", ab.OpText(), "should be *")
	a := ab.lhs
	AssertEqual(t, "a", p.NodeText(a), "should be name a")
	b := ab.rhs
	AssertEqual(t, "b", p.NodeText(b), "should be name b")

	c := expr.rhs
	AssertEqual(t, "c", p.NodeText(c), "should be name c")
}

func TestExprAssoc(t *testing.T) {
	ast, p, err := compile("a ** b ** c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)

	lhs := expr.lhs
	AssertEqual(t, N_NAME, lhs.Type(), "should be name")
	AssertEqual(t, "a", p.NodeText(lhs), "should be name a")

	rhs := expr.rhs
	AssertEqual(t, N_EXPR_BIN, rhs.Type(), "should be bin **")

	lhs = rhs.(*BinExpr).lhs
	AssertEqual(t, N_NAME, lhs.Type(), "should be name")
	AssertEqual(t, "b", p.NodeText(lhs), "should be name b")

	rhs = rhs.(*BinExpr).rhs
	AssertEqual(t, N_NAME, rhs.Type(), "should be name")
	AssertEqual(t, "c", p.NodeText(rhs), "should be name c")
}

func TestCond(t *testing.T) {
	ast, _, err := compile("a > 0 ? a : b", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CondExpr)
	test := expr.test.(*BinExpr)
	AssertEqual(t, ">", test.OpText(), "should be >")
}

func TestAssign(t *testing.T) {
	ast, p, err := compile("a = a > 0 ? a : b", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	a := expr.lhs.(*Ident)
	AssertEqual(t, "a", p.NodeText(a), "should be a")

	cond := expr.rhs.(*CondExpr)
	test := cond.test.(*BinExpr)
	AssertEqual(t, ">", test.OpText(), "should be >")
	AssertEqual(t, "a", p.NodeText(a), "should be a")
}

func TestMemberExprSubscript(t *testing.T) {
	ast, p, err := compile("a[b][c]", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	AssertEqual(t, "a", p.NodeText(ab.obj), "should be a")
	AssertEqual(t, "b", p.NodeText(ab.prop), "should be b")
	AssertEqual(t, "c", p.NodeText(expr.prop), "should be c")
}

func TestMemberExprDot(t *testing.T) {
	ast, p, err := compile("a.b.c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	ab := expr.obj.(*MemberExpr)

	AssertEqual(t, "a", p.NodeText(ab.obj), "should be a")
	AssertEqual(t, "b", p.NodeText(ab.prop), "should be b")
	AssertEqual(t, "c", p.NodeText(expr.prop), "should be c")
}

func TestUnaryExpr(t *testing.T) {
	ast, p, err := compile("a + void 0", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	a := expr.lhs.(*Ident)
	AssertEqual(t, "a", p.NodeText(a), "should be a")

	v0 := expr.rhs.(*UnaryExpr)
	AssertEqual(t, "void", v0.OpText(), "should be void")
	AssertEqual(t, "0", p.NodeText(v0.arg), "should be 0")
}

func TestUpdateExpr(t *testing.T) {
	ast, p, err := compile("a + ++b + c++", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*BinExpr)
	ab := expr.lhs.(*BinExpr)
	AssertEqual(t, "a", p.NodeText(ab.lhs), "should be a")

	u1 := ab.rhs.(*UpdateExpr)
	AssertEqual(t, "b", p.NodeText(u1.arg), "should be b")
	AssertEqual(t, true, u1.prefix, "should be prefix")

	u2 := expr.rhs.(*UpdateExpr)
	AssertEqual(t, "c", p.NodeText(u2.arg), "should be c")
	AssertEqual(t, false, u2.prefix, "should be postfix")
}

func TestNewExpr(t *testing.T) {
	ast, p, err := compile("new new a", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*NewExpr).callee.(*NewExpr)
	AssertEqual(t, "a", p.NodeText(expr.callee), "should be a")
}

func TestCallExpr(t *testing.T) {
	ast, p, err := compile("a()(c, b, ...a)", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)
	callee := expr.callee.(*CallExpr)
	AssertEqual(t, "a", p.NodeText(callee.callee), "should be b")

	params := expr.args
	AssertEqual(t, "c", p.NodeText(params[0]), "should be c")
	AssertEqual(t, "a", p.NodeText(params[2].(*Spread).arg), "should be b")
	AssertEqual(t, "b", p.NodeText(params[1]), "should be a")
}

func TestCallExprMem(t *testing.T) {
	ast, p, err := compile("a(b).c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	obj := expr.obj.(*CallExpr)
	callee := obj.callee.(*Ident)
	AssertEqual(t, "a", p.NodeText(callee), "should be a")
	AssertEqual(t, "c", p.NodeText(expr.prop), "should be c")

	params := obj.args
	AssertEqual(t, "b", p.NodeText(params[0]), "should be b")
}

func TestCallExprLit(t *testing.T) {
	_, _, err := compile("a('b')", nil)
	AssertEqual(t, nil, err, "should be prog ok")
}

func TestCallCascadeExpr(t *testing.T) {
	ast, p, err := compile("a[b][c]()[d][e]()", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*CallExpr)

	// a[b][c]()[d][e]
	expr0 := expr.callee.(*MemberExpr)

	// a[b][c]()[d]
	expr1 := expr0.obj.(*MemberExpr)
	e := expr0.prop.(*Ident)
	AssertEqual(t, "e", p.NodeText(e), "should be e")

	// a[b][c]()
	expr2 := expr1.obj.(*CallExpr)
	d := expr1.prop.(*Ident)
	AssertEqual(t, "d", p.NodeText(d), "should be d")

	// a[b][c]
	expr3 := expr2.callee.(*MemberExpr)
	c := expr3.prop.(*Ident)
	AssertEqual(t, "c", p.NodeText(c), "should be c")

	// a[b]
	expr4 := expr3.obj.(*MemberExpr)
	b := expr4.prop.(*Ident)
	AssertEqual(t, "b", p.NodeText(b), "should be b")
	AssertEqual(t, "a", p.NodeText(expr4.obj), "should be a")
}

func TestVarDec(t *testing.T) {
	ast, p, err := compile("var a = 1", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]
	id := varDec.(*VarDec).id.(*Ident)
	init := varDec.(*VarDec).init.(*NumLit)
	AssertEqual(t, "a", p.NodeText(id), "should be a")
	AssertEqual(t, "1", p.NodeText(init), "should be 1")
}

func TestVarDecArrPattern(t *testing.T) {
	ast, p, err := compile("var [a, b = 1, [c] = 1, [d = 1]] = e", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.(*VarDec).init.(*Ident)
	AssertEqual(t, "e", p.NodeText(init), "should be e")

	arr := varDec.(*VarDec).id.(*ArrPat)
	elem0 := arr.elems[0].(*Ident)
	AssertEqual(t, "a", p.NodeText(elem0), "should be a")

	elem1 := arr.elems[1].(*AssignPat)
	elem1Lhs := elem1.lhs.(*Ident)
	elem1Rhs := elem1.rhs.(*NumLit)
	AssertEqual(t, "b", p.NodeText(elem1Lhs), "should be b")
	AssertEqual(t, "1", p.NodeText(elem1Rhs), "should be 1")

	elem2 := arr.elems[2].(*AssignPat)
	elem2Lhs := elem2.lhs.(*ArrPat)
	elem2Rhs := elem2.rhs.(*NumLit)
	AssertEqual(t, "c", p.NodeText(elem2Lhs.elems[0]), "should be c")
	AssertEqual(t, "1", p.NodeText(elem2Rhs), "should be 1")

	elem3 := arr.elems[3].(*ArrPat)
	elem31 := elem3.elems[0].(*AssignPat)
	AssertEqual(t, "d", p.NodeText(elem31.lhs), "should be d")
	AssertEqual(t, "1", p.NodeText(elem31.rhs), "should be 1")
}

func TestVarDecArrPatternElision(t *testing.T) {
	ast, p, err := compile("var [a, , b, , , c, ,] = e", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	init := varDec.(*VarDec).init.(*Ident)
	AssertEqual(t, "e", p.NodeText(init), "should be e")

	arr := varDec.(*VarDec).id.(*ArrPat)
	AssertEqual(t, 7, len(arr.elems), "should be len 7")

	elem0 := arr.elems[0].(*Ident)
	AssertEqual(t, "a", p.NodeText(elem0), "should be a")

	elem1 := arr.elems[1]
	AssertEqual(t, nil, elem1, "should be nil")

	elem2 := arr.elems[2].(*Ident)
	AssertEqual(t, "b", p.NodeText(elem2), "should be b")

	elem3 := arr.elems[3]
	AssertEqual(t, nil, elem3, "should be nil")

	elem4 := arr.elems[4]
	AssertEqual(t, nil, elem4, "should be nil")

	elem5 := arr.elems[5].(*Ident)
	AssertEqual(t, "c", p.NodeText(elem5), "should be c")

	elem6 := arr.elems[6]
	AssertEqual(t, nil, elem6, "should be nil")
}

func TestArrLit(t *testing.T) {
	ast, p, err := compile("[a, , b, , , c, ,]", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	arrLit := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*ArrLit)
	AssertEqual(t, 7, len(arrLit.elems), "should be len 7")

	elem0 := arrLit.elems[0].(*Ident)
	AssertEqual(t, "a", p.NodeText(elem0), "should be a")

	elem1 := arrLit.elems[1]
	AssertEqual(t, nil, elem1, "should be nil")

	elem2 := arrLit.elems[2].(*Ident)
	AssertEqual(t, "b", p.NodeText(elem2), "should be b")

	elem3 := arrLit.elems[3]
	AssertEqual(t, nil, elem3, "should be nil")

	elem4 := arrLit.elems[4]
	AssertEqual(t, nil, elem4, "should be nil")

	elem5 := arrLit.elems[5].(*Ident)
	AssertEqual(t, "c", p.NodeText(elem5), "should be c")

	elem6 := arrLit.elems[6]
	AssertEqual(t, nil, elem6, "should be nil")
}

func TestObjLit(t *testing.T) {
	ast, p, err := compile(`var a = {...a, b, ...c, "d": 1, [e]: {f: 1}, ...g}`, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.(*VarDec).id.(*Ident)
	AssertEqual(t, "a", p.NodeText(id), "should be a")

	objLit := varDec.(*VarDec).init.(*ObjLit)
	AssertEqual(t, 6, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Spread)
	AssertEqual(t, "a", p.NodeText(prop0.arg), "should be ...a")

	prop1 := objLit.props[1].(*Prop)
	AssertEqual(t, "b", p.NodeText(prop1.key), "should be b")

	prop2 := objLit.props[2].(*Spread)
	AssertEqual(t, "c", p.NodeText(prop2.arg), "should be ...c")

	prop3 := objLit.props[3].(*Prop)
	AssertEqual(t, "d", prop3.key.(*StrLit).val, "should be d")
	AssertEqual(t, "1", p.NodeText(prop3.value), "should be 1")

	prop4 := objLit.props[4].(*Prop)
	AssertEqual(t, "e", p.NodeText(prop4.key), "should be e")
	AssertEqual(t, "f", p.NodeText(prop4.value.(*ObjLit).props[0].(*Prop).key), "should be f")
	AssertEqual(t, "1", p.NodeText(prop4.value.(*ObjLit).props[0].(*Prop).value), "should be 1")

	prop5 := objLit.props[5].(*Spread)
	AssertEqual(t, "g", p.NodeText(prop5.arg), "should be ...g")
}

func TestObjLitMethod(t *testing.T) {
	ast, p, err := compile(`
  var o = {
    a,
    [b] () {},
    c,
    e: () => {},
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	varDecStmt := ast.(*Prog).stmts[0].(*VarDecStmt)
	varDec := varDecStmt.decList[0]

	id := varDec.(*VarDec).id.(*Ident)
	AssertEqual(t, "o", p.NodeText(id), "should be o")

	objLit := varDec.(*VarDec).init.(*ObjLit)
	AssertEqual(t, 4, len(objLit.props), "should be len 6")

	prop0 := objLit.props[0].(*Prop)
	AssertEqual(t, "a", p.NodeText(prop0.key), "should be a")

	prop1 := objLit.props[1].(*Prop)
	AssertEqual(t, "b", p.NodeText(prop1.key), "should be b")
	_ = prop1.value.(*FnDec)

	prop2 := objLit.props[2].(*Prop)
	AssertEqual(t, "c", p.NodeText(prop2.key), "should be c")

	prop3 := objLit.props[3].(*Prop)
	AssertEqual(t, "e", p.NodeText(prop3.key), "should be e")
	_ = prop3.value.(*ArrowFn)
}

func TestFnDec(t *testing.T) {
	ast, p, err := compile(`
  function a({ b }) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	AssertEqual(t, "a", p.NodeText(id), "should be a")
}

func TestFnExpr(t *testing.T) {
	ast, p, err := compile(`
  let a = function a({ b }) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*FnDec)
	id := fn.id.(*Ident)
	AssertEqual(t, "a", p.NodeText(id), "should be a")
}

func TestAsyncFnDec(t *testing.T) {
	ast, p, err := compile(`
  async function a({ b }) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	AssertEqual(t, "a", p.NodeText(id), "should be a")
	AssertEqual(t, true, fn.async, "should be true")
}

func TestArrowFn(t *testing.T) {
	ast, _, err := compile(`
  a = () => {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	expr := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = expr.rhs.(*ArrowFn)
}

func TestDoWhileStmt(t *testing.T) {
	ast, _, err := compile(`
  do {} while(1)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*DoWhileStmt)
}

func TestWhileStmt(t *testing.T) {
	ast, _, err := compile(`
  while(1) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*WhileStmt)
}

func TestForStmt(t *testing.T) {
	ast, _, err := compile(`
  for(;;) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForStmt)
}

func TestForInStmt(t *testing.T) {
	ast, _, err := compile(`
  for (a in b) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)
}

func TestForOfStmt(t *testing.T) {
	ast, _, err := compile(`
  for (a of b) {}
  for await (a of b) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ForInOfStmt)

	forAwait := ast.(*Prog).stmts[1].(*ForInOfStmt)
	AssertEqual(t, true, forAwait.await, "should be await")
}

func TestIfStmt(t *testing.T) {
	ast, p, err := compile(`
  if (a) {} else b
  if (c) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	stmt := ast.(*Prog).stmts[0].(*IfStmt)
	AssertEqual(t, "a", p.NodeText(stmt.test), "should be a")

	stmt = ast.(*Prog).stmts[1].(*IfStmt)
	AssertEqual(t, "c", p.NodeText(stmt.test), "should be c")
}

func TestSwitchStmtEmpty(t *testing.T) {
	ast, _, err := compile(`
	switch (a) {
	}
	`, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*SwitchStmt)
}

func TestSwitchStmt(t *testing.T) {
	ast, p, err := compile(`
  switch (a) {
    case b in c:
      d
      e
    case f:
    default:
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.(*SwitchCase).test.(*BinExpr)
	AssertEqual(t, "b", p.NodeText(test0.lhs), "should be prog b")
	AssertEqual(t, "c", p.NodeText(test0.rhs), "should be prog c")

	cons00 := case0.(*SwitchCase).cons[0].(*ExprStmt)
	AssertEqual(t, "d", p.NodeText(cons00.expr), "should be prog d")

	cons01 := case0.(*SwitchCase).cons[1].(*ExprStmt)
	AssertEqual(t, "e", p.NodeText(cons01.expr), "should be prog e")

	case1 := stmt.cases[1]
	AssertEqual(t, "f", p.NodeText(case1.(*SwitchCase).test), "should be prog f")

	case2 := stmt.cases[2]
	AssertEqual(t, nil, case2.(*SwitchCase).test, "should be default")
}

func TestSwitchStmtDefaultMiddle(t *testing.T) {
	ast, p, err := compile(`
  switch (a) {
    case b in c:
      d
      e
    default:
    case f:
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt := ast.(*Prog).stmts[0].(*SwitchStmt)

	case0 := stmt.cases[0]
	test0 := case0.(*SwitchCase).test.(*BinExpr)
	AssertEqual(t, "b", p.NodeText(test0.lhs), "should be prog b")
	AssertEqual(t, "c", p.NodeText(test0.rhs), "should be prog c")

	cons00 := case0.(*SwitchCase).cons[0].(*ExprStmt)
	AssertEqual(t, "d", p.NodeText(cons00.expr), "should be prog d")

	cons01 := case0.(*SwitchCase).cons[1].(*ExprStmt)
	AssertEqual(t, "e", p.NodeText(cons01.expr), "should be prog e")

	case1 := stmt.cases[1]
	AssertEqual(t, nil, case1.(*SwitchCase).test, "should be default")

	case2 := stmt.cases[2]
	AssertEqual(t, "f", p.NodeText(case2.(*SwitchCase).test), "should be prog f")
}

func TestBrkStmt(t *testing.T) {
	ast, p, err := compile(`
  while(true) break
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	whileStmt := ast.(*Prog).stmts[0].(*WhileStmt)
	_ = whileStmt.body.(*BrkStmt)

	ast, p, err = compile(`
  a: while(true) break a;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	labelStmt := ast.(*Prog).stmts[0].(*LabelStmt)
	whileStmt = labelStmt.body.(*WhileStmt)
	stmt := whileStmt.body.(*BrkStmt)
	AssertEqual(t, "a", p.NodeText(stmt.label), "should be a")
}

func TestContStmt(t *testing.T) {
	ast, p, err := compile(`
  while(true) continue
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	whileStmt := ast.(*Prog).stmts[0].(*WhileStmt)
	_ = whileStmt.body.(*ContStmt)

	ast, _, err = compile(`
  a: while(true) continue a;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	labelStmt := ast.(*Prog).stmts[0].(*LabelStmt)
	whileStmt = labelStmt.body.(*WhileStmt)
	stmt := whileStmt.body.(*ContStmt)
	AssertEqual(t, "a", p.NodeText(stmt.label), "should be a")
}

func TestLabelStmt(t *testing.T) {
	ast, p, err := compile(`
  a:
  b
  c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	lbStmt := ast.(*Prog).stmts[0].(*LabelStmt)
	AssertEqual(t, "a", p.NodeText(lbStmt.label), "should be a")

	lbBody := lbStmt.body.(*ExprStmt)
	AssertEqual(t, "b", p.NodeText(lbBody.expr), "should be b")

	expr := ast.(*Prog).stmts[1].(*ExprStmt)
	AssertEqual(t, "c", p.NodeText(expr.expr), "should be c")

	ast, _, err = compile(`
  a: b
  c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	lbStmt = ast.(*Prog).stmts[0].(*LabelStmt)
	AssertEqual(t, "a", p.NodeText(lbStmt.label), "should be a")

	lbBody = lbStmt.body.(*ExprStmt)
	AssertEqual(t, "b", p.NodeText(lbBody.expr), "should be b")

	expr = ast.(*Prog).stmts[1].(*ExprStmt)
	AssertEqual(t, "c", p.NodeText(expr.expr), "should be c")
}

func TestRetStmt(t *testing.T) {
	ast, p, err := compile(`
  function a() { return a }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*FnDec)
	AssertEqual(t, "a", p.NodeText(stmt0.body.(*BlockStmt).body[0].(*RetStmt).arg), "should be a")

	ast, _, err = compile(`
  function a() {
    return
    a
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	stmt0 = ast.(*Prog).stmts[0].(*FnDec)
	AssertEqual(t, nil, stmt0.body.(*BlockStmt).body[0].(*RetStmt).arg, "should be nil")

	stmt1 := stmt0.body.(*BlockStmt).body[1].(*ExprStmt)
	AssertEqual(t, "a", p.NodeText(stmt1.expr), "should be a")
}

func TestThrowStmt(t *testing.T) {
	ast, p, err := compile(`
  throw a
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*ThrowStmt)
	AssertEqual(t, "a", p.NodeText(stmt0.arg), "should be a")

	_, p, err = compile(`
  throw
  a
  `, nil)
	AssertEqual(t, true, err != nil, "should be failed")
}

func TestTryStmt(t *testing.T) {
	ast, p, err := compile(`
  try {} catch(e) {} finally {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	stmt0 := ast.(*Prog).stmts[0].(*TryStmt)
	catch := stmt0.catch
	AssertEqual(t, "e", p.NodeText(catch.(*Catch).param), "should be e")

	AssertEqual(t, true, stmt0.fin != nil, "should have fin")

	ast, p, err = compile(`
  try {} finally {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt0 = ast.(*Prog).stmts[0].(*TryStmt)
	AssertEqual(t, true, stmt0.fin != nil, "should have fin")

	_, p, err = compile(`
  try {}
  `, nil)
	AssertEqual(t, true, err != nil, "should be err")
}

func TestDebugStmt(t *testing.T) {
	ast, _, err := compile(`
  a
  debugger
  b
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*ExprStmt)
	_ = ast.(*Prog).stmts[1].(*DebugStmt)
	_ = ast.(*Prog).stmts[2].(*ExprStmt)
}

func TestEmptyStmt(t *testing.T) {
	ast, _, err := compile(`
  ;a;;
  ;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	_ = ast.(*Prog).stmts[0].(*EmptyStmt)
	_ = ast.(*Prog).stmts[1].(*ExprStmt)
	_ = ast.(*Prog).stmts[2].(*EmptyStmt)
	_ = ast.(*Prog).stmts[3].(*EmptyStmt)
}

func TestClassStmt(t *testing.T) {
	ast, _, err := compile(`
  class a {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	_ = ast.(*Prog).stmts[0].(*ClassDec)
}

func TestClassField(t *testing.T) {
	ast, p, err := compile(`
  class a {
    #f1
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.(*ClassBody).elems[0].(*Field)
	AssertEqual(t, true, elem0.key.(*Ident).pvt, "should be pvt")
	AssertEqual(t, "f1", p.NodeText(elem0.key), "should be f1")
}

func TestClassMethod(t *testing.T) {
	ast, p, err := compile(`
  class a {
    [a] (b) {
      c
    }

    e
    #f () {}
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	cls := ast.(*Prog).stmts[0].(*ClassDec)
	elem0 := cls.body.(*ClassBody).elems[0].(*Method)
	AssertEqual(t, "a", p.NodeText(elem0.key), "should be a")
	AssertEqual(t, "b", p.NodeText(elem0.val.(*FnDec).params[0]), "should be b")

	elem1 := cls.body.(*ClassBody).elems[1].(*Field)
	AssertEqual(t, "e", p.NodeText(elem1.key), "should be e")

	elem2 := cls.body.(*ClassBody).elems[2].(*Method)
	AssertEqual(t, true, elem2.key.(*Ident).pvt, "should be pvt")
	AssertEqual(t, "f", p.NodeText(elem2.key), "should be f")
}

func TestSeqExpr(t *testing.T) {
	ast, p, err := compile(`
  a = (b, c)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	elem0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	seq := elem0.rhs.(*ParenExpr).expr.(*SeqExpr)
	AssertEqual(t, "b", p.NodeText(seq.elems[0]), "should be b")
	AssertEqual(t, "c", p.NodeText(seq.elems[1]), "should be c")
}

func TestClassExpr(t *testing.T) {
	ast, _, err := compile(`
  a = class {};
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*ClassDec)
}

func TestRegexpExpr(t *testing.T) {
	ast, _, err := compile(`
  a = /a/
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*RegLit)
}

func TestParenExpr(t *testing.T) {
	ast, _, err := compile(`
  a = (b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	stmt0 := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	_ = stmt0.rhs.(*ParenExpr).expr.(*Ident)
}

func TestTplExpr(t *testing.T) {
	ast, p, err := compile("tag`\na${b}c`", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	AssertEqual(t, "tag", p.NodeText(tag), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	AssertEqual(t, "\na", span0.val, "should be a")

	span1 := tpl.elems[1].(*Ident)
	AssertEqual(t, "b", p.NodeText(span1), "should be b")

	span2 := tpl.elems[2].(*StrLit)
	AssertEqual(t, "c", span2.val, "should be c")
}

func TestTplExprNest(t *testing.T) {
	ast, p, err := compile("tag`\na${ f`g\n${d}e` }c`", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	tpl := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*TplExpr)
	tag := tpl.tag.(*Ident)
	AssertEqual(t, "tag", p.NodeText(tag), "should be tag")

	span0 := tpl.elems[0].(*StrLit)
	AssertEqual(t, "\na", span0.val, "should be a")

	span2 := tpl.elems[2].(*StrLit)
	AssertEqual(t, "c", span2.val, "should be c")

	tpl = tpl.elems[1].(*TplExpr)
	tag = tpl.tag.(*Ident)
	AssertEqual(t, "f", p.NodeText(tag), "should be f")

	span0 = tpl.elems[0].(*StrLit)
	AssertEqual(t, "g\n", span0.val, "should be g")

	span2 = tpl.elems[2].(*StrLit)
	AssertEqual(t, "e", span2.val, "should be e")

	span1 := tpl.elems[1].(*Ident)
	AssertEqual(t, "d", p.NodeText(span1), "should be d")
}

func TestTplExprMember(t *testing.T) {
	ast, p, err := compile("tag`\na${b}c`[d]", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	member := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*MemberExpr)
	tpl := member.obj.(*TplExpr)
	tag := tpl.tag.(*Ident)
	AssertEqual(t, "tag", p.NodeText(tag), "should be tag")
}

func TestSuper(t *testing.T) {
	ast, _, err := compile("class a extends b { constructor() { super() } }", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	ctor := ast.(*Prog).stmts[0].(*ClassDec).body.(*ClassBody).elems[0].(*Method).val.(*FnDec)
	expr := ctor.body.(*BlockStmt).body[0].(*ExprStmt).expr
	call := expr.(*CallExpr)
	AssertEqual(t, N_SUPER, call.callee.Type(), "should be tag")
}

func TestImportCall(t *testing.T) {
	ast, p, err := compile("a = import(b)", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	assign := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	importCall := assign.rhs.(*ImportCall)
	AssertEqual(t, "b", p.NodeText(importCall.src), "should be b")
}

func TestMetaProp(t *testing.T) {
	ast, p, err := compile("a = import.meta", nil)
	AssertEqual(t, nil, err, "should be prog ok")
	assign := ast.(*Prog).stmts[0].(*ExprStmt).expr.(*AssignExpr)
	metaProp := assign.rhs.(*MetaProp)
	AssertEqual(t, "meta", p.NodeText(metaProp.prop), "should be meta")
}

func TestScopeBalance(t *testing.T) {
	parser := newParser("function a () {}", nil)
	parser.Prog()
	AssertEqual(t, 0, parser.symtab.Cur.Id, "scope should be balanced")
}

func TestLabelledUsage(t *testing.T) {
	ast, _, err := compile(`
LabelA: for (;;) {
  break LabelA;
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	AssertEqual(t, true, ast.(*Prog).stmts[0].(*LabelStmt).Used(), "should be meta")
}

func TestLabelledUsageCont(t *testing.T) {
	ast, _, err := compile(`
LabelA: for (;;) {
  for (;;) {
    continue LabelA;
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	AssertEqual(t, true, ast.(*Prog).stmts[0].(*LabelStmt).Used(), "should be meta")
}

func TestLabelledUsageContNested(t *testing.T) {
	ast, _, err := compile(`
LabelA: for (;;) {
  LabelB: for (;;) {
    continue LabelB;
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	AssertEqual(t, false, ast.(*Prog).stmts[0].(*LabelStmt).Used(), "should be meta")
	AssertEqual(t, true, ast.(*Prog).stmts[0].(*LabelStmt).body.(*ForStmt).Body().(*BlockStmt).body[0].(*LabelStmt).Used(), "should be meta")
}

func TestLabelledUsageNoUse(t *testing.T) {
	ast, _, err := compile(`
LabelA: for (;;) {
  break;
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	AssertEqual(t, false, ast.(*Prog).stmts[0].(*LabelStmt).Used(), "should be meta")
}

func TestLabelledUsageContNoUse(t *testing.T) {
	ast, _, err := compile(`
LabelA: for (;;) {
  LabelB: for (;;) {
    continue LabelB;
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")
	AssertEqual(t, false, ast.(*Prog).stmts[0].(*LabelStmt).Used(), "should be meta")
}

func TestFail1(t *testing.T) {
	testFail(t, "{", "Unexpected token `EOF` at (1:1)", nil)
}

func TestFail2(t *testing.T) {
	testFail(t, "}", "Unexpected token `}` at (1:0)", nil)
}

func TestFail3(t *testing.T) {
	testFail(t, "3ea", "Identifier directly after number at (1:2)", nil)
}

func TestFail4(t *testing.T) {
	testFail(t, "3in []", "Identifier directly after number at (1:1)", nil)
}

func TestFail5(t *testing.T) {
	testFail(t, "3e", "Invalid number at (1:2)", nil)
}

func TestFail6(t *testing.T) {
	testFail(t, "3e+", "Invalid number at (1:3)", nil)
}

func TestFail7(t *testing.T) {
	testFail(t, "3e-", "Invalid number at (1:3)", nil)
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
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "01a", "Identifier directly after number at (1:2)", opts)
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
	testFail(t, "/", "Unterminated regular expression at (1:0)", nil)
}

func TestFail19(t *testing.T) {
	testFail(t, "/test", "Unterminated regular expression at (1:0)", nil)
}

func TestFail20(t *testing.T) {
	testFail(t, "var x = /[a-z]/\\ux", "Bad character escape sequence at (1:15)", nil)
}

func TestFail21(t *testing.T) {
	testFail(t, "3 = 4", "Assigning to rvalue at (1:0)", nil)
}

func TestFail22(t *testing.T) {
	testFail(t, "func() = 4", "Assigning to rvalue at (1:0)", nil)
}

func TestFail23(t *testing.T) {
	testFail(t, "(1 + 1) = 10", "Assigning to rvalue at (1:1)", nil)
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
	testFail(t, "var x = /\n/", "Unterminated regular expression at (1:9)", nil)
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
		"Rest element must be last element at (1:18)", nil)
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
	testFail(t, "break\n",
		"Illegal break at (1:0)", nil)
}

func TestFail67(t *testing.T) {
	testFail(t, "break 1;",
		"Unexpected token at (1:6)", nil)
}

func TestFail68(t *testing.T) {
	testFail(t, "continue\n",
		"Illegal continue at (1:0)", nil)
}

func TestFail69(t *testing.T) {
	testFail(t, "continue 2;",
		"Unexpected token at (1:9)", nil)
}

func TestFail70(t *testing.T) {
	testFail(t, "throw",
		"Unexpected token `EOF` at (1:5)", nil)
}

func TestFail71(t *testing.T) {
	testFail(t, "throw;",
		"Unexpected token `;` at (1:5)", nil)
}

func TestFail72(t *testing.T) {
	testFail(t, "for (var i, i2 in {});",
		"Must have a single binding at (1:12)", nil)
}

func TestFail73(t *testing.T) {
	testFail(t, "for ((i in {}));",
		"Unexpected token `)` at (1:14)", nil)
}

func TestFail74(t *testing.T) {
	testFail(t, "for (i + 1 in {});",
		"Assigning to rvalue at (1:5)", nil)
}

func TestFail75(t *testing.T) {
	testFail(t, "for (+i in {});",
		"Assigning to rvalue at (1:5)", nil)
}

func TestFail76(t *testing.T) {
	testFail(t, "if(false)",
		"Unexpected token `EOF` at (1:9)", nil)
}

func TestFail77(t *testing.T) {
	testFail(t, "if(false) doThis(); else",
		"Unexpected token `EOF` at (1:24)", nil)
}

func TestFail78(t *testing.T) {
	testFail(t, "do",
		"Unexpected token `EOF` at (1:2)", nil)
}

func TestFail79(t *testing.T) {
	testFail(t, "while(false)",
		"Unexpected token `EOF` at (1:12)", nil)
}

func TestFail80(t *testing.T) {
	testFail(t, "for(;;)",
		"Unexpected token `EOF` at (1:7)", nil)
}

func TestFail81(t *testing.T) {
	testFail(t, "with(x)",
		"Unexpected token `EOF` at (1:7)", nil)
}

func TestFail82(t *testing.T) {
	testFail(t, "try { }",
		"Unexpected token `EOF` at (1:7)", nil)
}

func TestFail83(t *testing.T) {
	testFail(t, "‿ = 10",
		"Unexpected character at (1:0)", nil)
}

func TestFail84(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "if(true) let a = 1;",
		"Unexpected token `identifier` at (1:9)", opts)
}

func TestFail85(t *testing.T) {
	testFail(t, "switch (c) { default: default: }",
		"Multiple default clauses at (1:22)", nil)
}

func TestFail86(t *testing.T) {
	testFail(t, "new X().\"s\"",
		"Unexpected token `string` at (1:8)", nil)
}

func TestFail87(t *testing.T) {
	testFail(t, "/*",
		"Unterminated comment at (1:0)", nil)
}

func TestFail88(t *testing.T) {
	testFail(t, "/*\n\n\n",
		"Unterminated comment at (1:0)", nil)
}

func TestFail89(t *testing.T) {
	testFail(t, "/**",
		"Unterminated comment at (1:0)", nil)
}

func TestFail90(t *testing.T) {
	testFail(t, "/*\n\n*",
		"Unterminated comment at (1:0)", nil)
}

func TestFail91(t *testing.T) {
	testFail(t, "/*hello",
		"Unterminated comment at (1:0)", nil)
}

func TestFail92(t *testing.T) {
	testFail(t, "/*hello  *",
		"Unterminated comment at (1:0)", nil)
}

func TestFail93(t *testing.T) {
	testFail(t, "\n]",
		"Unexpected token `]` at (2:0)", nil)
}

func TestFail94(t *testing.T) {
	testFail(t, "\r]",
		"Unexpected token `]` at (2:0)", nil)
}

func TestFail95(t *testing.T) {
	testFail(t, "\r\n]",
		"Unexpected token `]` at (2:0)", nil)
}

func TestFail96(t *testing.T) {
	testFail(t, "\n\r]",
		"Unexpected token `]` at (3:0)", nil)
}

func TestFail97(t *testing.T) {
	testFail(t, "//\r\n]",
		"Unexpected token `]` at (2:0)", nil)
}

func TestFail98(t *testing.T) {
	testFail(t, "//\n\r]",
		"Unexpected token `]` at (3:0)", nil)
}

func TestFail99(t *testing.T) {
	testFail(t, "/a\\\n/",
		"Unterminated regular expression at (1:3)", nil)
}

func TestFail100(t *testing.T) {
	testFail(t, "//\r \n]",
		"Unexpected token `]` at (3:0)", nil)
}

func TestFail101(t *testing.T) {
	testFail(t, "/*\r\n*/]",
		"Unexpected token `]` at (2:2)", nil)
}

func TestFail102(t *testing.T) {
	testFail(t, "/*\n\r*/]",
		"Unexpected token `]` at (3:2)", nil)
}

func TestFail103(t *testing.T) {
	testFail(t, "/*\r \n*/]",
		"Unexpected token `]` at (3:2)", nil)
}

func TestFail104(t *testing.T) {
	testFail(t, "\\\\",
		"Expecting Unicode escape sequence \\uXXXX at (1:0)", nil)
}

func TestFail105(t *testing.T) {
	testFail(t, "\\u005c",
		"Invalid Unicode escape at (1:0)", nil)
}

func TestFail106(t *testing.T) {
	testFail(t, "\\x",
		"Expecting Unicode escape sequence \\uXXXX at (1:0)", nil)
}

func TestFail107(t *testing.T) {
	testFail(t, "\\u0000",
		"Invalid Unicode escape at (1:0)", nil)
}

func TestFail108(t *testing.T) {
	//lint:ignore ST1018 lhs is `\u200c`
	testFail(t, "‌ = []",
		"Unexpected character at (1:0)", nil)
}

func TestFail109(t *testing.T) {
	//lint:ignore ST1018 lhs is `\u200d`
	testFail(t, "‍ = []",
		"Unexpected character at (1:0)", nil)
}

func TestFail110(t *testing.T) {
	testFail(t, "\"\\",
		"Unterminated string constant at (1:0)", nil)
}

func TestFail111(t *testing.T) {
	testFail(t, "\"\\u",
		"Bad character escape sequence at (1:2)", nil)
}

func TestFail112(t *testing.T) {
	testFail(t, "return",
		"Illegal return at (1:0)", nil)
}

func TestFail113(t *testing.T) {
	testFail(t, "break",
		"Illegal break at (1:0)", nil)
}

func TestFail114(t *testing.T) {
	testFail(t, "continue",
		"Illegal continue at (1:0)", nil)
}

func TestFail115(t *testing.T) {
	testFail(t, "switch (x) { default: continue; }",
		"Illegal continue at (1:22)", nil)
}

func TestFail116(t *testing.T) {
	testFail(t, "do { x } *",
		"Unexpected token `*` at (1:9)", nil)
}

func TestFail117(t *testing.T) {
	testFail(t, "while (true) { break x; }",
		"Undefined label `x` at (1:21)", nil)
}

func TestFail118(t *testing.T) {
	testFail(t, "while (true) { continue x; }",
		"Undefined label `x` at (1:24)", nil)
}

func TestFail119(t *testing.T) {
	testFail(t, "x: while (true) { (function () { break x; }); }",
		"Undefined label `x` at (1:39)", nil)
}

func TestFail120(t *testing.T) {
	testFail(t, "x: while (true) { (function () { continue x; }); }",
		"Undefined label `x` at (1:42)", nil)
}

func TestFail121(t *testing.T) {
	testFail(t, "x: while (true) { (function () { break; }); }",
		"Illegal break at (1:33)", nil)
}

func TestFail122(t *testing.T) {
	testFail(t, "x: while (true) { (function () { continue; }); }",
		"Illegal continue at (1:33)", nil)
}

func TestFail123(t *testing.T) {
	testFail(t, "x: while (true) { x: while (true) { } }",
		"Label `x` already declared at (1:18)", nil)
}

func TestFail124(t *testing.T) {
	testFail(t, "(function () { 'use strict'; delete i; }())",
		"Deleting local variable in strict mode at (1:36)", nil)
}

func TestFail125(t *testing.T) {
	testFail(t, "function x() { '\\12'; 'use strict'; }",
		"Octal escape sequences are not allowed in strict mode at (1:15)", nil)
}

func TestFail126(t *testing.T) {
	testFail(t, "function hello() {'use strict'; var eval = 10; }",
		"Unexpected token `eval` at (1:36)", nil)
}

func TestFail127(t *testing.T) {
	testFail(t, "function hello() {'use strict'; var arguments = 10; }",
		"Unexpected token `arguments` at (1:36)", nil)
}

func TestFail128(t *testing.T) {
	testFail(t, "function hello() {'use strict'; try { } catch (arguments) { } }",
		"Unexpected token `arguments` at (1:47)", nil)
}

func TestFail129(t *testing.T) {
	testFail(t, "function hello() {'use strict'; try { } catch (arguments) { } }",
		"Unexpected token `arguments` at (1:47)", nil)
}

func TestFail130(t *testing.T) {
	testFail(t, "function hello() {'use strict'; eval = 10; }",
		"Assigning to `eval` in strict mode at (1:32)", nil)
}

func TestFail131(t *testing.T) {
	testFail(t, "function hello() {'use strict'; arguments = 10; }",
		"Assigning to `arguments` in strict mode at (1:32)", nil)
}

func TestFail132(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ++eval; }",
		"Assigning to rvalue at (1:34)", nil)
}

func TestFail133(t *testing.T) {
	testFail(t, "function hello() {'use strict'; --eval; }",
		"Assigning to rvalue at (1:34)", nil)
}

func TestFail134(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ++arguments; }",
		"Assigning to rvalue at (1:34)", nil)
}

func TestFail135(t *testing.T) {
	testFail(t, "function hello() {'use strict'; --arguments; }",
		"Assigning to rvalue at (1:34)", nil)
}

func TestFail136(t *testing.T) {
	testFail(t, "function hello() {'use strict'; eval++; }",
		"Assigning to rvalue at (1:32)", nil)
}

func TestFail137(t *testing.T) {
	testFail(t, "function hello() {'use strict'; eval--; }",
		"Assigning to rvalue at (1:32)", nil)
}

func TestFail138(t *testing.T) {
	testFail(t, "function hello() {'use strict'; arguments++; }",
		"Assigning to rvalue at (1:32)", nil)
}

func TestFail139(t *testing.T) {
	testFail(t, "function hello() {'use strict'; arguments--; }",
		"Assigning to rvalue at (1:32)", nil)
}

func TestFail140(t *testing.T) {
	testFail(t, "function hello() {'use strict'; function eval() { } }",
		"Unexpected token `eval` at (1:41)", nil)
}

func TestFail141(t *testing.T) {
	testFail(t, "function hello() {'use strict'; function arguments() { } }",
		"Unexpected token `arguments` at (1:41)", nil)
}

func TestFail142(t *testing.T) {
	testFail(t, "function eval() {'use strict'; }",
		"Unexpected token `eval` at (1:9)", nil)
}

func TestFail143(t *testing.T) {
	testFail(t, "function arguments() {'use strict'; }",
		"Unexpected token `arguments` at (1:9)", nil)
}

func TestFail144(t *testing.T) {
	testFail(t, "function hello() {'use strict'; (function eval() { }()) }",
		"Unexpected token `eval` at (1:42)", nil)
}

func TestFail145(t *testing.T) {
	testFail(t, "function hello() {'use strict'; (function arguments() { }()) }",
		"Unexpected token `arguments` at (1:42)", nil)
}

func TestFail146(t *testing.T) {
	testFail(t, "(function eval() {'use strict'; })()",
		"Unexpected token `eval` at (1:10)", nil)
}

func TestFail147(t *testing.T) {
	testFail(t, "(function arguments() {'use strict'; })()",
		"Unexpected token `arguments` at (1:10)", nil)
}

func TestFail148(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ({ s: function eval() { } }); }",
		"Unexpected token `eval` at (1:47)", nil)
}

func TestFail149(t *testing.T) {
	testFail(t, "(function package() {'use strict'; })()",
		"Invalid binding `package` at (1:10)", nil)
}

func TestFail150(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ({ i: 10, set s(eval) { } }); }",
		"Invalid binding `eval` at (1:48)", nil)
}

func TestFail151(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ({ set s(eval) { } }); }",
		"Invalid binding `eval` at (1:41)", nil)
}

func TestFail152(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ({ s: function s(eval) { } }); }",
		"Invalid binding `eval` at (1:49)", nil)
}

func TestFail153(t *testing.T) {
	testFail(t, "function hello(eval) {'use strict';}",
		"Invalid binding `eval` at (1:15)", nil)
}

func TestFail154(t *testing.T) {
	testFail(t, "function hello(arguments) {'use strict';}",
		"Invalid binding `arguments` at (1:15)", nil)
}

func TestFail155(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; function inner(eval) {} }",
		"Invalid binding `eval` at (1:48)", nil)
}

func TestFail156(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; function inner(arguments) {} }",
		"Invalid binding `arguments` at (1:48)", nil)
}

func TestFail157(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; \"\\1\"; }",
		"Octal escape sequences are not allowed in strict mode at (1:33)", nil)
}

func TestFail158(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; \"\\00\"; }",
		"Octal escape sequences are not allowed in strict mode at (1:33)", nil)
}

func TestFail159(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; \"\\000\"; }",
		"Octal escape sequences are not allowed in strict mode at (1:33)", nil)
}

func TestFail160(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; 021; }",
		"Octal literals are not allowed in strict mode at (1:33)", nil)
}

func TestFail161(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; ({ \"\\1\": 42 }); }",
		"Octal escape sequences are not allowed in strict mode at (1:36)", nil)

}

func TestFail162(t *testing.T) {
	testFail(t, "function hello() { 'use strict'; ({ 021: 42 }); }",
		"Octal literals are not allowed in strict mode at (1:36)", nil)
}

func TestFail163(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; function inner() { \"octal directive\\1\"; } }",
		"Octal escape sequences are not allowed in strict mode at (1:52)", nil)
}

func TestFail164(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var implements; }",
		"Invalid binding `implements` at (1:37)", nil)
}

func TestFail165(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var interface; }",
		"Invalid binding `interface` at (1:37)", nil)
}

func TestFail166(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var package; }",
		"Invalid binding `package` at (1:37)", nil)
}

func TestFail167(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var private; }",
		"Unexpected token `private` at (1:37)", nil)
}

func TestFail168(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var protected; }",
		"Invalid binding `protected` at (1:37)", nil)
}

func TestFail169(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var public; }",
		"Invalid binding `public` at (1:37)", nil)
}

func TestFail170(t *testing.T) {
	testFail(t, "function hello() { \"use strict\"; var static; }",
		"Invalid binding `static` at (1:37)", nil)
}

func TestFail171(t *testing.T) {
	testFail(t, "function hello(static) { \"use strict\"; }",
		"Invalid binding `static` at (1:15)", nil)
}

func TestFail172(t *testing.T) {
	testFail(t, "function static() { \"use strict\"; }",
		"Invalid binding `static` at (1:9)", nil)
}

func TestFail173(t *testing.T) {
	testFail(t, "\"use strict\"; function static() { }",
		"Invalid binding `static` at (1:23)", nil)
}

func TestFail174(t *testing.T) {
	testFail(t, "function a(t, t) { \"use strict\"; }",
		"Parameter name clash at (1:14)", nil)
}

func TestFail175(t *testing.T) {
	testFail(t, "function a(eval) { \"use strict\"; }",
		"Invalid binding `eval` at (1:11)", nil)
}

func TestFail176(t *testing.T) {
	testFail(t, "function a(package) { \"use strict\"; }",
		"Invalid binding `package` at (1:11)", nil)
}

func TestFail177(t *testing.T) {
	testFail(t, "function a() { \"use strict\"; function b(t, t) { }; }",
		"Parameter name clash at (1:43)", nil)
}

func TestFail178(t *testing.T) {
	testFail(t, "(function a(t, t) { \"use strict\"; })",
		"Parameter name clash at (1:15)", nil)
}

func TestFail179(t *testing.T) {
	testFail(t, "function a() { \"use strict\"; (function b(t, t) { }); }",
		"Parameter name clash at (1:44)", nil)
}

func TestFail180(t *testing.T) {
	testFail(t, "(function a(eval) { \"use strict\"; })",
		"Invalid binding `eval` at (1:12)", nil)
}

func TestFail181(t *testing.T) {
	testFail(t, "(function a(package) { \"use strict\"; })",
		"Invalid binding `package` at (1:12)", nil)
}

func TestFail182(t *testing.T) {
	testFail(t, "\"use strict\";function foo(){\"use strict\";}function bar(){var v = 015}",
		"Octal literals are not allowed in strict mode at (1:65)", nil)
}

func TestFail183(t *testing.T) {
	testFail(t, "var this = 10;", "Unexpected token `this` at (1:4)", nil)
}

func TestFail184(t *testing.T) {
	testFail(t, "throw\n10;", "Illegal newline after throw at (1:0)", nil)
}

func TestFail185(t *testing.T) {
	testFail(t, "const a;",
		"Const declarations require an initialization value at (1:6)", nil)
}

func TestFail186(t *testing.T) {
	testFail(t, "({ get prop(x) {} })",
		"Getter must not have any formal parameters at (1:12)", nil)
}

func TestFail187(t *testing.T) {
	testFail(t, "({ set prop() {} })",
		"Setter must have exactly one formal parameter at (1:11)", nil)
}

func TestFail188(t *testing.T) {
	testFail(t, "({ set prop(x, y) {} })",
		"Setter must have exactly one formal parameter at (1:11)", nil)
}

func TestFail189(t *testing.T) {
	testFail(t, "function(){}", "Unexpected token `(` at (1:8)", nil)
}

func TestFail190(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "07.5", "Unexpected token at (1:2)", opts)
}

func TestFail191(t *testing.T) {
	testFail(t, "\\u{74}rue",
		"Keyword must not contain escaped characters at (1:0)", nil)
}

func TestFail192(t *testing.T) {
	testFail(t, "export { X \\u0061s Y }",
		"Keyword must not contain escaped characters at (1:11)", nil)
}

func TestFail193(t *testing.T) {
	testFail(t, "import X fro\\u006d 'x'",
		"Unexpected token `identifier` at (1:9)", nil)
}

func TestFail194(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "le\\u0074 x = 5", "Unexpected token at (1:9)", opts)
}

func TestFail195(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "(function* () { y\\u0069eld 10 })",
		"Unexpected token at (1:27)", opts)
}

func TestFail196(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_STRICT)
	testFail(t, "(async function() { aw\\u0061it x })",
		"Keyword must not contain escaped characters at (1:20)", opts)
}

func TestFail197(t *testing.T) {
	testFail(t, "(\\u0061sync function() { await x })",
		"Keyword must not contain escaped characters at (1:1)", nil)
}

func TestFail198(t *testing.T) {
	testFail(t, "(\\u0061sync () => { await x })",
		"Keyword must not contain escaped characters at (1:1)", nil)
}

func TestFail199(t *testing.T) {
	testFail(t, "\\u0061sync x => { await x }",
		"Keyword must not contain escaped characters at (1:0)", nil)
}

func TestFail200(t *testing.T) {
	testFail(t, "class X { \\u0061sync function x() { await x } }",
		"Keyword must not contain escaped characters at (1:10)", nil)
}

func TestFail201(t *testing.T) {
	testFail(t, "class X { \\u0061sync x() { await x } }",
		"Keyword must not contain escaped characters at (1:10)", nil)
}

func TestFail202(t *testing.T) {
	testFail(t, "class X { static \\u0061sync x() { await x } }",
		"Keyword must not contain escaped characters at (1:17)", nil)
}

func TestFail203(t *testing.T) {
	testFail(t, "({ ge\\u0074 x() {} })",
		"Keyword must not contain escaped characters at (1:3)", nil)
}

func TestFail204(t *testing.T) {
	testFail(t, "export \\u0061sync function y() { await x }",
		"Keyword must not contain escaped characters at (1:7)", nil)
}

func TestFail205(t *testing.T) {
	testFail(t, "export default \\u0061sync function () { await x }",
		"Keyword must not contain escaped characters at (1:15)", nil)
}

func TestFail206(t *testing.T) {
	testFail(t, "({ \\u0061sync x() { await x } })",
		"Keyword must not contain escaped characters at (1:3)", nil)
}

func TestFail207(t *testing.T) {
	testFail(t, "for (x \\u006ff y) {}", "Unexpected token `identifier` at (1:7)", nil)
}

func TestFail208(t *testing.T) {
	testFail(t, "(x=1)=2", "Assigning to rvalue at (1:1)", nil)
}

func TestFail209(t *testing.T) {
	testFail(t, "let foo; try {} catch (foo) {} let foo;",
		"Identifier `foo` has already been declared at (1:35)", nil)
}

func TestFail210(t *testing.T) {
	testFail(t, "try {} catch (foo) { let foo; }",
		"Identifier `foo` has already been declared at (1:25)", nil)
}

func TestFail211(t *testing.T) {
	testFail(t, "try {} catch ([foo]) { var foo; }",
		"Identifier `foo` has already been declared at (1:27)", nil)
}

func TestFail212(t *testing.T) {
	testFail(t, "try {} catch ([foo, foo]) {}",
		"Identifier `foo` has already been declared at (1:20)", nil)
}

func TestFail213(t *testing.T) {
	testFail(t, "try {} catch ({ a: foo, b: { c: [foo] } }) {}",
		"Identifier `foo` has already been declared at (1:33)", nil)
}

func TestFail214(t *testing.T) {
	testFail(t, "try {} catch (foo) { function foo() {} }",
		"Identifier `foo` has already been declared at (1:30)", nil)
}

func TestFail215(t *testing.T) {
	testPass(t, "try {} catch (foo) {} var foo;", nil)
}

func TestFail216(t *testing.T) {
	testPass(t, "try {} catch (foo) {} let foo;", nil)
}

func TestFail217(t *testing.T) {
	testPass(t, "try {} catch (foo) { function x() { var foo; } }", nil)
}

func TestFail218(t *testing.T) {
	testPass(t, "'use strict'; let foo = function foo() {}", nil)
}

func TestFail219(t *testing.T) {
	testFail(t, "½", "Unexpected character at (1:0)", nil)
}

func TestFail220(t *testing.T) {
	testFail(t, "\"use strict\"\nfoo\n05", "Octal literals are not allowed in strict mode at (3:0)", nil)
}

func TestFail221(t *testing.T) {
	testFail(t, "\"use strict\"\n;(foo)\n05", "Octal literals are not allowed in strict mode at (3:0)", nil)
}

func TestFail222(t *testing.T) {
	testFail(t, "'use strict'\n!blah; 05", "Octal literals are not allowed in strict mode at (2:7)", nil)
}

func TestFail223(t *testing.T) {
	testFail(t, "var x = /[P QR]/\\u0067", "Unexpected token at (1:16)", nil)
}

func TestFail224(t *testing.T) {
	testFail(t, "let a = () => { 'use strict'; delete i; }",
		"Deleting local variable in strict mode at (1:37)", nil)
}

func TestFail225(t *testing.T) {
	testFail(t, "let a = () => { '\\12'; 'use strict'; }",
		"Octal escape sequences are not allowed in strict mode at (1:16)", nil)
}

func TestFail226(t *testing.T) {
	testFail(t, "(function () { 'use strict'; with (i); }())",
		"Strict mode code may not include a with statement at (1:29)", nil)
}

func TestFail227(t *testing.T) {
	testFail(t, "let hello = () => {'use strict'; var eval = 10; }",
		"Unexpected token `eval` at (1:37)", nil)
}

func TestFail228(t *testing.T) {
	testFail(t, "let hello = () => {'use strict'; try { } catch (eval) { } }",
		"Unexpected token `eval` at (1:48)", nil)
}

func TestFail229(t *testing.T) {
	testFail(t, "let a = (t, t) => { \"use strict\"; }",
		"Parameter name clash at (1:12)", nil)
}

func TestFail230(t *testing.T) {
	testFail(t, "let a = ({ t }) => { \"use strict\"; }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:9)", nil)
}

func TestFail231(t *testing.T) {
	testFail(t, "class { a = 1 }", "Class name is required at (1:6)", nil)
}

func TestFail232(t *testing.T) {
	testFail(t, "function a(1) {}", "Unexpected token `number` at (1:11)", nil)
}

func TestFail233(t *testing.T) {
	testFail(t, "function a([ a = { b = 1 } ]) {}",
		"Shorthand property assignments are valid only in destructuring patterns at (1:21)", nil)
}

func TestFail234(t *testing.T) {
	testFail(t, "let a = ([ a = { b = 1 } ]) => {}",
		"Shorthand property assignments are valid only in destructuring patterns at (1:19)", nil)
}

func TestFail235(t *testing.T) {
	testFail(t, "let a = ([ a = { b: { c = 1 } } ]) => {}",
		"Shorthand property assignments are valid only in destructuring patterns at (1:24)", nil)
}

func TestFail236(t *testing.T) {
	testFail(t, "f({x = 0})",
		"Shorthand property assignments are valid only in destructuring patterns at (1:5)", nil)
}

func TestFail237(t *testing.T) {
	testFail(t, "class c { f([ a = { b = 1 } ]) {} }",
		"Shorthand property assignments are valid only in destructuring patterns at (1:22)", nil)
}

func TestFail238(t *testing.T) {
	testFail(t, "({...})", "Unexpected token `}` at (1:5)", nil)
}

func TestFail239(t *testing.T) {
	testFail(t, "let {...obj1,} = foo",
		"Rest element must be last element at (1:12)", nil)
}

func TestFail240(t *testing.T) {
	testFail(t, "let {...obj1,a} = foo",
		"Rest element must be last element at (1:12)", nil)
}

func TestFail241(t *testing.T) {
	testFail(t, "let {...obj1,...obj2} = foo",
		"Rest element must be last element at (1:12)", nil)
}

func TestFail242(t *testing.T) {
	testFail(t, "let {...(obj)} = foo", "Unexpected token `(` at (1:8)", nil)
}

func TestFail243(t *testing.T) {
	testFail(t, "let {...(a,b)} = foo", "Unexpected token `(` at (1:8)", nil)
}

func TestFail244(t *testing.T) {
	testFail(t, "let {...{a,b}} = foo",
		"Binding pattern is not permitted as rest operator's argument at (1:8)", nil)
}

func TestFail245(t *testing.T) {
	testFail(t, "let {...[a,b]} = foo",
		"Binding pattern is not permitted as rest operator's argument at (1:8)", nil)
}

func TestFail246(t *testing.T) {
	testFail(t, "({...obj1,} = foo)",
		"Rest element must be last element at (1:9)", nil)
}

func TestFail247(t *testing.T) {
	testFail(t, "({...obj1,a} = foo)",
		"Rest element must be last element at (1:9)", nil)
}

func TestFail248(t *testing.T) {
	testFail(t, "({...obj1,...obj2} = foo)",
		"Rest element must be last element at (1:9)", nil)
}

func TestFail249(t *testing.T) {
	testFail(t, "({...(a,b)} = foo)", "Assigning to rvalue at (1:6)", nil)
}

func TestFail250(t *testing.T) {
	testFail(t, "({...{a,b}} = foo)", "Invalid rest operator's argument at (1:5)", nil)
}

func TestFail251(t *testing.T) {
	testFail(t, "({...[a,b]} = foo)", "Invalid rest operator's argument at (1:5)", nil)
}

func TestFail252(t *testing.T) {
	testFail(t, "({...(obj)}) => {}", "Invalid parenthesized assignment pattern at (1:5)", nil)
}

func TestFail253(t *testing.T) {
	testFail(t, "({...(obj)}) => {}", "Invalid parenthesized assignment pattern at (1:5)", nil)
}

func TestFail254(t *testing.T) {
	testFail(t, "({...(a,b)}) => {}", "Invalid parenthesized assignment pattern at (1:5)", nil)
}

func TestFail255(t *testing.T) {
	testFail(t, "({...{a,b}}) => {}", "Invalid rest operator's argument at (1:5)", nil)
}

func TestFail256(t *testing.T) {
	testFail(t, "({...[a,b]}) => {}", "Invalid rest operator's argument at (1:5)", nil)
}

func TestFail257(t *testing.T) {
	testFail(t, "({get x() {}}) => {}",
		"Object pattern can't contain getter or setter at (1:2)", nil)
}

func TestFail258(t *testing.T) {
	testFail(t, "let {...x, ...y} = {}",
		"Rest element must be last element at (1:9)", nil)
}

func TestFail259(t *testing.T) {
	testFail(t, "({...x,}) => z",
		"Rest element must be last element at (1:6)", nil)
}

func TestFail260(t *testing.T) {
	testFail(t, "function ({...x,}) { z }", "Unexpected token `(` at (1:9)", nil)
}

func TestFail261(t *testing.T) {
	testFail(t, "let {...{x, y}} = {}",
		"Binding pattern is not permitted as rest operator's argument at (1:8)", nil)
}

func TestFail262(t *testing.T) {
	testFail(t, "let {...{...{x, y}}} = {}",
		"Binding pattern is not permitted as rest operator's argument at (1:8)", nil)
}

func TestFail263(t *testing.T) {
	testFail(t, "0, {...rest, b} = {}",
		"Rest element must be last element at (1:11)", nil)
}

func TestFail264(t *testing.T) {
	testFail(t, "(([a, ...b = 0]) => {})", "Rest elements cannot have a default value at (1:9)", nil)
}

func TestFail265(t *testing.T) {
	testFail(t, "(({a, ...b = 0}) => {})", "Rest elements cannot have a default value at (1:9)", nil)
}

func TestFail266(t *testing.T) {
	testFail(t, "export const { foo, ...bar } = baz;\nexport const bar = 1;\n",
		"Identifier `bar` has already been declared at (2:13)", nil)
}

func TestFail267(t *testing.T) {
	testFail(t, "`\\unicode`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail268(t *testing.T) {
	testFail(t, "`\\u`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail269(t *testing.T) {
	testFail(t, "`\\u{`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail270(t *testing.T) {
	testFail(t, "`\\u{abcdx`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail271(t *testing.T) {
	testFail(t, "`\\u{abcdx}`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail272(t *testing.T) {
	testFail(t, "`\\xylophone`", "Bad character escape sequence at (1:1)", nil)
}

func TestFail275(t *testing.T) {
	testFail(t, "foo`\\unicode", "Unterminated template at (1:3)", nil)
}

func TestFail276(t *testing.T) {
	testFail(t, "foo`\\unicode\\`", "Unterminated template at (1:3)", nil)
}

func TestFail277(t *testing.T) {
	testFail(t, "(...a,) => a",
		"Rest element must be last element at (1:5)", nil)
}

func TestFail278(t *testing.T) {
	testFail(t, "function foo(...a,) { }",
		"Rest element must be last element at (1:17)", nil)
}

func TestFail279(t *testing.T) {
	testFail(t, "(function(...a,) { })",
		"Rest element must be last element at (1:14)", nil)
}

func TestFail280(t *testing.T) {
	testFail(t, "async (...a,) => a",
		"Rest element must be last element at (1:11)", nil)
}

func TestFail281(t *testing.T) {
	testFail(t, "({foo(...a,) {}})", "Rest element must be last element at (1:10)", nil)
}

func TestFail282(t *testing.T) {
	testFail(t, "class A {foo(...a,) {}}",
		"Rest element must be last element at (1:17)", nil)
}

func TestFail283(t *testing.T) {
	testFail(t, "class A {static foo(...a,) {}}",
		"Rest element must be last element at (1:24)", nil)
}

func TestFail284(t *testing.T) {
	testFail(t, "export default function foo(...a,) { }",
		"Rest element must be last element at (1:32)", nil)
}

func TestFail285(t *testing.T) {
	testFail(t, "export default (function foo(...a,) { })",
		"Rest element must be last element at (1:33)", nil)
}

func TestFail286(t *testing.T) {
	testFail(t, "export function foo(...a,) { }",
		"Rest element must be last element at (1:24)", nil)
}

func TestFail287(t *testing.T) {
	testFail(t, "function foo(,) { }", "Unexpected token `,` at (1:13)", nil)
}

func TestFail288(t *testing.T) {
	testFail(t, "(function(,) { })", "Unexpected token `,` at (1:10)", nil)
}

func TestFail289(t *testing.T) {
	testFail(t, "(,) => a", "Unexpected token `,` at (1:1)", nil)
}

func TestFail290(t *testing.T) {
	testFail(t, "async (,) => a", "Unexpected token `,` at (1:7)", nil)
}

func TestFail291(t *testing.T) {
	testFail(t, "({foo(,) {}})", "Unexpected token `,` at (1:6)", nil)
}

func TestFail292(t *testing.T) {
	testFail(t, "class A {foo(,) {}}", "Unexpected token `,` at (1:13)", nil)
}

func TestFail293(t *testing.T) {
	testFail(t, "class A {static foo(,) {}}", "Unexpected token `,` at (1:20)", nil)
}

func TestFail294(t *testing.T) {
	testFail(t, "(class {foo(,) {}})", "Unexpected token `,` at (1:12)", nil)
}

func TestFail295(t *testing.T) {
	testFail(t, "(class {static foo(,) {}})", "Unexpected token `,` at (1:19)", nil)
}

func TestFail296(t *testing.T) {
	testFail(t, "export default function foo(,) { }",
		"Unexpected token `,` at (1:28)", nil)
}

func TestFail297(t *testing.T) {
	testFail(t, "export default (function foo(,) { })",
		"Unexpected token `,` at (1:29)", nil)
}

func TestFail298(t *testing.T) {
	testFail(t, "export function foo(,) { }",
		"Unexpected token `,` at (1:20)", nil)
}

func TestFail299(t *testing.T) {
	testFail(t, "(a,)", "Unexpected trailing comma at (1:2)", nil)
}

func TestFail300(t *testing.T) {
	testFail(t, "({a} &&= b)", "Assigning to rvalue at (1:1)", nil)
}

func TestFail301(t *testing.T) {
	testFail(t, "({a} ||= b)", "Assigning to rvalue at (1:1)", nil)
}

func TestFail302(t *testing.T) {
	testFail(t, "({a} ??= b)", "Assigning to rvalue at (1:1)", nil)
}

func TestFail303(t *testing.T) {
	testFail(t, "/\u2029/", "Unterminated regular expression at (1:1)", nil)
}

func TestFail304(t *testing.T) {
	testFail(t, "/\u2028/", "Unterminated regular expression at (1:1)", nil)
}

func TestFail305(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, "foo`\\unicode`", "Bad character escape sequence at (1:4)", opts)
}

func TestFail306(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, "foo`\\xylophone`", "Bad character escape sequence at (1:4)", opts)
}

// cover some labeled statements
func TestFail307(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, "LabelA: let a = 0", "Unexpected token `let` at (1:8)", opts)
}

func TestFail308(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, `
LabelA: a = 1

for (;;) {
  continue LabelA;
}
`, "Undefined label `LabelA` at (5:11)", opts)
}

func TestFail309(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, `
LabelA: LabelB: for (;;) {
  LabelA: b = 1;
  break LabelA;
}
`, "Label `LabelA` already declared at (3:2)", opts)
}

func TestFail310(t *testing.T) {
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_BAD_ESCAPE_IN_TAGGED_TPL)
	testFail(t, `
LabelA: LabelB: for (;;) {
  LabelB: b = 1;
  break LabelA;
}
`, "Label `LabelB` already declared at (3:2)", opts)
}

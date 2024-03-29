package parser

// import (
// 	. "github.com/hsiaosiyuan0/mole/util"
// )

// func compileTs(code string, opts *ParserOpts) (Node, error) {
// 	if opts == nil {
// 		opts = NewParserOpts()
// 	}
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	p := newParser(code, opts)
// 	return p.Prog()
// }

// func TestTs(t *testing.T) {
// 	ast, err := compileTs("var a: (a: string | number, b: string) => number = () => 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().TsTyp().Type(), "should be ok")
// }

// func TestTs1(t *testing.T) {
// 	ast, err := compileTs("var a: string | number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_UNION_TYP, dec.id.(*Ident).ti.TypAnnot().TsTyp().Type(), "should be ok")
// }

// func TestTs2(t *testing.T) {
// 	ast, err := compileTs("var a: (Array<b> | number) = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_UNION_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs3(t *testing.T) {
// 	ast, err := compileTs("var a: ({ a = c }: { a: string | number }, b: string) => number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsFn := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsFnTyp)
// 	AssertEqual(t, "a", tsFn.params[0].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_LIT_OBJ, tsFn.params[0].(*ObjPat).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs4(t *testing.T) {
// 	// should be failed since `[...a, string|number]` is not a legal formal param
// 	_, err := compileTs("var a: ([string | number], a: string) => number = 1", nil)
// 	AssertEqual(t, "Unexpected token at (1:16)", err.Error(), "should be prog ok")
// }

// func TestTs5(t *testing.T) {
// 	_, err := compileTs("function fn(a: number, b: string) { }", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// }

// func TestTs6(t *testing.T) {
// 	ast, err := compileTs("var a: ({ b: Array<a>| number}) = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_LIT_OBJ, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	prop := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsObj).props[0].(*TsProp)
// 	AssertEqual(t, "b", prop.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_UNION_TYP, prop.val.Type(), "should be ok")
// }

// func TestTs7(t *testing.T) {
// 	_, err := compileTs("var a: ({ b: Array<a> | number, ...c }) = 1", nil)
// 	AssertEqual(t, "Unexpected token at (1:32)", err.Error(), "should be prog ok")
// }

// func TestTs8(t *testing.T) {
// 	ast, err := compileTs("var a: ({ [k: string]: { b: Array<a> | number, c } }) = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_LIT_OBJ, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	p0 := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsObj).props[0]
// 	AssertEqual(t, N_TS_IDX_SIG, p0.Type(), "should be ok")
// }

// func TestTs9(t *testing.T) {
// 	_, err := compileTs("var a: (string) => number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// }

// func TestTs10(t *testing.T) {
// 	_, err := compileTs("var a: (string<a>) => number = 1", nil)
// 	AssertEqual(t, "Unexpected token `<` at (1:14)", err.Error(), "should be prog ok")
// }

// func TestTs11(t *testing.T) {
// 	_, err := compileTs("var a: (string[][]) => number = 1", nil)
// 	AssertEqual(t, "Unexpected token at (1:14)", err.Error(), "should be prog ok")
// }

// func TestTs12(t *testing.T) {
// 	_, err := compileTs("var a: (string<a>|b) => number = 1", nil)
// 	AssertEqual(t, "Unexpected token `<` at (1:14)", err.Error(), "should be prog ok")
// }

// func TestTs13(t *testing.T) {
// 	ast, err := compileTs("var a: ({a}, {b}) => number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsFn := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsFnTyp)
// 	AssertEqual(t, "a", tsFn.params[0].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "b", tsFn.params[1].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
// }

// func TestTs14(t *testing.T) {
// 	ast, err := compileTs("var a: ([a, ...b]: number[], { c }: { c: string }) => number = () => 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsFn := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsFnTyp)
// 	AssertEqual(t, N_PAT_ARRAY, tsFn.params[0].Type(), "should be ok")
// 	AssertEqual(t, N_PAT_OBJ, tsFn.params[1].Type(), "should be ok")
// }

// func TestTs15(t *testing.T) {
// 	ast, err := compileTs("function f(a?: number) {}", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	fn := prog.stmts[0].(*FnDec)
// 	AssertEqual(t, true, fn.params[0].(*Ident).ti.Ques() != nil, "should be ok")
// }

// func TestTs16(t *testing.T) {
// 	ast, err := compileTs("function f(a: {a?: number}) {}", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	fn := prog.stmts[0].(*FnDec)
// 	AssertEqual(t, true, fn.params[0].(*Ident).ti.Ques() == nil, "should be ok")

// 	p0 := fn.params[0].(*Ident).ti.TypAnnot().tsTyp.(*TsObj).props[0].(*TsProp)
// 	AssertEqual(t, true, p0.ques != nil, "should be ok")
// }

// func TestTs17(t *testing.T) {
// 	ast, err := compileTs("var a: (a: {a?: number}) => number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsFn := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsFnTyp)
// 	AssertEqual(t, true, tsFn.params[0].(*Ident).ti.Ques() == nil, "should be ok")

// 	p0 := tsFn.params[0].(*Ident).ti.TypAnnot().tsTyp.(*TsObj).props[0].(*TsProp)
// 	AssertEqual(t, true, p0.ques != nil, "should be ok")
// }

// func TestTs18(t *testing.T) {
// 	ast, err := compileTs("var a: (a: {m?()}) => number = 1", nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_FN_TYP, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsFn := dec.id.(*Ident).ti.TypAnnot().tsTyp.(*TsFnTyp)
// 	AssertEqual(t, true, tsFn.params[0].(*Ident).ti.Ques() == nil, "should be ok")

// 	p0 := tsFn.params[0].(*Ident).ti.TypAnnot().tsTyp.(*TsObj).props[0]
// 	AssertEqual(t, N_TS_PROP, p0.Type(), "should be ok")
// 	AssertEqual(t, true, p0.(*TsProp).ques != nil, "should be ok")
// 	AssertEqual(t, "m", p0.(*TsProp).key.(*Ident).Text(), "should be ok")
// }

// func TestTs19(t *testing.T) {
// 	// PropertyDefinition
// 	ast, err := compileTs(`let a = {
//     m(b: { c: string }) { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	obj := dec.init.(*ObjLit)
// 	prop0 := obj.props[0].(*Prop)
// 	fn := prop0.value.(*FnDec)
// 	param0 := fn.params[0].(*Ident)
// 	typAnnot := param0.ti.TypAnnot().tsTyp.(*TsObj)
// 	AssertEqual(t, "c", typAnnot.props[0].(*TsProp).key.(*Ident).Text(), "should be ok")
// }

// func TestTs20(t *testing.T) {
// 	// AccessibilityModifier
// 	ast, err := compileTs(`class A {
//   constructor(public b: { c: string }) { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	md := dec.Body().(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_CTOR, md.kind, "should be ok")

// 	ti := md.val.(*FnDec).params[0].(*Ident).ti
// 	AssertEqual(t, ACC_MOD_PUB, ti.AccMod(), "should be ok")
// }

// func TestTs21(t *testing.T) {
// 	// AccessibilityModifier
// 	_, err := compileTs(`let a = {
//     m(public b: { c: string }) { }
// }`, nil)

// 	AssertEqual(t,
// 		"A parameter property is only allowed in a constructor implementation at (2:6)", err.Error(),
// 		"should be prog ok")
// }

// func TestTs22(t *testing.T) {
// 	// ArrowFn
// 	ast, err := compileTs(`let a = ({ b }: { b?: string }, c: Array<string> & number) => { }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	fn := dec.init.(*ArrowFn)
// 	param0 := fn.params[0].(*ObjPat)
// 	ti := param0.ti
// 	AssertEqual(t, N_TS_LIT_OBJ, ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	param1 := fn.params[1].(*Ident)
// 	ti = param1.ti
// 	AssertEqual(t, N_TS_INTERSECT_TYP, ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs23(t *testing.T) {
// 	// ReturnType
// 	ast, err := compileTs(`let a = ({ b }: { b?: string }, c: Array<string> & number): void => { }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	fn := dec.init.(*ArrowFn)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs24(t *testing.T) {
// 	// ReturnType
// 	ast, err := compileTs(`function fn(): void { }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	fn := prog.stmts[0].(*FnDec)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs25(t *testing.T) {
// 	// ReturnType
// 	ast, err := compileTs(`function fn(): { b?: string } { }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	fn := prog.stmts[0].(*FnDec)
// 	AssertEqual(t, N_TS_LIT_OBJ, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs26(t *testing.T) {
// 	// ReturnType
// 	ast, err := compileTs(`let a = {
//     m(): void { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	obj := dec.init.(*ObjLit)
// 	prop0 := obj.props[0].(*Prop)
// 	fn := prop0.value.(*FnDec)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs27(t *testing.T) {
// 	// ReturnType
// 	ast, err := compileTs(`class A {
//     m(): void { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, N_TS_VOID, m.val.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs28(t *testing.T) {
// 	// ReturnType & getter
// 	ast, err := compileTs(`class A {
//     get m(): string { return "" }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_GETTER, m.kind, "should be ok")
// 	AssertEqual(t, N_TS_STR, m.val.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs29(t *testing.T) {
// 	// Setter
// 	ast, err := compileTs(`class A {
//     set m(n: string) { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_SETTER, m.kind, "should be ok")

// 	param0 := m.val.(*FnDec).params[0].(*Ident)
// 	AssertEqual(t, N_TS_STR, param0.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs30(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`f<string | number, void>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
// 	AssertEqual(t, 2, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")

// 	AssertEqual(t, N_TS_VOID, c.ti.TypArgs().(*TsParamsInst).params[1].Type(), "should be ok")
// }

// func TestTs31(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`f<string>()<number>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
// 	AssertEqual(t, 1, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_NUM, c.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")

// 	c = prog.stmts[0].(*ExprStmt).expr.(*CallExpr).callee.(*CallExpr)
// 	AssertEqual(t, 1, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_STR, c.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")
// }

// func TestTs32(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`f<string>().f<number>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
// 	AssertEqual(t, 1, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_NUM, c.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")

// 	m := prog.stmts[0].(*ExprStmt).expr.(*CallExpr).callee.(*MemberExpr)
// 	c = m.obj.(*CallExpr)
// 	AssertEqual(t, 1, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_STR, c.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")
// }

// func TestTs33(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`new f<string>()<number>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
// 	AssertEqual(t, 1, len(c.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_NUM, c.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")

// 	n := c.callee.(*NewExpr)
// 	AssertEqual(t, 1, len(n.ti.TypArgs().(*TsParamsInst).params), "should be ok")
// 	AssertEqual(t, N_TS_STR, n.ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")
// }

// func TestTs34(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`class A {
//     m<T, R>(): void { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_METHOD, m.kind, "should be ok")
// 	AssertEqual(t, N_TS_VOID, m.val.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tp := m.val.(*FnDec).ti.TypParams()
// 	AssertEqual(t, 2, len(tp.(*TsParamsDec).params), "should be ok")
// 	AssertEqual(t, "R", tp.(*TsParamsDec).params[1].(*TsParam).name.(*Ident).Text(), "should be ok")
// }

// func TestTs35(t *testing.T) {
// 	// arguments
// 	ast, err := compileTs(`let a = {
//     m<T>(): void { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	obj := dec.init.(*ObjLit)
// 	prop0 := obj.props[0].(*Prop)
// 	fn := prop0.value.(*FnDec)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// 	AssertEqual(t, "T", fn.ti.TypParams().(*TsParamsDec).params[0].(*TsParam).name.(*Ident).Text(), "should be ok")
// }

// func TestTs36(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	_, err := compileTs(`let f = <T>(a: T) => {}`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")
// }

// func TestTs37(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = { m: <T, R>(a: T): void => { a++ } }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	obj := dec.init.(*ObjLit)
// 	prop0 := obj.props[0].(*Prop)
// 	fn := prop0.value.(*ArrowFn)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// 	AssertEqual(t, 2, len(fn.ti.TypParams().(*TsParamsDec).params), "should be ok")
// 	AssertEqual(t, "R", fn.ti.TypParams().(*TsParamsDec).params[1].(*TsParam).name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, 1, len(fn.body.(*BlockStmt).body), "should be ok")
// }

// func TestTs38(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`class A {
//     m<T>() { }
// }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_METHOD, m.kind, "should be ok")
// 	typAnnot := m.val.(*FnDec).ti.TypAnnot()
// 	AssertEqual(t, nil, typAnnot, "should be ok")

// 	tp := m.val.(*FnDec).ti.TypParams()
// 	AssertEqual(t, 1, len(tp.(*TsParamsDec).params), "should be ok")
// 	AssertEqual(t, "T", tp.(*TsParamsDec).params[0].(*TsParam).name.(*Ident).Text(), "should be ok")
// }

// func TestTs39(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`class A {
//     set a<T>(a: T) { }
// }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_SETTER, m.kind, "should be ok")
// 	AssertEqual(t, "a", m.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, nil, m.val.(*FnDec).ti.TypAnnot(), "should be ok")

// 	tp := m.val.(*FnDec).ti.TypParams()
// 	AssertEqual(t, 1, len(tp.(*TsParamsDec).params), "should be ok")
// 	AssertEqual(t, "T", tp.(*TsParamsDec).params[0].(*TsParam).name.(*Ident).Text(), "should be ok")
// }

// func TestTs40(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number>b`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")
// 	AssertEqual(t, "b", ta.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs41(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number>b++`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs42(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number>++b`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, true, up.prefix, "should be ok")
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs43(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number><string>b++`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_STR, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs44(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number><string><boolean>b++`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_STR, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_BOOL, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs45(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = <number><string><boolean>++b`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	ta := dec.init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_STR, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_BOOL, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, true, up.prefix, "should be ok")
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs46(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = 1 + <number><string>b++`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	bin := dec.init.(*BinExpr)
// 	ta := bin.rhs.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NUM, ta.des.Type(), "should be ok")

// 	ta = ta.arg.(*TsTypAssert)
// 	AssertEqual(t, N_TS_STR, ta.des.Type(), "should be ok")

// 	up := ta.arg.(*UpdateExpr)
// 	AssertEqual(t, "b", up.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs47(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)
// 	ast, err := compileTs(`let a = { m: <T, R extends string>(a: T): void => { a++ } }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	obj := dec.init.(*ObjLit)
// 	prop0 := obj.props[0].(*Prop)
// 	fn := prop0.value.(*ArrowFn)
// 	AssertEqual(t, N_TS_VOID, fn.ti.TypAnnot().tsTyp.Type(), "should be ok")
// 	AssertEqual(t, 2, len(fn.ti.TypParams().(*TsParamsDec).params), "should be ok")
// 	AssertEqual(t, "R", fn.ti.TypParams().(*TsParamsDec).params[1].(*TsParam).name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "string", fn.ti.TypParams().(*TsParamsDec).params[1].(*TsParam).cons.(*TsPredef).Text(), "should be ok")
// 	AssertEqual(t, 1, len(fn.body.(*BlockStmt).body), "should be ok")
// }

// func TestTs48(t *testing.T) {
// 	// arguments
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	_, err := compileTs(`class A {
//     constructor<T>(): T { }
// }`, opts)

// 	AssertEqual(t,
// 		"Type parameters cannot appear on a constructor declaration at (2:15)",
// 		err.Error(), "should be prog ok")
// }

// func TestTs49(t *testing.T) {
// 	// TypeAliasDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`type a = string | number`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsTypDec)
// 	AssertEqual(t, "a", dec.name.(*Ident).Text(), "should be ok")

// 	AssertEqual(t, N_TS_UNION_TYP, dec.ti.TypAnnot().tsTyp.Type(), "should be ok")
// 	AssertEqual(t, "string", dec.ti.TypAnnot().TsTyp().(*TsUnionTyp).elems[0].(*TsPredef).Text(), "should be ok")
// }

// func TestTs50(t *testing.T) {
// 	// SimpleVariableDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`let a: number = 1`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, "a", dec.id.(*Ident).Text(), "should be ok")

// 	AssertEqual(t, N_TS_NUM, dec.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs51(t *testing.T) {
// 	// DestructuringLexicalBinding
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`let { a }: { a: number } = { a: 1 }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)

// 	prop := dec.id.(*ObjPat).props[0].(*Prop)
// 	AssertEqual(t, "a", prop.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_LIT_OBJ, dec.id.(*ObjPat).ti.TypAnnot().tsTyp.Type(), "should be ok")

// 	tsProp := dec.id.(*ObjPat).ti.TypAnnot().tsTyp.(*TsObj).props[0].(*TsProp)
// 	AssertEqual(t, "a", tsProp.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, tsProp.val.(*TsTypAnnot).tsTyp.Type(), "should be ok")
// }

// func TestTs52(t *testing.T) {
// 	// FunctionDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`function f()
// function f(): any {}`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	sig := prog.stmts[0].(*FnDec)
// 	fn := prog.stmts[1].(*FnDec)
// 	AssertEqual(t, true, sig.IsSig(), "should be ok")
// 	AssertEqual(t, true, sig.id.(*Ident).Text() == fn.id.(*Ident).Text(), "should be ok")
// }

// func TestTs53(t *testing.T) {
// 	// FunctionDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	_, err := compileTs(`function f()
// function f1()
// function f(): any {}`, opts)

// 	AssertEqual(t,
// 		"Function implementation is missing or not immediately following the declaration at (1:0)",
// 		err.Error(), "should be prog ok")
// }

// func TestTs54(t *testing.T) {
// 	// FunctionDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	_, err := compileTs(`function f()
// function f1(): any {}`, opts)

// 	AssertEqual(t,
// 		"Function implementation name must be `f` at (2:9)",
// 		err.Error(), "should be prog ok")
// }

// func TestTs55(t *testing.T) {
// 	// InterfaceDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`interface A<T> extends C<R>, D<S> { b }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	itf := prog.stmts[0].(*TsInterface)
// 	AssertEqual(t, 2, len(itf.supers), "should be ok")
// 	AssertEqual(t, "b", itf.body.(*TsInterfaceBody).body[0].(*TsProp).key.(*Ident).Text(), "should be ok")
// }

// func TestTs56(t *testing.T) {
// 	// EnumDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`const enum A { m1, m2 = "m2" }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	enum := prog.stmts[0].(*TsEnum)
// 	AssertEqual(t, true, enum.cons, "should be ok")
// 	AssertEqual(t, 2, len(enum.items), "should be ok")
// }

// func TestTs57(t *testing.T) {
// 	// EnumDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`enum A { m1, m2 = "a" + "b" }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	enum := prog.stmts[0].(*TsEnum)
// 	AssertEqual(t, false, enum.cons, "should be ok")
// 	AssertEqual(t, 2, len(enum.items), "should be ok")
// 	AssertEqual(t, N_EXPR_BIN, enum.items[1].(*TsEnumMember).val.Type(), "should be ok")
// }

// func TestTs58(t *testing.T) {
// 	// ImportAliasDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`import a = b.c`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ia := prog.stmts[0].(*TsImportAlias)
// 	AssertEqual(t, "a", ia.name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NS_NAME, ia.val.Type(), "should be ok")
// 	AssertEqual(t, "c", ia.val.(*TsNsName).rhs.(*Ident).Text(), "should be ok")
// }

// func TestTs59(t *testing.T) {
// 	// ImportAliasDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export import a = b.c`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ia := prog.stmts[0].(*TsImportAlias)
// 	AssertEqual(t, true, ia.export, "should be ok")
// 	AssertEqual(t, "a", ia.name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NS_NAME, ia.val.Type(), "should be ok")
// 	AssertEqual(t, "c", ia.val.(*TsNsName).rhs.(*Ident).Text(), "should be ok")
// }

// func TestTs60(t *testing.T) {
// 	// NamespaceDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`namespace b { export const c = 1}`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ns := prog.stmts[0].(*TsNS)
// 	AssertEqual(t, 1, len(ns.body.(*BlockStmt).body), "should be ok")
// }

// func TestTs61(t *testing.T) {
// 	// NamespaceDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export namespace b { export const c = 1}`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*ExportDec)
// 	ns := ep.dec.(*TsNS)
// 	AssertEqual(t, 1, len(ns.body.(*BlockStmt).body), "should be ok")
// }

// func TestTs62(t *testing.T) {
// 	// export TypeAliasDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export type a = string | number;`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*ExportDec)
// 	dec := ep.dec.(*TsTypDec)
// 	AssertEqual(t, 1, len(prog.stmts), "should be ok")
// 	AssertEqual(t, "a", dec.name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_UNION_TYP, dec.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs63(t *testing.T) {
// 	// export InterfaceDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export interface A<T> extends C<R>, D<S> { b }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*ExportDec)
// 	itf := ep.dec.(*TsInterface)
// 	AssertEqual(t, 2, len(itf.supers), "should be ok")
// 	AssertEqual(t, "b", itf.body.(*TsInterfaceBody).body[0].(*TsProp).key.(*Ident).Text(), "should be ok")
// }

// func TestTs64(t *testing.T) {
// 	// ImportRequireDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`import a = require('test');`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	n := prog.stmts[0].(*TsImportRequire)
// 	AssertEqual(t, "test", n.expr.(*CallExpr).args[0].(*StrLit).Text(), "should be ok")
// }

// func TestTs65(t *testing.T) {
// 	// export ImportRequireDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export import a = require('test');`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*ExportDec)
// 	n := ep.dec.(*TsImportRequire)
// 	AssertEqual(t, 1, len(prog.stmts), "should be ok")
// 	AssertEqual(t, "test", n.expr.(*CallExpr).args[0].(*StrLit).Text(), "should be ok")
// }

// func TestTs66(t *testing.T) {
// 	// ExportAssignment
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export = a`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*TsExportAssign)
// 	AssertEqual(t, "a", ep.expr.(*Ident).Text(), "should be ok")
// }

// func TestTs67(t *testing.T) {
// 	// AmbientDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare enum Enum { A = 1, B, C = 2, }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, 3, len(dec.inner.(*TsEnum).items), "should be ok")
// }

// func TestTs68(t *testing.T) {
// 	// AmbientDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`export declare enum Enum { A = 1, B, C = 2, }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	ep := prog.stmts[0].(*ExportDec)
// 	dec := ep.dec.(*TsDec)
// 	AssertEqual(t, "Enum", dec.inner.(*TsEnum).name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, 3, len(dec.inner.(*TsEnum).items), "should be ok")
// }

// func TestTs69(t *testing.T) {
// 	// AmbientVariableDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare let a`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	td := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, N_TS_DEC_VAR_DEC, td.Type(), "should be ok")
// 	AssertEqual(t, "a", td.inner.(*VarDecStmt).decList[0].(*VarDec).id.(*Ident).Text(), "should be ok")
// }

// func TestTs70(t *testing.T) {
// 	// AmbientFunctionDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare function a();`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	td := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, N_TS_DEC_FN, td.Type(), "should be ok")
// 	AssertEqual(t, "a", td.inner.(*FnDec).id.(*Ident).Text(), "should be ok")
// }

// func TestTs71(t *testing.T) {
// 	// AmbientFunctionDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare function a(): number;`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	td := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, N_TS_DEC_FN, td.Type(), "should be ok")
// 	AssertEqual(t, "a", td.inner.(*FnDec).id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, td.inner.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs72(t *testing.T) {
// 	// AmbientTypeAliasDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare type a = number;`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	td := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, N_TS_DEC_TYP_DEC, td.Type(), "should be ok")
// 	AssertEqual(t, "a", td.inner.(*TsTypDec).name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, td.inner.(*TsTypDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs73(t *testing.T) {
// 	// AmbientNamespaceDeclaration
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.Off(FEAT_JSX)

// 	ast, err := compileTs(`declare namespace a { type a = number; }`, opts)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	td := prog.stmts[0].(*TsDec)
// 	AssertEqual(t, N_TS_DEC_NS, td.Type(), "should be ok")

// 	dec := td.inner.(*TsNS).body.(*BlockStmt).body[0].(*TsTypDec)
// 	AssertEqual(t, "a", dec.name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, dec.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs74(t *testing.T) {
// 	// Class TypeParams
// 	ast, err := compileTs(`class A<T> {
//     m(): void { }
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", dec.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "T", dec.ti.TypParams().(*TsParamsDec).params[0].(*TsParam).name.(*Ident).Text(), "should be ok")

// 	m := dec.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, N_TS_VOID, m.val.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs75(t *testing.T) {
// 	// AmbientClassDeclaration ConstructorSignature
// 	ast, err := compileTs(`declare class a {
//     constructor()
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	cls := dec.inner.(*ClassDec)
// 	AssertEqual(t, "a", cls.id.(*Ident).Text(), "should be ok")

// 	ctor := cls.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_CTOR, ctor.kind, "should be ok")
// }

// func TestTs76(t *testing.T) {
// 	// AmbientClassDeclaration MethodSignature
// 	ast, err := compileTs(`declare class a {
//     c(): any;
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	cls := dec.inner.(*ClassDec)
// 	AssertEqual(t, "a", cls.id.(*Ident).Text(), "should be ok")

// 	m := cls.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_METHOD, m.kind, "should be ok")
// 	AssertEqual(t, N_TS_ANY, m.val.(*FnDec).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs77(t *testing.T) {
// 	// AmbientClassDeclaration IndexSignature
// 	ast, err := compileTs(`declare class a {
//     [k: string]: number
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	cls := dec.inner.(*ClassDec)
// 	AssertEqual(t, "a", cls.id.(*Ident).Text(), "should be ok")

// 	idx := cls.body.(*ClassBody).elems[0].(*Field)
// 	AssertEqual(t, "k", idx.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, idx.ti.typAnnot.tsTyp.Type(), "should be ok")
// }

// func TestTs78(t *testing.T) {
// 	// Class AccessibilityModifier
// 	ast, err := compileTs(`declare class a {
//     public static b: number
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	cls := dec.inner.(*ClassDec)
// 	AssertEqual(t, "a", cls.id.(*Ident).Text(), "should be ok")

// 	f := cls.body.(*ClassBody).elems[0].(*Field)
// 	AssertEqual(t, true, f.static, "should be ok")
// 	AssertEqual(t, ACC_MOD_PUB, f.ti.AccMod(), "should be ok")
// 	AssertEqual(t, "b", f.key.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, f.ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs79(t *testing.T) {
// 	// Class AccessibilityModifier
// 	ast, err := compileTs(`class a {
//     public ['test']() {}
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	cls := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "a", cls.id.(*Ident).Text(), "should be ok")

// 	m := cls.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, false, m.static, "should be ok")
// 	AssertEqual(t, ACC_MOD_PUB, m.ti.AccMod(), "should be ok")
// 	AssertEqual(t, "test", m.key.(*StrLit).Text(), "should be ok")
// }

// func TestTs80(t *testing.T) {
// 	// Class AccessibilityModifier
// 	ast, err := compileTs(`class A {
//     private constructor() {}
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	cls := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, "A", cls.id.(*Ident).Text(), "should be ok")

// 	m := cls.body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, PK_CTOR, m.kind, "should be ok")
// 	AssertEqual(t, ACC_MOD_PRI, m.ti.AccMod(), "should be ok")
// }

// func TestTs81(t *testing.T) {
// 	// AmbientModuleDeclaration
// 	ast, err := compileTs(`declare module 'a' {
//     let a: number;
// }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*TsDec)
// 	blk := dec.inner.(*BlockStmt)
// 	AssertEqual(t, N_TS_DEC_MODULE, dec.Type(), "should be ok")
// 	AssertEqual(t, "a", dec.name.(*StrLit).Text(), "should be ok")

// 	vds := blk.body[0].(*VarDecStmt)
// 	vd := vds.decList[0].(*VarDec)
// 	AssertEqual(t, "a", vd.id.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_NUM, vd.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs82(t *testing.T) {
// 	// function with typParams
// 	ast, err := compileTs(`function a<T>() {}`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	dec := prog.stmts[0].(*FnDec)
// 	AssertEqual(t, "T", dec.ti.TypParams().(*TsParamsDec).params[0].(*TsParam).name.(*Ident).Text(), "should be ok")
// }

// func TestTs83(t *testing.T) {
// 	ast, err := compileTs(`4 + async<number>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, N_TS_NUM, expr.rhs.(*CallExpr).ti.TypArgs().(*TsParamsInst).params[0].Type(), "should be ok")
// }

// func TestTs84(t *testing.T) {
// 	ast, err := compileTs(`type a = ({ a }?: { a: string }) => void`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*TsTypDec)
// 	fn := expr.ti.TypAnnot().TsTyp().(*TsFnTyp)
// 	AssertEqual(t, true, fn.params[0].(*ObjPat).ti.Ques() != nil, "should be ok")
// }

// func TestTs85(t *testing.T) {
// 	// predicate types
// 	ast, _ := compileTs(`(x: any): x is string => true;`, nil)
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	ast, _ = compileTs(`(x: any): asserts x is string => true;`, nil)
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// }

// func TestTs86(t *testing.T) {
// 	// predicate types
// 	ast, err := compileTs(`x < y`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, T_LT, expr.op, "should be ok")
// 	AssertEqual(t, "y", expr.rhs.(*Ident).Text(), "should be ok")
// }

// func TestTs87(t *testing.T) {
// 	// predicate types
// 	ast, err := compileTs(`x < y < z<a>()`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, T_LT, expr.op, "should be ok")
// 	AssertEqual(t, "z", expr.rhs.(*CallExpr).callee.(*Ident).Text(), "should be ok")
// }

// func TestTs88(t *testing.T) {
// 	// predicate types
// 	ast, err := compileTs(`x < z<a>() > y`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, T_GT, expr.op, "should be ok")
// 	AssertEqual(t, "z", expr.lhs.(*BinExpr).rhs.(*CallExpr).callee.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "y", expr.rhs.(*Ident).Text(), "should be ok")
// }

// func TestTs89(t *testing.T) {
// 	// cast
// 	ast, err := compileTs(`x as any as T;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")

// 	ast, err = compileTs(`x as boolean <= y;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")

// 	ast, err = compileTs(`x === 1 as number;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "1", expr.rhs.(*BinExpr).lhs.(*NumLit).Text(), "should be ok")

// 	ast, err = compileTs(`x as boolean ?? y;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "y", expr.rhs.(*Ident).Text(), "should be ok")

// 	ast, err = compileTs(`x < 1 as A<string>`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "1", expr.lhs.(*BinExpr).rhs.(*NumLit).Text(), "should be ok")
// }

// func TestTs90(t *testing.T) {
// 	// cast to TypRef
// 	ast, err := compileTs(`x < 1 as A<string>`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "1", expr.lhs.(*BinExpr).rhs.(*NumLit).Text(), "should be ok")

// 	ast, err = compileTs(`x < b as A<string>`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*ExprStmt).expr.(*BinExpr)
// 	AssertEqual(t, "x", expr.lhs.(*BinExpr).lhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "b", expr.lhs.(*BinExpr).rhs.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, "A", expr.rhs.(*TsRef).name.(*Ident).Text(), "should be ok")
// 	AssertEqual(t, N_TS_STR, expr.rhs.(*TsRef).args.(*TsParamsInst).params[0].Type(), "should be ok")
// }

// func TestTs91(t *testing.T) {
// 	// assert const
// 	ast, err := compileTs(`let v1 = 'abc' as const;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*BinExpr)
// 	AssertEqual(t, "const", expr.rhs.(*TsRef).name.(*Ident).Text(), "should be ok")
// }

// func TestTs92(t *testing.T) {
// 	// assert const
// 	ast, err := compileTs(`let q1 = <const> 10;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*TsTypAssert)
// 	AssertEqual(t, "const", expr.des.(*TsRef).name.(*Ident).Text(), "should be ok")
// }

// func TestTs93(t *testing.T) {
// 	// ts literal
// 	ast, err := compileTs(`let q1 = <1> 10;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*TsTypAssert)
// 	AssertEqual(t, "1", expr.des.(*TsLit).lit.(*NumLit).Text(), "should be ok")

// 	ast, err = compileTs(`let q1 = <true> 10;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*TsTypAssert)
// 	AssertEqual(t, "true", expr.des.(*TsLit).lit.(*BoolLit).Text(), "should be ok")

// 	ast, err = compileTs(`let q1 = <null> 10;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	expr = prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init.(*TsTypAssert)
// 	AssertEqual(t, N_TS_NULL, expr.des.(*TsPredef).typ, "should be ok")
// }

// func TestTs94(t *testing.T) {
// 	// non-null
// 	ast, err := compileTs(`x!;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*TsNoNull)
// 	AssertEqual(t, "x", expr.arg.(*Ident).Text(), "should be ok")
// }

// func TestTs95(t *testing.T) {
// 	// non-null
// 	ast, err := compileTs(`x!.y;`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr.(*MemberExpr)
// 	AssertEqual(t, "x", expr.obj.(*TsNoNull).arg.(*Ident).Text(), "should be ok")
// }

// func TestTs96(t *testing.T) {
// 	// ts union type
// 	ast, err := compileTs(`let a:  a | b & c | d`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
// 	AssertEqual(t, N_TS_UNION_TYP, expr.id.(*Ident).ti.TypAnnot().tsTyp.Type(), "should be ok")
// }

// func TestTs97(t *testing.T) {
// 	// abstract class
// 	ast, err := compileTs(`abstract class C1 {}`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ClassDec)
// 	AssertEqual(t, true, expr.Abstract(), "should be ok")

// 	ast, err = compileTs(`abstract class C1 {
//     abstract f();
//   }`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog = ast.(*Prog)
// 	prop := prog.stmts[0].(*ClassDec).body.(*ClassBody).elems[0].(*Method)
// 	AssertEqual(t, true, prop.ti.Abstract(), "should be ok")
// }

// func TestTs98(t *testing.T) {
// 	// abstract class
// 	ast, err := compileTs(`export abstract class C3 {}`, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")
// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExportDec).dec.(*ClassDec)
// 	AssertEqual(t, true, expr.Abstract(), "should be ok")
// }

// func TestTs99(t *testing.T) {
// 	// jsx
// 	opts := NewParserOpts()
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	p := newParser(`let a = <T>a</T>`, opts)

// 	ast, err := p.Prog()
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	jsx := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init
// 	AssertEqual(t, N_JSX_ELEM, jsx.Type(), "should be ok")
// }

// func TestTs100(t *testing.T) {
// 	// jsx
// 	opts := NewParserOpts()
// 	// jsx is turned on within default options
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	p := newParser(`<T,>() => 1;`, opts)

// 	ast, err := p.Prog()
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr
// 	AssertEqual(t, N_EXPR_ARROW, expr.Type(), "should be ok")
// }

// func TestTs101(t *testing.T) {
// 	// jsx
// 	opts := NewParserOpts()
// 	// jsx is turned on within default options
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	p := newParser(`<T extends string>() => 1;`, opts)

// 	ast, err := p.Prog()
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*ExprStmt).expr
// 	AssertEqual(t, N_EXPR_ARROW, expr.Type(), "should be ok")
// }

// func TestTs102(t *testing.T) {
// 	// jsx
// 	opts := NewParserOpts()
// 	// jsx feature is turned on within the default options
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	p := newParser(`<T>() => 1;`, opts)

// 	_, err := p.Prog()
// 	AssertEqual(t, true, err != nil, "should not pass")
// }

// func TestTs103(t *testing.T) {
// 	// jsx
// 	opts := NewParserOpts()
// 	// jsx feature is turned on within the default options
// 	opts.Feature = opts.Feature.On(FEAT_TS)
// 	p := newParser(`let a = <T>t</T>`, opts)

// 	ast, err := p.Prog()
// 	AssertEqual(t, nil, err, "should be prog ok")
// 	AssertEqual(t, true, ast != nil, "should be prog ok")

// 	prog := ast.(*Prog)
// 	expr := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec).init
// 	AssertEqual(t, N_JSX_ELEM, expr.Type(), "should be ok")
// }

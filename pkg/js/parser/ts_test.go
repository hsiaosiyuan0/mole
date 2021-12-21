package parser

import (
	"fmt"
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func compileTs(code string, opts *ParserOpts) (Node, error) {
	if opts == nil {
		opts = NewParserOpts()
	}
	opts.Feature = opts.Feature.On(FEAT_TS)
	p := newParser(code, opts)
	return p.Prog()
}

func TestTs(t *testing.T) {
	ast, err := compileTs("var a: (a: string | number, b: string) => number = () => 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")
}

func TestTs1(t *testing.T) {
	ast, err := compileTs("var a: string | number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_UNION_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")
}

func TestTs2(t *testing.T) {
	ast, err := compileTs("var a: (Array<b> | number) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_UNION_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")
}

func TestTs3(t *testing.T) {
	ast, err := compileTs("var a: ({ a = c }: { a: string | number }, b: string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	tsFn := dec.id.(*Ident).ti.typAnnot.(*TsFnTyp)
	assert.Equal(t, "a", tsFn.params[0].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_OBJ, tsFn.params[0].(*ObjPat).ti.typAnnot.Type(), "should be ok")
}

func TestTs4(t *testing.T) {
	// should be failed since `[...a, string|number]` is not a legal formal param
	_, err := compileTs("var a: ([string | number], a: string) => number = 1", nil)
	assert.Equal(t, "Unexpected token at (1:16)", err.Error(), "should be prog ok")
}

func TestTs5(t *testing.T) {
	_, err := compileTs("function fn(a: number, b: string) { }", nil)
	assert.Equal(t, nil, err, "should be prog ok")
}

func TestTs6(t *testing.T) {
	ast, err := compileTs("var a: ({ b: Array<a>| number}) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_OBJ, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	prop := dec.id.(*Ident).ti.typAnnot.(*TsObj).props[0].(*TsProp)
	assert.Equal(t, "b", prop.key.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_UNION_TYP, prop.val.Type(), "should be ok")
}

func TestTs7(t *testing.T) {
	_, err := compileTs("var a: ({ b: Array<a> | number, ...c }) = 1", nil)
	assert.Equal(t, "Unexpected token at (1:32)", err.Error(), "should be prog ok")
}

func TestTs8(t *testing.T) {
	ast, err := compileTs("var a: ({ [k: string]: { b: Array<a> | number, c } }) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_OBJ, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	p0 := dec.id.(*Ident).ti.typAnnot.(*TsObj).props[0]
	assert.Equal(t, N_TS_IDX_SIG, p0.Type(), "should be ok")
}

func TestTs9(t *testing.T) {
	_, err := compileTs("var a: (string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
}

func TestTs10(t *testing.T) {
	_, err := compileTs("var a: (string<a>) => number = 1", nil)
	assert.Equal(t, "Unexpected token `<` at (1:14)", err.Error(), "should be prog ok")
}

func TestTs11(t *testing.T) {
	_, err := compileTs("var a: (string[][]) => number = 1", nil)
	assert.Equal(t, "Unexpected token at (1:14)", err.Error(), "should be prog ok")
}

func TestTs12(t *testing.T) {
	_, err := compileTs("var a: (string<a>|b) => number = 1", nil)
	assert.Equal(t, "Unexpected token `<` at (1:14)", err.Error(), "should be prog ok")
}

func TestTs13(t *testing.T) {
	ast, err := compileTs("var a: ({a}, {b}) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	tsFn := dec.id.(*Ident).ti.typAnnot.(*TsFnTyp)
	assert.Equal(t, "a", tsFn.params[0].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
	assert.Equal(t, "b", tsFn.params[1].(*ObjPat).props[0].(*Prop).key.(*Ident).Text(), "should be ok")
}

func TestTs14(t *testing.T) {
	ast, err := compileTs("var a: ([a, ...b]: number[], { c }: { c: string }) => number = () => 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	tsFn := dec.id.(*Ident).ti.typAnnot.(*TsFnTyp)
	assert.Equal(t, N_PAT_ARRAY, tsFn.params[0].Type(), "should be ok")
	assert.Equal(t, N_PAT_OBJ, tsFn.params[1].Type(), "should be ok")
}

func TestTs15(t *testing.T) {
	ast, err := compileTs("function f(a?: number) {}", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	fn := prog.stmts[0].(*FnDec)
	assert.Equal(t, true, fn.params[0].(*Ident).ti.ques != nil, "should be ok")
}

func TestTs16(t *testing.T) {
	ast, err := compileTs("function f(a: {a?: number}) {}", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	fn := prog.stmts[0].(*FnDec)
	assert.Equal(t, true, fn.params[0].(*Ident).ti.ques == nil, "should be ok")

	p0 := fn.params[0].(*Ident).ti.typAnnot.(*TsObj).props[0].(*TsProp)
	assert.Equal(t, true, p0.ques != nil, "should be ok")
}

func TestTs17(t *testing.T) {
	ast, err := compileTs("var a: (a: {a?: number}) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	tsFn := dec.id.(*Ident).ti.typAnnot.(*TsFnTyp)
	assert.Equal(t, true, tsFn.params[0].(*Ident).ti.ques == nil, "should be ok")

	p0 := tsFn.params[0].(*Ident).ti.typAnnot.(*TsObj).props[0].(*TsProp)
	assert.Equal(t, true, p0.ques != nil, "should be ok")
}

func TestTs18(t *testing.T) {
	ast, err := compileTs("var a: (a: {m?()}) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_FN_TYP, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")

	tsFn := dec.id.(*Ident).ti.typAnnot.(*TsFnTyp)
	assert.Equal(t, true, tsFn.params[0].(*Ident).ti.ques == nil, "should be ok")

	p0 := tsFn.params[0].(*Ident).ti.typAnnot.(*TsObj).props[0]
	assert.Equal(t, N_TS_PROP, p0.Type(), "should be ok")
	assert.Equal(t, true, p0.(*TsProp).ques != nil, "should be ok")
	assert.Equal(t, "m", p0.(*TsProp).key.(*Ident).Text(), "should be ok")
}

func TestTs19(t *testing.T) {
	// PropertyDefinition
	ast, err := compileTs(`let a = {
    m(b: { c: string }) { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	obj := dec.init.(*ObjLit)
	prop0 := obj.props[0].(*Prop)
	fn := prop0.value.(*FnDec)
	param0 := fn.params[0].(*Ident)
	typAnnot := param0.ti.typAnnot.(*TsObj)
	assert.Equal(t, "c", typAnnot.props[0].(*TsProp).key.(*Ident).Text(), "should be ok")
}

func TestTs20(t *testing.T) {
	// AccessibilityModifier
	ast, err := compileTs(`class A {
  constructor(public b: { c: string }) { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	md := dec.Body().(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_CTOR, md.kind, "should be ok")

	ti := md.value.(*FnDec).params[0].(*Ident).ti
	assert.Equal(t, ACC_MOD_PUB, ti.accMod, "should be ok")
}

func TestTs21(t *testing.T) {
	// AccessibilityModifier
	_, err := compileTs(`let a = {
    m(public b: { c: string }) { }
}`, nil)

	assert.Equal(t,
		"A parameter property is only allowed in a constructor implementation at (2:6)", err.Error(),
		"should be prog ok")
}

func TestTs22(t *testing.T) {
	// ArrowFn
	ast, err := compileTs(`let a = ({ b }: { b?: string }, c: Array<string> & number) => { }`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	fn := dec.init.(*ArrowFn)
	param0 := fn.params[0].(*ObjPat)
	ti := param0.ti
	assert.Equal(t, N_TS_OBJ, ti.typAnnot.Type(), "should be ok")

	param1 := fn.params[1].(*Ident)
	ti = param1.ti
	assert.Equal(t, N_TS_INTERSEC_TYP, ti.typAnnot.Type(), "should be ok")
}

func TestTs23(t *testing.T) {
	// ReturnType
	ast, err := compileTs(`let a = ({ b }: { b?: string }, c: Array<string> & number): void => { }`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	fn := dec.init.(*ArrowFn)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
}

func TestTs24(t *testing.T) {
	// ReturnType
	ast, err := compileTs(`function fn(): void { }`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	fn := prog.stmts[0].(*FnDec)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
}

func TestTs25(t *testing.T) {
	// ReturnType
	ast, err := compileTs(`function fn(): { b?: string } { }`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	fn := prog.stmts[0].(*FnDec)
	assert.Equal(t, N_TS_OBJ, fn.ti.typAnnot.Type(), "should be ok")
}

func TestTs26(t *testing.T) {
	// ReturnType
	ast, err := compileTs(`let a = {
    m(): void { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	obj := dec.init.(*ObjLit)
	prop0 := obj.props[0].(*Prop)
	fn := prop0.value.(*FnDec)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
}

func TestTs27(t *testing.T) {
	// ReturnType
	ast, err := compileTs(`class A {
    m(): void { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, N_TS_VOID, m.value.(*FnDec).ti.typAnnot.Type(), "should be ok")
}

func TestTs28(t *testing.T) {
	// ReturnType & getter
	ast, err := compileTs(`class A {
    get m(): string { return "" }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_GETTER, m.kind, "should be ok")
	assert.Equal(t, N_TS_STR, m.value.(*FnDec).ti.typAnnot.Type(), "should be ok")
}

func TestTs29(t *testing.T) {
	// Setter
	ast, err := compileTs(`class A {
    set m(n: string) { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_SETTER, m.kind, "should be ok")

	param0 := m.value.(*FnDec).params[0].(*Ident)
	assert.Equal(t, N_TS_STR, param0.ti.typAnnot.Type(), "should be ok")
}

func TestTs30(t *testing.T) {
	// arguments
	ast, err := compileTs(`f<string | number, void>()`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
	assert.Equal(t, 2, len(c.ti.typArgs), "should be ok")

	assert.Equal(t, N_TS_VOID, c.ti.typArgs[1].Type(), "should be ok")
}

func TestTs31(t *testing.T) {
	// arguments
	ast, err := compileTs(`f<string>()<number>()`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
	assert.Equal(t, 1, len(c.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_NUM, c.ti.typArgs[0].Type(), "should be ok")

	c = prog.stmts[0].(*ExprStmt).expr.(*CallExpr).callee.(*CallExpr)
	assert.Equal(t, 1, len(c.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_STR, c.ti.typArgs[0].Type(), "should be ok")
}

func TestTs32(t *testing.T) {
	// arguments
	ast, err := compileTs(`f<string>().f<number>()`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
	assert.Equal(t, 1, len(c.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_NUM, c.ti.typArgs[0].Type(), "should be ok")

	m := prog.stmts[0].(*ExprStmt).expr.(*CallExpr).callee.(*MemberExpr)
	c = m.obj.(*CallExpr)
	assert.Equal(t, 1, len(c.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_STR, c.ti.typArgs[0].Type(), "should be ok")
}

func TestTs33(t *testing.T) {
	// arguments
	ast, err := compileTs(`new f<string>()<number>()`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	c := prog.stmts[0].(*ExprStmt).expr.(*CallExpr)
	assert.Equal(t, 1, len(c.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_NUM, c.ti.typArgs[0].Type(), "should be ok")

	n := c.callee.(*NewExpr)
	assert.Equal(t, 1, len(n.ti.typArgs), "should be ok")
	assert.Equal(t, N_TS_STR, n.ti.typArgs[0].Type(), "should be ok")
}

func TestTs34(t *testing.T) {
	// arguments
	ast, err := compileTs(`class A {
    m<T, R>(): void { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_METHOD, m.kind, "should be ok")
	assert.Equal(t, N_TS_VOID, m.value.(*FnDec).ti.typAnnot.Type(), "should be ok")

	tp := m.value.(*FnDec).ti.typParams
	assert.Equal(t, 2, len(tp), "should be ok")
	assert.Equal(t, "R", tp[1].(*TsParam).name.(*Ident).Text(), "should be ok")
}

func TestTs35(t *testing.T) {
	// arguments
	ast, err := compileTs(`let a = {
    m<T>(): void { }
}`, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	obj := dec.init.(*ObjLit)
	prop0 := obj.props[0].(*Prop)
	fn := prop0.value.(*FnDec)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
	assert.Equal(t, "T", fn.ti.typParams[0].(*TsParam).name.(*Ident).Text(), "should be ok")
}

func TestTs36(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let f = <T>(a: T) => {}`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	fmt.Println(ast)
}

func TestTs37(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = { m: <T, R>(a: T): void => { a++ } }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	obj := dec.init.(*ObjLit)
	prop0 := obj.props[0].(*Prop)
	fn := prop0.value.(*ArrowFn)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
	assert.Equal(t, 2, len(fn.ti.typParams), "should be ok")
	assert.Equal(t, "R", fn.ti.typParams[1].(*TsParam).name.(*Ident).Text(), "should be ok")
	assert.Equal(t, 1, len(fn.body.(*BlockStmt).body), "should be ok")
}

func TestTs38(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`class A {
    m<T>() { }
}`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_METHOD, m.kind, "should be ok")
	assert.Equal(t, nil, m.value.(*FnDec).ti.typAnnot, "should be ok")

	tp := m.value.(*FnDec).ti.typParams
	assert.Equal(t, 1, len(tp), "should be ok")
	assert.Equal(t, "T", tp[0].(*TsParam).name.(*Ident).Text(), "should be ok")
}

func TestTs39(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`class A {
    set a<T>(a: T) { }
}`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*ClassDec)
	assert.Equal(t, "A", dec.id.(*Ident).Text(), "should be ok")

	m := dec.body.(*ClassBody).elems[0].(*Method)
	assert.Equal(t, PK_SETTER, m.kind, "should be ok")
	assert.Equal(t, "a", m.key.(*Ident).Text(), "should be ok")
	assert.Equal(t, nil, m.value.(*FnDec).ti.typAnnot, "should be ok")

	tp := m.value.(*FnDec).ti.typParams
	assert.Equal(t, 1, len(tp), "should be ok")
	assert.Equal(t, "T", tp[0].(*TsParam).name.(*Ident).Text(), "should be ok")
}

func TestTs40(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number>b`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")
	assert.Equal(t, "b", ta.arg.(*Ident).Text(), "should be ok")
}

func TestTs41(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number>b++`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs42(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number>++b`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, true, up.prefix, "should be ok")
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs43(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number><string>b++`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_STR, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs44(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number><string><boolean>b++`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_STR, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_BOOL, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs45(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = <number><string><boolean>++b`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	ta := dec.init.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_STR, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_BOOL, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, true, up.prefix, "should be ok")
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs46(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = 1 + <number><string>b++`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	bin := dec.init.(*BinExpr)
	ta := bin.rhs.(*TsTypAssert)
	assert.Equal(t, N_TS_NUM, ta.des.Type(), "should be ok")

	ta = ta.arg.(*TsTypAssert)
	assert.Equal(t, N_TS_STR, ta.des.Type(), "should be ok")

	up := ta.arg.(*UpdateExpr)
	assert.Equal(t, "b", up.arg.(*Ident).Text(), "should be ok")
}

func TestTs47(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)
	ast, err := compileTs(`let a = { m: <T, R extends string>(a: T): void => { a++ } }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	obj := dec.init.(*ObjLit)
	prop0 := obj.props[0].(*Prop)
	fn := prop0.value.(*ArrowFn)
	assert.Equal(t, N_TS_VOID, fn.ti.typAnnot.Type(), "should be ok")
	assert.Equal(t, 2, len(fn.ti.typParams), "should be ok")
	assert.Equal(t, "R", fn.ti.typParams[1].(*TsParam).name.(*Ident).Text(), "should be ok")
	assert.Equal(t, "string", fn.ti.typParams[1].(*TsParam).cons.(*TsPredef).Text(), "should be ok")
	assert.Equal(t, 1, len(fn.body.(*BlockStmt).body), "should be ok")
}

func TestTs48(t *testing.T) {
	// arguments
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	_, err := compileTs(`class A {
    constructor<T>(): T { }
}`, opts)

	assert.Equal(t,
		"Type parameters cannot appear on a constructor declaration at (2:15)",
		err.Error(), "should be prog ok")
}

func TestTs49(t *testing.T) {
	// TypeAliasDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`type a = string | number`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*TsTypDec)
	assert.Equal(t, "a", dec.name.(*Ident).Text(), "should be ok")

	assert.Equal(t, N_TS_UNION_TYP, dec.ti.typAnnot.Type(), "should be ok")
	assert.Equal(t, "string", dec.ti.typAnnot.(*TsUnionTyp).lhs.(*TsPredef).Text(), "should be ok")
}

func TestTs50(t *testing.T) {
	// SimpleVariableDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`let a: number = 1`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)
	assert.Equal(t, "a", dec.id.(*Ident).Text(), "should be ok")

	assert.Equal(t, N_TS_NUM, dec.id.(*Ident).ti.typAnnot.Type(), "should be ok")
}

func TestTs51(t *testing.T) {
	// DestructuringLexicalBinding
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`let { a }: { a: number } = { a: 1 }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	dec := prog.stmts[0].(*VarDecStmt).decList[0].(*VarDec)

	prop := dec.id.(*ObjPat).props[0].(*Prop)
	assert.Equal(t, "a", prop.key.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_OBJ, dec.id.(*ObjPat).ti.typAnnot.Type(), "should be ok")

	tsProp := dec.id.(*ObjPat).ti.typAnnot.(*TsObj).props[0].(*TsProp)
	assert.Equal(t, "a", tsProp.key.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_NUM, tsProp.val.Type(), "should be ok")
}

func TestTs52(t *testing.T) {
	// FunctionDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`function f()
function f(): any {}`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	sig := prog.stmts[0].(*FnDec)
	fn := prog.stmts[1].(*FnDec)
	assert.Equal(t, true, sig.IsSig(), "should be ok")
	assert.Equal(t, true, sig.id.(*Ident).Text() == fn.id.(*Ident).Text(), "should be ok")
}

func TestTs53(t *testing.T) {
	// FunctionDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	_, err := compileTs(`function f()
function f1()
function f(): any {}`, opts)

	assert.Equal(t,
		"Function implementation is missing or not immediately following the declaration at (1:0)",
		err.Error(), "should be prog ok")
}

func TestTs54(t *testing.T) {
	// FunctionDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	_, err := compileTs(`function f()
function f1(): any {}`, opts)

	assert.Equal(t,
		"Function implementation name must be `f` at (2:9)",
		err.Error(), "should be prog ok")
}

func TestTs55(t *testing.T) {
	// InterfaceDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`interface A<T> extends C<R>, D<S> { b }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	itf := prog.stmts[0].(*TsInferface)
	assert.Equal(t, 2, len(itf.supers), "should be ok")
	assert.Equal(t, "b", itf.body.(*TsObj).props[0].(*TsProp).key.(*Ident).Text(), "should be ok")
}

func TestTs56(t *testing.T) {
	// EnumDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`const enum A { m1, m2 = "m2" }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	enum := prog.stmts[0].(*TsEnum)
	assert.Equal(t, true, enum.cons, "should be ok")
	assert.Equal(t, 2, len(enum.mems), "should be ok")
}

func TestTs57(t *testing.T) {
	// EnumDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`enum A { m1, m2 = "a" + "b" }`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	enum := prog.stmts[0].(*TsEnum)
	assert.Equal(t, false, enum.cons, "should be ok")
	assert.Equal(t, 2, len(enum.mems), "should be ok")
	assert.Equal(t, N_EXPR_BIN, enum.mems[1].(*TsEnumMember).val.Type(), "should be ok")
}

func TestTs58(t *testing.T) {
	// ImportAliasDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`import a = b.c`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	ia := prog.stmts[0].(*TsImportAlias)
	assert.Equal(t, "a", ia.name.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_NS_NAME, ia.val.Type(), "should be ok")
	assert.Equal(t, "c", ia.val.(*TsNsName).rhs.(*Ident).Text(), "should be ok")
}

func TestTs59(t *testing.T) {
	// ImportAliasDeclaration
	opts := NewParserOpts()
	opts.Feature = opts.Feature.Off(FEAT_JSX)

	ast, err := compileTs(`export import a = b.c`, opts)
	assert.Equal(t, nil, err, "should be prog ok")

	prog := ast.(*Prog)
	ia := prog.stmts[0].(*TsImportAlias)
	assert.Equal(t, true, ia.export, "should be ok")
	assert.Equal(t, "a", ia.name.(*Ident).Text(), "should be ok")
	assert.Equal(t, N_TS_NS_NAME, ia.val.Type(), "should be ok")
	assert.Equal(t, "c", ia.val.(*TsNsName).rhs.(*Ident).Text(), "should be ok")
}

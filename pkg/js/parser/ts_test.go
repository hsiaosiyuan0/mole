package parser

import (
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
	assert.Equal(t, "Unexpected token at (1:14)", err.Error(), "should be prog ok")
}

func TestTs11(t *testing.T) {
	_, err := compileTs("var a: (string[][]) => number = 1", nil)
	assert.Equal(t, "Unexpected token at (1:14)", err.Error(), "should be prog ok")
}

func TestTs12(t *testing.T) {
	_, err := compileTs("var a: (string<a>|b) => number = 1", nil)
	assert.Equal(t, "Unexpected token at (1:14)", err.Error(), "should be prog ok")
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

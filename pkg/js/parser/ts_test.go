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
	ast, err := compileTs("var a: (a:string|number,a:string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs1(t *testing.T) {
	ast, err := compileTs("var a: string|number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs2(t *testing.T) {
	ast, err := compileTs("var a: (string<a>|number) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs3(t *testing.T) {
	ast, err := compileTs("var a: ({a = c}:string|number,a:string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs4(t *testing.T) {
	// should be failed since `[...a, string|number]` is not a legal formal param
	ast, err := compileTs("var a: ([...a, string|number],a:string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs5(t *testing.T) {
	ast, err := compileTs("function fn (a: number,b:string) {}", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs6(t *testing.T) {
	ast, err := compileTs("var a: ({b: string<a>|number}) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs7(t *testing.T) {
	ast, err := compileTs("var a: ({b: string<a>|number, ...c}) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs8(t *testing.T) {
	ast, err := compileTs("var a: ({[k: string]: {b: string<a>|number, c}}) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs9(t *testing.T) {
	ast, err := compileTs("var a: ({[k: string]: {b: string<a>|number, c}}) = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs10(t *testing.T) {
	ast, err := compileTs("var a: (string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs11(t *testing.T) {
	ast, err := compileTs("var a: (string<a>) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs12(t *testing.T) {
	ast, err := compileTs("var a: (string[][]) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs13(t *testing.T) {
	ast, err := compileTs("var a: (string<a>|b) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs14(t *testing.T) {
	ast, err := compileTs("var a: ({a}, {b}) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

func TestTs15(t *testing.T) {
	ast, err := compileTs("var a: ([a, ...b]: number[], { c }: { c: string }) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

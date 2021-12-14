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
	ast, err := compileTs("var a: ([...a, string|number],a:string) => number = 1", nil)
	assert.Equal(t, nil, err, "should be prog ok")
	fmt.Println(ast)
}

package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestJSXFail1(t *testing.T) {
	TestFail(t, "var x = <div>one</div><div>two</div>;",
		"Adjacent JSX elements must be wrapped in an enclosing tag at (1:22)", nil)
}

func TestJSXFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_JSX_NS)
	TestFail(t, "<a:b.c />", "Unexpected token at (1:4)", opts)
}

func TestJSXFail3(t *testing.T) {
	TestFail(t, "<ns:div />", "Unexpected token `:` at (1:3)", nil)
}

func TestJSXFail4(t *testing.T) {
	TestFail(t, "<div ns:attr />", "Unexpected token `:` at (1:7)", nil)
}

func TestJSXFail5(t *testing.T) {
	TestFail(t, "<A>foo{</A>", "Unexpected token at (1:8)", nil)
}

func TestJSXFail6(t *testing.T) {
	TestFail(t, "<A>foo<</A>", "Unexpected token `<` at (1:7)", nil)
}

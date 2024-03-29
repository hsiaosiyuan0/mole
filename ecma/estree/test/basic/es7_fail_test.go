package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestEs7thFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_POW)
	TestFail(t, "x **= 42", "Unexpected token at (1:2)", opts)
}

func TestEs7thFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_POW)
	TestFail(t, "x ** y", "Unexpected token at (1:2)", opts)
}

func TestEs7thFail3(t *testing.T) {
	TestFail(t, "delete o.p ** 2;",
		"Unary operator `delete` used immediately before exponentiation expression at (1:7)", nil)
}

func TestEs7thFail4(t *testing.T) {
	TestFail(t, "void 2 ** 2;",
		"Unary operator `void` used immediately before exponentiation expression at (1:5)", nil)
}

func TestEs7thFail5(t *testing.T) {
	TestFail(t, "typeof 2 ** 2;",
		"Unary operator `typeof` used immediately before exponentiation expression at (1:7)", nil)
}

func TestEs7thFail6(t *testing.T) {
	TestFail(t, "~3 ** 2;",
		"Unary operator `~` used immediately before exponentiation expression at (1:1)", nil)
}

func TestEs7thFail7(t *testing.T) {
	TestFail(t, "!1 ** 2;",
		"Unary operator `!` used immediately before exponentiation expression at (1:1)", nil)
}

func TestEs7thFail8(t *testing.T) {
	TestFail(t, "-2** 2;",
		"Unary operator `-` used immediately before exponentiation expression at (1:1)", nil)
}

func TestEs7thFail9(t *testing.T) {
	TestFail(t, "+2** 2;",
		"Unary operator `+` used immediately before exponentiation expression at (1:1)", nil)
}

func TestEs7thFail10(t *testing.T) {
	TestFail(t, "-(i--) ** 2",
		"Unary operator `-` used immediately before exponentiation expression at (1:2)", nil)
}

func TestEs7thFail11(t *testing.T) {
	TestFail(t, "+(i--) ** 2",
		"Unary operator `+` used immediately before exponentiation expression at (1:2)", nil)
}

func TestEs7thFail12(t *testing.T) {
	TestFail(t, "x %* y", "Unexpected token `*` at (1:3)", nil)
}

func TestEs7thFail13(t *testing.T) {
	TestFail(t, "x %*= y", "Unexpected token `*=` at (1:3)", nil)
}

func TestEs7thFail14(t *testing.T) {
	TestFail(t, "function foo(a=2) { 'use strict'; }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:13)", nil)
}

func TestEs7thFail15(t *testing.T) {
	TestFail(t, "(a=2) => { 'use strict'; }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:1)", nil)
}

func TestEs7thFail16(t *testing.T) {
	TestFail(t, "function foo({a}) { 'use strict'; }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:13)", nil)
}

func TestEs7thFail17(t *testing.T) {
	TestFail(t, "({a}) => { 'use strict'; }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:1)", nil)
}

func TestEs7thFail18(t *testing.T) {
	TestFail(t,
		"'use strict'; if(x) function f() {}",
		"function declarations can't appear in single-statement context at (1:20)", nil)
}

func TestEs7thFail19(t *testing.T) {
	TestFail(t, "'use strict'; function y(x = 1) { 'use strict' }",
		"Illegal 'use strict' directive in function with non-simple parameter list at (1:25)", nil)
}

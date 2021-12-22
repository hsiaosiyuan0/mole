package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := parser.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compileProg(code string, opts *parser.ParserOpts) (parser.Node, error) {
	p := newParser(code, opts)
	return p.Prog()
}

func testFail(t *testing.T, code, errMs string, opts *parser.ParserOpts) {
	ast, err := compileProg(code, opts)
	if err == nil {
		t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
	}
	assert.Equal(t, errMs, err.Error(), "")
}

func TestHarmonyFail1(t *testing.T) {
	testFail(t, "([a.a]) => 42", "Assigning to rvalue at (1:2)", nil)
}

func TestHarmonyFail2(t *testing.T) {
	testFail(t, "() => {}()", "Unexpected token at (1:8)", nil)
}

func TestHarmonyFail3(t *testing.T) {
	testFail(t, "(a) => {}()", "Unexpected token at (1:9)", nil)
}

func TestHarmonyFail4(t *testing.T) {
	testFail(t, "a => {}()", "Unexpected token at (1:7)", nil)
}

func TestHarmonyFail5(t *testing.T) {
	testFail(t, "console.log(typeof () => {});", "Malformed arrow function parameter list at (1:19)", nil)
}

func TestHarmonyFail6(t *testing.T) {
	testFail(t, "x = { method() { super(); } }", "super() call outside constructor of a subclass at (1:17)", nil)
}

func TestHarmonyFail7(t *testing.T) {
	testFail(t, "export var await", "Invalid binding `await` at (1:11)", nil)
}

func TestHarmonyFail8(t *testing.T) {
	testFail(t, "export new Foo();", "Unexpected token `new` at (1:7)", nil)
}

func TestHarmonyFail9(t *testing.T) {
	testFail(t, "export typeof foo;", "Unexpected token `typeof` at (1:7)", nil)
}

func TestHarmonyFail10(t *testing.T) {
	testFail(t, "export *", "Unexpected token `EOF` at (1:8)", nil)
}

func TestHarmonyFail11(t *testing.T) {
	testFail(t, "export { encrypt }", "Export `encrypt` is not defined at (1:9)", nil)
}

func TestHarmonyFail12(t *testing.T) {
	testFail(t, "class Test {}; export default class Test {}", "Identifier `Test` has already been declared at (1:36)", nil)
}

func TestHarmonyFail13(t *testing.T) {
	testFail(t, "export { encrypt, encrypt }", "Duplicate export `encrypt` at (1:18)", nil)
}

func TestHarmonyFail14(t *testing.T) {
	testFail(t, "var encrypt; export { encrypt }; export { encrypt }", "Duplicate export `encrypt` at (1:42)", nil)
}

func TestHarmonyFail15(t *testing.T) {
	testFail(t, "export { decrypt as encrypt }; function encrypt() {}", "Export `decrypt` is not defined at (1:9)", nil)
}

func TestHarmonyFail16(t *testing.T) {
	testFail(t, "export { encrypt }; if (true) function encrypt() {}",
		"function declarations can't appear in single-statement context at (1:30)", nil)
}

func TestHarmonyFail17(t *testing.T) {
	testFail(t, "{ function encrypt() {} } export { encrypt }", "Export `encrypt` is not defined at (1:35)", nil)
}

func TestHarmonyFail18(t *testing.T) {
	testFail(t, "export { default }", "Unexpected token `default` at (1:9)", nil)
}

func TestHarmonyFail19(t *testing.T) {
	testFail(t, "export { if }", "Unexpected token `if` at (1:9)", nil)
}

func TestHarmonyFail20(t *testing.T) {
	testFail(t, "export { default as foo }", "Unexpected token `default` at (1:9)", nil)
}

func TestHarmonyFail21(t *testing.T) {
	testFail(t, "export { if as foo }", "Unexpected token `if` at (1:9)", nil)
}

func TestHarmonyFail22(t *testing.T) {
	testFail(t, "import default from \"foo\"", "Unexpected token `default` at (1:7)", nil)
}

func TestHarmonyFail23(t *testing.T) {
	testFail(t, "import { class } from 'foo'", "Unexpected token `class` at (1:9)", nil)
}

func TestHarmonyFail24(t *testing.T) {
	testFail(t, "import { class, var } from 'foo'", "Unexpected token `class` at (1:9)", nil)
}

func TestHarmonyFail25(t *testing.T) {
	testFail(t, "import { a as class } from 'foo'", "Unexpected token `class` at (1:14)", nil)
}

func TestHarmonyFail26(t *testing.T) {
	testFail(t, "import * as class from 'foo'", "Unexpected token `class` at (1:12)", nil)
}

func TestHarmonyFail27(t *testing.T) {
	testFail(t, "import { enum } from 'foo'", "Unexpected token `enum` at (1:9)", nil)
}

func TestHarmonyFail28(t *testing.T) {
	testFail(t, "import { a as enum } from 'foo'", "Unexpected token `enum` at (1:14)", nil)
}

func TestHarmonyFail29(t *testing.T) {
	testFail(t, "import * as enum from 'foo'", "Unexpected token `enum` at (1:12)", nil)
}

func TestHarmonyFail30(t *testing.T) {
	testFail(t, "() => { class a extends b { static get prototype(){} } }",
		"Classes can't have a static field named `prototype` at (1:39)", nil)
}

func TestHarmonyFail31(t *testing.T) {
	testFail(t, "class a extends b { static set prototype(a){} }",
		"Classes can't have a static field named `prototype` at (1:31)", nil)
}

func TestHarmonyFail32(t *testing.T) {
	testFail(t, "class a { static prototype(a){} }",
		"Classes can't have a static field named `prototype` at (1:17)", nil)
}

func TestHarmonyFail33(t *testing.T) {
	testFail(t, "function *g() { (x = yield) => {} }", "Yield expression cannot be a default value at (1:21)", nil)
}

func TestHarmonyFail34(t *testing.T) {
	testFail(t, "function *g() { ({x = yield}) => {} }", "Yield expression cannot be a default value at (1:22)", nil)
}

func TestHarmonyFail35(t *testing.T) {
	testFail(t, "(class { *static x() {} })", "Unexpected token `identifier` at (1:17)", nil)
}

func TestHarmonyFail36(t *testing.T) {
	testFail(t, "(class A {constructor() { super() }})", "super() call outside constructor of a subclass at (1:26)", nil)
}

func TestHarmonyFail37(t *testing.T) {
	testFail(t, "(class A extends B { constructor() { function f() { super() } } })", "super() call outside constructor of a subclass at (1:52)", nil)
}

func TestHarmonyFail38(t *testing.T) {
	testFail(t, "(class A extends B { method() { super() } })", "super() call outside constructor of a subclass at (1:32)", nil)
}

func TestHarmonyFail39(t *testing.T) {
	testFail(t, "class A { constructor() {} 'constructor'() {} }",
		"Duplicate constructor in the same class at (1:27)", nil)
}

func TestHarmonyFail40(t *testing.T) {
	testFail(t, "class A { get constructor() {} }", "Constructor can't have get/set modifier at (1:14)", nil)
}

func TestHarmonyFail41(t *testing.T) {
	testFail(t, "class A { *constructor() {} }", "Constructor can't be a generator at (1:11)", nil)
}

func TestHarmonyFail42(t *testing.T) {
	testFail(t, "\"use strict\"; (class A extends B { static constructor() { super() }})",
		"super() call outside constructor of a subclass at (1:58)", nil)
}

func TestHarmonyFail43(t *testing.T) {
	testFail(t, "({[x]})", "A computed property name must have property initialization at (1:5)", nil)
}

func TestHarmonyFail44(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "function x(...[ a, b ]){}", "Binding pattern is not permitted as rest operator's argument at (1:14)", opts)
}

func TestHarmonyFail45(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(([...[ a, b ]]) => {})", "Binding pattern is not permitted as rest operator's argument at (1:6)", opts)
}

func TestHarmonyFail46(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "function x({ a: { w, x }, b: [y, z] }, ...[a, b, c]){}",
		"Binding pattern is not permitted as rest operator's argument at (1:42)", opts)
}

func TestHarmonyFail47(t *testing.T) {
	testFail(t, "(function ({ a(){} }) {})", "Invalid destructuring assignment target at (1:13)", nil)
}

func TestHarmonyFail48(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(function x(...[ a, b ]){})", "Binding pattern is not permitted as rest operator's argument at (1:15)", opts)
}

func TestHarmonyFail49(t *testing.T) {
	testFail(t, "var a = { set foo(...v) {} };",
		"Setter cannot use rest params at (1:18)", nil)
}

func TestHarmonyFail50(t *testing.T) {
	testFail(t, "class a { set foo(...v) {} };",
		"Setter cannot use rest params at (1:18)", nil)
}

func TestHarmonyFail51(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(function x({ a: { w, x }, b: [y, z] }, ...[a, b, c]){})",
		"Binding pattern is not permitted as rest operator's argument at (1:43)", opts)
}

func TestHarmonyFail52(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(...[a, b]) => {}", "Binding pattern is not permitted as rest operator's argument at (1:4)", opts)
}

func TestHarmonyFail53(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(a, ...[b]) => {}", "Binding pattern is not permitted as rest operator's argument at (1:7)", opts)
}

// Harmony Invalid syntax

func TestHarmonyFail54(t *testing.T) {
	testFail(t, "0o", "Expected number in radix 8 at (1:2)", nil)
}

func TestHarmonyFail55(t *testing.T) {
	testFail(t, "0o1a", "Identifier directly after number at (1:3)", nil)
}

func TestHarmonyFail56(t *testing.T) {
	testFail(t, "0o9", "Expected number in radix 8 at (1:2)", nil)
}

func TestHarmonyFail57(t *testing.T) {
	testFail(t, "0O18", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail58(t *testing.T) {
	testFail(t, "0b", "Expected number in radix 2 at (1:2)", nil)
}

func TestHarmonyFail59(t *testing.T) {
	testFail(t, "0b1a", "Identifier directly after number at (1:3)", nil)
}

func TestHarmonyFail60(t *testing.T) {
	testFail(t, "0b9", "Expected number in radix 2 at (1:2)", nil)
}

func TestHarmonyFail61(t *testing.T) {
	testFail(t, "0b18", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail62(t *testing.T) {
	testFail(t, "0b12", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail63(t *testing.T) {
	testFail(t, "0B", "Expected number in radix 2 at (1:2)", nil)
}

func TestHarmonyFail64(t *testing.T) {
	testFail(t, "0B1a", "Identifier directly after number at (1:3)", nil)
}

func TestHarmonyFail65(t *testing.T) {
	testFail(t, "0B9", "Expected number in radix 2 at (1:2)", nil)
}

func TestHarmonyFail66(t *testing.T) {
	testFail(t, "0B18", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail67(t *testing.T) {
	testFail(t, "0B12", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail68(t *testing.T) {
	testFail(t, "\"\\u{110000}\"", "Code point out of bounds at (1:2)", nil)
}

func TestHarmonyFail69(t *testing.T) {
	testFail(t, "\"\\u{}\"", "Bad character escape sequence at (1:2)", nil)
}

func TestHarmonyFail70(t *testing.T) {
	testFail(t, "\"\\u{FFFF\"", "Bad character escape sequence at (1:2)", nil)
}

func TestHarmonyFail71(t *testing.T) {
	testFail(t, "\"\\u{FFZ}\"", "Bad character escape sequence at (1:2)", nil)
}

func TestHarmonyFail72(t *testing.T) {
	testFail(t, "[v] += ary", "Assigning to rvalue at (1:0)", nil)
}

func TestHarmonyFail73(t *testing.T) {
	testFail(t, "[2] = 42", "Assigning to rvalue at (1:1)", nil)
}

func TestHarmonyFail74(t *testing.T) {
	testFail(t, "({ obj:20 }) = 42", "Invalid parenthesized assignment pattern at (1:0)", nil)
}

func TestHarmonyFail75(t *testing.T) {
	testFail(t, "( { get x() {} } = 0)",
		"Object pattern can't contain getter or setter at (1:4)", nil)
}

func TestHarmonyFail76(t *testing.T) {
	testFail(t, "x \n is y", "Unexpected token at (2:4)", nil)
}

func TestHarmonyFail77(t *testing.T) {
	testFail(t, "x \n isnt y", "Unexpected token at (2:6)", nil)
}

func TestHarmonyFail78(t *testing.T) {
	testFail(t, "function default() {}", "Unexpected token `default` at (1:9)", nil)
}

func TestHarmonyFail79(t *testing.T) {
	testFail(t, "function hello() {'use strict'; ({ i: 10, s(eval) { } }); }",
		"Invalid binding `eval` at (1:44)", nil)
}

func TestHarmonyFail80(t *testing.T) {
	testFail(t, "function a() { \"use strict\"; ({ b(t, t) { } }); }", "Parameter name clash at (1:37)", nil)
}

func TestHarmonyFail81(t *testing.T) {
	testFail(t, "var super", "Unexpected token `super` at (1:4)", nil)
}

func TestHarmonyFail82(t *testing.T) {
	testFail(t, "var default", "Unexpected token `default` at (1:4)", nil)
}

func TestHarmonyFail83(t *testing.T) {
	testFail(t, "let default", "Unexpected token `default` at (1:4)", nil)
}

func TestHarmonyFail84(t *testing.T) {
	testFail(t, "const default", "Unexpected token `default` at (1:6)", nil)
}

func TestHarmonyFail85(t *testing.T) {
	testFail(t, "\"use strict\"; ({ v: eval } = obj)",
		"Assigning to `eval` in strict mode at (1:20)", nil)
}

func TestHarmonyFail86(t *testing.T) {
	testFail(t, "\"use strict\"; ({ v: arguments } = obj)",
		"Assigning to `arguments` in strict mode at (1:20)", nil)
}

func TestHarmonyFail87(t *testing.T) {
	testFail(t, "for (let x = 42 in list) process(x);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail88(t *testing.T) {
	testFail(t, "for (const x = 42 in list) process(x);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail89(t *testing.T) {
	testFail(t, "for (let x = 42 of list) process(x);",
		"for-of loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail90(t *testing.T) {
	testFail(t, "for (const x = 42 of list) process(x);",
		"for-of loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail91(t *testing.T) {
	testFail(t, "for (var x = 42 of list) process(x);",
		"for-of loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail92(t *testing.T) {
	testFail(t, "for (var {x} = 42 of list) process(x);",
		"for-of loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail93(t *testing.T) {
	testFail(t, "for (var [x] = 42 of list) process(x);",
		"for-of loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail94(t *testing.T) {
	testFail(t, "var x; for (x = 42 of list) process(x);", "Assigning to rvalue at (1:12)", nil)
}

func TestHarmonyFail95(t *testing.T) {
	testFail(t, "import foo", "Unexpected token `EOF` at (1:10)", nil)
}

func TestHarmonyFail96(t *testing.T) {
	testFail(t, "import { foo, bar }", "Unexpected token `EOF` at (1:19)", nil)
}

func TestHarmonyFail97(t *testing.T) {
	testFail(t, "import foo from bar", "Unexpected token `identifier` at (1:16)", nil)
}

func TestHarmonyFail98(t *testing.T) {
	testFail(t, "((a)) => 42", "Invalid parenthesized assignment pattern at (1:1)", nil)
}

func TestHarmonyFail99(t *testing.T) {
	testFail(t, "(a, (b)) => 42", "Invalid parenthesized assignment pattern at (1:4)", nil)
}

func TestHarmonyFail100(t *testing.T) {
	testFail(t, "\"use strict\"; (eval = 10) => 42",
		"Invalid binding `eval` at (1:15)", nil)
}

func TestHarmonyFail101(t *testing.T) {
	testFail(t, "\"use strict\"; eval => 42",
		"Invalid binding `eval` at (1:14)", nil)
}

func TestHarmonyFail102(t *testing.T) {
	testFail(t, "\"use strict\"; arguments => 42",
		"Invalid binding `arguments` at (1:14)", nil)
}

func TestHarmonyFail103(t *testing.T) {
	testFail(t, "\"use strict\"; (eval, a) => 42",
		"Invalid binding `eval` at (1:15)", nil)
}

func TestHarmonyFail104(t *testing.T) {
	testFail(t, "\"use strict\"; (arguments, a) => 42",
		"Invalid binding `arguments` at (1:15)", nil)
}

func TestHarmonyFail105(t *testing.T) {
	testFail(t, "\"use strict\"; (eval, a = 10) => 42",
		"Invalid binding `eval` at (1:15)", nil)
}

func TestHarmonyFail106(t *testing.T) {
	testFail(t, "(a, a) => 42", "Parameter name clash at (1:4)", nil)
}

func TestHarmonyFail107(t *testing.T) {
	testFail(t, "function foo(a, a = 2) {}", "Parameter name clash at (1:16)", nil)
}

func TestHarmonyFail108(t *testing.T) {
	testFail(t, "\"use strict\"; (a, a) => 42", "Parameter name clash at (1:18)", nil)
}

func TestHarmonyFail109(t *testing.T) {
	testFail(t, "\"use strict\"; (a) => 00", "Octal literals are not allowed in strict mode at (1:21)", nil)
}

func TestHarmonyFail110(t *testing.T) {
	testFail(t, "() <= 42", "Unexpected token `<=` at (1:1)", nil)
}

func TestHarmonyFail111(t *testing.T) {
	testFail(t, "(10) => 00", "Unexpected token at (1:1)", nil)
}

func TestHarmonyFail112(t *testing.T) {
	testFail(t, "(10, 20) => 00", "Unexpected token at (1:1)", nil)
}

func TestHarmonyFail113(t *testing.T) {
	testFail(t, "yield v", "Unexpected token `yield` at (1:0)", nil)
}

func TestHarmonyFail114(t *testing.T) {
	testFail(t, "yield 10", "Unexpected token `yield` at (1:0)", nil)
}

func TestHarmonyFail115(t *testing.T) {
	testFail(t, "void { [1, 2]: 3 };", "Unexpected token `,` at (1:9)", nil)
}

func TestHarmonyFail116(t *testing.T) {
	testFail(t, "let [this] = [10]", "Unexpected token `this` at (1:5)", nil)
}

func TestHarmonyFail117(t *testing.T) {
	testFail(t, "let {this} = x", "Unexpected token `this` at (1:5)", nil)
}

func TestHarmonyFail118(t *testing.T) {
	testFail(t, "let [function] = [10]", "Unexpected token `function` at (1:5)", nil)
}

func TestHarmonyFail119(t *testing.T) {
	testFail(t, "let [function] = x", "Unexpected token `function` at (1:5)", nil)
}

func TestHarmonyFail120(t *testing.T) {
	testFail(t, "([function] = [10])", "Unexpected token `]` at (1:10)", nil)
}

func TestHarmonyFail121(t *testing.T) {
	testFail(t, "([this] = [10])", "Assigning to rvalue at (1:2)", nil)
}

func TestHarmonyFail122(t *testing.T) {
	testFail(t, "({this} = x)", "Unexpected token `this` at (1:2)", nil)
}

func TestHarmonyFail123(t *testing.T) {
	testFail(t, "var x = {this}", "Unexpected token `this` at (1:9)", nil)
}

func TestHarmonyFail124(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "(function () { yield 10 })", "Unexpected token at (1:21)", opts)
}

func TestHarmonyFail125(t *testing.T) {
	testFail(t, "let let", "Invalid binding `let` at (1:4)", nil)
}

func TestHarmonyFail126(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "const let", "let is disallowed as a lexically bound name at (1:6)", opts)
}

func TestHarmonyFail127(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "let { let } = {};", "let is disallowed as a lexically bound name at (1:6)", opts)
}

func TestHarmonyFail128(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "const { let } = {};", "let is disallowed as a lexically bound name at (1:8)", opts)
}

func TestHarmonyFail129(t *testing.T) {
	testFail(t, "let [let] = [];", "Invalid binding `let` at (1:5)", nil)
}

func TestHarmonyFail130(t *testing.T) {
	testFail(t, "const [let] = [];", "Invalid binding `let` at (1:7)", nil)
}

func TestHarmonyFail131(t *testing.T) {
	testFail(t, "'use strict'; let + 1", "Unexpected token `+` at (1:18)", nil)
}

func TestHarmonyFail132(t *testing.T) {
	testFail(t, "'use strict'; let let", "Invalid binding `let` at (1:18)", nil)
}

func TestHarmonyFail133(t *testing.T) {
	testFail(t, "'use strict'; const let", "Invalid binding `let` at (1:20)", nil)
}

func TestHarmonyFail134(t *testing.T) {
	testFail(t, "'use strict'; let { let } = {};", "Unexpected token `let` at (1:20)", nil)
}

func TestHarmonyFail135(t *testing.T) {
	testFail(t, "'use strict'; const { let } = {};", "Unexpected token `let` at (1:22)", nil)
}

func TestHarmonyFail136(t *testing.T) {
	testFail(t, "'use strict'; let [let] = [];", "Invalid binding `let` at (1:19)", nil)
}

func TestHarmonyFail137(t *testing.T) {
	testFail(t, "'use strict'; const [let] = [];", "Invalid binding `let` at (1:21)", nil)
}

func TestHarmonyFail138(t *testing.T) {
	testFail(t, "(function() { \"use strict\"; f(yield v) })", "Unexpected token `yield` at (1:30)", nil)
}

func TestHarmonyFail139(t *testing.T) {
	testFail(t, "var obj = { *test** }", "Unexpected token `**` at (1:17)", nil)
}

func TestHarmonyFail140(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "class A extends yield B { }", "Unexpected token `yield` at (1:16)", opts)
}

func TestHarmonyFail141(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "class default", "Unexpected token `default` at (1:6)", opts)
}

func TestHarmonyFail142(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "class let {}", "Invalid binding `let` at (1:6)", opts)
}

func TestHarmonyFail143(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "`test", "Unterminated template at (1:0)", opts)
}

func TestHarmonyFail144(t *testing.T) {
	testFail(t, "switch `test`", "Unexpected token `template head` at (1:7)", nil)
}

func TestHarmonyFail145(t *testing.T) {
	testFail(t, "`hello ${10 `test`", "Unexpected token `EOF` at (1:18)", nil)
}

func TestHarmonyFail146(t *testing.T) {
	testFail(t, "`hello ${10;test`", "Unexpected token `;` at (1:11)", nil)
}

func TestHarmonyFail147(t *testing.T) {
	testFail(t, "function a() 1 // expression closure is not supported", "Unexpected token `number` at (1:13)", nil)
}

func TestHarmonyFail148(t *testing.T) {
	testFail(t, "({ \"chance\" }) = obj", "Invalid parenthesized assignment pattern at (1:0)", nil)
}

func TestHarmonyFail149(t *testing.T) {
	testFail(t, "({ 42 }) = obj", "Invalid parenthesized assignment pattern at (1:0)", nil)
}

func TestHarmonyFail150(t *testing.T) {
	testFail(t, "function f(a, ...b, c)", "Rest element must be last element at (1:18)", nil)
}

func TestHarmonyFail151(t *testing.T) {
	testFail(t, "function f(a, ...b = 0)", "Rest elements cannot have a default value at (1:17)", nil)
}

func TestHarmonyFail152(t *testing.T) {
	testFail(t, "(([a, ...b = 0]) => {})", "Rest elements cannot have a default value at (1:9)", nil)
}

func TestHarmonyFail153(t *testing.T) {
	testFail(t, "[a, ...b = 0] = []", "Rest elements cannot have a default value at (1:7)", nil)
}

func TestHarmonyFail154(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "function x(...{ a }){}", "Binding pattern is not permitted as rest operator's argument at (1:14)", opts)
}

func TestHarmonyFail155(t *testing.T) {
	testFail(t, "\"use strict\"; function x(a, { a }){}", "Parameter name clash at (1:30)", nil)
}

func TestHarmonyFail156(t *testing.T) {
	testFail(t, "\"use strict\"; function x({ b: { a } }, [{ b: { a } }]){}",
		"Parameter name clash at (1:47)", nil)
}

func TestHarmonyFail157(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "\"use strict\"; function x(a, ...[a]){}",
		"Binding pattern is not permitted as rest operator's argument at (1:31)", opts)
}

func TestHarmonyFail158(t *testing.T) {
	testFail(t, "(...a, b) => {}", "Rest element must be last element at (1:5)", nil)
}

func TestHarmonyFail159(t *testing.T) {
	testFail(t, "([ 5 ]) => {}", "Assigning to rvalue at (1:3)", nil)
}

func TestHarmonyFail160(t *testing.T) {
	testFail(t, "({ 5 }) => {}", "Unexpected token at (1:3)", nil)
}

func TestHarmonyFail161(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BINDING_REST_ELEM_NESTED)
	testFail(t, "(...[ 5 ]) => {}", "Binding pattern is not permitted as rest operator's argument at (1:4)", opts)
}

func TestHarmonyFail162(t *testing.T) {
	testFail(t, "[...a, b] = c",
		"Rest element must be last element at (1:5)", nil)
}

func TestHarmonyFail163(t *testing.T) {
	testFail(t, "({ t(eval) { \"use strict\"; } });",
		"Invalid binding `eval` at (1:5)", nil)
}

func TestHarmonyFail164(t *testing.T) {
	testFail(t, "\"use strict\"; `${test}\\02`;", "Octal escape sequences are not allowed in template strings at (1:22)", nil)
}

func TestHarmonyFail165(t *testing.T) {
	testFail(t, "if (1) import \"acorn\";",
		"'import' and 'export' may only appear at the top level at (1:7)", nil)
}

func TestHarmonyFail166(t *testing.T) {
	testFail(t, "[...a, ] = b", "Rest element must be last element at (1:5)", nil)
}

func TestHarmonyFail167(t *testing.T) {
	testFail(t, "if (b,...a, );", "Unexpected token `...` at (1:6)", nil)
}

func TestHarmonyFail168(t *testing.T) {
	testFail(t, "(b, ...a)", "Unexpected token at (1:4)", nil)
}

func TestHarmonyFail169(t *testing.T) {
	testFail(t, "switch (cond) { case 10: let a = 20; ",
		"Unexpected token `EOF` at (1:37)", nil)
}

func TestHarmonyFail170(t *testing.T) {
	testFail(t, "\"use strict\"; (eval) => 42", "Invalid binding `eval` at (1:15)", nil)
}

func TestHarmonyFail171(t *testing.T) {
	testFail(t, "(eval) => { \"use strict\"; 42 }", "Invalid binding `eval` at (1:1)", nil)
}

func TestHarmonyFail172(t *testing.T) {
	testFail(t, "({ get test() { } }) => 42",
		"Object pattern can't contain getter or setter at (1:3)", nil)
}

func TestHarmonyFail173(t *testing.T) {
	testFail(t, "obj = {x = 0}",
		"Shorthand property assignments are valid only in destructuring patterns at (1:9)", nil)
}

func TestHarmonyFail174(t *testing.T) {
	testFail(t, "f({x = 0})",
		"Shorthand property assignments are valid only in destructuring patterns at (1:5)", nil)
}

func TestHarmonyFail175(t *testing.T) {
	testFail(t, "(localVar |= defaultValue) => {}",
		"Unexpected token at (1:10)", nil)
}

func TestHarmonyFail176(t *testing.T) {
	testFail(t, "let [x]", "Complex binding patterns require an initialization value at (1:7)", nil)
}

func TestHarmonyFail178(t *testing.T) {
	testFail(t, "var _𖫵 = 11;", "Unexpected token at (1:5)", nil)
}

func TestHarmonyFail179(t *testing.T) {
	testFail(t, "var 𫠞_ = 12;", "Unexpected character at (1:4)", nil)
}

func TestHarmonyFail180(t *testing.T) {
	testFail(t, "var 𫠟_ = 10;", "Unexpected character at (1:4)", nil)
}

func TestHarmonyFail181(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "if (1) let x = 10;", "Unexpected token `identifier` at (1:7)", opts)
}

func TestHarmonyFail182(t *testing.T) {
	testFail(t, "for (;;) const x = 10;", "Unexpected token `const` at (1:9)", nil)
}

func TestHarmonyFail183(t *testing.T) {
	testFail(t, "while (1) function foo(){}",
		"function declarations can't appear in single-statement context at (1:10)", nil)
}

func TestHarmonyFail184(t *testing.T) {
	testFail(t, "if (1) ; else class Cls {}", "Unexpected token `class` at (1:14)", nil)
}

func TestHarmonyFail185(t *testing.T) {
	testFail(t, "'use strict'; [...eval] = arr",
		"Assigning to `eval` in strict mode at (1:18)", nil)
}

func TestHarmonyFail186(t *testing.T) {
	testFail(t, "'use strict'; ({eval = defValue} = obj)", "Assigning to `eval` in strict mode at (1:16)", nil)
}

func TestHarmonyFail187(t *testing.T) {
	testFail(t, "[...eval] = arr",
		"Assigning to `eval` in strict mode at (1:4)", nil)
}

func TestHarmonyFail188(t *testing.T) {
	testFail(t, "function* y({yield}) {}", "Unexpected token `yield` at (1:13)", nil)
}

func TestHarmonyFail189(t *testing.T) {
	testFail(t, "new.prop", "The only valid meta property for new is `new.target` at (1:4)", nil)
}

func TestHarmonyFail190(t *testing.T) {
	testFail(t, "new.target", "`new.target` can only be used in functions at (1:0)", nil)
}

func TestHarmonyFail191(t *testing.T) {
	testFail(t, "let y = () => new.target",
		"`new.target` can only be used in functions at (1:14)", nil)
}

func TestHarmonyFail192(t *testing.T) {
	testFail(t, "`\\07`", "Octal escape sequences are not allowed in template strings at (1:1)", nil)
}

func TestHarmonyFail193(t *testing.T) {
	testFail(t, "(function(){ 'use strict'; '\\07'; })", "Octal escape sequences are not allowed in strict mode at (1:27)", nil)
}

func TestHarmonyFail194(t *testing.T) {
	testFail(t, "x = { method() 42 }", "Unexpected token at (1:15)", nil)
}

func TestHarmonyFail195(t *testing.T) {
	testFail(t, "x = { get method() 42 }", "Unexpected token at (1:19)", nil)
}

func TestHarmonyFail196(t *testing.T) {
	testFail(t, "x = { set method(val) v = val }", "Unexpected token at (1:22)", nil)
}

func TestHarmonyFail197(t *testing.T) {
	testFail(t, "/\\u{110000}/u", "Code point out of bounds at (1:3)", nil)
}

func TestHarmonyFail198(t *testing.T) {
	testFail(t, "super", "'super' is only allowed in object methods and classes at (1:0)", nil)
}

func TestHarmonyFail199(t *testing.T) {
	testFail(t, "class A { get prop(x) {} }", "Getter must not have any formal parameters at (1:19)", nil)
}

func TestHarmonyFail200(t *testing.T) {
	testFail(t, "class A { set prop() {} }", "Setter must have exactly one formal parameter at (1:18)", nil)
}

func TestHarmonyFail201(t *testing.T) {
	testFail(t, "class A { set prop(x, y) {} }", "Setter must have exactly one formal parameter at (1:18)", nil)
}

func TestHarmonyFail202(t *testing.T) {
	testFail(t, "({ __proto__: 1, __proto__: 2 })", "Redefinition of property at (1:17)", nil)
}

func TestHarmonyFail203(t *testing.T) {
	testFail(t, "({ '__proto__': 1, __proto__: 2 })", "Redefinition of property at (1:19)", nil)
}

func TestHarmonyFail204(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "var await = 0", "Unexpected token `await` at (1:4)", opts)
}

func TestHarmonyFail205(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_CHK_REGEXP_FLAGS)
	testFail(t, "/[a-z]/s", "Invalid regular expression flag at (1:0)", opts)
}

func TestHarmonyFail206(t *testing.T) {
	testFail(t, "[...x in y] = []", "Invalid rest operator's argument at (1:4)", nil)
}

func TestHarmonyFail207(t *testing.T) {
	testFail(t, "export let x = a; export function x() {}",
		"Identifier `x` has already been declared at (1:34)", nil)
}

func TestHarmonyFail208(t *testing.T) {
	testFail(t, "export let [{x = 2}] = a; export {x}",
		"Duplicate export `x` at (1:34)", nil)
}

func TestHarmonyFail209(t *testing.T) {
	testFail(t, "export default 100; export default 3",
		"Duplicate export `default` at (1:27)", nil)
}

func TestHarmonyFail210(t *testing.T) {
	testFail(t, "function foo() { 'use strict'; return {yield} }",
		"Unexpected token `yield` at (1:39)", nil)
}

func TestHarmonyFail211(t *testing.T) {
	testFail(t, "function foo() { 'use strict'; var {arguments} = {} }",
		"Unexpected token `arguments` at (1:36)", nil)
}

func TestHarmonyFail212(t *testing.T) {
	testFail(t, "function foo() { 'use strict'; var {eval} = {} }",
		"Unexpected token `eval` at (1:36)", nil)
}

func TestHarmonyFail213(t *testing.T) {
	testFail(t, "function foo() { 'use strict'; var {arguments = 1} = {} }",
		"Unexpected token `arguments` at (1:36)", nil)
}

func TestHarmonyFail214(t *testing.T) {
	testFail(t, "function foo() { 'use strict'; var {eval = 1} = {} }",
		"Unexpected token `eval` at (1:36)", nil)
}

func TestHarmonyFail215(t *testing.T) {
	testFail(t, "function* wrap() { function* foo(a = 1 + (yield)) {} }",
		"Yield expression can't be used in parameter at (1:42)", nil)
}

func TestHarmonyFail216(t *testing.T) {
	testFail(t, "function* wrap() { return (a = 1 + (yield)) => a }",
		"Yield expression cannot be a default value at (1:36)", nil)
}

func TestHarmonyFail217(t *testing.T) {
	testFail(t, "(function* g() {\nfor (yield '' in {}; ; ) ;\n }", "Assigning to rvalue at (2:5)", nil)
}

func TestHarmonyFail218(t *testing.T) {
	testFail(t, "(function* yield() {})", "Invalid binding `yield` at (1:11)", nil)
}

func TestHarmonyFail219(t *testing.T) {
	testFail(t, "function* wrap() {\nfunction* yield() {}\n}",
		"Invalid binding `yield` at (2:10)", nil)
}

// Forbid yield expressions in default parameters:

func TestHarmonyFail220(t *testing.T) {
	testFail(t, "function* foo(a = yield b) {}",
		"Yield expression can't be used in parameter at (1:18)", nil)
}

func TestHarmonyFail221(t *testing.T) {
	testFail(t, "(function* foo(a = yield b) {})",
		"Yield expression can't be used in parameter at (1:19)", nil)
}

func TestHarmonyFail222(t *testing.T) {
	testFail(t, "({*foo(a = yield b) {}})",
		"Yield expression can't be used in parameter at (1:11)", nil)
}

func TestHarmonyFail223(t *testing.T) {
	testFail(t, "(class {*foo(a = yield b) {}})",
		"Yield expression can't be used in parameter at (1:17)", nil)
}

func TestHarmonyFail224(t *testing.T) {
	testFail(t, "function* foo(a = class extends (yield b) {}) {}",
		"Yield expression can't be used in parameter at (1:33)", nil)
}

func TestHarmonyFail225(t *testing.T) {
	testFail(t, "function* wrap() {\n(a = yield b) => a\n}",
		"Yield expression cannot be a default value at (2:5)", nil)
}

func TestHarmonyFail226(t *testing.T) {
	testFail(t, "class B { constructor(a = super()) { return a }}",
		"super() call outside constructor of a subclass at (1:26)", nil)
}

func TestHarmonyFail227(t *testing.T) {
	testFail(t, "function* wrap() {\n({a = yield b} = obj) => a\n}",
		"Yield expression cannot be a default value at (2:6)", nil)
}

func TestHarmonyFail228(t *testing.T) {
	testFail(t, "({*foo: 1})", "Unexpected token `:` at (1:6)", nil)
}

func TestHarmonyFail229(t *testing.T) {
	testFail(t, "export { default} from './y.js';\nexport default 42;",
		"Duplicate export `default` at (2:7)", nil)
}

func TestHarmonyFail230(t *testing.T) {
	testFail(t, "export * from foo", "Unexpected token `identifier` at (1:14)", nil)
}

func TestHarmonyFail231(t *testing.T) {
	testFail(t, "export { bar } from foo", "Unexpected token `identifier` at (1:20)", nil)
}

func TestHarmonyFail232(t *testing.T) {
	testFail(t, "foo: class X {}", "Unexpected token `class` at (1:5)", nil)
}

func TestHarmonyFail233(t *testing.T) {
	testFail(t, "'use strict'; bar: function x() {}",
		"function declarations can't appear in single-statement context at (1:19)", nil)
}

func TestHarmonyFail234(t *testing.T) {
	testFail(t, "'use strict'; bar: function* x() {}",
		"function declarations can't appear in single-statement context at (1:19)", nil)
}

func TestHarmonyFail235(t *testing.T) {
	testFail(t, "bar: function* x() {}",
		"function declarations can't appear in single-statement context at (1:5)", nil)
}

func TestHarmonyFail236(t *testing.T) {
	testFail(t, "({x, y}) = {}", "Invalid parenthesized assignment pattern at (1:0)", nil)
}

func TestHarmonyFail237(t *testing.T) {
	testFail(t, "var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:17)", nil)
}

func TestHarmonyFail238(t *testing.T) {
	testFail(t, "{ var foo = 1; let foo = 1; }",
		"Identifier `foo` has already been declared at (1:19)", nil)
}

func TestHarmonyFail239(t *testing.T) {
	testFail(t, "let bar; var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:26)", nil)
}

func TestHarmonyFail240(t *testing.T) {
	testFail(t, "{ let bar; var foo = 1; let foo = 1; }",
		"Identifier `foo` has already been declared at (1:28)", nil)
}

func TestHarmonyFail241(t *testing.T) {
	testFail(t, "let foo = 1; var foo = 1;",
		"Identifier `foo` has already been declared at (1:17)", nil)
}

func TestHarmonyFail242(t *testing.T) {
	testFail(t, "let bar; let foo = 1; var foo = 1;",
		"Identifier `foo` has already been declared at (1:26)", nil)
}

func TestHarmonyFail243(t *testing.T) {
	testFail(t, "{ let bar; let foo = 1; var foo = 1; }",
		"Identifier `foo` has already been declared at (1:28)", nil)
}

func TestHarmonyFail244(t *testing.T) {
	testFail(t, "let foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:17)", nil)
}

func TestHarmonyFail245(t *testing.T) {
	testFail(t, "var foo = 1; const foo = 1;",
		"Identifier `foo` has already been declared at (1:19)", nil)
}

func TestHarmonyFail246(t *testing.T) {
	testFail(t, "const foo = 1; var foo = 1;",
		"Identifier `foo` has already been declared at (1:19)", nil)
}

func TestHarmonyFail247(t *testing.T) {
	testFail(t, "var [foo] = [1]; let foo = 1;",
		"Identifier `foo` has already been declared at (1:21)", nil)
}

func TestHarmonyFail248(t *testing.T) {
	testFail(t, "var [{ bar: [foo] }] = x; let {foo} = 1;",
		"Identifier `foo` has already been declared at (1:31)", nil)
}

func TestHarmonyFail249(t *testing.T) {
	testFail(t, "if (x) var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:24)", nil)
}

func TestHarmonyFail250(t *testing.T) {
	testFail(t, "if (x) {} else var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:32)", nil)
}

func TestHarmonyFail251(t *testing.T) {
	testFail(t, "if (x) var foo = 1; else {} let foo = 1;",
		"Identifier `foo` has already been declared at (1:32)", nil)
}

func TestHarmonyFail252(t *testing.T) {
	testFail(t, "if (x) {} else if (y) {} else var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:47)", nil)
}

func TestHarmonyFail253(t *testing.T) {
	testFail(t, "while (x) var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:27)", nil)
}

func TestHarmonyFail254(t *testing.T) {
	testFail(t, "do var foo = 1; while (x) let foo = 1;",
		"Identifier `foo` has already been declared at (1:30)", nil)
}

func TestHarmonyFail255(t *testing.T) {
	testFail(t, "for (;;) var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:26)", nil)
}

func TestHarmonyFail256(t *testing.T) {
	testFail(t, "for (const x of y) var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:36)", nil)
}

func TestHarmonyFail257(t *testing.T) {
	testFail(t, "for (const x in y) var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:36)", nil)
}

func TestHarmonyFail258(t *testing.T) {
	testFail(t, "label: var foo = 1; let foo = 1;",
		"Identifier `foo` has already been declared at (1:24)", nil)
}

func TestHarmonyFail259(t *testing.T) {
	testFail(t, "switch (x) { case 0: var foo = 1 } let foo = 1;",
		"Identifier `foo` has already been declared at (1:39)", nil)
}

func TestHarmonyFail260(t *testing.T) {
	testFail(t, "try { var foo = 1; } catch (e) {} let foo = 1;",
		"Identifier `foo` has already been declared at (1:38)", nil)
}

func TestHarmonyFail261(t *testing.T) {
	testFail(t, "function foo() {} let foo = 1;",
		"Identifier `foo` has already been declared at (1:22)", nil)
}

func TestHarmonyFail262(t *testing.T) {
	testFail(t, "{ var foo = 1; } let foo = 1;",
		"Identifier `foo` has already been declared at (1:21)", nil)
}

func TestHarmonyFail263(t *testing.T) {
	testFail(t, "let foo = 1; { var foo = 1; }",
		"Identifier `foo` has already been declared at (1:19)", nil)
}

func TestHarmonyFail264(t *testing.T) {
	testFail(t, "let foo = 1; function x(foo) {} { var foo = 1; }",
		"Identifier `foo` has already been declared at (1:38)", nil)
}

func TestHarmonyFail265(t *testing.T) {
	testFail(t, "if (x) { if (y) var foo = 1; } let foo = 1;",
		"Identifier `foo` has already been declared at (1:35)", nil)
}

func TestHarmonyFail266(t *testing.T) {
	testFail(t, "var foo = 1; function x() {} let foo = 1;",
		"Identifier `foo` has already been declared at (1:33)", nil)
}

func TestHarmonyFail267(t *testing.T) {
	testFail(t, "{ let foo = 1; { let foo = 2; } let foo = 1; }",
		"Identifier `foo` has already been declared at (1:36)", nil)
}

func TestHarmonyFail268(t *testing.T) {
	testFail(t, "for (var foo of y) {} let foo = 1;",
		"Identifier `foo` has already been declared at (1:26)", nil)
}

func TestHarmonyFail269(t *testing.T) {
	testFail(t, "function x(foo) { let foo = 1; }",
		"Identifier `foo` has already been declared at (1:22)", nil)
}

func TestHarmonyFail270(t *testing.T) {
	testFail(t, "var [...foo] = x; let foo = 1;",
		"Identifier `foo` has already been declared at (1:22)", nil)
}

func TestHarmonyFail271(t *testing.T) {
	testFail(t, "foo => { let foo; }",
		"Identifier `foo` has already been declared at (1:13)", nil)
}

func TestHarmonyFail272(t *testing.T) {
	testFail(t, "({ x(foo) { let foo; } })", "Identifier `foo` has already been declared at (1:16)", nil)
}

func TestHarmonyFail273(t *testing.T) {
	testFail(t, "try {} catch (foo) { let foo = 1; }",
		"Identifier `foo` has already been declared at (1:25)", nil)
}

func TestHarmonyFail274(t *testing.T) {
	testFail(t, "(x) => {} + 2", "Unexpected token at (1:10)", nil)
}

func TestHarmonyFail275(t *testing.T) {
	testFail(t, "'use strict'; { function f() {} function f() {} }",
		"Identifier `f` has already been declared at (1:41)", nil)
}

func TestHarmonyFail276(t *testing.T) {
	testFail(t, "{ function f() {} function* f() {} }",
		"Identifier `f` has already been declared at (1:28)", nil)
}

func TestHarmonyFail278(t *testing.T) {
	testFail(t, "{ function* f() {} function f() {} }",
		"Identifier `f` has already been declared at (1:28)", nil)
}

func TestHarmonyFail279(t *testing.T) {
	testFail(t, "class A extends B { constructor() { super } }",
		"Unexpected token `super` at (1:36)", nil)
}

func TestHarmonyFail280(t *testing.T) {
	testFail(t, "class A extends B { constructor() { super; } }",
		"Unexpected token `super` at (1:36)", nil)
}

func TestHarmonyFail281(t *testing.T) {
	testFail(t, "class A extends B { constructor() { (super)() } }",
		"Unexpected token `super` at (1:37)", nil)
}

func TestHarmonyFail282(t *testing.T) {
	testFail(t, "class A extends B { foo() { (super).foo } }",
		"Unexpected token `super` at (1:29)", nil)
}

func TestHarmonyFail283(t *testing.T) {
	testFail(t, "for (let x of y, z) {}", "Unexpected token `,` at (1:15)", nil)
}

func TestHarmonyFail284(t *testing.T) {
	testFail(t, "[...foo, bar] = b",
		"Rest element must be last element at (1:7)", nil)
}

func TestHarmonyFail285(t *testing.T) {
	testFail(t, "for (let [...foo, bar] in qux);",
		"Rest element must be last element at (1:16)", nil)
}

func TestHarmonyFail286(t *testing.T) {
	testFail(t, "var f;\nfunction f() {}",
		"Identifier `f` has already been declared at (2:9)", nil)
}

func TestHarmonyFail287(t *testing.T) {
	testFail(t, "({ a = 42, b: c = d })",
		"Shorthand property assignments are valid only in destructuring patterns at (1:5)", nil)
}

func TestHarmonyFail288(t *testing.T) {
	testFail(t, "for (let a = b => b in c;;);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail289(t *testing.T) {
	testFail(t, "for (let a = b => c => d in e;;);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail290(t *testing.T) {
	testFail(t, "for (var a = b => c in d;;);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail291(t *testing.T) {
	testFail(t, "for (var a = b => c => d in e;;);",
		"for-in loop variable declaration may not have an initializer at (1:5)", nil)
}

func TestHarmonyFail292(t *testing.T) {
	testFail(t, "for (x => x in y;;);", "Assigning to rvalue at (1:5)", nil)
}

func TestHarmonyFail293(t *testing.T) {
	testFail(t, "for (x => y => y in z;;);", "Assigning to rvalue at (1:5)", nil)
}

func TestHarmonyFail294(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_OPT_CATCH_PARAM)
	testFail(t, "try {} catch {}", "Unexpected token `{` at (1:13)", opts)
}

func TestHarmonyFail295(t *testing.T) {}

func TestHarmonyFail296(t *testing.T) {}

func TestHarmonyFail297(t *testing.T) {}

func TestHarmonyFail298(t *testing.T) {}

func TestHarmonyFail299(t *testing.T) {}

func TestHarmonyFail300(t *testing.T) {}

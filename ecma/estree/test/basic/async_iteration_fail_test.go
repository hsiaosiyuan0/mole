package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestAsyncIterationFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_GLOBAL_ASYNC)
	TestFail(t, "for await (x of xs);", "Unexpected token `await` at (1:4)", opts)
}

func TestAsyncIterationFail2(t *testing.T) {
	TestFail(t, "function f() { for await (x of xs); }", "Unexpected token `await` at (1:19)", nil)
}

func TestAsyncIterationFail3(t *testing.T) {
	TestFail(t, "f = function() { for await (x of xs); }", "Unexpected token `await` at (1:21)", nil)
}

func TestAsyncIterationFail4(t *testing.T) {
	TestFail(t, "f = () => { for await (x of xs); }", "Unexpected token `await` at (1:16)", nil)
}

func TestAsyncIterationFail5(t *testing.T) {
	TestFail(t, "async function f() { () => { for await (x of xs); } }", "Unexpected token `await` at (1:33)", nil)
}

func TestAsyncIterationFail6(t *testing.T) {
	TestFail(t, "async function f() { for await (x in xs); }", "Unexpected token `in` at (1:34)", nil)
}

func TestAsyncIterationFail7(t *testing.T) {
	TestFail(t, "async function f() { for await (x;;); }", "Unexpected token `;` at (1:33)", nil)
}

func TestAsyncIterationFail8(t *testing.T) {
	TestFail(t, "async function f() { for await (let x = 0;;); }", "Unexpected token `;` at (1:41)", nil)
}

func TestAsyncIterationFail9(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_ITERATION)
	TestFail(t, "async function f() { for await (x of xs); }", "Unexpected token `await` at (1:25)", opts)
}

func TestAsyncIterationFail10(t *testing.T) {
	TestFail(t, "async function* f() { () => await a; }",
		"Cannot use keyword 'await' outside an async function at (1:28)", nil)
}

func TestAsyncIterationFail11(t *testing.T) {
	TestFail(t, "async function* f() { () => yield a; }", "Unexpected token `yield` at (1:28)", nil)
}

func TestAsyncIterationFail12(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_AWAIT)
	TestFail(t, "async function* f() { await a; yield b; }", "Unexpected token at (1:6)", opts)
}

func TestAsyncIterationFail13(t *testing.T) {
	TestFail(t, "f = async function*() { () => await a; }",
		"Cannot use keyword 'await' outside an async function at (1:30)", nil)
}

func TestAsyncIterationFail14(t *testing.T) {
	TestFail(t, "f = async function*() { () => yield a; }", "Unexpected token `yield` at (1:30)", nil)
}

func TestAsyncIterationFail15(t *testing.T) {
	TestFail(t, "obj = { async\n* f() {} }", "Unexpected token `*` at (2:0)", nil)
}

func TestAsyncIterationFail16(t *testing.T) {
	TestFail(t, "obj = { *async f() {}", "Unexpected token `identifier` at (1:15)", nil)
}

func TestAsyncIterationFail17(t *testing.T) {
	TestFail(t, "obj = { *async* f() {}", "Unexpected token `*` at (1:14)", nil)
}

func TestAsyncIterationFail18(t *testing.T) {
	TestFail(t, "obj = { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:27)", nil)
}

func TestAsyncIterationFail19(t *testing.T) {
	TestFail(t, "obj = { async* f() { () => yield a; } }", "Unexpected token `yield` at (1:27)", nil)
}

func TestAsyncIterationFail20(t *testing.T) {
	// skipped, it's legal in chrome and ff
	// TestFail(t, "class A { async\n* f() {} }", "Unexpected token (2:0)", nil)
}

func TestAsyncIterationFail21(t *testing.T) {
	TestFail(t, "class A { *async f() {} }", "Unexpected token `identifier` at (1:17)", nil)
}

func TestAsyncIterationFail22(t *testing.T) {
	TestFail(t, "class A { *async* f() {} }", "Unexpected token `*` at (1:16)", nil)
}

func TestAsyncIterationFail23(t *testing.T) {
	TestFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestAsyncIterationFail24(t *testing.T) {
	TestFail(t, "class A { async* f() { () => yield a; } }",
		"Unexpected token `yield` at (1:29)", nil)
}

func TestAsyncIterationFail25(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_AWAIT)
	TestFail(t, "f = async function*() { await a; yield b; }", "Unexpected token at (1:10)", opts)
}

func TestAsyncIterationFail26(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_AWAIT)
	TestFail(t, "obj = { async* f() { await a; yield b; } }", "Unexpected token `*` at (1:13)", opts)
}

func TestAsyncIterationFail27(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_AWAIT)
	TestFail(t, "class A { async* f() { await a; yield b; } }", "Unexpected token `*` at (1:15)", opts)
}

func TestAsyncIterationFail28(t *testing.T) {
	TestFail(t, "({ \\u0061sync *method(){} })",
		"Keyword must not contain escaped characters at (1:3)", nil)
}

func TestAsyncIterationFail29(t *testing.T) {
	TestFail(t, "void \\u0061sync function* f(){};", "Keyword must not contain escaped characters at (1:5)", nil)
}

func TestAsyncIterationFail30(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	TestFail(t, "for ( ; false; ) async function* g() {}",
		"function declarations can't appear in single-statement context at (1:17)", opts)
}

func TestAsyncIterationFail31(t *testing.T) {
	TestFail(t, "({async\n    foo() { }})", "Unexpected token `identifier` at (2:4)", nil)
}

func TestAsyncIterationFail32(t *testing.T) {
	TestFail(t, "for (async of [1]) {}",
		"The left-hand side of a for-of loop may not be 'async' at (1:5)", nil)
}

func TestAsyncIterationFail33(t *testing.T) {
	TestFail(t, "async function f() { for await (;;); }", "Unexpected token `;` at (1:32)", nil)
}

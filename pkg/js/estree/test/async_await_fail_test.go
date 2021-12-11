package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestAsyncAwaitFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_GENERATOR)
	testFail(t, "async function* foo() { }", "Unexpected token `*` at (1:14)", opts)
}

func TestAsyncAwaitFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE)
	testFail(t, "async function wrap() {\nasync function await() { }\n}",
		"Invalid binding `await` at (2:15)", opts)
}

func TestAsyncAwaitFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "async function foo(await) { }",
		"Invalid binding `await` at (1:19)", opts)
}

func TestAsyncAwaitFail4(t *testing.T) {
	testFail(t, "for ( ; false; ) async function* g() {}",
		"function declarations can't appear in single-statement context at (1:17)", nil)
}

func TestAsyncAwaitFail5(t *testing.T) {
	testFail(t, "async function foo() { return {await} }",
		"Unexpected token `await` at (1:31)", nil)
}

func TestAsyncAwaitFail6(t *testing.T) {
	testFail(t, "(async\nfunction foo() { })", "Unexpected token `function` at (2:0)", nil)
}

func TestAsyncAwaitFail7(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_GENERATOR)
	testFail(t, "(async function* foo() { })",
		"Unexpected token `*` at (1:15)", opts)
}

func TestAsyncAwaitFail8(t *testing.T) {
	testFail(t, "(async function await() { })",
		"Invalid binding `await` at (1:16)", nil)
}

func TestAsyncAwaitFail9(t *testing.T) {
	testFail(t, "for ( ; false; ) async function* g() {}",
		"function declarations can't appear in single-statement context at (1:17)", nil)
}

func TestAsyncAwaitFail10(t *testing.T) {
	testFail(t, "(async function foo(await) { })",
		"Invalid binding `await` at (1:20)", nil)
}

func TestAsyncAwaitFail11(t *testing.T) {
	testFail(t, "(async function foo() { return {await} })",
		"Unexpected token `await` at (1:32)", nil)
}

func TestAsyncAwaitFail12(t *testing.T) {
	testFail(t, "async ({a = b})",
		"Shorthand property assignments are valid only in destructuring patterns at (1:10)", nil)
}

func TestAsyncAwaitFail13(t *testing.T) {
	testFail(t, "async\n() => a", "Unexpected token at (2:3)", nil)
}

func TestAsyncAwaitFail14(t *testing.T) {
	testFail(t, "async a\n=> a", "Unexpected token `=>` at (2:0)", nil)
}

func TestAsyncAwaitFail15(t *testing.T) {
	testFail(t, "async ()\n=> a", "Unexpected token `=>` at (2:0)", nil)
}

func TestAsyncAwaitFail16(t *testing.T) {
	testFail(t, "async await => 1",
		"Unexpected token at (1:6)", nil)
}

func TestAsyncAwaitFail17(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "async (await) => 1",
		"Invalid binding `await` at (1:7)", opts)
}

func TestAsyncAwaitFail18(t *testing.T) {
	testFail(t, "async (...await) => 1",
		"Invalid binding `await` at (1:10)", nil)
}

func TestAsyncAwaitFail19(t *testing.T) {
	testFail(t, "async ({await}) => 1",
		"Unexpected token `await` at (1:8)", nil)
}

func TestAsyncAwaitFail20(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE)
	testFail(t, "async ({a: await}) => 1",
		"Invalid binding `await` at (1:11)", opts)
}

func TestAsyncAwaitFail21(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "async ([await]) => 1",
		"Invalid binding `await` at (1:8)", opts)
}

func TestAsyncAwaitFail22(t *testing.T) {
	testFail(t, "async ([...await]) => 1",
		"Unexpected token `]` at (1:16)", nil)
}

func TestAsyncAwaitFail23(t *testing.T) {
	testFail(t, "async (b = {await}) => 1",
		"Unexpected token `await` at (1:12)", nil)
}

func TestAsyncAwaitFail24(t *testing.T) {
	testFail(t, "async (b = {a: await}) => 1",
		"Unexpected token `}` at (1:20)", nil)
}

func TestAsyncAwaitFail25(t *testing.T) {
	testFail(t, "async (b = [await]) => 1",
		"Unexpected token `]` at (1:17)", nil)
}

func TestAsyncAwaitFail26(t *testing.T) {
	testFail(t, "async (b = [...await]) => 1",
		"Unexpected token `]` at (1:20)", nil)
}

func TestAsyncAwaitFail27(t *testing.T) {
	testFail(t, "async (b = class await {}) => 1",
		"Invalid binding `await` at (1:17)", nil)
}

func TestAsyncAwaitFail28(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "async (b = (await) => {}) => 1",
		"Invalid binding `await` at (1:12)", opts)
}

func TestAsyncAwaitFail29(t *testing.T) {
	testFail(t, "async (await, b = async()) => 2",
		"Invalid binding `await` at (1:7)", nil)
}

func TestAsyncAwaitFail30(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "async (await, b = async () => {}) => 1",
		"Invalid binding `await` at (1:7)", opts)
}

func TestAsyncAwaitFail31(t *testing.T) {
	testFail(t, "({async\nfoo() { }})", "Unexpected token `identifier` at (2:0)", nil)
}

func TestAsyncAwaitFail32(t *testing.T) {
	testFail(t, "({async get foo() { }})", "Unexpected token `identifier` at (1:12)", nil)
}

func TestAsyncAwaitFail33(t *testing.T) {
	testFail(t, "({async set foo(value) { }})", "Unexpected token `identifier` at (1:12)", nil)
}

func TestAsyncAwaitFail34(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_GENERATOR)
	testFail(t, "({async* foo() { }})", "Unexpected token `*` at (1:7)", opts)
}

func TestAsyncAwaitFail35(t *testing.T) {
	testFail(t, "({async foo() { var await }})",
		"Invalid binding `await` at (1:20)", nil)
}

func TestAsyncAwaitFail36(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "({async foo(await) { }})",
		"Invalid binding `await` at (1:12)", opts)
}

func TestAsyncAwaitFail37(t *testing.T) {
	testFail(t, "({async foo() { return {await} }})",
		"Unexpected token `await` at (1:24)", nil)
}

func TestAsyncAwaitFail38(t *testing.T) {
	testFail(t, "({async foo: 1})", "Unexpected token `:` at (1:11)", nil)
}

func TestAsyncAwaitFail39(t *testing.T) {
	// below is legal in both chrome and ff
	// testFail(t, "class A {async\nfoo() { }}", "Unexpected token (2:0)", nil)
}

func TestAsyncAwaitFail40(t *testing.T) {
	// below is legal in both chrome and ff
	// testFail(t, "class A {static async\nfoo() { }}", "Unexpected token (2:0)", nil)
}

func TestAsyncAwaitFail41(t *testing.T) {
	testFail(t, "class A {async constructor() { }}",
		"Constructor can't be a async at (1:15)", nil)
}

func TestAsyncAwaitFail42(t *testing.T) {
	testFail(t, "class A {async get foo() { }}", "Unexpected token `identifier` at (1:19)", nil)
}

func TestAsyncAwaitFail43(t *testing.T) {
	testFail(t, "class A {async set foo(value) { }}", "Unexpected token `identifier` at (1:19)", nil)
}

func TestAsyncAwaitFail44(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_GENERATOR)
	testFail(t, "class A {async* foo() { }}", "Unexpected token `*` at (1:14)", opts)
}

func TestAsyncAwaitFail45(t *testing.T) {
	testFail(t, "class A {static async get foo() { }}", "Unexpected token `identifier` at (1:26)", nil)
}

func TestAsyncAwaitFail46(t *testing.T) {
	testFail(t, "class A {static async set foo(value) { }}", "Unexpected token `identifier` at (1:26)", nil)
}

func TestAsyncAwaitFail47(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_ASYNC_GENERATOR)
	testFail(t, "class A {static async* foo() { }}", "Unexpected token `*` at (1:21)", opts)
}

func TestAsyncAwaitFail48(t *testing.T) {
	testFail(t, "class A {async foo() { var await }}",
		"Invalid binding `await` at (1:27)", nil)
}

func TestAsyncAwaitFail49(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "class A {async foo(await) { }}",
		"Invalid binding `await` at (1:19)", opts)
}

func TestAsyncAwaitFail50(t *testing.T) {
	testFail(t, "class A {async foo() { return {await} }}",
		"Unexpected token `await` at (1:31)", nil)
}

func TestAsyncAwaitFail51(t *testing.T) {
	testFail(t, "await", "Unexpected token `EOF` at (1:5)", nil)
}

func TestAsyncAwaitFail52(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "await a", "Unexpected token at (1:6)", opts)
}

func TestAsyncAwaitFail53(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "await a",
		"Cannot use keyword 'await' outside an async function at (1:0)", opts)
}

func TestAsyncAwaitFail54(t *testing.T) {
	testFail(t, "async function foo() { await }",
		"Unexpected token `}` at (1:29)", nil)
}

func TestAsyncAwaitFail55(t *testing.T) {
	testFail(t, "(async function foo() { await })",
		"Unexpected token `}` at (1:30)", nil)
}

func TestAsyncAwaitFail56(t *testing.T) {
	testFail(t, "async () => await", "Unexpected token `EOF` at (1:17)", nil)
}

func TestAsyncAwaitFail57(t *testing.T) {
	testFail(t, "({async foo() { await }})", "Unexpected token `}` at (1:22)", nil)
}

func TestAsyncAwaitFail58(t *testing.T) {
	testFail(t, "(class {async foo() { await }})", "Unexpected token `}` at (1:28)", nil)
}

func TestAsyncAwaitFail59(t *testing.T) {
	testFail(t, "async function foo(a = await b) {}",
		"Await expression can't be used in parameter at (1:23)", nil)
}

func TestAsyncAwaitFail60(t *testing.T) {
	testFail(t, "(async function foo(a = await b) {})",
		"Await expression can't be used in parameter at (1:24)", nil)
}

func TestAsyncAwaitFail61(t *testing.T) {
	testFail(t, "async (a = await b) => {}",
		"Await expression can't be used in parameter at (1:11)", nil)
}

func TestAsyncAwaitFail62(t *testing.T) {
	testFail(t, "for ( ; false; ) async function* g() {}",
		"function declarations can't appear in single-statement context at (1:17)", nil)
}

func TestAsyncAwaitFail63(t *testing.T) {
	testFail(t, "async function wrapper() {\nasync (a = await b) => {}\n}",
		"Await expression can't be used in parameter at (2:11)", nil)
}

func TestAsyncAwaitFail64(t *testing.T) {
	testFail(t, "({async foo(a = await b) {}})",
		"Await expression cannot be a default value at (1:16)", nil)
}

func TestAsyncAwaitFail65(t *testing.T) {
	testFail(t, "(class {async foo(a = await b) {}})",
		"Await expression cannot be a default value at (1:22)", nil)
}

func TestAsyncAwaitFail66(t *testing.T) {
	testFail(t, "async function foo(a = class extends (await b) {}) {}",
		"Await expression can't be used in parameter at (1:38)", nil)
}

func TestAsyncAwaitFail67(t *testing.T) {
	testFail(t, "async function wrap() {\n(a = await b) => a\n}",
		"Await expression cannot be a default value at (2:5)", nil)
}

func TestAsyncAwaitFail68(t *testing.T) {
	testFail(t, "async function wrap() {\n({a = await b} = obj) => a\n}",
		"Await expression cannot be a default value at (2:6)", nil)
}

func TestAsyncAwaitFail69(t *testing.T) {
	testFail(t, "function* wrap() {\nasync(a = yield b) => a\n}",
		"Unexpected token `yield` at (2:10)", nil)
}

func TestAsyncAwaitFail70(t *testing.T) {
	testFail(t, "if (x) async function f() {}",
		"function declarations can't appear in single-statement context at (1:7)", nil)
}

func TestAsyncAwaitFail71(t *testing.T) {
	testFail(t, "(async)(a) => 12", "Unexpected token at (1:11)", nil)
}

func TestAsyncAwaitFail72(t *testing.T) {
	testFail(t, "f = async ((x)) => x", "Invalid parenthesized assignment pattern at (1:11)", nil)
}

func TestAsyncAwaitFail73(t *testing.T) {
	testFail(t, "abc: async function a() {}",
		"function declarations can't appear in single-statement context at (1:5)", nil)
}

func TestAsyncAwaitFail74(t *testing.T) {
	testFail(t, "(async() => { await 4 ** 2 })()",
		"Unary operator used immediately before exponentiation expression at (1:20)", nil)
}

func TestAsyncAwaitFail75(t *testing.T) {
	testFail(t, "async() => (await 1 ** 3)",
		"Unary operator used immediately before exponentiation expression at (1:18)", nil)
}

func TestAsyncAwaitFail76(t *testing.T) {
	testFail(t, "async() => await 5 ** 6",
		"Unary operator used immediately before exponentiation expression at (1:17)", nil)
}

func TestAsyncAwaitFail77(t *testing.T) {
	testFail(t, "async() => await (5) ** 6",
		"Unary operator used immediately before exponentiation expression at (1:18)", nil)
}

func TestAsyncAwaitFail78(t *testing.T) {
	testFail(t, "4 + async() => 2", "Unexpected token `=>` at (1:12)", nil)
}

func TestAsyncAwaitFail79(t *testing.T) {
	testFail(t, "async function𝐬 f() {}", "Unexpected token `identifier` at (1:16)", nil)
}

func TestAsyncAwaitFail80(t *testing.T) {
	testFail(t, "console.log( -2 ** 4 )",
		"Unary operator used immediately before exponentiation expression at (1:14)", nil)
}

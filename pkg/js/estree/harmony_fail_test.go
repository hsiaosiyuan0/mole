package estree

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

func compileProg(code string, opts *parser.ParserOpts) (Node, error) {
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
	testFail(t, "export var await", "Unexpected strict mode reserved word at (1:11)", nil)
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
	testFail(t, "export { encrypt }; if (true) function encrypt() {}", "In strict mode code, functions can only be declared at top level or inside a block at (1:30)", nil)
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
	testFail(t, "() => { class a extends b { static get prototype(){} } }", "Classes may not have a static property named `prototype` at (1:39)", nil)
}

func TestHarmonyFail31(t *testing.T) {
	testFail(t, "class a extends b { static set prototype(a){} }", "Classes may not have a static property named `prototype` at (1:31)", nil)
}

func TestHarmonyFail32(t *testing.T) {
	testFail(t, "class a { static prototype(a){} }", "Classes may not have a static property named `prototype` at (1:17)", nil)
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

func TestHarmonyFail43(t *testing.T) {}

func TestHarmonyFail44(t *testing.T) {}

func TestHarmonyFail45(t *testing.T) {}

func TestHarmonyFail46(t *testing.T) {}

func TestHarmonyFail47(t *testing.T) {}

func TestHarmonyFail48(t *testing.T) {}

func TestHarmonyFail49(t *testing.T) {}

func TestHarmonyFail50(t *testing.T) {}

func TestHarmonyFail51(t *testing.T) {}

func TestHarmonyFail52(t *testing.T) {}

func TestHarmonyFail53(t *testing.T) {}

func TestHarmonyFail54(t *testing.T) {}

func TestHarmonyFail55(t *testing.T) {}

func TestHarmonyFail56(t *testing.T) {}

func TestHarmonyFail57(t *testing.T) {}

func TestHarmonyFail58(t *testing.T) {}

func TestHarmonyFail59(t *testing.T) {}

func TestHarmonyFail60(t *testing.T) {}

func TestHarmonyFail61(t *testing.T) {}

func TestHarmonyFail62(t *testing.T) {}

func TestHarmonyFail63(t *testing.T) {}

func TestHarmonyFail64(t *testing.T) {}

func TestHarmonyFail65(t *testing.T) {}

func TestHarmonyFail66(t *testing.T) {}

func TestHarmonyFail67(t *testing.T) {}

func TestHarmonyFail68(t *testing.T) {}

func TestHarmonyFail69(t *testing.T) {}

func TestHarmonyFail70(t *testing.T) {}

func TestHarmonyFail71(t *testing.T) {}

func TestHarmonyFail72(t *testing.T) {}

func TestHarmonyFail73(t *testing.T) {}

func TestHarmonyFail74(t *testing.T) {}

func TestHarmonyFail75(t *testing.T) {}

func TestHarmonyFail76(t *testing.T) {}

func TestHarmonyFail77(t *testing.T) {}

func TestHarmonyFail78(t *testing.T) {}

func TestHarmonyFail79(t *testing.T) {}

func TestHarmonyFail80(t *testing.T) {}

func TestHarmonyFail81(t *testing.T) {}

func TestHarmonyFail82(t *testing.T) {}

func TestHarmonyFail83(t *testing.T) {}

func TestHarmonyFail84(t *testing.T) {}

func TestHarmonyFail85(t *testing.T) {}

func TestHarmonyFail86(t *testing.T) {}

func TestHarmonyFail87(t *testing.T) {}

func TestHarmonyFail88(t *testing.T) {}

func TestHarmonyFail89(t *testing.T) {}

func TestHarmonyFail90(t *testing.T) {}

func TestHarmonyFail91(t *testing.T) {}

func TestHarmonyFail92(t *testing.T) {}

func TestHarmonyFail93(t *testing.T) {}

func TestHarmonyFail94(t *testing.T) {}

func TestHarmonyFail95(t *testing.T) {}

func TestHarmonyFail96(t *testing.T) {}

func TestHarmonyFail97(t *testing.T) {}

func TestHarmonyFail98(t *testing.T) {}

func TestHarmonyFail99(t *testing.T) {}

func TestHarmonyFail100(t *testing.T) {}

func TestHarmonyFail101(t *testing.T) {}

func TestHarmonyFail102(t *testing.T) {}

func TestHarmonyFail103(t *testing.T) {}

func TestHarmonyFail104(t *testing.T) {}

func TestHarmonyFail105(t *testing.T) {}

func TestHarmonyFail106(t *testing.T) {}

func TestHarmonyFail107(t *testing.T) {}

func TestHarmonyFail108(t *testing.T) {}

func TestHarmonyFail109(t *testing.T) {}

func TestHarmonyFail110(t *testing.T) {}

func TestHarmonyFail111(t *testing.T) {}

func TestHarmonyFail112(t *testing.T) {}

func TestHarmonyFail113(t *testing.T) {}

func TestHarmonyFail114(t *testing.T) {}

func TestHarmonyFail115(t *testing.T) {}

func TestHarmonyFail116(t *testing.T) {}

func TestHarmonyFail117(t *testing.T) {}

func TestHarmonyFail118(t *testing.T) {}

func TestHarmonyFail119(t *testing.T) {}

func TestHarmonyFail120(t *testing.T) {}

func TestHarmonyFail121(t *testing.T) {}

func TestHarmonyFail122(t *testing.T) {}

func TestHarmonyFail123(t *testing.T) {}

func TestHarmonyFail124(t *testing.T) {}

func TestHarmonyFail125(t *testing.T) {}

func TestHarmonyFail126(t *testing.T) {}

func TestHarmonyFail127(t *testing.T) {}

func TestHarmonyFail128(t *testing.T) {}

func TestHarmonyFail129(t *testing.T) {}

func TestHarmonyFail130(t *testing.T) {}

func TestHarmonyFail131(t *testing.T) {}

func TestHarmonyFail132(t *testing.T) {}

func TestHarmonyFail133(t *testing.T) {}

func TestHarmonyFail134(t *testing.T) {}

func TestHarmonyFail135(t *testing.T) {}

func TestHarmonyFail136(t *testing.T) {}

func TestHarmonyFail137(t *testing.T) {}

func TestHarmonyFail138(t *testing.T) {}

func TestHarmonyFail139(t *testing.T) {}

func TestHarmonyFail140(t *testing.T) {}

func TestHarmonyFail141(t *testing.T) {}

func TestHarmonyFail142(t *testing.T) {}

func TestHarmonyFail143(t *testing.T) {}

func TestHarmonyFail144(t *testing.T) {}

func TestHarmonyFail145(t *testing.T) {}

func TestHarmonyFail146(t *testing.T) {}

func TestHarmonyFail147(t *testing.T) {}

func TestHarmonyFail148(t *testing.T) {}

func TestHarmonyFail149(t *testing.T) {}

func TestHarmonyFail150(t *testing.T) {}

func TestHarmonyFail151(t *testing.T) {}

func TestHarmonyFail152(t *testing.T) {}

func TestHarmonyFail153(t *testing.T) {}

func TestHarmonyFail154(t *testing.T) {}

func TestHarmonyFail155(t *testing.T) {}

func TestHarmonyFail156(t *testing.T) {}

func TestHarmonyFail157(t *testing.T) {}

func TestHarmonyFail158(t *testing.T) {}

func TestHarmonyFail159(t *testing.T) {}

func TestHarmonyFail160(t *testing.T) {}

func TestHarmonyFail161(t *testing.T) {}

func TestHarmonyFail162(t *testing.T) {}

func TestHarmonyFail163(t *testing.T) {}

func TestHarmonyFail164(t *testing.T) {}

func TestHarmonyFail165(t *testing.T) {}

func TestHarmonyFail166(t *testing.T) {}

func TestHarmonyFail167(t *testing.T) {}

func TestHarmonyFail168(t *testing.T) {}

func TestHarmonyFail169(t *testing.T) {}

func TestHarmonyFail170(t *testing.T) {}

func TestHarmonyFail171(t *testing.T) {}

func TestHarmonyFail172(t *testing.T) {}

func TestHarmonyFail173(t *testing.T) {}

func TestHarmonyFail174(t *testing.T) {}

func TestHarmonyFail175(t *testing.T) {}

func TestHarmonyFail176(t *testing.T) {}

func TestHarmonyFail178(t *testing.T) {}

func TestHarmonyFail179(t *testing.T) {}

func TestHarmonyFail180(t *testing.T) {}

func TestHarmonyFail181(t *testing.T) {}

func TestHarmonyFail182(t *testing.T) {}

func TestHarmonyFail183(t *testing.T) {}

func TestHarmonyFail184(t *testing.T) {}

func TestHarmonyFail185(t *testing.T) {}

func TestHarmonyFail186(t *testing.T) {}

func TestHarmonyFail187(t *testing.T) {}

func TestHarmonyFail188(t *testing.T) {}

func TestHarmonyFail189(t *testing.T) {}

func TestHarmonyFail190(t *testing.T) {}

func TestHarmonyFail191(t *testing.T) {}

func TestHarmonyFail192(t *testing.T) {}

func TestHarmonyFail193(t *testing.T) {}

func TestHarmonyFail194(t *testing.T) {}

func TestHarmonyFail195(t *testing.T) {}

func TestHarmonyFail196(t *testing.T) {}

func TestHarmonyFail197(t *testing.T) {}

func TestHarmonyFail198(t *testing.T) {}

func TestHarmonyFail199(t *testing.T) {}

func TestHarmonyFail200(t *testing.T) {}

func TestHarmonyFail201(t *testing.T) {}

func TestHarmonyFail202(t *testing.T) {}

func TestHarmonyFail203(t *testing.T) {}

func TestHarmonyFail204(t *testing.T) {}

func TestHarmonyFail205(t *testing.T) {}

func TestHarmonyFail206(t *testing.T) {}

func TestHarmonyFail207(t *testing.T) {}

func TestHarmonyFail208(t *testing.T) {}

func TestHarmonyFail209(t *testing.T) {}

func TestHarmonyFail210(t *testing.T) {}

func TestHarmonyFail211(t *testing.T) {}

func TestHarmonyFail212(t *testing.T) {}

func TestHarmonyFail213(t *testing.T) {}

func TestHarmonyFail214(t *testing.T) {}

func TestHarmonyFail215(t *testing.T) {}

func TestHarmonyFail216(t *testing.T) {}

func TestHarmonyFail217(t *testing.T) {}

func TestHarmonyFail218(t *testing.T) {}

func TestHarmonyFail219(t *testing.T) {}

func TestHarmonyFail220(t *testing.T) {}

func TestHarmonyFail221(t *testing.T) {}

func TestHarmonyFail222(t *testing.T) {}

func TestHarmonyFail223(t *testing.T) {}

func TestHarmonyFail224(t *testing.T) {}

func TestHarmonyFail225(t *testing.T) {}

func TestHarmonyFail226(t *testing.T) {}

func TestHarmonyFail227(t *testing.T) {}

func TestHarmonyFail228(t *testing.T) {}

func TestHarmonyFail229(t *testing.T) {}

func TestHarmonyFail230(t *testing.T) {}

func TestHarmonyFail231(t *testing.T) {}

func TestHarmonyFail232(t *testing.T) {}

func TestHarmonyFail233(t *testing.T) {}

func TestHarmonyFail234(t *testing.T) {}

func TestHarmonyFail235(t *testing.T) {}

func TestHarmonyFail236(t *testing.T) {}

func TestHarmonyFail237(t *testing.T) {}

func TestHarmonyFail238(t *testing.T) {}

func TestHarmonyFail239(t *testing.T) {}

func TestHarmonyFail240(t *testing.T) {}

func TestHarmonyFail241(t *testing.T) {}

func TestHarmonyFail242(t *testing.T) {}

func TestHarmonyFail243(t *testing.T) {}

func TestHarmonyFail244(t *testing.T) {}

func TestHarmonyFail245(t *testing.T) {}

func TestHarmonyFail246(t *testing.T) {}

func TestHarmonyFail247(t *testing.T) {}

func TestHarmonyFail248(t *testing.T) {}

func TestHarmonyFail249(t *testing.T) {}

func TestHarmonyFail250(t *testing.T) {}

func TestHarmonyFail251(t *testing.T) {}

func TestHarmonyFail252(t *testing.T) {}

func TestHarmonyFail253(t *testing.T) {}

func TestHarmonyFail254(t *testing.T) {}

func TestHarmonyFail255(t *testing.T) {}

func TestHarmonyFail256(t *testing.T) {}

func TestHarmonyFail257(t *testing.T) {}

func TestHarmonyFail258(t *testing.T) {}

func TestHarmonyFail259(t *testing.T) {}

func TestHarmonyFail260(t *testing.T) {}

func TestHarmonyFail261(t *testing.T) {}

func TestHarmonyFail262(t *testing.T) {}

func TestHarmonyFail263(t *testing.T) {}

func TestHarmonyFail264(t *testing.T) {}

func TestHarmonyFail265(t *testing.T) {}

func TestHarmonyFail266(t *testing.T) {}

func TestHarmonyFail267(t *testing.T) {}

func TestHarmonyFail268(t *testing.T) {}

func TestHarmonyFail269(t *testing.T) {}

func TestHarmonyFail270(t *testing.T) {}

func TestHarmonyFail271(t *testing.T) {}

func TestHarmonyFail272(t *testing.T) {}

func TestHarmonyFail273(t *testing.T) {}

func TestHarmonyFail274(t *testing.T) {}

func TestHarmonyFail275(t *testing.T) {}

func TestHarmonyFail276(t *testing.T) {}

func TestHarmonyFail278(t *testing.T) {}

func TestHarmonyFail279(t *testing.T) {}

func TestHarmonyFail280(t *testing.T) {}

func TestHarmonyFail281(t *testing.T) {}

func TestHarmonyFail282(t *testing.T) {}

func TestHarmonyFail283(t *testing.T) {}

func TestHarmonyFail284(t *testing.T) {}

func TestHarmonyFail285(t *testing.T) {}

func TestHarmonyFail286(t *testing.T) {}

func TestHarmonyFail287(t *testing.T) {}

func TestHarmonyFail288(t *testing.T) {}

func TestHarmonyFail289(t *testing.T) {}

func TestHarmonyFail290(t *testing.T) {}

func TestHarmonyFail291(t *testing.T) {}

func TestHarmonyFail292(t *testing.T) {}

func TestHarmonyFail293(t *testing.T) {}

func TestHarmonyFail294(t *testing.T) {}

func TestHarmonyFail295(t *testing.T) {}

func TestHarmonyFail296(t *testing.T) {}

func TestHarmonyFail297(t *testing.T) {}

func TestHarmonyFail298(t *testing.T) {}

func TestHarmonyFail299(t *testing.T) {}

func TestHarmonyFail300(t *testing.T) {}

package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestClassFeatureFail1(t *testing.T) {
	testFail(t, "async function f() { class C { aaa = await } }",
		"Cannot use keyword 'await' outside an async function at (1:37)", nil)
}

func TestClassFeatureFail2(t *testing.T) {
	testFail(t, "function* f() { class C { aaa = yield } }",
		"Unexpected token `yield` at (1:32)", nil)
}

func TestClassFeatureFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_CLASS_PUB_FIELD)
	testFail(t, "class C { aaa }", "Unexpected token `}` at (1:14)", opts)
}

func TestClassFeatureFail4(t *testing.T) {
	testFail(t, "class C { [super.bbb] = 0 }",
		"'super' is only allowed in object methods and classes at (1:11)", nil)
}

func TestClassFeatureFail5(t *testing.T) {
	testFail(t, "class C { aaa, bbb }", "Unexpected token `,` at (1:13)", nil)
}

func TestClassFeatureFail6(t *testing.T) {
	testFail(t, "class C { get aaa }", "Unexpected token `}` at (1:18)", nil)
}

func TestClassFeatureFail7(t *testing.T) {
	testFail(t, "class C { set aaa }", "Unexpected token `}` at (1:18)", nil)
}

func TestClassFeatureFail8(t *testing.T) {
	testFail(t, "class C { *aaa }", "Unexpected token `}` at (1:15)", nil)
}

func TestClassFeatureFail9(t *testing.T) {
	testFail(t, "class C { async aaa }", "Unexpected token `}` at (1:20)", nil)
}

func TestClassFeatureFail10(t *testing.T) {
	testFail(t, "class C { async*aaa }", "Unexpected token `}` at (1:20)", nil)
}

func TestClassFeatureFail11(t *testing.T) {
	testFail(t, "class C { aaa bbb }", "Unexpected token `identifier` at (1:14)", nil)
}

func TestClassFeatureFail12(t *testing.T) {
	testFail(t, "class C { aaa = 0, 1 }", "Unexpected token `,` at (1:17)", nil)
}

func TestClassFeatureFail13(t *testing.T) {
	testFail(t, "class C { constructor }",
		"Classes can't have a field named `constructor` at (1:10)", nil)
}

func TestClassFeatureFail14(t *testing.T) {
	testFail(t, "class C { static constructor }",
		"Classes can't have a field named `constructor` at (1:17)", nil)
}

func TestClassFeatureFail15(t *testing.T) {
	testFail(t, "class C { static prototype }",
		"Classes can't have a static field named `prototype` at (1:17)", nil)
}

func TestClassFeatureFail16(t *testing.T) {
	testFail(t, "class C { aaa = arguments }",
		"Unexpected token `arguments` at (1:16)", nil)
}

func TestClassFeatureFail17(t *testing.T) {
	testFail(t, "class C { aaa = { arguments } }",
		"Unexpected token `arguments` at (1:18)", nil)
}

func TestClassFeatureFail18(t *testing.T) {
	testFail(t, "class C { aaa = { arguments: { arguments } } }",
		"Unexpected token `arguments` at (1:31)", nil)
}

func TestClassFeatureFail19(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_CLASS_PRV)
	testFail(t, "class C { #aaa }", "Unexpected character at (1:10)", opts)
}

func TestClassFeatureFail20(t *testing.T) {
	testFail(t, "class C { # aaa }", "Unexpected character at (1:11)", nil)
}

func TestClassFeatureFail21(t *testing.T) {
	testFail(t, "class C { #+aaa }", "Unexpected character at (1:11)", nil)
}

func TestClassFeatureFail22(t *testing.T) {
	testFail(t, "class C { #👍 }", "Unexpected character at (1:11)", nil)
}

func TestClassFeatureFail23(t *testing.T) {
	testFail(t, "class C { #constructor }",
		"Classes can't have a field named `constructor` at (1:10)", nil)
}

func TestClassFeatureFail24(t *testing.T) {
	testFail(t, "class C { static #constructor }",
		"Classes can't have a field named `constructor` at (1:17)", nil)
}

func TestClassFeatureFail25(t *testing.T) {
	testFail(t, "class C { #a; #a }",
		"Identifier `#a` has already been declared at (1:14)", nil)
}

func TestClassFeatureFail26(t *testing.T) {
	testFail(t, "class C { #a; static #a }",
		"Identifier `#a` has already been declared at (1:21)", nil)
}

func TestClassFeatureFail27(t *testing.T) {
	testFail(t, "for ( ; false; ) async function* g() {}",
		"function declarations can't appear in single-statement context at (1:17)", nil)
}

func TestClassFeatureFail28(t *testing.T) {
	testFail(t, "class C { static #a; #a }",
		"Identifier `#a` has already been declared at (1:21)", nil)
}

func TestClassFeatureFail29(t *testing.T) {
	testFail(t, "class C { #a(){} #a }",
		"Identifier `#a` has already been declared at (1:17)", nil)
}

func TestClassFeatureFail30(t *testing.T) {
	testFail(t, "class C { #a; #a(){} }",
		"Identifier `#a` has already been declared at (1:14)", nil)
}

func TestClassFeatureFail31(t *testing.T) {
	testFail(t, "class C { #a(){} #a(){} }",
		"Identifier `#a` has already been declared at (1:17)", nil)
}

func TestClassFeatureFail32(t *testing.T) {
	testFail(t, "class C { #a(){} static #a(){} }",
		"Identifier `#a` has already been declared at (1:24)", nil)
}

func TestClassFeatureFail33(t *testing.T) {
	testFail(t, "class C { static #a(){} #a(){} }",
		"Identifier `#a` has already been declared at (1:24)", nil)
}

func TestClassFeatureFail34(t *testing.T) {
	testFail(t, "class C { get #a(){} static set #a(x){} }",
		"Identifier `#a` has already been declared at (1:32)", nil)
}

func TestClassFeatureFail35(t *testing.T) {
	testFail(t, "class C { set #a(x){} static get #a(){} }",
		"Identifier `#a` has already been declared at (1:33)", nil)
}

func TestClassFeatureFail36(t *testing.T) {
	testFail(t, "class C { static get #a(){} set #a(x){} }",
		"Identifier `#a` has already been declared at (1:32)", nil)
}

func TestClassFeatureFail37(t *testing.T) {
	testFail(t, "class C { static set #a(x){} get #a(){} }",
		"Identifier `#a` has already been declared at (1:33)", nil)
}

func TestClassFeatureFail38(t *testing.T) {
	testFail(t, "class C { #a; get #a(){} }",
		"Identifier `#a` has already been declared at (1:18)", nil)
}

func TestClassFeatureFail39(t *testing.T) {
	testFail(t, "class C { #a; set #a(x){} }",
		"Identifier `#a` has already been declared at (1:18)", nil)
}

func TestClassFeatureFail40(t *testing.T) {
	testFail(t, "class C { #a(){}; get #a(){} }",
		"Identifier `#a` has already been declared at (1:22)", nil)
}

func TestClassFeatureFail41(t *testing.T) {
	testFail(t, "class C { get #a(){} #a }",
		"Identifier `#a` has already been declared at (1:21)", nil)
}

func TestClassFeatureFail42(t *testing.T) {
	testFail(t, "class C { set #a(x){} #a }",
		"Identifier `#a` has already been declared at (1:22)", nil)
}

func TestClassFeatureFail43(t *testing.T) {
	testFail(t, "class C { get #a(){} #a(){} }",
		"Identifier `#a` has already been declared at (1:21)", nil)
}

func TestClassFeatureFail44(t *testing.T) {
	testFail(t, "class C { set #a(x){} #a(){} }",
		"Identifier `#a` has already been declared at (1:22)", nil)
}

func TestClassFeatureFail45(t *testing.T) {
	testFail(t, "class C { #a(){}; set #a(x){} }",
		"Identifier `#a` has already been declared at (1:22)", nil)
}

func TestClassFeatureFail46(t *testing.T) {
	testFail(t, "class C extends Base { f() { return super.#aaa } }",
		"Unexpected private field at (1:42)", nil)
}

func TestClassFeatureFail47(t *testing.T) {
	testFail(t, "class C { #aaa; f() { delete this.#aaa } }",
		"Private fields can not be deleted at (1:34)", nil)
}

func TestClassFeatureFail48(t *testing.T) {
	testFail(t, "class C { #aaa; f() { delete obj?.#aaa } }",
		"Private fields can not be deleted at (1:34)", nil)
}

func TestClassFeatureFail49(t *testing.T) {
	testFail(t, "class C { #aaa; f() { delete obj?.p.#aaa } }",
		"Private fields can not be deleted at (1:36)", nil)
}

func TestClassFeatureFail50(t *testing.T) {
	testFail(t, "const obj = #aaa", "Unexpected token `private identifier` at (1:12)", nil)
}

func TestClassFeatureFail51(t *testing.T) {
	testFail(t, "const obj = { #aaa }", "Unexpected token `private identifier` at (1:14)", nil)
}

func TestClassFeatureFail52(t *testing.T) {
	testFail(t, "class C { #aaa; f() { #aaa } }",
		"Unexpected token `private identifier` at (1:22)", nil)
}

func TestClassFeatureFail53(t *testing.T) {
	testFail(t, "class C { #aaa; f() { return { #aaa: 1 } } }",
		"Unexpected token `private identifier` at (1:31)", nil)
}

func TestClassFeatureFail54(t *testing.T) {
	testFail(t, "class C { #a; a = this.#b; }",
		"Private field `#b` must be declared in an enclosing class at (1:23)", nil)
}

func TestClassFeatureFail55(t *testing.T) {
	testFail(t, "class C { a = this.#b; #a; }",
		"Private field `#b` must be declared in an enclosing class at (1:19)", nil)
}

func TestClassFeatureFail56(t *testing.T) {
	testFail(t, "class C { #a; [this.#b]; }",
		"Private field `#b` must be declared in an enclosing class at (1:20)", nil)
}

func TestClassFeatureFail57(t *testing.T) {
	testFail(t, "class C { [this.#b]; #a; }",
		"Private field `#b` must be declared in an enclosing class at (1:16)", nil)
}

func TestClassFeatureFail58(t *testing.T) {
	testFail(t, "class C { #a; f(){ this.#b } }",
		"Private field `#b` must be declared in an enclosing class at (1:24)", nil)
}

func TestClassFeatureFail59(t *testing.T) {
	testFail(t, "class C { f(){ this.#b } #a; }",
		"Private field `#b` must be declared in an enclosing class at (1:20)", nil)
}

func TestClassFeatureFail60(t *testing.T) {
	testFail(t, "obj.#aaa",
		"Private field `#aaa` must be declared in an enclosing class at (1:4)", nil)
}

func TestClassFeatureFail61(t *testing.T) {
	testFail(t, "function F() { obj.#aaa }",
		"Private field `#aaa` must be declared in an enclosing class at (1:19)", nil)
}

func TestClassFeatureFail62(t *testing.T) {
	testFail(t, "class Outer { Inner = class { f(obj) { obj.#nonexist } #inner; }; #outer; }",
		"Private field `#nonexist` must be declared in an enclosing class at (1:43)", nil)
}

func TestClassFeatureFail63(t *testing.T) {
	testFail(t, `class C {
    static #z
    static {
      this.#y = {}
    }
  }`, "Private field `#y` must be declared in an enclosing class at (4:11)", nil)
}

func TestClassFeatureFail64(t *testing.T) {
	testFail(t, `let zRead
  class C {
    static #z
    static {
      zRead = () => this.#y
    }
  }`, "Private field `#y` must be declared in an enclosing class at (5:25)", nil)
}

func TestClassFeatureFail65(t *testing.T) {
	testFail(t, "class C { 'constructor' }",
		"Classes can't have a field named `constructor` at (1:10)", nil)
}

func TestClassFeatureFail66(t *testing.T) {
	testFail(t, "class A { [#c]() {} }",
		"Unexpected token `private identifier` at (1:11)", nil)
}

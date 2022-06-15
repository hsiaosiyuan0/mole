package linter

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestGetterReturn1(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() {} };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn2(t *testing.T) {
	r := lint(t, "var foo = { get\n bar () {} };", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn3(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() { if (baz) { return true; } } };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn25(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() { if (baz) { return true; } else { return; } } };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

func TestGetterReturn26(t *testing.T) {
	r := lintCb(t, `
  var foo = { get bar() { if (baz) { return true; } else { return; } } };
  `, func(c *Config) {
		c.Rules["getter-return"] = []interface{}{
			2,
			map[string]interface{}{
				"allowImplicit": true,
			},
		}
	})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestGetterReturn4(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() { ~function () { return true; } } };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn5(t *testing.T) {
	r := lintCb(t, `
  var foo = { get bar() { return; } };
  `, func(c *Config) {
		c.Rules["getter-return"] = []interface{}{
			2,
			map[string]interface{}{
				"allowImplicit": true,
			},
		}
	}, &GetterReturn{})

	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestGetterReturn6(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() { return; } };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn7(t *testing.T) {
	r := lint(t, `
  var foo = { get bar() { if (baz) { return; } } };
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn8(t *testing.T) {
	r := lint(t, `
  class foo { get bar(){} }
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn9(t *testing.T) {
	r := lint(t, "var foo = class {\n  static get\nbar(){} }", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn10(t *testing.T) {
	r := lintCb(t, `
  class foo { get bar(){ return; }}
  `, func(c *Config) {
		c.Rules["getter-return"] = []interface{}{
			2,
			map[string]interface{}{
				"allowImplicit": true,
			},
		}
	}, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestGetterReturn11(t *testing.T) {
	r := lint(t, `
  class foo { get bar(){ ~function () { return true; }()}}
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn12(t *testing.T) {
	r := lint(t, `
  class foo { get bar(){if (baz) {return true;} } }
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn13(t *testing.T) {
	r := lint(t, `
  Object.defineProperty(foo, 'bar', { get: function (){}});
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn14(t *testing.T) {
	r := lint(t, `
  Object.defineProperty(foo, 'bar', { get: function getFoo (){}});
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn15(t *testing.T) {
	r := lint(t, `
  Object.defineProperty(foo, 'bar', { get(){} });
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn16(t *testing.T) {
	r := lint(t, `
  Object.defineProperty(foo, 'bar', { get: () => {}});
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn17(t *testing.T) {
	r := lint(t, "Object.defineProperty(foo, \"bar\", { get: function (){if(bar) {return true;}}});", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn18(t *testing.T) {
	r := lint(t, "Object.defineProperty(foo, \"bar\", { get: function (){ ~function () { return true; }()}});", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn19(t *testing.T) {
	r := lint(t, `
  Object.defineProperties(foo, { bar: { get: function (){if(bar) {return true;}}}});
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn20(t *testing.T) {
	r := lint(t, `
  Object.defineProperties(foo, { bar: { get: function () {~function () { return true; }()}} });
  `, &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn21(t *testing.T) {
	r := lint(t, "Object.defineProperty(foo, \"bar\", { get: function (){}});", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn22(t *testing.T) {
	r := lint(t, "(Object?.defineProperty)(foo, 'bar', { get: function (){} });", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn23(t *testing.T) {
	r := lint(t, "Object?.defineProperty(foo, 'bar', { get: function (){} });", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

func TestGetterReturn24(t *testing.T) {
	r := lint(t, "(Object?.defineProperty)(foo, 'bar', { get: function (){ return } });", &GetterReturn{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "expected to return a value", r.Diagnoses[0].Msg, "should be ok")
}

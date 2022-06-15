package linter

import (
	"plugin"
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func newCfg() *Config {
	cfg := &Config{
		Rules: map[string]interface{}{},
	}

	cfg.cwd = ""
	cfg.plugins = map[string]*plugin.Plugin{}

	cfg.ruleFacts = map[string]RuleFact{}
	cfg.rulesCfg = map[string]*RuleCfg{}
	cfg.ruleFactsLang = map[string]map[string]RuleFact{}

	cfg.IgnorePatterns = []string{}
	return cfg
}

func lintCb(t *testing.T, code string, cb func(*Config), rfs ...RuleFact) *Reports {
	cfg := newCfg()

	for _, rf := range rfs {
		cfg.AddRuleFact(rf)
	}

	cb(cfg)

	linter, err := NewLinter("", cfg, true)
	if err != nil {
		t.Fatal(err)
	}

	u, err := NewJsUnit("test.js", code, linter.cfg)
	if err != nil {
		t.Fatal(err)
	}

	u.linter = linter
	u.initRules().enableAllRules(false)
	u.ana.Analyze()

	return linter.mrkReports()
}

func lint(t *testing.T, code string, rfs ...RuleFact) *Reports {
	cfg := newCfg()

	for _, rf := range rfs {
		cfg.AddRuleFact(rf)
	}

	linter, err := NewLinter("", cfg, true)
	if err != nil {
		t.Fatal(err)
	}

	u, err := NewJsUnit("test.js", code, linter.cfg)
	if err != nil {
		t.Fatal(err)
	}

	u.linter = linter
	u.initRules().enableAllRules(false)
	u.ana.Analyze()

	return linter.mrkReports()
}

func TestUnreachable1(t *testing.T) {
	r := lint(t, `function foo() { return x; var x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable2(t *testing.T) {
	r := lint(t, `function foo() { return x; var x, y = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable3(t *testing.T) {
	r := lint(t, `while (true) { continue; var x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable4(t *testing.T) {
	r := lint(t, `function foo() { return; x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable5(t *testing.T) {
	r := lint(t, `function foo() { throw error; x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable6(t *testing.T) {
	r := lint(t, `while (true) { break; x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable7(t *testing.T) {
	r := lint(t, `while (true) { continue; x = 1; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable8(t *testing.T) {
	r := lint(t, `function foo() { switch (foo) { case 1: return; x = 1; } }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable9(t *testing.T) {
	r := lint(t, `function foo() { switch (foo) { case 1: throw e; x = 1; } }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable10(t *testing.T) {
	r := lint(t, `while (true) { switch (foo) { case 1: break; x = 1; } }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable11(t *testing.T) {
	r := lint(t, `while (true) { switch (foo) { case 1: continue; x = 1; } }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable12(t *testing.T) {
	r := lint(t, `var x = 1; throw 'uh oh'; var y = 2;`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable13(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; if (x) { return; } else { throw e; } x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable14(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; if (x) return; else throw -1; x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable15(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; try { return; } finally {} x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable16(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; try { } finally { return; } x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable17(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; do { return; } while (x); x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable18(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; while (x) { if (x) break; else continue; x = 2; } }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable19(t *testing.T) {
	r := lint(t, `function foo() { var x = 1; for (;;) { if (x) continue; } x = 2; }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable21(t *testing.T) {
	r := lint(t, `const arrow_direction = arrow => {  switch (arrow) { default: throw new Error();  }; g() }`, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable22(t *testing.T) {
	r := lint(t, `
function foo() {
  return;

  a();  // ← ERROR: Unreachable code. (no-unreachable)

  b()   // ↑ ';' token is included in the unreachable code, so this statement will be merged.
  // comment
  c();  // ↑ ')' token is included in the unreachable code, so this statement will be merged.
}
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable23(t *testing.T) {
	r := lint(t, `
function foo() {
  return;

  a();

  if (b()) {
      c()
  } else {
      d()
  }
}
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable24(t *testing.T) {
	r := lint(t, `
function foo() {
  if (a) {
      return
      b();
      c();
  } else {
      throw err
      d();
  }
}
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 2, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable25(t *testing.T) {
	r := lint(t, `
function foo() {
  if (a) {
      return
      b();
      c();
  } else {
      throw err
      d();
  }
  e();
}
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 3, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable26(t *testing.T) {
	r := lint(t, `
function foo() {
  try {
      return;
      let a = 1;
  } catch (err) {
      return err;
  }
}
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestUnreachable27(t *testing.T) {
	r := lint(t, `
  while (true) { }
  x = 1;
  `, &NoUnreachable{})
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

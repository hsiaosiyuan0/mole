package lint_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/lint"
	"github.com/hsiaosiyuan0/mole/plugin"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestResolvePlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", "resolve_plugin")
	util.ShellInDir(dir, "npm", "i", "go-cross-ci-demo@0.0.5")

	p, err := plugin.Resolve(dir, "go-cross-ci-demo")
	if err != nil {
		t.Fatal(err)
	}

	v, err := p.Lookup("V")
	if err != nil {
		t.Fatal(err)
	}

	f, err := p.Lookup("F")
	if err != nil {
		t.Fatal(err)
	}

	*v.(*int) = 7
	util.AssertEqualString(t, "Hello, number 7", f.(func() string)(), "should be ok")
}

func TestLoadJsCfg(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	cf := filepath.Join(basepath, "asset", "resolve_plugin", ".eslintrc.js")

	cfg, err := lint.NewConfig(cf, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.Plugins), "should be ok")
}

func TestLoadJsCfgInDir(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", "resolve_plugin")

	cfg, err := lint.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.Plugins), "should be ok")
}

func TestIgnorePattern(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", "resolve_plugin")

	cfg, err := lint.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg.InitIgPatterns()

	util.AssertEqual(t, true, cfg.IsIgnored(filepath.Join(dir, "test.js")), "should be ok")
	util.AssertEqual(t, true, cfg.IsIgnored(filepath.Join(dir, "node_modules", "test1.js")), "should be ok")
	util.AssertEqual(t, false, cfg.IsIgnored(filepath.Join(dir, "node_modules", "test.js")), "should be ok")
	util.AssertEqual(t, false, cfg.IsIgnored(filepath.Join(dir, ".eslintrc.js")), "should be ok")
}

func TestRegisterPlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", "register_plugin")
	util.ShellInDir(dir, "go", "build", "-buildmode=plugin", fmt.Sprintf("-o=node_modules/go-cross-ci-demo/build/go-cross-ci-demo-%s-%s", runtime.GOOS, runtime.GOARCH))

	cfg, err := lint.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = cfg.Init(); err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.RuleFactsLang()[lint.RL_JS]), "should be ok")
}

func TestProcess(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", "register_plugin")
	util.ShellInDir(dir, "go", "build", "-buildmode=plugin", fmt.Sprintf("-o=node_modules/go-cross-ci-demo/build/go-cross-ci-demo-%s-%s", runtime.GOOS, runtime.GOARCH))

	linter, err := lint.NewLinter(dir, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	r := linter.Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow the use of `alert`, `confirm`, and `prompt`", r.Diagnoses[0].Msg, "should be ok")
}

func mkrLinter(t *testing.T, rule string) *lint.Linter {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "asset", rule)

	linter, err := lint.NewLinter(dir, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	return linter
}

func TestNoAlert(t *testing.T) {
	r := mkrLinter(t, "no_alert").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow the use of `alert`, `confirm`, and `prompt`", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoUnreachable(t *testing.T) {
	r := mkrLinter(t, "no_unreachable").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestIgnore(t *testing.T) {
	r := mkrLinter(t, "ignore").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].File(), "test1.js"), "should be ok")
}

func TestIgnoreFile(t *testing.T) {
	r := mkrLinter(t, "ignore_file").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].File(), "test1.js"), "should be ok")
}

func TestIgnoreRoot(t *testing.T) {
	r := mkrLinter(t, "ignore_root").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].File(), "a/test.js"), "should be ok")
}

func TestIgnoreNestOverride(t *testing.T) {
	// the nested `.eslintignore` needs a `.eslintrc.js` to active the nested config resolution
	r := mkrLinter(t, "ignore_nested").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].File(), "a/test.js"), "should be ok")
}

func TestDisableAll(t *testing.T) {
	r := mkrLinter(t, "disable_all").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestEnableAll(t *testing.T) {
	r := mkrLinter(t, "enable_all").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

func TestDisableRules(t *testing.T) {
	r := mkrLinter(t, "disable_rules").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

func TestDisableNextLine(t *testing.T) {
	r := mkrLinter(t, "disable_next_line").Process()
	util.AssertEqual(t, 0, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

type PanicByNum struct{}

func (n *PanicByNum) Name() string {
	return "PanicInCallExpr"
}

func (n *PanicByNum) Meta() *lint.Meta {
	return &lint.Meta{
		Lang: []string{lint.RL_JS},
		Kind: lint.RK_LINT_SEMANTIC,
		Docs: lint.Docs{
			Desc: "",
			Url:  "",
		},
	}
}

func (n *PanicByNum) Options() *plugin.Options {
	return nil
}

func (n *PanicByNum) Validate() *validator.Validate {
	return nil
}

func (n *PanicByNum) Validates() map[int]plugin.Validate {
	return nil
}

func (n *PanicByNum) Create(rc *lint.RuleCtx) map[parser.NodeType]walk.ListenFn {
	return map[parser.NodeType]walk.ListenFn{
		walk.NodeBeforeEvent(parser.N_LIT_NUM): func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			panic("panic by num")
		},
	}
}

// panic in one unit should interrupt the other units's routines
func TestPanicInRule(t *testing.T) {
	lin := mkrLinter(t, "panic_in_rule")
	lin.Config().AddRuleFacts([]lint.RuleFact{&PanicByNum{}})

	r := lin.Process()
	util.AssertEqual(t, 1, len(r.Abnormals), "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

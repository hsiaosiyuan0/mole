package plugin_test

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/linter"
	"github.com/hsiaosiyuan0/mole/plugin"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestResolvePlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "resolve_plugin")
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

	cf := path.Join(basepath, "asset", "resolve_plugin", ".eslintrc.js")

	cfg, err := linter.NewConfig(cf, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.Plugins), "should be ok")
}

func TestLoadJsCfgInDir(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "resolve_plugin")

	cfg, err := linter.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.Plugins), "should be ok")
}

func TestIgnorePattern(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "resolve_plugin")

	cfg, err := linter.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg.InitIgPatterns()

	util.AssertEqual(t, true, cfg.IsIgnored(path.Join(dir, "test.js")), "should be ok")
	util.AssertEqual(t, true, cfg.IsIgnored(path.Join(dir, "node_modules", "test1.js")), "should be ok")
	util.AssertEqual(t, false, cfg.IsIgnored(path.Join(dir, "node_modules", "test.js")), "should be ok")
	util.AssertEqual(t, false, cfg.IsIgnored(path.Join(dir, ".eslintrc.js")), "should be ok")
}

func TestRegisterPlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "register_plugin")
	util.ShellInDir(dir, "go", "build", "-buildmode=plugin", fmt.Sprintf("-o=node_modules/go-cross-ci-demo/build/go-cross-ci-demo-%s-%s", runtime.GOOS, runtime.GOARCH))

	cfg, err := linter.LoadCfgInDir(dir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = cfg.Init(); err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(cfg.RuleFact()[linter.RL_JS]), "should be ok")
}

func TestProcess(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "register_plugin")
	util.ShellInDir(dir, "go", "build", "-buildmode=plugin", fmt.Sprintf("-o=node_modules/go-cross-ci-demo/build/go-cross-ci-demo-%s-%s", runtime.GOOS, runtime.GOARCH))

	linter, err := linter.NewLinter(dir, nil, true)
	if err != nil {
		t.Fatal(err)
	}

	r := linter.Process()
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow the use of `alert`, `confirm`, and `prompt`", r.Diagnoses[0].Msg, "should be ok")
}

func lint(t *testing.T, rule string) *linter.Reports {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", rule)

	linter, err := linter.NewLinter(dir, nil, false)
	if err != nil {
		t.Fatal(err)
	}

	return linter.Process()
}

func TestNoAlert(t *testing.T) {
	r := lint(t, "no_alert")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow the use of `alert`, `confirm`, and `prompt`", r.Diagnoses[0].Msg, "should be ok")
}

func TestNoUnreachable(t *testing.T) {
	r := lint(t, "no_unreachable")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, "disallow unreachable code", r.Diagnoses[0].Msg, "should be ok")
}

func TestIgnore(t *testing.T) {
	r := lint(t, "ignore")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].Loc.Source(), "test1.js"), "should be ok")
}

func TestIgnoreFile(t *testing.T) {
	r := lint(t, "ignore_file")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].Loc.Source(), "test1.js"), "should be ok")
}

func TestIgnoreRoot(t *testing.T) {
	r := lint(t, "ignore_root")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].Loc.Source(), "a/test.js"), "should be ok")
}

func TestIgnoreNestOverride(t *testing.T) {
	// the nested `.eslintignore` needs a `.eslintrc.js` to active the nested config resolution
	r := lint(t, "ignore_nested")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(r.Diagnoses[0].Loc.Source(), "a/test.js"), "should be ok")
}

func TestDisableAll(t *testing.T) {
	r := lint(t, "disable_all")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 0, len(r.Diagnoses), "should be ok")
}

func TestEnableAll(t *testing.T) {
	r := lint(t, "enable_all")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

func TestDisableRules(t *testing.T) {
	r := lint(t, "disable_rules")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

func TestDisableNextLine(t *testing.T) {
	r := lint(t, "disable_next_line")
	util.AssertEqual(t, true, r.InternalError == nil, "should be ok")
	util.AssertEqual(t, 1, len(r.Diagnoses), "should be ok")
}

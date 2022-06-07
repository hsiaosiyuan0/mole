package plugin_basic

import (
	"path"
	"path/filepath"
	"plugin"
	"runtime"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/linter"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestGoPlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "go_plugin")
	util.ShellInDir(dir, "go", "build", "-buildmode=plugin")

	p, err := plugin.Open(path.Join(dir, "go_plugin.so"))
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

func TestResolvePlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "asset", "resolve_plugin")
	util.ShellInDir(dir, "npm", "i", "go-cross-ci-demo@0.0.5")

	p, err := linter.Resolve(dir, "go-cross-ci-demo")
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

package plugin_basic

import (
	"path"
	"path/filepath"
	"plugin"
	"runtime"
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestGoPlugin(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "go_plugin")
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

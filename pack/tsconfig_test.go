package pack

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestTsConfig(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "tsconfig", "sub", "sub")
	c, err := NewTsConfig(dir, "tsconfig.json")
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, filepath.IsAbs(c.CompilerOptions.BaseUrl), "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(c.CompilerOptions.BaseUrl, "/test/asset/tsconfig/src"), "should be ok")
	util.AssertEqual(t, true, filepath.IsAbs(c.CompilerOptions.Paths["jquery"][0]), "should be ok")
}

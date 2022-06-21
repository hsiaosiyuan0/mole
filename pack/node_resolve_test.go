package pack

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestNodeResolveIndex(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("./a", dir)
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "a/index.js"), "should be ok")
}

func TestNodeResolveSelf(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("node-resolve-index/b1.js", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")

	file, err = r.Resolve("node-resolve-index/b2", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveImport(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("#dep", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveModuleMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("mimic1", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "index1.js"), "should be ok")
}

func TestNodeResolveExports(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("mimic2/a", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveExportsMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, err := r.Resolve("mimic4", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "a.js"), "should be ok")
}

func TestNodeResolvePathMaps(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	jsConfig, err := NewTsConfig(dir, "jsconfig.json")
	if err != nil {
		t.Fatal(err)
	}

	pathMaps, err := jsConfig.PathMaps()
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, pathMaps)
	file, err := r.Resolve("@/app", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "app.js"), "should be ok")

	file, err = r.Resolve("@page/page-a", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "src/page/page-a.js"), "should be ok")

	file, err = r.Resolve("@scope/a", path.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "@scope/a/index.js"), "should be ok")
}

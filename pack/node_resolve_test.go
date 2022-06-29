package pack

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/util"
)

func TestNodeResolveIndex(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, pi, err := r.Resolve("./a", dir)
	util.AssertEqual(t, "node-resolve-index", pi.Name, "should be ok")
	util.AssertEqual(t, "0.0.1", pi.Version, "should be ok")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "a/index.js"), "should be ok")
}

func TestNodeResolveSelf(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, _, err := r.Resolve("node-resolve-index/b1.js", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")

	file, _, err = r.Resolve("node-resolve-index/b2", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveImport(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, _, err := r.Resolve("#dep", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveModuleMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, _, err := r.Resolve("mimic1", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "index1.js"), "should be ok")
}

func TestNodeResolveExports(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, pi, err := r.Resolve("mimic2/a", filepath.Join(dir, "a"))
	util.AssertEqual(t, "mimic2", pi.Name, "should be ok")
	util.AssertEqual(t, "0.0.1", pi.Version, "should be ok")
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "b.js"), "should be ok")
}

func TestNodeResolveExportsMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, _, err := r.Resolve("mimic4", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "a.js"), "should be ok")
}

func TestNodeResolveDirIdx(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(nil, nil, nil, nil, pkiLoader, false, nil)
	file, _, err := r.Resolve("mimic5", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "index.js"), "should be ok")
}

func TestNodeResolvePathMaps(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(nil)
	_, err := pkiLoader.Load(filepath.Join(dir, "package.json"))
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
	file, _, err := r.Resolve("@/app", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "app.js"), "should be ok")

	file, _, err = r.Resolve("@page/page-a", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "src/page/page-a.js"), "should be ok")

	file, _, err = r.Resolve("@scope/a", filepath.Join(dir, "a"))
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file[0], "@scope/a/index.js"), "should be ok")
}

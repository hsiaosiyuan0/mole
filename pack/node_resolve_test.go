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

	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "./a", dir, nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "a/index.js"), "should be ok")
}

func TestNodeResolveSelf(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "node-resolve-index/b1.js", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "b.js"), "should be ok")

	r = NewNodeResolver(pki, "node-resolve-index/b2", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err = r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "b.js"), "should be ok")
}

func TestNodeResolveImport(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "#dep", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "b.js"), "should be ok")
}

func TestNodeResolveModuleMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "mimic1", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "index1.js"), "should be ok")
}

func TestNodeResolveExports(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "mimic2/a", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "b.js"), "should be ok")
}

func TestNodeResolveExportsMain(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := path.Join(basepath, "test", "asset", "node-resolve-index")
	pkiLoader := NewPkginfoLoader(1024, 10)
	pki, err := pkiLoader.Load(path.Join(dir, "package.json"))
	if err != nil {
		t.Fatal(err)
	}

	r := NewNodeResolver(pki, "mimic4", path.Join(dir, "a"), nil, nil, nil, nil, pkiLoader)
	file, err := r.Resolve()
	util.AssertEqual(t, nil, err, "should be ok")
	util.AssertEqual(t, true, strings.HasSuffix(file, "a.js"), "should be ok")
}

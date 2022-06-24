package pack

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestParseDep(t *testing.T) {
	deps, err := parseDep("", `
  require('a.js')
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseDepRebound(t *testing.T) {
	deps, err := parseDep("", `
  require = a
  require('a.js')
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 0, len(deps), "should be ok")
}

func TestParseDepValShadow(t *testing.T) {
	deps, err := parseDep("", `
function f() {
  var require = a
  require('a.js')
}
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 0, len(deps), "should be ok")
}

func TestParseDepAfterValShadow(t *testing.T) {
	deps, err := parseDep("", `
function f() {
  var require = a
  require('a.js')
}
require('a.js')
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseDepImport(t *testing.T) {
	deps, err := parseDep("", `
import('a.js')
`, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseCondImport1(t *testing.T) {
	deps, err := parseDep("", `
if (process.env.NODE_ENV === 'production') {
  // DCE check should happen before ReactDOM bundle executes so that
  // DevTools can report bad minification during injection.
  checkDCE();
  module.exports = require('./cjs/react-dom.production.min.js');
} else {
  module.exports = require('./cjs/react-dom.development.js');
}
`, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "development",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
	util.AssertEqual(t, "./cjs/react-dom.development.js", deps[0], "should be ok")
}

func TestParseCondImport2(t *testing.T) {
	deps, err := parseDep("", `
if (process.env.NODE_ENV === 'production') {
  // DCE check should happen before ReactDOM bundle executes so that
  // DevTools can report bad minification during injection.
  checkDCE();
  module.exports = require('./cjs/react-dom.production.min.js');
} else {
  module.exports = require('./cjs/react-dom.development.js');
}
`, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "production",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
	util.AssertEqual(t, "./cjs/react-dom.production.min.js", deps[0], "should be ok")
}

func TestDepScanner(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "dep-scanner")
	util.ShellInDir(dir, "npm", "i")

	opts := NewDepScannerOpts()
	opts.dir = dir
	opts.entries = append(opts.entries, "./src/index.js")

	err := opts.SetTsconfig(opts.dir, "jsconfig.json", true)
	if err != nil {
		t.Fatal(err)
	}

	opts.SerVars(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "TEST",
			},
		},
	})

	s := NewDepScanner(opts)
	err = s.Run()
	if err != nil {
		t.Fatal(err)
	}

	cnt := 0
	for _, m := range s.modules {
		if m != nil {
			cnt += 1
			fmt.Println(m.File())
		}
	}

	peCnt := 0
	for _, err := range s.minors {
		if _, ok := err.(*parser.ParserError); ok {
			peCnt += 1
		}
	}

	util.AssertEqual(t, true, cnt > 0, "should be ok")
	util.AssertEqual(t, 0, peCnt, "should be ok")
}

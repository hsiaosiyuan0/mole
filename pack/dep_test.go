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
	p, err := parse("", `
  require('a.js')
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseDepRebound(t *testing.T) {
	p, err := parse("", `
  require = a
  require('a.js')
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 0, len(deps), "should be ok")
}

func TestParseDepValShadow(t *testing.T) {
	p, err := parse("", `
  function f() {
    var require = a
    require('a.js')
  }
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 0, len(deps), "should be ok")
}

func TestParseDepAfterValShadow(t *testing.T) {
	p, err := parse("", `
  function f() {
    var require = a
    require('a.js')
  }
  require('a.js')
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseDepImport(t *testing.T) {
	p, err := parse("", `
  import('a.js')
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
}

func TestParseCondImport1(t *testing.T) {
	vars := map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "development",
			},
		},
	}

	p, err := parse("", `
if (process.env.NODE_ENV === 'production') {
  // DCE check should happen before ReactDOM bundle executes so that
  // DevTools can report bad minification during injection.
  checkDCE();
  module.exports = require('./cjs/react-dom.production.min.js');
} else {
  module.exports = require('./cjs/react-dom.development.js');
}
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, vars, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
	util.AssertEqual(t, "./cjs/react-dom.development.js", deps[0].file, "should be ok")
}

func TestParseCondImport2(t *testing.T) {
	vars := map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "production",
			},
		},
	}

	p, err := parse("", `
if (process.env.NODE_ENV === 'production') {
  // DCE check should happen before ReactDOM bundle executes so that
  // DevTools can report bad minification during injection.
  checkDCE();
  module.exports = require('./cjs/react-dom.production.min.js');
} else {
  module.exports = require('./cjs/react-dom.development.js');
}
`, parser.NewParserOpts(), true)

	if err != nil {
		t.Fatal(err)
	}

	deps, _, err := walkDep(p, vars, nil)
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 1, len(deps), "should be ok")
	util.AssertEqual(t, "./cjs/react-dom.production.min.js", deps[0].file, "should be ok")
}

func TestDepScanner(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "dep-scanner")
	util.ShellInDir(dir, "npm", "ci")

	opts := NewDepScannerOpts()
	opts.Dir = dir
	opts.Entries = append(opts.Entries, "./src/index.js")

	err := opts.SetTsconfig(opts.Dir, "jsconfig.json", true)
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
	err = s.ResolveDeps()
	if err != nil {
		t.Fatal(err)
	}

	cnt := 0
	for _, m := range s.fileModules {
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

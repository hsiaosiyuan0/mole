package pack

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/util"
)

func TestDepScanner(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	dir := filepath.Join(basepath, "test", "asset", "dep-scanner")
	util.ShellInDir(dir, "npm", "i")

	opts := NewDepScannerOpts()
	opts.dir = dir
	opts.entries = append(opts.entries, "src/index.js")

	err := opts.SetTsconfig(opts.dir, "jsconfig.json", true)
	if err != nil {
		t.Fatal(err)
	}

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

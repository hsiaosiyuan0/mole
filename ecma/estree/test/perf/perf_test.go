package estree_test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
)

func compileToESTree(code string, toEstree bool) error {
	s := span.NewSource("", code)
	opts := parser.NewParserOpts()
	p := parser.NewParser(s, opts)
	ast, err := p.Prog()
	if err != nil {
		return err
	}

	if toEstree {
		ctx := estree.NewConvertCtx(p)
		ctx.LineCol = false
		estree.ConvertProg(ast.(*parser.Prog), ctx)
	}
	return nil
}

func BenchmarkParsingToESTree(t *testing.B) {
	libs := []*struct {
		name string
		code string
	}{
		{"angular.js", ""},
		{"backbone.js", ""},
		{"ember.js", ""},
		{"jquery.js", ""},
		{"react-dom.js", ""},
		{"react.js", ""},
	}

	_, fileName, _, _ := runtime.Caller(0)
	for _, lib := range libs {
		b, err := ioutil.ReadFile(filepath.Join(path.Dir(fileName), "asset", lib.name))
		if err != nil {
			t.Fatal(err)
		}
		lib.code = string(b)
	}

	for _, lib := range libs {
		t.Run(lib.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := compileToESTree(lib.code, true)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParsing(t *testing.B) {
	libs := []*struct {
		name string
		code string
	}{
		{"angular.js", ""},
		{"backbone.js", ""},
		{"ember.js", ""},
		{"jquery.js", ""},
		{"react-dom.js", ""},
		{"react.js", ""},
	}

	_, fileName, _, _ := runtime.Caller(0)
	for _, lib := range libs {
		b, err := ioutil.ReadFile(filepath.Join(path.Dir(fileName), "asset", lib.name))
		if err != nil {
			t.Fatal(err)
		}
		lib.code = string(b)
	}

	for _, lib := range libs {
		t.Run(lib.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := compileToESTree(lib.code, false)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

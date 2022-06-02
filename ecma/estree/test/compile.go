package estree_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	internal "github.com/hsiaosiyuan0/mole/util"
)

func NewParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func CompileProg(code string, opts *parser.ParserOpts) (parser.Node, error) {
	p := NewParser(code, opts)
	return p.Prog()
}

func Compile(code string) (string, error) {
	return CompileWithOpts(code, parser.NewParserOpts())
}

func CompileWithOpts(code string, opts *parser.ParserOpts) (string, error) {
	s := span.NewSource("", code)
	p := parser.NewParser(s, opts)
	ast, err := p.Prog()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(estree.ConvertProg(ast.(*parser.Prog), estree.NewConvertCtx()))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	return out.String(), nil
}

func TestFail(t *testing.T, code, errMs string, opts *parser.ParserOpts) {
	ast, err := CompileProg(code, opts)
	if err == nil {
		t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
	}
	internal.AssertEqual(t, errMs, err.Error(), "")
}

type Fixture struct {
	name   string
	input  string
	output string
	opts   string
}

func ScanFixtures(name string) (map[string]*Fixture, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	fxs := make(map[string]*Fixture)
	err := filepath.Walk(path.Join(basepath, "fixture", name), func(pathStr string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() {
				return nil
			}
			dir := filepath.Dir(pathStr)
			name := strings.Trim(strings.Replace(dir, path.Join(basepath, "fixture"), "", 1), string(os.PathSeparator))
			fx := fxs[name]
			if fx == nil {
				fx = &Fixture{name, "", "", ""}
				fxs[name] = fx
			}
			n := info.Name()
			if strings.HasPrefix(n, "input") {
				fx.input = pathStr
			} else if strings.HasPrefix(n, "output") {
				fx.output = pathStr
			} else if strings.HasPrefix(n, "options") {
				fx.opts = pathStr
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return fxs, nil
}

func RunFixture(t *testing.T, input, output, opts string, defaultOpts *parser.ParserOpts) {
	code, err := ioutil.ReadFile(input)
	if err != nil {
		t.Fatalf("failed to read fixture code at: %s\nerror: %v", input, err)
	}
	out, err := ioutil.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read fixture output at: %s\nerror: %v", output, err)
	}
	jsonObj := make(map[string]interface{})
	if err = json.Unmarshal(out, &jsonObj); err != nil {
		t.Fatalf("failed to decode fixture output at: %s\nerror: %v", output, err)
	}

	theOpts := defaultOpts.Clone()
	if opts != "" {
		optsRaw, err := ioutil.ReadFile(opts)
		if err != nil {
			t.Fatalf("failed to read options at: %s\nerror: %v", opts, err)
		}
		optsObj := make(map[string]interface{})
		if err = json.Unmarshal(optsRaw, &optsObj); err != nil {
			t.Fatalf("failed to decode options at: %s\nerror: %v", opts, err)
		}
		theOpts.MergeJson(optsObj)
	}

	ast, err := CompileWithOpts(string(code), theOpts)
	if jsonObj["throws"] != nil {
		if err == nil {
			t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
		}
		internal.AssertEqual(t, jsonObj["throws"].(string), err.Error(), "")
	} else {
		if err != nil {
			t.Fatalf("failed to parse fixture at: %s\nerror: %v", input, err)
		}
		internal.AssertEqualJson(t, string(out), ast)
	}
}

func RunFixtures(t *testing.T, name string, defaultOpts *parser.ParserOpts) {
	fxs, err := ScanFixtures(name)
	if err != nil {
		t.Fatalf("failed to run fixture [%s] %v", name, err)
	}

	t.Logf("Running %d fixtures in [%s]...", len(fxs), name)
	for _, fx := range fxs {
		t.Run(fx.name, func(t *testing.T) {
			RunFixture(t, fx.input, fx.output, fx.opts, defaultOpts)
		})
	}
}

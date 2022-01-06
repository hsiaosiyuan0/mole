package estree_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/internal"
)

type Fixture struct {
	name   string
	input  string
	output string
}

func scanFixtures(name string) (map[string]*Fixture, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	fxs := make(map[string]*Fixture)
	err := filepath.Walk(path.Join(basepath, name), func(path string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() {
				return nil
			}
			dir := filepath.Dir(path)
			name := strings.Trim(strings.Replace(dir, basepath, "", 1), string(os.PathSeparator))
			fx := fxs[name]
			if fx == nil {
				fx = &Fixture{name, "", ""}
				fxs[name] = fx
			}
			if strings.HasPrefix(info.Name(), "input") {
				fx.input = path
			} else {
				fx.output = path
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return fxs, nil
}

func runFixture(t *testing.T, input, output string, opts *parser.ParserOpts) {
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

	ast, err := CompileWithOpts(string(code), opts)
	if jsonObj["throws"] != nil {
		if err == nil {
			t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
		}
		AssertEqual(t, jsonObj["throws"].(string), err.Error(), "")
	} else {
		if err != nil {
			t.Fatalf("failed to parse fixture at: %s\nerror: %v", input, err)
		}
		AssertEqualJson(t, string(out), ast)
	}
}

func runFixtures(t *testing.T, name string, opts *parser.ParserOpts) {
	if opts == nil {
		opts = parser.NewParserOpts()
	}

	fxs, err := scanFixtures(name)
	if err != nil {
		t.Fatalf("failed to run fixture [%s] %v", name, err)
	}

	t.Logf("Running %d fixtures in [%s]...", len(fxs), name)
	for _, fx := range fxs {
		t.Run(fx.name, func(t *testing.T) {
			runFixture(t, fx.input, fx.output, opts)
		})
	}
}

func TestFixture_es2015(t *testing.T) {
	runFixtures(t, "es2015", nil)
}

func TestFixture_ts(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_TS)
	runFixtures(t, "typescript", opts)
}

func TestFixture_tsManually(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_TS)
	runFixtures(t, "typescript/cast/as-const", opts)
}

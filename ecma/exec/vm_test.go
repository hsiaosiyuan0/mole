package exec

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compile(code string, opts *parser.ParserOpts) (*parser.Parser, parser.Node, *parser.SymTab, error) {
	p := newParser(code, opts)
	ast, err := p.Prog()
	if err != nil {
		return nil, nil, nil, err
	}
	return p, ast, p.Symtab(), nil
}

func TestExecExprAdd(t *testing.T) {
	_, ast, _, err := compile(`
  1 + 2
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 3.0, res, "should be ok")
}

func TestExecNull(t *testing.T) {
	_, ast, _, err := compile(`
  null
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, nil, res, "should be ok")
}

func TestExecUndef(t *testing.T) {
	_, ast, _, err := compile(`
  undefined
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, _ := ee.Exec(nil).GetResult()
	util.AssertEqual(t, nil, res, "should be ok")
}

func TestExecExprAddStr(t *testing.T) {
	_, ast, _, err := compile(`
  1 + '2'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, "12", res, "should be ok")
}

func TestExecExprEqual(t *testing.T) {
	_, ast, _, err := compile(`
  1 + 2 == 3
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecExprMemExpr(t *testing.T) {
	_, ast, _, err := compile(`
  process.env.ENV
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}).GetResult()

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, "development", res, "should be ok")
}

func TestExecExprMemExpr2(t *testing.T) {
	_, ast, _, err := compile(`
  process['env']['ENV']
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}).GetResult()

	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, "development", res, "should be ok")
}

func TestExecExprMemExprEqual(t *testing.T) {
	_, ast, _, err := compile(`
  process['env']['ENV'] === 'development'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecExprToBool(t *testing.T) {
	_, ast, _, err := compile(`
  !(process['env']['ENV'] === 'development')
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, false, res, "should be ok")
}

func TestExecExprToBool2(t *testing.T) {
	_, ast, _, err := compile(`
  !process['not set']
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecBool(t *testing.T) {
	_, ast, _, err := compile(`
  true
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecLogic(t *testing.T) {
	_, ast, _, err := compile(`
  process.env.ENV === 'development' || process.env.ENV === 'production'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	// 1
	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")

	// 2
	res, err = ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "production",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecArrLit(t *testing.T) {
	_, ast, _, err := compile(`
  [1, 2]
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, 2, len(res.([]interface{})), "should be ok")
}

func TestExecArrLitIdx(t *testing.T) {
	_, ast, _, err := compile(`
  [1, 2][0] == 1
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecArrLitNoIdx(t *testing.T) {
	_, ast, _, err := compile(`
  [1, 2][3] == null
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecArrIncludes(t *testing.T) {
	_, ast, _, err := compile(`
  [1, 2].includes(1)
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(nil).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

func TestExecArrIncludes2(t *testing.T) {
	_, ast, _, err := compile(`
  ["REG", "ONLINE"].includes(process.env.ENV)
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	// 1
	ee := NewExprEvaluator(ast)
	res, err := ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "TEST",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, false, res, "should be ok")

	// 2
	res, err = ee.Exec(map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "ONLINE",
			},
		},
	}).GetResult()
	if err != nil {
		t.Fatal(err)
	}

	util.AssertEqual(t, true, res, "should be ok")
}

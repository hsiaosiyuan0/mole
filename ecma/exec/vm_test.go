package exec

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
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
	_, ast, symtab, err := compile(`
  1 + 2
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)
	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, 3.0, res, "should be ok")

	ee.Release()
}

func TestExecExprAddStr(t *testing.T) {
	_, ast, symtab, err := compile(`
  1 + '2'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)
	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, "12", res, "should be ok")

	ee.Release()
}

func TestExecExprEqual(t *testing.T) {
	_, ast, symtab, err := compile(`
  1 + 2 == 3
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)
	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.Release()
}

func TestExecExprMemExpr(t *testing.T) {
	_, ast, symtab, err := compile(`
  process.env.ENV
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, "development", res, "should be ok")

	ee.Release()
}

func TestExecExprMemExpr2(t *testing.T) {
	_, ast, symtab, err := compile(`
  process['env']['ENV']
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, "development", res, "should be ok")

	ee.Release()
}

func TestExecExprMemExprEqual(t *testing.T) {
	_, ast, symtab, err := compile(`
  process['env']['ENV'] === 'development'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.Release()
}

func TestExecExprToBool(t *testing.T) {
	_, ast, symtab, err := compile(`
  !(process['env']['ENV'] === 'development')
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, false, res, "should be ok")

	ee.Release()
}

func TestExecExprToBool2(t *testing.T) {
	_, ast, symtab, err := compile(`
  !process['not set']
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.Release()
}

func TestExecBool(t *testing.T) {
	_, ast, symtab, err := compile(`
  true
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())

	res := ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.Release()
}

func TestExecLogic(t *testing.T) {
	_, ast, symtab, err := compile(`
  process.env.ENV === 'development' || process.env.ENV === 'production'
  `, nil)
	util.AssertEqual(t, nil, err, "should pass")

	ctx := walk.NewWalkCtx(ast, symtab)
	ee := NewExprEvaluator(ctx)

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "development",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())
	res := ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.vars = map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"ENV": "production",
			},
		},
	}

	walk.VisitNode(ast, "", ctx.VisitorCtx())
	res = ee.GetResult()
	util.AssertEqual(t, true, res, "should be ok")

	ee.Release()
}

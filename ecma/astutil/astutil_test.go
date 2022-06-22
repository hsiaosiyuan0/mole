package astutil

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

func compile(code string, opts *parser.ParserOpts) (parser.Node, error) {
	p := newParser(code, opts)
	return p.Prog()
}

func TestIfStmtToSwitchBranches(t *testing.T) {
	ast, err := compile(`
  if (a) {
    b
  } else if (c) {
    d
  } else {
    e
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bs := IfStmtToSwitchBranches(ast.(*parser.Prog).Body()[0].(*parser.IfStmt))
	util.AssertEqual(t, 3, len(bs), "should be ok")
}

func TestIsNodeContains(t *testing.T) {
	ast, err := compile(`
  if (a) {
    b
  } else if (c) {
    d
  } else {
    e
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	bs := IfStmtToSwitchBranches(ifSmt)
	ok := IsNodeContains(ifSmt, bs[len(bs)-1].body.(*parser.BlockStmt).Body()[0])
	util.AssertEqual(t, true, ok, "should be ok")
}

func TestSelectIf(t *testing.T) {
	ast, err := compile(`
  if (a) {
    b
  } else if (c) {
    d
  } else {
    e
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	// 1
	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := SelectTrueBranches(ifSmt, map[string]interface{}{
		"a": 0,
		"c": 0,
	})

	e := nodes[0].(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident).Text()
	util.AssertEqual(t, "e", e, "should be ok")

	// 2
	nodes = SelectTrueBranches(ifSmt, map[string]interface{}{
		"a": 0,
		"c": 1,
	})
	d := nodes[0].(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident).Text()
	util.AssertEqual(t, "d", d, "should be ok")
}

func TestSelectBin(t *testing.T) {
	ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 1,
	})

	a := nodes[0].(*parser.Ident).Text()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestSelectBin1(t *testing.T) {
	ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	// 1
	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 1,
		"b": 1,
	})

	a := nodes[0].(*parser.Ident).Text()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestSelectBin2(t *testing.T) {
	ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 0,
		"b": 1,
	})

	a := nodes[0].(*parser.Ident).Text()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestCollectNodesInTrueBranches(t *testing.T) {
	ast, err := compile(`
  a && require("a.js")
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := CollectNodesInTrueBranches(bin, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"require": 1,
	})

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Text()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches1(t *testing.T) {
	ast, err := compile(`
  if (a) {
    a && require("a.js")
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"require": 1,
	})

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Text()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches2(t *testing.T) {
	ast, err := compile(`
  if (a) {
   if (b) {
     a && require("a.js")
   }
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"b":       1,
		"require": 1,
	})

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Text()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches3(t *testing.T) {
	ast, err := compile(`
  if (a) {
   if (b) {
     a && require("a.js")
   }
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"require": 1,
	})

	util.AssertEqual(t, 0, len(nodes), "should be ok")
}

func TestCollectNodesInTrueBranches4(t *testing.T) {
	ast, err := compile(`
  if (process.env.NODE_ENV === 'production') {
    module.exports = require('./cjs/react.production.min.js');
  } else {
    module.exports = require('./cjs/react.development.js');
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	// 1
	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "development",
			},
		},
	})
	util.AssertEqual(t, "./cjs/react.development.js", nodes[0].(*parser.CallExpr).Args()[0].(*parser.StrLit).Text(), "should be ok")

	// 2
	nodes = CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "production",
			},
		},
	})
	util.AssertEqual(t, "./cjs/react.production.min.js", nodes[0].(*parser.CallExpr).Args()[0].(*parser.StrLit).Text(), "should be ok")
}

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

func compile(code string, opts *parser.ParserOpts) (*parser.Parser, parser.Node, error) {
	p := newParser(code, opts)
	prg, err := p.Prog()
	return p, prg, err
}

func TestIfStmtToSwitchBranches(t *testing.T) {
	_, ast, err := compile(`
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
	_, ast, err := compile(`
  if (a) {
    b
  } else if (c) {
    d
  } else {
    e
  }

  f
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	prog := ast.(*parser.Prog)
	ifSmt := prog.Body()[0].(*parser.IfStmt)
	bs := IfStmtToSwitchBranches(ifSmt)
	ok := IsNodeContains(ifSmt, bs[len(bs)-1].body.(*parser.BlockStmt).Body()[0])
	util.AssertEqual(t, true, ok, "should be ok")

	ok = IsNodeContains(ifSmt, prog.Body()[1])
	util.AssertEqual(t, false, ok, "should be ok")
}

func TestSelectIf(t *testing.T) {
	p, ast, err := compile(`
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
	}, p)

	e := nodes[0].(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident).Val()
	util.AssertEqual(t, "e", e, "should be ok")

	// 2
	nodes = SelectTrueBranches(ifSmt, map[string]interface{}{
		"a": 0,
		"c": 1,
	}, p)
	d := nodes[0].(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident).Val()
	util.AssertEqual(t, "d", d, "should be ok")
}

func TestSelectBin(t *testing.T) {
	p, ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 1,
	}, p)

	a := nodes[0].(*parser.Ident).Val()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestSelectBin1(t *testing.T) {
	p, ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	// 1
	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 1,
		"b": 1,
	}, p)

	a := nodes[0].(*parser.Ident).Val()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestSelectBin2(t *testing.T) {
	p, ast, err := compile(`
  a && b
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := SelectTrueBranches(bin, map[string]interface{}{
		"a": 0,
		"b": 1,
	}, p)

	a := nodes[0].(*parser.Ident).Val()
	util.AssertEqual(t, "a", a, "should be ok")
}

func TestCollectNodesInTrueBranches(t *testing.T) {
	p, ast, err := compile(`
  a && require("a.js")
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	bin := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	nodes := CollectNodesInTrueBranches(bin, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"require": 1,
	}, p)

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Val()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches1(t *testing.T) {
	p, ast, err := compile(`
  if (a) {
    a && require("a.js")
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"require": 1,
	}, p)

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Val()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches2(t *testing.T) {
	p, ast, err := compile(`
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
	}, p)

	require := nodes[0].(*parser.CallExpr).Callee().(*parser.Ident).Val()
	util.AssertEqual(t, "require", require, "should be ok")
}

func TestCollectNodesInTrueBranches3(t *testing.T) {
	p, ast, err := compile(`
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
	}, p)

	util.AssertEqual(t, 0, len(nodes), "should be ok")
}

func TestCollectNodesInTrueBranches4(t *testing.T) {
	p, ast, err := compile(`
  if (a) {
   require("a1.js")
   if (b) {
     a && require("a2.js")
     require("a3.js")
   }
   require("a4.js")
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"b":       1,
		"require": 1,
	}, p)

	util.AssertEqual(t, 4, len(nodes), "should be ok")
}

func TestCollectNodesInTrueBranches5(t *testing.T) {
	p, ast, err := compile(`
  if (a) {
   require("a1.js")
   if (b) {
     a && require("a2.js")
     require("a3.js")
   }
  }
  require("a4.js")
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	nodes := CollectNodesInTrueBranches(ast, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"a":       1,
		"b":       1,
		"require": 1,
	}, p)

	util.AssertEqual(t, 4, len(nodes), "should be ok")
}

func TestCollectNodesInTrueBranches6(t *testing.T) {
	p, ast, err := compile(`
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
	}, p)
	util.AssertEqual(t, "./cjs/react.development.js", nodes[0].(*parser.CallExpr).Args()[0].(*parser.StrLit).Val(), "should be ok")

	// 2
	nodes = CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "production",
			},
		},
	}, p)
	util.AssertEqual(t, "./cjs/react.production.min.js", nodes[0].(*parser.CallExpr).Args()[0].(*parser.StrLit).Val(), "should be ok")
}

func TestCollectNodesInTrueBranches7(t *testing.T) {
	p, ast, err := compile(`
  if (process.env.NODE_ENV !== "production") {
    (function() {
  'use strict';

  var React = require('react');
  var _assign = require('object-assign');
  var Scheduler = require('scheduler');
  var checkPropTypes = require('prop-types/checkPropTypes');
  var tracing = require('scheduler/tracing');
    })()
  }
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	// 1
	ifSmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	nodes := CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "production",
			},
		},
	}, p)
	util.AssertEqual(t, 0, len(nodes), "should be ok")

	// 2
	nodes = CollectNodesInTrueBranches(ifSmt, []parser.NodeType{parser.N_EXPR_CALL}, map[string]interface{}{
		"process": map[string]interface{}{
			"env": map[string]interface{}{
				"NODE_ENV": "development",
			},
		},
	}, p)
	util.AssertEqual(t, 6, len(nodes), "should be ok")
}

func TestBuildFnDepGraph(t *testing.T) {
	p, ast, err := compile(`
const imgSuffix = (picUrl) => ImageSuffix.directSuffix(picUrl, { paramWidth: 50 })
function test() {
  imgSuffix()
}
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	graph := BuildFnDepGraph(ast, p.Symtab())
	stmts := ast.(*parser.Prog).Body()

	imgSuffix := graph.Nodes[stmts[0].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()]
	test := graph.Nodes[stmts[1].(*parser.FnDec)]
	util.AssertEqual(t, true, test.Deps[0] == imgSuffix, "should be prog ok")
}

func TestIsFnDepsOnNode(t *testing.T) {
	p, ast, err := compile(`
const fn0 = () => {}
const fn1 = () => fn0()
const fn2 = () => fn1()
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	graph := BuildFnDepGraph(ast, p.Symtab())
	stmts := ast.(*parser.Prog).Body()

	fn2 := stmts[2].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()
	fn0 := stmts[0]
	util.AssertEqual(t, true, IsFnDepsOnNode(graph, fn2, fn0), "should be prog ok")
}

func TestIsFnDepsOnNode1(t *testing.T) {
	p, ast, err := compile(`
function fn0() {}
const fn1 = () => fn0()
const fn2 = () => fn1()
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	graph := BuildFnDepGraph(ast, p.Symtab())
	stmts := ast.(*parser.Prog).Body()

	fn2 := stmts[2].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()
	fn0 := stmts[0]
	util.AssertEqual(t, true, IsFnDepsOnNode(graph, fn2, fn0), "should be prog ok")
}

func TestIsFnDepsOnNode2(t *testing.T) {
	p, ast, err := compile(`
import ImageSuffix from '@music/mobile-image'

function fn0() { ImageSuffix() }
const fn1 = () => fn0()
const fn2 = () => fn1()
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	graph := BuildFnDepGraph(ast, p.Symtab())
	stmts := ast.(*parser.Prog).Body()

	fn2 := stmts[3].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()
	ifx := stmts[0]
	util.AssertEqual(t, true, IsFnDepsOnNode(graph, fn2, ifx), "should be prog ok")
}

func TestIsFnDepsOnNode3(t *testing.T) {
	p, ast, err := compile(`
import ImageSuffix from '@music/mobile-image'

function fn0() { ImageSuffix() }
const fn1 = () => () => fn0();
const fn2 = () => fn1()
`, nil)
	util.AssertEqual(t, nil, err, "should be prog ok")

	graph := BuildFnDepGraph(ast, p.Symtab())
	stmts := ast.(*parser.Prog).Body()

	fn2 := stmts[3].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()
	ifx := stmts[0]
	util.AssertEqual(t, true, IsFnDepsOnNode(graph, fn2, ifx), "should be prog ok")
}

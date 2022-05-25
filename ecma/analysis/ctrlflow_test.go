package analysis

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
	. "github.com/hsiaosiyuan0/mole/util"
)

func newParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := span.NewSource("", code)
	return parser.NewParser(s, opts)
}

func compile(code string, opts *parser.ParserOpts) (parser.Node, *parser.SymTab, error) {
	p := newParser(code, opts)
	ast, err := p.Prog()
	if err != nil {
		return nil, nil, err
	}
	return ast, p.Symtab(), nil
}

func TestCtrlflow_Basic(t *testing.T) {
	ast, symtab, err := compile(`
a
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()

	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident)
	block := astNodeMap[expr]

	AssertEqual(t, 5, len(block.Nodes), "should be ok")
	AssertEqual(t, "Prog:enter", nodeToString(block.Nodes[0]), "should be ok")
	AssertEqual(t, "ExprStmt:enter", nodeToString(block.Nodes[1]), "should be ok")
	AssertEqual(t, "Ident(a)", nodeToString(block.Nodes[2]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(block.Nodes[3]), "should be ok")
	AssertEqual(t, "Prog:exit", nodeToString(block.Nodes[4]), "should be ok")
}

func TestCtrlflow_Logic(t *testing.T) {
	ast, symtab, err := compile(`
 a && b
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	a := astNodeMap[expr.Lhs()]
	b := astNodeMap[expr.Rhs()]
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "BinExpr(&&):exit", nodeToString(b.Nodes[1]), "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "ExprStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
}

func TestCtrlflow_LogicMix(t *testing.T) {
	ast, symtab, err := compile(`
 a && b || c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	ab := expr.Lhs().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]
	c := astNodeMap[expr.Rhs()]
	AssertEqual(t, c, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "BinExpr(&&):exit", nodeToString(b.Nodes[1]), "should be ok")
	AssertEqual(t, "BinExpr(||):exit", nodeToString(c.Nodes[1]), "should be ok")

	bTrue := b.OutJmpEdge(ET_JMP_T)
	exit := bTrue.Dst.Nodes[0]
	AssertEqual(t, "ExprStmt:exit", nodeToString(exit), "should be ok")
}

func TestCtrlflow_IfStmt(t *testing.T) {
	ast, symtab, err := compile(`
  a;
  if (b) c;
  else d;
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	b := astNodeMap[stmt.Test()]
	c := astNodeMap[stmt.Cons().(*parser.ExprStmt).Expr()]
	d := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	AssertEqual(t, d, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(c.Nodes[2]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(d.Nodes[2]), "should be ok")

	ifExit := c.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "IfStmt:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_IfBlkStmt(t *testing.T) {
	ast, symtab, err := compile(`
  a;
  if (b) {
    c;
    d;
  }
  else e
  f;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr()]
	e := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	cons := stmt.Cons().(*parser.BlockStmt).Body()
	c := astNodeMap[cons[0].(*parser.ExprStmt).Expr()]
	d := astNodeMap[cons[1].(*parser.ExprStmt).Expr()]
	AssertEqual(t, c, d, "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(d.Nodes[6]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(d.Nodes[7]), "should be ok")

	AssertEqual(t, e, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, a.OutSeqEdge().Dst, "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "IfStmt:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_IfBlk2Stmt(t *testing.T) {
	ast, symtab, err := compile(`
  a;
  if (b) {
    c;
    d;
  }
  else {
    e;
  }
  f;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr()]
	b := astNodeMap[stmt.Test()]
	e := astNodeMap[stmt.Alt().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()]
	cons := stmt.Cons().(*parser.BlockStmt).Body()
	c := astNodeMap[cons[0].(*parser.ExprStmt).Expr()]
	d := astNodeMap[cons[1].(*parser.ExprStmt).Expr()]
	AssertEqual(t, a, b, "should be ok")
	AssertEqual(t, c, d, "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(d.Nodes[6]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(d.Nodes[7]), "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(e.Nodes[3]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(e.Nodes[4]), "should be ok")

	AssertEqual(t, e, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, a.OutSeqEdge().Dst, "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "IfStmt:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_IfLogic(t *testing.T) {
	ast, symtab, err := compile(`
  if (a && b) c
  else d
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	ab := stmt.Test().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]
	c := astNodeMap[stmt.Cons().(*parser.ExprStmt).Expr()]
	d := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(c.Nodes[2]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(d.Nodes[2]), "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "IfStmt:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_IfLogicMix(t *testing.T) {
	ast, symtab, err := compile(`
  a
  if (b || c && d) e
  else f
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr()]
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	bcd := stmt.Test().(*parser.BinExpr)
	b := astNodeMap[bcd.Lhs()]
	cd := bcd.Rhs().(*parser.BinExpr)
	c := astNodeMap[cd.Lhs()]
	d := astNodeMap[cd.Rhs()]
	e := astNodeMap[stmt.Cons().(*parser.ExprStmt).Expr()]
	f := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	AssertEqual(t, a, b, "should be ok")

	AssertEqual(t, e, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")
	AssertEqual(t, c, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, f, c.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, f, d.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, e, d.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "ExprStmt:exit", nodeToString(e.Nodes[2]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(f.Nodes[2]), "should be ok")

	ifExit := e.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "IfStmt:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_UpdateExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a++
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.UpdateExpr).Arg()]

	AssertEqual(t, "UpdateExpr(++):exit", nodeToString(a.Nodes[4]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "Prog:exit", nodeToString(a.Nodes[6]), "should be ok")
}

func TestCtrlflow_VarDecStmt(t *testing.T) {
	ast, symtab, err := compile(`
  let a = 1, c = d
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	varDecStmt := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt)
	varDec0 := varDecStmt.DecList()[0].(*parser.VarDec)
	varDec1 := varDecStmt.DecList()[1].(*parser.VarDec)
	a := astNodeMap[varDec0.Id()]
	c := astNodeMap[varDec1.Id()]

	AssertEqual(t, a, c, "should be ok")

	AssertEqual(t, "VarDecStmt:enter", nodeToString(a.Nodes[1]), "should be ok")
	AssertEqual(t, "VarDec:enter", nodeToString(a.Nodes[2]), "should be ok")
	AssertEqual(t, "VarDec:exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "VarDec:enter", nodeToString(a.Nodes[6]), "should be ok")
	AssertEqual(t, "VarDec:exit", nodeToString(a.Nodes[9]), "should be ok")
	AssertEqual(t, "VarDecStmt:exit", nodeToString(a.Nodes[10]), "should be ok")
	AssertEqual(t, "Prog:exit", nodeToString(a.Nodes[11]), "should be ok")
}

func TestCtrlflow_ForStmt(t *testing.T) {
	ast, symtab, err := compile(`
  for (let b = 1; b < c; b++) {
    d;
  }
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	init := stmt.Init().(*parser.VarDecStmt)
	test := stmt.Test().(*parser.BinExpr)
	update := stmt.Update().(*parser.UpdateExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	b := astNodeMap[init.DecList()[0].(*parser.VarDec).Id()]
	bc := astNodeMap[test.Lhs()]
	bu := astNodeMap[update.Arg()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "VarDecStmt:exit", nodeToString(b.Nodes[7]), "should be ok")
	AssertEqual(t, "BinExpr(<):enter", nodeToString(bc.Nodes[0]), "should be ok")
	AssertEqual(t, "BinExpr(<):exit", nodeToString(bc.Nodes[3]), "should be ok")

	AssertEqual(t, bc, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, e, bc.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, "ForStmt:exit", nodeToString(e.Nodes[0]), "should be ok")
	AssertEqual(t, "ExprStmt:enter", nodeToString(e.Nodes[1]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(e.Nodes[3]), "should be ok")

	AssertEqual(t, bu, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "UpdateExpr(++):exit", nodeToString(bu.Nodes[7]), "should be ok")

	AssertEqual(t, d, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, bc, bu.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_ForStmtOmitInit(t *testing.T) {
	ast, symtab, err := compile(`
  for (; b < c; b++) {
    d;
  }
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	update := stmt.Update().(*parser.UpdateExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	bc := astNodeMap[test.Lhs()]
	bu := astNodeMap[update.Arg()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "BinExpr(<):enter", nodeToString(bc.Nodes[0]), "should be ok")
	AssertEqual(t, "BinExpr(<):exit", nodeToString(bc.Nodes[3]), "should be ok")

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "ForStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")

	AssertEqual(t, bc, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, e, bc.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, "ForStmt:exit", nodeToString(e.Nodes[0]), "should be ok")
	AssertEqual(t, "ExprStmt:enter", nodeToString(e.Nodes[1]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(e.Nodes[3]), "should be ok")

	AssertEqual(t, bu, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "UpdateExpr(++):exit", nodeToString(bu.Nodes[7]), "should be ok")

	AssertEqual(t, d, bu, "should be ok")
	AssertEqual(t, d, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, bc, bu.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_ForStmtOmitInitUpdate(t *testing.T) {
	ast, symtab, err := compile(`
  for (; b < c;) {
    d;
  }
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	bc := astNodeMap[test.Lhs()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "ForStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, bc, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, e, bc.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, bc, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_ForStmtOmitInitUpdate_TestLogic(t *testing.T) {
	ast, symtab, err := compile(`
  for (; b && c;) {
    d;
  }
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	b := astNodeMap[test.Lhs()]
	c := astNodeMap[test.Rhs()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "ForStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, e, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, e, c.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, b, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, "BlockStmt:enter", nodeToString(d.Nodes[0]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(d.Nodes[4]), "should be ok")

	AssertEqual(t, "ForStmt:exit", nodeToString(e.Nodes[0]), "should be ok")
}

func TestCtrlflow_ForStmtOmitAll(t *testing.T) {
	ast, symtab, err := compile(`
  for (; ;) {
    d;
  }
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "ForStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, d, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, d, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")

	AssertEqual(t, nil, e, "should be ok")
}

func TestCtrlflow_While(t *testing.T) {
	ast, symtab, err := compile(`
  while(a) b;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.Ident)
	body := stmt.Body().(*parser.ExprStmt).Expr()
	a := astNodeMap[test]
	b := astNodeMap[body]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "WhileStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst.Nodes[0]
	AssertEqual(t, "WhileStmt:exit", nodeToString(exit), "should be ok")

	AssertEqual(t, a, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_WhileBodyBlk(t *testing.T) {
	ast, symtab, err := compile(`
  while(a) {
    b;
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.Ident)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test]
	b := astNodeMap[body]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "WhileStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst.Nodes[0]
	AssertEqual(t, "WhileStmt:exit", nodeToString(exit), "should be ok")

	AssertEqual(t, "BlockStmt:enter", nodeToString(b.Nodes[0]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(b.Nodes[4]), "should be ok")
}

func TestCtrlflow_WhileLogicOr(t *testing.T) {
	ast, symtab, err := compile(`
  while((a + b) || c) {
    d;
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test.Lhs().(*parser.ParenExpr).Expr().(*parser.BinExpr).Lhs()]
	b := astNodeMap[test.Lhs().(*parser.ParenExpr).Expr().(*parser.BinExpr).Rhs()]
	c := astNodeMap[test.Rhs()]
	d := astNodeMap[body]

	AssertEqual(t, "Prog:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "WhileStmt:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a, "should be ok")

	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")
	AssertEqual(t, c, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "BinExpr(+):exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "ParenExpr:exit", nodeToString(a.Nodes[6]), "should be ok")

	exit := c.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "WhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, "BlockStmt:enter", nodeToString(d.Nodes[0]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(d.Nodes[4]), "should be ok")
}

func TestCtrlflow_ParenExpr(t *testing.T) {
	ast, symtab, err := compile(`
  (a + b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.ParenExpr).Expr().(*parser.BinExpr).Lhs()]
	AssertEqual(t, 10, len(a.Nodes), "should be ok")
	AssertEqual(t, "BinExpr(+):exit", nodeToString(a.Nodes[6]), "should be ok")
	AssertEqual(t, "ParenExpr:exit", nodeToString(a.Nodes[7]), "should be ok")
}

func TestCtrlflow_ParenExprLogic(t *testing.T) {
	ast, symtab, err := compile(`
  (a && b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	ab := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.ParenExpr).Expr().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]

	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "BinExpr(&&):exit", nodeToString(b.Nodes[1]), "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "ParenExpr:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "ExprStmt:exit", nodeToString(exit.Nodes[1]), "should be ok")
}

func TestCtrlflow_DoWhile(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.DoWhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.Ident)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	b := astNodeMap[test]
	a := astNodeMap[body]

	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, b, "should be ok")
	AssertEqual(t, a, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "DoWhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
}

func TestCtrlflow_DoWhileLogicOr(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b || c)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.DoWhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	b := astNodeMap[test.Lhs()]
	c := astNodeMap[test.Rhs()]
	a := astNodeMap[body]

	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, b, "should be ok")

	AssertEqual(t, a, b.OutJmpEdge(ET_JMP_T).Dst, "should be ok")
	AssertEqual(t, a, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, c.OutJmpEdge(ET_LOOP).Dst, "should be ok")

	exit := c.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "DoWhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "Prog:exit", nodeToString(exit.Nodes[1]), "should be ok")
}

func TestCtrlflow_DoWhileLogicAnd(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b && c)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.DoWhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	b := astNodeMap[test.Lhs()]
	c := astNodeMap[test.Rhs()]
	a := astNodeMap[body]

	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, b, "should be ok")

	exit := b.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "DoWhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "Prog:exit", nodeToString(exit.Nodes[1]), "should be ok")

	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, exit, c.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, a, c.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_WhileLogicMix(t *testing.T) {
	ast, symtab, err := compile(`
  while(a || b && c) {
    d
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test.Lhs()]
	b := astNodeMap[test.Rhs().(*parser.BinExpr).Lhs()]
	c := astNodeMap[test.Rhs().(*parser.BinExpr).Rhs()]
	d := astNodeMap[body]
	exit := c.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "WhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")

	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, exit, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_ContinueBasicEntry(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a) {
    continue LabelA
    d
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.LabelStmt).Body().(*parser.WhileStmt)
	begin := ana.Graph().Head
	a := astNodeMap[stmt.Test()]
	cont := astNodeMap[stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ContStmt).Label()]
	d := astNodeMap[stmt.Body().(*parser.BlockStmt).Body()[1].(*parser.ExprStmt).Expr()]
	exit := a.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "WhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "LabelStmt:exit", nodeToString(exit.Nodes[1]), "should be ok")

	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, cont, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, cont.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, d, cont.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, true, cont.HasOutCut(), "should be ok")

	AssertEqual(t, a, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, true, d.HasOutCut(), "should be ok")
}

func TestCtrlflow_Continue(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a || b && c) {
    continue LabelA
    d
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.LabelStmt).Body().(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	a := astNodeMap[test.Lhs()]
	b := astNodeMap[test.Rhs().(*parser.BinExpr).Lhs()]
	c := astNodeMap[test.Rhs().(*parser.BinExpr).Rhs()]
	cont := astNodeMap[stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ContStmt).Label()]
	d := astNodeMap[stmt.Body().(*parser.BlockStmt).Body()[1].(*parser.ExprStmt).Expr()]
	exit := c.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "WhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "LabelStmt:exit", nodeToString(exit.Nodes[1]), "should be ok")

	AssertEqual(t, a, begin.OutSeqEdge().Dst.NextBlk(), "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, cont, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")

	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, exit, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, cont, c.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, exit, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, a, cont.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, d, cont.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, true, cont.HasOutCut(), "should be ok")

	AssertEqual(t, begin.OutSeqEdge().Dst, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, true, d.HasOutCut(), "should be ok")
}

func TestCtrlflow_ContinueOuter(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a) {
    while(b) {
      continue LabelA
      c
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	while1 := ast.(*parser.Prog).Body()[0].(*parser.LabelStmt).Body().(*parser.WhileStmt)
	while2 := while1.Body().(*parser.BlockStmt).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	a := astNodeMap[while1.Test()]
	b := astNodeMap[while2.Test()]
	c := astNodeMap[while2.Body().(*parser.BlockStmt).Body()[1].(*parser.ExprStmt).Expr()]
	cont := astNodeMap[while2.Body().(*parser.BlockStmt).Body()[0].(*parser.ContStmt).Label()]
	exit := a.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "WhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "LabelStmt:exit", nodeToString(exit.Nodes[1]), "should be ok")

	innerWhileEnter := a.OutSeqEdge().Dst
	AssertEqual(t, "BlockStmt:enter", nodeToString(innerWhileEnter.Nodes[0]), "should be ok")
	AssertEqual(t, "WhileStmt:enter", nodeToString(innerWhileEnter.Nodes[1]), "should be ok")

	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, b, innerWhileEnter.OutSeqEdge().Dst, "should be ok")

	innerWhileExit := b.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "WhileStmt:exit", nodeToString(innerWhileExit.Nodes[0]), "should be ok")
	AssertEqual(t, "BlockStmt:exit", nodeToString(innerWhileExit.Nodes[1]), "should be ok")
	AssertEqual(t, cont, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, cont.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, c, cont.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, true, cont.HasOutCut(), "should be ok")

	AssertEqual(t, b, c.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, true, c.HasOutCut(), "should be ok")
}

func TestCtrlflow_ContinueDoWhile(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: do {
  a;
  do {
    continue LabelA;
    d;
  } while (c);
} while (b);
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, _, _, astNodeMap := ana.Graph().NodesEdges()
	doWhile1 := ast.(*parser.Prog).Body()[0].(*parser.LabelStmt).Body().(*parser.DoWhileStmt)
	doWhile2 := doWhile1.Body().(*parser.BlockStmt).Body()[1].(*parser.DoWhileStmt)
	begin := ana.Graph().Head
	a := astNodeMap[doWhile1.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()]
	b := astNodeMap[doWhile1.Test()]
	c := astNodeMap[doWhile2.Test()]
	cont := astNodeMap[doWhile2.Body().(*parser.BlockStmt).Body()[0].(*parser.ContStmt).Label()]
	d := astNodeMap[doWhile2.Body().(*parser.BlockStmt).Body()[1].(*parser.ExprStmt).Expr()]
	exit := b.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "DoWhileStmt:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "LabelStmt:exit", nodeToString(exit.Nodes[1]), "should be ok")

	AssertEqual(t, a, begin.OutSeqEdge().Dst.NextBlk(), "should be ok")
	AssertEqual(t, cont, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, cont.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, d, cont.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, true, cont.HasOutCut(), "should be ok")

	AssertEqual(t, c, d, "should be ok")
	AssertEqual(t, cont, c.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, true, c.HasOutCut(), "should be ok")
	AssertEqual(t, b, c.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, begin.OutSeqEdge().Dst, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_ContinueFor(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for(let a = 1; a < 10; a++) {
  for(let b = a; b < 10; b++) {
    if (b > 3) {
      continue LabelA
      c
    }
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nLabelStmt:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\n"];
loc2_16_54[label="Ident(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
loc2_23_156_0[label="BinExpr(<):enter\nIdent(a)\nNumLit(10)\nBinExpr(<):exit\n"];
loc2_36_156_0[label="BlockStmt:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(b)\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\n"];
loc2_8_156_1[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
loc3_17_156_0[label="BinExpr(<):enter\nIdent(b)\nNumLit(10)\nBinExpr(<):exit\n"];
loc3_2_156_1[label="ForStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\n"];
loc3_30_156_0[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(b)\nNumLit(3)\nBinExpr(>):exit\n"];
loc4_15_156_0[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
loc4_4_156_1[label="IfStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
loc6_6_156_0[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_16_54 [xlabel="",color="black"];
loc2_16_54->loc2_23_156_0 [xlabel="",color="black"];
loc2_23_156_0->loc2_36_156_0 [xlabel="",color="black"];
loc2_23_156_0->loc2_8_156_1 [xlabel="F",color="orange"];
loc2_36_156_0->loc3_17_156_0 [xlabel="",color="black"];
loc2_8_156_1->final [xlabel="",color="black"];
loc3_17_156_0->loc3_2_156_1 [xlabel="F",color="orange"];
loc3_17_156_0->loc3_30_156_0 [xlabel="",color="black"];
loc3_2_156_1:s->loc2_23_156_0:ne [xlabel="L",color="orange"];
loc3_30_156_0->loc4_15_156_0 [xlabel="",color="black"];
loc3_30_156_0->loc4_4_156_1 [xlabel="F",color="orange"];
loc4_15_156_0:s->loc2_16_54:ne [xlabel="L",color="orange"];
loc4_15_156_0->loc6_6_156_0 [xlabel="",color="red"];
loc4_4_156_1:s->loc3_17_156_0:ne [xlabel="L",color="orange"];
loc6_6_156_0->loc4_4_156_1 [xlabel="",color="red"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ForIn(t *testing.T) {
	ast, symtab, err := compile(`
for(a in b) {
    c
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nForInOfStmt:enter\n"];
loc2_0_156_1[label="ForInOfStmt:exit\nProg:exit\n"];
loc2_12_156_0[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc2_4_54[label="Ident(a)\nIdent(b)\n"];
loc0->loc2_4_54 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc2_12_156_0:s->loc2_4_54:ne [xlabel="L",color="orange"];
loc2_4_54->loc2_0_156_1 [xlabel="F",color="orange"];
loc2_4_54->loc2_12_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ForInLet(t *testing.T) {
	ast, symtab, err := compile(`
for(let a in b) {
    c
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nForInOfStmt:enter\n"];
loc2_0_156_1[label="ForInOfStmt:exit\nProg:exit\n"];
loc2_16_156_0[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc2_4_156_0[label="VarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\nIdent(b)\n"];
loc0->loc2_4_156_0 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc2_16_156_0:s->loc2_4_156_0:ne [xlabel="L",color="orange"];
loc2_4_156_0->loc2_0_156_1 [xlabel="F",color="orange"];
loc2_4_156_0->loc2_16_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForIn(t *testing.T) {
	ast, symtab, err := compile(`
s
LabelA: for(a in b) {
  for(c in d) {
    if (e && f) {
      continue LabelA
    }
    g
  }
  h
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nIdent(s)\nExprStmt:exit\nLabelStmt:enter\nForInOfStmt:enter\n"];
loc3_12_54[label="Ident(a)\nIdent(b)\n"];
loc3_20_156_0[label="BlockStmt:enter\nForInOfStmt:enter\n"];
loc3_8_156_1[label="ForInOfStmt:exit\nLabelStmt:exit\nProg:exit\n"];
loc4_14_156_0[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(&&):enter\nIdent(e)\n"];
loc4_2_156_1[label="ForInOfStmt:exit\nExprStmt:enter\nIdent(h)\nExprStmt:exit\nBlockStmt:exit\n"];
loc4_6_54[label="Ident(c)\nIdent(d)\n"];
loc5_13_54[label="Ident(f)\nBinExpr(&&):exit\n"];
loc5_16_156_0[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
loc5_16_156_1[label="BlockStmt:exit\n"];
loc5_4_156_1[label="IfStmt:exit\nExprStmt:enter\nIdent(g)\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc3_12_54 [xlabel="",color="black"];
loc3_12_54->loc3_20_156_0 [xlabel="",color="black"];
loc3_12_54->loc3_8_156_1 [xlabel="F",color="orange"];
loc3_20_156_0->loc4_6_54 [xlabel="",color="black"];
loc3_8_156_1->final [xlabel="",color="black"];
loc4_14_156_0->loc5_13_54 [xlabel="",color="black"];
loc4_14_156_0->loc5_4_156_1 [xlabel="F",color="orange"];
loc4_2_156_1:s->loc3_12_54:ne [xlabel="L",color="orange"];
loc4_6_54->loc4_14_156_0 [xlabel="",color="black"];
loc4_6_54->loc4_2_156_1 [xlabel="F",color="orange"];
loc5_13_54->loc5_16_156_0 [xlabel="",color="black"];
loc5_13_54->loc5_4_156_1 [xlabel="F",color="orange"];
loc5_16_156_0:s->loc3_12_54:ne [xlabel="L",color="orange"];
loc5_16_156_0->loc5_16_156_1 [xlabel="",color="red"];
loc5_16_156_1->loc5_4_156_1 [xlabel="",color="red"];
loc5_4_156_1:s->loc4_6_54:ne [xlabel="L",color="orange"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueNoLabel(t *testing.T) {
	ast, symtab, err := compile(`
while(a) {
  if (a > 10) continue;
  a--
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nWhileStmt:enter\n"];
loc2_0_156_1[label="WhileStmt:exit\nProg:exit\n"];
loc2_6_54[label="Ident(a)\n"];
loc2_9_156_0[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
loc3_14_156_0[label="ContStmt:enter\nContStmt:exit\n"];
loc3_2_156_1[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_6_54 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc2_6_54->loc2_0_156_1 [xlabel="F",color="orange"];
loc2_6_54->loc2_9_156_0 [xlabel="",color="black"];
loc2_9_156_0->loc3_14_156_0 [xlabel="",color="black"];
loc2_9_156_0->loc3_2_156_1 [xlabel="F",color="orange"];
loc3_14_156_0:s->loc2_6_54:ne [xlabel="L",color="orange"];
loc3_14_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1:s->loc2_6_54:ne [xlabel="L",color="orange"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakNoLabelForNest(t *testing.T) {
	ast, symtab, err := compile(`
while(a) {
  if (a > 10) break;
  a--
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nWhileStmt:enter\n"];
loc2_0_156_1[label="WhileStmt:exit\nProg:exit\n"];
loc2_6_54[label="Ident(a)\n"];
loc2_9_156_0[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
loc3_14_156_0[label="BrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_6_54 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc2_6_54->loc2_0_156_1 [xlabel="F",color="orange"];
loc2_6_54->loc2_9_156_0 [xlabel="",color="black"];
loc2_9_156_0->loc3_14_156_0 [xlabel="",color="black"];
loc2_9_156_0->loc3_2_156_1 [xlabel="F",color="orange"];
loc3_14_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_14_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1:s->loc2_6_54:ne [xlabel="L",color="orange"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakNoLabelFor(t *testing.T) {
	ast, symtab, err := compile(`
for (;a;) {
  if (a > 10) {
    break;
    c
  }
  a--
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nForStmt:enter\n"];
loc2_0_156_1[label="ForStmt:exit\nProg:exit\n"];
loc2_10_156_0[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
loc2_6_54[label="Ident(a)\n"];
loc3_14_156_0[label="BlockStmt:enter\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
loc5_4_156_0[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_6_54 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc2_10_156_0->loc3_14_156_0 [xlabel="",color="black"];
loc2_10_156_0->loc3_2_156_1 [xlabel="F",color="orange"];
loc2_6_54->loc2_0_156_1 [xlabel="F",color="orange"];
loc2_6_54->loc2_10_156_0 [xlabel="",color="black"];
loc3_14_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_14_156_0->loc5_4_156_0 [xlabel="",color="red"];
loc3_2_156_1:s->loc2_6_54:ne [xlabel="L",color="orange"];
loc5_4_156_0->loc3_2_156_1 [xlabel="",color="red"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakLabelWhile(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a) {
    for(;b;) {
      break LabelA
      c
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nLabelStmt:enter\nWhileStmt:enter\n"];
loc2_10_156_1[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
loc2_16_54[label="Ident(a)\n"];
loc2_19_156_0[label="BlockStmt:enter\nForStmt:enter\n"];
loc3_13_156_0[label="BlockStmt:enter\nBrkStmt:enter\nIdent(LabelA)\nBrkStmt:exit\n"];
loc3_4_156_1[label="ForStmt:exit\nBlockStmt:exit\n"];
loc3_9_54[label="Ident(b)\n"];
loc5_6_156_0[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_16_54 [xlabel="",color="black"];
loc2_10_156_1->final [xlabel="",color="black"];
loc2_16_54->loc2_10_156_1 [xlabel="F",color="orange"];
loc2_16_54->loc2_19_156_0 [xlabel="",color="black"];
loc2_19_156_0->loc3_9_54 [xlabel="",color="black"];
loc3_13_156_0->loc2_10_156_1 [xlabel="U",color="orange"];
loc3_13_156_0->loc5_6_156_0 [xlabel="",color="red"];
loc3_4_156_1:s->loc2_16_54:ne [xlabel="L",color="orange"];
loc3_9_54->loc3_13_156_0 [xlabel="",color="black"];
loc3_9_54->loc3_4_156_1 [xlabel="F",color="orange"];
loc5_6_156_0:s->loc3_9_54:ne [xlabel="L",color="red"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakNoLabelWhile(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a) {
    for(;b;) {
      break
      c
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nLabelStmt:enter\nWhileStmt:enter\n"];
loc2_10_156_1[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
loc2_16_54[label="Ident(a)\n"];
loc2_19_156_0[label="BlockStmt:enter\nForStmt:enter\n"];
loc3_13_156_0[label="BlockStmt:enter\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_4_156_1[label="ForStmt:exit\nBlockStmt:exit\n"];
loc3_9_54[label="Ident(b)\n"];
loc5_6_156_0[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc0->loc2_16_54 [xlabel="",color="black"];
loc2_10_156_1->final [xlabel="",color="black"];
loc2_16_54->loc2_10_156_1 [xlabel="F",color="orange"];
loc2_16_54->loc2_19_156_0 [xlabel="",color="black"];
loc2_19_156_0->loc3_9_54 [xlabel="",color="black"];
loc3_13_156_0->loc3_4_156_1 [xlabel="U",color="orange"];
loc3_13_156_0->loc5_6_156_0 [xlabel="",color="red"];
loc3_4_156_1:s->loc2_16_54:ne [xlabel="L",color="orange"];
loc3_9_54->loc3_13_156_0 [xlabel="",color="black"];
loc3_9_54->loc3_4_156_1 [xlabel="F",color="orange"];
loc5_6_156_0:s->loc3_9_54:ne [xlabel="L",color="red"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CallExpr(t *testing.T) {
	ast, symtab, err := compile(`
  fn()
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nCallExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CallExprArgs(t *testing.T) {
	ast, symtab, err := compile(`
  fn(1, 2, a && b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nNumLit(1)\nNumLit(2)\nBinExpr(&&):enter\nIdent(a)\n"];
loc2_16_54[label="Ident(b)\nBinExpr(&&):exit\n"];
loc2_2_156_1[label="CallExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->loc2_16_54 [xlabel="",color="black"];
loc0->loc2_2_156_1 [xlabel="F",color="orange"];
loc2_16_54->loc2_2_156_1 [xlabel="",color="black"];
loc2_2_156_1->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CallExprArgsLoc(t *testing.T) {
	ast, symtab, err := compile(`
  fn(1, a || b && c, d && e)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nNumLit(1)\nBinExpr(||):enter\nIdent(a)\n"];
loc2_13_156_0[label="BinExpr(&&):enter\nIdent(b)\n"];
loc2_18_54[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
loc2_21_156_0[label="BinExpr(&&):enter\nIdent(d)\n"];
loc2_26_54[label="Ident(e)\nBinExpr(&&):exit\n"];
loc2_2_156_1[label="CallExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->loc2_13_156_0 [xlabel="",color="black"];
loc0->loc2_21_156_0 [xlabel="T",color="orange"];
loc2_13_156_0->loc2_18_54 [xlabel="",color="black"];
loc2_13_156_0->loc2_21_156_0 [xlabel="F",color="orange"];
loc2_18_54->loc2_21_156_0 [xlabel="",color="black"];
loc2_21_156_0->loc2_26_54 [xlabel="",color="black"];
loc2_21_156_0->loc2_2_156_1 [xlabel="F",color="orange"];
loc2_26_54->loc2_2_156_1 [xlabel="",color="black"];
loc2_2_156_1->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_Nullish(t *testing.T) {
	ast, symtab, err := compile(`
  const foo = null ?? 'default string';
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(foo)\nBinExpr(??):enter\nNullLit\n"];
loc2_22_31[label="StrLit\nBinExpr(??):exit\n"];
loc2_8_156_1[label="VarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
loc0->loc2_22_31 [xlabel="",color="black"];
loc0->loc2_8_156_1 [xlabel="T",color="orange"];
loc2_22_31->loc2_8_156_1 [xlabel="",color="black"];
loc2_8_156_1->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_FnDecOuter(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nFnDec\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_FnDecBody(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  a
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_0_156_0[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
loc2_0_156_0->final [xlabel="",color="black"];
initial->loc2_0_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ReturnNoArg(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  a;
  return;
  b
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_0_156_0[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
loc2_0_156_1[label="FnDec:exit\n"];
loc5_2_156_0[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\n"];
loc2_0_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc2_0_156_0->loc5_2_156_0 [xlabel="",color="red"];
loc2_0_156_1->final [xlabel="",color="black"];
loc5_2_156_0->loc2_0_156_1 [xlabel="",color="red"];
initial->loc2_0_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ReturnArg(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  a;
  return a ?? b;
  c
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_0_156_0[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nRetStmt:enter\nBinExpr(??):enter\nIdent(a)\n"];
loc2_0_156_1[label="FnDec:exit\n"];
loc4_14_54[label="Ident(b)\nBinExpr(??):exit\n"];
loc4_2_156_1[label="RetStmt:exit\n"];
loc5_2_156_0[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
loc2_0_156_0->loc4_14_54 [xlabel="",color="black"];
loc2_0_156_0->loc4_2_156_1 [xlabel="T",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc4_14_54->loc4_2_156_1 [xlabel="",color="black"];
loc4_2_156_1->loc2_0_156_1 [xlabel="U",color="orange"];
loc4_2_156_1->loc5_2_156_0 [xlabel="",color="red"];
loc5_2_156_0->loc2_0_156_1 [xlabel="",color="red"];
initial->loc2_0_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_FnDecParam(t *testing.T) {
	ast, symtab, err := compile(`
function f(a, b) {}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nFnDec\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
  `, ana.Graph().Dot(), "global should be ok")

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_0_156_0[label="FnDec:enter\nIdent(f)\nIdent(a)\nIdent(b)\nBlockStmt:enter\nBlockStmt:exit\nFnDec:exit\n"];
loc2_0_156_0->final [xlabel="",color="black"];
initial->loc2_0_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "fn should be ok")
}

func TestCtrlflow_FnExpr(t *testing.T) {
	ast, symtab, err := compile(`
fn = function f(a, b) { c }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(fn)\nFnDec\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
    `, ana.Graph().Dot(), "global should be ok")

	fn := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.AssignExpr).Rhs()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_5_156_0[label="FnDec:enter\nIdent(f)\nIdent(a)\nIdent(b)\nBlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
loc2_5_156_0->final [xlabel="",color="black"];
initial->loc2_5_156_0 [xlabel="",color="black"];
}
	`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ArrLit(t *testing.T) {
	ast, symtab, err := compile(`
a = [b, c, d ?? e]
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nIdent(b)\nIdent(c)\nBinExpr(??):enter\nIdent(d)\n"];
loc2_16_54[label="Ident(e)\nBinExpr(??):exit\n"];
loc2_4_156_1[label="ArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->loc2_16_54 [xlabel="",color="black"];
loc0->loc2_4_156_1 [xlabel="T",color="orange"];
loc2_16_54->loc2_4_156_1 [xlabel="",color="black"];
loc2_4_156_1->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ObjLit(t *testing.T) {
	ast, symtab, err := compile(`
a = { b: 1, c: d }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(c)\nIdent(d)\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ObjLitNest(t *testing.T) {
	ast, symtab, err := compile(`
a = { b: {c: 1} }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ObjLitFn(t *testing.T) {
	ast, symtab, err := compile(`
a = { b: function (a, b) {} }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nFnDec\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
	`, ana.Graph().Dot(), "global should be ok")

	rhs := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.AssignExpr).Rhs().(*parser.ObjLit)
	fn := rhs.Props()[0].(*parser.Prop).Val()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_9_156_0[label="FnDec:enter\nIdent(a)\nIdent(b)\nBlockStmt:enter\nBlockStmt:exit\nFnDec:exit\n"];
loc2_9_156_0->final [xlabel="",color="black"];
initial->loc2_9_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "fn should be ok")

}

func TestCtrlflow_EmptyExpr(t *testing.T) {
	ast, symtab, err := compile(`
a;
;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DebugStmt(t *testing.T) {
	ast, symtab, err := compile(`
a
debugger
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nDebugStmt\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitNull(t *testing.T) {
	ast, symtab, err := compile(`
a = null
null
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nNullLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nNullLit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitBool(t *testing.T) {
	ast, symtab, err := compile(`
a = true
false
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nBoolLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nBoolLit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitNum(t *testing.T) {
	ast, symtab, err := compile(`
a = 1
1.1
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nNumLit(1)\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nNumLit(1.1)\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitStr(t *testing.T) {
	ast, symtab, err := compile(`
a = "str"
'str'
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nStrLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nStrLit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitArr(t *testing.T) {
	ast, symtab, err := compile(`
a = [1, 2, 3];
[1, 2, 3]
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nNumLit(1)\nNumLit(2)\nNumLit(3)\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nArrLit:enter\nNumLit(1)\nNumLit(2)\nNumLit(3)\nArrLit:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitArrNest(t *testing.T) {
	ast, symtab, err := compile(`
a = [1, 2, [4, 5, 6]];
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:enter\nNumLit(4)\nNumLit(5)\nNumLit(6)\nArrLit:exit\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitObj(t *testing.T) {
	ast, symtab, err := compile(`
a = {b: { c: 1 } };
({b: { c: 1 } })
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nParenExpr:enter\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitObjNest(t *testing.T) {
	ast, symtab, err := compile(`
a = {b: { c: 1, d: [1, { f: 2}] } };
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(d)\nArrLit:enter\nNumLit(1)\nObjLit:enter\nProp:enter\nIdent(f)\nNumLit(2)\nProp:exit\nObjLit:exit\nArrLit:exit\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitRexexp(t *testing.T) {
	ast, symtab, err := compile(`
a = /reg/;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nRegLit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
loc0->final [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCase(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1();
  case 2: doSth2();
  case 3: doSth3();
  default: doSthDefault();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc6_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc3_10_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_10_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="F",color="orange"];
loc5_10_156_0->loc6_11_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc6_2_156_0 [xlabel="F",color="orange"];
loc6_11_156_0->final [xlabel="",color="black"];
loc6_2_156_0->loc6_11_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBlk(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1();
  case 2: {
    doSth21();
    doSth22();
  }
  case 3: doSth3();
  default: doSthDefault();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_10_156_0[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth21)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth22)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc8_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc8_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc9_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc9_2_156_0[label="SwitchCase:enter\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc3_10_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_10_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc8_2_156_0 [xlabel="F",color="orange"];
loc8_10_156_0->loc9_11_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc9_2_156_0 [xlabel="F",color="orange"];
loc9_11_156_0->final [xlabel="",color="black"];
loc9_2_156_0->loc9_11_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseFallthrough(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1:
  case 2: {
    doSth21();
    doSth22();
  }
  case 3: doSth3();
  default: doSthDefault();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc4_10_156_0[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth21)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth22)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc8_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc8_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc9_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc9_2_156_0[label="SwitchCase:enter\n"];
loc0->loc3_2_156_1 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc3_2_156_1->loc4_10_156_0 [xlabel="",color="black"];
loc4_10_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc8_2_156_0 [xlabel="F",color="orange"];
loc8_10_156_0->loc9_11_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc9_2_156_0 [xlabel="F",color="orange"];
loc9_11_156_0->final [xlabel="",color="black"];
loc9_2_156_0->loc9_11_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseDefaultFirst(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  default: doSthDefault();
  case 1: doSth1();
  case 2: doSth2();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\n"];
loc3_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(1)\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc0->loc4_2_156_0 [xlabel="",color="black"];
loc3_11_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_10_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="F",color="orange"];
loc5_10_156_0->final [xlabel="",color="black"];
loc5_2_156_0->loc3_11_156_0 [xlabel="F",color="orange"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseDefaultInMiddle(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1();
  default: doSthDefault();
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc6_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc3_10_156_0->loc4_11_156_0 [xlabel="",color="black"];
loc4_11_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="",color="black"];
loc5_10_156_0->loc6_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc6_2_156_0 [xlabel="F",color="orange"];
loc6_10_156_0->final [xlabel="",color="black"];
loc6_2_156_0->loc4_11_156_0 [xlabel="F",color="orange"];
loc6_2_156_0->loc6_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseDefaultInMiddleFallthrough(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1:
  default: doSthDefault();
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc4_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc6_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_2_156_1 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc3_2_156_1->loc4_11_156_0 [xlabel="",color="black"];
loc4_11_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="",color="black"];
loc5_10_156_0->loc6_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc6_2_156_0 [xlabel="F",color="orange"];
loc6_10_156_0->final [xlabel="",color="black"];
loc6_2_156_0->loc4_11_156_0 [xlabel="F",color="orange"];
loc6_2_156_0->loc6_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBasic(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1(); break;
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc4_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_10_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_10_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1->loc4_10_156_0 [xlabel="",color="red"];
loc4_10_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="F",color="orange"];
loc5_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc5_2_156_0->final [xlabel="F",color="orange"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBreak(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1(); break;
  default: doSthDefault();
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc4_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc6_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_10_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_10_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1->loc4_11_156_0 [xlabel="",color="red"];
loc4_11_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="",color="black"];
loc5_10_156_0->loc6_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc6_2_156_0 [xlabel="F",color="orange"];
loc6_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc6_2_156_0->loc4_11_156_0 [xlabel="F",color="orange"];
loc6_2_156_0->loc6_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBreakInBlock(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: {
    doSth1(); break;
  }
  default: doSthDefault();
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_10_156_0[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_10_156_1[label="BlockStmt:exit\nSwitchCase:exit\n"];
loc6_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\n"];
loc7_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc7_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc8_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc8_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc6_2_156_0 [xlabel="F",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_10_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_10_156_0->loc3_10_156_1 [xlabel="",color="red"];
loc3_10_156_1->loc6_11_156_0 [xlabel="",color="red"];
loc6_11_156_0->loc7_10_156_0 [xlabel="",color="black"];
loc6_2_156_0->loc7_2_156_0 [xlabel="",color="black"];
loc7_10_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc7_2_156_0->loc7_10_156_0 [xlabel="",color="black"];
loc7_2_156_0->loc8_2_156_0 [xlabel="F",color="orange"];
loc8_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc8_2_156_0->loc6_11_156_0 [xlabel="F",color="orange"];
loc8_2_156_0->loc8_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBreakSiblingBlock(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: {
    doSth1();
  }
  break;
  default: doSthDefault();
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_10_156_0[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc7_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc7_2_156_0[label="SwitchCase:enter\n"];
loc8_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc8_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc9_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc9_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc7_2_156_0 [xlabel="F",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_10_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_10_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1->loc7_11_156_0 [xlabel="",color="red"];
loc7_11_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc7_2_156_0->loc8_2_156_0 [xlabel="",color="black"];
loc8_10_156_0->loc9_10_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc8_10_156_0 [xlabel="",color="black"];
loc8_2_156_0->loc9_2_156_0 [xlabel="F",color="orange"];
loc9_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc9_2_156_0->loc7_11_156_0 [xlabel="F",color="orange"];
loc9_2_156_0->loc9_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBreakDefaultFirst(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  default: doSthDefault(); break;
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc3_2_156_1[label="SwitchCase:exit\n"];
loc4_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc4_2_156_0 [xlabel="",color="black"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_11_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc3_11_156_0->loc3_2_156_1 [xlabel="",color="red"];
loc3_2_156_1->loc4_10_156_0 [xlabel="",color="red"];
loc4_10_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc4_10_156_0 [xlabel="",color="black"];
loc4_2_156_0->loc5_2_156_0 [xlabel="F",color="orange"];
loc5_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc5_2_156_0->loc3_11_156_0 [xlabel="F",color="orange"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseBreakDefaultMiddle(t *testing.T) {
	ast, symtab, err := compile(`
switch(a) {
  case 1: doSth1();
  default: doSthDefault(); break;
  case 2: doSth2();
  case 3: doSth3();
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc2_0_156_1[label="SwitchStmt:exit\nProg:exit\n"];
loc3_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc4_11_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
loc4_2_156_0[label="SwitchCase:enter\n"];
loc4_2_156_1[label="SwitchCase:exit\n"];
loc5_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc5_2_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc6_10_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc6_2_156_0[label="SwitchCase:enter\nNumLit(3)\n"];
loc0->loc3_10_156_0 [xlabel="",color="black"];
loc0->loc4_2_156_0 [xlabel="F",color="orange"];
loc2_0_156_1->final [xlabel="",color="black"];
loc3_10_156_0->loc4_11_156_0 [xlabel="",color="black"];
loc4_11_156_0->loc2_0_156_1 [xlabel="U",color="orange"];
loc4_11_156_0->loc4_2_156_1 [xlabel="",color="red"];
loc4_2_156_0->loc5_2_156_0 [xlabel="",color="black"];
loc4_2_156_1->loc5_10_156_0 [xlabel="",color="red"];
loc5_10_156_0->loc6_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc5_10_156_0 [xlabel="",color="black"];
loc5_2_156_0->loc6_2_156_0 [xlabel="F",color="orange"];
loc6_10_156_0->loc2_0_156_1 [xlabel="",color="black"];
loc6_2_156_0->loc4_11_156_0 [xlabel="F",color="orange"];
loc6_2_156_0->loc6_10_156_0 [xlabel="",color="black"];
initial->loc0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SwitchCaseReturnDefaultMiddle(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  switch(a) {
    case 1: {
      doSth1(); return;
    }
    default: doSthDefault();
    case 2: doSth2();
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
loc2_0_156_0[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
loc4_12_156_0[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
loc4_12_156_1[label="BlockStmt:exit\nSwitchCase:exit\n"];
loc7_13_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
loc7_4_156_0[label="SwitchCase:enter\n"];
loc8_12_156_0[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
loc8_4_156_0[label="SwitchCase:enter\nNumLit(2)\n"];
loc2_0_156_0->loc4_12_156_0 [xlabel="",color="black"];
loc2_0_156_0->loc7_4_156_0 [xlabel="F",color="orange"];
loc4_12_156_0->final [xlabel="U",color="orange"];
loc4_12_156_0->loc4_12_156_1 [xlabel="",color="red"];
loc4_12_156_1->loc7_13_156_0 [xlabel="",color="red"];
loc7_13_156_0->loc8_12_156_0 [xlabel="",color="black"];
loc7_4_156_0->loc8_4_156_0 [xlabel="",color="black"];
loc8_12_156_0->final [xlabel="",color="black"];
loc8_4_156_0->loc7_13_156_0 [xlabel="F",color="orange"];
loc8_4_156_0->loc8_12_156_0 [xlabel="",color="black"];
initial->loc2_0_156_0 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

package analysis

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	. "github.com/hsiaosiyuan0/mole/fuzz"
	"github.com/hsiaosiyuan0/mole/span"
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

	_, _, astNodeMap := ana.Graph().NodesEdges()

	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.Ident)
	block := astNodeMap[expr]

	AssertEqual(t, 5, len(block.Nodes), "should be ok")
	AssertEqual(t, "N_PROG:enter", nodeToString(block.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:enter", nodeToString(block.Nodes[1]), "should be ok")
	AssertEqual(t, "N_NAME(a)", nodeToString(block.Nodes[2]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(block.Nodes[3]), "should be ok")
	AssertEqual(t, "N_PROG:exit", nodeToString(block.Nodes[4]), "should be ok")
}

func TestCtrlflow_Logic(t *testing.T) {
	ast, symtab, err := compile(`
 a && b
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	a := astNodeMap[expr.Lhs()]
	b := astNodeMap[expr.Rhs()]
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "N_EXPR_BIN(&&):exit", nodeToString(b.Nodes[1]), "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(exit.Nodes[0]), "should be ok")
}

func TestCtrlflow_LogicMix(t *testing.T) {
	ast, symtab, err := compile(`
 a && b || c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	expr := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.BinExpr)
	ab := expr.Lhs().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]
	c := astNodeMap[expr.Rhs()]
	AssertEqual(t, c, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "N_EXPR_BIN(&&):exit", nodeToString(b.Nodes[1]), "should be ok")
	AssertEqual(t, "N_EXPR_BIN(||):exit", nodeToString(c.Nodes[1]), "should be ok")

	bTrue := b.OutJmpEdge(ET_JMP_T)
	exit := bTrue.Dst.Nodes[0]
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(exit), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	b := astNodeMap[stmt.Test()]
	c := astNodeMap[stmt.Cons().(*parser.ExprStmt).Expr()]
	d := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	AssertEqual(t, d, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(c.Nodes[2]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(d.Nodes[2]), "should be ok")

	ifExit := c.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "N_STMT_IF:exit", nodeToString(ifExit), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr()]
	e := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	cons := stmt.Cons().(*parser.BlockStmt).Body()
	c := astNodeMap[cons[0].(*parser.ExprStmt).Expr()]
	d := astNodeMap[cons[1].(*parser.ExprStmt).Expr()]
	AssertEqual(t, c, d, "should be ok")

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(d.Nodes[6]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(d.Nodes[7]), "should be ok")

	AssertEqual(t, e, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, a.OutSeqEdge().Dst, "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "N_STMT_IF:exit", nodeToString(ifExit), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[1].(*parser.IfStmt)
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr()]
	b := astNodeMap[stmt.Test()]
	e := astNodeMap[stmt.Alt().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()]
	cons := stmt.Cons().(*parser.BlockStmt).Body()
	c := astNodeMap[cons[0].(*parser.ExprStmt).Expr()]
	d := astNodeMap[cons[1].(*parser.ExprStmt).Expr()]
	AssertEqual(t, a, b, "should be ok")
	AssertEqual(t, c, d, "should be ok")

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(d.Nodes[6]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(d.Nodes[7]), "should be ok")

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(e.Nodes[3]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(e.Nodes[4]), "should be ok")

	AssertEqual(t, e, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, a.OutSeqEdge().Dst, "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "N_STMT_IF:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_IfLogic(t *testing.T) {
	ast, symtab, err := compile(`
  if (a && b) c
  else d
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.IfStmt)
	ab := stmt.Test().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]
	c := astNodeMap[stmt.Cons().(*parser.ExprStmt).Expr()]
	d := astNodeMap[stmt.Alt().(*parser.ExprStmt).Expr()]

	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(c.Nodes[2]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(d.Nodes[2]), "should be ok")

	ifExit := a.OutJmpEdge(ET_JMP_F).Dst.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "N_STMT_IF:exit", nodeToString(ifExit), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
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

	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(e.Nodes[2]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(f.Nodes[2]), "should be ok")

	ifExit := e.OutSeqEdge().Dst.Nodes[0]
	AssertEqual(t, "N_STMT_IF:exit", nodeToString(ifExit), "should be ok")
}

func TestCtrlflow_UpdateExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a++
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.UpdateExpr).Arg()]

	AssertEqual(t, "N_EXPR_UPDATE(++):exit", nodeToString(a.Nodes[4]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "N_PROG:exit", nodeToString(a.Nodes[6]), "should be ok")
}

func TestCtrlflow_VarDecStmt(t *testing.T) {
	ast, symtab, err := compile(`
  let a = 1, c = d
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	varDecStmt := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt)
	varDec0 := varDecStmt.DecList()[0].(*parser.VarDec)
	varDec1 := varDecStmt.DecList()[1].(*parser.VarDec)
	a := astNodeMap[varDec0.Id()]
	c := astNodeMap[varDec1.Id()]

	AssertEqual(t, a, c, "should be ok")

	AssertEqual(t, "N_STMT_VAR_DEC:enter", nodeToString(a.Nodes[1]), "should be ok")
	AssertEqual(t, "N_VAR_DEC:enter", nodeToString(a.Nodes[2]), "should be ok")
	AssertEqual(t, "N_VAR_DEC:exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "N_VAR_DEC:enter", nodeToString(a.Nodes[6]), "should be ok")
	AssertEqual(t, "N_VAR_DEC:exit", nodeToString(a.Nodes[9]), "should be ok")
	AssertEqual(t, "N_STMT_VAR_DEC:exit", nodeToString(a.Nodes[10]), "should be ok")
	AssertEqual(t, "N_PROG:exit", nodeToString(a.Nodes[11]), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
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

	AssertEqual(t, "N_STMT_VAR_DEC:exit", nodeToString(b.Nodes[7]), "should be ok")
	AssertEqual(t, "N_EXPR_BIN(<):enter", nodeToString(bc.Nodes[0]), "should be ok")
	AssertEqual(t, "N_EXPR_BIN(<):exit", nodeToString(bc.Nodes[3]), "should be ok")

	AssertEqual(t, bc, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, e, bc.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, "N_STMT_FOR:exit", nodeToString(e.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:enter", nodeToString(e.Nodes[1]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(e.Nodes[3]), "should be ok")

	AssertEqual(t, bu, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "N_EXPR_UPDATE(++):exit", nodeToString(bu.Nodes[7]), "should be ok")

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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	update := stmt.Update().(*parser.UpdateExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	bc := astNodeMap[test.Lhs()]
	bu := astNodeMap[update.Arg()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "N_EXPR_BIN(<):enter", nodeToString(bc.Nodes[0]), "should be ok")
	AssertEqual(t, "N_EXPR_BIN(<):exit", nodeToString(bc.Nodes[3]), "should be ok")

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_FOR:enter", nodeToString(begin.Nodes[1]), "should be ok")

	AssertEqual(t, bc, begin.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, e, bc.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, "N_STMT_FOR:exit", nodeToString(e.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:enter", nodeToString(e.Nodes[1]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(e.Nodes[3]), "should be ok")

	AssertEqual(t, bu, bc.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "N_EXPR_UPDATE(++):exit", nodeToString(bu.Nodes[7]), "should be ok")

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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	bc := astNodeMap[test.Lhs()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_FOR:enter", nodeToString(begin.Nodes[1]), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	b := astNodeMap[test.Lhs()]
	c := astNodeMap[test.Rhs()]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_FOR:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, e, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, e, c.OutJmpEdge(ET_JMP_F).Dst, "should be ok")
	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, b, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:enter", nodeToString(d.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(d.Nodes[4]), "should be ok")

	AssertEqual(t, "N_STMT_FOR:exit", nodeToString(e.Nodes[0]), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.ForStmt)
	begin := ana.Graph().Head
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	d := astNodeMap[body]
	e := astNodeMap[ast.(*parser.Prog).Body()[1].(*parser.ExprStmt).Expr()]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_FOR:enter", nodeToString(begin.Nodes[1]), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.Ident)
	body := stmt.Body().(*parser.ExprStmt).Expr()
	a := astNodeMap[test]
	b := astNodeMap[body]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_WHILE:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst.Nodes[0]
	AssertEqual(t, "N_STMT_WHILE:exit", nodeToString(exit), "should be ok")

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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.Ident)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test]
	b := astNodeMap[body]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_WHILE:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, b.OutJmpEdge(ET_LOOP).Dst, "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst.Nodes[0]
	AssertEqual(t, "N_STMT_WHILE:exit", nodeToString(exit), "should be ok")

	AssertEqual(t, "N_STMT_BLOCK:enter", nodeToString(b.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(b.Nodes[4]), "should be ok")
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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test.Lhs().(*parser.ParenExpr).Expr().(*parser.BinExpr).Lhs()]
	b := astNodeMap[test.Lhs().(*parser.ParenExpr).Expr().(*parser.BinExpr).Rhs()]
	c := astNodeMap[test.Rhs()]
	d := astNodeMap[body]

	AssertEqual(t, "N_PROG:enter", nodeToString(begin.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_WHILE:enter", nodeToString(begin.Nodes[1]), "should be ok")
	AssertEqual(t, b, a, "should be ok")

	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")
	AssertEqual(t, c, a.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, "N_EXPR_BIN(+):exit", nodeToString(a.Nodes[5]), "should be ok")
	AssertEqual(t, "N_EXPR_PAREN:exit", nodeToString(a.Nodes[6]), "should be ok")

	exit := c.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "N_STMT_WHILE:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, a, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:enter", nodeToString(d.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_BLOCK:exit", nodeToString(d.Nodes[4]), "should be ok")
}

func TestCtrlflow_ParenExpr(t *testing.T) {
	ast, symtab, err := compile(`
  (a + b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	a := astNodeMap[ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.ParenExpr).Expr().(*parser.BinExpr).Lhs()]
	AssertEqual(t, 10, len(a.Nodes), "should be ok")
	AssertEqual(t, "N_EXPR_BIN(+):exit", nodeToString(a.Nodes[6]), "should be ok")
	AssertEqual(t, "N_EXPR_PAREN:exit", nodeToString(a.Nodes[7]), "should be ok")
}

func TestCtrlflow_ParenExprLogic(t *testing.T) {
	ast, symtab, err := compile(`
  (a && b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
	ab := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.ParenExpr).Expr().(*parser.BinExpr)
	a := astNodeMap[ab.Lhs()]
	b := astNodeMap[ab.Rhs()]

	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, "N_EXPR_BIN(&&):exit", nodeToString(b.Nodes[1]), "should be ok")

	exit := a.OutJmpEdge(ET_JMP_F).Dst
	AssertEqual(t, "N_EXPR_PAREN:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "N_STMT_EXPR:exit", nodeToString(exit.Nodes[1]), "should be ok")
}

func TestCtrlflow_DoWhile(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
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
	AssertEqual(t, "N_STMT_DO_WHILE:exit", nodeToString(exit.Nodes[0]), "should be ok")
}

func TestCtrlflow_DoWhileLogicOr(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b || c)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
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
	AssertEqual(t, "N_STMT_DO_WHILE:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "N_PROG:exit", nodeToString(exit.Nodes[1]), "should be ok")
}

func TestCtrlflow_DoWhileLogicAnd(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b && c)
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	_, _, astNodeMap := ana.Graph().NodesEdges()
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
	AssertEqual(t, "N_STMT_DO_WHILE:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, "N_PROG:exit", nodeToString(exit.Nodes[1]), "should be ok")

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

	_, _, astNodeMap := ana.Graph().NodesEdges()
	stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	begin := ana.Graph().Head
	test := stmt.Test().(*parser.BinExpr)
	body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	a := astNodeMap[test.Lhs()]
	b := astNodeMap[test.Rhs().(*parser.BinExpr).Lhs()]
	c := astNodeMap[test.Rhs().(*parser.BinExpr).Rhs()]
	d := astNodeMap[body]
	exit := c.OutJmpEdge(ET_JMP_F).Dst

	AssertEqual(t, "N_STMT_WHILE:exit", nodeToString(exit.Nodes[0]), "should be ok")
	AssertEqual(t, a, begin.OutSeqEdge().Dst, "should be ok")

	AssertEqual(t, b, a.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, d, a.OutJmpEdge(ET_JMP_T).Dst, "should be ok")

	AssertEqual(t, c, b.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, exit, b.OutJmpEdge(ET_JMP_F).Dst, "should be ok")

	AssertEqual(t, d, c.OutSeqEdge().Dst, "should be ok")
	AssertEqual(t, a, d.OutJmpEdge(ET_LOOP).Dst, "should be ok")
}

func TestCtrlflow_Continue(t *testing.T) {
	ast, symtab, err := compile(`
  while(a || b && c) {
    continue
    d
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	// fmt.Println(ana.Graph().Dot())

	// _, _, astNodeMap := ana.Graph().NodesEdges()
	// stmt := ast.(*parser.Prog).Body()[0].(*parser.WhileStmt)
	// begin := ana.Graph().Head
	// test := stmt.Test().(*parser.BinExpr)
	// body := stmt.Body().(*parser.BlockStmt).Body()[0].(*parser.ExprStmt).Expr()
	// a := astNodeMap[test.Lhs()]
	// b := astNodeMap[test.Rhs().(*parser.BinExpr).Lhs()]
	// c := astNodeMap[test.Rhs().(*parser.BinExpr).Rhs()]
	// d := astNodeMap[body]
	// exit := c.OutEdge(EK_JMP_FALSE).Dst

	// AssertEqual(t, a, begin.OutEdge(EK_SEQ).Dst, "should be ok")
	// AssertEqual(t, b, a.OutEdge(EK_SEQ).Dst, "should be ok")
	// AssertEqual(t, d, a.OutEdge(EK_JMP_TRUE).Dst, "should be ok")

	// AssertEqual(t, c, b.OutEdge(EK_SEQ).Dst, "should be ok")
	// AssertEqual(t, exit, b.OutEdge(EK_JMP_FALSE).Dst, "should be ok")

	// AssertEqual(t, d, c.OutEdge(EK_SEQ).Dst, "should be ok")
	// AssertEqual(t, a, d.OutEdge(EK_LOOP).Dst, "should be ok")
}

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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_Logic(t *testing.T) {
	ast, symtab, err := compile(`
 a && b
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
b0[label="Prog:enter\nExprStmt:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b6[label="Ident(b)\n"];
b9[label="BinExpr(&&):exit\nExprStmt:exit\nProg:exit\n"];
b0->b6 [xlabel="",color="black"];
b0->b9 [xlabel="F",color="orange"];
b6->b9 [xlabel="",color="black"];
b9->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LogicMix(t *testing.T) {
	ast, symtab, err := compile(`
 a && b || c
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
b0[label="Prog:enter\nExprStmt:enter\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(a)\n"];
b11[label="Ident(c)\n"];
b14[label="BinExpr(||):exit\nExprStmt:exit\nProg:exit\n"];
b7[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b11 [xlabel="F",color="orange"];
b0->b7 [xlabel="",color="black"];
b11->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
b7->b11 [xlabel="",color="black"];
b7->b14 [xlabel="T",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nIfStmt:enter\nIdent(b)\n"];
b13[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\n"];
b19[label="IfStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b8[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\n"];
b0->b13 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b13->b19 [xlabel="",color="black"];
b19->final [xlabel="",color="black"];
b8->b19 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nIfStmt:enter\nIdent(b)\n"];
b17[label="ExprStmt:enter\nIdent(e)\nExprStmt:exit\n"];
b23[label="IfStmt:exit\nExprStmt:enter\nIdent(f)\nExprStmt:exit\nProg:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b17 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b17->b23 [xlabel="",color="black"];
b23->final [xlabel="",color="black"];
b8->b23 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nIfStmt:enter\nIdent(b)\n"];
b17[label="BlockStmt:enter\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nBlockStmt:exit\n"];
b24[label="IfStmt:exit\nExprStmt:enter\nIdent(f)\nExprStmt:exit\nProg:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b17 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b17->b24 [xlabel="",color="black"];
b24->final [xlabel="",color="black"];
b8->b24 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_IfLogic(t *testing.T) {
	ast, symtab, err := compile(`
  if (a && b) c
  else d
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
b0[label="Prog:enter\nIfStmt:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b10[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\n"];
b15[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\n"];
b21[label="IfStmt:exit\nProg:exit\n"];
b6[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b15 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b10->b21 [xlabel="",color="black"];
b15->b21 [xlabel="",color="black"];
b21->final [xlabel="",color="black"];
b6->b10 [xlabel="",color="black"];
b6->b15 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nIfStmt:enter\nBinExpr(||):enter\nIdent(b)\n"];
b11[label="Ident(d)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b18[label="ExprStmt:enter\nIdent(e)\nExprStmt:exit\n"];
b23[label="ExprStmt:enter\nIdent(f)\nExprStmt:exit\n"];
b29[label="IfStmt:exit\nProg:exit\n"];
b9[label="BinExpr(&&):enter\nIdent(c)\n"];
b0->b18 [xlabel="T",color="orange"];
b0->b9 [xlabel="",color="black"];
b11->b18 [xlabel="",color="black"];
b11->b23 [xlabel="F",color="orange"];
b18->b29 [xlabel="",color="black"];
b23->b29 [xlabel="",color="black"];
b29->final [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
b9->b23 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_UpdateExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a++
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
b0[label="Prog:enter\nExprStmt:enter\nUpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_VarDecStmt(t *testing.T) {
	ast, symtab, err := compile(`
  let a = 1, c = d
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDec:enter\nIdent(c)\nIdent(d)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(b)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b12[label="BinExpr(<):enter\nIdent(b)\nIdent(c)\nBinExpr(<):exit\n"];
b22[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
b29[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b0->b12 [xlabel="",color="black"];
b12->b22 [xlabel="",color="black"];
b12->b29 [xlabel="F",color="orange"];
b22:s->b12:ne [xlabel="L",color="orange"];
b29->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nForStmt:enter\n"];
b14[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
b21[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b4[label="BinExpr(<):enter\nIdent(b)\nIdent(c)\nBinExpr(<):exit\n"];
b0->b4 [xlabel="",color="black"];
b14:s->b4:ne [xlabel="L",color="orange"];
b21->final [xlabel="",color="black"];
b4->b14 [xlabel="",color="black"];
b4->b21 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nForStmt:enter\n"];
b10[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b17[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b4[label="BinExpr(<):enter\nIdent(b)\nIdent(c)\nBinExpr(<):exit\n"];
b0->b4 [xlabel="",color="black"];
b10:s->b4:ne [xlabel="L",color="orange"];
b17->final [xlabel="",color="black"];
b4->b10 [xlabel="",color="black"];
b4->b17 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
  `, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nForStmt:enter\n"];
b10[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b17[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b4[label="BinExpr(&&):enter\nIdent(b)\n"];
b6[label="Ident(c)\nBinExpr(&&):exit\n"];
b0->b4 [xlabel="",color="black"];
b10:s->b4:ne [xlabel="L",color="orange"];
b17->final [xlabel="",color="black"];
b4->b17 [xlabel="F",color="orange"];
b4->b6 [xlabel="",color="black"];
b6->b10 [xlabel="",color="black"];
b6->b17 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nForStmt:enter\n"];
b11[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b11->final [xlabel="",color="red"];
b4->b11 [xlabel="",color="red"];
b4:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ForStmtLitTest(t *testing.T) {
	ast, symtab, err := compile(`
  for (; 1 ;) {
    d;
  }
  e;
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
b0[label="Prog:enter\nForStmt:enter\n"];
b12[label="ForStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b4[label="NumLit(1)\n"];
b5[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b4->b12 [xlabel="",color="red"];
b4->b5 [xlabel="",color="black"];
b5:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_While(t *testing.T) {
	ast, symtab, err := compile(`
  while(a) b;
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b11[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b11->final [xlabel="",color="black"];
b4->b11 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nWhileStmt:enter\n"];
b12[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="BlockStmt:enter\nExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b4->b12 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nWhileStmt:enter\n"];
b14[label="Ident(c)\nBinExpr(||):exit\n"];
b18[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b25[label="WhileStmt:exit\nProg:exit\n"];
b4[label="BinExpr(||):enter\nParenExpr:enter\nBinExpr(+):enter\nIdent(a)\nIdent(b)\nBinExpr(+):exit\nParenExpr:exit\n"];
b0->b4 [xlabel="",color="black"];
b14->b18 [xlabel="",color="black"];
b14->b25 [xlabel="F",color="orange"];
b18:s->b4:ne [xlabel="L",color="orange"];
b25->final [xlabel="",color="black"];
b4->b14 [xlabel="",color="black"];
b4->b18 [xlabel="T",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_WhileTestLit(t *testing.T) {
	ast, symtab, err := compile(`
  while(1) {
    d;
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b12[label="WhileStmt:exit\nProg:exit\n"];
b4[label="NumLit(1)\nBlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="red"];
b4->b12 [xlabel="",color="red"];
b4:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nWhileStmt:enter\n"];
b15[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b22[label="WhileStmt:exit\nProg:exit\n"];
b4[label="BinExpr(||):enter\nIdent(a)\n"];
b6[label="BinExpr(&&):enter\nIdent(b)\n"];
b8[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b0->b4 [xlabel="",color="black"];
b15:s->b4:ne [xlabel="L",color="orange"];
b22->final [xlabel="",color="black"];
b4->b15 [xlabel="T",color="orange"];
b4->b6 [xlabel="",color="black"];
b6->b22 [xlabel="F",color="orange"];
b6->b8 [xlabel="",color="black"];
b8->b15 [xlabel="",color="black"];
b8->b22 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_WhileCont(t *testing.T) {
	ast, symtab, err := compile(`
while (a) {
  continue
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b10[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b15[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="BlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b10:s->b4:ne [xlabel="L",color="red"];
b15->final [xlabel="",color="black"];
b4->b15 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5->b10 [xlabel="",color="red"];
b5:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_WhileLitCont(t *testing.T) {
	ast, symtab, err := compile(`
while (1) {
  continue
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b10[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b15[label="WhileStmt:exit\nProg:exit\n"];
b4[label="NumLit(1)\nBlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b10->b15 [xlabel="",color="red"];
b10:s->b4:ne [xlabel="L",color="red"];
b15->final [xlabel="",color="red"];
b4->b10 [xlabel="",color="red"];
b4:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ParenExpr(t *testing.T) {
	ast, symtab, err := compile(`
  (a + b)
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
b0[label="Prog:enter\nExprStmt:enter\nParenExpr:enter\nBinExpr(+):enter\nIdent(a)\nIdent(b)\nBinExpr(+):exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ParenExprLogic(t *testing.T) {
	ast, symtab, err := compile(`
  (a && b)
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
b0[label="Prog:enter\nExprStmt:enter\nParenExpr:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b11[label="ParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b7[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b11 [xlabel="F",color="orange"];
b0->b7 [xlabel="",color="black"];
b11->final [xlabel="",color="black"];
b7->b11 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhile(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b12[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nIdent(b)\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b4->b12 [xlabel="F",color="orange"];
b4:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhileLogicOr(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b || c)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b12[label="Ident(c)\nBinExpr(||):exit\n"];
b17[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nBinExpr(||):enter\nIdent(b)\n"];
b0->b4 [xlabel="",color="black"];
b12->b17 [xlabel="F",color="orange"];
b12:s->b4:ne [xlabel="L",color="orange"];
b17->final [xlabel="",color="black"];
b4->b12 [xlabel="",color="black"];
b4:s->b4:ne [xlabel="T,L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhileLogicAnd(t *testing.T) {
	ast, symtab, err := compile(`
  do { a } while(b && c)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b12[label="Ident(c)\nBinExpr(&&):exit\n"];
b17[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nBinExpr(&&):enter\nIdent(b)\n"];
b0->b4 [xlabel="",color="black"];
b12->b17 [xlabel="F",color="orange"];
b12:s->b4:ne [xlabel="L",color="orange"];
b17->final [xlabel="",color="black"];
b4->b12 [xlabel="",color="black"];
b4->b17 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhileLit(t *testing.T) {
	ast, symtab, err := compile(`
  do {
    d;
  } while(1)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b12[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\nNumLit(1)\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b4->b12 [xlabel="",color="red"];
b4:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhileCont(t *testing.T) {
	ast, symtab, err := compile(`
do {
  continue
  c
} while(a)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b13[label="Ident(a)\n"];
b15[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b9[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b13->b15 [xlabel="F",color="orange"];
b13:s->b4:ne [xlabel="L",color="orange"];
b15->final [xlabel="",color="black"];
b4:s->b13:ne [xlabel="L",color="orange"];
b4->b9 [xlabel="",color="red"];
b9->b13 [xlabel="",color="red"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_DoWhileLitCont(t *testing.T) {
	ast, symtab, err := compile(`
do {
  continue
  c
} while(1)
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
b0[label="Prog:enter\nDoWhileStmt:enter\n"];
b13[label="NumLit(1)\n"];
b15[label="DoWhileStmt:exit\nProg:exit\n"];
b4[label="BlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b9[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b13->b15 [xlabel="",color="red"];
b13:s->b4:ne [xlabel="L",color="orange"];
b15->final [xlabel="",color="black"];
b4:s->b13:ne [xlabel="L",color="orange"];
b4->b9 [xlabel="",color="red"];
b9->b13 [xlabel="",color="red"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nWhileStmt:enter\n"];
b13[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b18[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="Ident(a)\n"];
b7[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b0->b6 [xlabel="",color="black"];
b13:s->b6:ne [xlabel="L",color="red"];
b18->final [xlabel="",color="black"];
b6->b18 [xlabel="F",color="orange"];
b6->b7 [xlabel="",color="black"];
b7->b13 [xlabel="",color="red"];
b7:s->b6:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
  `, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nWhileStmt:enter\n"];
b10[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b17[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b23[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b28[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="BinExpr(||):enter\nIdent(a)\n"];
b8[label="BinExpr(&&):enter\nIdent(b)\n"];
b0->b6 [xlabel="",color="black"];
b10->b17 [xlabel="",color="black"];
b10->b28 [xlabel="F",color="orange"];
b17->b23 [xlabel="",color="red"];
b17:s->b6:ne [xlabel="L",color="orange"];
b23:s->b6:ne [xlabel="L",color="red"];
b28->final [xlabel="",color="black"];
b6->b17 [xlabel="T",color="orange"];
b6->b8 [xlabel="",color="black"];
b8->b10 [xlabel="",color="black"];
b8->b28 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nWhileStmt:enter\n"];
b10[label="Ident(b)\n"];
b11[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b17[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b22[label="WhileStmt:exit\nBlockStmt:exit\n"];
b25[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="Ident(a)\n"];
b7[label="BlockStmt:enter\nWhileStmt:enter\n"];
b0->b6 [xlabel="",color="black"];
b10->b11 [xlabel="",color="black"];
b10->b22 [xlabel="F",color="orange"];
b11->b17 [xlabel="",color="red"];
b11:s->b6:ne [xlabel="L",color="orange"];
b17:s->b10:ne [xlabel="L",color="red"];
b22:s->b6:ne [xlabel="L",color="orange"];
b25->final [xlabel="",color="black"];
b6->b25 [xlabel="F",color="orange"];
b6->b7 [xlabel="",color="black"];
b7->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
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

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nDoWhileStmt:enter\n"];
b12[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b18[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\nIdent(c)\n"];
b24[label="DoWhileStmt:exit\nBlockStmt:exit\n"];
b26[label="Ident(b)\n"];
b28[label="DoWhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="BlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nDoWhileStmt:enter\n"];
b0->b6 [xlabel="",color="black"];
b12->b18 [xlabel="",color="red"];
b12:s->b26:ne [xlabel="L",color="orange"];
b18:s->b12:ne [xlabel="L",color="red"];
b18->b24 [xlabel="F",color="red"];
b24->b26 [xlabel="",color="red"];
b26->b28 [xlabel="F",color="orange"];
b26:s->b6:ne [xlabel="L",color="orange"];
b28->final [xlabel="",color="black"];
b6->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForBasic(t *testing.T) {
	ast, symtab, err := compile(`
for (let a = 1; a < 10; a++) {
  continue
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
b0[label="Prog:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b12[label="BinExpr(<):enter\nIdent(a)\nNumLit(10)\nBinExpr(<):exit\n"];
b18[label="UpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\n"];
b22[label="BlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b27[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b32[label="ForStmt:exit\nProg:exit\n"];
b0->b12 [xlabel="",color="black"];
b12->b22 [xlabel="",color="black"];
b12->b32 [xlabel="F",color="orange"];
b18:s->b12:ne [xlabel="L",color="orange"];
b22:s->b18:ne [xlabel="L",color="orange"];
b22->b27 [xlabel="",color="red"];
b27->b18 [xlabel="",color="red"];
b32->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForListTestLabel(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for (let a = 1; 1; a++) {
  continue LabelA
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b14[label="NumLit(1)\n"];
b15[label="UpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\n"];
b19[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b25[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b30[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b14->b19 [xlabel="",color="black"];
b14->b30 [xlabel="",color="red"];
b15:s->b14:ne [xlabel="L",color="orange"];
b19:s->b15:ne [xlabel="L",color="orange"];
b19->b25 [xlabel="",color="red"];
b25->b15 [xlabel="",color="red"];
b30->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForBasicLabel(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for (let a = 1; a < 10; a++) {
  continue LabelA
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b14[label="BinExpr(<):enter\nIdent(a)\nNumLit(10)\nBinExpr(<):exit\n"];
b20[label="UpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\n"];
b24[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b30[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b35[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b14->b24 [xlabel="",color="black"];
b14->b35 [xlabel="F",color="orange"];
b20:s->b14:ne [xlabel="L",color="orange"];
b24:s->b20:ne [xlabel="L",color="orange"];
b24->b30 [xlabel="",color="red"];
b30->b20 [xlabel="",color="red"];
b35->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueFor(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for (let a = 1; a < 10; a++) {
  for (let b = a; b < 10; b++) {
    continue LabelA
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b14[label="BinExpr(<):enter\nIdent(a)\nNumLit(10)\nBinExpr(<):exit\n"];
b20[label="UpdateExpr(++):enter\nIdent(a)\nUpdateExpr(++):exit\n"];
b24[label="BlockStmt:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(b)\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\n"];
b35[label="BinExpr(<):enter\nIdent(b)\nNumLit(10)\nBinExpr(<):exit\n"];
b45[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b51[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
b56[label="ForStmt:exit\nBlockStmt:exit\n"];
b59[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b14->b24 [xlabel="",color="black"];
b14->b59 [xlabel="F",color="orange"];
b20:s->b14:ne [xlabel="L",color="orange"];
b24->b35 [xlabel="",color="black"];
b35->b45 [xlabel="",color="black"];
b35->b56 [xlabel="F",color="orange"];
b45:s->b20:ne [xlabel="L",color="orange"];
b45->b51 [xlabel="",color="red"];
b51:s->b35:ne [xlabel="L",color="red"];
b56->b20 [xlabel="",color="black"];
b59->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForNoUpdate(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for (let a = 1; a < 10; ) {
  for (let b = a; b < 10; b++) {
    continue LabelA
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b14[label="BinExpr(<):enter\nIdent(a)\nNumLit(10)\nBinExpr(<):exit\n"];
b20[label="BlockStmt:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(b)\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\n"];
b31[label="BinExpr(<):enter\nIdent(b)\nNumLit(10)\nBinExpr(<):exit\n"];
b41[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b47[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
b52[label="ForStmt:exit\nBlockStmt:exit\n"];
b55[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b14->b20 [xlabel="",color="black"];
b14->b55 [xlabel="F",color="orange"];
b20->b31 [xlabel="",color="black"];
b31->b41 [xlabel="",color="black"];
b31->b52 [xlabel="F",color="orange"];
b41:s->b14:ne [xlabel="L",color="orange"];
b41->b47 [xlabel="",color="red"];
b47:s->b31:ne [xlabel="L",color="red"];
b52:s->b14:ne [xlabel="L",color="orange"];
b55->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueForNoTest(t *testing.T) {
	ast, symtab, err := compile(`
LabelA: for (let a = 1; ; ) {
  for (let b = a; b < 10; b++) {
    continue LabelA
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nNumLit(1)\nVarDec:exit\nVarDecStmt:exit\n"];
b14[label="BlockStmt:enter\nForStmt:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(b)\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\n"];
b25[label="BinExpr(<):enter\nIdent(b)\nNumLit(10)\nBinExpr(<):exit\n"];
b35[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b41[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nUpdateExpr(++):enter\nIdent(b)\nUpdateExpr(++):exit\n"];
b46[label="ForStmt:exit\nBlockStmt:exit\n"];
b49[label="ForStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b14->b25 [xlabel="",color="black"];
b25->b35 [xlabel="",color="black"];
b25->b46 [xlabel="F",color="orange"];
b35:s->b14:ne [xlabel="L",color="orange"];
b35->b41 [xlabel="",color="red"];
b41:s->b25:ne [xlabel="L",color="red"];
b46:s->b14:ne [xlabel="L",color="orange"];
b46->b49 [xlabel="",color="red"];
b49->final [xlabel="",color="red"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ForIn(t *testing.T) {
	ast, symtab, err := compile(`
for (a in b) {
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
b0[label="Prog:enter\nForInOfStmt:enter\n"];
b12[label="ForInOfStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\nIdent(b)\n"];
b6[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b4->b12 [xlabel="F",color="orange"];
b4->b6 [xlabel="",color="black"];
b6:s->b4:ne [xlabel="L",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ForInLet(t *testing.T) {
	ast, symtab, err := compile(`
for (let a in b) {
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
b0[label="Prog:enter\nForInOfStmt:enter\n"];
b12[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b18[label="ForInOfStmt:exit\nProg:exit\n"];
b4[label="VarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\nIdent(b)\n"];
b0->b4 [xlabel="",color="black"];
b12:s->b4:ne [xlabel="L",color="orange"];
b18->final [xlabel="",color="black"];
b4->b12 [xlabel="",color="black"];
b4->b18 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nIdent(s)\nExprStmt:exit\nLabelStmt:enter\nIdent(LabelA)\nForInOfStmt:enter\n"];
b11[label="BlockStmt:enter\nForInOfStmt:enter\n"];
b14[label="Ident(c)\nIdent(d)\n"];
b16[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(&&):enter\nIdent(e)\n"];
b21[label="Ident(f)\nBinExpr(&&):exit\n"];
b25[label="BlockStmt:enter\nContStmt:enter\nIdent(LabelA)\nContStmt:exit\n"];
b31[label="BlockStmt:exit\n"];
b33[label="IfStmt:exit\nExprStmt:enter\nIdent(g)\nExprStmt:exit\nBlockStmt:exit\n"];
b38[label="ForInOfStmt:exit\nExprStmt:enter\nIdent(h)\nExprStmt:exit\nBlockStmt:exit\n"];
b43[label="ForInOfStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b9[label="Ident(a)\nIdent(b)\n"];
b0->b9 [xlabel="",color="black"];
b11->b14 [xlabel="",color="black"];
b14->b16 [xlabel="",color="black"];
b14->b38 [xlabel="F",color="orange"];
b16->b21 [xlabel="",color="black"];
b16->b33 [xlabel="F",color="orange"];
b21->b25 [xlabel="",color="black"];
b21->b33 [xlabel="F",color="orange"];
b25->b31 [xlabel="",color="red"];
b25:s->b9:ne [xlabel="L",color="orange"];
b31->b33 [xlabel="",color="red"];
b33:s->b14:ne [xlabel="L",color="orange"];
b38:s->b9:ne [xlabel="L",color="orange"];
b43->final [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
b9->b43 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ContinueNoLabel(t *testing.T) {
	ast, symtab, err := compile(`
while(a) {
  if (a > 10) {
    continue;
    b
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b14[label="BlockStmt:enter\nContStmt:enter\nContStmt:exit\n"];
b19[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\n"];
b24[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
b33[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
b0->b4 [xlabel="",color="black"];
b14->b19 [xlabel="",color="red"];
b14:s->b4:ne [xlabel="L",color="orange"];
b19->b24 [xlabel="",color="red"];
b24:s->b4:ne [xlabel="L",color="orange"];
b33->final [xlabel="",color="black"];
b4->b33 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5->b14 [xlabel="",color="black"];
b5->b24 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nWhileStmt:enter\n"];
b14[label="BrkStmt:enter\nBrkStmt:exit\n"];
b18[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
b27[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
b0->b4 [xlabel="",color="black"];
b14->b18 [xlabel="",color="red"];
b14->b27 [xlabel="U",color="orange"];
b18:s->b4:ne [xlabel="L",color="orange"];
b27->final [xlabel="",color="black"];
b4->b27 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5->b14 [xlabel="",color="black"];
b5->b18 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakNoLabelFor(t *testing.T) {
	ast, symtab, err := compile(`
for (; a; ) {
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
b0[label="Prog:enter\nForStmt:enter\n"];
b14[label="BlockStmt:enter\nBrkStmt:enter\nBrkStmt:exit\n"];
b18[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b23[label="IfStmt:exit\nExprStmt:enter\nUpdateExpr(--):enter\nIdent(a)\nUpdateExpr(--):exit\nExprStmt:exit\nBlockStmt:exit\n"];
b32[label="ForStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="BlockStmt:enter\nIfStmt:enter\nBinExpr(>):enter\nIdent(a)\nNumLit(10)\nBinExpr(>):exit\n"];
b0->b4 [xlabel="",color="black"];
b14->b18 [xlabel="",color="red"];
b14->b32 [xlabel="U",color="orange"];
b18->b23 [xlabel="",color="red"];
b23:s->b4:ne [xlabel="L",color="orange"];
b32->final [xlabel="",color="black"];
b4->b32 [xlabel="F",color="orange"];
b4->b5 [xlabel="",color="black"];
b5->b14 [xlabel="",color="black"];
b5->b23 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_BreakLabelWhile(t *testing.T) {
	ast, symtab, err := compile(`
  LabelA: while(a) {
    for( ; b; ) {
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nWhileStmt:enter\n"];
b10[label="Ident(b)\n"];
b11[label="BlockStmt:enter\nBrkStmt:enter\nIdent(LabelA)\nBrkStmt:exit\n"];
b16[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b21[label="ForStmt:exit\nBlockStmt:exit\n"];
b24[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="Ident(a)\n"];
b7[label="BlockStmt:enter\nForStmt:enter\n"];
b0->b6 [xlabel="",color="black"];
b10->b11 [xlabel="",color="black"];
b10->b21 [xlabel="F",color="orange"];
b11->b16 [xlabel="",color="red"];
b11->b24 [xlabel="U",color="orange"];
b16:s->b10:ne [xlabel="L",color="red"];
b21:s->b6:ne [xlabel="L",color="orange"];
b24->final [xlabel="",color="black"];
b6->b24 [xlabel="F",color="orange"];
b6->b7 [xlabel="",color="black"];
b7->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nLabelStmt:enter\nIdent(LabelA)\nWhileStmt:enter\n"];
b10[label="Ident(b)\n"];
b11[label="BlockStmt:enter\nBrkStmt:enter\nBrkStmt:exit\n"];
b15[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b20[label="ForStmt:exit\nBlockStmt:exit\n"];
b23[label="WhileStmt:exit\nLabelStmt:exit\nProg:exit\n"];
b6[label="Ident(a)\n"];
b7[label="BlockStmt:enter\nForStmt:enter\n"];
b0->b6 [xlabel="",color="black"];
b10->b11 [xlabel="",color="black"];
b10->b20 [xlabel="F",color="orange"];
b11->b15 [xlabel="",color="red"];
b11->b20 [xlabel="U",color="orange"];
b15:s->b10:ne [xlabel="L",color="red"];
b20:s->b6:ne [xlabel="L",color="orange"];
b23->final [xlabel="",color="black"];
b6->b23 [xlabel="F",color="orange"];
b6->b7 [xlabel="",color="black"];
b7->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nCallExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nNumLit(1)\nNumLit(2)\nBinExpr(&&):enter\nIdent(a)\n"];
b10[label="Ident(b)\nBinExpr(&&):exit\n"];
b14[label="CallExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b10 [xlabel="",color="black"];
b0->b14 [xlabel="F",color="orange"];
b10->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CallExprArgsBin(t *testing.T) {
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
b0[label="Prog:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fn)\nNumLit(1)\nBinExpr(||):enter\nIdent(a)\n"];
b11[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b18[label="BinExpr(&&):enter\nIdent(d)\n"];
b20[label="Ident(e)\nBinExpr(&&):exit\n"];
b24[label="CallExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="BinExpr(&&):enter\nIdent(b)\n"];
b0->b18 [xlabel="T",color="orange"];
b0->b9 [xlabel="",color="black"];
b11->b18 [xlabel="",color="black"];
b18->b20 [xlabel="",color="black"];
b18->b24 [xlabel="F",color="orange"];
b20->b24 [xlabel="",color="black"];
b24->final [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
b9->b18 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(foo)\nBinExpr(??):enter\nNullLit\n"];
b12[label="VarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b8[label="StrLit\nBinExpr(??):exit\n"];
b0->b12 [xlabel="T",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nFnDec\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b9[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b9->final [xlabel="",color="black"];
initial->b9 [xlabel="",color="black"];
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
b10[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\n"];
b14[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b15[label="FnDec:exit\n"];
b10->b15 [xlabel="",color="red"];
b14->b10 [xlabel="",color="red"];
b14->b15 [xlabel="U",color="orange"];
b15->final [xlabel="",color="black"];
initial->b14 [xlabel="",color="black"];
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
b11[label="Ident(b)\nBinExpr(??):exit\n"];
b15[label="RetStmt:exit\n"];
b16[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b20[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nRetStmt:enter\nBinExpr(??):enter\nIdent(a)\n"];
b21[label="FnDec:exit\n"];
b11->b15 [xlabel="",color="black"];
b15->b16 [xlabel="",color="red"];
b15->b21 [xlabel="U",color="orange"];
b16->b21 [xlabel="",color="red"];
b20->b11 [xlabel="",color="black"];
b20->b15 [xlabel="T",color="orange"];
b21->final [xlabel="",color="black"];
initial->b20 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ReturnArgBinNested(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  a;
  return a && b || c;
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
b12[label="Ident(b)\nBinExpr(&&):exit\n"];
b16[label="Ident(c)\nBinExpr(||):exit\n"];
b20[label="RetStmt:exit\n"];
b21[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b25[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nRetStmt:enter\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(a)\n"];
b26[label="FnDec:exit\n"];
b12->b16 [xlabel="",color="black"];
b12->b20 [xlabel="T",color="orange"];
b16->b20 [xlabel="",color="black"];
b20->b21 [xlabel="",color="red"];
b20->b26 [xlabel="U",color="orange"];
b21->b26 [xlabel="",color="red"];
b25->b12 [xlabel="",color="black"];
b25->b16 [xlabel="F",color="orange"];
b26->final [xlabel="",color="black"];
initial->b25 [xlabel="",color="black"];
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

	fn := ast.(*parser.Prog).Body()[0]
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b8[label="FnDec:enter\nIdent(f)\nIdent(a)\nIdent(b)\nBlockStmt:enter\nBlockStmt:exit\nFnDec:exit\n"];
b8->final [xlabel="",color="black"];
initial->b8 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_FnExpr(t *testing.T) {
	ast, symtab, err := compile(`
fn = function f(a, b) { c }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.AssignExpr).Rhs()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b11[label="FnDec:enter\nIdent(f)\nIdent(a)\nIdent(b)\nBlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b11->final [xlabel="",color="black"];
initial->b11 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nIdent(b)\nIdent(c)\nBinExpr(??):enter\nIdent(d)\n"];
b11[label="Ident(e)\nBinExpr(??):exit\n"];
b15[label="ArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b11 [xlabel="",color="black"];
b0->b15 [xlabel="T",color="orange"];
b11->b15 [xlabel="",color="black"];
b15->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(c)\nIdent(d)\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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

	rhs := ast.(*parser.Prog).Body()[0].(*parser.ExprStmt).Expr().(*parser.AssignExpr).Rhs().(*parser.ObjLit)
	fn := rhs.Props()[0].(*parser.Prop).Val()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b7[label="FnDec:enter\nIdent(a)\nIdent(b)\nBlockStmt:enter\nBlockStmt:exit\nFnDec:exit\n"];
b7->final [xlabel="",color="black"];
initial->b7 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nDebugStmt\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nNullLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nNullLit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nBoolLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nBoolLit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nNumLit(1)\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nNumLit(1.1)\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nStrLit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nStrLit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nNumLit(1)\nNumLit(2)\nNumLit(3)\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nArrLit:enter\nNumLit(1)\nNumLit(2)\nNumLit(3)\nArrLit:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:enter\nNumLit(4)\nNumLit(5)\nNumLit(6)\nArrLit:exit\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitArrSpread(t *testing.T) {
	ast, symtab, err := compile(`
  a = [1, 2, ...d];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nArrLit:enter\nNumLit(1)\nNumLit(2)\nSpread:enter\nIdent(d)\nSpread:exit\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitObj(t *testing.T) {
	ast, symtab, err := compile(`
a = { b: { c: 1 }, ...d };
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nSpread:enter\nIdent(d)\nSpread:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitObjSpread(t *testing.T) {
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nExprStmt:enter\nParenExpr:enter\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nObjLit:enter\nProp:enter\nIdent(b)\nObjLit:enter\nProp:enter\nIdent(c)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(d)\nArrLit:enter\nNumLit(1)\nObjLit:enter\nProp:enter\nIdent(f)\nNumLit(2)\nProp:exit\nObjLit:exit\nArrLit:exit\nProp:exit\nObjLit:exit\nProp:exit\nObjLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_LitRegexp(t *testing.T) {
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nRegLit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b19[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b28[label="SwitchCase:enter\nNumLit(2)\n"];
b32[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b41[label="SwitchCase:enter\nNumLit(3)\n"];
b45[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b53[label="SwitchCase:enter\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b28 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b19->b32 [xlabel="",color="black"];
b28->b19 [xlabel="",color="black"];
b28->b41 [xlabel="F",color="orange"];
b32->b45 [xlabel="",color="black"];
b41->b32 [xlabel="",color="black"];
b41->b53 [xlabel="F",color="orange"];
b45->final [xlabel="",color="black"];
b53->b45 [xlabel="",color="black"];
b6->b19 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b19[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth21)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth22)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nSwitchCase:exit\n"];
b36[label="SwitchCase:enter\nNumLit(2)\n"];
b40[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b49[label="SwitchCase:enter\nNumLit(3)\n"];
b53[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b61[label="SwitchCase:enter\n"];
b0->b36 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b19->b40 [xlabel="",color="black"];
b36->b19 [xlabel="",color="black"];
b36->b49 [xlabel="F",color="orange"];
b40->b53 [xlabel="",color="black"];
b49->b40 [xlabel="",color="black"];
b49->b61 [xlabel="F",color="orange"];
b53->final [xlabel="",color="black"];
b61->b53 [xlabel="",color="black"];
b6->b19 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b12[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth21)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth22)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nSwitchCase:exit\n"];
b29[label="SwitchCase:enter\nNumLit(2)\n"];
b33[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b42[label="SwitchCase:enter\nNumLit(3)\n"];
b46[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b54[label="SwitchCase:enter\n"];
b9[label="SwitchCase:exit\n"];
b0->b29 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b12->b33 [xlabel="",color="black"];
b29->b12 [xlabel="",color="black"];
b29->b42 [xlabel="F",color="orange"];
b33->b46 [xlabel="",color="black"];
b42->b33 [xlabel="",color="black"];
b42->b54 [xlabel="F",color="orange"];
b46->final [xlabel="",color="black"];
b54->b46 [xlabel="",color="black"];
b9->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\n"];
b18[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b27[label="SwitchCase:enter\nNumLit(1)\n"];
b31[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b40[label="SwitchCase:enter\nNumLit(2)\n"];
b43[label="SwitchStmt:exit\nProg:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b27 [xlabel="",color="black"];
b18->b31 [xlabel="",color="black"];
b27->b18 [xlabel="",color="black"];
b27->b40 [xlabel="F",color="orange"];
b31->b43 [xlabel="",color="black"];
b40->b31 [xlabel="",color="black"];
b40->b6 [xlabel="F",color="orange"];
b43->final [xlabel="",color="black"];
b6->b18 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b19[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b27[label="SwitchCase:enter\n"];
b31[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b40[label="SwitchCase:enter\nNumLit(2)\n"];
b44[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b53[label="SwitchCase:enter\nNumLit(3)\n"];
b56[label="SwitchStmt:exit\nProg:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b27 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b19->b31 [xlabel="",color="black"];
b27->b40 [xlabel="",color="black"];
b31->b44 [xlabel="",color="black"];
b40->b31 [xlabel="",color="black"];
b40->b53 [xlabel="F",color="orange"];
b44->b56 [xlabel="",color="black"];
b53->b19 [xlabel="F",color="orange"];
b53->b44 [xlabel="",color="black"];
b56->final [xlabel="",color="black"];
b6->b19 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b12[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b20[label="SwitchCase:enter\n"];
b24[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b33[label="SwitchCase:enter\nNumLit(2)\n"];
b37[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b46[label="SwitchCase:enter\nNumLit(3)\n"];
b49[label="SwitchStmt:exit\nProg:exit\n"];
b9[label="SwitchCase:exit\n"];
b0->b20 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b12->b24 [xlabel="",color="black"];
b20->b33 [xlabel="",color="black"];
b24->b37 [xlabel="",color="black"];
b33->b24 [xlabel="",color="black"];
b33->b46 [xlabel="F",color="orange"];
b37->b49 [xlabel="",color="black"];
b46->b12 [xlabel="F",color="orange"];
b46->b37 [xlabel="",color="black"];
b49->final [xlabel="",color="black"];
b9->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b18[label="SwitchCase:exit\n"];
b21[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b30[label="SwitchCase:enter\nNumLit(2)\n"];
b34[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b43[label="SwitchCase:enter\nNumLit(3)\n"];
b46[label="SwitchStmt:exit\n"];
b47[label="Prog:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b0->b30 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b18->b21 [xlabel="",color="red"];
b21->b34 [xlabel="",color="black"];
b30->b21 [xlabel="",color="black"];
b30->b43 [xlabel="F",color="orange"];
b34->b46 [xlabel="",color="black"];
b43->b34 [xlabel="",color="black"];
b43->b46 [xlabel="F",color="orange"];
b46->b47 [xlabel="",color="black"];
b47->final [xlabel="",color="black"];
b6->b18 [xlabel="",color="red"];
b6->b46 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b18[label="SwitchCase:exit\n"];
b21[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b29[label="SwitchCase:enter\n"];
b33[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b42[label="SwitchCase:enter\nNumLit(2)\n"];
b46[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b55[label="SwitchCase:enter\nNumLit(3)\n"];
b58[label="SwitchStmt:exit\n"];
b59[label="Prog:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b0->b29 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b18->b21 [xlabel="",color="red"];
b21->b33 [xlabel="",color="black"];
b29->b42 [xlabel="",color="black"];
b33->b46 [xlabel="",color="black"];
b42->b33 [xlabel="",color="black"];
b42->b55 [xlabel="F",color="orange"];
b46->b58 [xlabel="",color="black"];
b55->b21 [xlabel="F",color="orange"];
b55->b46 [xlabel="",color="black"];
b58->b59 [xlabel="",color="black"];
b59->final [xlabel="",color="black"];
b6->b18 [xlabel="",color="red"];
b6->b58 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b18[label="BlockStmt:exit\nSwitchCase:exit\n"];
b23[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b31[label="SwitchCase:enter\n"];
b35[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b44[label="SwitchCase:enter\nNumLit(2)\n"];
b48[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b57[label="SwitchCase:enter\nNumLit(3)\n"];
b6[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b60[label="SwitchStmt:exit\n"];
b61[label="Prog:exit\n"];
b0->b31 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b18->b23 [xlabel="",color="red"];
b23->b35 [xlabel="",color="black"];
b31->b44 [xlabel="",color="black"];
b35->b48 [xlabel="",color="black"];
b44->b35 [xlabel="",color="black"];
b44->b57 [xlabel="F",color="orange"];
b48->b60 [xlabel="",color="black"];
b57->b23 [xlabel="F",color="orange"];
b57->b48 [xlabel="",color="black"];
b60->b61 [xlabel="",color="black"];
b61->final [xlabel="",color="black"];
b6->b18 [xlabel="",color="red"];
b6->b60 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b20[label="SwitchCase:exit\n"];
b23[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b31[label="SwitchCase:enter\n"];
b35[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b44[label="SwitchCase:enter\nNumLit(2)\n"];
b48[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b57[label="SwitchCase:enter\nNumLit(3)\n"];
b6[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b60[label="SwitchStmt:exit\n"];
b61[label="Prog:exit\n"];
b0->b31 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b20->b23 [xlabel="",color="red"];
b23->b35 [xlabel="",color="black"];
b31->b44 [xlabel="",color="black"];
b35->b48 [xlabel="",color="black"];
b44->b35 [xlabel="",color="black"];
b44->b57 [xlabel="F",color="orange"];
b48->b60 [xlabel="",color="black"];
b57->b23 [xlabel="F",color="orange"];
b57->b48 [xlabel="",color="black"];
b60->b61 [xlabel="",color="black"];
b61->final [xlabel="",color="black"];
b6->b20 [xlabel="",color="red"];
b6->b60 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\n"];
b17[label="SwitchCase:exit\n"];
b20[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b29[label="SwitchCase:enter\nNumLit(2)\n"];
b33[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b42[label="SwitchCase:enter\nNumLit(3)\n"];
b45[label="SwitchStmt:exit\n"];
b46[label="Prog:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b0->b29 [xlabel="",color="black"];
b17->b20 [xlabel="",color="red"];
b20->b33 [xlabel="",color="black"];
b29->b20 [xlabel="",color="black"];
b29->b42 [xlabel="F",color="orange"];
b33->b45 [xlabel="",color="black"];
b42->b33 [xlabel="",color="black"];
b42->b6 [xlabel="F",color="orange"];
b45->b46 [xlabel="",color="black"];
b46->final [xlabel="",color="black"];
b6->b17 [xlabel="",color="red"];
b6->b45 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
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
b0[label="Prog:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b19[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nBrkStmt:enter\nBrkStmt:exit\n"];
b29[label="SwitchCase:enter\n"];
b30[label="SwitchCase:exit\n"];
b33[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b42[label="SwitchCase:enter\nNumLit(2)\n"];
b46[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b55[label="SwitchCase:enter\nNumLit(3)\n"];
b58[label="SwitchStmt:exit\n"];
b59[label="Prog:exit\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b29 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b19->b30 [xlabel="",color="red"];
b19->b58 [xlabel="U",color="orange"];
b29->b42 [xlabel="",color="black"];
b30->b33 [xlabel="",color="red"];
b33->b46 [xlabel="",color="black"];
b42->b33 [xlabel="",color="black"];
b42->b55 [xlabel="F",color="orange"];
b46->b58 [xlabel="",color="black"];
b55->b19 [xlabel="F",color="orange"];
b55->b46 [xlabel="",color="black"];
b58->b59 [xlabel="",color="black"];
b59->final [xlabel="",color="black"];
b6->b19 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
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
b20[label="BlockStmt:exit\nSwitchCase:exit\n"];
b25[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b33[label="SwitchCase:enter\n"];
b37[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b46[label="SwitchCase:enter\nNumLit(2)\n"];
b49[label="SwitchStmt:exit\nBlockStmt:exit\n"];
b51[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b52[label="FnDec:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b20->b25 [xlabel="",color="red"];
b25->b37 [xlabel="",color="black"];
b33->b46 [xlabel="",color="black"];
b37->b49 [xlabel="",color="black"];
b46->b25 [xlabel="F",color="orange"];
b46->b37 [xlabel="",color="black"];
b49->b52 [xlabel="",color="black"];
b51->b33 [xlabel="F",color="orange"];
b51->b8 [xlabel="",color="black"];
b52->final [xlabel="",color="black"];
b8->b20 [xlabel="",color="red"];
b8->b52 [xlabel="U",color="orange"];
initial->b51 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryCatchBasic(t *testing.T) {
	ast, symtab, err := compile(`
try {
  doSth1();
  doSth2();
} catch (error) {
  log(error);
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b19[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b33[label="TryStmt:exit\nProg:exit\n"];
b0->b19 [xlabel="E",color="orange"];
b0->b33 [xlabel="",color="black"];
b19->b33 [xlabel="",color="black"];
b33->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TryCatchIf(t *testing.T) {
	ast, symtab, err := compile(`
try {
  if (a) {
    doSth1();
  }
  doSth2();
} catch (error) {
  log(error);
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nIfStmt:enter\nIdent(a)\n"];
b18[label="IfStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b26[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b40[label="TryStmt:exit\nProg:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b18 [xlabel="F",color="orange"];
b0->b26 [xlabel="E",color="orange"];
b0->b8 [xlabel="",color="black"];
b18->b26 [xlabel="E",color="orange"];
b18->b40 [xlabel="",color="black"];
b26->b40 [xlabel="",color="black"];
b40->final [xlabel="",color="black"];
b8->b18 [xlabel="",color="black"];
b8->b26 [xlabel="E",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TryCatchFin(t *testing.T) {
	ast, symtab, err := compile(`
try {
  doSth1();
  doSth2();
} catch (error) {
  log(error);
} finally {
  fin();
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b19[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b33[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nProg:exit\n"];
b0->b19 [xlabel="E",color="orange"];
b0->b33 [xlabel="",color="black"];
b19->b33 [xlabel="",color="black"];
b33->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TryReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    return;
    doSth2();
  } catch (error) {
    log(error);
  }
  afterFin()
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
b16[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b23[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b45[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b46[label="FnDec:exit\n"];
b16->b23 [xlabel="E",color="red"];
b16->b37 [xlabel="",color="red"];
b23->b37 [xlabel="",color="black"];
b37->b46 [xlabel="",color="black"];
b45->b16 [xlabel="",color="red"];
b45->b23 [xlabel="E",color="orange"];
b45->b46 [xlabel="U",color="orange"];
b46->final [xlabel="",color="black"];
initial->b45 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryReturnFin(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    return;
    doSth2();
  } catch (error) {
    log(error);
  } finally {
    fin()
  }
  afterFin()
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
b16[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b23[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b46[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b54[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b55[label="FnDec:exit\n"];
b16->b23 [xlabel="E",color="red"];
b16->b37 [xlabel="",color="red"];
b23->b37 [xlabel="",color="black"];
b37->b46 [xlabel="",color="black"];
b37->b55 [xlabel="P",color="orange"];
b46->b55 [xlabel="",color="black"];
b54->b16 [xlabel="",color="red"];
b54->b23 [xlabel="E",color="orange"];
b54->b37 [xlabel="U,P",color="orange"];
b55->final [xlabel="",color="black"];
initial->b54 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryReturnFinIf(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    return;
    doSth2();
  } catch (error) {
    log(error);
  } finally {
    if (a) {
      fin1()
    }
    fin2()
  }
  afterFin()
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
b16[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b23[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="BlockStmt:enter\nIfStmt:enter\nIdent(a)\n"];
b41[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b51[label="IfStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(fin2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b59[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b67[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b68[label="FnDec:exit\n"];
b16->b23 [xlabel="E",color="red"];
b16->b37 [xlabel="",color="red"];
b23->b37 [xlabel="",color="black"];
b37->b41 [xlabel="",color="black"];
b37->b51 [xlabel="F",color="orange"];
b41->b51 [xlabel="",color="black"];
b51->b59 [xlabel="",color="black"];
b51->b68 [xlabel="P",color="orange"];
b59->b68 [xlabel="",color="black"];
b67->b16 [xlabel="",color="red"];
b67->b23 [xlabel="E",color="orange"];
b67->b37 [xlabel="U,P",color="orange"];
b68->final [xlabel="",color="black"];
initial->b67 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryCatchReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    doSth2();
  } catch (error) {
    return
    log(error);
  }
  afterFin()
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
b21[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nRetStmt:enter\nRetStmt:exit\n"];
b27[label="ExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b45[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b46[label="FnDec:exit\n"];
b21->b27 [xlabel="",color="red"];
b21->b46 [xlabel="U",color="orange"];
b27->b37 [xlabel="",color="red"];
b37->b46 [xlabel="",color="black"];
b45->b21 [xlabel="E",color="orange"];
b45->b37 [xlabel="",color="black"];
b46->final [xlabel="",color="black"];
initial->b45 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryCatchReturnFin(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    doSth2();
  } catch (error) {
    return
    log(error);
  } finally {
    fin()
  }
  afterFin()
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
b21[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nRetStmt:enter\nRetStmt:exit\n"];
b27[label="ExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b46[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b54[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b55[label="FnDec:exit\n"];
b21->b27 [xlabel="",color="red"];
b21->b37 [xlabel="U,P",color="orange"];
b27->b37 [xlabel="",color="red"];
b37->b46 [xlabel="",color="black"];
b37->b55 [xlabel="P",color="orange"];
b46->b55 [xlabel="",color="black"];
b54->b21 [xlabel="E",color="orange"];
b54->b37 [xlabel="",color="black"];
b55->final [xlabel="",color="black"];
initial->b54 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryFinReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1();
    doSth2();
  } catch (error) {
    return
    log(error);
  } finally {
    fin()
    return
  }
  afterFin()
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
b21[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nRetStmt:enter\nRetStmt:exit\n"];
b27[label="ExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b47[label="BlockStmt:exit\n"];
b48[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(afterFin)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b56[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b57[label="FnDec:exit\n"];
b21->b27 [xlabel="",color="red"];
b21->b37 [xlabel="U,P",color="orange"];
b27->b37 [xlabel="",color="red"];
b37->b47 [xlabel="",color="red"];
b37->b57 [xlabel="U",color="orange"];
b47->b48 [xlabel="",color="red"];
b47->b57 [xlabel="P",color="red"];
b48->b57 [xlabel="",color="red"];
b56->b21 [xlabel="E",color="orange"];
b56->b37 [xlabel="",color="black"];
b57->final [xlabel="",color="black"];
initial->b56 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryNested(t *testing.T) {
	ast, symtab, err := compile(`
try {
  doSth1()

  try {
    doSth2()
  } catch (error) {

  }
} catch (error) {
  log(error)
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\n"];
b12[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b22[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b29[label="TryStmt:exit\nBlockStmt:exit\n"];
b31[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b45[label="TryStmt:exit\nProg:exit\n"];
b0->b12 [xlabel="",color="black"];
b0->b31 [xlabel="E",color="orange"];
b12->b22 [xlabel="E",color="orange"];
b12->b29 [xlabel="",color="black"];
b22->b29 [xlabel="",color="black"];
b22->b31 [xlabel="E",color="orange"];
b29->b45 [xlabel="",color="black"];
b31->b45 [xlabel="",color="black"];
b45->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TryNestedTryReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1()
  
    try {
      doSth2()
      return
    } catch (error) {
  
    }
  } catch (error) {
    log(error)
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
b14[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b25[label="BlockStmt:exit\n"];
b26[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b33[label="TryStmt:exit\nBlockStmt:exit\n"];
b35[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b49[label="TryStmt:exit\nBlockStmt:exit\n"];
b51[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\n"];
b52[label="FnDec:exit\n"];
b14->b25 [xlabel="",color="red"];
b14->b26 [xlabel="E",color="orange"];
b14->b52 [xlabel="U",color="orange"];
b25->b33 [xlabel="",color="red"];
b26->b33 [xlabel="",color="black"];
b26->b35 [xlabel="E",color="orange"];
b33->b49 [xlabel="",color="black"];
b35->b49 [xlabel="",color="black"];
b49->b52 [xlabel="",color="black"];
b51->b14 [xlabel="",color="black"];
b51->b35 [xlabel="E",color="orange"];
b52->final [xlabel="",color="black"];
initial->b51 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryNestedTryFinReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1()

    try {
      doSth2()
      return
    } catch (error) {

    } finally {
      fin1()
    }
  } catch (error) {
    log(error)
  } finally {
    fin2()
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
b14[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b25[label="BlockStmt:exit\n"];
b26[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b33[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b42[label="TryStmt:exit\nBlockStmt:exit\n"];
b44[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b58[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\n"];
b69[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\n"];
b70[label="FnDec:exit\n"];
b14->b25 [xlabel="",color="red"];
b14->b26 [xlabel="E",color="orange"];
b14->b33 [xlabel="U,P",color="orange"];
b25->b33 [xlabel="",color="red"];
b26->b33 [xlabel="",color="black"];
b26->b44 [xlabel="E",color="orange"];
b33->b42 [xlabel="",color="black"];
b33->b44 [xlabel="E",color="orange"];
b33->b58 [xlabel="P",color="orange"];
b42->b58 [xlabel="",color="black"];
b44->b58 [xlabel="",color="black"];
b58->b70 [xlabel="",color="black"];
b69->b14 [xlabel="",color="black"];
b69->b44 [xlabel="E",color="orange"];
b70->final [xlabel="",color="black"];
initial->b69 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryNestedNested(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1()

    try {
      doSth2()

      try {
        doSth3()
      } catch (error) {
        log1()
      } finally {
        fin1()
      }
    } catch (error) {
      log2()
    } finally {
      fin2()
    }
  } catch (error) {
    log3(error)
  } finally {
    fin3()
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
b106[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\n"];
b14[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\n"];
b23[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b33[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b46[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\n"];
b57[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b70[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\n"];
b81[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log3)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b95[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin3)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b106->b14 [xlabel="",color="black"];
b106->b81 [xlabel="E",color="orange"];
b14->b23 [xlabel="",color="black"];
b14->b57 [xlabel="E",color="orange"];
b23->b33 [xlabel="E",color="orange"];
b23->b46 [xlabel="",color="black"];
b33->b46 [xlabel="",color="black"];
b33->b57 [xlabel="E",color="orange"];
b46->b57 [xlabel="E",color="orange"];
b46->b70 [xlabel="",color="black"];
b57->b70 [xlabel="",color="black"];
b57->b81 [xlabel="E",color="orange"];
b70->b81 [xlabel="E",color="orange"];
b70->b95 [xlabel="",color="black"];
b81->b95 [xlabel="",color="black"];
b95->final [xlabel="",color="black"];
initial->b106 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_TryNestedNestedTryFinReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    doSth1()

    try {
      doSth2()

      try {
        doSth3()
        return
      } catch (error) {
  
      } finally {
        fin1()
      }
    } catch (error) {

    } finally {
      fin2()
    }
  } catch (error) {
    log(error)
  } finally {
    fin3()
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
b14[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\n"];
b23[label="TryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b34[label="BlockStmt:exit\n"];
b35[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b42[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin1)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b51[label="TryStmt:exit\nBlockStmt:exit\n"];
b53[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b60[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin2)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\n"];
b71[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b85[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(fin3)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nTryStmt:exit\nBlockStmt:exit\n"];
b96[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\n"];
b97[label="FnDec:exit\n"];
b14->b23 [xlabel="",color="black"];
b14->b53 [xlabel="E",color="orange"];
b23->b34 [xlabel="",color="red"];
b23->b35 [xlabel="E",color="orange"];
b23->b42 [xlabel="U,P",color="orange"];
b34->b42 [xlabel="",color="red"];
b35->b42 [xlabel="",color="black"];
b35->b53 [xlabel="E",color="orange"];
b42->b51 [xlabel="",color="black"];
b42->b53 [xlabel="E",color="orange"];
b42->b60 [xlabel="P",color="orange"];
b51->b60 [xlabel="",color="black"];
b53->b60 [xlabel="",color="black"];
b53->b71 [xlabel="E",color="orange"];
b60->b71 [xlabel="E",color="orange"];
b60->b85 [xlabel="",color="black"];
b71->b85 [xlabel="",color="black"];
b85->b97 [xlabel="",color="black"];
b96->b14 [xlabel="",color="black"];
b96->b71 [xlabel="E",color="orange"];
b97->final [xlabel="",color="black"];
initial->b96 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_Throw(t *testing.T) {
	ast, symtab, err := compile(`
try {
  throw a
  a
} catch (error) {
  log(error)
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nThrowStmt:enter\nIdent(a)\n"];
b13[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b27[label="TryStmt:exit\nProg:exit\n"];
b8[label="ThrowStmt:exit\n"];
b9[label="ExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b13 [xlabel="E",color="orange"];
b0->b8 [xlabel="",color="black"];
b13->b27 [xlabel="",color="black"];
b27->final [xlabel="",color="black"];
b8->b13 [xlabel="U",color="orange"];
b8->b9 [xlabel="",color="red"];
b9->b13 [xlabel="E",color="red"];
b9->b27 [xlabel="",color="red"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ThrowBin(t *testing.T) {
	ast, symtab, err := compile(`
try {
  throw a && b || c
  a
} catch (error) {
  log(error)
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\nThrowStmt:enter\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(a)\n"];
b10[label="Ident(b)\nBinExpr(&&):exit\n"];
b14[label="Ident(c)\nBinExpr(||):exit\n"];
b18[label="ThrowStmt:exit\n"];
b19[label="ExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\n"];
b23[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b37[label="TryStmt:exit\nProg:exit\n"];
b0->b10 [xlabel="",color="black"];
b0->b14 [xlabel="F",color="orange"];
b0->b23 [xlabel="E",color="orange"];
b10->b14 [xlabel="",color="black"];
b10->b18 [xlabel="T",color="orange"];
b10->b23 [xlabel="E",color="orange"];
b14->b18 [xlabel="",color="black"];
b14->b23 [xlabel="E",color="orange"];
b18->b19 [xlabel="",color="red"];
b18->b23 [xlabel="U",color="orange"];
b19->b23 [xlabel="E",color="red"];
b19->b37 [xlabel="",color="red"];
b23->b37 [xlabel="",color="black"];
b37->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ThrowNested(t *testing.T) {
	ast, symtab, err := compile(`
try {
  try {
    throw a
  } catch (error) {
    log1(error)
  }
} catch (error) {
  log2(error)
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\n"];
b11[label="ThrowStmt:exit\n"];
b12[label="BlockStmt:exit\n"];
b13[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log1)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b27[label="TryStmt:exit\nBlockStmt:exit\n"];
b29[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log2)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b43[label="TryStmt:exit\nProg:exit\n"];
b5[label="TryStmt:enter\nBlockStmt:enter\nThrowStmt:enter\nIdent(a)\n"];
b0->b5 [xlabel="",color="black"];
b11->b12 [xlabel="",color="red"];
b11->b13 [xlabel="U",color="orange"];
b12->b27 [xlabel="",color="red"];
b13->b27 [xlabel="",color="black"];
b13->b29 [xlabel="E",color="orange"];
b27->b43 [xlabel="",color="black"];
b29->b43 [xlabel="",color="black"];
b43->final [xlabel="",color="black"];
b5->b11 [xlabel="",color="black"];
b5->b13 [xlabel="E",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ThrowNestedReturn(t *testing.T) {
	ast, symtab, err := compile(`
function f() {
  try {
    if (a) {
      return a;
    } else {
      throw b;
    }
  } catch (err) {
    // do nothing.
  }

  foo();
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
b10[label="BlockStmt:enter\nRetStmt:enter\nIdent(a)\nRetStmt:exit\n"];
b15[label="BlockStmt:exit\n"];
b16[label="BlockStmt:enter\nThrowStmt:enter\nIdent(b)\n"];
b20[label="ThrowStmt:exit\n"];
b21[label="BlockStmt:exit\n"];
b23[label="IfStmt:exit\nBlockStmt:exit\n"];
b25[label="Catch:enter\nIdent(err)\nBlockStmt:enter\nBlockStmt:exit\nCatch:exit\n"];
b32[label="TryStmt:exit\nExprStmt:enter\nCallExpr:enter\nIdent(foo)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\n"];
b40[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nTryStmt:enter\nBlockStmt:enter\nIfStmt:enter\nIdent(a)\n"];
b41[label="FnDec:exit\n"];
b10->b15 [xlabel="",color="red"];
b10->b25 [xlabel="E",color="orange"];
b10->b41 [xlabel="U",color="orange"];
b15->b23 [xlabel="",color="red"];
b16->b20 [xlabel="",color="black"];
b16->b25 [xlabel="E",color="orange"];
b20->b21 [xlabel="",color="red"];
b20->b25 [xlabel="U",color="orange"];
b21->b23 [xlabel="",color="red"];
b23->b32 [xlabel="",color="red"];
b25->b32 [xlabel="",color="black"];
b32->b41 [xlabel="",color="black"];
b40->b10 [xlabel="",color="black"];
b40->b16 [xlabel="F",color="orange"];
b40->b25 [xlabel="E",color="orange"];
b41->final [xlabel="",color="black"];
initial->b40 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ThrowCatch(t *testing.T) {
	ast, symtab, err := compile(`
try {
  try {
    throw 1
    a
  } catch (error) {
    throw error
    b
  }
  c
} catch (error) {
  log2(error)
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
b0[label="Prog:enter\nTryStmt:enter\nBlockStmt:enter\n"];
b12[label="ExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\n"];
b16[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nThrowStmt:enter\nIdent(error)\n"];
b22[label="ThrowStmt:exit\n"];
b23[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b29[label="TryStmt:exit\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b34[label="Catch:enter\nIdent(error)\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(log2)\nIdent(error)\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nCatch:exit\n"];
b48[label="TryStmt:exit\nProg:exit\n"];
b5[label="TryStmt:enter\nBlockStmt:enter\nThrowStmt:enter\nNumLit(1)\nThrowStmt:exit\n"];
b0->b5 [xlabel="",color="black"];
b12->b16 [xlabel="E",color="red"];
b12->b29 [xlabel="",color="red"];
b16->b22 [xlabel="",color="black"];
b16->b34 [xlabel="E",color="orange"];
b22->b23 [xlabel="",color="red"];
b22->b34 [xlabel="U",color="orange"];
b23->b29 [xlabel="",color="red"];
b23->b34 [xlabel="E",color="red"];
b29->b34 [xlabel="E",color="red"];
b29->b48 [xlabel="",color="red"];
b34->b48 [xlabel="",color="black"];
b48->final [xlabel="",color="black"];
b5->b12 [xlabel="",color="red"];
b5->b16 [xlabel="U",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ThrowBare(t *testing.T) {
	ast, symtab, err := compile(`
throw 1
a
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
b0[label="Prog:enter\nThrowStmt:enter\nNumLit(1)\nThrowStmt:exit\n"];
b6[label="ExprStmt:enter\nIdent(a)\nExprStmt:exit\n"];
b9[label="Prog:exit\n"];
b0->b6 [xlabel="",color="red"];
b0->b9 [xlabel="U",color="orange"];
b6->b9 [xlabel="",color="red"];
b9->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_NewExpr(t *testing.T) {
	ast, symtab, err := compile(`
  new fn()
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
b0[label="Prog:enter\nExprStmt:enter\nNewExpr:enter\nIdent(fn)\nNewExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_NewExprArgs(t *testing.T) {
	ast, symtab, err := compile(`
  new fn(1, 2, a && b)
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
b0[label="Prog:enter\nExprStmt:enter\nNewExpr:enter\nIdent(fn)\nNumLit(1)\nNumLit(2)\nBinExpr(&&):enter\nIdent(a)\n"];
b10[label="Ident(b)\nBinExpr(&&):exit\n"];
b14[label="NewExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b10 [xlabel="",color="black"];
b0->b14 [xlabel="F",color="orange"];
b10->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_NewExprArgsBin(t *testing.T) {
	ast, symtab, err := compile(`
  new fn(1, a || b && c, d && e)
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
b0[label="Prog:enter\nExprStmt:enter\nNewExpr:enter\nIdent(fn)\nNumLit(1)\nBinExpr(||):enter\nIdent(a)\n"];
b11[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b18[label="BinExpr(&&):enter\nIdent(d)\n"];
b20[label="Ident(e)\nBinExpr(&&):exit\n"];
b24[label="NewExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="BinExpr(&&):enter\nIdent(b)\n"];
b0->b18 [xlabel="T",color="orange"];
b0->b9 [xlabel="",color="black"];
b11->b18 [xlabel="",color="black"];
b18->b20 [xlabel="",color="black"];
b18->b24 [xlabel="F",color="orange"];
b20->b24 [xlabel="",color="black"];
b24->final [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
b9->b18 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a.b
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nIdent(a)\nIdent(b)\nMemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprChain(t *testing.T) {
	ast, symtab, err := compile(`
  a.b.c
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(a)\nIdent(b)\nMemberExpr:exit\nIdent(c)\nMemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprCompute(t *testing.T) {
	ast, symtab, err := compile(`
  (a && b).c
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nParenExpr:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b12[label="ParenExpr:exit\nIdent(c)\nMemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b8[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprSubscript(t *testing.T) {
	ast, symtab, err := compile(`
  a[b]
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nIdent(a)\nIdent(b)\nMemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprSubscriptBin(t *testing.T) {
	ast, symtab, err := compile(`
  a[b && c]
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nIdent(a)\nBinExpr(&&):enter\nIdent(b)\n"];
b12[label="MemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b8[label="Ident(c)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprSubscriptBinBin(t *testing.T) {
	ast, symtab, err := compile(`
  a[b && c][e || f]
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(a)\nBinExpr(&&):enter\nIdent(b)\n"];
b13[label="MemberExpr:exit\nBinExpr(||):enter\nIdent(e)\n"];
b17[label="Ident(f)\nBinExpr(||):exit\n"];
b21[label="MemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="Ident(c)\nBinExpr(&&):exit\n"];
b0->b13 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b13->b17 [xlabel="",color="black"];
b13->b21 [xlabel="T",color="orange"];
b17->b21 [xlabel="",color="black"];
b21->final [xlabel="",color="black"];
b9->b13 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MemberExprMixChain(t *testing.T) {
	ast, symtab, err := compile(`
  a.b[c && d].e
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
b0[label="Prog:enter\nExprStmt:enter\nMemberExpr:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(a)\nIdent(b)\nMemberExpr:exit\nBinExpr(&&):enter\nIdent(c)\n"];
b13[label="Ident(d)\nBinExpr(&&):exit\n"];
b17[label="MemberExpr:exit\nIdent(e)\nMemberExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b13 [xlabel="",color="black"];
b0->b17 [xlabel="F",color="orange"];
b13->b17 [xlabel="",color="black"];
b17->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SeqExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a, b, c
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
b0[label="Prog:enter\nExprStmt:enter\nSeqExpr:enter\nIdent(a)\nIdent(b)\nIdent(c)\nSeqExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SeqExprBin(t *testing.T) {
	ast, symtab, err := compile(`
  a, b && c || d, e || f
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
b0[label="Prog:enter\nExprStmt:enter\nSeqExpr:enter\nIdent(a)\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(b)\n"];
b13[label="Ident(d)\nBinExpr(||):exit\n"];
b17[label="BinExpr(||):enter\nIdent(e)\n"];
b19[label="Ident(f)\nBinExpr(||):exit\n"];
b23[label="SeqExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="Ident(c)\nBinExpr(&&):exit\n"];
b0->b13 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b13->b17 [xlabel="",color="black"];
b17->b19 [xlabel="",color="black"];
b17->b23 [xlabel="T",color="orange"];
b19->b23 [xlabel="",color="black"];
b23->final [xlabel="",color="black"];
b9->b13 [xlabel="",color="black"];
b9->b17 [xlabel="T",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_SeqExprParen(t *testing.T) {
	ast, symtab, err := compile(`
  (a, b && c || d, e || f)
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
b0[label="Prog:enter\nExprStmt:enter\nParenExpr:enter\nSeqExpr:enter\nIdent(a)\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(b)\n"];
b10[label="Ident(c)\nBinExpr(&&):exit\n"];
b14[label="Ident(d)\nBinExpr(||):exit\n"];
b18[label="BinExpr(||):enter\nIdent(e)\n"];
b20[label="Ident(f)\nBinExpr(||):exit\n"];
b24[label="SeqExpr:exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b10 [xlabel="",color="black"];
b0->b14 [xlabel="F",color="orange"];
b10->b14 [xlabel="",color="black"];
b10->b18 [xlabel="T",color="orange"];
b14->b18 [xlabel="",color="black"];
b18->b20 [xlabel="",color="black"];
b18->b24 [xlabel="T",color="orange"];
b20->b24 [xlabel="",color="black"];
b24->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ThisExpr(t *testing.T) {
	ast, symtab, err := compile(`
  function f() {
    this
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
b9[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nThisExpr\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b9->final [xlabel="",color="black"];
initial->b9 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_Unary(t *testing.T) {
	ast, symtab, err := compile(`
  !a
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
b0[label="Prog:enter\nExprStmt:enter\nUnaryExpr(!):enter\nIdent(a)\nUnaryExpr(!):exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_UnaryParen(t *testing.T) {
	ast, symtab, err := compile(`
  !(a || b && c)
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
b0[label="Prog:enter\nExprStmt:enter\nUnaryExpr(!):enter\nParenExpr:enter\nBinExpr(||):enter\nIdent(a)\n"];
b10[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b17[label="ParenExpr:exit\nUnaryExpr(!):exit\nExprStmt:exit\nProg:exit\n"];
b8[label="BinExpr(&&):enter\nIdent(b)\n"];
b0->b17 [xlabel="T",color="orange"];
b0->b8 [xlabel="",color="black"];
b10->b17 [xlabel="",color="black"];
b17->final [xlabel="",color="black"];
b8->b10 [xlabel="",color="black"];
b8->b17 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a ? b : c
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nIdent(a)\n"];
b6[label="Ident(b)\n"];
b7[label="Ident(c)\n"];
b8[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b6 [xlabel="",color="black"];
b0->b7 [xlabel="F",color="orange"];
b6->b8 [xlabel="",color="black"];
b7->b8 [xlabel="",color="black"];
b8->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprTestBin(t *testing.T) {
	ast, symtab, err := compile(`
  a && b ? c : d
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b11[label="Ident(c)\n"];
b12[label="Ident(d)\n"];
b13[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b7[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b7 [xlabel="",color="black"];
b11->b13 [xlabel="",color="black"];
b12->b13 [xlabel="",color="black"];
b13->final [xlabel="",color="black"];
b7->b11 [xlabel="",color="black"];
b7->b12 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprTestBinOr(t *testing.T) {
	ast, symtab, err := compile(`
  a || b ? c : d
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nBinExpr(||):enter\nIdent(a)\n"];
b11[label="Ident(c)\n"];
b12[label="Ident(d)\n"];
b13[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b7[label="Ident(b)\nBinExpr(||):exit\n"];
b0->b11 [xlabel="T",color="orange"];
b0->b7 [xlabel="",color="black"];
b11->b13 [xlabel="",color="black"];
b12->b13 [xlabel="",color="black"];
b13->final [xlabel="",color="black"];
b7->b11 [xlabel="",color="black"];
b7->b12 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprTestBinMix(t *testing.T) {
	ast, symtab, err := compile(`
  a || b && c ? d : e
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nBinExpr(||):enter\nIdent(a)\n"];
b16[label="Ident(d)\n"];
b17[label="Ident(e)\n"];
b18[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b7[label="BinExpr(&&):enter\nIdent(b)\n"];
b9[label="Ident(c)\nBinExpr(&&):exit\nBinExpr(||):exit\n"];
b0->b16 [xlabel="T",color="orange"];
b0->b7 [xlabel="",color="black"];
b16->b18 [xlabel="",color="black"];
b17->b18 [xlabel="",color="black"];
b18->final [xlabel="",color="black"];
b7->b17 [xlabel="F",color="orange"];
b7->b9 [xlabel="",color="black"];
b9->b16 [xlabel="",color="black"];
b9->b17 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprTestBinMixAnd(t *testing.T) {
	ast, symtab, err := compile(`
  a && b || c ? d : e
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(a)\n"];
b12[label="Ident(c)\nBinExpr(||):exit\n"];
b16[label="Ident(d)\n"];
b17[label="Ident(e)\n"];
b18[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b8[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->b16 [xlabel="",color="black"];
b12->b17 [xlabel="F",color="orange"];
b16->b18 [xlabel="",color="black"];
b17->b18 [xlabel="",color="black"];
b18->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
b8->b16 [xlabel="T",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprNestedLeft(t *testing.T) {
	ast, symtab, err := compile(`
  a ? b ? c : d : c
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nIdent(a)\n"];
b10[label="CondExpr:exit\n"];
b12[label="Ident(c)\n"];
b13[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b6[label="CondExpr:enter\nIdent(b)\n"];
b8[label="Ident(c)\n"];
b9[label="Ident(d)\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b10->b13 [xlabel="",color="black"];
b12->b13 [xlabel="",color="black"];
b13->final [xlabel="",color="black"];
b6->b8 [xlabel="",color="black"];
b6->b9 [xlabel="F",color="orange"];
b8->b10 [xlabel="",color="black"];
b9->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_CondExprNestedRight(t *testing.T) {
	ast, symtab, err := compile(`
  a ? b : c ? d : e
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
b0[label="Prog:enter\nExprStmt:enter\nCondExpr:enter\nIdent(a)\n"];
b10[label="Ident(e)\n"];
b11[label="CondExpr:exit\n"];
b13[label="CondExpr:exit\nExprStmt:exit\nProg:exit\n"];
b6[label="Ident(b)\n"];
b7[label="CondExpr:enter\nIdent(c)\n"];
b9[label="Ident(d)\n"];
b0->b6 [xlabel="",color="black"];
b0->b7 [xlabel="F",color="orange"];
b10->b11 [xlabel="",color="black"];
b11->b13 [xlabel="",color="black"];
b13->final [xlabel="",color="black"];
b6->b13 [xlabel="",color="black"];
b7->b10 [xlabel="F",color="orange"];
b7->b9 [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_AssignExpr(t *testing.T) {
	ast, symtab, err := compile(`
  a = 1
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nNumLit(1)\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_AssignExprBin(t *testing.T) {
	ast, symtab, err := compile(`
  a = b && c
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nIdent(a)\nBinExpr(&&):enter\nIdent(b)\n"];
b12[label="AssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b8[label="Ident(c)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportCall(t *testing.T) {
	ast, symtab, err := compile(`
  import("a")
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
b0[label="Prog:enter\nExprStmt:enter\nImportCall:enter\nStrLit\nImportCall:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_AwaitImportExpr(t *testing.T) {
	ast, symtab, err := compile(`
  await import("a")
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
b0[label="Prog:enter\nExprStmt:enter\nUnaryExpr(await):enter\nImportCall:enter\nStrLit\nImportCall:exit\nUnaryExpr(await):exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_WithStmt(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)

	ast, symtab, err := compile(`
  with (a) {
    b
  }
  `, opts)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nWithStmt:enter\nIdent(a)\nBlockStmt:enter\nExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\nWithStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_WithStmtBin(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)

	ast, symtab, err := compile(`
  with (a || b) {
    c
  }
  `, opts)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nWithStmt:enter\nBinExpr(||):enter\nIdent(a)\n"];
b10[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nWithStmt:exit\nProg:exit\n"];
b6[label="Ident(b)\nBinExpr(||):exit\n"];
b0->b10 [xlabel="T",color="orange"];
b0->b6 [xlabel="",color="black"];
b10->final [xlabel="",color="black"];
b6->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportEffect(t *testing.T) {
	ast, symtab, err := compile(`
  import "a";
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
b0[label="Prog:enter\nImportDec:enter\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportNs(t *testing.T) {
	ast, symtab, err := compile(`
  import * as a from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec(Namespace):enter\nIdent(a)\nImportSpec(Namespace):exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportDefault(t *testing.T) {
	ast, symtab, err := compile(`
  import a from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec(Default):enter\nIdent(a)\nImportSpec(Default):exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportDefaultMix(t *testing.T) {
	ast, symtab, err := compile(`
  import a, * as b from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec(Default):enter\nIdent(a)\nImportSpec(Default):exit\nImportSpec(Namespace):enter\nIdent(b)\nImportSpec(Namespace):exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportNamed(t *testing.T) {
	ast, symtab, err := compile(`
  import { a } from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec:enter\nIdent(a)\nIdent(a)\nImportSpec:exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportDefaultMixNamed(t *testing.T) {
	ast, symtab, err := compile(`
  import a, { b } from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec(Default):enter\nIdent(a)\nImportSpec(Default):exit\nImportSpec:enter\nIdent(b)\nIdent(b)\nImportSpec:exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ImportNamedAs(t *testing.T) {
	ast, symtab, err := compile(`
  import { a as b } from "a";
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
b0[label="Prog:enter\nImportDec:enter\nImportSpec:enter\nIdent(a)\nIdent(b)\nImportSpec:exit\nStrLit\nImportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportIndividual(t *testing.T) {
	ast, symtab, err := compile(`
  export let a, b, c;
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
b0[label="Prog:enter\nExportDec:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDec:enter\nIdent(c)\nVarDec:exit\nVarDecStmt:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportFn(t *testing.T) {
	ast, symtab, err := compile(`
  export function f() {}
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
b0[label="Prog:enter\nExportDec:enter\nFnDec\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportClass(t *testing.T) {
	ast, symtab, err := compile(`
  export class A {}
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
b0[label="Prog:enter\nExportDec:enter\nClassDec:enter\nIdent(A)\nClassBody:enter\nClassBody:exit\nClassDec:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportList(t *testing.T) {
	ast, symtab, err := compile(`
  let a, b, c;
  export { a, b, c }
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDec:enter\nIdent(c)\nVarDec:exit\nVarDecStmt:exit\nExportDec:enter\nExportSpec:enter\nIdent(a)\nIdent(a)\nExportSpec:exit\nExportSpec:enter\nIdent(b)\nIdent(b)\nExportSpec:exit\nExportSpec:enter\nIdent(c)\nIdent(c)\nExportSpec:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportRename(t *testing.T) {
	ast, symtab, err := compile(`
  let a, b, c;
  export { a as aa, b as bb, c as cc }
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDec:enter\nIdent(c)\nVarDec:exit\nVarDecStmt:exit\nExportDec:enter\nExportSpec:enter\nIdent(aa)\nIdent(a)\nExportSpec:exit\nExportSpec:enter\nIdent(bb)\nIdent(b)\nExportSpec:exit\nExportSpec:enter\nIdent(cc)\nIdent(c)\nExportSpec:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportDestructArr(t *testing.T) {
	ast, symtab, err := compile(`
  export const [ a ] = arr
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
b0[label="Prog:enter\nExportDec:enter\nVarDecStmt:enter\nVarDec:enter\nArrPat:enter\nIdent(a)\nArrPat:exit\nIdent(arr)\nVarDec:exit\nVarDecStmt:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportDestructObj(t *testing.T) {
	ast, symtab, err := compile(`
  export const { a, b: c } = obj
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
b0[label="Prog:enter\nExportDec:enter\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(a)\nIdent(a)\nProp:exit\nProp:enter\nIdent(b)\nIdent(c)\nProp:exit\nObjPat:exit\nIdent(obj)\nVarDec:exit\nVarDecStmt:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportDefault(t *testing.T) {
	ast, symtab, err := compile(`
  export default a
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
b0[label="Prog:enter\nExportDec(Default):enter\nIdent(a)\nExportDec(Default):exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportDefaultFn(t *testing.T) {
	ast, symtab, err := compile(`
  export default function f() {}
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
b0[label="Prog:enter\nExportDec(Default):enter\nFnDec\nExportDec(Default):exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportAsDefault(t *testing.T) {
	ast, symtab, err := compile(`
  let a;
  export { a as default }
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDecStmt:exit\nExportDec:enter\nExportSpec:enter\nIdent(default)\nIdent(a)\nExportSpec:exit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportNs(t *testing.T) {
	ast, symtab, err := compile(`
  export * from "a"
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
b0[label="Prog:enter\nExportDec(All):enter\nStrLit\nExportDec(All):exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportNsAs(t *testing.T) {
	ast, symtab, err := compile(`
  export * as a from "a"
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
b0[label="Prog:enter\nExportDec(All):enter\nExportSpec(Namespace):enter\nIdent(a)\nExportSpec(Namespace):exit\nStrLit\nExportDec(All):exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ExportNamed(t *testing.T) {
	ast, symtab, err := compile(`
  export { a, b, c as cc, default, default as d } from "a"
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
b0[label="Prog:enter\nExportDec:enter\nExportSpec:enter\nIdent(a)\nIdent(a)\nExportSpec:exit\nExportSpec:enter\nIdent(b)\nIdent(b)\nExportSpec:exit\nExportSpec:enter\nIdent(cc)\nIdent(c)\nExportSpec:exit\nExportSpec:enter\nIdent(default)\nIdent(default)\nExportSpec:exit\nExportSpec:enter\nIdent(d)\nIdent(default)\nExportSpec:exit\nStrLit\nExportDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArr(t *testing.T) {
	ast, symtab, err := compile(`
  let [a, b] = [1, 2]
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nArrPat:enter\nIdent(a)\nIdent(b)\nArrPat:exit\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrAssign(t *testing.T) {
	ast, symtab, err := compile(`
  let a, b;
  [a, b] = [1, 2]
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDecStmt:exit\nExprStmt:enter\nAssignExpr:enter\nArrPat:enter\nIdent(a)\nIdent(b)\nArrPat:exit\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrDefault(t *testing.T) {
	ast, symtab, err := compile(`
  let [a = 3, b = 4] = [1, 2]
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nArrPat:enter\nAssignPat:enter\nIdent(a)\n"];
b13[label="NumLit(4)\n"];
b14[label="AssignPat:exit\nArrPat:exit\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b8[label="NumLit(3)\n"];
b9[label="AssignPat:exit\nAssignPat:enter\nIdent(b)\n"];
b0->b8 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b13->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
b8->b9 [xlabel="",color="black"];
b9->b13 [xlabel="F",color="orange"];
b9->b14 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrNested(t *testing.T) {
	ast, symtab, err := compile(`
  let [a = e && f, [ b, c = 2 ]] = arr
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nArrPat:enter\nAssignPat:enter\nIdent(a)\n"];
b10[label="Ident(f)\nBinExpr(&&):exit\n"];
b14[label="AssignPat:exit\nArrPat:enter\nIdent(b)\nAssignPat:enter\nIdent(c)\n"];
b20[label="NumLit(2)\n"];
b21[label="AssignPat:exit\nArrPat:exit\nArrPat:exit\nIdent(arr)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b8[label="BinExpr(&&):enter\nIdent(e)\n"];
b0->b14 [xlabel="",color="black"];
b0->b8 [xlabel="F",color="orange"];
b10->b14 [xlabel="",color="black"];
b14->b20 [xlabel="F",color="orange"];
b14->b21 [xlabel="",color="black"];
b20->b21 [xlabel="",color="black"];
b21->final [xlabel="",color="black"];
b8->b10 [xlabel="",color="black"];
b8->b14 [xlabel="F",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrRest(t *testing.T) {
	ast, symtab, err := compile(`
  let [a, b, ...rest] = arr
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nArrPat:enter\nIdent(a)\nIdent(b)\nRestPat:enter\nIdent(rest)\nRestPat:exit\nArrPat:exit\nIdent(arr)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrDiscardMiddle(t *testing.T) {
	ast, symtab, err := compile(`
  [a,,b] = [1, 2]
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nArrPat:enter\nIdent(a)\nIdent(b)\nArrPat:exit\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatArrDiscardAll(t *testing.T) {
	ast, symtab, err := compile(`
  [,,] = [1, 2]
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
b0[label="Prog:enter\nExprStmt:enter\nAssignExpr:enter\nArrPat:enter\nArrPat:exit\nArrLit:enter\nNumLit(1)\nNumLit(2)\nArrLit:exit\nAssignExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObj(t *testing.T) {
	ast, symtab, err := compile(`
  let {a, b} = { a: 1, b: 2 }
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(a)\nIdent(a)\nProp:exit\nProp:enter\nIdent(b)\nIdent(b)\nProp:exit\nObjPat:exit\nObjLit:enter\nProp:enter\nIdent(a)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(b)\nNumLit(2)\nProp:exit\nObjLit:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjAssign(t *testing.T) {
	ast, symtab, err := compile(`
  let a, b;
  ({ a, b } = { a: 1, b: 2 })
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDecStmt:exit\nExprStmt:enter\nParenExpr:enter\nAssignExpr:enter\nObjPat:enter\nProp:enter\nIdent(a)\nIdent(a)\nProp:exit\nProp:enter\nIdent(b)\nIdent(b)\nProp:exit\nObjPat:exit\nObjLit:enter\nProp:enter\nIdent(a)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(b)\nNumLit(2)\nProp:exit\nObjLit:exit\nAssignExpr:exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjDefault(t *testing.T) {
	ast, symtab, err := compile(`
  let a, b;
  ({ a = 3, b = 4 } = { a: 1, b: 2 })
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nVarDec:exit\nVarDec:enter\nIdent(b)\nVarDec:exit\nVarDecStmt:exit\nExprStmt:enter\nParenExpr:enter\nAssignExpr:enter\nObjPat:enter\nProp:enter\nIdent(a)\nAssignPat:enter\nIdent(a)\n"];
b21[label="NumLit(3)\n"];
b22[label="AssignPat:exit\nProp:exit\nProp:enter\nIdent(b)\nAssignPat:enter\nIdent(b)\n"];
b30[label="NumLit(4)\n"];
b31[label="AssignPat:exit\nProp:exit\nObjPat:exit\nObjLit:enter\nProp:enter\nIdent(a)\nNumLit(1)\nProp:exit\nProp:enter\nIdent(b)\nNumLit(2)\nProp:exit\nObjLit:exit\nAssignExpr:exit\nParenExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b21 [xlabel="F",color="orange"];
b0->b22 [xlabel="",color="black"];
b21->b22 [xlabel="",color="black"];
b22->b30 [xlabel="F",color="orange"];
b22->b31 [xlabel="",color="black"];
b30->b31 [xlabel="",color="black"];
b31->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjRename(t *testing.T) {
	ast, symtab, err := compile(`
  let { a: b } = obj
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(a)\nIdent(b)\nProp:exit\nObjPat:exit\nIdent(obj)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjNested(t *testing.T) {
	ast, symtab, err := compile(`
  let { a: { b } } = obj
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(a)\nObjPat:enter\nProp:enter\nIdent(b)\nIdent(b)\nProp:exit\nObjPat:exit\nProp:exit\nObjPat:exit\nIdent(obj)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjCompute(t *testing.T) {
	ast, symtab, err := compile(`
  let key = "a"
  let { [key]: foo } = obj
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(key)\nStrLit\nVarDec:exit\nVarDecStmt:exit\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(key)\nIdent(foo)\nProp:exit\nObjPat:exit\nIdent(obj)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_PatObjRest(t *testing.T) {
	ast, symtab, err := compile(`
  let { a, b, ...rest } = obj
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nObjPat:enter\nProp:enter\nIdent(a)\nIdent(a)\nProp:exit\nProp:enter\nIdent(b)\nIdent(b)\nProp:exit\nRestPat:enter\nIdent(rest)\nRestPat:exit\nObjPat:exit\nIdent(obj)\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ArrowFn(t *testing.T) {
	ast, symtab, err := compile(`
  let f = () => {}
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(f)\nArrowFn\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ArrowFnExpr(t *testing.T) {
	ast, symtab, err := compile(`
  let f = () => b
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(f)\nArrowFn\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ArrowFnExprBody(t *testing.T) {
	ast, symtab, err := compile(`
  let f = () => a && b || c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	varDec := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt)
	fn := varDec.DecList()[0].(*parser.VarDec).Init()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b13[label="ArrowFn:enter\nBinExpr(||):enter\nBinExpr(&&):enter\nIdent(a)\n"];
b14[label="ArrowFn:exit\n"];
b5[label="Ident(b)\nBinExpr(&&):exit\n"];
b9[label="Ident(c)\nBinExpr(||):exit\n"];
b13->b5 [xlabel="",color="black"];
b13->b9 [xlabel="F",color="orange"];
b14->final [xlabel="",color="black"];
b5->b14 [xlabel="T",color="orange"];
b5->b9 [xlabel="",color="black"];
b9->b14 [xlabel="",color="black"];
initial->b13 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ArrowFnArgs(t *testing.T) {
	ast, symtab, err := compile(`
  let f = (a, b) => {
    c
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	varDec := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt)
	fn := varDec.DecList()[0].(*parser.VarDec).Init()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b10[label="ArrowFn:enter\nIdent(a)\nIdent(b)\nBlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\nArrowFn:exit\n"];
b10->final [xlabel="",color="black"];
initial->b10 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ArrowFnRet(t *testing.T) {
	ast, symtab, err := compile(`
  let f = (a, b) => {
    return
    c
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	varDec := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt)
	fn := varDec.DecList()[0].(*parser.VarDec).Init()
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b12[label="ArrowFn:enter\nIdent(a)\nIdent(b)\nBlockStmt:enter\nRetStmt:enter\nRetStmt:exit\n"];
b13[label="ArrowFn:exit\n"];
b8[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\nBlockStmt:exit\n"];
b12->b13 [xlabel="U",color="orange"];
b12->b8 [xlabel="",color="red"];
b13->final [xlabel="",color="black"];
b8->b13 [xlabel="",color="red"];
initial->b12 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_Yield(t *testing.T) {
	ast, symtab, err := compile(`
  function* f() {
    yield 1
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
b13[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nYieldExpr:enter\nNumLit(1)\nYieldExpr:exit\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b13->final [xlabel="",color="black"];
initial->b13 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_YieldStar(t *testing.T) {
	ast, symtab, err := compile(`
  function* f() {
    yield* ff()
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
b16[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nYieldExpr:enter\nCallExpr:enter\nIdent(ff)\nCallExpr:exit\nYieldExpr:exit\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b16->final [xlabel="",color="black"];
initial->b16 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_OptChain(t *testing.T) {
	ast, symtab, err := compile(`
  obj?.prop
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nMemberExpr:enter\nIdent(obj)\n"];
b10[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b7[label="Ident(prop)\nMemberExpr:exit\n"];
b0->b10 [xlabel="N",color="orange"];
b0->b7 [xlabel="",color="black"];
b10->final [xlabel="",color="black"];
b7->b10 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainCompute(t *testing.T) {
	ast, symtab, err := compile(`
  (a || b)?.prop
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nMemberExpr:enter\nParenExpr:enter\nBinExpr(||):enter\nIdent(a)\n"];
b13[label="ParenExpr:exit\n"];
b15[label="Ident(prop)\nMemberExpr:exit\n"];
b18[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="Ident(b)\nBinExpr(||):exit\n"];
b0->b13 [xlabel="T",color="orange"];
b0->b9 [xlabel="",color="black"];
b13->b15 [xlabel="",color="black"];
b13->b18 [xlabel="N",color="orange"];
b15->b18 [xlabel="",color="black"];
b18->final [xlabel="",color="black"];
b9->b13 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainMember(t *testing.T) {
	ast, symtab, err := compile(`
  obj.val?.prop
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(obj)\nIdent(val)\nMemberExpr:exit\n"];
b11[label="Ident(prop)\nMemberExpr:exit\n"];
b14[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b11 [xlabel="",color="black"];
b0->b14 [xlabel="N",color="orange"];
b11->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainCallee(t *testing.T) {
	ast, symtab, err := compile(`
  obj.func?.()
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nCallExpr:enter\nMemberExpr:enter\nIdent(obj)\nIdent(func)\nMemberExpr:exit\n"];
b11[label="CallExpr:exit\n"];
b13[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b11 [xlabel="",color="black"];
b0->b13 [xlabel="F",color="orange"];
b11->b13 [xlabel="",color="black"];
b13->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainCalleeArgs(t *testing.T) {
	ast, symtab, err := compile(`
  obj.func?.(args)
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nCallExpr:enter\nMemberExpr:enter\nIdent(obj)\nIdent(func)\nMemberExpr:exit\n"];
b11[label="Ident(args)\nCallExpr:exit\n"];
b14[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b11 [xlabel="",color="black"];
b0->b14 [xlabel="F",color="orange"];
b11->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainIdx(t *testing.T) {
	ast, symtab, err := compile(`
  obj.arr?.[index]
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(obj)\nIdent(arr)\nMemberExpr:exit\n"];
b11[label="Ident(index)\nMemberExpr:exit\n"];
b14[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b11 [xlabel="",color="black"];
b0->b14 [xlabel="N",color="orange"];
b11->b14 [xlabel="",color="black"];
b14->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_OptChainNested(t *testing.T) {
	ast, symtab, err := compile(`
  a?.b?.c()
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
b0[label="Prog:enter\nExprStmt:enter\nChainExpr:enter\nCallExpr:enter\nMemberExpr:enter\nMemberExpr:enter\nIdent(a)\n"];
b12[label="Ident(c)\nMemberExpr:exit\nCallExpr:exit\n"];
b17[label="ChainExpr:exit\nExprStmt:exit\nProg:exit\n"];
b9[label="Ident(b)\nMemberExpr:exit\n"];
b0->b17 [xlabel="N",color="orange"];
b0->b9 [xlabel="",color="black"];
b12->b17 [xlabel="",color="black"];
b17->final [xlabel="",color="black"];
b9->b12 [xlabel="",color="black"];
b9->b17 [xlabel="N",color="orange"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TplStr(t *testing.T) {
	ast, symtab, err := compile("`string text`", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nTplExpr:enter\nStrLit\nTplExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TplStrExpr(t *testing.T) {
	ast, symtab, err := compile("`string text ${a && b} string ${c} text`", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nTplExpr:enter\nStrLit\nBinExpr(&&):enter\nIdent(a)\n"];
b12[label="StrLit\nIdent(c)\nStrLit\nTplExpr:exit\nExprStmt:exit\nProg:exit\n"];
b8[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->final [xlabel="",color="black"];
b8->b12 [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TplStrExprLast(t *testing.T) {
	ast, symtab, err := compile("`string text ${a} string ${b} text ${c && d}`", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nTplExpr:enter\nStrLit\nIdent(a)\nStrLit\nIdent(b)\nStrLit\nBinExpr(&&):enter\nIdent(c)\n"];
b12[label="Ident(d)\nBinExpr(&&):exit\n"];
b16[label="StrLit\nTplExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->b12 [xlabel="",color="black"];
b0->b16 [xlabel="F",color="orange"];
b12->b16 [xlabel="",color="black"];
b16->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_TplStrFn(t *testing.T) {
	ast, symtab, err := compile("tagFn`string text ${a} string text`", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nExprStmt:enter\nTplExpr:enter\nIdent(tagFn)\nStrLit\nIdent(a)\nStrLit\nTplExpr:exit\nExprStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_MetaProp(t *testing.T) {
	ast, symtab, err := compile(`
  function f() {
    new.target
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
b13[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nMetaProp:enter\nIdent(new)\nIdent(target)\nMetaProp:exit\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b13->final [xlabel="",color="black"];
initial->b13 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_JsxBasic(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div></div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxOpen:exit\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxSelfClosed(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div/>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen(Closed):enter\nJsxIdent(div)\nJsxOpen(Closed):exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxExpr(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div>a {b && c} d</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxOpen:exit\nJsxText\nJsxExprSpan:enter\nBinExpr(&&):enter\nIdent(b)\n"];
b15[label="Ident(c)\nBinExpr(&&):exit\n"];
b19[label="JsxExprSpan:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->b15 [xlabel="",color="black"];
b0->b19 [xlabel="F",color="orange"];
b15->b19 [xlabel="",color="black"];
b19->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxEmpty(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div>{/* empty */}</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxOpen:exit\nJsxExprSpan:enter\nJsxEmpty\nJsxExprSpan:exit\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxMember(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <A.b/>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen(Closed):enter\nJsxMember:enter\nJsxIdent(A)\nIdent(b)\nJsxMember:exit\nJsxOpen(Closed):exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxNs(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_JSX_NS)

	ast, symtab, err := compile(`
  let a = <div:a>text</div:a>
  `, opts)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxNsName\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxNsName\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxSpread(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div>
    {...e}
  </div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxOpen:exit\nJsxText\nJsxSpreadChild:enter\nIdent(e)\nJsxSpreadChild:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttr(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div attr="str">text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxAttr:enter\nJsxIdent(attr)\nStrLit\nJsxAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttrExpr(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div attr={1}>text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxAttr:enter\nJsxIdent(attr)\nJsxExprSpan:enter\nNumLit(1)\nJsxExprSpan:exit\nJsxAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttrExprBin(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div attr={a && b}>text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxAttr:enter\nJsxIdent(attr)\nJsxExprSpan:enter\nBinExpr(&&):enter\nIdent(a)\n"];
b14[label="Ident(b)\nBinExpr(&&):exit\n"];
b18[label="JsxExprSpan:exit\nJsxAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->b14 [xlabel="",color="black"];
b0->b18 [xlabel="F",color="orange"];
b14->b18 [xlabel="",color="black"];
b18->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttrFlag(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div attr>text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxAttr:enter\nJsxIdent(attr)\nJsxAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttrSpread(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div {...props}>text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxSpreadAttr:enter\nSpread:enter\nIdent(props)\nSpread:exit\nJsxSpreadAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_JsxAttrMixed(t *testing.T) {
	ast, symtab, err := compile(`
  let a = <div attr0 attr1={true} attr2="a" {...b}>text</div>
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nJsxElem:enter\nJsxOpen:enter\nJsxIdent(div)\nJsxAttr:enter\nJsxIdent(attr0)\nJsxAttr:exit\nJsxAttr:enter\nJsxIdent(attr1)\nJsxExprSpan:enter\nBoolLit\nJsxExprSpan:exit\nJsxAttr:exit\nJsxAttr:enter\nJsxIdent(attr2)\nStrLit\nJsxAttr:exit\nJsxSpreadAttr:enter\nSpread:enter\nIdent(b)\nSpread:exit\nJsxSpreadAttr:exit\nJsxOpen:exit\nJsxText\nJsxClose:enter\nJsxIdent(div)\nJsxClose:exit\nJsxElem:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ClassStmt(t *testing.T) {
	ast, symtab, err := compile(`
class A {}
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
b0[label="Prog:enter\nClassDec:enter\nIdent(A)\nClassBody:enter\nClassBody:exit\nClassDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ClassField(t *testing.T) {
	ast, symtab, err := compile(`
class A {
  a;
  b = 1
  c = a && b
  d = function f() {}
  e = () => {}
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
b0[label="Prog:enter\nClassDec:enter\nIdent(A)\nClassBody:enter\nField:enter\nIdent(a)\nField:exit\nField:enter\nIdent(b)\nNumLit(1)\nField:exit\nField:enter\nIdent(c)\nBinExpr(&&):enter\nIdent(a)\n"];
b19[label="Ident(b)\nBinExpr(&&):exit\n"];
b23[label="Field:exit\nField:enter\nIdent(d)\nFnDec\nField:exit\nField:enter\nIdent(e)\nArrowFn\nField:exit\nClassBody:exit\nClassDec:exit\nProg:exit\n"];
b0->b19 [xlabel="",color="black"];
b0->b23 [xlabel="F",color="orange"];
b19->b23 [xlabel="",color="black"];
b23->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ClassFieldFn(t *testing.T) {
	ast, symtab, err := compile(`
class A {
  a = function f() {
    b
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Field).Val().(*parser.FnDec)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b9[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b9->final [xlabel="",color="black"];
initial->b9 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ClassFieldArrowFn(t *testing.T) {
	ast, symtab, err := compile(`
class A {
  a = () => {
    b
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Field).Val().(*parser.ArrowFn)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b8[label="ArrowFn:enter\nBlockStmt:enter\nExprStmt:enter\nIdent(b)\nExprStmt:exit\nBlockStmt:exit\nArrowFn:exit\n"];
b8->final [xlabel="",color="black"];
initial->b8 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ClassFieldArrowFnExpr(t *testing.T) {
	ast, symtab, err := compile(`
class A {
  a = () => a?.b()
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Field).Val().(*parser.ArrowFn)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b11[label="ChainExpr:exit\nArrowFn:exit\n"];
b14[label="ArrowFn:enter\nChainExpr:enter\nCallExpr:enter\nMemberExpr:enter\nIdent(a)\n"];
b6[label="Ident(b)\nMemberExpr:exit\nCallExpr:exit\n"];
b11->final [xlabel="",color="black"];
b14->b11 [xlabel="N",color="orange"];
b14->b6 [xlabel="",color="black"];
b6->b11 [xlabel="",color="black"];
initial->b14 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ClassExpr(t *testing.T) {
	ast, symtab, err := compile(`
  let a = class A {}
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
b0[label="Prog:enter\nVarDecStmt:enter\nVarDec:enter\nIdent(a)\nClassDec:enter\nIdent(A)\nClassBody:enter\nClassBody:exit\nClassDec:exit\nVarDec:exit\nVarDecStmt:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ClassMethod(t *testing.T) {
	ast, symtab, err := compile(`
  class A {
    f () {}
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
b0[label="Prog:enter\nClassDec:enter\nIdent(A)\nClassBody:enter\nMethod:enter\nIdent(f)\nFnDec\nMethod:exit\nClassBody:exit\nClassDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

func TestCtrlflow_ClassMethodBdy(t *testing.T) {
	ast, symtab, err := compile(`
  class A {
    f () {
      a
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Method).Val().(*parser.FnDec)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b8[label="FnDec:enter\nBlockStmt:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b8->final [xlabel="",color="black"];
initial->b8 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ClassCtor(t *testing.T) {
	ast, symtab, err := compile(`
  class A extends B {
    constructor () {
      super()
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Method).Val().(*parser.FnDec)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b11[label="FnDec:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nSuper\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b11->final [xlabel="",color="black"];
initial->b11 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_ClassSuper(t *testing.T) {
	ast, symtab, err := compile(`
  class A {
    f () {
      super.a()
    }
  }
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()

	fn := ast.(*parser.Prog).Body()[0].(*parser.ClassDec).Body().(*parser.ClassBody).Elems()[0].(*parser.Method).Val().(*parser.FnDec)
	fnGraph := ana.AnalysisCtx().GraphOf(fn)

	AssertEqualString(t, `
digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
b15[label="FnDec:enter\nBlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nMemberExpr:enter\nSuper\nIdent(a)\nMemberExpr:exit\nCallExpr:exit\nExprStmt:exit\nBlockStmt:exit\nFnDec:exit\n"];
b15->final [xlabel="",color="black"];
initial->b15 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

func TestCtrlflow_StaticBlk(t *testing.T) {
	ast, symtab, err := compile(`
  class A {
    static {
      a
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
b0[label="Prog:enter\nClassDec:enter\nClassBody:enter\nStaticBlock:enter\nExprStmt:enter\nIdent(a)\nExprStmt:exit\nStaticBlock:exit\nClassBody:exit\nClassDec:exit\nProg:exit\n"];
b0->final [xlabel="",color="black"];
initial->b0 [xlabel="",color="black"];
}
`, ana.Graph().Dot(), "should be ok")
}

// func TestCtrlflow_Demo(t *testing.T) {
// 	ast, symtab, err := compile(`
//   var a = function f() { if (baz) { return true; } return 1 }
//   `, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	ana := NewAnalysis(ast, symtab)
// 	ana.Analyze()

// 	fn := ast.(*parser.Prog).Body()[0].(*parser.VarDecStmt).DecList()[0].(*parser.VarDec).Init()
// 	fnGraph := ana.AnalysisCtx().GraphOf(fn)

// 	AssertEqualString(t, `

// 	`, fnGraph.Dot(), "should be ok")

// 	// 	AssertEqualString(t, `

// 	// `, ana.Graph().Dot(), "should be ok")
// }

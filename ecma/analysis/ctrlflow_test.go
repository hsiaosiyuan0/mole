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
b10[label="ExprStmt:exit\nProg:exit\n"];
b6[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b10 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b10->final [xlabel="",color="black"];
b6->b10 [xlabel="",color="black"];
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
b12[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\n"];
b17[label="IfStmt:exit\nExprStmt:enter\nIdent(e)\nExprStmt:exit\nProg:exit\n"];
b8[label="ExprStmt:enter\nIdent(c)\nExprStmt:exit\n"];
b0->b12 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b12->b17 [xlabel="",color="black"];
b17->final [xlabel="",color="black"];
b8->b17 [xlabel="",color="black"];
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
b22[label="IfStmt:exit\nExprStmt:enter\nIdent(f)\nExprStmt:exit\nProg:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nIdent(c)\nExprStmt:exit\nExprStmt:enter\nIdent(d)\nExprStmt:exit\nBlockStmt:exit\n"];
b0->b17 [xlabel="F",color="orange"];
b0->b8 [xlabel="",color="black"];
b17->b22 [xlabel="",color="black"];
b22->final [xlabel="",color="black"];
b8->b22 [xlabel="",color="black"];
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
b14[label="ExprStmt:enter\nIdent(d)\nExprStmt:exit\n"];
b19[label="IfStmt:exit\nProg:exit\n"];
b6[label="Ident(b)\nBinExpr(&&):exit\n"];
b0->b14 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b10->b19 [xlabel="",color="black"];
b14->b19 [xlabel="",color="black"];
b19->final [xlabel="",color="black"];
b6->b10 [xlabel="",color="black"];
b6->b14 [xlabel="F",color="orange"];
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
b22[label="ExprStmt:enter\nIdent(f)\nExprStmt:exit\n"];
b27[label="IfStmt:exit\nProg:exit\n"];
b9[label="BinExpr(&&):enter\nIdent(c)\n"];
b0->b18 [xlabel="T",color="orange"];
b0->b9 [xlabel="",color="black"];
b11->b18 [xlabel="",color="black"];
b11->b22 [xlabel="F",color="orange"];
b18->b27 [xlabel="",color="black"];
b22->b27 [xlabel="",color="black"];
b27->final [xlabel="",color="black"];
b9->b11 [xlabel="",color="black"];
b9->b22 [xlabel="F",color="orange"];
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
b10[label="WhileStmt:exit\nProg:exit\n"];
b4[label="Ident(a)\n"];
b5[label="ExprStmt:enter\nIdent(b)\nExprStmt:exit\n"];
b0->b4 [xlabel="",color="black"];
b10->final [xlabel="",color="black"];
b4->b10 [xlabel="F",color="orange"];
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
b4:s->b4:ne [xlabel="T",color="orange"];
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
b24->b26 [xlabel="",color="black"];
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
b31[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b40[label="SwitchCase:enter\nNumLit(2)\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSthDefault)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b27 [xlabel="",color="black"];
b18->b31 [xlabel="",color="black"];
b27->b18 [xlabel="",color="black"];
b27->b40 [xlabel="F",color="orange"];
b31->final [xlabel="",color="black"];
b40->b31 [xlabel="",color="black"];
b40->b6 [xlabel="F",color="orange"];
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
b44[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b53[label="SwitchCase:enter\nNumLit(3)\n"];
b6[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\n"];
b0->b27 [xlabel="F",color="orange"];
b0->b6 [xlabel="",color="black"];
b19->b31 [xlabel="",color="black"];
b27->b40 [xlabel="",color="black"];
b31->b44 [xlabel="",color="black"];
b40->b31 [xlabel="",color="black"];
b40->b53 [xlabel="F",color="orange"];
b44->final [xlabel="",color="black"];
b53->b19 [xlabel="F",color="orange"];
b53->b44 [xlabel="",color="black"];
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
b37[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth3)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\nProg:exit\n"];
b46[label="SwitchCase:enter\nNumLit(3)\n"];
b9[label="SwitchCase:exit\n"];
b0->b20 [xlabel="F",color="orange"];
b0->b9 [xlabel="",color="black"];
b12->b24 [xlabel="",color="black"];
b20->b33 [xlabel="",color="black"];
b24->b37 [xlabel="",color="black"];
b33->b24 [xlabel="",color="black"];
b33->b46 [xlabel="F",color="orange"];
b37->final [xlabel="",color="black"];
b46->b12 [xlabel="F",color="orange"];
b46->b37 [xlabel="",color="black"];
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
b37[label="ExprStmt:enter\nCallExpr:enter\nIdent(doSth2)\nCallExpr:exit\nExprStmt:exit\nSwitchCase:exit\nSwitchStmt:exit\n"];
b46[label="SwitchCase:enter\nNumLit(2)\n"];
b50[label="BlockStmt:exit\n"];
b51[label="FnDec:enter\nIdent(f)\nBlockStmt:enter\nSwitchStmt:enter\nIdent(a)\nSwitchCase:enter\nNumLit(1)\n"];
b52[label="FnDec:exit\n"];
b8[label="BlockStmt:enter\nExprStmt:enter\nCallExpr:enter\nIdent(doSth1)\nCallExpr:exit\nExprStmt:exit\nRetStmt:enter\nRetStmt:exit\n"];
b20->b25 [xlabel="",color="red"];
b25->b37 [xlabel="",color="black"];
b33->b46 [xlabel="",color="black"];
b37->b50 [xlabel="",color="black"];
b46->b25 [xlabel="F",color="orange"];
b46->b37 [xlabel="",color="black"];
b50->b52 [xlabel="",color="black"];
b51->b33 [xlabel="F",color="orange"];
b51->b8 [xlabel="",color="black"];
b52->final [xlabel="",color="black"];
b8->b20 [xlabel="",color="red"];
b8->b52 [xlabel="U",color="orange"];
initial->b51 [xlabel="",color="black"];
}
`, fnGraph.Dot(), "should be ok")
}

// func TestCtrlflow_TryCatchBasic(t *testing.T) {
// 	ast, symtab, err := compile(`
// try {
//   doSth1();
//   doSth2();
// } catch (error) {
//   log(error);
// }
//   `, nil)
// 	AssertEqual(t, nil, err, "should be prog ok")

// 	ana := NewAnalysis(ast, symtab)
// 	ana.Analyze()

// 	AssertEqualString(t, `

// `, ana.Graph().Dot(), "should be ok")
// }

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

func TestLogic(t *testing.T) {
	ast, symtab, err := compile(`
 a && b || c
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()
}

func TestIfStmt(t *testing.T) {
	ast, symtab, err := compile(`
  a;
  if (b) c;
  else d;
  e;
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ana := NewAnalysis(ast, symtab)
	ana.Analyze()
}

func TestIfBlkStmt(t *testing.T) {
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
}

func TestIfBlk2Stmt(t *testing.T) {
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
}

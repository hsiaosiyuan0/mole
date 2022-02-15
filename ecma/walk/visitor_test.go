package walk

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

func TestVisitor(t *testing.T) {
	ast, symtab, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	names := make([]string, 0)

	ctx := NewWalkCtx(ast, symtab)
	AddVisitor(&ctx.Visitors, VK_NAME, func(node parser.Node, key string, ctx *WalkCtx) {
		n := node.(*parser.Ident)
		names = append(names, n.Text())
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 3, len(names), "should be 3 names")
}

func TestVisitorFnScope(t *testing.T) {
	ast, symtab, err := compile(`
function a() {
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	AddVisitor(&ctx.Visitors, VK_STMT_BLOCK, func(node parser.Node, key string, ctx *WalkCtx) {
		fnScopeId = ctx.ScopeId()
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 1, fnScopeId, "should be ok")
}

func TestVisitorFnNestedScope(t *testing.T) {
	ast, symtab, err := compile(`
function a() {
  function b() {
  }
}
  `, nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ctx := NewWalkCtx(ast, symtab)

	fnScopeId := 0
	AddVisitor(&ctx.Visitors, VK_STMT_BLOCK, func(node parser.Node, key string, ctx *WalkCtx) {
		if ctx.ScopeId() > fnScopeId {
			fnScopeId = ctx.ScopeId()
		}
	})

	VisitNode(ast, "", ctx)
	AssertEqual(t, 2, fnScopeId, "should be ok")
}

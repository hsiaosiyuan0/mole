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

func compile(code string, opts *parser.ParserOpts) (parser.Node, error) {
	p := newParser(code, opts)
	return p.Prog()
}

func TestVisitor(t *testing.T) {
	ast, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	names := make([]string, 0)

	ctx := NewWalkCtx(ast, nil)
	AddVisitor(&ctx.Visitors, VK_NAME, func(node parser.Node, ctx *WalkCtx) {
		n := node.(*parser.Ident)
		names = append(names, n.Text())
	})

	VisitNode(ast, ctx)
	AssertEqual(t, 3, len(names), "should be 3 names")
}

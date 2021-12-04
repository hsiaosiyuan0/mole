package indent

import (
	"strings"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

type Indenter struct {
	opts *parser.ParserOpts
	vi   *parser.NodeFns
	sb   strings.Builder
}

func NewIndenter(opts *parser.ParserOpts) *Indenter {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	vi := parser.NewVisitorImpl()

	vi[parser.N_EXPR_BIN] = func(n parser.Node, v interface{}, ctx *parser.TraverseCtx) parser.ContOrStop {
		node := n.(*parser.BinExpr)
		it := ctx.Extra.(*Indenter)
		parser.VisitNode(node.Lhs(), v, ctx, false)
		it.sb.WriteByte(' ')
		it.sb.WriteString(node.OpText())
		it.sb.WriteByte(' ')
		parser.VisitNode(node.Rhs(), v, ctx, false)
		return true
	}

	vi[parser.N_NAME] = func(n parser.Node, v interface{}, ctx *parser.TraverseCtx) parser.ContOrStop {
		id := n.(*parser.Ident)
		it := ctx.Extra.(*Indenter)
		it.sb.WriteString(id.Text())
		return true
	}
	return &Indenter{opts, vi, strings.Builder{}}
}

func (it *Indenter) Process(code string, file string) (string, error) {
	s := parser.NewSource(file, code)
	p := parser.NewParser(s, it.opts)
	ast, err := p.Prog()
	if err != nil {
		return "nil", err
	}

	ctx := &parser.TraverseCtx{}
	ctx.Extra = it
	parser.VisitNode(ast, it.vi, ctx, false)
	return it.sb.String(), nil
}

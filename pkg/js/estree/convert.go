package estree

import (
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func pos(p *parser.Pos) *Position {
	return &Position{Line: p.Line(), Column: p.Column()}
}

func loc(s *parser.Loc) *SrcLoc {
	return &SrcLoc{
		Source: s.Source(),
		Start:  pos(s.Begin()),
		End:    pos(s.End()),
		Range:  rng(s),
	}
}

func rng(s *parser.Loc) *SrcRange {
	rng := s.Range()
	return &SrcRange{
		Start: rng.Start(),
		End:   rng.End(),
	}
}

func program(n *parser.Prog) *Program {
	body := make([]Node, len(n.Body()))
	for i, s := range n.Body() {
		body[i] = convert(s)
	}
	return &Program{
		Type: "Program",
		Loc:  loc(n.Loc()),
		Body: body,
	}
}

func convert(node parser.Node) Node {
	if node == nil {
		return nil
	}
	switch node.Type() {
	case parser.N_STMT_EXPR:
		expr := convert(node.(*parser.ExprStmt).Expr())
		return &ExpressionStatement{
			Type:       "ExpressionStatement",
			Loc:        loc(node.Loc()),
			Expression: expr,
		}
	case parser.N_EXPR_NEW:
		expr := convert(node.(*parser.NewExpr).Expr())
		return &NewExpression{
			Type:   "NewExpression",
			Loc:    loc(node.Loc()),
			Callee: expr,
		}
	case parser.N_NAME:
		id := node.(*parser.Ident)
		name := id.Text()
		return &Identifier{
			Type: "Identifier",
			Loc:  loc(id.Loc()),
			Name: name,
		}
	case parser.N_EXPR_THIS:
		return &ThisExpression{
			Type: "ThisExpression",
			Loc:  loc(node.Loc()),
		}
	case parser.N_LIT_NULL:
		return &Literal{
			Type: "Literal",
			Loc:  loc(node.Loc()),
		}
	case parser.N_LIT_NUM:
		return &Literal{
			Type:  "Literal",
			Loc:   loc(node.Loc()),
			Value: node.(*parser.NumLit).ToFloat(),
		}
	case parser.N_LIT_REGEXP:
		regexp := node.(*parser.RegexpLit)
		return &RegExpLiteral{
			Type:   "Literal",
			Loc:    loc(node.Loc()),
			Regexp: &Regexp{regexp.Pattern(), regexp.Flags()},
		}
	case parser.N_EXPR_BIN:
		bin := node.(*parser.BinExpr)
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
		op := bin.Op().Text()
		return &BinaryExpression{
			Type:     "BinaryExpression",
			Loc:      loc(node.Loc()),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	case parser.N_EXPR_ASSIGN:
		bin := node.(*parser.AssignExpr)
		lhs := convert(bin.Lhs())
		rhs := convert(bin.Rhs())
		op := bin.Op().Text()
		return &AssignmentExpression{
			Type:     "AssignmentExpression",
			Loc:      loc(node.Loc()),
			Operator: op,
			Left:     lhs,
			Right:    rhs,
		}
	}
	return nil
}

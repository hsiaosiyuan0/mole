package estree

import (
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func NewPosition(p *parser.Pos) *Position {
	return &Position{Line: p.Line(), Column: p.Column()}
}

func NewSrcLoc(s *parser.Loc) *SrcLoc {
	return &SrcLoc{
		Source: s.Source(),
		Start:  NewPosition(s.Begin()),
		End:    NewPosition(s.End()),
	}
}

func NewProgram(n *parser.Prog) *Program {
	body := make([]Node, len(n.Body()))
	for i, s := range n.Body() {
		body[i] = ToESNode(s)
	}
	return &Program{
		Type: "program",
		Loc:  NewSrcLoc(n.Loc()),
		Body: body,
	}
}

func ToESNode(node parser.Node) Node {
	if node == nil {
		return nil
	}
	switch node.Type() {
	case parser.N_STMT_EXPR:
		expr := ToESNode(node.(*parser.ExprStmt).Expr())
		return &ExpressionStatement{
			Type:       "ExpressionStatement",
			Loc:        NewSrcLoc(node.Loc()),
			Expression: expr,
		}
	case parser.N_EXPR_NEW:
		expr := ToESNode(node.(*parser.NewExpr).Expr())
		return &NewExpression{
			Type:   "NewExpression",
			Loc:    NewSrcLoc(node.Loc()),
			Callee: expr,
		}
	case parser.N_NAME:
		name := node.(*parser.Ident).Text()
		return &Identifier{
			Type: "Identifier",
			Name: name,
		}
	}
	return nil
}

package lint

import "github.com/hsiaosiyuan0/mole/ecma/parser"

func GetStaticPropertyName(node parser.Node) string {
	node = parser.UnParen(node)
	switch node.Type() {
	case parser.N_EXPR_MEMBER:
		n := node.(*parser.MemberExpr)
		if n.Prop().Type() == parser.N_NAME {
			return n.Prop().(*parser.Ident).Text()
		}
	case parser.N_EXPR_CHAIN:
		n := node.(*parser.ChainExpr)
		return GetStaticPropertyName(n.Expr())
	}
	return ""
}

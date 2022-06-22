package parser

func GetStaticPropertyName(node Node) string {
	node = UnParen(node)
	switch node.Type() {
	case N_EXPR_MEMBER:
		n := node.(*MemberExpr)
		if n.Prop().Type() == N_NAME {
			return n.Prop().(*Ident).Text()
		}
	case N_EXPR_CHAIN:
		n := node.(*ChainExpr)
		return GetStaticPropertyName(n.Expr())
	}
	return ""
}

func GetName(node Node) string {
	if node.Type() != N_NAME {
		return ""
	}
	return node.(*Ident).Text()
}

package parser

var builtinTyp = map[string]NodeType{
	"any":     N_TYP_ANY,
	"number":  N_TYP_NUM,
	"boolean": N_TYP_BOOL,
	"string":  N_TYP_STR,
	"symbol":  N_TYP_SYM,
	"void":    N_TYP_VOID,
}

func (p *Parser) typ() (Node, error) {
	return p.primaryTyp()
}

func (p *Parser) parenTyp() (Node, error) {
	return p.typ()
}

func (p *Parser) primaryTyp() (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	av := ahead.value
	if av == T_PAREN_L {
		return p.parenTyp()
	}
	if av == T_NAME {
		name := ahead.Text()
		if typ, ok := builtinTyp[name]; ok {
			p.lexer.Next()
			return &BuiltinType{typ, p.finLoc(loc)}, nil
		}
	}
	return p.primaryTyp()
}

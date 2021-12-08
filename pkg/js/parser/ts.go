package parser

var builtinTyp = map[string]NodeType{
	"any":     N_TS_ANY,
	"number":  N_TS_NUM,
	"boolean": N_TS_BOOL,
	"string":  N_TS_STR,
	"symbol":  N_TS_SYM,
	"void":    N_TS_VOID,
}

func (p *Parser) tsTyp() (Node, error) {
	return p.tsPrimaryTyp()
}

func (p *Parser) tsParenTyp() (Node, error) {
	p.lexer.Next()
	typ, err := p.tsTyp()
	if err != nil {
		return nil, err
	}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}
	return typ, nil
}

func (p *Parser) tsTypName(ns Node) (Node, error) {
	if ns == nil {
		var err error
		ns, err = p.ident(nil)
		if err != nil {
			return nil, err
		}
	}
	for {
		if p.lexer.Peek().value == T_DOT {
			p.lexer.Next()
			id, err := p.ident(nil)
			if err != nil {
				return nil, err
			}
			ns = &TsNsName{N_TS_NS_NAME, p.finLoc(ns.Loc().Clone()), ns, id}
		} else {
			break
		}
	}
	return ns, nil
}

func (p *Parser) tsTypRef(ns Node) (Node, error) {
	name, err := p.tsTypName(ns)
	if err != nil {
		return nil, err
	}
	// type args
	return &TsTypeRef{N_TS_TYP_REF, p.finLoc(ns.Loc().Clone()), name, nil}, nil
}

func (p *Parser) tsArrType(typ Node) (Node, error) {
	for {
		if p.lexer.Peek().value == T_BRACKET_L {
			p.lexer.Next()
			if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
				return nil, err
			}
			typ = &TsArrType{N_TS_ARR_TYP, p.finLoc(typ.Loc().Clone()), typ}
		} else {
			break
		}
	}
	return typ, nil
}

func (p *Parser) tsTupleType() (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	args := make([]Node, 0, 1)
	for {
		if p.lexer.Peek().value == T_BRACKET_R {
			break
		}
		arg, err := p.tsTyp()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return &TsTupleType{N_TS_TUPLE_TYP, p.finLoc(loc), args}, nil
}

func (p *Parser) tsObjType() (Node, error) {
	return nil, nil
}

func (p *Parser) tsPrimaryTyp() (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	av := ahead.value
	if av == T_PAREN_L {
		// paren type
		return p.tsParenTyp()
	}

	var err error
	var node Node
	if av == T_NAME {
		id, err := p.ident(nil)
		if err != nil {
			return nil, err
		}

		name := id.Text()
		if typ, ok := builtinTyp[name]; ok {
			// predef
			node = &TsPredefType{typ, id.loc}
		} else {
			node = id
		}

		node, err = p.tsTypRef(node)
		if err != nil {
			return nil, err
		}
	} else if av == T_BRACE_L {
		// obj type
		node, err = p.tsObjType()
		if err != nil {
			return nil, err
		}
	} else if av == T_BRACKET_L {
		// tuple type
		node, err = p.tsTupleType()
		if err != nil {
			return nil, err
		}
	} else if av == T_TYPE_OF {
		// type query
		p.lexer.Next()
		name, err := p.tsTypName(nil)
		if err != nil {
			return nil, err
		}

		node = &TsTypeQuery{N_TS_TYP_QUERY, p.finLoc(loc), name}
	} else if av == T_THIS {
		// this type
		p.lexer.Next()
		node = &TsThisType{N_TS_THIS_TYP, p.finLoc(loc)}
	}

	if node != nil {
		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_BRACKET_L && !ahead.afterLineTerminator {
			// array type
			node, err = p.tsArrType(node)
			if err != nil {
				return nil, err
			}
		}
		return node, nil
	}

	return nil, p.errorTok(ahead)
}

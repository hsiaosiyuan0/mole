package parser

var builtinTyp = map[string]NodeType{
	"any":     N_TS_ANY,
	"number":  N_TS_NUM,
	"boolean": N_TS_BOOL,
	"string":  N_TS_STR,
	"symbol":  N_TS_SYM,
	"void":    N_TS_VOID,
}

func (p *Parser) ts() bool {
	return p.feat&FEAT_TS != 0
}

func (p *Parser) tsTypAnnot() (Node, error) {
	if p.ts() && p.lexer.Peek().value == T_COLON {
		p.lexer.Next()
		return p.tsPrimary()
	}
	return nil, nil
}

// `var a:  () => number | string = 1` equals:
// `var a:  () => (number | string) = 1`
func (p *Parser) tsTyp() (Node, error) {
	return p.tsPrimary()
}

func (p *Parser) tsParen() (Node, error) {
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
		ns, err = p.ident(nil, false)
		if err != nil {
			return nil, err
		}
	}
	for {
		if p.lexer.Peek().value == T_DOT {
			p.lexer.Next()
			id, err := p.ident(nil, false)
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

func (p *Parser) tsRef(ns Node) (Node, error) {
	name, err := p.tsTypName(ns)
	if err != nil {
		return nil, err
	}
	// type args
	return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, nil}, nil
}

func (p *Parser) tsArr(typ Node) (Node, error) {
	for {
		if p.lexer.Peek().value == T_BRACKET_L {
			p.lexer.Next()
			if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
				return nil, err
			}
			typ = &TsArr{N_TS_ARR, p.finLoc(typ.Loc().Clone()), typ}
		} else {
			break
		}
	}
	return typ, nil
}

func (p *Parser) tsTuple() (Node, error) {
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
	return &TsTuple{N_TS_TUPLE, p.finLoc(loc), args}, nil
}

// returns `PrimaryType` or `FunctionType`
func (p *Parser) tsFnTyp() (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)

	tv := tok.value
	var typParams []Node
	var err error
	if tv == T_LET {
		typParams, _, err = p.tsTypParams()
		if err != nil {
			return nil, err
		}
		if _, err = p.nextMustTok(T_PAREN_L); err != nil {
			return nil, err
		}
	}

	params, parenL, err := p.paramList()
	if err != nil {
		return nil, err
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	if av != T_ARROW && typParams == nil {
		// ParenthesizedType
		pl := len(params)
		if pl == 1 {

		}
		if pl == 0 {
			return nil, p.errorAtLoc(parenL, ERR_UNEXPECTED_TOKEN)
		}
		if pl > 1 {
			return nil, p.errorAtLoc(params[pl-1].Loc(), ERR_UNEXPECTED_TOKEN)
		}
	}

	if _, err := p.nextMustTok(T_ARROW); err != nil {
		return nil, err
	}

	retTyp, err := p.tsTyp()
	if err != nil {
		return nil, err
	}

	return &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), typParams, params, retTyp}, nil
}

// returns `FunctionType` or `PrimaryType` since `FunctionType`
// is conflicts with `ParenthesizedType`
func (p *Parser) tsPrimary() (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	av := ahead.value
	if av == T_PAREN_L {
		// paren type
		return p.tsFnTyp()
	}

	var err error
	var node Node
	if av == T_NAME {
		id, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}

		name := id.Text()
		if typ, ok := builtinTyp[name]; ok {
			// predef
			node = &TsPredef{typ, id.loc}
		} else {
			node = id
		}

		node, err = p.tsRef(node)
		if err != nil {
			return nil, err
		}
	} else if av == T_BRACE_L {
		// obj type
		node, err = p.tsObj()
		if err != nil {
			return nil, err
		}
	} else if av == T_BRACKET_L {
		// tuple type
		node, err = p.tsTuple()
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

		node = &TsQuery{N_TS_QUERY, p.finLoc(loc), name}
	} else if av == T_THIS {
		// this type
		p.lexer.Next()
		node = &TsThis{N_TS_THIS, p.finLoc(loc)}
	}

	if node != nil {
		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_BRACKET_L && !ahead.afterLineTerminator {
			// array type
			node, err = p.tsArr(node)
			if err != nil {
				return nil, err
			}
		}
		return node, nil
	}

	return nil, p.errorTok(ahead)
}

func (p *Parser) tsObj() (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	props := make([]Node, 0, 1)
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACE_R {
			p.lexer.Next()
			break
		}
		prop, err := p.tsProp()
		if err != nil {
			return nil, err
		}
		props = append(props, prop)
		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_COMMA || av == T_SEMI {
			p.lexer.Next()
		}
	}
	return &TsObj{N_TS_OBJ, p.finLoc(loc), props}, nil
}

// func (p *Parser) paramList() ([]Node, *Loc, error) {
// 	scope := p.scope()
// 	p.checkName = false
// 	scope.AddKind(SPK_FORMAL_PARAMS)
// 	args, loc, err := p.argList(false, false)
// 	scope.EraseKind(SPK_FORMAL_PARAMS)
// 	p.checkName = true
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	params, err := p.argsToParams(args, false)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return params, loc, nil
// }

func (p *Parser) tsProp() (Node, error) {
	ahead := p.lexer.Peek()
	av := ahead.value
	loc := p.locFromTok(ahead)

	if av == T_LT {
		ps, _, err := p.tsTypParams()
		if err != nil {
			return nil, err
		}
		return p.tsCallSig(ps, loc)
	} else if av == T_NEW {
		// ConstructSignature
		typParams, _, err := p.tsTypParams()
		if err != nil {
			return nil, err
		}
		params, _, err := p.paramList()
		if err != nil {
			return nil, err
		}
		retTyp, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		return &TsNewSig{N_TS_NEW_SIG, p.finLoc(loc), typParams, params, retTyp}, nil
	} else if av == T_BRACKET_L {
		// IndexSignature
		p.lexer.Next()
		id, err := p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		if _, err = p.nextMustTok(T_COLON); err != nil {
			return nil, err
		}
		key, err := p.tsTyp()
		if err != nil {
			return nil, err
		}
		val, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		return &TsIdxSig{N_TS_IDX_SIG, p.finLoc(loc), id, key, val}, nil
	} else if av == T_PAREN_L {
		// CallSignature
		return p.tsCallSig(nil, loc)
	}

	// PropertySignature or MethodSignature
	name, err := p.tsPropName()
	if err != nil {
		return nil, err
	}
	var ques *Loc
	if p.lexer.Peek().value == T_COLON {
		ques = p.locFromTok(p.lexer.Next())
	}
	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	if typAnnot != nil {
		return &TsProp{N_TS_PROP, p.finLoc(name.Loc().Clone()), name, typAnnot, ques}, nil
	}

	// MethodSignature is deserved
	callSig, err := p.tsCallSig(nil, nil)
	if err != nil {
		return nil, err
	}
	return &TsProp{N_TS_PROP, p.finLoc(name.Loc().Clone()), name, callSig, ques}, nil
}

func (p *Parser) tsPropName() (Node, error) {
	tok := p.lexer.Peek()
	loc := p.locFromTok(tok)

	switch tok.value {
	case T_NUM:
		p.lexer.Next()
		return &NumLit{N_LIT_NUM, p.finLoc(loc), nil}, nil
	case T_STRING:
		p.lexer.Next()
		legacyOctalEscapeSeq := tok.HasLegacyOctalEscapeSeq()
		if p.scope().IsKind(SPK_STRICT) && legacyOctalEscapeSeq {
			return nil, p.errorAtLoc(p.finLoc(loc), ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE)
		}
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), legacyOctalEscapeSeq, nil}, nil
	case T_NAME:
		return p.ident(nil, false)
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) tsTypParam() (Node, error) {
	id, err := p.ident(nil, false)
	if err != nil {
		return nil, err
	}

	var cons Node
	if p.lexer.Peek().value == T_EXTENDS {
		p.lexer.Next()
		cons, err = p.tsTyp()
		if err != nil {
			return nil, err
		}
	}
	return &TsParam{N_TS_PARAM, p.finLoc(id.loc.Clone()), id, cons}, nil
}

func (p *Parser) tsTypParams() ([]Node, *Loc, error) {
	loc := p.locFromTok(p.lexer.Next())
	ps := make([]Node, 0, 1)
	for {
		pa, err := p.tsTypParam()
		if err != nil {
			return nil, nil, err
		}
		ps = append(ps, pa)

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_COMMA {
			p.lexer.Next()
		} else if av == T_GT {
			p.lexer.Next()
			break
		}
	}
	return ps, loc, nil
}

func (p *Parser) tsTypArgList() ([]Node, error) {
	p.lexer.Next()
	args := make([]Node, 0, 1)
	for {
		arg, err := p.tsTyp()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_COMMA {
			p.lexer.Next()
		} else if av == T_GT {
			p.lexer.Next()
			break
		}
	}
	return args, nil
}

func (p *Parser) tsCallSig(typParams []Node, loc *Loc) (Node, error) {
	if typParams == nil && loc == nil {
		var err error
		typParams, loc, err = p.tsTypParams()
		if err != nil {
			return nil, err
		}
	}
	params, _, err := p.paramList()
	if err != nil {
		return nil, err
	}
	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsCallSig{N_TS_CALL_SIG, p.finLoc(loc), typParams, params, typAnnot}, nil
}

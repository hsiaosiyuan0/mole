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

func (p *Parser) tsTypAnnot() (Node, *Loc, error) {
	if p.ts() && p.lexer.Peek().value == T_COLON {
		loc := p.locFromTok(p.lexer.Next())
		node, err := p.tsTyp(false)
		if err != nil {
			return nil, nil, err
		}
		return node, loc, nil
	}
	return nil, nil, nil
}

// for dealing with the ambiguous between `ParenthesizedType` and the `formalParamList`:
// `var a: ({a = c}:string|number,a:string) => number = 1`
//
// set `rough` to `true` to parse `{a = c, b?: number}` as a super set consists of `tsTyp`
// and `bindingPattern`, the parsed result will be judged by later process by the time
func (p *Parser) tsTyp(rough bool) (Node, error) {
	if p.lexer.Peek().value == T_NEW {
		return p.tsConstructTyp()
	}
	return p.tsUnionOrIntersecType(nil, rough)
}

func (p *Parser) tsUnionOrIntersecType(lhs Node, rough bool) (Node, error) {
	var err error
	if lhs == nil {
		lhs, err = p.tsPrimary(rough)
		if err != nil {
			return nil, err
		}
	}

	if lhs.Type() == N_TS_FN_TYP {
		return lhs, nil
	}

	var rhs Node
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BIT_OR {
			p.lexer.Next()
			rhs, err = p.tsPrimary(rough)
			if err != nil {
				return nil, err
			}
			lhs = &TsUnionTyp{N_TS_UNION_TYP, p.finLoc(lhs.Loc()), lhs, rhs}
		} else if av == T_BIT_AND {
			p.lexer.Next()
			rhs, err = p.tsPrimary(rough)
			if err != nil {
				return nil, err
			}
			lhs = &TsIntersecTyp{N_TS_INTERSEC_TYP, p.finLoc(lhs.Loc()), lhs, rhs}
		} else {
			break
		}
	}
	return lhs, nil
}

func (p *Parser) tsConstructTyp() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	typParams, _, err := p.tsTypParams()
	if err != nil {
		return nil, err
	}
	params, _, err := p.paramList()
	if err != nil {
		return nil, err
	}
	if _, err = p.nextMustTok(T_ARROW); err != nil {
		return nil, err
	}
	retTyp, err := p.tsTyp(false)
	if err != nil {
		return nil, err
	}
	return &TsNewSig{N_TS_NEW_SIG, p.finLoc(loc), typParams, params, retTyp}, nil
}

// convert rough param to normal param
func (p *Parser) tsRoughParamToParam(node Node) (Node, error) {
	// param := node.(*TsRoughParam)
	// name := param.name
	// switch name.Type() {
	// case N_TS_ANY:
	// case N_TS_NUM:
	// case N_TS_BOOL:
	// case N_TS_STR:
	// case N_TS_SYM:
	// case N_TS_VOID:
	// case N_TS_REF:
	// case N_TS_OBJ:
	// case N_TS_ARR:
	// case N_TS_TUPLE:

	// }
	// return nil, p.errorAtLoc(node.Loc(), ERR_UNEXPECTED_TOKEN)
	return node, nil
}

func (p *Parser) tsRoughParamToTyp(node Node) (Node, error) {
	// param := node.(*TsRoughParam)
	// name := param.name
	// switch name.Type() {
	// case N_TS_ANY:
	// case N_TS_NUM:
	// case N_TS_BOOL:
	// case N_TS_STR:
	// case N_TS_SYM:
	// case N_TS_VOID:
	// case N_TS_REF:
	// case N_TS_OBJ:

	// }
	// return nil, p.errorAtLoc(node.Loc(), ERR_UNEXPECTED_TOKEN)
	return node, nil
}

// `ParenthesizedType` or`FunctionType`
func (p *Parser) tsParen() (Node, error) {
	params, loc, err := p.paramList()
	if err != nil {
		return nil, err
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	if av == T_ARROW {
		// check the first param
		if len(params) >= 1 {
			params[0], err = p.tsRoughParamToParam(params[0])
			if err != nil {
				return nil, err
			}
		}
		return p.tsFnTyp(params, loc)
	}
	if len(params) == 1 {
		param := params[0].(*TsRoughParam)
		if param.colon != nil {
			return nil, p.errorAtLoc(param.colon, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToTyp(param)
	} else if len(params) == 0 {
		return &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), nil, params, nil}, nil
	}
	return nil, p.errorTok(ahead)
}

// returns `PrimaryType` or `FunctionType`
func (p *Parser) tsFnTyp(params []Node, parenL *Loc) (Node, error) {
	var typParams []Node
	var err error
	var loc *Loc
	if parenL != nil {
		loc = parenL
	}
	if params == nil {
		tok := p.lexer.Next()
		loc = p.locFromTok(tok)
		tv := tok.value

		if tv == T_LT {
			typParams, _, err = p.tsTypParams()
			if err != nil {
				return nil, err
			}
			if _, err = p.nextMustTok(T_PAREN_L); err != nil {
				return nil, err
			}
		}

		params, _, err = p.paramList()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.nextMustTok(T_ARROW); err != nil {
		return nil, err
	}

	retTyp, err := p.tsTyp(false)
	if err != nil {
		return nil, err
	}

	return &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), typParams, params, retTyp}, nil
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
	if p.lexer.Peek().value != T_LT {
		return name, nil
	}
	args, err := p.tsTypArgs()
	if err != nil {
		return nil, err
	}
	return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, args}, nil
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

func (p *Parser) tsTuple(rough bool) (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	args := make([]Node, 0, 1)

	var arg Node
	var err error
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACKET_R {
			p.lexer.Next()
			break
		} else if av == T_DOT_TRI {
			arg, err = p.patternRest(true, true)
		} else {
			arg, err = p.tsTyp(rough)
		}
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
	}
	return &TsTuple{N_TS_TUPLE, p.finLoc(loc), args}, nil
}

// returns `FunctionType` or `PrimaryType` since `FunctionType`
// is conflicts with `ParenthesizedType`
func (p *Parser) tsPrimary(rough bool) (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	av := ahead.value
	if av == T_PAREN_L {
		// paren type
		return p.tsParen()
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
		node, err = p.tsObj(rough)
		if err != nil {
			return nil, err
		}
	} else if av == T_BRACKET_L {
		// tuple type
		node, err = p.tsTuple(rough)
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

func (p *Parser) tsObj(rough bool) (Node, error) {
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
		prop, err := p.tsProp(rough)
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

func (p *Parser) tsNewSig(loc *Loc) (Node, error) {
	p.lexer.Next()
	typParams, _, err := p.tsTypParams()
	if err != nil {
		return nil, err
	}
	params, _, err := p.paramList()
	if err != nil {
		return nil, err
	}
	retTyp, _, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsNewSig{N_TS_NEW_SIG, p.finLoc(loc), typParams, params, retTyp}, nil
}

func (p *Parser) tsIdxSig(loc *Loc) (Node, error) {
	p.lexer.Next()
	id, err := p.ident(nil, false)
	if err != nil {
		return nil, err
	}
	if _, err = p.nextMustTok(T_COLON); err != nil {
		return nil, err
	}
	key, err := p.tsTyp(false)
	if err != nil {
		return nil, err
	}
	if _, err = p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}
	val, _, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsIdxSig{N_TS_IDX_SIG, p.finLoc(loc), id, key, val}, nil
}

func (p *Parser) tsProp(rough bool) (Node, error) {
	ahead := p.lexer.Peek()
	av := ahead.value
	loc := p.locFromTok(ahead)

	if av == T_LT {
		ps, _, err := p.tsTypParams()
		if err != nil {
			return nil, err
		}
		return p.tsCallSig(ps, loc)
	}
	if av == T_NEW {
		// ConstructSignature
		return p.tsNewSig(loc)
	}
	if !rough && av == T_BRACKET_L {
		// IndexSignature
		return p.tsIdxSig(loc)
	}
	if av == T_PAREN_L {
		// CallSignature
		return p.tsCallSig(nil, loc)
	}
	if rough && av == T_DOT_TRI {
		binding, err := p.patternRest(false, true)
		if err != nil {
			return nil, err
		}
		return binding, nil
	}

	if !rough {
		// PropertySignature or MethodSignature
		name, err := p.tsPropName()
		if err != nil {
			return nil, err
		}
		var ques *Loc
		if p.lexer.Peek().value == T_HOOK {
			ques = p.locFromTok(p.lexer.Next())
		}
		typAnnot, _, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}
		if typAnnot != nil {
			return &TsProp{N_TS_PROP, p.finLoc(name.Loc().Clone()), name, typAnnot, ques, nil}, nil
		}

		// MethodSignature is deserved
		callSig, err := p.tsCallSig(nil, nil)
		if err != nil {
			return nil, err
		}
		return &TsProp{N_TS_PROP, p.finLoc(name.Loc().Clone()), name, callSig, ques, nil}, nil
	}

	key, compute, err := p.propName(false, true)
	if err != nil {
		return nil, err
	}
	if key.Type() == N_PROP {
		return key, nil
	}

	var ques *Loc
	if p.lexer.Peek().value == T_HOOK {
		ques = p.locFromTok(p.lexer.Next())
	}

	tok := p.lexer.Peek()
	assign := tok.value == T_ASSIGN
	var value Node
	if tok.value == T_COLON {
		p.lexer.Next()
		value, err = p.tsTyp(rough)
		if err != nil {
			return nil, err
		}
	} else if assign {
		if key.Type() == N_NAME {
			return p.patternAssign(key, true)
		}
		return nil, p.errorTok(tok)
	}
	return &TsProp{N_TS_PROP, p.finLoc(key.Loc().Clone()), key, value, ques, compute}, nil
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
		cons, err = p.tsTyp(false)
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

func (p *Parser) tsTypArgs() ([]Node, error) {
	p.lexer.Next()
	args := make([]Node, 0, 1)
	for {
		arg, err := p.tsTyp(false)
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
	typAnnot, _, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsCallSig{N_TS_CALL_SIG, p.finLoc(loc), typParams, params, typAnnot}, nil
}

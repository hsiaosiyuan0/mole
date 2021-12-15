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
	if !p.ts() {
		return nil, nil, nil
	}
	ahead := p.lexer.Peek()
	if ahead.value == T_COLON {
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
			loc := p.locFromTok(p.lexer.Next())
			rhs, err = p.tsPrimary(rough)
			if err != nil {
				return nil, err
			}
			lhs = &TsUnionTyp{N_TS_UNION_TYP, p.finLoc(lhs.Loc()), lhs, loc, rhs}
		} else if av == T_BIT_AND {
			loc := p.locFromTok(p.lexer.Next())
			rhs, err = p.tsPrimary(rough)
			if err != nil {
				return nil, err
			}
			lhs = &TsIntersecTyp{N_TS_INTERSEC_TYP, p.finLoc(lhs.Loc()), lhs, loc, rhs}
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
	params, _, err := p.paramList(false)
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

func (p *Parser) tsIsPrimitive(typ NodeType) bool {
	switch typ {
	case N_TS_ANY, N_TS_NUM, N_TS_BOOL, N_TS_STR, N_TS_SYM, N_TS_VOID, N_TS_THIS:
		return true
	}
	return false
}

func (p *Parser) tsNodeTypAnnot(binding Node, typAnnot Node) bool {
	if typAnnot == nil {
		return true
	}

	switch binding.Type() {
	case N_NAME:
		binding.(*Ident).typAnnot = typAnnot
		return true
	case N_PAT_ARRAY:
		binding.(*ArrPat).typAnnot = typAnnot
		return true
	case N_PAT_OBJ:
		binding.(*ObjPat).typAnnot = typAnnot
		return true
	case N_LIT_ARR:
		binding.(*ArrLit).typAnnot = typAnnot
		return true
	case N_LIT_OBJ:
		binding.(*ObjLit).typAnnot = typAnnot
		return true
	case N_SPREAD:
		binding.(*Spread).typAnnot = typAnnot
		return true
	}
	return false
}

// `RoughParam` is something like `a:b` which `a` is a rough-type and `b` is typAnnot
// convert rough param to formal param needs to process `a` in above example - in other
// words convert ts-type-node to js-node
func (p *Parser) tsRoughParamToParam(node Node) (Node, error) {
	var err error
	n := node
	if node.Type() == N_TS_ROUGH_PARAM {
		param := node.(*TsRoughParam)
		if param.name.Type() == N_TS_THIS && param.typAnnot == nil {
			return nil, p.errorAtLoc(param.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		fp, err := p.tsRoughParamToParam(param.name)
		if err != nil {
			return nil, err
		}

		if ok := p.tsNodeTypAnnot(fp, param.typAnnot); !ok {
			return nil, p.errorAtLoc(fp.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		return fp, nil
	}

	switch n.Type() {
	case N_TS_ANY, N_TS_NUM, N_TS_BOOL, N_TS_STR, N_TS_SYM:
		d := n.(*TsPredef)
		return &Ident{N_NAME, d.loc, d.loc.Text(), false, false, nil, false, nil}, nil
	case N_TS_VOID:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	case N_TS_REF:
		r := n.(*TsRef)
		if r.HasArgs() {
			return nil, p.errorAtLoc(r.lt, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToParam(r.name)
	case N_TS_OBJ:
		o := n.(*TsObj)
		props := make([]Node, len(o.props))
		for i, pn := range o.props {
			props[i], err = p.tsRoughParamToParam(pn)
			if err != nil {
				return nil, err
			}
		}
		return &ObjPat{N_PAT_OBJ, o.loc, props, nil, nil}, nil
	case N_TS_PROP:
		pn := n.(*TsProp)
		if pn.computeLoc != nil {
			return nil, p.errorAtLoc(pn.computeLoc, ERR_UNEXPECTED_TOKEN)
		}
		if pn.ques != nil {
			return nil, p.errorAtLoc(pn.ques, ERR_UNEXPECTED_TOKEN)
		}
		var val Node
		if pn.val != nil {
			val, err = p.tsRoughParamToParam(pn.val)
			if err != nil {
				return nil, err
			}
		}
		return &Prop{N_PROP, pn.loc, pn.key, nil, val, false, false, val == nil, false, PK_INIT}, nil
	case N_TS_ARR:
		a := n.(*TsArr)
		if p.tsIsPrimitive(a.arg.Type()) {
			return nil, p.errorAtLoc(a.bracket, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToParam(a.arg)
	case N_TS_TUPLE:
		t := n.(*TsTuple)
		elems := make([]Node, len(t.args))
		for i, arg := range t.args {
			elems[i], err = p.tsRoughParamToParam(arg)
			if err != nil {
				return nil, err
			}
		}
		return &ArrPat{N_PAT_ARRAY, t.loc, elems, nil, nil}, nil
	case N_TS_PAREN:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	case N_TS_THIS:
		t := n.(*TsThis)
		return &Ident{N_NAME, t.loc, t.loc.Text(), false, false, nil, true, nil}, nil
	case N_TS_NS_NAME:
		s := n.(*TsNsName)
		return nil, p.errorAtLoc(s.dot, ERR_UNEXPECTED_TOKEN)
	case N_TS_CALL_SIG, N_TS_NEW_SIG, N_TS_FN_TYP:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	case N_TS_UNION_TYP:
		u := n.(*TsUnionTyp)
		if p.tsIsPrimitive(u.lhs.Type()) {
			return nil, p.errorAtLoc(u.op, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToParam(u.lhs)
	case N_TS_INTERSEC_TYP:
		i := n.(*TsIntersecTyp)
		if p.tsIsPrimitive(i.lhs.Type()) {
			return nil, p.errorAtLoc(i.op, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToParam(i.lhs)
	}
	return node, nil
}

func (p *Parser) tsPropToIdxSig(prop *TsProp) (Node, error) {
	if prop.key.Type() != N_NAME {
		return nil, p.errorAtLoc(prop.key.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	name := prop.key.(*Ident)
	if name.typAnnot == nil {
		return nil, p.errorAtLoc(name.loc, ERR_UNEXPECTED_TOKEN)
	}
	switch name.typAnnot.Type() {
	case N_TS_NUM, N_TS_STR, N_TS_SYM:
		break
	default:
		return nil, p.errorAtLoc(name.typAnnot.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	vt := prop.val.Type()
	if vt < N_TS_ANY || vt > N_TS_ROUGH_PARAM {
		return nil, p.errorAtLoc(prop.val.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	typAnnot, err := p.tsRoughParamToTyp(prop.val)
	if err != nil {
		return nil, err
	}
	return &TsIdxSig{N_TS_IDX_SIG, prop.loc, name, name.typAnnot, typAnnot}, nil
}

// `RoughParam` is something like `a:b` which `a` is a rough-type and `b` is typAnnot
// `rough-type` is a mixed node consists of ts-type-node and js-node, especially in tsObj
// and tsTuple
func (p *Parser) tsRoughParamToTyp(node Node) (Node, error) {
	var err error
	n := node
	if node.Type() == N_TS_ROUGH_PARAM {
		param := node.(*TsRoughParam)
		if param.colon != nil {
			return nil, p.errorAtLoc(param.colon, ERR_UNEXPECTED_TOKEN)
		}
		n = param.name
	}

	switch n.Type() {
	case N_TS_OBJ:
		obj := n.(*TsObj)
		for i, prop := range obj.props {
			obj.props[i], err = p.tsRoughParamToTyp(prop)
			if err != nil {
				return nil, err
			}
		}
		return obj, nil
	case N_TS_PROP:
		pn := n.(*TsProp)
		var prop Node
		var err error
		if pn.computeLoc != nil {
			prop, err = p.tsPropToIdxSig(pn)
			if err != nil {
				return nil, err
			}
		}
		return prop, nil
	case N_TS_TUPLE:
		arr := n.(*TsTuple)
		for i, arg := range arr.args {
			arr.args[i], err = p.tsRoughParamToTyp(arg)
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	case N_PAT_REST:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	return n, nil
}

// `ParenthesizedType` or`FunctionType`
func (p *Parser) tsParen(keepParen bool) (Node, error) {
	params, loc, err := p.paramList(true)
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
		typ, err := p.tsFnTyp(params, loc)
		if err != nil {
			return nil, err
		}
		if !keepParen {
			return typ, nil
		}
		return &TsParen{N_TS_PARAM, p.finLoc(loc), typ}, nil
	}

	if len(params) == 1 {
		param := params[0].(*TsRoughParam)
		if param.colon != nil {
			return nil, p.errorAtLoc(param.colon, ERR_UNEXPECTED_TOKEN)
		}
		typ, err := p.tsRoughParamToTyp(param)
		if err != nil {
			return nil, err
		}
		if !keepParen {
			return typ, nil
		}
		return &TsParen{N_TS_PARAM, p.finLoc(loc), typ}, nil
	}

	if len(params) == 0 {
		typ := &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), nil, params, nil}
		if !keepParen {
			return typ, nil
		}
		return &TsParen{N_TS_PARAM, p.finLoc(loc), typ}, nil
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

		params, _, err = p.paramList(false)
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
			loc := p.locFromTok(p.lexer.Next())
			id, err := p.ident(nil, false)
			if err != nil {
				return nil, err
			}
			ns = &TsNsName{N_TS_NS_NAME, p.finLoc(ns.Loc().Clone()), ns, loc, id}
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
	args, lt, err := p.tsTypArgs()
	if err != nil {
		return nil, err
	}
	return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, lt, args}, nil
}

func (p *Parser) tsArr(typ Node) (Node, error) {
	for {
		if p.lexer.Peek().value == T_BRACKET_L {
			loc := p.locFromTok(p.lexer.Next())
			if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
				return nil, err
			}
			typ = &TsArr{N_TS_ARR, p.finLoc(typ.Loc().Clone()), loc, typ}
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
		return p.tsParen(rough)
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
	params, _, err := p.paramList(false)
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

	key, compute, err := p.propName(false, true, true)
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

func (p *Parser) tsTypArgs() ([]Node, *Loc, error) {
	loc := p.locFromTok(p.lexer.Next())
	args := make([]Node, 0, 1)
	for {
		arg, err := p.tsTyp(false)
		if err != nil {
			return nil, nil, err
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
	return args, loc, nil
}

func (p *Parser) tsCallSig(typParams []Node, loc *Loc) (Node, error) {
	if typParams == nil && loc == nil {
		var err error
		typParams, loc, err = p.tsTypParams()
		if err != nil {
			return nil, err
		}
	}
	params, _, err := p.paramList(false)
	if err != nil {
		return nil, err
	}
	typAnnot, _, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsCallSig{N_TS_CALL_SIG, p.finLoc(loc), typParams, params, typAnnot}, nil
}

package parser

import "fmt"

var builtinTyp = map[string]NodeType{
	"any":     N_TS_ANY,
	"number":  N_TS_NUM,
	"boolean": N_TS_BOOL,
	"string":  N_TS_STR,
	"symbol":  N_TS_SYM,
	"void":    N_TS_VOID,
	"never":   N_TS_NEVER,
	"unknown": N_TS_UNKNOWN,
}

func (p *Parser) ts() bool {
	return p.feat&FEAT_TS != 0
}

func (p *Parser) newTypInfo() *TypInfo {
	if p.ts() {
		return &TypInfo{ACC_MOD_PUB, nil, nil, nil, nil}
	}
	return nil
}

func (p *Parser) tsTypAnnot() (Node, error) {
	if !p.ts() {
		return nil, nil
	}
	ahead := p.lexer.Peek()
	if ahead.value == T_COLON {
		loc := p.locFromTok(p.lexer.Next())
		node, err := p.tsTyp(false)
		if err != nil {
			return nil, err
		}
		return &TsTypAnnot{N_TS_TYP_ANNOT, p.finLoc(loc), node}, nil
	}
	return nil, nil
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
			lhs = &TsUnionTyp{N_TS_UNION_TYP, p.finLoc(lhs.Loc().Clone()), lhs, loc, rhs}
		} else if av == T_BIT_AND {
			loc := p.locFromTok(p.lexer.Next())
			rhs, err = p.tsPrimary(rough)
			if err != nil {
				return nil, err
			}
			lhs = &TsIntersecTyp{N_TS_INTERSEC_TYP, p.finLoc(lhs.Loc().Clone()), lhs, loc, rhs}
		} else {
			break
		}
	}
	return lhs, nil
}

func (p *Parser) tsConstructTyp() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	params, typParams, _, err := p.paramList(false, false, true)
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

func (p *Parser) tsExprHasTypAnnot(node Node) bool {
	switch node.Type() {
	case N_LIT_ARR, N_LIT_OBJ, N_EXPR_THIS, N_NAME, N_PAT_REST, N_PAT_ARRAY, N_PAT_ASSIGN, N_PAT_OBJ:
		return true
	}
	return false
}

func (p *Parser) tsQues() *Loc {
	var ques *Loc
	if p.lexer.Peek().value == T_HOOK {
		ques = p.locFromTok(p.lexer.Next())
	}
	return ques
}

func (p *Parser) tsNodeTypAnnot(binding Node, typAnnot Node, accMod ACC_MOD, ques *Loc) bool {
	if wt, ok := binding.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti == nil {
			return false
		}

		if typAnnot != nil {
			ti.typAnnot = typAnnot
		}
		if ques != nil {
			ti.ques = ques
		}
		ti.accMod = accMod
		return true
	}
	return false
}

// `TypeAssertionExpression`
func (p *Parser) tsTypAssert(node Node, typArgs Node) (Node, error) {
	if typArgs == nil {
		return node, nil
	}

	args := typArgs.(*TsParamsInst).params
	if len(args) == 0 {
		return node, nil
	}

	if !p.ts() {
		return nil, p.errorAtLoc(typArgs.Loc(), ERR_UNEXPECTED_TOKEN)
	}

	if len(args) > 1 {
		return nil, p.errorAtLoc(args[1].Loc(), ERR_UNEXPECTED_TOKEN)
	}

	return &TsTypAssert{N_TS_TYP_ASSERT, p.finLoc(typArgs.Loc().Clone()), args[0], node}, nil
}

func (p *Parser) tsAdvanceHook() *Loc {
	var loc *Loc
	if p.lexer.Peek().value == T_HOOK {
		loc = p.locFromTok(p.lexer.Next())
	}
	return loc
}

// `RoughParam` is something like `a:b` which `a` is a rough-type and `b` is typAnnot
// convert rough param to formal param needs to process `a` in above example - in other
// words convert ts-type-node to js-node
func (p *Parser) tsRoughParamToParam(node Node) (Node, error) {
	var err error
	n := node
	if node.Type() == N_TS_ROUGH_PARAM {
		param := node.(*TsRoughParam)
		if param.name.Type() == N_TS_THIS && param.ti.typAnnot == nil {
			return nil, p.errorAtLoc(param.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		fp, err := p.tsRoughParamToParam(param.name)
		if err != nil {
			return nil, err
		}

		ti := param.ti
		if ok := p.tsNodeTypAnnot(fp, ti.typAnnot, ti.accMod, ti.ques); !ok {
			return nil, p.errorAtLoc(fp.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		return fp, nil
	}

	switch n.Type() {
	case N_TS_ANY, N_TS_NUM, N_TS_BOOL, N_TS_STR, N_TS_SYM:
		d := n.(*TsPredef)
		ti := p.newTypInfo()
		ti.ques = d.ques
		return &Ident{N_NAME, d.loc, d.loc.Text(), false, false, nil, false, ti}, nil
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
		return &ObjPat{N_PAT_OBJ, o.loc, props, nil, p.newTypInfo()}, nil
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
		return &Prop{N_PROP, pn.loc, pn.key, nil, val, false, false, val == nil, false, PK_INIT, ACC_MOD_NONE}, nil
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
		return &ArrPat{N_PAT_ARRAY, t.loc, elems, nil, p.newTypInfo()}, nil
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
	if name.ti.typAnnot == nil {
		return nil, p.errorAtLoc(name.loc, ERR_UNEXPECTED_TOKEN)
	}
	switch name.ti.typAnnot.(*TsTypAnnot).tsTyp.Type() {
	case N_TS_NUM, N_TS_STR, N_TS_SYM:
		break
	default:
		return nil, p.errorAtLoc(name.ti.typAnnot.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	vt := prop.val.Type()
	if vt < N_TS_ANY || vt > N_TS_ROUGH_PARAM {
		return nil, p.errorAtLoc(prop.val.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	typAnnot, err := p.tsRoughParamToTyp(prop.val, false)
	if err != nil {
		return nil, err
	}
	return &TsIdxSig{N_TS_IDX_SIG, prop.loc, name, name.ti.typAnnot, typAnnot}, nil
}

// `RoughParam` is something like `a:b` which `a` is a rough-type and `b` is typAnnot
// `rough-type` is a mixed node consists of ts-type-node and js-node, especially in tsObj
// and tsTuple
func (p *Parser) tsRoughParamToTyp(node Node, raise bool) (Node, error) {
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
	case N_NAME:
		n := node.(*Ident)
		name := n.Text()
		if typ, ok := builtinTyp[name]; ok {
			// predef
			return &TsPredef{typ, n.loc, nil}, nil
		}
		return &TsRef{N_TS_REF, n.Loc().Clone(), n, nil, nil}, nil
	case N_TS_OBJ:
		obj := n.(*TsObj)
		for i, prop := range obj.props {
			obj.props[i], err = p.tsRoughParamToTyp(prop, raise)
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
			return prop, nil
		}
		return pn, nil
	case N_TS_TUPLE:
		arr := n.(*TsTuple)
		for i, arg := range arr.args {
			arr.args[i], err = p.tsRoughParamToTyp(arg, raise)
			if err != nil {
				return nil, err
			}
		}
		return arr, nil
	case N_PAT_REST:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	}

	if raise {
		return nil, p.errorAtLoc(node.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	return n, nil
}

// `ParenthesizedType` or`FunctionType`
func (p *Parser) tsParen(keepParen bool) (Node, error) {
	params, _, loc, err := p.paramList(true, false, false)
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
		typ, err := p.tsRoughParamToTyp(param, false)
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
	var typParams Node
	var err error
	var loc *Loc
	if parenL != nil {
		loc = parenL
	}
	if params == nil {
		params, typParams, _, err = p.paramList(false, false, true)
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

// `typePredicates` and `assertPredicate`
func (p *Parser) tsTypPredicate(name Node, asserts bool) (Node, error) {
	loc := name.Loc().Clone()
	var err error
	if asserts {
		name, err = p.ident(nil, true)
		if err != nil {
			return nil, err
		}
	}

	ahead := p.lexer.Peek()
	av := ahead.value

	var typ Node
	if av == T_NAME && ahead.Text() == "is" {
		p.lexer.Next()

		typ, err = p.tsTyp(false)
		if err != nil {
			return nil, err
		}
		typ = &TsTypAnnot{N_TS_TYP_ANNOT, typ.Loc().Clone(), typ}
	}

	return &TsTypPredicate{N_TS_TYP_PREDICATE, p.finLoc(loc), name, typ, asserts}, nil
}

func (p *Parser) tsRef(ns Node) (Node, error) {
	name, err := p.tsTypName(ns)
	if err != nil {
		return nil, err
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	asserts := p.isName(name, "asserts", false, false)
	// type predicates
	if (asserts && av == T_NAME) || ahead.Text() == "is" {
		return p.tsTypPredicate(name, asserts)
	}

	if av != T_LT {
		return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, nil, nil}, nil
	}
	args, err := p.tsTypArgs()
	if err != nil {
		return nil, err
	}
	return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, args.Loc().Clone(), args}, nil
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

func (p *Parser) tsPredefOrRef(tok *Token) (Node, error) {
	if tok == nil {
		tok = p.lexer.Peek()
	}
	tv := tok.value
	name := "void"

	var node Node
	var err error
	var loc *Loc
	if tv == T_NAME {
		node, err = p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		name = node.(*Ident).Text()
		loc = node.Loc()
	} else {
		// void
		loc = p.locFromTok(p.lexer.Next())
	}

	if typ, ok := builtinTyp[name]; ok {
		// predef
		return &TsPredef{typ, loc, nil}, nil
	}

	return p.tsRef(node)
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
	if av == T_NAME || av == T_VOID {
		node, err = p.tsPredefOrRef(ahead)
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
	params, typParams, _, err := p.paramList(false, false, true)
	if err != nil {
		return nil, err
	}
	retTyp, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	return &TsNewSig{N_TS_NEW_SIG, p.finLoc(loc), typParams, params, retTyp}, nil
}

func (p *Parser) tsIdxSig(loc *Loc) (Node, error) {
	tok := p.lexer.Next()
	if loc == nil {
		loc = p.locFromTok(tok)
	}
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
	val, err := p.tsTypAnnot()
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
		ps, err := p.tsTypParams()
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
		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}

		av := p.lexer.Peek().value
		if av != T_PAREN_L && av != T_LT {
			return &TsProp{N_TS_PROP, p.finLoc(name.Loc().Clone()), name, typAnnot, ques, nil}, nil
		}

		// MethodSignature is deserved
		callSig, err := p.tsCallSig(nil, p.locFromTok(p.lexer.Peek()))
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
		return &StrLit{N_LIT_STR, p.finLoc(loc), tok.Text(), legacyOctalEscapeSeq, nil, nil}, nil
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
	return &TsParam{N_TS_PARAM, p.finLoc(id.loc.Clone()), id, cons, nil}, nil
}

func (p *Parser) tsTryTypParams() (Node, error) {
	if !p.ts() || p.lexer.Peek().value != T_LT {
		return nil, nil
	}
	return p.tsTypParams()
}

func (p *Parser) tsTypParams() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	ps := make([]Node, 0, 1)
	for {
		pa, err := p.tsTypParam()
		if err != nil {
			return nil, err
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
	return &TsParamsDec{N_TS_PARAM_DEC, p.finLoc(loc), ps}, nil
}

// returned nodes are the superset of typeArgs which also includes `Constraint`
//
// the caller maybe required to check the returned nodes are:
// - valid type args, by further doing a `tsCheckTypArgs` subroutine
// - valid type params, by further doing a `tsTypArgsToTypParams` subroutine
//
// this method returns type params instead of type args since the former is
// a superset of the later and in the calling point of this method, there is
// no enough information to determine whether the type params or args is satisfied
// so like others method to resolve the ambiguities in the grammar - a rough firstly
// parsing is introduced and construct rough nodes to let the later processes to
// do checking or transformation to produce more precise results
//
// `async` will bypassed to method `tsTryTypArgsAfterAsync` to handle the async ambiguity:
//
// ```
// async < a, b;
// async<T>() == 0;
// ```
func (p *Parser) tsTryTypArgs(asyncLoc *Loc) (Node, error) {
	if !p.ts() || p.lexer.Peek().value != T_LT {
		return nil, nil
	}
	if asyncLoc != nil {
		return p.tsTryTypArgsAfterAsync(asyncLoc)
	}
	return p.tsTypArgs()
}

// consider ambiguity of the first `<` in below example:
//
// ```
// async < a, b;
// async<T>() == 0;
// ```
//
// for avoiding lookbehind the process should accept the input as seqExpr then try to
// tansform the subtree of seqExpr to typArgs if its followed by `>`
func (p *Parser) tsTryTypArgsAfterAsync(asyncLoc *Loc) (Node, error) {
	name := &Ident{N_NAME, asyncLoc, asyncLoc.Text(), false, false, nil, true, p.newTypInfo()}
	binExpr, err := p.binExpr(name, 0, false, false, true)
	ltLoc := binExpr.(*BinExpr).opLoc.Clone()
	if err != nil {
		return nil, err
	}
	seq, err := p.seqExpr(binExpr, true)
	if err != nil {
		return nil, err
	}
	if p.lexer.Peek().value == T_GT {
		// to type args
		return p.tsSeqExprToTypArgs(seq, ltLoc)
	}
	return seq, nil
}

func (p *Parser) tsSeqExprToTypArgs(node Node, loc *Loc) (Node, error) {
	p.lexer.Next() // `>`
	p.finLoc(loc)

	var nodes []Node
	nt := node.Type()
	if nt == N_EXPR_BIN {
		nodes = []Node{node.(*BinExpr).rhs}
	} else if nt == N_EXPR_SEQ {
		els := node.(*SeqExpr).elems
		nodes = []Node{els[0].(*BinExpr).rhs}
		if len(els) > 1 {
			nodes = append(nodes, els[1:]...)
		}
	}

	var err error
	args := make([]Node, len(nodes))
	for i, n := range nodes {
		args[i], err = p.tsRoughParamToTyp(n, true)
		if err != nil {
			return nil, err
		}
	}
	return &TsParamsInst{N_TS_PARAM_INST, loc, args}, nil
}

func (p *Parser) tsCheckTypArgs(node Node) error {
	if node == nil {
		return nil
	}

	nodes := node.(*TsParamsInst).params
	for i, arg := range nodes {
		if arg.Type() == N_TS_PARAM {
			pn := arg.(*TsParam)
			if pn.cons != nil {
				return p.errorAtLoc(arg.(*TsParam).cons.Loc(), ERR_UNEXPECTED_TOKEN)
			}
			nodes[i] = pn.name
		}
	}
	return nil
}

func (p *Parser) tsTypArgsToTypParams(node Node) (Node, error) {
	if node == nil {
		return nil, nil
	}
	nodes := node.(*TsParamsInst).params

	var err error
	for i, n := range nodes {
		n, err = p.tsRoughParamToParam(n)
		if n.Type() == N_NAME {
			n = &TsParam{N_TS_PARAM, n.Loc().Clone(), n, nil, nil}
		}
		nodes[i] = n
		if err != nil {
			return nil, err
		}
	}
	return &TsParamsDec{N_TS_PARAM_DEC, node.Loc(), nodes}, nil
}

func (p *Parser) tsPredefToName(node Node) (Node, error) {
	node, err := p.tsRoughParamToParam(node)
	if err != nil {
		return nil, err
	}
	if node.Type() != N_NAME {
		return nil, p.errorAtLoc(node.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	return node, nil
}

func (p *Parser) tsTypArgs() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	args := make([]Node, 0, 1)
	for {
		arg, err := p.tsTyp(false)
		if p.lexer.Peek().value == T_EXTENDS {
			id, err := p.tsPredefToName(arg)
			if err != nil {
				return nil, err
			}

			p.lexer.Next()
			cons, err := p.tsTyp(false)
			if err != nil {
				return nil, err
			}
			arg = &TsParam{N_TS_PARAM, p.finLoc(id.Loc().Clone()), id, cons, nil}
		}
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
	return &TsParamsInst{N_TS_PARAM_DEC, p.finLoc(loc), args}, nil
}

func (p *Parser) tsCallSig(typParams Node, loc *Loc) (Node, error) {
	params, tp, _, err := p.paramList(false, false, typParams == nil && loc == nil)
	if err != nil {
		return nil, err
	}
	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	if tp != nil {
		typParams = tp
	}
	return &TsCallSig{N_TS_CALL_SIG, p.finLoc(loc), typParams, params, typAnnot}, nil
}

func (p *Parser) aheadIsTsTypDec(tok *Token) bool {
	if p.ts() && tok.value == T_NAME {
		return tok.Text() == "type"
	}
	return false
}

func (p *Parser) tsTypDec() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	name, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}
	params, err := p.tsTryTypParams()
	if err != nil {
		return nil, err
	}

	if _, err = p.nextMustTok(T_ASSIGN); err != nil {
		return nil, err
	}

	typAnnot, err := p.tsTyp(false)
	if err != nil {
		return nil, err
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	ti := p.newTypInfo()
	ti.typParams = params
	ti.typAnnot = typAnnot
	return &TsTypDec{N_TS_TYP_DEC, p.finLoc(loc), name, ti}, nil
}

func (p *Parser) tsIsFnSigValid(name string) error {
	if p.lastTsFnSig == nil {
		return nil
	}

	if p.lastTsFnSig.id.(*Ident).Text() == name {
		return nil
	}
	return p.errorAtLoc(p.lastTsFnSig.loc, ERR_FN_SIG_MISSING_IMPL)
}

func (p *Parser) tsIsFnImplValid(id Node) error {
	if p.lastTsFnSig == nil {
		return nil
	}

	ep := p.lastTsFnSig.id.(*Ident).Text()
	act := id.(*Ident).Text()
	if ep == act {
		return nil
	}
	return p.errorAtLoc(id.(*Ident).loc, fmt.Sprintf(ERR_TPL_INVALID_FN_IMPL_NAME, ep))
}

func (p *Parser) aheadIsTsItf(tok *Token) bool {
	return p.ts() && tok.value == T_INTERFACE
}

func (p *Parser) tsItfExtClause() ([]Node, error) {
	if p.lexer.Peek().value == T_EXTENDS {
		p.lexer.Next()
	}
	ns := make([]Node, 0, 1)
	for {
		tr, err := p.tsPredefOrRef(nil)
		if err != nil {
			return nil, err
		}

		tt := tr.Type()
		if tt != N_TS_REF {
			if tt >= N_TS_ANY && tt <= N_TS_SYM {
				return nil, p.errorAtLoc(tr.Loc(), fmt.Sprintf(ERR_TPL_USE_TYP_AS_VALUE, tr.Loc().Text()))
			}
			return nil, p.errorAtLoc(tr.Loc(), ERR_UNEXPECTED_TOKEN)
		}
		ns = append(ns, tr)

		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACE_L {
			break
		} else if av == T_COMMA {
			p.lexer.Next()
		}
	}
	return ns, nil
}

func (p *Parser) tsItf() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	name, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}

	params, err := p.tsTryTypParams()
	if err != nil {
		return nil, err
	}

	supers, err := p.tsItfExtClause()
	if err != nil {
		return nil, err
	}

	body, err := p.tsObj(false)
	if err != nil {
		return nil, err
	}
	return &TsInferface{N_TS_INTERFACE, p.finLoc(loc), name, params, supers, body}, nil
}

func (p *Parser) tsEnumBody() ([]Node, error) {
	if _, err := p.nextMustTok(T_BRACE_L); err != nil {
		return nil, err
	}
	mems := make([]Node, 0, 1)
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACE_R {
			p.lexer.Next()
			break
		}
		name, err := p.tsPropName()
		if err != nil {
			return nil, err
		}
		var val Node
		if p.lexer.Peek().value == T_ASSIGN {
			p.lexer.Next()
			val, err = p.assignExpr(false, false, false)
			if err != nil {
				return nil, err
			}
		}
		mems = append(mems, &TsEnumMember{N_TS_ENUM_MEMBER, p.finLoc(name.Loc().Clone()), name, val})
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
	}
	return mems, nil
}

func (p *Parser) aheadIsTsEnum(tok *Token) bool {
	if !p.ts() {
		return false
	}
	if tok == nil {
		tok = p.lexer.Peek()
	}
	return tok.value == T_ENUM
}

// `loc` is the loc of the preceding `const`
func (p *Parser) tsEnum(loc *Loc) (Node, error) {
	cons := loc != nil
	tok := p.lexer.Next() // enum
	if loc == nil {
		loc = p.locFromTok(tok)
	}

	name, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}

	mems, err := p.tsEnumBody()
	if err != nil {
		return nil, err
	}
	return &TsEnum{N_TS_ENUM, p.finLoc(loc), name, mems, cons}, nil
}

// `ImportAliasDeclaration` or `ImportRequireDeclaration`
func (p *Parser) tsImportAlias(loc *Loc, name Node, export bool) (Node, error) {
	p.lexer.Next() // `=`

	val, err := p.tsTypName(nil)
	if err != nil {
		return nil, err
	}

	var node Node
	if val.Type() == N_NAME && val.(*Ident).Text() == "require" {
		call, _, err := p.callExpr(val, true, false, nil)
		if err != nil {
			return nil, err
		}
		node = &TsImportRequire{N_TS_IMPORT_REQUIRE, p.finLoc(loc), name, call.(*CallExpr).args}
	} else {
		node = &TsImportAlias{N_TS_IMPORT_ALIAS, p.finLoc(loc), name, val, export}
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) aheadIsTsNS(tok *Token) bool {
	return p.ts() && tok.value == T_NAME && tok.Text() == "namespace"
}

func (p *Parser) tsNS() (Node, error) {
	loc := p.locFromTok(p.lexer.Next()) // `namespace`

	name, err := p.tsTypName(nil)
	if err != nil {
		return nil, err
	}

	blk, err := p.blockStmt(true)
	if err != nil {
		return nil, err
	}

	return &TsNS{N_TS_IMPORT_ALIAS, p.finLoc(loc), name, blk.body}, nil
}

func (p *Parser) aheadIsTsDec(tok *Token) bool {
	return p.ts() && tok.value == T_NAME && tok.Text() == "declare"
}

func (p *Parser) aheadIsModDec(tok *Token) bool {
	return p.ts() && tok.value == T_NAME && tok.Text() == "module"
}

func (p *Parser) tsModDec() (*TsDec, error) {
	p.lexer.Next() // `module`

	tok := p.lexer.Next()
	if tok.value != T_STRING {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_UNEXPECTED_TOKEN)
	}

	name := &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(tok)), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil, nil}

	blk, err := p.blockStmt(true)
	if err != nil {
		return nil, err
	}
	return &TsDec{N_TS_DEC_MODULE, p.finLoc(name.loc.Clone()), name, blk}, nil
}

func (p *Parser) tsDec() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())

	tok := p.lexer.Peek()
	tv := tok.value

	var err error
	typ := N_ILLEGAL
	dec := &TsDec{typ, nil, nil, nil}
	if ok, kind := p.aheadIsVarDec(tok); ok {
		dec.inner, err = p.varDecStmt(kind, false)
		typ = N_TS_DEC_VAR_DEC
	} else if tv == T_FUNC {
		dec.inner, err = p.fnDec(false, nil, false, true)
		typ = N_TS_DEC_FN
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tv, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return nil, p.errorAt(tv, &tok.begin, ERR_ASYNC_IN_AMBIENT)
	}

	if tv == T_CLASS {
		dec.inner, err = p.classDec(false, false, true)
		typ = N_TS_DEC_CLASS
	} else if p.aheadIsTsItf(tok) {
		dec.inner, err = p.tsItf()
		typ = N_TS_DEC_INTERFACE
	} else if p.aheadIsTsTypDec(tok) {
		dec.inner, err = p.tsTypDec()
		typ = N_TS_DEC_TYP_DEC
	} else if p.aheadIsTsEnum(tok) {
		dec.inner, err = p.tsEnum(nil)
		typ = N_TS_DEC_ENUM
	} else if p.aheadIsTsNS(tok) {
		dec.inner, err = p.tsNS()
		typ = N_TS_DEC_NS
	} else if p.aheadIsModDec(tok) {
		dec, err = p.tsModDec()
		typ = N_TS_DEC_MODULE
	}

	if err != nil {
		return nil, err
	}

	if err = p.checkAmbient(typ, dec.inner); err != nil {
		return nil, err
	}

	dec.typ = typ
	dec.loc = p.finLoc(loc)

	return dec, nil
}

func (p *Parser) checkAmbient(typ NodeType, dec Node) error {
	switch typ {
	case N_TS_DEC_VAR_DEC:
		n := dec.(*VarDecStmt)
		for _, v := range n.decList {
			init := v.(*VarDec).init
			if init != nil {
				return p.errorAtLoc(init.Loc(), ERR_INIT_NOT_ALLOWED)
			}
		}
	}
	return nil
}

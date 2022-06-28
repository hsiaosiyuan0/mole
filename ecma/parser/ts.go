package parser

import (
	"errors"
	"fmt"
)

var builtinTyp = map[string]NodeType{
	"any":       N_TS_ANY,
	"number":    N_TS_NUM,
	"boolean":   N_TS_BOOL,
	"string":    N_TS_STR,
	"symbol":    N_TS_SYM,
	"object":    N_TS_OBJ,
	"void":      N_TS_VOID,
	"never":     N_TS_NEVER,
	"unknown":   N_TS_UNKNOWN,
	"undefined": N_TS_UNDEF,
	"null":      N_TS_NULL,
	"bigint":    N_TS_BIGINT,
	"intrinsic": N_TS_INTRINSIC,
}

// indicates the closing `>` is missing, so the processed `<` should be considered
// as the LessThan operator. produced in `tsTypArgs`
type ErrTypArgMissingGT struct {
	line uint32
	col  uint32
}

func (e *ErrTypArgMissingGT) Error() string {
	return "missing the closing `>`"
}

var errTypArgMissingGT = &ErrTypArgMissingGT{}

// indicates the current position should be re-entered as `jsx`. produced in `tsTypArgs`
var errTypArgMaybeJsx = errors.New("maybe jsx")

func (p *Parser) newTypInfo() *TypInfo {
	if p.ts || p.feat&FEAT_DECORATOR != 0 {
		return NewTypInfo()
	}
	return nil
}

func (p *Parser) tsTypAnnot() (Node, error) {
	if !p.ts {
		return nil, nil
	}
	ahead := p.lexer.Peek()
	av := ahead.value
	if av == T_COLON {
		loc := p.locFromTok(p.lexer.Next())
		node, err := p.tsTyp(false, false, true)
		if err != nil {
			return nil, err
		}
		return &TsTypAnnot{N_TS_TYP_ANNOT, p.finLoc(loc), node}, nil
	}
	return nil, nil
}

// for dealing with the ambiguous between `ParenthesizedType` and the `formalParamList`, eg:
//
// ```ts
// var a: ({ a = c }: string | number) => number = 1
// var a: (string | number) = 1
// ```
//
// the ambiguous in above code is that the `(` can be either the leading token of `formalParamList`
// or `ParenthesizedType`
//
// the manner used here is by setting the `rough` argument to `true` to parse `{a = c, b?: number}`
// as a superset consists of `tsTyp` and `bindingPattern`, the parsed result will be judged by later
// process by the time via method such as `tsRoughParamToParam`
func (p *Parser) tsTyp(rough bool, canConst bool, canCond bool) (Node, error) {
	if p.lexer.Peek().value == T_NEW {
		return p.tsConstructTyp(nil, false)
	}
	return p.tsUnionOrIntersectType(nil, 0, rough, canConst, canCond)
}

func (p *Parser) tsUnionOrIntersectType(lhs Node, minPcd int, rough bool, canConst bool, canCond bool) (Node, error) {
	var err error
	if lhs == nil {
		ahead := p.lexer.Peek()
		av := ahead.value

		// `type K = | number | string`
		if av != T_BIT_OR && av != T_BIT_AND {
			lhs, err = p.tsPrimary(rough, canConst, canCond)
			if err != nil {
				return nil, err
			}
		}
	}

	if lhs != nil && lhs.Type() == N_TS_FN_TYP {
		return lhs, nil
	}

	// only `intrinsic` at the first place of the rvalue of the type alias can
	// be parsed as `intrinsic` keywords or else return as ident
	p.scope().EraseKind(SPK_TS_MAY_INTRINSIC)

	var rhs Node
	var elems []Node
	var firstOp *Loc
	var nt NodeType
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av != T_BIT_OR && av != T_BIT_AND {
			break
		}

		kind := TokenKinds[av]
		pcd := kind.Pcd
		if pcd < minPcd {
			break
		}

		if firstOp == nil {
			firstOp = p.locFromTok(p.lexer.Next())
		} else {
			p.lexer.Next()
		}

		if nt == N_ILLEGAL {
			if av == T_BIT_OR {
				nt = N_TS_UNION_TYP
			} else {
				nt = N_TS_INTERSECT_TYP
			}
			elems = make([]Node, 0, 1)
			if lhs != nil {
				elems = append(elems, lhs)
			}
		}

		rhs, err = p.tsPrimary(rough, canConst, canCond)
		if err != nil {
			return nil, err
		}

		ahead = p.lexer.Peek()
		av = ahead.value
		kind = TokenKinds[av]
		for (av == T_BIT_OR || av == T_BIT_AND) && kind.Pcd > pcd {
			pcd = kind.Pcd
			rhs, err = p.tsUnionOrIntersectType(rhs, pcd, rough, canConst, canCond)
			if err != nil {
				return nil, err
			}
			ahead = p.lexer.Peek()
			av = ahead.value
			kind = TokenKinds[av]
		}

		elems = append(elems, rhs)
		if lhs != nil && kind.Pcd < pcd {
			if nt == N_TS_UNION_TYP {
				lhs = &TsUnionTyp{N_TS_UNION_TYP, p.finLoc(lhs.Loc().Clone()), firstOp, elems, nil}
			} else {
				lhs = &TsIntersectTyp{N_TS_INTERSECT_TYP, p.finLoc(lhs.Loc().Clone()), firstOp, elems, nil}
			}
			nt = N_ILLEGAL
		}
	}

	return lhs, nil
}

func (p *Parser) tsConstructTyp(loc *Loc, abstract bool) (Node, error) {
	tok := p.lexer.Next()
	if loc == nil {
		loc = p.locFromTok(tok)
	}
	params, typParams, _, err := p.paramList(false, PK_NONE, true)
	if err != nil {
		return nil, err
	}

	arrow, err := p.nextMustTok(T_ARROW)
	if err != nil {
		return nil, err
	}
	tiLoc := p.locFromTok(arrow)

	retTyp, err := p.tsTyp(false, false, false)
	if err != nil {
		return nil, err
	}
	ti := &TsTypAnnot{N_TS_TYP_ANNOT, p.finLoc(tiLoc), retTyp}
	return &TsNewSig{N_TS_NEW, p.finLoc(loc), typParams, params, ti, abstract}, nil
}

func (p *Parser) tsIsPrimitive(typ NodeType) bool {
	switch typ {
	case N_TS_ANY, N_TS_NUM, N_TS_BOOL, N_TS_STR, N_TS_SYM, N_TS_VOID, N_TS_THIS:
		return true
	}
	return false
}

func (p *Parser) tsIsNameLike(typ NodeType) bool {
	switch typ {
	case N_TS_ANY, N_TS_BOOL, N_TS_SYM, N_TS_VOID, N_TS_THIS, N_TS_NEVER, N_TS_UNKNOWN, N_TS_UNDEF, N_TS_REF:
		return true
	}
	return false
}

func (p *Parser) tsExprHasTypAnnot(node Node) bool {
	if !p.ts {
		return false
	}
	switch node.Type() {
	case N_LIT_ARR, N_LIT_OBJ, N_EXPR_THIS, N_NAME, N_PAT_REST, N_PAT_ARRAY, N_PAT_ASSIGN, N_PAT_OBJ:
		return true
	}
	return false
}

func (p *Parser) tsAdvanceHook(ep bool) (ques, not *Loc) {
	ahead := p.lexer.Peek()
	av := ahead.value
	if av == T_HOOK {
		ques = p.locFromTok(p.lexer.Next())
	} else if ep && av == T_NOT {
		not = p.locFromTok(p.lexer.Next())
	}
	return
}

func (p *Parser) tsNodeTypAnnot(binding Node, typAnnot Node, accMod ACC_MOD,
	beginLoc *Loc, abstract, readonly, override, declare bool, ques *Loc) bool {
	if wt, ok := binding.(NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti == nil {
			return false
		}

		if typAnnot != nil {
			ti.SetTypAnnot(typAnnot)
		}
		if ques != nil {
			ti.SetQues(ques)
		}
		ti.SetAccMod(accMod)
		ti.SetBeginLoc(beginLoc)
		ti.SetAbstract(abstract)
		ti.SetReadonly(readonly)
		ti.SetOverride(override)
		ti.SetDeclare(declare)
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

	if !p.ts {
		return nil, p.errorAtLoc(typArgs.Loc(), ERR_UNEXPECTED_TOKEN)
	}

	if len(args) > 1 {
		return nil, p.errorAtLoc(args[1].Loc(), ERR_UNEXPECTED_TOKEN)
	}

	return &TsTypAssert{N_TS_TYP_ASSERT, p.finLoc(typArgs.Loc().Clone()), args[0], node, nil}, nil
}

// `RoughParam` is something like `a:b` which `a` is a rough-type and `b` is typAnnot
// convert rough param to formal param needs to process `a` in above example - in other
// words convert ts-type-node to js-node
func (p *Parser) tsRoughParamToParam(node Node) (Node, error) {
	var err error
	n := node
	if node.Type() == N_TS_ROUGH_PARAM {
		param := node.(*TsRoughParam)
		if param.name.Type() == N_TS_THIS && param.ti.TypAnnot() == nil {
			return nil, p.errorAtLoc(param.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		fp, err := p.tsRoughParamToParam(param.name)
		if err != nil {
			return nil, err
		}

		ti := param.ti
		if ok := p.tsNodeTypAnnot(fp, ti.TypAnnot(), ti.AccMod(), ti.BeginLoc(),
			ti.Abstract(), ti.Readonly(), ti.Override(), ti.Declare(), ti.Ques()); !ok {
			return nil, p.errorAtLoc(fp.Loc(), ERR_UNEXPECTED_TOKEN)
		}

		return fp, nil
	}

	switch n.Type() {
	case N_TS_ANY, N_TS_NUM, N_TS_BOOL, N_TS_STR, N_TS_SYM:
		d := n.(*TsPredef)
		ti := p.newTypInfo()
		ti.SetQues(d.ques)
		return &Ident{N_NAME, d.loc, d.loc.Text(), false, false, nil, false, ti}, nil
	case N_TS_VOID:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	case N_TS_REF:
		r := n.(*TsRef)
		if r.HasArgs() {
			return nil, p.errorAtLoc(r.lt, ERR_UNEXPECTED_TOKEN)
		}
		return p.tsRoughParamToParam(r.name)
	case N_TS_LIT_OBJ:
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
		return &Ident{N_NAME, t.loc, t.loc.Text(), false, false, nil, true, p.newTypInfo()}, nil
	case N_TS_NS_NAME:
		s := n.(*TsNsName)
		return nil, p.errorAtLoc(s.dot, ERR_UNEXPECTED_TOKEN)
	case N_TS_CALL_SIG, N_TS_NEW_SIG, N_TS_FN_TYP:
		return nil, p.errorAtLoc(n.Loc(), ERR_UNEXPECTED_TOKEN)
	case N_TS_UNION_TYP:
		u := n.(*TsUnionTyp)
		return nil, p.errorAtLoc(u.op, ERR_UNEXPECTED_TOKEN)
	case N_TS_INTERSECT_TYP:
		i := n.(*TsIntersectTyp)
		return nil, p.errorAtLoc(i.op, ERR_UNEXPECTED_TOKEN)
	case N_TS_REST:
		n := n.(*TsRest)
		if n.arg.Type() != N_TS_REF {
			return nil, p.errorAtLoc(n.arg.Loc(), ERR_UNEXPECTED_TOKEN)
		}
		ref := n.arg.(*TsRef)
		if ref.args != nil {
			return nil, p.errorAtLoc(ref.args.Loc(), ERR_UNEXPECTED_TOKEN)
		}
		if ref.lt != nil {
			return nil, p.errorAtLoc(ref.lt, ERR_UNEXPECTED_TOKEN)
		}
		name := ref.name
		if name.Type() != N_NAME {
			return nil, p.errorAtLoc(name.Loc(), ERR_UNEXPECTED_TOKEN)
		}
		return &RestPat{N_PAT_REST, n.loc, n.arg, nil, p.newTypInfo()}, nil
	}
	return node, nil
}

func (p *Parser) tsPropToIdxSig(prop *TsProp) (Node, error) {
	if prop.key.Type() != N_NAME {
		return nil, p.errorAtLoc(prop.key.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	name := prop.key.(*Ident)
	if name.ti.TypAnnot() == nil {
		return nil, p.errorAtLoc(name.loc, ERR_UNEXPECTED_TOKEN)
	}
	switch name.ti.TypAnnot().tsTyp.Type() {
	case N_TS_NUM, N_TS_STR, N_TS_SYM:
		break
	default:
		return nil, p.errorAtLoc(name.ti.TypAnnot().Loc(), ERR_UNEXPECTED_TOKEN)
	}
	vt := prop.val.Type()
	if vt < N_TS_ANY || vt > N_TS_ROUGH_PARAM {
		return nil, p.errorAtLoc(prop.val.Loc(), ERR_UNEXPECTED_TOKEN)
	}
	typAnnot, err := p.tsRoughParamToTyp(prop.val, false)
	if err != nil {
		return nil, err
	}
	return &TsIdxSig{N_TS_IDX_SIG, prop.loc, name, typAnnot, nil}, nil
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
			return &TsPredef{typ, n.loc, nil, nil}, nil
		}
		return &TsRef{N_TS_REF, n.Loc().Clone(), n, nil, nil, nil}, nil
	case N_TS_LIT_OBJ:
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
	typParams, err := p.tsTryTypParams()
	if err != nil {
		return nil, err
	}

	params, _, loc, err := p.paramList(true, PK_NONE, false)
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
		typ, err := p.tsFnTyp(typParams, params, loc)
		if err != nil {
			return nil, err
		}
		if !keepParen {
			return typ, nil
		}
		return &TsParen{N_TS_PAREN, p.finLoc(loc), typ, nil}, nil
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
		return &TsParen{N_TS_PAREN, p.finLoc(loc), typ, nil}, nil
	}

	if len(params) == 0 {
		typ := &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), typParams, params, nil, nil}
		if !keepParen {
			return typ, nil
		}
		return &TsParen{N_TS_PAREN, p.finLoc(loc), typ, nil}, nil
	}
	return nil, p.errorTok(ahead)
}

// returns `PrimaryType` or `FunctionType`
func (p *Parser) tsFnTyp(typParams Node, params []Node, parenL *Loc) (Node, error) {
	var err error
	var loc *Loc
	if parenL != nil {
		loc = parenL
	}
	if typParams != nil {
		loc = typParams.Loc().Clone()
	}
	if params == nil {
		params, typParams, _, err = p.paramList(false, PK_NONE, true)
		if err != nil {
			return nil, err
		}
	}

	arrow, err := p.nextMustTok(T_ARROW)
	if err != nil {
		return nil, err
	}
	tiLoc := p.locFromTok(arrow)

	retTyp, err := p.tsTyp(false, false, false)
	if err != nil {
		return nil, err
	}

	ti := &TsTypAnnot{N_TS_TYP_ANNOT, p.finLoc(tiLoc), retTyp}
	return &TsFnTyp{N_TS_FN_TYP, p.finLoc(loc), typParams, params, ti, nil}, nil
}

func (p *Parser) tsTypName(ns Node) (Node, error) {
	if ns == nil {
		var err error
		ns, err = p.identWithKw(nil, false)
		if err != nil {
			return nil, err
		}
	}
	for {
		if p.lexer.Peek().value == T_DOT {
			loc := p.locFromTok(p.lexer.Next())
			id, err := p.identWithKw(nil, false)
			if err != nil {
				return nil, err
			}
			ns = &TsNsName{N_TS_NS_NAME, p.finLoc(ns.Loc().Clone()), ns, loc, id, nil}
		} else {
			break
		}
	}
	return ns, nil
}

// `typePredicates` and `assertPredicate`
func (p *Parser) tsTypPredicate(name Node, asserts bool, this bool) (Node, error) {
	loc := name.Loc().Clone()
	var err error
	if asserts {
		if this {
			name = &TsPredef{N_TS_THIS, p.finLoc(p.locFromTok(p.lexer.Next())), nil, nil}
		} else {
			name, err = p.ident(nil, true)
		}
		if err != nil {
			return nil, err
		}
	}

	ahead := p.lexer.Peek()
	av := ahead.value

	var typ Node
	if av == T_NAME && ahead.Text() == "is" {
		p.lexer.Next()

		typ, err = p.tsTyp(false, false, false)
		if err != nil {
			return nil, err
		}
		typ = &TsTypAnnot{N_TS_TYP_ANNOT, typ.Loc().Clone(), typ}
	}

	return &TsTypPredicate{N_TS_TYP_PREDICATE, p.finLoc(loc), name, typ, asserts, nil}, nil
}

func (p *Parser) tsRef(ns Node) (Node, error) {
	name, err := p.tsTypName(ns)
	if err != nil {
		return nil, err
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	asserts := p.isName(name, "asserts", false, false)
	// `assertPredicate` or `typePredicates`
	if (asserts && (av == T_NAME || av == T_THIS)) || ahead.Text() == "is" {
		return p.tsTypPredicate(name, asserts, av == T_THIS)
	}

	if av != T_LT {
		return &TsRef{N_TS_REF, p.finLoc(name.Loc().Clone()), name, nil, nil, nil}, nil
	}

	args, err := p.tsTypArgs(false, true)
	if err != nil {
		return nil, err
	}
	return &TsRef{N_TS_REF, p.finLoc(ns.Loc().Clone()), name, args.Loc().Clone(), args, nil}, nil
}

func (p *Parser) tsArr(typ Node) (Node, error) {
	for {
		if p.lexer.Peek().value == T_BRACKET_L {
			loc := p.locFromTok(p.lexer.Next())
			if p.lexer.Peek().value != T_BRACKET_R {
				idx, err := p.tsTyp(false, false, true)
				if err != nil {
					return nil, err
				}
				if _, err := p.nextMustTok(T_BRACKET_R); err != nil {
					return nil, err
				}
				typ = &TsIdxAccess{N_TS_IDX_ACCESS, p.finLoc(typ.Loc().Clone()), typ, idx, nil}
			} else {
				p.lexer.Next() // `]`
				typ = &TsArr{N_TS_ARR, p.finLoc(typ.Loc().Clone()), loc, typ, nil}
			}
		} else {
			break
		}
	}
	return typ, nil
}

func (p *Parser) tsToBeTupleLabel(node Node) Node {
	if node.Type() == N_TS_REF {
		n := node.(*TsRef)
		if n.lt == nil && n.args == nil && n.name.Type() == N_NAME {
			return n.name
		}
	}
	if pd, ok := node.(*TsPredef); ok {
		return &Ident{N_NAME, pd.loc, pd.loc.Text(), false, false, nil, true, p.newTypInfo()}
	}
	return nil
}

func (p *Parser) tsTupleNamedMember(orig Node, label Node, opt bool) (Node, error) {
	p.lexer.Next() // `:` or `?`
	if opt {
		if p.lexer.Peek().value != T_COLON {
			return &TsOpt{N_TS_OPT, p.finLoc(orig.Loc().Clone()), orig}, nil
		}
		p.lexer.Next()
	}
	arg, err := p.tsTyp(false, false, true)
	if err != nil {
		return nil, err
	}
	return &TsTupleNamedMember{N_TS_TUPLE_NAMED_MEMBER, p.finLoc(label.Loc().Clone()), label, opt, arg}, nil
}

func (p *Parser) tsTupleMember(rough bool) (Node, error) {
	ahead := p.lexer.Peek()
	av := ahead.value
	if av == T_DOT_TRI {
		loc := p.locFromTok(p.lexer.Next())
		arg, err := p.tsTupleMember(rough)
		if err != nil {
			return nil, err
		}

		argTyp := arg.Type()
		if !rough && argTyp != N_TS_ARR && argTyp != N_TS_TUPLE {
			if arg.Type() != N_TS_TUPLE_NAMED_MEMBER {
				return nil, p.errorAtLoc(arg.Loc(), ERR_REST_TYPE_SHOULD_BE_ARRAY)
			}
			mb := arg.(*TsTupleNamedMember)
			mt := mb.val.Type()
			if mt != N_TS_ARR && mt != N_TS_TUPLE {
				return nil, p.errorAtLoc(arg.Loc(), ERR_REST_TYPE_SHOULD_BE_ARRAY)
			}
		}
		return &TsRest{N_TS_REST, p.finLoc(loc), arg}, nil
	}

	arg, err := p.tsTyp(rough, false, false)
	if err != nil {
		return nil, err
	}

	ahead = p.lexer.Peek()
	av = ahead.value
	if av == T_COLON || av == T_HOOK {
		if label := p.tsToBeTupleLabel(arg); label != nil {
			return p.tsTupleNamedMember(arg, label, av == T_HOOK)
		}
		if av == T_HOOK {
			if _, ok := arg.(InParenNode); ok {
				p.lexer.Next() // `?`
				return &TsOpt{N_TS_OPT, p.finLoc(arg.Loc().Clone()), arg}, nil
			}
		}
		return nil, p.errorAtLoc(arg.Loc(), ERR_TUPLE_LABEL_SHOULD_BE_SIMPLE)
	}
	return arg, nil
}

func (p *Parser) tsTuple(rough bool) (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok)
	args := make([]Node, 0, 1)

	var arg Node
	var err error
	named := false
	opt := false
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACKET_R {
			p.lexer.Next()
			break
		} else {
			arg, err = p.tsTupleMember(rough)
			if err != nil {
				return nil, err
			}

			argTyp := arg.Type()
			if argTyp == N_TS_TUPLE_NAMED_MEMBER {
				named = true
				md := arg.(*TsTupleNamedMember)
				if !md.opt && opt {
					return nil, p.errorAtLoc(arg.Loc(), ERR_TUPLE_OPT_SHOULD_AFTER_REQUIRED)
				}
				opt = md.opt
			} else if argTyp == N_TS_OPT {
				opt = true
			} else if opt && argTyp != N_TS_REST {
				return nil, p.errorAtLoc(arg.Loc(), ERR_TUPLE_OPT_SHOULD_AFTER_REQUIRED)
			}
		}
		args = append(args, arg)
		if p.lexer.Peek().value == T_COMMA {
			p.lexer.Next()
		}
		if named {
			for _, mb := range args {
				typ := mb.Type()
				if typ != N_TS_TUPLE_NAMED_MEMBER && typ != N_TS_REST {
					return nil, p.errorAtLoc(mb.Loc(), ERR_TUPLE_NAMED_SHOULD_ALL_NAMED)
				}
				if typ == N_TS_REST {
					n := mb.(*TsRest)
					if n.arg.Type() != N_TS_TUPLE_NAMED_MEMBER {
						return nil, p.errorAtLoc(mb.Loc(), ERR_TUPLE_NAMED_SHOULD_ALL_NAMED)
					}
				}
			}
		}
	}
	return &TsTuple{N_TS_TUPLE, p.finLoc(loc), args, nil}, nil
}

func (p *Parser) tsPredefOrRef(tok *Token) (Node, error) {
	if tok == nil {
		tok = p.lexer.Peek()
	}
	tv := tok.value

	var node Node
	var err error
	var loc *Loc
	var name string
	if tv == T_NAME {
		node, err = p.ident(nil, false)
		if err != nil {
			return nil, err
		}
		name = node.(*Ident).Text()
		loc = node.Loc()
	} else if tv == T_VOID {
		name = "void"
		loc = p.locFromTok(p.lexer.Next())
	}

	if name == "intrinsic" && !p.scope().IsKind(SPK_TS_MAY_INTRINSIC) {
		node = &Ident{N_NAME, node.Loc(), "intrinsic", false, tok.ContainsEscape(), nil, tok.IsKw(), p.newTypInfo()}
	} else if typ, ok := builtinTyp[name]; ok {
		// predef
		if p.lexer.Peek().value != T_DOT {
			return &TsPredef{typ, loc, nil, nil}, nil
		}
	}

	return p.tsRef(node)
}

func (p *Parser) tsTypOp(op *Token) (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	opStr := op.Text()
	arg, err := p.tsTyp(false, false, true)
	if err != nil {
		return nil, err
	}
	return &TsTypOp{N_TS_TYP_OP, p.finLoc(loc), opStr, arg, nil}, nil
}

func (p *Parser) tsAheadIsRo(ahead *Token) (bool, error) {
	loc := p.locFromTok(ahead)
	if ahead.value != T_NAME || ahead.ContainsEscape() || ahead.Text() != "readonly" {
		return false, nil
	}
	ahead2nd := p.lexer.Peek2nd()
	if ahead2nd.value == T_BRACKET_L ||
		(ahead2nd.value == T_NAME && !ahead2nd.ContainsEscape() && ahead2nd.Text() != "readonly") {
		return true, nil
	}
	if ahead2nd.Kind().StartExpr {
		return false, p.errorAtLoc(loc, ERR_READONLY_ONLY_PERMITTED_ON_ARRAY_TUPLE)
	}
	return false, nil
}

// returns `FunctionType` or `PrimaryType` since `FunctionType`
// is conflicts with `ParenthesizedType`
func (p *Parser) tsPrimary(rough bool, canConst, canCond bool) (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	av := ahead.value
	if av == T_PAREN_L || av == T_LT {
		// paren type
		node, err := p.tsParen(true)
		if err != nil {
			return nil, err
		}
		if p.lexer.Peek().value == T_BRACKET_L {
			// array type
			node, err = p.tsArr(node)
			if err != nil {
				return nil, err
			}
		}
		if node.Type() == N_TS_PAREN {
			arg := node.(*TsParen).arg
			if pn, ok := arg.(InParenNode); ok {
				pn.SetOuterParen(node.Loc())
				node = arg
			}
		}
		return node, nil
	}

	var err error
	var node Node
	if av == T_NAME || av == T_VOID || (av == T_CONST && canConst) {
		if ok, _, new := p.tsAheadIsAbstract(ahead, false, false, true); ok && new {
			loc := p.locFromTok(p.lexer.Next())
			return p.tsConstructTyp(loc, true)
		}
		if av == T_NAME && !ahead.ContainsEscape() {
			str := ahead.Text()
			if str == "infer" {
				node, err = p.tsInfer()
			} else if str == "keyof" || str == "unique" {
				node, err = p.tsTypOp(ahead)
			} else {
				ok, err := p.tsAheadIsRo(ahead)
				if err != nil {
					return nil, err
				}
				if ok {
					loc := p.locFromTok(p.lexer.Next())
					arg, err := p.tsPrimary(false, false, false)
					if err != nil {
						return nil, err
					}
					if arg.Type() != N_TS_ARR && arg.Type() != N_TS_TUPLE {
						return nil, p.errorAtLoc(loc, ERR_READONLY_ONLY_PERMITTED_ON_ARRAY_TUPLE)
					}
					return &TsTypOp{N_TS_TYP_OP, p.finLoc(loc), "readonly", arg, nil}, nil
				}
			}
		}
		if node == nil {
			node, err = p.tsPredefOrRef(ahead)
		}
		if err != nil {
			return nil, err
		}
	} else if ahead.IsLit(true) {
		lit, err := p.primaryExpr(false)
		if err != nil {
			return nil, err
		}
		if lit.Type() == N_LIT_NULL {
			return &TsPredef{N_TS_NULL, lit.Loc().Clone(), nil, nil}, nil
		}
		return &TsLit{N_TS_LIT, lit.Loc().Clone(), lit, nil}, nil
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
		arg, err := p.tsTyp(false, false, canCond)
		if err != nil {
			return nil, err
		}
		node = &TsTypQuery{N_TS_TYP_QUERY, p.finLoc(loc), arg, nil}
	} else if av == T_THIS {
		// this type
		p.lexer.Next()
		node = &TsThis{N_TS_THIS, p.finLoc(loc), nil}
		ahead := p.lexer.Peek()
		if ahead.value == T_NAME && ahead.Text() == "is" {
			node, err = p.tsTypPredicate(node, false, false)
			if err != nil {
				return nil, err
			}
		}
	} else if av == T_IMPORT {
		// `let a: import("a")<a>;`
		return p.tsImport()
	} else if av == T_SUB {
		loc := p.locFromTok(p.lexer.Next())
		tok, err := p.nextMustTok(T_NUM)
		if err != nil {
			return nil, err
		}
		numLoc := p.locFromTok(tok)
		arg := &NumLit{N_LIT_NUM, p.finLoc(numLoc), nil}
		un := &UnaryExpr{N_EXPR_UNARY, p.finLoc(loc), T_SUB, arg, nil}
		return &TsLit{N_TS_LIT, un.Loc().Clone(), un, nil}, nil
	} else if av == T_TPL_HEAD {
		tpl, err := p.tplExpr(nil, true)
		if err != nil {
			return nil, err
		}
		return &TsLit{N_TS_LIT, tpl.Loc().Clone(), tpl, nil}, nil
	}

	if node != nil {
		ahead = p.lexer.Peek()
		av = ahead.value
		if node.Type() != N_TS_INTRINSIC && av == T_BRACKET_L && !ahead.afterLineTerm {
			// array type
			node, err = p.tsArr(node)
			if err != nil {
				return nil, err
			}
		} else if av == T_EXTENDS && canCond {
			return p.tsConditional(node)
		}
		return node, nil
	}

	return nil, p.errorTok(ahead)
}

func (p *Parser) tsImport() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	if _, err := p.nextMustTok(T_PAREN_L); err != nil {
		return nil, err
	}
	tok, err := p.nextMustTok(T_STRING)
	if err != nil {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_IMPORT_ARG_SHOULD_BE_STR)
	}
	strLoc := p.locFromTok(tok)
	arg := &StrLit{N_LIT_STR, p.finLoc(strLoc), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil, p.newTypInfo()}
	if _, err := p.nextMustTok(T_PAREN_R); err != nil {
		return nil, err
	}
	var qualifier Node
	if p.lexer.Peek().value == T_DOT {
		p.lexer.Next()
		qualifier, err = p.tsTypName(nil)
		if err != nil {
			return nil, err
		}
	}
	typArgs, err := p.tsTryTypArgs(nil, true)
	if err != nil {
		return nil, err
	}
	return &TsImportType{N_TS_IMPORT_TYP, p.finLoc(loc), arg, qualifier, typArgs, nil}, nil
}

func (p *Parser) tsMapped(loc *Loc, add, sub bool, roLoc *Loc, key Node) (Node, error) {
	var err error
	ro := 0
	if roLoc != nil {
		ro = 1
	}

	// if key is not nil then the preceding of key already been processed, so there
	// is no need to process `readonly` again
	if roLoc == nil && key == nil {
		if add || sub {
			p.lexer.Next()
			if _, err := p.nextMustName("readonly", false); err != nil {
				return nil, err
			}
			if add {
				ro = 2
			} else {
				ro = 3
			}
		} else {
			ahead := p.lexer.Peek()
			if ahead.value == T_NAME && !ahead.ContainsEscape() && ahead.Text() == "readonly" {
				p.lexer.Next()
				ro = 1
			}
		}
	}

	if key == nil {
		if _, err := p.nextMustTok(T_BRACKET_L); err != nil {
			return nil, err
		}

		key, err = p.tsTypParam(nil, false, true)
		if err != nil {
			return nil, err
		}
	}

	ahead := p.lexer.Peek()
	av := ahead.value
	var name Node
	if av == T_NAME && !ahead.ContainsEscape() && ahead.Text() == "as" {
		p.lexer.Next()
		name, err = p.tsTyp(false, false, false)
		if err != nil {
			return nil, err
		}
	}

	ahead = p.lexer.Peek()
	av = ahead.value
	if av == T_BRACKET_R {
		p.lexer.Next()
	} else {
		return nil, p.errorTok(ahead)
	}

	ahead = p.lexer.Peek()
	av = ahead.value
	opt := 0
	if av == T_ADD || av == T_SUB {
		p.lexer.Next()
		if _, err := p.nextMustTok(T_HOOK); err != nil {
			return nil, err
		}
		if av == T_ADD {
			opt = 2
		} else {
			opt = 3
		}
	} else if p.lexer.Peek().value == T_HOOK {
		p.lexer.Next()
		opt = 1
	}

	var val Node
	if p.lexer.Peek().value == T_COLON {
		p.lexer.Next()
		val, err = p.tsTyp(false, false, true)
		if err != nil {
			return nil, err
		}
	}

	p.advanceIfSemi(false)
	if _, err := p.nextMustTok(T_BRACE_R); err != nil {
		return nil, err
	}
	return &TsMapped{N_TS_MAPPED, p.finLoc(loc), ro, opt, key, name, val, nil}, nil
}

// this method can be used to parse
// - tsObj
// - the body of tsItf
// - mapped type - should be used with `tsMapped` method to cover all possible cases
func (p *Parser) tsObj(rough bool) (Node, error) {
	tok := p.lexer.Next() // `{`
	loc := p.locFromTok(tok)
	props := make([]Node, 0, 1)

	ahead := p.lexer.Peek()
	av := ahead.value
	if av == T_ADD || av == T_SUB {
		return p.tsMapped(loc, av == T_ADD, av == T_SUB, nil, nil)
	}

	i := 0
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_BRACE_R {
			p.lexer.Next()
			break
		}
		prop, err := p.tsProp(rough, i == 0, loc)
		if err != nil {
			return nil, err
		}
		if prop.Type() == N_TS_MAPPED {
			return prop, nil
		}
		props = append(props, prop)
		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_COMMA || av == T_SEMI {
			p.lexer.Next()
		}
	}
	return &TsObj{N_TS_LIT_OBJ, p.finLoc(loc), props, nil}, nil
}

func (p *Parser) tsNewSig(loc *Loc) (Node, error) {
	p.lexer.Next()
	params, typParams, _, err := p.paramList(false, PK_NONE, true)
	if err != nil {
		return nil, err
	}
	retTyp, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	p.advanceIfSemi(false)
	return &TsNewSig{N_TS_NEW_SIG, p.finLoc(loc), typParams, params, retTyp, false}, nil
}

// the returned node can be one of these types:
// - TsIdxSig
// - TsProp which holds the method declaration
// - TsMapped
func (p *Parser) tsIdxSig(loc *Loc, braceL *Loc, roLoc *Loc, canMapped bool) (Node, error) {
	tok := p.lexer.Next()
	bracketL := p.locFromTok(tok)

	if loc == nil {
		loc = p.locFromTok(tok)
	}

	p.scope().AddKind(SPK_NOT_IN)
	key, err := p.binExpr(nil, 0, false, false, false, false)
	p.scope().EraseKind(SPK_NOT_IN)
	if err != nil {
		return nil, err
	}

	if canMapped && key.Type() == N_NAME {
		key, err = p.tsTypParam(key, false, true)
		if err != nil {
			return nil, err
		}
		if key.(*TsParam).cons != nil {
			// return mapped
			return p.tsMapped(braceL, false, false, roLoc, key)
		}
		key = key.(*TsParam).name
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	if typAnnot != nil {
		if wt, ok := key.(NodeWithTypInfo); ok {
			wt.TypInfo().SetTypAnnot(typAnnot)
		} else {
			return nil, p.errorAtLoc(typAnnot.Loc(), ERR_UNEXPECTED_TOKEN)
		}
	}
	if _, err = p.nextMustTok(T_BRACKET_R); err != nil {
		return nil, err
	}

	// callSig can be optional, eg:
	//
	// ```
	// interface I {
	//   [Symbol.iterator]?(): number; // legal
	//   [a: string]?(): number;       // `?` is illegal
	// }
	// ```
	ques, _ := p.tsAdvanceHook(false)
	if ques != nil {
		if typAnnot != nil {
			return nil, p.errorAtLoc(ques, ERR_UNEXPECTED_TOKEN)
		}
	}

	val, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}

	// handle the callSig: `[Symbol.iterator](): void;`
	if val == nil {
		ahead := p.lexer.Peek()
		if p.aheadIsArgList(ahead) {
			if roLoc != nil {
				return nil, p.errorAtLoc(roLoc, ERR_INVALID_RO_MODIFIER_IN_TS_OBJ)
			}

			val, err = p.tsCallSig(nil, p.locFromTok(ahead), PK_INIT)
			if err != nil {
				return nil, err
			}
			return &TsProp{N_TS_PROP, p.finLoc(loc), key, val, ques, PK_METHOD, nil, false}, nil
		}
	}
	p.advanceIfSemi(false)

	if typAnnot != nil {
		return &TsIdxSig{N_TS_IDX_SIG, p.finLoc(loc), key, val, ques}, nil
	}
	return &TsProp{N_TS_PROP, p.finLoc(loc), key, val, ques, PK_INIT, bracketL, roLoc != nil}, nil
}

var modifiers = map[string]int{
	"private":   1,
	"public":    1,
	"protected": 1,
	"static":    1,
	"declare":   1,
	"abstract":  1,
}

func (p *Parser) tsPropReadonly(ahead *Token) *Loc {
	loc := p.locFromTok(ahead)
	if ahead.value != T_NAME || ahead.ContainsEscape() || ahead.Text() != "readonly" {
		return nil
	}
	ahead2nd := p.lexer.Peek2nd()
	if ahead2nd.value == T_BRACKET_L {
		p.lexer.Next()
		return loc
	}
	if _, _, ok := ahead2nd.CanBePropKey(); ok {
		return loc
	}
	return nil
}

func (p *Parser) tsProp(rough bool, canMapped bool, braceL *Loc) (Node, error) {
	ahead := p.lexer.Peek()
	loc := p.locFromTok(ahead)
	roLoc := p.tsPropReadonly(ahead)

	ahead = p.lexer.Peek()
	av := ahead.value
	if av == T_LT {
		if roLoc != nil {
			return nil, p.errorAtLoc(roLoc, ERR_INVALID_RO_MODIFIER_IN_TS_OBJ)
		}
		ps, err := p.tsTypParams()
		if err != nil {
			return nil, err
		}
		return p.tsCallSig(ps, loc, PK_NONE)
	}
	if av == T_NEW {
		if roLoc != nil {
			return nil, p.errorAtLoc(roLoc, ERR_INVALID_RO_MODIFIER_IN_TS_OBJ)
		}
		// ConstructSignature
		return p.tsNewSig(loc)
	}
	if !rough && av == T_BRACKET_L {
		// IndexSignature
		return p.tsIdxSig(loc, braceL, roLoc, canMapped)
	}
	if av == T_PAREN_L {
		if roLoc != nil {
			return nil, p.errorAtLoc(roLoc, ERR_INVALID_RO_MODIFIER_IN_TS_OBJ)
		}

		// CallSignature
		return p.tsCallSig(nil, loc, PK_NONE)
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

		kind := PK_INIT
		var ro *Loc
		if ques == nil && name.Type() == N_NAME {
			s := name.(*Ident).Text()
			ahead := p.lexer.Peek()
			if ahead.value == T_NAME && !ahead.afterLineTerm {
				if s == "set" || s == "get" {
					if s == "set" {
						kind = PK_SETTER
					} else {
						kind = PK_GETTER
					}

					// accessor should be followed by propName
					name, err = p.tsPropName()
					if err != nil {
						return nil, err
					}
				} else if _, ok := modifiers[s]; ok {
					return nil, p.errorAtLoc(name.Loc(), fmt.Sprintf(ERR_TPL_MODIFIER_ON_TYPE_MEMBER, s))
				} else if s == "readonly" {
					ro = name.Loc()
					name, err = p.tsPropName()
					if err != nil {
						return nil, err
					}
					name.(NodeWithTypInfo).TypInfo().SetReadonly(true)
				}
			}
		}

		typAnnot, err := p.tsTypAnnot()
		if err != nil {
			return nil, err
		}

		ahead := p.lexer.Peek()
		av := ahead.value
		if av != T_PAREN_L && av != T_LT {
			p.advanceIfSemi(false)
			return &TsProp{N_TS_PROP, p.finLoc(loc), name, typAnnot, ques, kind, nil, roLoc != nil}, nil
		}

		// MethodSignature is deserved
		if ro != nil {
			// method cannot be decorated by `readonly`
			return nil, p.errorAtLoc(ro, fmt.Sprintf(ERR_TPL_MODIFIER_ON_TYPE_MEMBER, "readonly"))
		}

		if kind == PK_INIT {
			kind = PK_METHOD
		}

		if roLoc != nil {
			return nil, p.errorAtLoc(roLoc, ERR_INVALID_RO_MODIFIER_IN_TS_OBJ)
		}
		callSig, err := p.tsCallSig(nil, p.locFromTok(p.lexer.Peek()), kind)
		if err != nil {
			return nil, err
		}

		return &TsProp{N_TS_PROP, p.finLoc(loc), name, callSig, ques, kind, nil, false}, nil
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
		value, err = p.tsTyp(rough, false, false)
		if err != nil {
			return nil, err
		}
	} else if assign {
		if key.Type() == N_NAME {
			return p.patternAssign(key, true)
		}
		return nil, p.errorTok(tok)
	}
	return &TsProp{N_TS_PROP, p.finLoc(key.Loc().Clone()), key, value, ques, PK_INIT, compute, roLoc != nil}, nil
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
	if keyName, kw, ok := tok.CanBePropKey(); ok {
		p.lexer.Next()
		return &Ident{N_NAME, p.finLoc(loc), keyName, false, tok.ContainsEscape(), nil, kw, p.newTypInfo()}, nil
	}
	return nil, p.errorTok(tok)
}

func (p *Parser) tsTypParam(id Node, ext, in bool) (Node, error) {
	var err error
	if id == nil {
		id, err = p.ident(nil, false)
		if err != nil {
			return nil, err
		}
	}
	var cons, val Node
	ahead := p.lexer.Peek()
	av := ahead.value
	if ext && av == T_EXTENDS {
		p.lexer.Next()
		cons, err = p.tsTyp(false, false, false)
		if err != nil {
			return nil, err
		}
		if p.lexer.Peek().value == T_ASSIGN {
			p.lexer.Next()
			val, err = p.tsTyp(false, false, false)
			if err != nil {
				return nil, err
			}
		}
	} else if in && av == T_NAME && ahead.Text() == "in" && !ahead.ContainsEscape() {
		p.lexer.Next()
		cons, err = p.tsTyp(false, false, true)
		if err != nil {
			return nil, err
		}
	}
	return &TsParam{N_TS_PARAM, p.finLoc(id.Loc().Clone()), id, cons, val}, nil
}

func (p *Parser) tsTryTypParams() (Node, error) {
	if !p.ts || p.lexer.Peek().value != T_LT {
		return nil, nil
	}
	return p.tsTypParams()
}

func (p *Parser) tsTypParams() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	ps := make([]Node, 0, 1)
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_GT {
			p.lexer.Next()
			break
		} else if av == T_COMMA {
			p.lexer.Next()
		}

		pa, err := p.tsTypParam(nil, true, false)
		if err != nil {
			return nil, err
		}
		ps = append(ps, pa)
	}
	if len(ps) == 0 {
		return nil, p.errorAtLoc(loc, ERR_EMPTY_TYPE_PARAM_LIST)
	}
	return &TsParamsDec{N_TS_PARAM_DEC, p.finLoc(loc), ps}, nil
}

// returned nodes maybe:
// - the superset of typeArgs which also includes `Constraint`
// - jsxELem if the `FEAT_JSX` is turned on
//
// the caller maybe required to check if the returned node is one of:
// - valid typeArgs, by further doing a `tsCheckTypArgs` subroutine
// - valid typeParams, by further doing a `tsTypArgsToTypParams` subroutine
// - jsxElem
//
// this method returns typeParams instead of typeArgs since the former is
// a superset of the later and in the calling point of this method, there is
// no enough information to determine whether the typeParams or typeArgs is satisfied
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
//
// jsxElem maybe the results, consider below cases:
//
// ````
// <T>() => {}         // illegal jsxElem missing the closing tag
// <T,>() => {}        // arrowExor with typParams `<T,>`
// <T extends D> => {} // arrowExor with typParams `<T extends D>`
// ```
func (p *Parser) tsTryTypArgs(asyncLoc *Loc, noJsx bool) (Node, error) {
	ahead := p.lexer.Peek()
	if ahead.value != T_LT {
		return nil, nil
	}
	if !p.ts && p.feat&FEAT_JSX != 0 {
		return p.jsx(true, false)
	}
	if asyncLoc != nil {
		return p.tsTryTypArgsAfterAsync(asyncLoc)
	}

	p.pushState()
	loc := p.locFromTok(ahead)
	node, err := p.tsTypArgs(true, noJsx || p.scope().IsKind(SPK_CLASS_EXTEND_SUPER))
	if err != nil {
		if err != errTypArgMaybeJsx {
			return nil, err
		}

		p.popState()
		jsx, err := p.jsx(true, true)
		if err != nil {
			if pe, ok := err.(*ParserError); ok {
				if pe.msg == ERR_UNTERMINATED_JSX_CONTENTS {
					pe.msg = ERR_JSX_TS_LT_AMBIGUITY
					pe.line = loc.begin.Line
					pe.col = loc.begin.Col
				}
			}
			return nil, err
		}
		return jsx, nil
	}
	p.discardState()

	return node, nil
}

// consider ambiguity of the first `<` in below example:
//
// ```
// async < a, b;
// async<T>() == 0;
// ```
//
// for avoiding lookbehind the process should accept the input as seqExpr then try to
// transform the subtree of seqExpr to typArgs if its followed by `>`
func (p *Parser) tsTryTypArgsAfterAsync(asyncLoc *Loc) (Node, error) {
	name := &Ident{N_NAME, asyncLoc, asyncLoc.Text(), false, false, nil, true, p.newTypInfo()}
	binExpr, err := p.binExpr(name, 0, false, false, true, false)
	if err != nil {
		return nil, err
	}
	ltLoc := binExpr.(*BinExpr).opLoc.Clone()
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

func (p *Parser) tsTypArgs(canConst bool, noJsx bool) (Node, error) {
	p.lexer.PushMode(LM_TS_TYP_ARG, true)
	defer p.lexer.PopMode()

	loc := p.locFromTok(p.lexer.Next()) // `<`
	errTypArgMissingGT.line = loc.begin.Line
	errTypArgMissingGT.col = loc.begin.Col

	args := make([]Node, 0, 1)
	jsx := p.feat&FEAT_JSX != 0
	for {
		ahead := p.lexer.Peek()
		av := ahead.value
		if av == T_GT {
			if !noJsx && jsx && len(args) == 0 {
				return nil, errTypArgMaybeJsx
			}

			p.lexer.Next()
			break
		} else if av == T_NAME || ahead.IsLit(true) || ahead.IsCtxKw() || av == T_LT || av == T_PAREN_L || av == T_BRACE_L {
			// next is typ, fallthrough to below `p.tsTyp` to handle this branch
		} else {
			return nil, errTypArgMissingGT
		}

		arg, err := p.tsTyp(false, canConst, false)
		if err != nil {
			if av == T_LT {
				return nil, errTypArgMissingGT
			}
			return nil, err
		}
		nameLike := p.tsIsNameLike(arg.Type())

		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_EXTENDS {
			id, err := p.tsPredefToName(arg)
			if err != nil {
				return nil, err
			}

			p.lexer.Next() // consume `extends`
			cons, err := p.tsTyp(false, false, false)
			if err != nil {
				return nil, err
			}
			arg = &TsParam{N_TS_PARAM, p.finLoc(id.Loc().Clone()), id, cons, nil}
		} else if !noJsx && jsx && (av == T_NAME || ahead.IsKw() || av == T_DIV || av == T_BRACE_L || (av == T_GT && nameLike && len(args) == 0)) {
			return nil, errTypArgMaybeJsx
		}

		args = append(args, arg)

		ahead = p.lexer.Peek()
		if ahead.value == T_COMMA {
			p.lexer.Next()
		} else if ahead.value == T_GT {
			p.lexer.Next()
			break
		} else {
			return nil, errTypArgMissingGT
		}
	}
	if len(args) == 0 {
		return nil, p.errorAtLoc(loc, ERR_TYPE_ARG_EMPTY)
	}
	return &TsParamsInst{N_TS_PARAM_INST, p.finLoc(loc), args}, nil
}

func (p *Parser) tsCallSig(typParams Node, loc *Loc, kind PropKind) (Node, error) {
	if kind == PK_METHOD || kind == PK_GETTER || kind == PK_SETTER {
		p.scope().AddKind(SPK_METHOD)
	}

	params, tp, _, err := p.paramList(false, kind, typParams == nil)
	if err != nil {
		return nil, err
	}
	if tp != nil && (kind == PK_GETTER || kind == PK_SETTER) {
		return nil, p.errorAtLoc(tp.Loc(), ERR_ACCESSOR_WITH_TYPE_PARAMS)
	}

	opts := NewTsCheckParamOpts()
	opts.getter = kind == PK_GETTER
	opts.setter = kind == PK_SETTER
	opts.loc = loc
	if err = p.tsCheckParams(params, opts); err != nil {
		return nil, err
	}

	typAnnot, err := p.tsTypAnnot()
	if err != nil {
		return nil, err
	}
	if typAnnot != nil && opts.setter {
		return nil, p.errorAtLoc(typAnnot.Loc(), ERR_SETTER_WITH_RET_TYP)
	}

	if tp != nil {
		typParams = tp
	}
	p.advanceIfSemi(false)
	return &TsCallSig{N_TS_CALL_SIG, p.finLoc(loc), typParams, params, typAnnot}, nil
}

func (p *Parser) aheadIsTsTypDec(tok *Token, asterisk bool) bool {
	if !p.ts || tok.value != T_NAME || tok.Text() != "type" || tok.ContainsEscape() {
		return false
	}
	ahead := p.lexer.Peek2nd()
	return !ahead.afterLineTerm && (ahead.Kind().StartExpr || asterisk && ahead.value == T_MUL)
}

func (p *Parser) tsTypDec(loc *Loc, skipDec bool) (Node, error) {
	name, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}

	ref := NewRef()
	ref.Def = name
	ref.BindKind = BK_LET
	ref.Typ = RDT_TYPE
	if err := p.addLocalBinding(nil, ref, true, ref.Def.Text()); err != nil {
		return nil, err
	}

	ti := p.newTypInfo()
	if skipDec {
		return name, nil
	}

	params, err := p.tsTryTypParams()
	if err != nil {
		return nil, err
	}

	ahead := p.lexer.Peek()
	if ahead.value == T_ASSIGN {
		p.lexer.Next()
	} else {
		return nil, p.errorTok(ahead)
	}

	p.scope().AddKind(SPK_TS_MAY_INTRINSIC).AddKind(SPK_TS_DECLARE)
	typAnnot, err := p.tsTyp(false, false, true)
	p.scope().EraseKind(SPK_TS_MAY_INTRINSIC).EraseKind(SPK_TS_DECLARE)
	if err != nil {
		return nil, err
	}

	if err := p.advanceIfSemi(true); err != nil {
		return nil, err
	}

	ti.SetTypParams(params)
	ti.SetTypAnnot(typAnnot)
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
	if !p.ts || p.lastTsFnSig == nil {
		return nil
	}

	ecp := p.nameOfNode(p.lastTsFnSig.id)
	act := p.nameOfNode(id)
	if ecp == act {
		return nil
	}
	return p.errorAtLoc(id.(*Ident).loc, fmt.Sprintf(ERR_TPL_INVALID_FN_IMPL_NAME, ecp))
}

func (p *Parser) aheadIsTsItf(tok *Token) bool {
	if !p.ts || tok.value != T_INTERFACE {
		return false
	}
	ahead := p.lexer.Peek2nd()
	return !ahead.afterLineTerm
}

func (p *Parser) tsItfExtClause() ([]Node, error) {
	ns := make([]Node, 0, 1)
	var loc *Loc
	if p.lexer.Peek().value == T_EXTENDS {
		loc = p.locFromTok(p.lexer.Next())
	} else {
		return ns, nil
	}

	if p.lexer.Peek().value == T_BRACE_L {
		return nil, p.errorAtLoc(loc, ERR_EXTEND_LIST_EMPTY)
	}

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
	ref := NewRef()
	ref.Def = name
	ref.BindKind = BK_CONST
	ref.Typ = RDT_ITF | RDT_TYPE
	if err := p.addLocalBinding(nil, ref, true, ref.Def.Text()); err != nil {
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

	scope := p.symtab.EnterScope(false, false, true)
	scope.AddKind(SPK_TS_INTERFACE)
	body, err := p.tsObj(false)
	if err != nil {
		return nil, err
	}
	p.symtab.LeaveScope()

	itfBody := &TsInterfaceBody{
		typ:  N_TS_INTERFACE_BODY,
		loc:  body.(*TsObj).loc,
		body: body.(*TsObj).props,
	}
	return &TsInterface{N_TS_INTERFACE, p.finLoc(loc), name, params, supers, itfBody}, nil
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
			val, err = p.assignExpr(false, false, false, false)
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
	if !p.ts {
		return false
	}
	if tok == nil {
		tok = p.lexer.Peek()
	}
	return tok.value == T_ENUM
}

// `loc` is the loc of the preceding `const`
func (p *Parser) tsEnum(loc *Loc, cst bool) (Node, error) {
	cons := loc != nil
	tok := p.lexer.Next() // enum
	if loc == nil {
		loc = p.locFromTok(tok)
	}

	name, err := p.ident(nil, true)
	if err != nil {
		return nil, err
	}
	ref := NewRef()
	ref.Def = name
	ref.BindKind = BK_CONST
	ref.Typ = RDT_TYPE
	if cst {
		ref.Typ = ref.Typ.On(RDT_CONST_ENUM)
	} else {
		ref.Typ = ref.Typ.On(RDT_ENUM)
	}
	if err := p.addLocalBinding(nil, ref, true, ref.Def.Text()); err != nil {
		return nil, err
	}

	mems, err := p.tsEnumBody()
	if err != nil {
		return nil, err
	}
	return &TsEnum{N_TS_ENUM, p.finLoc(loc), name, mems, cons}, nil
}

// produces either `ImportAliasDeclaration` or `ImportRequireDeclaration`
//
// `typ` means the caller has met the `type` keyword: `export import type A = B.C;`
func (p *Parser) tsImportAlias(loc *Loc, name Node, export bool) (Node, error) {
	var err error
	ahead := p.lexer.Peek()
	var typ *Loc
	if ahead.value == T_NAME && name.(*Ident).Text() == "type" && !name.(*Ident).containsEscape {
		typ = name.Loc()
		name, err = p.ident(nil, true)
		if err != nil {
			return nil, err
		}
	}

	return p.tsImportAliasRsh(loc, typ, name, export)
}

func (p *Parser) tsImportAliasRsh(loc *Loc, typ *Loc, name Node, export bool) (Node, error) {
	if _, err := p.nextMustTok(T_ASSIGN); err != nil {
		return nil, err
	}

	val, err := p.tsTypName(nil)
	if err != nil {
		return nil, err
	}

	var node Node
	if val.Type() == N_NAME && val.(*Ident).Text() == "require" {
		call, _, err := p.callExpr(val, true, false, nil, false)
		if err != nil {
			return nil, err
		}

		// the arguments count of `require` should be one with type `StrLit`
		ce := call.(*CallExpr)
		if len(ce.args) == 0 {
			return nil, p.errorAt(p.lexer.PrevTok(), p.lexer.PrevTokBegin(), ERR_IMPORT_REQUIRE_STR_LIT_DESERVED)
		}
		if ce.args[0].Type() != N_LIT_STR {
			return nil, p.errorAtLoc(ce.args[0].Loc(), ERR_IMPORT_REQUIRE_STR_LIT_DESERVED)
		}

		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
		node = &TsImportRequire{N_TS_IMPORT_REQUIRE, p.finLoc(loc), name, call, nil}
	} else {
		if typ != nil {
			// `export import type A = B.C;`
			return nil, p.errorAtLoc(typ, ERR_IMPORT_TYPE_IN_IMPORT_ALIAS)
		}
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
		node = &TsImportAlias{N_TS_IMPORT_ALIAS, p.finLoc(loc), name, val, export}
	}

	ref := NewRef()
	ref.Def = name.(*Ident)
	ref.BindKind = BK_LET
	ref.Typ = RDT_TYPE
	if err := p.addLocalBinding(nil, ref, true, ref.Def.Text()); err != nil {
		return nil, err
	}

	return node, nil
}

func (p *Parser) aheadIsTsNS(tok *Token) bool {
	if !p.ts || tok.value != T_NAME {
		return false
	}
	str := tok.Text()
	ahead := p.lexer.Peek2nd()
	return (str == "namespace" || str == "as") && !ahead.afterLineTerm
}

func (p *Parser) tsNS() (Node, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok) // `namespace` or `as`
	as := tok.Text() == "as"
	if as {
		if _, err := p.nextMustName("namespace", false); err != nil {
			return nil, err
		}
	}

	name, err := p.tsTypName(nil)
	if err != nil {
		return nil, err
	}

	var blk Node
	if !as {
		blk, err = p.blockStmt(true, SPK_TS_MODULE)
		if err != nil {
			return nil, err
		}
	} else {
		p.advanceIfSemi(false)
	}

	// for namespace with qualified name, it will be split to
	// multiple modules and those modules will be constructed to
	// a tree structure whose children keep the order in that
	// qualified name
	n := name
	var mod *TsNS
	if n.Type() == N_TS_NS_NAME {
		modChain := make([]Node, 0, 2)

		for {
			ns := n.(*TsNsName)
			nestName := ns.rhs
			mod = &TsNS{N_TS_NAMESPACE, nil, nestName, nil, false}
			ml := len(modChain)
			if ml == 0 {
				mod.body = blk
			} else {
				mc, last := modChain[:ml-1], modChain[ml-1]
				mod.body = last
				modChain = mc
			}
			mod.loc = NewLocFromSpan(nestName, mod.body)
			n = ns.lhs
			if n.Type() == N_NAME {
				mod = &TsNS{N_TS_NAMESPACE, NewLocFromSpan(n, mod), n, mod, false}
				mod.loc.begin = loc.begin
				mod.loc.rng.start = loc.rng.start
				break
			} else {
				modChain = append(modChain, mod)
			}
		}
	}

	def := name
	if mod != nil {
		def = mod.name
	}
	ref := NewRef()
	ref.Def = def.(*Ident)
	ref.BindKind = BK_CONST
	ref.Typ = RDT_NS | RDT_TYPE
	if err := p.addLocalBinding(nil, ref, true, ref.Def.Text()); err != nil {
		return nil, err
	}

	if mod != nil {
		return mod, nil
	}
	return &TsNS{N_TS_NAMESPACE, p.finLoc(loc), name, blk, as}, nil
}

func (p *Parser) aheadIsTsDec(tok *Token) bool {
	if !p.ts || tok.value != T_NAME || tok.Text() != "declare" {
		return false
	}
	ahead := p.lexer.Peek2nd()
	return !ahead.afterLineTerm
}

func (p *Parser) aheadIsModDec(tok *Token) bool {
	if !p.ts || tok.value != T_NAME {
		return false
	}
	str := tok.Text()
	ahead := p.lexer.Peek2nd()
	return (str == "module" || str == "global") && !ahead.afterLineTerm
}

func (p *Parser) tsModDec() (*TsDec, error) {
	tok := p.lexer.Next()
	loc := p.locFromTok(tok) // `module` or `global`
	global := tok.Text() == "global"

	tok = p.lexer.Peek()
	var name Node
	var err error
	var str bool
	if tok.value == T_STRING {
		loc := p.locFromTok(p.lexer.Next())
		if !p.scope().IsKind(SPK_TS_DECLARE) {
			return nil, p.errorAtLoc(loc, ERR_ONLY_AMBIENT_MOD_WITH_STR_NAME)
		}
		str = true
		name = &StrLit{N_LIT_STR, p.finLoc(p.locFromTok(tok)), tok.Text(), tok.HasLegacyOctalEscapeSeq(), nil, nil}
	} else if global {
		name = &Ident{N_NAME, p.finLoc(loc.Clone()), "global", false, false, nil, true, p.newTypInfo()}
	} else {
		name, err = p.identStrict(nil, false, false)
		if err != nil {
			return nil, err
		}
	}

	var blk Node
	if str && p.lexer.Peek().value != T_BRACE_L {
		p.advanceIfSemi(false)
	} else {
		blk, err = p.blockStmt(true, SPK_TS_MODULE)
		if err != nil {
			return nil, err
		}
	}

	typ := N_TS_DEC_MODULE
	if global {
		typ = N_TS_DEC_GLOBAL
	}
	return &TsDec{typ, p.finLoc(loc), name, blk}, nil
}

func (p *Parser) tsDec() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())

	tok := p.lexer.Peek()
	tv := tok.value

	scope := p.scope()
	scope.AddKind(SPK_TS_DECLARE)

	var err error
	typ := N_ILLEGAL
	dec := &TsDec{typ, nil, nil, nil}
	if ok, kind := p.aheadIsVarDec(tok); ok {
		dec.inner, err = p.varDecStmt(kind, false)
		typ = N_TS_DEC_VAR_DEC
		if dec.inner.Type() == N_TS_ENUM {
			typ = N_TS_DEC_ENUM
		}
	} else if tv == T_FUNC {
		dec.inner, err = p.fnDec(false, nil, false)
		if err != nil {
			return nil, err
		}
		typ = N_TS_DEC_FN
		if err := p.advanceIfSemi(true); err != nil {
			return nil, err
		}
	} else if p.aheadIsAsync(tok, false, false) {
		if tok.ContainsEscape() {
			return nil, p.errorAt(tv, &tok.begin, ERR_ESCAPE_IN_KEYWORD)
		}
		return nil, p.errorAt(tv, &tok.begin, ERR_ASYNC_IN_AMBIENT)
	} else if tv == T_CLASS {
		dec.inner, err = p.classDec(false, false, true, false)
		typ = N_TS_DEC_CLASS
	} else if p.aheadIsTsItf(tok) {
		dec.inner, err = p.tsItf()
		typ = N_TS_DEC_INTERFACE
	} else if p.aheadIsTsTypDec(tok, false) {
		loc := p.locFromTok(p.lexer.Next())
		dec.inner, err = p.tsTypDec(loc, false)
		typ = N_TS_DEC_TYP_DEC
	} else if p.aheadIsTsEnum(tok) {
		dec.inner, err = p.tsEnum(nil, false)
		typ = N_TS_DEC_ENUM
	} else if p.aheadIsTsNS(tok) {
		dec.inner, err = p.tsNS()
		typ = N_TS_DEC_NS
	} else if p.aheadIsModDec(tok) {
		dec, err = p.tsModDec()
		if dec != nil {
			typ = dec.typ
		}
	} else if ok, itf, _ := p.tsAheadIsAbstract(tok, false, false, false); ok {
		if itf {
			return nil, p.errorAtLoc(p.locFromTok(tok), ERR_ABSTRACT_AT_INVALID_POSITION)
		}
		dec.inner, err = p.classDec(false, false, true, true)
		typ = N_TS_DEC_CLASS
	} else {
		return nil, p.errorAt(tok.value, &tok.begin, ERR_EXPORT_DECLARE_MISSING_DECLARATION)
	}

	if err != nil {
		return nil, err
	}

	if err = p.checkAmbient(typ, dec.inner); err != nil {
		return nil, err
	}

	dec.typ = typ
	dec.loc = p.finLoc(loc)
	scope.EraseKind(SPK_TS_DECLARE)

	return dec, nil
}

func (p *Parser) checkAmbient(typ NodeType, dec Node) error {
	switch typ {
	case N_TS_DEC_VAR_DEC:
		n := dec.(*VarDecStmt)
		for _, v := range n.decList {
			init := v.(*VarDec).init
			if init != nil {
				return p.errorAtLoc(init.Loc(), ERR_INIT_IN_ALLOWED_CTX)
			}
		}
	}
	return nil
}

func (p *Parser) tsNoNull(node Node) Node {
	if !p.ts {
		return node
	}

	ahead := p.lexer.Peek()
	if ahead.afterLineTerm || ahead.value != T_NOT {
		return node
	}

	p.lexer.NextRevise(T_TS_NO_NULL)
	return &TsNoNull{N_TS_NO_NULL, p.finLoc(node.Loc().Clone()), node, nil}
}

// the ts node which are valid at the left hand side of assignExpr
func (p *Parser) isTsLhs(node Node) bool {
	if !p.ts {
		return false
	}

	nt := node.Type()
	return nt == N_TS_NO_NULL || nt == N_TS_TYP_ASSERT
}

// `new` for expr: `let x: abstract new () => void = X;`
func (p *Parser) tsAheadIsAbstract(tok *Token, prop bool, pvt bool, new bool) (bool, bool, bool) {
	if p.ts && IsName(tok, "abstract", false) {
		ahead := p.lexer.Peek2nd()
		if ahead.afterLineTerm {
			return false, false, false
		}
		av := ahead.value
		if av == T_CLASS ||
			av == T_INTERFACE ||
			(new && av == T_NEW) ||
			(p.aheadIsArgList(ahead) && !prop) ||
			av == T_MUL {
			return true, av == T_INTERFACE, av == T_NEW
		}
		_, _, canProp := ahead.CanBePropKey()
		if prop && (av == T_BRACKET_L || av == T_NAME || av == T_STRING || canProp) {
			return true, false, false
		}
		if pvt && av == T_NAME_PVT {
			return true, false, false
		}
		if av == T_NAME {
			if p.scope().IsKind(SPK_NOT_IN) && (IsName(ahead, "in", false) || IsName(ahead, "of", false)) {
				return false, false, false
			}
			return true, false, false
		}
	}
	return false, false, false
}

type ModifierNameLoc struct {
	name string
	loc  *Loc
	skip []string
}

func (p *Parser) tsCheckLabeledOrder(orders []ModifierNameLoc) error {
	stuff := []ModifierNameLoc{}
	for _, od := range orders {
		if od.loc != nil {
			stuff = append(stuff, od)
		}
	}
	for i := 0; i < len(stuff)-1; i++ {
		a := stuff[i]
		b := stuff[i+1]
		if !a.loc.Before(b.loc) {
			skipped := false
			if a.skip != nil {
				for _, s := range a.skip {
					if s == b.name {
						skipped = true
						break
					}
				}
			}
			if !skipped {
				return p.errorAtLoc(a.loc, fmt.Sprintf(ERR_TPL_INVALID_MODIFIER_ORDER, a.name, b.name))
			}
		}
	}
	return nil
}

func (p *Parser) tsModifierOrder(staticLoc, overrideLoc, readonlyLoc, accessLoc, abstractLoc, declareLoc *Loc, accMod ACC_MOD, mayStaticBlock bool) error {
	if staticLoc != nil && abstractLoc != nil {
		return p.errorAtLoc(abstractLoc, ERR_ABSTRACT_MIXED_WITH_STATIC)
	}
	if declareLoc != nil && overrideLoc != nil {
		return p.errorAtLoc(overrideLoc, ERR_DECLARE_MIXED_WITH_OVERRIDE)
	}

	if staticLoc != nil && p.lexer.Peek().value == T_BRACE_L &&
		(accessLoc != nil || overrideLoc != nil || readonlyLoc != nil || declareLoc != nil) {
		if mayStaticBlock {
			return p.errorAtLoc(staticLoc, ERR_STATIC_BLOCK_WITH_MODIFIER)
		}
		return p.errorAtLoc(p.locFromTok(p.lexer.Peek()), ERR_UNEXPECTED_TOKEN)
	}

	orders := []ModifierNameLoc{
		{accMod.String(), accessLoc, nil},
		{"abstract", abstractLoc, []string{"readonly"}},
		{"static", staticLoc, nil},
		{"override", overrideLoc, nil},
		{"readonly", readonlyLoc, nil},
	}
	return p.tsCheckLabeledOrder(orders)
}

// process the implementation list in the class declaration
func (p *Parser) tsImplements() ([]Node, error) {
	ahead := p.lexer.Peek()
	av := ahead.value
	if av != T_IMPLEMENTS {
		return nil, nil
	}
	implLoc := p.locFromTok(p.lexer.Next())
	impl := []Node{}
	for {
		ahead = p.lexer.Peek()
		av = ahead.value
		if av == T_BRACE_L {
			break
		}
		typ, err := p.tsTyp(false, false, false)
		if err != nil {
			return nil, err
		}
		impl = append(impl, typ)
	}
	if len(impl) == 0 {
		return nil, p.errorAtLoc(implLoc, ERR_IMPLEMENT_LIST_EMPTY)
	}
	return impl, nil
}

func (p *Parser) tsConditional(chk Node) (Node, error) {
	loc := chk.Loc().Clone()

	var err error
	if _, err = p.nextMustTok(T_EXTENDS); err != nil {
		return nil, err
	}
	ext, err := p.tsTyp(false, false, false)
	if err != nil {
		return nil, err
	}
	if _, err = p.nextMustTok(T_HOOK); err != nil {
		return nil, err
	}
	trueTyp, err := p.tsTyp(false, false, true)
	if err != nil {
		return nil, err
	}
	if _, err = p.nextMustTok(T_COLON); err != nil {
		return nil, err
	}
	falseTyp, err := p.tsTyp(false, false, true)
	if err != nil {
		return nil, err
	}
	return &TsCondType{N_TS_COND, p.finLoc(loc), chk, ext, trueTyp, falseTyp, nil}, nil
}

func (p *Parser) tsInfer() (Node, error) {
	loc := p.locFromTok(p.lexer.Next())
	arg, err := p.tsTypName(nil)
	if err != nil {
		return nil, err
	}
	return &TsTypInfer{N_TS_TYP_INFER, p.finLoc(loc), arg, nil}, nil
}

func (p *Parser) tsCheckParams(params []Node, opts *TsCheckParamOpts) error {
	if opts.setter && len(params) == 0 {
		return p.errorAtLoc(opts.loc, ERR_SETTER_MISSING_PARAM)
	}

	scope := p.scope()
	opts.inDeclare = scope.IsKind(SPK_TS_DECLARE) || scope.IsKind(SPK_TS_INTERFACE) || p.feat&FEAT_DTS != 0
	for _, param := range params {
		if err := p.tsCheckParam(param, opts); err != nil {
			return err
		}
	}
	return nil
}

type TsCheckParamOpts struct {
	// in ts declaration scope
	inDeclare bool
	// has implementation
	impl bool
	// setter
	setter bool
	// getter
	getter bool
	// raise error at this LoC
	loc *Loc
}

func NewTsCheckParamOpts() *TsCheckParamOpts {
	return &TsCheckParamOpts{}
}

func (p *Parser) tsCheckParam(node Node, opts *TsCheckParamOpts) error {
	var ti *TypInfo
	if wt, ok := node.(NodeWithTypInfo); ok {
		ti = wt.TypInfo()
	}

	if opts.getter {
		return p.errorAtLoc(node.Loc(), ERR_GETTER_WITH_PARAMS)
	}

	if opts.setter && ti != nil && ti.Optional() {
		return p.errorAtLoc(ti.Ques(), ERR_SETTER_WITH_PARAM_OPTIONAL)
	}

	switch node.Type() {
	case N_PAT_ARRAY:
		// `[]?` in `declare function foo([]?): void` is legal
		// `[]?` in `function foo([]?): void {}` is illegal
		n := node.(*ArrPat)
		if n.ti.ques != nil && (!opts.inDeclare || opts.impl) {
			return p.errorAtLoc(n.ti.ques, ERR_BINDING_PATTERN_REQUIRE_IN_IMPL)
		}
	case N_PAT_REST:
		if opts.setter {
			return p.errorAtLoc(node.Loc(), ERR_SETTER_WITH_REST_PARAM)
		}
	}
	return nil
}

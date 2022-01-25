package estree

import "github.com/hsiaosiyuan0/mole/ecma/parser"

func convertTsTyp(node parser.Node) Node {
	if node == nil {
		return nil
	}

	switch node.Type() {
	case parser.N_NAME:
		return ident(node)
	case parser.N_TS_TYP_ANNOT:
		n := node.(*parser.TsTypAnnot)
		return &TSTypeAnnotation{
			Type:           "TSTypeAnnotation",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			TypeAnnotation: convertTsTyp(n.TsTyp()),
		}
	case parser.N_TS_NUM:
		return &TSNumberKeyword{
			Type:  "TSNumberKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_STR:
		return &TSStringKeyword{
			Type:  "TSStringKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_ANY:
		return &TSAnyKeyword{
			Type:  "TSAnyKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_BOOL:
		return &TSBooleanKeyword{
			Type:  "TSBooleanKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_VOID:
		return &TSVoidKeyword{
			Type:  "TSVoidKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_THIS:
		return &TSThisType{
			Type:  "TSThisType",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_UNKNOWN:
		return &TSUnknownKeyword{
			Type:  "TSUnknownKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_OBJ:
		return &TSObjectKeyword{
			Type:  "TSObjectKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_REF:
		n := node.(*parser.TsRef)
		return &TSTypeReference{
			Type:           "TSTypeReference",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			TypeName:       convert(n.Name()),
			TypeParameters: convert(n.ParamsInst()),
		}
	case parser.N_TS_PARAM_INST:
		n := node.(*parser.TsParamsInst)
		return &TSTypeParameterInstantiation{
			Type:   "TSTypeParameterInstantiation",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Params: elems(n.Params()),
		}
	case parser.N_TS_PARAM_DEC:
		n := node.(*parser.TsParamsDec)
		return &TSTypeParameterDeclaration{
			Type:   "TSTypeParameterDeclaration",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Params: elems(n.Params()),
		}
	case parser.N_TS_PARAM:
		n := node.(*parser.TsParam)
		return &TSTypeParameter{
			Type:       "TSTypeParameter",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Name:       convert(n.Name()),
			Constraint: convert(n.Cons()),
			Default:    convert(n.Default()),
		}
	case parser.N_TS_ARR:
		n := node.(*parser.TsArr)
		return &TSArrayType{
			Type:        "TSArrayType",
			Start:       start(node.Loc()),
			End:         end(node.Loc()),
			Loc:         loc(node.Loc()),
			ElementType: convertTsTyp(n.Arg()),
		}
	case parser.N_TS_LIT_OBJ:
		n := node.(*parser.TsObj)
		return &TSTypeLiteral{
			Type:    "TSTypeLiteral",
			Start:   start(node.Loc()),
			End:     end(node.Loc()),
			Loc:     loc(node.Loc()),
			Members: elems(n.Props()),
		}
	case parser.N_TS_CALL_SIG:
		n := node.(*parser.TsCallSig)
		return &TSFunctionDeclaration{
			Type:           "FunctionExpression",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Params:         fnParams(n.Params()),
			TypeParameters: convertTsTyp(n.TypParams()),
			ReturnType:     convertTsTyp(n.RetTyp()),
		}
	case parser.N_TS_PROP:
		n := node.(*parser.TsProp)
		if n.IsMethod() {
			return &TSMethodSignature{
				Type:     "TSMethodSignature",
				Start:    start(node.Loc()),
				End:      end(node.Loc()),
				Loc:      loc(node.Loc()),
				Key:      convert(n.Key()),
				Value:    convert(n.Val()),
				Computed: n.Computed(),
				Optional: n.Optional(),
			}
		}
		return &TSPropertySignature{
			Type:           "TSPropertySignature",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			Key:            convertTsTyp(n.Key()),
			Optional:       n.Optional(),
			Computed:       n.Computed(),
			TypeAnnotation: convertTsTyp(n.Val()),
		}
	case parser.N_TS_TYP_PREDICATE:
		n := node.(*parser.TsTypPredicate)
		return &TSTypePredicate{
			Type:           "TSTypePredicate",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			ParameterName:  convert(n.Name()),
			TypeAnnotation: convertTsTyp(n.Typ()),
			Asserts:        n.Asserts(),
		}
	case parser.N_TS_DEC_FN:
		n := node.(*parser.TsDec).Inner().(*parser.FnDec)
		ti := n.TypInfo()
		lc := parser.LocWithTypeInfo(node, false)
		return &TSDeclareFunction{
			Type:           "TSDeclareFunction",
			Start:          start(lc),
			End:            end(lc),
			Loc:            loc(lc),
			Id:             convert(n.Id()),
			Params:         fnParams(n.Params()),
			Body:           convert(n.Body()),
			Generator:      false,
			Async:          n.Async(),
			TypeParameters: typParams(ti),
			ReturnType:     typAnnot(ti),
		}
	case parser.N_TS_TYP_ASSERT:
		n := node.(*parser.TsTypAssert)
		return &TSTypeAssertion{
			Type:           "TSTypeAssertion",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			Expression:     convert(n.Expr()),
			TypeAnnotation: convertTsTyp(n.Typ()),
		}
	case parser.N_TS_NO_NULL:
		n := node.(*parser.TsNoNull)
		return &TSNonNullExpression{
			Type:       "TSNonNullExpression",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: convert(n.Arg()),
		}
	case parser.N_TS_UNION_TYP:
		n := node.(*parser.TsUnionTyp)
		return &TSUnionType{
			Type:  "TSUnionType",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Types: elems(n.Elems()),
		}
	case parser.N_TS_INTERSEC_TYP:
		n := node.(*parser.TsIntersecTyp)
		return &TSIntersectionType{
			Type:  "TSIntersectionType",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Types: elems(n.Elems()),
		}
	case parser.N_TS_DEC_CLASS:
		n := node.(*parser.TsDec)
		cls := n.Inner().(*parser.ClassDec)
		return &TSClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      start(n.Loc()),
			End:        end(cls.Loc()),
			Loc:        loc(n.Loc()),
			Id:         convert(cls.Id()),
			SuperClass: convert(cls.Super()),
			Body:       convert(cls.Body()),
			Declare:    true,
			Abstract:   cls.Abstract(),
		}
	case parser.N_TS_IDX_SIG:
		n := node.(*parser.TsIdxSig)
		if wt, ok := n.Key().(parser.NodeWithTypInfo); ok {
			ti := wt.TypInfo()
			if ti == nil {
				ti = parser.NewTypInfo()
				wt.SetTypInfo(ti)
			}
			ti.SetTypAnnot(n.KeyType())
		}
		return &TSIndexSignature{
			Type:           "TSIndexSignature",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Parameters:     elems([]parser.Node{n.Key()}),
			TypeAnnotation: convertTsTyp(n.Value()),
		}
	case parser.N_TS_NS_NAME:
		n := node.(*parser.TsNsName)
		return &TSQualifiedName{
			Type:  "TSQualifiedName",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Left:  convertTsTyp(n.Lhs()),
			Right: convertTsTyp(n.Rhs()),
		}
	case parser.N_TS_DEC_VAR_DEC:
		n := node.(*parser.TsDec)
		varDec := n.Inner().(*parser.VarDecStmt)
		return &TSVariableDeclaration{
			Type:         "VariableDeclaration",
			Start:        start(n.Loc()),
			End:          end(n.Loc()),
			Loc:          loc(n.Loc()),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList()),
			Declare:      true,
		}
	case parser.N_TS_INTERFACE_BODY:
		n := node.(*parser.TsInferfaceBody)
		return &TSInterfaceBody{
			Type:  "TSInterfaceBody",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Body:  elems(n.Body()),
		}
	case parser.N_TS_DEC_INTERFACE:
		n := node.(*parser.TsDec)
		itf := n.Inner().(*parser.TsInferface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             convert(itf.Id()),
			TypeParameters: convertTsTyp(itf.TypParams()),
			Extends:        extends(itf.Supers()),
			Body:           convert(itf.Body()),
			Declare:        true,
		}
	case parser.N_TS_INTERFACE:
		n := node.(*parser.TsInferface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             convert(n.Id()),
			TypeParameters: convertTsTyp(n.TypParams()),
			Extends:        extends(n.Supers()),
			Body:           convert(n.Body()),
			Declare:        false,
		}
	}

	return nil
}

func extends(exts []parser.Node) []Node {
	if exts == nil {
		return nil
	}
	ret := make([]Node, len(exts))
	for i, ext := range exts {
		ret[i] = exprWithTypArg(ext)
	}
	return ret
}

func exprWithTypArg(node parser.Node) Node {
	var typParams Node
	if wt, ok := node.(parser.NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti != nil {
			typParams = typAnnot(ti)
		}
	}
	lc := parser.LocWithTypeInfo(node, false)
	return &TSExpressionWithTypeArguments{
		Type:           "TSExpressionWithTypeArguments",
		Start:          start(lc),
		End:            end(lc),
		Loc:            loc(lc),
		Expression:     convert(node),
		TypeParameters: typParams,
	}
}

func typAnnot(ti *parser.TypInfo) Node {
	ta := ti.TypAnnot()
	if ta == nil {
		return nil
	}
	return convertTsTyp(ta)
}

func optional(ti *parser.TypInfo) bool {
	if ti == nil {
		return false
	}
	return ti.Optional()
}

func typParams(ti *parser.TypInfo) Node {
	psDec := ti.TypParams()
	if psDec == nil {
		return nil
	}

	ps := psDec.(*parser.TsParamsDec).Params()

	psLen := len(ps)
	if psLen == 0 {
		return nil
	}

	ret := make([]Node, len(ps))
	for i, p := range ps {
		ret[i] = convertTsTyp(p)
	}

	return &TSTypeParameterDeclaration{
		Type:   "TSTypeParameterDeclaration",
		Start:  start(psDec.Loc()),
		End:    end(psDec.Loc()),
		Loc:    loc(psDec.Loc()),
		Params: ret,
	}
}

func typArgs(ti *parser.TypInfo) Node {
	psInst := ti.TypArgs()
	if psInst == nil {
		return nil
	}

	ps := psInst.(*parser.TsParamsInst).Params()

	psLen := len(ps)
	if psLen == 0 {
		return nil
	}

	ret := make([]Node, len(ps))
	for i, p := range ps {
		ret[i] = convertTsTyp(p)
	}

	return &TSTypeParameterInstantiation{
		Type:   "TSTypeParameterInstantiation",
		Start:  start(psInst.Loc()),
		End:    end(psInst.Loc()),
		Loc:    loc(psInst.Loc()),
		Params: ret,
	}
}

func tsParamProp(node parser.Node) Node {
	ti, ok := isTyParamProp(node)
	if !ok {
		return nil
	}

	lc := parser.LocWithTypeInfo(node, true)
	return &TSParameterProperty{
		Type:          "TSParameterProperty",
		Start:         start(lc),
		End:           end(lc),
		Loc:           loc(lc),
		Parameter:     convert(node),
		Readonly:      ti.Readonly(),
		Override:      ti.Override(),
		Accessibility: ti.AccMod().String(),
	}
}

func isTyParamProp(node parser.Node) (*parser.TypInfo, bool) {
	wt, ok := node.(parser.NodeWithTypInfo)
	if !ok {
		return nil, false
	}
	ti := wt.TypInfo()
	if ti == nil {
		return nil, false
	}
	return ti, ti.Readonly() || ti.Override() || ti.AccMod() != parser.ACC_MOD_NONE
}

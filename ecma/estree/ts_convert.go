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
	case parser.N_TS_REF:
		n := node.(*parser.TsRef)
		return &TSTypeReference{
			Type:     "TSTypeReference",
			Start:    start(node.Loc()),
			End:      end(node.Loc()),
			Loc:      loc(node.Loc()),
			TypeName: convert(n.Name()),
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
		lc := parser.LocWithTypeInfo(node)
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
	}

	return nil
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

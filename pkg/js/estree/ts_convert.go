package estree

import "github.com/hsiaosiyuan0/mole/pkg/js/parser"

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

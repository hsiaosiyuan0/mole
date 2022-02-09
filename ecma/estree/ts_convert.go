package estree

import "github.com/hsiaosiyuan0/mole/ecma/parser"

func ConvertTsTyp(node parser.Node, ctx *ConvertCtx) Node {
	if node == nil {
		return nil
	}

	switch node.Type() {
	case parser.N_NAME:
		return ident(node, ctx)
	case parser.N_TS_TYP_ANNOT:
		n := node.(*parser.TsTypAnnot)
		return &TSTypeAnnotation{
			Type:           "TSTypeAnnotation",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			TypeAnnotation: ConvertTsTyp(n.TsTyp(), ctx),
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
	case parser.N_TS_INTRINSIC:
		return &TSIntrinsicKeyword{
			Type:  "TSIntrinsicKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_NEVER:
		return &TSNeverKeyword{
			Type:  "TSNeverKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_SYM:
		return &TSSymbolKeyword{
			Type:  "TSSymbolKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_UNDEF:
		return &TSUndefinedKeyword{
			Type:  "TSUndefinedKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_BIGINT:
		return &TSBigIntKeyword{
			Type:  "TSBigIntKeyword",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
		}
	case parser.N_TS_NULL:
		return &TSNullKeyword{
			Type:  "TSNullKeyword",
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
			TypeName:       Convert(n.Name(), ctx),
			TypeParameters: Convert(n.ParamsInst(), ctx),
		}
	case parser.N_TS_PARAM_INST:
		n := node.(*parser.TsParamsInst)
		return &TSTypeParameterInstantiation{
			Type:   "TSTypeParameterInstantiation",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Params: elems(n.Params(), ctx),
		}
	case parser.N_TS_PARAM_DEC:
		n := node.(*parser.TsParamsDec)
		return &TSTypeParameterDeclaration{
			Type:   "TSTypeParameterDeclaration",
			Start:  start(node.Loc()),
			End:    end(node.Loc()),
			Loc:    loc(node.Loc()),
			Params: elems(n.Params(), ctx),
		}
	case parser.N_TS_PARAM:
		n := node.(*parser.TsParam)
		return &TSTypeParameter{
			Type:       "TSTypeParameter",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Name:       Convert(n.Name(), ctx),
			Constraint: Convert(n.Cons(), ctx),
			Default:    Convert(n.Default(), ctx),
		}
	case parser.N_TS_ARR:
		n := node.(*parser.TsArr)
		return &TSArrayType{
			Type:        "TSArrayType",
			Start:       start(node.Loc()),
			End:         end(node.Loc()),
			Loc:         loc(node.Loc()),
			ElementType: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_LIT_OBJ:
		n := node.(*parser.TsObj)
		return &TSTypeLiteral{
			Type:    "TSTypeLiteral",
			Start:   start(node.Loc()),
			End:     end(node.Loc()),
			Loc:     loc(node.Loc()),
			Members: elems(n.Props(), ctx),
		}
	case parser.N_TS_CALL_SIG:
		n := node.(*parser.TsCallSig)
		return &TSCallSignatureDeclaration{
			Type:           "TSCallSignatureDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
		}
	case parser.N_TS_PROP:
		n := node.(*parser.TsProp)
		if n.IsMethod() {
			return &TSMethodSignature{
				Type:     "TSMethodSignature",
				Start:    start(node.Loc()),
				End:      end(node.Loc()),
				Loc:      loc(node.Loc()),
				Key:      Convert(n.Key(), ctx),
				Value:    Convert(n.Val(), ctx),
				Computed: n.Computed(),
				Optional: n.Optional(),
				Kind:     n.Kind().ToString(),
			}
		}
		return &TSPropertySignature{
			Type:           "TSPropertySignature",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			Key:            Convert(n.Key(), ctx),
			Optional:       n.Optional(),
			Computed:       n.Computed(),
			TypeAnnotation: ConvertTsTyp(n.Val(), ctx),
			Kind:           n.Kind().ToString(),
			Readonly:       n.Readonly(),
		}
	case parser.N_TS_TYP_PREDICATE:
		n := node.(*parser.TsTypPredicate)
		return &TSTypePredicate{
			Type:           "TSTypePredicate",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			ParameterName:  Convert(n.Name(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Typ(), ctx),
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
			Id:             Convert(n.Id(), ctx),
			Params:         fnParams(n.Params(), ctx),
			Body:           Convert(n.Body(), ctx),
			Generator:      false,
			Async:          n.Async(),
			TypeParameters: typParams(ti, ctx),
			ReturnType:     typAnnot(ti, ctx),
		}
	case parser.N_TS_TYP_ASSERT:
		n := node.(*parser.TsTypAssert)
		return &TSTypeAssertion{
			Type:           "TSTypeAssertion",
			Start:          start(node.Loc()),
			End:            end(node.Loc()),
			Loc:            loc(node.Loc()),
			Expression:     Convert(n.Expr(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Typ(), ctx),
		}
	case parser.N_TS_NO_NULL:
		n := node.(*parser.TsNoNull)
		return &TSNonNullExpression{
			Type:       "TSNonNullExpression",
			Start:      start(node.Loc()),
			End:        end(node.Loc()),
			Loc:        loc(node.Loc()),
			Expression: Convert(n.Arg(), ctx),
		}
	case parser.N_TS_UNION_TYP:
		n := node.(*parser.TsUnionTyp)
		return &TSUnionType{
			Type:  "TSUnionType",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Types: elems(n.Elems(), ctx),
		}
	case parser.N_TS_INTERSEC_TYP:
		n := node.(*parser.TsIntersecTyp)
		return &TSIntersectionType{
			Type:  "TSIntersectionType",
			Start: start(node.Loc()),
			End:   end(node.Loc()),
			Loc:   loc(node.Loc()),
			Types: elems(n.Elems(), ctx),
		}
	case parser.N_TS_DEC_CLASS:
		n := node.(*parser.TsDec)
		cls := n.Inner().(*parser.ClassDec)
		return &TSClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      start(n.Loc()),
			End:        end(cls.Loc()),
			Loc:        loc(n.Loc()),
			Id:         Convert(cls.Id(), ctx),
			SuperClass: Convert(cls.Super(), ctx),
			Body:       Convert(cls.Body(), ctx),
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
			Parameters:     elems([]parser.Node{n.Key()}, ctx),
			TypeAnnotation: ConvertTsTyp(n.Val(), ctx),
		}
	case parser.N_TS_NS_NAME:
		n := node.(*parser.TsNsName)
		return &TSQualifiedName{
			Type:  "TSQualifiedName",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Left:  ConvertTsTyp(n.Lhs(), ctx),
			Right: ConvertTsTyp(n.Rhs(), ctx),
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
			Declarations: declarations(varDec.DecList(), ctx),
			Declare:      true,
		}
	case parser.N_TS_DEC_INTERFACE:
		n := node.(*parser.TsDec)
		itf := n.Inner().(*parser.TsInferface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             Convert(itf.Id(), ctx),
			TypeParameters: ConvertTsTyp(itf.TypParams(), ctx),
			Extends:        extends(itf.Supers(), ctx),
			Body:           Convert(itf.Body(), ctx),
			Declare:        true,
		}
	case parser.N_TS_INTERFACE_BODY:
		n := node.(*parser.TsInferfaceBody)
		scope := ctx.enter()
		defer ctx.leave()
		scope.Flag = scope.Flag.On(CSF_INTERFACE)

		return &TSInterfaceBody{
			Type:  "TSInterfaceBody",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Body:  elems(n.Body(), ctx),
		}
	case parser.N_TS_INTERFACE:
		n := node.(*parser.TsInferface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             Convert(n.Id(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			Extends:        extends(n.Supers(), ctx),
			Body:           Convert(n.Body(), ctx),
			Declare:        false,
		}
	case parser.N_TS_ENUM:
		n := node.(*parser.TsEnum)
		return &TSEnumDeclaration{
			Type:    "TSEnumDeclaration",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Id:      Convert(n.Id(), ctx),
			Const:   n.Const(),
			Members: elems(n.Members(), ctx),
		}
	case parser.N_TS_DEC_ENUM:
		n := node.(*parser.TsDec)
		enum := n.Inner().(*parser.TsEnum)
		return &TSEnumDeclaration{
			Type:    "TSEnumDeclaration",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Id:      Convert(enum.Id(), ctx),
			Const:   enum.Const(),
			Members: elems(enum.Members(), ctx),
			Declare: true,
		}
	case parser.N_TS_ENUM_MEMBER:
		n := node.(*parser.TsEnumMember)
		return &TSEnumMember{
			Type:        "TSEnumMember",
			Start:       start(n.Loc()),
			End:         end(n.Loc()),
			Loc:         loc(n.Loc()),
			Id:          Convert(n.Key(), ctx),
			Initializer: Convert(n.Val(), ctx),
		}
	case parser.N_TS_TYP_DEC:
		n := node.(*parser.TsTypDec)
		return &TSTypeAliasDeclaration{
			Type:           "TSTypeAliasDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             Convert(n.Id(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			TypeAnnotation: typAnnot(n.TypInfo(), ctx),
			Declare:        false,
		}
	case parser.N_TS_DEC_TYP_DEC:
		n := node.(*parser.TsDec)
		dec := n.Inner().(*parser.TsTypDec)
		return &TSTypeAliasDeclaration{
			Type:           "TSTypeAliasDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Id:             Convert(dec.Id(), ctx),
			TypeAnnotation: typAnnot(dec.TypInfo(), ctx),
			Declare:        true,
		}
	case parser.N_TS_DEC_MODULE, parser.N_TS_DEC_GLOBAL:
		n := node.(*parser.TsDec)
		return &TSModuleDeclaration{
			Type:    "TSModuleDeclaration",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Id:      Convert(n.Name(), ctx),
			Body:    Convert(n.Inner(), ctx),
			Declare: true,
			Global:  n.Type() == parser.N_TS_DEC_GLOBAL,
		}
	case parser.N_TS_DEC_NS:
		n := node.(*parser.TsDec)
		ns := n.Inner().(*parser.TsNS)
		return &TSModuleDeclaration{
			Type:    "TSModuleDeclaration",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Id:      Convert(ns.Id(), ctx),
			Body:    Convert(ns.Body(), ctx),
			Declare: true,
		}
	case parser.N_TS_NAMESPACE:
		n := node.(*parser.TsNS)
		return &TSModuleDeclaration{
			Type:    "TSModuleDeclaration",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Id:      Convert(n.Id(), ctx),
			Body:    Convert(n.Body(), ctx),
			Declare: false,
		}
	case parser.N_TS_EXPORT_ASSIGN:
		n := node.(*parser.TsExportAssign)
		return &TSExportAssignment{
			Type:       "TSExportAssignment",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			Expression: Convert(n.Expr(), ctx),
		}
	case parser.N_TS_LIT:
		n := node.(*parser.TsLit)
		return &TSLiteralType{
			Type:    "TSLiteralType",
			Start:   start(n.Loc()),
			End:     end(n.Loc()),
			Loc:     loc(n.Loc()),
			Literal: Convert(n.Lit(), ctx),
		}
	case parser.N_TS_IMPORT_ALIAS:
		n := node.(*parser.TsImportAlias)
		return &TSImportEqualsDeclaration{
			Type:            "TSImportEqualsDeclaration",
			Start:           start(n.Loc()),
			End:             end(n.Loc()),
			Loc:             loc(n.Loc()),
			Id:              Convert(n.Name(), ctx),
			ModuleReference: Convert(n.Val(), ctx),
			IsExport:        n.Export(),
		}
	case parser.N_TS_IMPORT_REQUIRE:
		n := node.(*parser.TsImportRequire)
		expr := n.Expr().(*parser.CallExpr)
		return &TSImportEqualsDeclaration{
			Type:  "TSImportEqualsDeclaration",
			Start: start(n.Loc()),
			End:   end(n.Loc()),
			Loc:   loc(n.Loc()),
			Id:    Convert(n.Name(), ctx),
			ModuleReference: &TSExternalModuleReference{
				Type:       "TSExternalModuleReference",
				Start:      start(expr.Loc()),
				End:        end(expr.Loc()),
				Loc:        loc(expr.Loc()),
				Expression: Convert(expr.Args()[0], ctx),
			},
			IsExport: false,
		}
	case parser.N_TS_NEW_SIG:
		n := node.(*parser.TsNewSig)
		return &TSConstructSignatureDeclaration{
			Type:           "TSConstructSignatureDeclaration",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
			Abstract:       n.Abstract(),
		}
	case parser.N_TS_NEW:
		n := node.(*parser.TsNewSig)
		return &TSConstructorType{
			Type:           "TSConstructorType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
			Abstract:       n.Abstract(),
		}
	case parser.N_TS_FN_TYP:
		n := node.(*parser.TsFnTyp)
		return &TSFunctionType{
			Type:           "TSFunctionType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
		}
	case parser.N_TS_IMPORT_TYP:
		n := node.(*parser.TsImportType)
		return &TSImportType{
			Type:           "TSImportType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Argument:       Convert(n.Arg(), ctx),
			Qualifier:      ConvertTsTyp(n.Qualifier(), ctx),
			TypeParameters: ConvertTsTyp(n.TypArg(), ctx),
		}
	case parser.N_TS_TYP_QUERY:
		n := node.(*parser.TsTypQuery)
		return &TSTypeQuery{
			Type:     "TSTypeQuery",
			Start:    start(n.Loc()),
			End:      end(n.Loc()),
			Loc:      loc(n.Loc()),
			ExprName: Convert(n.Arg(), ctx),
		}
	case parser.N_TS_COND:
		n := node.(*parser.TsCondType)
		return &TSConditionalType{
			Type:        "TSConditionalType",
			Start:       start(n.Loc()),
			End:         end(n.Loc()),
			Loc:         loc(n.Loc()),
			CheckType:   Convert(n.CheckTyp(), ctx),
			ExtendsType: Convert(n.ExtTyp(), ctx),
			TrueType:    Convert(n.TrueTyp(), ctx),
			FalseType:   Convert(n.FalseTyp(), ctx),
		}
	case parser.N_TS_TYP_INFER:
		n := node.(*parser.TsTypInfer)
		return &TSInferType{
			Type:          "TSInferType",
			Start:         start(n.Loc()),
			End:           end(n.Loc()),
			Loc:           loc(n.Loc()),
			TypeParameter: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_PAREN:
		n := node.(*parser.TsParen)
		return &TSParenthesizedType{
			Type:           "TSParenthesizedType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_IDX_ACCESS:
		n := node.(*parser.TsIdxAccess)
		return &TSIndexedAccessType{
			Type:       "TSIndexedAccessType",
			Start:      start(n.Loc()),
			End:        end(n.Loc()),
			Loc:        loc(n.Loc()),
			ObjectType: ConvertTsTyp(n.Obj(), ctx),
			IndexType:  ConvertTsTyp(n.Idx(), ctx),
		}
	case parser.N_TS_MAPPED:
		n := node.(*parser.TsMapped)
		return &TSMappedType{
			Type:           "TSMappedType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Readonly:       n.ReadonlyFmt(),
			Optional:       n.OptionalFmt(),
			TypeParameter:  ConvertTsTyp(n.Key(), ctx),
			NameType:       ConvertTsTyp(n.Name(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Val(), ctx),
		}
	case parser.N_TS_TYP_OP:
		n := node.(*parser.TsTypOp)
		return &TSTypeOperator{
			Type:           "TSTypeOperator",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			Operator:       n.Op(),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_TUPLE:
		n := node.(*parser.TsTuple)
		return &TSTupleType{
			Type:         "TSTupleType",
			Start:        start(n.Loc()),
			End:          end(n.Loc()),
			Loc:          loc(n.Loc()),
			ElementTypes: elems(n.Args(), ctx),
		}
	case parser.N_TS_TUPLE_NAMED_MEMBER:
		n := node.(*parser.TsTupleNamedMember)
		return &TSNamedTupleMember{
			Type:        "TSNamedTupleMember",
			Start:       start(n.Loc()),
			End:         end(n.Loc()),
			Loc:         loc(n.Loc()),
			Optional:    n.Opt(),
			Label:       Convert(n.Label(), ctx),
			ElementType: ConvertTsTyp(n.Val(), ctx),
		}
	case parser.N_TS_REST:
		n := node.(*parser.TsRest)
		return &TSRestType{
			Type:           "TSRestType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_OPT:
		n := node.(*parser.TsOpt)
		return &TSOptionalType{
			Type:           "TSOptionalType",
			Start:          start(n.Loc()),
			End:            end(n.Loc()),
			Loc:            loc(n.Loc()),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	}

	return nil
}

func extends(exts []parser.Node, ctx *ConvertCtx) []Node {
	if exts == nil {
		return nil
	}
	ret := make([]Node, len(exts))
	for i, ext := range exts {
		ret[i] = exprWithTypArg(ext, ctx)
	}
	return ret
}

func exprWithTypArg(node parser.Node, ctx *ConvertCtx) Node {
	expr := node
	var typParams Node
	if node.Type() == parser.N_TS_REF {
		n := node.(*parser.TsRef)
		expr = n.Name()
		typParams = ConvertTsTyp(n.ParamsInst(), ctx)
	} else if wt, ok := node.(parser.NodeWithTypInfo); ok {
		ti := wt.TypInfo()
		if ti != nil {
			typParams = typAnnot(ti, ctx)
		}
	}
	lc := parser.LocWithTypeInfo(node, false)
	return &TSExpressionWithTypeArguments{
		Type:           "TSExpressionWithTypeArguments",
		Start:          start(lc),
		End:            end(lc),
		Loc:            loc(lc),
		Expression:     Convert(expr, ctx),
		TypeParameters: typParams,
	}
}

func typAnnot(ti *parser.TypInfo, ctx *ConvertCtx) Node {
	ta := ti.TypAnnot()
	if ta == nil {
		return nil
	}
	return ConvertTsTyp(ta, ctx)
}

func optional(ti *parser.TypInfo) bool {
	if ti == nil {
		return false
	}
	return ti.Optional()
}

func typParams(ti *parser.TypInfo, ctx *ConvertCtx) Node {
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
		ret[i] = ConvertTsTyp(p, ctx)
	}

	return &TSTypeParameterDeclaration{
		Type:   "TSTypeParameterDeclaration",
		Start:  start(psDec.Loc()),
		End:    end(psDec.Loc()),
		Loc:    loc(psDec.Loc()),
		Params: ret,
	}
}

func typArgs(ti *parser.TypInfo, ctx *ConvertCtx) Node {
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
		ret[i] = ConvertTsTyp(p, ctx)
	}

	return &TSTypeParameterInstantiation{
		Type:   "TSTypeParameterInstantiation",
		Start:  start(psInst.Loc()),
		End:    end(psInst.Loc()),
		Loc:    loc(psInst.Loc()),
		Params: ret,
	}
}

func tsParamProp(node parser.Node, ctx *ConvertCtx) Node {
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
		Parameter:     Convert(node, ctx),
		Readonly:      ti.Readonly(),
		Override:      ti.Override(),
		Accessibility: ti.AccMod().String(),
		Decorators:    elems(ti.Decorators(), ctx),
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

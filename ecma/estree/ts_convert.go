package estree

import (
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

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
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
			TypeAnnotation: ConvertTsTyp(n.TsTyp(), ctx),
		}
	case parser.N_TS_NUM:
		return &TSNumberKeyword{
			Type:  "TSNumberKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_STR:
		return &TSStringKeyword{
			Type:  "TSStringKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_ANY:
		return &TSAnyKeyword{
			Type:  "TSAnyKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_BOOL:
		return &TSBooleanKeyword{
			Type:  "TSBooleanKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_VOID:
		return &TSVoidKeyword{
			Type:  "TSVoidKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_INTRINSIC:
		return &TSIntrinsicKeyword{
			Type:  "TSIntrinsicKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_NEVER:
		return &TSNeverKeyword{
			Type:  "TSNeverKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_SYM:
		return &TSSymbolKeyword{
			Type:  "TSSymbolKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_UNDEF:
		return &TSUndefinedKeyword{
			Type:  "TSUndefinedKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_BIGINT:
		return &TSBigIntKeyword{
			Type:  "TSBigIntKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_NULL:
		return &TSNullKeyword{
			Type:  "TSNullKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_THIS:
		return &TSThisType{
			Type:  "TSThisType",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_UNKNOWN:
		return &TSUnknownKeyword{
			Type:  "TSUnknownKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_OBJ:
		return &TSObjectKeyword{
			Type:  "TSObjectKeyword",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
		}
	case parser.N_TS_REF:
		n := node.(*parser.TsRef)
		return &TSTypeReference{
			Type:           "TSTypeReference",
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
			TypeName:       Convert(n.Name(), ctx),
			TypeParameters: Convert(n.ParamsInst(), ctx),
		}
	case parser.N_TS_PARAM_INST:
		n := node.(*parser.TsParamsInst)
		return &TSTypeParameterInstantiation{
			Type:   "TSTypeParameterInstantiation",
			Start:  int(node.Range().Lo),
			End:    int(node.Range().Hi),
			Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
			Params: elems(n.Params(), ctx),
		}
	case parser.N_TS_PARAM_DEC:
		n := node.(*parser.TsParamsDec)
		return &TSTypeParameterDeclaration{
			Type:   "TSTypeParameterDeclaration",
			Start:  int(node.Range().Lo),
			End:    int(node.Range().Hi),
			Loc:    locOfNode(node, ctx.Parser.Source(), ctx),
			Params: elems(n.Params(), ctx),
		}
	case parser.N_TS_PARAM:
		n := node.(*parser.TsParam)
		return &TSTypeParameter{
			Type:       "TSTypeParameter",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Name:       Convert(n.Name(), ctx),
			Constraint: Convert(n.Cons(), ctx),
			Default:    Convert(n.Default(), ctx),
		}
	case parser.N_TS_ARR:
		n := node.(*parser.TsArr)
		return &TSArrayType{
			Type:        "TSArrayType",
			Start:       int(node.Range().Lo),
			End:         int(node.Range().Hi),
			Loc:         locOfNode(node, ctx.Parser.Source(), ctx),
			ElementType: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_LIT_OBJ:
		n := node.(*parser.TsObj)
		return &TSTypeLiteral{
			Type:    "TSTypeLiteral",
			Start:   int(node.Range().Lo),
			End:     int(node.Range().Hi),
			Loc:     locOfNode(node, ctx.Parser.Source(), ctx),
			Members: elems(n.Props(), ctx),
		}
	case parser.N_TS_CALL_SIG:
		n := node.(*parser.TsCallSig)
		return &TSCallSignatureDeclaration{
			Type:           "TSCallSignatureDeclaration",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
		}
	case parser.N_TS_PROP:
		n := node.(*parser.TsProp)
		if n.IsMethod() {
			return &TSMethodSignature{
				Type:     "TSMethodSignature",
				Start:    int(node.Range().Lo),
				End:      int(node.Range().Hi),
				Loc:      locOfNode(node, ctx.Parser.Source(), ctx),
				Key:      Convert(n.Key(), ctx),
				Value:    Convert(n.Val(), ctx),
				Computed: n.Computed(),
				Optional: n.Optional(),
				Kind:     n.Kind().ToString(),
			}
		}
		return &TSPropertySignature{
			Type:           "TSPropertySignature",
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
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
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
			ParameterName:  Convert(n.Name(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Typ(), ctx),
			Asserts:        n.Asserts(),
		}
	case parser.N_TS_DEC_FN:
		n := node.(*parser.TsDec).Inner().(*parser.FnDec)
		ti := n.TypInfo()
		rng, loc := locWithTypeInfo(node, false, ctx.Parser.Source(), ctx)
		return &TSDeclareFunction{
			Type:           "TSDeclareFunction",
			Start:          int(rng.Lo),
			End:            int(rng.Hi),
			Loc:            loc,
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
			Start:          int(node.Range().Lo),
			End:            int(node.Range().Hi),
			Loc:            locOfNode(node, ctx.Parser.Source(), ctx),
			Expression:     Convert(n.Expr(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Typ(), ctx),
		}
	case parser.N_TS_NO_NULL:
		n := node.(*parser.TsNoNull)
		return &TSNonNullExpression{
			Type:       "TSNonNullExpression",
			Start:      int(node.Range().Lo),
			End:        int(node.Range().Hi),
			Loc:        locOfNode(node, ctx.Parser.Source(), ctx),
			Expression: Convert(n.Arg(), ctx),
		}
	case parser.N_TS_UNION_TYP:
		n := node.(*parser.TsUnionTyp)
		return &TSUnionType{
			Type:  "TSUnionType",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Types: elems(n.Elems(), ctx),
		}
	case parser.N_TS_INTERSECT_TYP:
		n := node.(*parser.TsIntersectTyp)
		return &TSIntersectionType{
			Type:  "TSIntersectionType",
			Start: int(node.Range().Lo),
			End:   int(node.Range().Hi),
			Loc:   locOfNode(node, ctx.Parser.Source(), ctx),
			Types: elems(n.Elems(), ctx),
		}
	case parser.N_TS_DEC_CLASS:
		n := node.(*parser.TsDec)
		cls := n.Inner().(*parser.ClassDec)
		return &TSClassDeclaration{
			Type:       "ClassDeclaration",
			Start:      int(n.Range().Lo),
			End:        int(cls.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Parameters:     elems([]parser.Node{n.Key()}, ctx),
			TypeAnnotation: ConvertTsTyp(n.Val(), ctx),
		}
	case parser.N_TS_NS_NAME:
		n := node.(*parser.TsNsName)
		return &TSQualifiedName{
			Type:  "TSQualifiedName",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
			Left:  ConvertTsTyp(n.Lhs(), ctx),
			Right: ConvertTsTyp(n.Rhs(), ctx),
		}
	case parser.N_TS_DEC_VAR_DEC:
		n := node.(*parser.TsDec)
		varDec := n.Inner().(*parser.VarDecStmt)
		return &TSVariableDeclaration{
			Type:         "VariableDeclaration",
			Start:        int(n.Range().Lo),
			End:          int(n.Range().Hi),
			Loc:          locOfNode(n, ctx.Parser.Source(), ctx),
			Kind:         varDec.Kind(),
			Declarations: declarations(varDec.DecList(), ctx),
			Declare:      true,
		}
	case parser.N_TS_DEC_INTERFACE:
		n := node.(*parser.TsDec)
		itf := n.Inner().(*parser.TsInterface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Id:             Convert(itf.Id(), ctx),
			TypeParameters: ConvertTsTyp(itf.TypParams(), ctx),
			Extends:        extends(itf.Supers(), ctx),
			Body:           Convert(itf.Body(), ctx),
			Declare:        true,
		}
	case parser.N_TS_INTERFACE_BODY:
		n := node.(*parser.TsInterfaceBody)
		scope := ctx.enter()
		defer ctx.leave()
		scope.Flag = scope.Flag.On(CSF_INTERFACE)

		return &TSInterfaceBody{
			Type:  "TSInterfaceBody",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
			Body:  elems(n.Body(), ctx),
		}
	case parser.N_TS_INTERFACE:
		n := node.(*parser.TsInterface)
		return &TSInterfaceDeclaration{
			Type:           "TSInterfaceDeclaration",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
			Id:      Convert(n.Id(), ctx),
			Const:   n.Const(),
			Members: elems(n.Members(), ctx),
		}
	case parser.N_TS_DEC_ENUM:
		n := node.(*parser.TsDec)
		enum := n.Inner().(*parser.TsEnum)
		return &TSEnumDeclaration{
			Type:    "TSEnumDeclaration",
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
			Id:      Convert(enum.Id(), ctx),
			Const:   enum.Const(),
			Members: elems(enum.Members(), ctx),
			Declare: true,
		}
	case parser.N_TS_ENUM_MEMBER:
		n := node.(*parser.TsEnumMember)
		return &TSEnumMember{
			Type:        "TSEnumMember",
			Start:       int(n.Range().Lo),
			End:         int(n.Range().Hi),
			Loc:         locOfNode(n, ctx.Parser.Source(), ctx),
			Id:          Convert(n.Key(), ctx),
			Initializer: Convert(n.Val(), ctx),
		}
	case parser.N_TS_TYP_DEC:
		n := node.(*parser.TsTypDec)
		return &TSTypeAliasDeclaration{
			Type:           "TSTypeAliasDeclaration",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Id:             Convert(dec.Id(), ctx),
			TypeAnnotation: typAnnot(dec.TypInfo(), ctx),
			Declare:        true,
		}
	case parser.N_TS_DEC_MODULE, parser.N_TS_DEC_GLOBAL:
		n := node.(*parser.TsDec)
		return &TSModuleDeclaration{
			Type:    "TSModuleDeclaration",
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
			Id:      Convert(ns.Id(), ctx),
			Body:    Convert(ns.Body(), ctx),
			Declare: true,
		}
	case parser.N_TS_NAMESPACE:
		n := node.(*parser.TsNS)
		return &TSModuleDeclaration{
			Type:    "TSModuleDeclaration",
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
			Id:      Convert(n.Id(), ctx),
			Body:    Convert(n.Body(), ctx),
			Declare: false,
		}
	case parser.N_TS_EXPORT_ASSIGN:
		n := node.(*parser.TsExportAssign)
		return &TSExportAssignment{
			Type:       "TSExportAssignment",
			Start:      int(n.Range().Lo),
			End:        int(n.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
			Expression: Convert(n.Expr(), ctx),
		}
	case parser.N_TS_LIT:
		n := node.(*parser.TsLit)
		return &TSLiteralType{
			Type:    "TSLiteralType",
			Start:   int(n.Range().Lo),
			End:     int(n.Range().Hi),
			Loc:     locOfNode(n, ctx.Parser.Source(), ctx),
			Literal: Convert(n.Lit(), ctx),
		}
	case parser.N_TS_IMPORT_ALIAS:
		n := node.(*parser.TsImportAlias)
		return &TSImportEqualsDeclaration{
			Type:            "TSImportEqualsDeclaration",
			Start:           int(n.Range().Lo),
			End:             int(n.Range().Hi),
			Loc:             locOfNode(n, ctx.Parser.Source(), ctx),
			Id:              Convert(n.Name(), ctx),
			ModuleReference: Convert(n.Val(), ctx),
			IsExport:        n.Export(),
		}
	case parser.N_TS_IMPORT_REQUIRE:
		n := node.(*parser.TsImportRequire)
		expr := n.Expr().(*parser.CallExpr)
		return &TSImportEqualsDeclaration{
			Type:  "TSImportEqualsDeclaration",
			Start: int(n.Range().Lo),
			End:   int(n.Range().Hi),
			Loc:   locOfNode(n, ctx.Parser.Source(), ctx),
			Id:    Convert(n.Name(), ctx),
			ModuleReference: &TSExternalModuleReference{
				Type:       "TSExternalModuleReference",
				Start:      int(expr.Range().Lo),
				End:        int(expr.Range().Hi),
				Loc:        locOfNode(expr, ctx.Parser.Source(), ctx),
				Expression: Convert(expr.Args()[0], ctx),
			},
			IsExport: false,
		}
	case parser.N_TS_NEW_SIG:
		n := node.(*parser.TsNewSig)
		return &TSConstructSignatureDeclaration{
			Type:           "TSConstructSignatureDeclaration",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
			Abstract:       n.Abstract(),
		}
	case parser.N_TS_NEW:
		n := node.(*parser.TsNewSig)
		return &TSConstructorType{
			Type:           "TSConstructorType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
			Abstract:       n.Abstract(),
		}
	case parser.N_TS_FN_TYP:
		n := node.(*parser.TsFnTyp)
		return &TSFunctionType{
			Type:           "TSFunctionType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Params:         fnParams(n.Params(), ctx),
			TypeParameters: ConvertTsTyp(n.TypParams(), ctx),
			ReturnType:     ConvertTsTyp(n.RetTyp(), ctx),
		}
	case parser.N_TS_IMPORT_TYP:
		n := node.(*parser.TsImportType)
		return &TSImportType{
			Type:           "TSImportType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Argument:       Convert(n.Arg(), ctx),
			Qualifier:      ConvertTsTyp(n.Qualifier(), ctx),
			TypeParameters: ConvertTsTyp(n.TypArg(), ctx),
		}
	case parser.N_TS_TYP_QUERY:
		n := node.(*parser.TsTypQuery)
		return &TSTypeQuery{
			Type:     "TSTypeQuery",
			Start:    int(n.Range().Lo),
			End:      int(n.Range().Hi),
			Loc:      locOfNode(n, ctx.Parser.Source(), ctx),
			ExprName: Convert(n.Arg(), ctx),
		}
	case parser.N_TS_COND:
		n := node.(*parser.TsCondType)
		return &TSConditionalType{
			Type:        "TSConditionalType",
			Start:       int(n.Range().Lo),
			End:         int(n.Range().Hi),
			Loc:         locOfNode(n, ctx.Parser.Source(), ctx),
			CheckType:   Convert(n.CheckTyp(), ctx),
			ExtendsType: Convert(n.ExtTyp(), ctx),
			TrueType:    Convert(n.TrueTyp(), ctx),
			FalseType:   Convert(n.FalseTyp(), ctx),
		}
	case parser.N_TS_TYP_INFER:
		n := node.(*parser.TsTypInfer)
		return &TSInferType{
			Type:          "TSInferType",
			Start:         int(n.Range().Lo),
			End:           int(n.Range().Hi),
			Loc:           locOfNode(n, ctx.Parser.Source(), ctx),
			TypeParameter: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_PAREN:
		n := node.(*parser.TsParen)
		return &TSParenthesizedType{
			Type:           "TSParenthesizedType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_IDX_ACCESS:
		n := node.(*parser.TsIdxAccess)
		return &TSIndexedAccessType{
			Type:       "TSIndexedAccessType",
			Start:      int(n.Range().Lo),
			End:        int(n.Range().Hi),
			Loc:        locOfNode(n, ctx.Parser.Source(), ctx),
			ObjectType: ConvertTsTyp(n.Obj(), ctx),
			IndexType:  ConvertTsTyp(n.Idx(), ctx),
		}
	case parser.N_TS_MAPPED:
		n := node.(*parser.TsMapped)
		return &TSMappedType{
			Type:           "TSMappedType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			Operator:       n.Op(),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_TUPLE:
		n := node.(*parser.TsTuple)
		return &TSTupleType{
			Type:         "TSTupleType",
			Start:        int(n.Range().Lo),
			End:          int(n.Range().Hi),
			Loc:          locOfNode(n, ctx.Parser.Source(), ctx),
			ElementTypes: elems(n.Args(), ctx),
		}
	case parser.N_TS_TUPLE_NAMED_MEMBER:
		n := node.(*parser.TsTupleNamedMember)
		return &TSNamedTupleMember{
			Type:        "TSNamedTupleMember",
			Start:       int(n.Range().Lo),
			End:         int(n.Range().Hi),
			Loc:         locOfNode(n, ctx.Parser.Source(), ctx),
			Optional:    n.Opt(),
			Label:       Convert(n.Label(), ctx),
			ElementType: ConvertTsTyp(n.Val(), ctx),
		}
	case parser.N_TS_REST:
		n := node.(*parser.TsRest)
		return &TSRestType{
			Type:           "TSRestType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
			TypeAnnotation: ConvertTsTyp(n.Arg(), ctx),
		}
	case parser.N_TS_OPT:
		n := node.(*parser.TsOpt)
		return &TSOptionalType{
			Type:           "TSOptionalType",
			Start:          int(n.Range().Lo),
			End:            int(n.Range().Hi),
			Loc:            locOfNode(n, ctx.Parser.Source(), ctx),
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
	rng, loc := locWithTypeInfo(node, false, ctx.Parser.Source(), ctx)
	return &TSExpressionWithTypeArguments{
		Type:           "TSExpressionWithTypeArguments",
		Start:          int(rng.Lo),
		End:            int(rng.Hi),
		Loc:            loc,
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
		Start:  int(psDec.Range().Lo),
		End:    int(psDec.Range().Hi),
		Loc:    locOfNode(psDec, ctx.Parser.Source(), ctx),
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
		Start:  int(psInst.Range().Lo),
		End:    int(psInst.Range().Hi),
		Loc:    locOfNode(psInst, ctx.Parser.Source(), ctx),
		Params: ret,
	}
}

func tsParamProp(node parser.Node, ctx *ConvertCtx) Node {
	ti, ok := isTyParamProp(node)
	if !ok {
		return nil
	}

	rng, loc := locWithTypeInfo(node, true, ctx.Parser.Source(), ctx)
	return &TSParameterProperty{
		Type:          "TSParameterProperty",
		Start:         int(rng.Lo),
		End:           int(rng.Hi),
		Loc:           loc,
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

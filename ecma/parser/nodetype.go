//go:generate go run github.com/hsiaosiyuan0/mole/script/nodetype_gen -d=../parser

package parser

type NodeType uint16

const (
	N_ILLEGAL NodeType = iota
	N_PROG             // #[visitor(Prog)]

	N_STMT_BEGIN
	N_STMT_EMPTY
	N_STMT_EXPR      // #[visitor(ExprStmt)]
	N_STMT_VAR_DEC   // #[visitor(VarDecStmt)]
	N_STMT_FN        // #[visitor(FnDec)]
	N_STMT_BLOCK     // #[visitor(BlockStmt)]
	N_STMT_DO_WHILE  // #[visitor(DoWhileStmt)]
	N_STMT_WHILE     // #[visitor(WhileStmt)]
	N_STMT_FOR       // #[visitor(ForStmt)]
	N_STMT_FOR_IN_OF // #[visitor(ForInOfStmt)]
	N_STMT_IF        // #[visitor(IfStmt)]
	N_STMT_SWITCH    // #[visitor(SwitchStmt)]
	N_STMT_BRK       // #[visitor(BrkStmt)]
	N_STMT_CONT      // #[visitor(ContStmt)]
	N_STMT_LABEL     // #[visitor(LabelStmt)]
	N_STMT_RET       // #[visitor(RetStmt)]
	N_STMT_THROW     // #[visitor(ThrowStmt)]
	N_STMT_TRY       // #[visitor(TryStmt)]
	N_STMT_DEBUG     // #[visitor(DebugStmt)]
	N_STMT_WITH      // #[visitor(WithStmt)]
	N_STMT_CLASS     // #[visitor(ClassDec)]
	N_STMT_IMPORT    // #[visitor(ImportDec)]
	N_STMT_EXPORT    // #[visitor(ExportDec)]
	N_STMT_END

	N_EXPR_BEGIN
	N_LIT_BEGIN
	N_LIT_NULL   // #[visitor(NullLit)]
	N_LIT_BOOL   // #[visitor(BoolLit)]
	N_LIT_NUM    // #[visitor(NumLit)]
	N_LIT_STR    // #[visitor(StrLit)]
	N_LIT_ARR    // #[visitor(ArrLit)]
	N_LIT_OBJ    // #[visitor(ObjLit)]
	N_LIT_REGEXP // #[visitor(RegLit)]
	N_LIT_END

	N_EXPR_NEW    // #[visitor(NewExpr)]
	N_EXPR_MEMBER // #[visitor(MemberExpr)]
	N_EXPR_CALL   // #[visitor(CallExpr)]
	N_EXPR_BIN    // #[visitor(BinExpr)]
	N_EXPR_UNARY  // #[visitor(UnaryExpr)]
	N_EXPR_UPDATE // #[visitor(UpdateExpr)]
	N_EXPR_COND   // #[visitor(CondExpr)]
	N_EXPR_ASSIGN // #[visitor(AssignExpr)]
	N_EXPR_FN     // #[visitor(FnDec)]
	N_EXPR_THIS   // #[visitor(ThisExpr)]
	N_EXPR_PAREN  // #[visitor(ParenExpr)]
	N_EXPR_ARROW  // #[visitor(ArrowFn)]
	N_EXPR_SEQ    // #[visitor(SeqExpr)]
	N_EXPR_CLASS  // #[visitor(ClassDec)]
	N_EXPR_TPL    // #[visitor(TplExpr)]
	N_EXPR_YIELD  // #[visitor(YieldExpr)]
	N_EXPR_CHAIN  // #[visitor(ChainExpr)]
	N_JSX_ELEM    // #[visitor(JsxElem)]
	N_NAME        // #[visitor(Ident)]
	N_IMPORT_CALL // #[visitor(ImportCall)]
	N_META_PROP   // #[visitor(MetaProp)]
	N_DECORATOR   // #[visitor(Decorator)]
	N_SPREAD      // #[visitor(Spread)]
	N_EXPR_END

	N_VAR_DEC      // #[visitor(VarDec)]
	N_PAT_REST     // #[visitor(RestPat)]
	N_PAT_ARRAY    // #[visitor(ArrPat)]
	N_PAT_ASSIGN   // #[visitor(AssignPat)]
	N_PAT_OBJ      // #[visitor(ObjPat)]
	N_PROP         // #[visitor(Prop)]
	N_SWITCH_CASE  // #[visitor(SwitchCase)]
	N_CATCH        // #[visitor(Catch)]
	N_CLASS_BODY   // #[visitor(ClassBody)]
	N_STATIC_BLOCK // #[visitor(StaticBlock)]
	N_METHOD       // #[visitor(Method)]
	N_FIELD        // #[visitor(Field)]
	N_SUPER        // #[visitor(Super)]
	N_IMPORT_SPEC  // #[visitor(ImportSpec)]
	N_EXPORT_SPEC  // #[visitor(ExportDec)]

	N_JSX_BEGIN
	N_JSX_ID           // #[visitor(JsxIdent)]
	N_JSX_MEMBER       // #[visitor(JsxMember)]
	N_JSX_NS           // #[visitor(JsxNsName)]
	N_JSX_ATTR_SPREAD  // #[visitor(JsxSpreadAttr)]
	N_JSX_CHILD_SPREAD // #[visitor(JsxSpreadChild)]
	N_JSX_OPEN         // #[visitor(JsxOpen)]
	N_JSX_CLOSE        // #[visitor(JsxClose)]
	N_JSX_EMPTY        // #[visitor(JsxEmpty)]
	N_JSX_EXPR_SPAN    // #[visitor(JsxExprSpan)]
	N_JSX_TXT          // #[visitor(JsxText)]
	N_JSX_ATTR         // #[visitor(JsxAttr)]
	N_JSX_END

	N_TS_BEGIN
	N_TS_TYP_ANNOT          // #[visitor(TsTypAnnot)]
	N_TS_ANY                // #[visitor(TsPredef)]
	N_TS_NUM                // #[visitor(TsPredef)]
	N_TS_BOOL               // #[visitor(TsPredef)]
	N_TS_STR                // #[visitor(TsPredef)]
	N_TS_SYM                // #[visitor(TsPredef)]
	N_TS_OBJ                // #[visitor(TsPredef)]
	N_TS_VOID               // #[visitor(TsPredef)]
	N_TS_NEVER              // #[visitor(TsPredef)]
	N_TS_UNKNOWN            // #[visitor(TsPredef)]
	N_TS_UNDEF              // #[visitor(TsPredef)]
	N_TS_BIGINT             // #[visitor(TsPredef)]
	N_TS_INTRINSIC          // #[visitor(TsPredef)]
	N_TS_NULL               // #[visitor(TsPredef)]
	N_TS_LIT                // #[visitor(TsLit)]
	N_TS_REF                // #[visitor(TsRef)]
	N_TS_LIT_OBJ            // #[visitor(TsObj)]
	N_TS_ARR                // #[visitor(TsArr)]
	N_TS_IDX_ACCESS         // #[visitor(TsIdxAccess)]
	N_TS_TUPLE              // #[visitor(TsTuple)]
	N_TS_REST               // #[visitor(TsRest)]
	N_TS_TUPLE_NAMED_MEMBER // #[visitor(TsTupleNamedMember)]
	N_TS_OPT                // #[visitor(TsOpt)]
	N_TS_TYP_QUERY          // #[visitor(TsTypQuery)]
	N_TS_COND               // #[visitor(TsCondType)]
	N_TS_TYP_OP             // #[visitor(TsTypOp)]
	N_TS_MAPPED             // #[visitor(TsMapped)]
	N_TS_TYP_INFER          // #[visitor(TsTypInfer)]
	N_TS_PAREN              // #[visitor(TsParen)]
	N_TS_THIS               // #[visitor(TsThis)]
	N_TS_NS_NAME            // #[visitor(TsNsName)]
	N_TS_PARAM              // #[visitor(TsParam)]
	N_TS_PARAM_DEC          // #[visitor(TsParamsDec)]
	N_TS_PARAM_INST         // #[visitor(TsParamsInst)]
	N_TS_PROP               // #[visitor(TsProp)]
	N_TS_CALL_SIG           // #[visitor(TsCallSig)]
	N_TS_NEW_SIG            // #[visitor(TsNewSig)]
	N_TS_IDX_SIG            // #[visitor(TsIdxSig)]
	N_TS_FN_TYP             // #[visitor(TsFnTyp)]
	N_TS_NEW                // #[visitor(TsNewSig)]
	N_TS_UNION_TYP          // #[visitor(TsUnionTyp)]
	N_TS_INTERSECT_TYP      // #[visitor(TsIntersectTyp)]
	N_TS_ROUGH_PARAM        // #[visitor(TsRoughParam)]
	N_TS_TYP_ASSERT         // #[visitor(TsTypAssert)]
	N_TS_TYP_DEC            // #[visitor(TsTypDec)]
	N_TS_INTERFACE          // #[visitor(TsInterface)]
	N_TS_INTERFACE_BODY     // #[visitor(TsInterfaceBody)]
	N_TS_ENUM               // #[visitor(TsEnum)]
	N_TS_ENUM_MEMBER        // #[visitor(TsEnumMember)]
	N_TS_IMPORT_ALIAS       // #[visitor(TsImportAlias)]
	N_TS_NAMESPACE          // #[visitor(TsNS)]
	N_TS_IMPORT_REQUIRE     // #[visitor(TsImportRequire)]
	N_TS_IMPORT_TYP         // #[visitor(TsImportType)]
	N_TS_EXPORT_ASSIGN      // #[visitor(TsExportAssign)]

	N_TS_DEC_VAR_DEC   // #[visitor(TsDec)]
	N_TS_DEC_FN        // #[visitor(TsDec)]
	N_TS_DEC_ENUM      // #[visitor(TsDec)]
	N_TS_DEC_CLASS     // #[visitor(TsDec)]
	N_TS_DEC_NS        // #[visitor(TsDec)]
	N_TS_DEC_MODULE    // #[visitor(TsDec)]
	N_TS_DEC_GLOBAL    // #[visitor(TsDec)]
	N_TS_DEC_INTERFACE // #[visitor(TsDec)]
	N_TS_DEC_TYP_DEC   // #[visitor(TsDec)]

	N_TS_TYP_PREDICATE // #[visitor(TsTypPredicate)]
	N_TS_NO_NULL       // #[visitor(TsNoNull)]
	N_TS_END

	N_NODE_DEF_END
)

func (nt NodeType) IsExpr() bool {
	return nt > N_EXPR_BEGIN && nt < N_EXPR_END
}

func (nt NodeType) IsStmt() bool {
	return nt > N_STMT_BEGIN && nt < N_STMT_END
}

func (nt NodeType) IsLit() bool {
	return nt > N_LIT_BEGIN && nt < N_LIT_END
}

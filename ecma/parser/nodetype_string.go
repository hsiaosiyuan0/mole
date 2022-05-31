// Code generated by script/nodetype_gen. DO NOT EDIT.

//go:generate go run github.com/hsiaosiyuan0/mole/script/nodetype_gen -d=../parser

package parser

var nodetypeStrings = map[NodeType]string{}

func init() {
	nodetypeStrings[N_CATCH] = "Catch"
	nodetypeStrings[N_CLASS_BODY] = "ClassBody"
	nodetypeStrings[N_DECORATOR] = "Decorator"
	nodetypeStrings[N_EXPORT_SPEC] = "ExportSpec"
	nodetypeStrings[N_EXPR_ARROW] = "ArrowFn"
	nodetypeStrings[N_EXPR_ASSIGN] = "AssignExpr"
	nodetypeStrings[N_EXPR_BIN] = "BinExpr"
	nodetypeStrings[N_EXPR_CALL] = "CallExpr"
	nodetypeStrings[N_EXPR_CHAIN] = "ChainExpr"
	nodetypeStrings[N_EXPR_CLASS] = "ClassDec"
	nodetypeStrings[N_EXPR_COND] = "CondExpr"
	nodetypeStrings[N_EXPR_FN] = "FnDec"
	nodetypeStrings[N_EXPR_MEMBER] = "MemberExpr"
	nodetypeStrings[N_EXPR_NEW] = "NewExpr"
	nodetypeStrings[N_EXPR_PAREN] = "ParenExpr"
	nodetypeStrings[N_EXPR_SEQ] = "SeqExpr"
	nodetypeStrings[N_EXPR_THIS] = "ThisExpr"
	nodetypeStrings[N_EXPR_TPL] = "TplExpr"
	nodetypeStrings[N_EXPR_UNARY] = "UnaryExpr"
	nodetypeStrings[N_EXPR_UPDATE] = "UpdateExpr"
	nodetypeStrings[N_EXPR_YIELD] = "YieldExpr"
	nodetypeStrings[N_FIELD] = "Field"
	nodetypeStrings[N_IMPORT_CALL] = "ImportCall"
	nodetypeStrings[N_IMPORT_SPEC] = "ImportSpec"
	nodetypeStrings[N_JSX_ATTR] = "JsxAttr"
	nodetypeStrings[N_JSX_ATTR_SPREAD] = "JsxSpreadAttr"
	nodetypeStrings[N_JSX_CHILD_SPREAD] = "JsxSpreadChild"
	nodetypeStrings[N_JSX_CLOSE] = "JsxClose"
	nodetypeStrings[N_JSX_ELEM] = "JsxElem"
	nodetypeStrings[N_JSX_EMPTY] = "JsxEmpty"
	nodetypeStrings[N_JSX_EXPR_SPAN] = "JsxExprSpan"
	nodetypeStrings[N_JSX_ID] = "JsxIdent"
	nodetypeStrings[N_JSX_MEMBER] = "JsxMember"
	nodetypeStrings[N_JSX_NS] = "JsxNsName"
	nodetypeStrings[N_JSX_OPEN] = "JsxOpen"
	nodetypeStrings[N_JSX_TXT] = "JsxText"
	nodetypeStrings[N_LIT_ARR] = "ArrLit"
	nodetypeStrings[N_LIT_BOOL] = "BoolLit"
	nodetypeStrings[N_LIT_NULL] = "NullLit"
	nodetypeStrings[N_LIT_NUM] = "NumLit"
	nodetypeStrings[N_LIT_OBJ] = "ObjLit"
	nodetypeStrings[N_LIT_REGEXP] = "RegLit"
	nodetypeStrings[N_LIT_STR] = "StrLit"
	nodetypeStrings[N_META_PROP] = "MetaProp"
	nodetypeStrings[N_METHOD] = "Method"
	nodetypeStrings[N_NAME] = "Ident"
	nodetypeStrings[N_PAT_ARRAY] = "ArrPat"
	nodetypeStrings[N_PAT_ASSIGN] = "AssignPat"
	nodetypeStrings[N_PAT_OBJ] = "ObjPat"
	nodetypeStrings[N_PAT_REST] = "RestPat"
	nodetypeStrings[N_PROG] = "Prog"
	nodetypeStrings[N_PROP] = "Prop"
	nodetypeStrings[N_SPREAD] = "Spread"
	nodetypeStrings[N_STATIC_BLOCK] = "StaticBlock"
	nodetypeStrings[N_STMT_BLOCK] = "BlockStmt"
	nodetypeStrings[N_STMT_BRK] = "BrkStmt"
	nodetypeStrings[N_STMT_CLASS] = "ClassDec"
	nodetypeStrings[N_STMT_CONT] = "ContStmt"
	nodetypeStrings[N_STMT_DEBUG] = "DebugStmt"
	nodetypeStrings[N_STMT_DO_WHILE] = "DoWhileStmt"
	nodetypeStrings[N_STMT_EXPORT] = "ExportDec"
	nodetypeStrings[N_STMT_EXPR] = "ExprStmt"
	nodetypeStrings[N_STMT_FN] = "FnDec"
	nodetypeStrings[N_STMT_FOR] = "ForStmt"
	nodetypeStrings[N_STMT_FOR_IN_OF] = "ForInOfStmt"
	nodetypeStrings[N_STMT_IF] = "IfStmt"
	nodetypeStrings[N_STMT_IMPORT] = "ImportDec"
	nodetypeStrings[N_STMT_LABEL] = "LabelStmt"
	nodetypeStrings[N_STMT_RET] = "RetStmt"
	nodetypeStrings[N_STMT_SWITCH] = "SwitchStmt"
	nodetypeStrings[N_STMT_THROW] = "ThrowStmt"
	nodetypeStrings[N_STMT_TRY] = "TryStmt"
	nodetypeStrings[N_STMT_VAR_DEC] = "VarDecStmt"
	nodetypeStrings[N_STMT_WHILE] = "WhileStmt"
	nodetypeStrings[N_STMT_WITH] = "WithStmt"
	nodetypeStrings[N_SUPER] = "Super"
	nodetypeStrings[N_SWITCH_CASE] = "SwitchCase"
	nodetypeStrings[N_TS_ANY] = "TsPredef"
	nodetypeStrings[N_TS_ARR] = "TsArr"
	nodetypeStrings[N_TS_BIGINT] = "TsPredef"
	nodetypeStrings[N_TS_BOOL] = "TsPredef"
	nodetypeStrings[N_TS_CALL_SIG] = "TsCallSig"
	nodetypeStrings[N_TS_COND] = "TsCondType"
	nodetypeStrings[N_TS_DEC_CLASS] = "TsDec"
	nodetypeStrings[N_TS_DEC_ENUM] = "TsDec"
	nodetypeStrings[N_TS_DEC_FN] = "TsDec"
	nodetypeStrings[N_TS_DEC_GLOBAL] = "TsDec"
	nodetypeStrings[N_TS_DEC_INTERFACE] = "TsDec"
	nodetypeStrings[N_TS_DEC_MODULE] = "TsDec"
	nodetypeStrings[N_TS_DEC_NS] = "TsDec"
	nodetypeStrings[N_TS_DEC_TYP_DEC] = "TsDec"
	nodetypeStrings[N_TS_DEC_VAR_DEC] = "TsDec"
	nodetypeStrings[N_TS_ENUM] = "TsEnum"
	nodetypeStrings[N_TS_ENUM_MEMBER] = "TsEnumMember"
	nodetypeStrings[N_TS_EXPORT_ASSIGN] = "TsExportAssign"
	nodetypeStrings[N_TS_FN_TYP] = "TsFnTyp"
	nodetypeStrings[N_TS_IDX_ACCESS] = "TsIdxAccess"
	nodetypeStrings[N_TS_IDX_SIG] = "TsIdxSig"
	nodetypeStrings[N_TS_IMPORT_ALIAS] = "TsImportAlias"
	nodetypeStrings[N_TS_IMPORT_REQUIRE] = "TsImportRequire"
	nodetypeStrings[N_TS_IMPORT_TYP] = "TsImportType"
	nodetypeStrings[N_TS_INTERFACE] = "TsInterface"
	nodetypeStrings[N_TS_INTERFACE_BODY] = "TsInterfaceBody"
	nodetypeStrings[N_TS_INTERSECT_TYP] = "TsIntersectTyp"
	nodetypeStrings[N_TS_INTRINSIC] = "TsPredef"
	nodetypeStrings[N_TS_LIT] = "TsLit"
	nodetypeStrings[N_TS_LIT_OBJ] = "TsObj"
	nodetypeStrings[N_TS_MAPPED] = "TsMapped"
	nodetypeStrings[N_TS_NAMESPACE] = "TsNS"
	nodetypeStrings[N_TS_NEVER] = "TsPredef"
	nodetypeStrings[N_TS_NEW] = "TsNewSig"
	nodetypeStrings[N_TS_NEW_SIG] = "TsNewSig"
	nodetypeStrings[N_TS_NO_NULL] = "TsNoNull"
	nodetypeStrings[N_TS_NS_NAME] = "TsNsName"
	nodetypeStrings[N_TS_NULL] = "TsPredef"
	nodetypeStrings[N_TS_NUM] = "TsPredef"
	nodetypeStrings[N_TS_OBJ] = "TsPredef"
	nodetypeStrings[N_TS_OPT] = "TsOpt"
	nodetypeStrings[N_TS_PARAM] = "TsParam"
	nodetypeStrings[N_TS_PARAM_DEC] = "TsParamsDec"
	nodetypeStrings[N_TS_PARAM_INST] = "TsParamsInst"
	nodetypeStrings[N_TS_PAREN] = "TsParen"
	nodetypeStrings[N_TS_PROP] = "TsProp"
	nodetypeStrings[N_TS_REF] = "TsRef"
	nodetypeStrings[N_TS_REST] = "TsRest"
	nodetypeStrings[N_TS_ROUGH_PARAM] = "TsRoughParam"
	nodetypeStrings[N_TS_STR] = "TsPredef"
	nodetypeStrings[N_TS_SYM] = "TsPredef"
	nodetypeStrings[N_TS_THIS] = "TsThis"
	nodetypeStrings[N_TS_TUPLE] = "TsTuple"
	nodetypeStrings[N_TS_TUPLE_NAMED_MEMBER] = "TsTupleNamedMember"
	nodetypeStrings[N_TS_TYP_ANNOT] = "TsTypAnnot"
	nodetypeStrings[N_TS_TYP_ASSERT] = "TsTypAssert"
	nodetypeStrings[N_TS_TYP_DEC] = "TsTypDec"
	nodetypeStrings[N_TS_TYP_INFER] = "TsTypInfer"
	nodetypeStrings[N_TS_TYP_OP] = "TsTypOp"
	nodetypeStrings[N_TS_TYP_PREDICATE] = "TsTypPredicate"
	nodetypeStrings[N_TS_TYP_QUERY] = "TsTypQuery"
	nodetypeStrings[N_TS_UNDEF] = "TsPredef"
	nodetypeStrings[N_TS_UNION_TYP] = "TsUnionTyp"
	nodetypeStrings[N_TS_UNKNOWN] = "TsPredef"
	nodetypeStrings[N_TS_VOID] = "TsPredef"
	nodetypeStrings[N_VAR_DEC] = "VarDec"
}

func (nt NodeType) String() string {
	return nodetypeStrings[nt]
}

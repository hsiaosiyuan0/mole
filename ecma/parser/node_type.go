package parser

type NodeType uint16

const (
	N_ILLEGAL NodeType = iota
	N_PROG

	N_STMT_BEGIN
	N_STMT_EXPR
	N_STMT_EMPTY
	N_STMT_VAR_DEC
	N_STMT_FN
	N_STMT_BLOCK
	N_STMT_DO_WHILE
	N_STMT_WHILE
	N_STMT_FOR
	N_STMT_FOR_IN_OF
	N_STMT_IF
	N_STMT_SWITCH
	N_STMT_BRK
	N_STMT_CONT
	N_STMT_LABEL
	N_STMT_RET
	N_STMT_THROW
	N_STMT_TRY
	N_STMT_DEBUG
	N_STMT_WITH
	N_STMT_CLASS
	N_STMT_IMPORT
	N_STMT_EXPORT
	N_STMT_END

	N_EXPR_BEGIN

	N_LIT_BEGIN
	N_LIT_NULL
	N_LIT_BOOL
	N_LIT_NUM
	N_LIT_STR
	N_LIT_ARR
	N_LIT_OBJ
	N_LIT_REGEXP
	N_LIT_END

	N_EXPR_NEW
	N_EXPR_MEMBER
	N_EXPR_CALL
	N_EXPR_BIN
	N_EXPR_UNARY
	N_EXPR_UPDATE
	N_EXPR_COND
	N_EXPR_ASSIGN
	N_EXPR_FN
	N_EXPR_THIS
	N_EXPR_PAREN
	N_EXPR_ARROW
	N_EXPR_SEQ
	N_EXPR_CLASS
	N_EXPR_TPL
	N_EXPR_YIELD
	N_EXPR_CHAIN
	N_JSX_ELEM
	N_NAME
	N_IMPORT_CALL
	N_META_PROP
	N_SPREAD

	N_EXPR_END

	N_VAR_DEC
	N_PAT_REST
	N_PAT_ARRAY
	N_PAT_ASSIGN
	N_PAT_OBJ
	N_PROP
	N_SWITCH_CASE
	N_CATCH
	N_ClASS_BODY
	N_STATIC_BLOCK
	N_METHOD
	N_FIELD
	N_SUPER
	N_IMPORT_SPEC
	N_EXPORT_SPEC

	N_JSX_ID
	N_JSX_MEMBER
	N_JSX_NS
	N_JSX_ATTR_SPREAD
	N_JSX_CHILD_SPREAD
	N_JSX_OPEN
	N_JSX_CLOSE
	N_JSX_EMPTY
	N_JSX_EXPR_SPAN
	N_JSX_TXT
	N_JSX_ATTR

	N_TS_TYP_ANNOT
	N_TS_ANY
	N_TS_NUM
	N_TS_BOOL
	N_TS_STR
	N_TS_SYM
	N_TS_OBJ
	N_TS_VOID
	N_TS_NEVER
	N_TS_UNKNOWN
	N_TS_UNDEF
	N_TS_NULL
	N_TS_LIT
	N_TS_REF
	N_TS_LIT_OBJ
	N_TS_ARR
	N_TS_TUPLE
	N_TS_QUERY
	N_TS_PAREN
	N_TS_THIS
	N_TS_NS_NAME
	N_TS_PARAM
	N_TS_PARAM_DEC
	N_TS_PARAM_INST
	N_TS_ARG
	N_TS_PROP
	N_TS_CALL_SIG
	N_TS_NEW_SIG
	N_TS_IDX_SIG
	N_TS_FN_TYP
	N_TS_UNION_TYP
	N_TS_INTERSEC_TYP
	N_TS_ROUGH_PARAM
	N_TS_TYP_ASSERT
	N_TS_TYP_DEC
	N_TS_INTERFACE
	N_TS_INTERFACE_BODY
	N_TS_ENUM
	N_TS_ENUM_MEMBER
	N_TS_IMPORT_ALIAS
	N_TS_NAMESPACE
	N_TS_IMPORT_REQUIRE
	N_TS_EXPORT_ASSIGN

	N_TS_DEC_VAR_DEC
	N_TS_DEC_FN
	N_TS_DEC_ENUM
	N_TS_DEC_CLASS
	N_TS_DEC_NS
	N_TS_DEC_MODULE
	N_TS_DEC_GLOBAL
	N_TS_DEC_INTERFACE
	N_TS_DEC_TYP_DEC

	N_TS_TYP_PREDICATE
	N_TS_AS
	N_TS_NO_NULL

	N_NODE_DEF_END
)

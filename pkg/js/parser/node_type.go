package parser

type NodeType int

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

	N_NODE_DEF_END
)
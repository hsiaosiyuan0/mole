package parser

import "log"

type TraverseCtx struct {
	RecordDepth ContOrStop
	Depth       int

	RecordPath ContOrStop
	Path       []string

	Extra interface{}
}

func NewTraverseCtx() *TraverseCtx {
	return &TraverseCtx{
		Path: make([]string, 0),
	}
}

func (c *TraverseCtx) IncDepth() {
	c.Depth += 1
}

func (c *TraverseCtx) DecDepth() {
	c.Depth -= 1
}

func (c *TraverseCtx) AppendPath(part string) {
	c.Path = append(c.Path, part)
}

func (c *TraverseCtx) PopPath() {
	c.Path = c.Path[:len(c.Path)-1]
}

type ContOrStop bool

const (
	TRAVERSE_CONT ContOrStop = true
	TRAVERSE_STOP ContOrStop = false
)

const (
	N_L_BEGIN NodeType = N_NODE_DEF_END + iota

	N_PROG_BEFORE = N_L_BEGIN + iota
	N_PROG_AFTER  = N_L_BEGIN + iota

	N_STMT_EXPR_BEFORE      = N_L_BEGIN + iota
	N_STMT_EXPR_AFTER       = N_L_BEGIN + iota
	N_STMT_EMPTY_BEFORE     = N_L_BEGIN + iota
	N_STMT_EMPTY_AFTER      = N_L_BEGIN + iota
	N_STMT_VAR_DEC_BEFORE   = N_L_BEGIN + iota
	N_STMT_VAR_DEC_AFTER    = N_L_BEGIN + iota
	N_STMT_FN_BEFORE        = N_L_BEGIN + iota
	N_STMT_FN_AFTER         = N_L_BEGIN + iota
	N_STMT_BLOCK_BEFORE     = N_L_BEGIN + iota
	N_STMT_BLOCK_AFTTER     = N_L_BEGIN + iota
	N_STMT_DO_WHILE_BEFORE  = N_L_BEGIN + iota
	N_STMT_DO_WHILE_AFTER   = N_L_BEGIN + iota
	N_STMT_WHILE_BEFORE     = N_L_BEGIN + iota
	N_STMT_WHILE_AFTER      = N_L_BEGIN + iota
	N_STMT_FOR_BEFORE       = N_L_BEGIN + iota
	N_STMT_FOR_AFTER        = N_L_BEGIN + iota
	N_STMT_FOR_IN_OF_BEFORE = N_L_BEGIN + iota
	N_STMT_FOR_IN_OF_AFTER  = N_L_BEGIN + iota
	N_STMT_IF_BEFORE        = N_L_BEGIN + iota
	N_STMT_IF_AFTER         = N_L_BEGIN + iota
	N_STMT_SWITCH_BEFORE    = N_L_BEGIN + iota
	N_STMT_SWITCH_AFTER     = N_L_BEGIN + iota
	N_STMT_BRK_BEFORE       = N_L_BEGIN + iota
	N_STMT_BRK_AFTER        = N_L_BEGIN + iota
	N_STMT_CONT_BEFORE      = N_L_BEGIN + iota
	N_STMT_CONT_AFTER       = N_L_BEGIN + iota
	N_STMT_LABEL_BEFORE     = N_L_BEGIN + iota
	N_STMT_LABEL_AFTER      = N_L_BEGIN + iota
	N_STMT_RET_BEFORE       = N_L_BEGIN + iota
	N_STMT_RET_AFTER        = N_L_BEGIN + iota
	N_STMT_THROW_BEFORE     = N_L_BEGIN + iota
	N_STMT_THROW_AFTER      = N_L_BEGIN + iota
	N_STMT_TRY_BEFORE       = N_L_BEGIN + iota
	N_STMT_TRY_AFTER        = N_L_BEGIN + iota
	N_STMT_DEBUG_BEFORE     = N_L_BEGIN + iota
	N_STMT_DEBUG_AFTER      = N_L_BEGIN + iota
	N_STMT_WITH_BEFORE      = N_L_BEGIN + iota
	N_STMT_WITH_AFTER       = N_L_BEGIN + iota
	N_STMT_CLASS_BEFORE     = N_L_BEGIN + iota
	N_STMT_CLASS_AFTER      = N_L_BEGIN + iota
	N_STMT_IMPORT_BEFORE    = N_L_BEGIN + iota
	N_STMT_IMPORT_AFTER     = N_L_BEGIN + iota
	N_STMT_EXPORT_BEFORE    = N_L_BEGIN + iota
	N_STMT_EXPORT_AFTER     = N_L_BEGIN + iota

	N_LIT_NULL_BEFORE   = N_L_BEGIN + iota
	N_LIT_NULL_AFTER    = N_L_BEGIN + iota
	N_LIT_BOOL_BEFORE   = N_L_BEGIN + iota
	N_LIT_BOOL_AFTER    = N_L_BEGIN + iota
	N_LIT_NUM_BEFORE    = N_L_BEGIN + iota
	N_LIT_NUM_AFTER     = N_L_BEGIN + iota
	N_LIT_STR_BEFORE    = N_L_BEGIN + iota
	N_LIT_STR_AFTER     = N_L_BEGIN + iota
	N_LIT_ARR_BEFORE    = N_L_BEGIN + iota
	N_LIT_ARR_AFTER     = N_L_BEGIN + iota
	N_LIT_OBJ_BEFORE    = N_L_BEGIN + iota
	N_LIT_OBJ_AFTER     = N_L_BEGIN + iota
	N_LIT_REGEXP_BEFORE = N_L_BEGIN + iota
	N_LIT_REGEXP_AFTER  = N_L_BEGIN + iota

	N_EXPR_NEW_BEFORE    = N_L_BEGIN + iota
	N_EXPR_NEW_AFTER     = N_L_BEGIN + iota
	N_EXPR_MEMBER_BEFORE = N_L_BEGIN + iota
	N_EXPR_MEMBER_AFTER  = N_L_BEGIN + iota
	N_EXPR_CALL_BEFORE   = N_L_BEGIN + iota
	N_EXPR_CALL_AFTER    = N_L_BEGIN + iota
	N_EXPR_BIN_BEFORE    = N_L_BEGIN + iota
	N_EXPR_BIN_AFTER     = N_L_BEGIN + iota
	N_EXPR_UNARY_BEFORE  = N_L_BEGIN + iota
	N_EXPR_UNARY_AFTER   = N_L_BEGIN + iota
	N_EXPR_UPDATE_BEFORE = N_L_BEGIN + iota
	N_EXPR_UPDATE_AFTER  = N_L_BEGIN + iota
	N_EXPR_COND_BEFORE   = N_L_BEGIN + iota
	N_EXPR_COND_AFTER    = N_L_BEGIN + iota
	N_EXPR_ASSIGN_BEFORE = N_L_BEGIN + iota
	N_EXPR_ASSIGN_AFTER  = N_L_BEGIN + iota
	N_EXPR_FN_BEFORE     = N_L_BEGIN + iota
	N_EXPR_FN_AFTER      = N_L_BEGIN + iota
	N_EXPR_THIS_BEFORE   = N_L_BEGIN + iota
	N_EXPR_THIS_AFTER    = N_L_BEGIN + iota
	N_EXPR_PAREN_BEFORE  = N_L_BEGIN + iota
	N_EXPR_PAREN_AFTER   = N_L_BEGIN + iota
	N_EXPR_ARROW_BEFORE  = N_L_BEGIN + iota
	N_EXPR_ARROW_AFTER   = N_L_BEGIN + iota
	N_EXPR_SEQ_BEFORE    = N_L_BEGIN + iota
	N_EXPR_SEQ_AFTER     = N_L_BEGIN + iota
	N_EXPR_CLASS_BEFORE  = N_L_BEGIN + iota
	N_EXPR_CLASS_AFTER   = N_L_BEGIN + iota
	N_EXPR_TPL_BEFORE    = N_L_BEGIN + iota
	N_EXPR_TPL_AFTER     = N_L_BEGIN + iota
	N_EXPR_YIELD_BEFORE  = N_L_BEGIN + iota
	N_EXPR_YIELD_AFTER   = N_L_BEGIN + iota
	N_EXPR_CHAIN_BEFORE  = N_L_BEGIN + iota
	N_EXPR_CHAIN_AFTER   = N_L_BEGIN + iota
	N_JSX_ELEM_BEFORE    = N_L_BEGIN + iota
	N_JSX_ELEM_AFTER     = N_L_BEGIN + iota
	N_NAME_BEFORE        = N_L_BEGIN + iota
	N_NAME_AFTER         = N_L_BEGIN + iota
	N_IMPORT_CALL_BEFORE = N_L_BEGIN + iota
	N_IMPORT_CALL_AFTER  = N_L_BEGIN + iota
	N_META_PROP_BEFORE   = N_L_BEGIN + iota
	N_META_PROP_AFTER    = N_L_BEGIN + iota
	N_SPREAD_BEFORE      = N_L_BEGIN + iota
	N_SPREAD_AFTER       = N_L_BEGIN + iota

	N_VAR_DEC_BEFORE      = N_L_BEGIN + iota
	N_VAR_DEC_AFTER       = N_L_BEGIN + iota
	N_PAT_REST_BEFORE     = N_L_BEGIN + iota
	N_PAT_REST_AFTER      = N_L_BEGIN + iota
	N_PAT_ARRAY_BEFORE    = N_L_BEGIN + iota
	N_PAT_ARRAY_AFTER     = N_L_BEGIN + iota
	N_PAT_ASSIGN_BEFORE   = N_L_BEGIN + iota
	N_PAT_ASSIGN_AFTER    = N_L_BEGIN + iota
	N_PAT_OBJ_BEFORE      = N_L_BEGIN + iota
	N_PAT_OBJ_AFTER       = N_L_BEGIN + iota
	N_PROP_BEFORE         = N_L_BEGIN + iota
	N_PROP_AFTER          = N_L_BEGIN + iota
	N_SWITCH_CASE_BEFORE  = N_L_BEGIN + iota
	N_SWITCH_CASE_AFTER   = N_L_BEGIN + iota
	N_CATCH_BEFORE        = N_L_BEGIN + iota
	N_CATCH_AFTER         = N_L_BEGIN + iota
	N_ClASS_BODY_BEFORE   = N_L_BEGIN + iota
	N_ClASS_BODY_AFTER    = N_L_BEGIN + iota
	N_STATIC_BLOCK_BEFORE = N_L_BEGIN + iota
	N_STATIC_BLOCK_AFTER  = N_L_BEGIN + iota
	N_METHOD_BEFORE       = N_L_BEGIN + iota
	N_METHOD_AFTER        = N_L_BEGIN + iota
	N_FIELD_BEFORE        = N_L_BEGIN + iota
	N_FIELD_AFTER         = N_L_BEGIN + iota
	N_SUPER_BEFORE        = N_L_BEGIN + iota
	N_SUPER_AFTER         = N_L_BEGIN + iota
	N_IMPORT_SPEC_BEFORE  = N_L_BEGIN + iota
	N_IMPORT_SPEC_AFTER   = N_L_BEGIN + iota
	N_EXPORT_SPEC_BEFORE  = N_L_BEGIN + iota
	N_EXPORT_SPEC_AFTER   = N_L_BEGIN + iota

	N_JSX_ID_BEFORE           = N_L_BEGIN + iota
	N_JSX_ID_AFTER            = N_L_BEGIN + iota
	N_JSX_MEMBER_BEFORE       = N_L_BEGIN + iota
	N_JSX_MEMBER_AFTER        = N_L_BEGIN + iota
	N_JSX_NS_BEFORE           = N_L_BEGIN + iota
	N_JSX_NS_AFTER            = N_L_BEGIN + iota
	N_JSX_ATTR_SPREAD_BEFORE  = N_L_BEGIN + iota
	N_JSX_ATTR_SPREAD_AFTER   = N_L_BEGIN + iota
	N_JSX_CHILD_SPREAD_BEFORE = N_L_BEGIN + iota
	N_JSX_CHILD_SPREAD_AFTER  = N_L_BEGIN + iota
	N_JSX_OPEN_BEFORE         = N_L_BEGIN + iota
	N_JSX_OPEN_AFTER          = N_L_BEGIN + iota
	N_JSX_CLOSE_BEFORE        = N_L_BEGIN + iota
	N_JSX_CLOSE_AFTER         = N_L_BEGIN + iota
	N_JSX_EMPTY_BEFORE        = N_L_BEGIN + iota
	N_JSX_EMPTY_AFTER         = N_L_BEGIN + iota
	N_JSX_EXPR_SPAN_BEFORE    = N_L_BEGIN + iota
	N_JSX_EXPR_SPAN_AFTER     = N_L_BEGIN + iota
	N_JSX_TXT_BEFORE          = N_L_BEGIN + iota
	N_JSX_TXT_AFTER           = N_L_BEGIN + iota
	N_JSX_ATTR_BEFORE         = N_L_BEGIN + iota
	N_JSX_ATTR_AFTER          = N_L_BEGIN + iota

	N_L_END = N_L_BEGIN + iota
)

// the type of `v` is `func VisitorImpl`
type NodeFn = func(n Node, v interface{}, ctx *TraverseCtx) ContOrStop

// for the `func VisitFn` can be routed in O(1)
type NodeFns = [N_L_END]NodeFn

func NewVisitorImpl() *NodeFns {
	n := DefaultListenerImpl
	return &n
}

func VisitNode(n Node, v interface{}, ctx *TraverseCtx, silent bool) ContOrStop {
	if n == nil {
		return TRAVERSE_CONT
	}

	vi := v.(*NodeFns)
	typ := n.Type()
	fn := vi[typ]
	if fn == nil {
		log.Fatalf("Impl does not exist for NodeType %d", typ)
	}

	return fn(n, v, ctx)
}

func VisitNodes(ns []Node, v interface{}, ctx *TraverseCtx, silent bool) ContOrStop {
	vi := v.(*NodeFns)
	for _, n := range ns {
		if !VisitNode(n, vi, ctx, silent) {
			return TRAVERSE_STOP
		}
	}
	return TRAVERSE_CONT
}

var DefaultVisitorImpl NodeFns = [N_L_END]NodeFn{
	nil, // N_ILLEGAL
	VisitProg,

	nil, // N_STMT_BEGIN
	VisitExprStmt,
	VisitEmptyStmt,
	VisitVarDecStmt,
	VisitFnDec,
	VisitBlockStmt,
	VisitDoWhileStmt,
	VisitWhileStmt,
	VisitForStmt,
	VisitForInOfStmt,
	VisitIfStmt,
	VisitSwitchStmt,
	VisitBrkStmt,
	VisitContStmt,
	VisitLabelStmt,
	VisitRetStmt,
	VisitThrowStmt,
	VisitTryStmt,
	VisitDebugStmt,
	VisitWithStmt,
	VisitClassDec,
	VisitImportDec,
	VisitExportDec,
	nil, // N_STMT_END

	nil, // N_EXPR_BEGIN

	nil, // N_LIT_BEGIN
	VisitNull,
	VisitBool,
	VisitNum,
	VisitStr,
	VisitArr,
	VisitObj,
	VisitReg,
	nil, // N_LIT_END

	VisitNewExpr,
	VisitMemberExpr,
	VisitCallExpr,
	VisitBinExpr,
	VisitUnaryExpr,
	VisitUpdateExpr,
	VisitCondExpr,
	VisitAssignExpr,
	VisitFnExpr,
	VisitThisExpr,
	VisitParenExpr,
	VisitArrowFn,
	VisitSeqExpr,
	VisitClassExpr,
	VisitTplExpr,
	VisitYieldExpr,
	VisitOptChainExpr,
	VisitJsxElem,
	VisitIdent,
	VisitImportCall,
	VisitMetaProp,
	VisitSpread,
	nil, // N_EXPR_END

	VisitVarDec,
	VisitRestPat,
	VisitArrPat,
	VisitAssignPat,
	VisitObjPat,
	VisitProp,
	VisitSwitchCase,
	VisitCatch,
	VisitClassBody,
	VisitStaticBlock,
	VisitMethod,
	VisitField,
	VisitSuper,
	VisitImportSpec,
	VisitExportSpec,

	VisitJsxIdent,
	VisitJsxMemberExpr,
	VisitJsxNsName,
	VisitJsxAttrSpread,
	VisitJsxChildSpread,
	VisitJsxOpen,
	VisitJsxClose,
	VisitJsxEmpty,
	VisitJsxExprSpan,
	VisitJsxText,
	VisitJsxAttr,
}

func VisitProg(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Prog)
	return VisitNodes(an.stmts, v, ctx, false)
}

func VisitExprStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExprStmt)
	return VisitNode(an.expr, v, ctx, false)
}
func VisitEmptyStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitVarDecStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDecStmt)
	return VisitNodes(an.decList, v, ctx, false)
}
func VisitFnDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.params, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitBlockStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BlockStmt)
	return VisitNodes(an.body, v, ctx, false)
}
func VisitDoWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*DoWhileStmt)
	if !VisitNode(an.body, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.test, v, ctx, false)
}
func VisitWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WhileStmt)
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitForStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForStmt)
	if !VisitNode(an.init, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.update, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitForInOfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForInOfStmt)
	if !VisitNode(an.left, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.right, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitIfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*IfStmt)
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.cons, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.alt, v, ctx, false)
}
func VisitSwitchStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchStmt)
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.cases, v, ctx, false)
}
func VisitBrkStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BrkStmt)
	return VisitNode(an.label, v, ctx, false)
}
func VisitContStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ContStmt)
	return VisitNode(an.label, v, ctx, false)
}
func VisitLabelStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*LabelStmt)
	if !VisitNode(an.label, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitRetStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RetStmt)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitThrowStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ThrowStmt)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitTryStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TryStmt)
	if !VisitNode(an.try, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.catch, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.fin, v, ctx, false)
}
func VisitDebugStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitWithStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WithStmt)
	if !VisitNode(an.expr, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitClassDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.super, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitImportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportDec)
	if !VisitNodes(an.specs, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.src, v, ctx, false)
}
func VisitExportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportDec)
	if !VisitNodes(an.specs, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.src, v, ctx, false)
}

func VisitNull(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitBool(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitNum(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitStr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitArr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitObj(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitReg(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}

func VisitNewExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*NewExpr)
	if !VisitNode(an.callee, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.args, v, ctx, false)
}
func VisitMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MemberExpr)
	if !VisitNode(an.obj, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx, false)
}
func VisitCallExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CallExpr)
	if !VisitNode(an.callee, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.args, v, ctx, false)
}
func VisitBinExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BinExpr)
	if !VisitNode(an.lhs, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx, false)
}
func VisitUnaryExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UnaryExpr)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitUpdateExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UpdateExpr)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitCondExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CondExpr)
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.cons, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.alt, v, ctx, false)
}
func VisitAssignExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignExpr)
	if !VisitNode(an.lhs, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx, false)
}
func VisitFnExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.params, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitThisExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitParenExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ParenExpr)
	return VisitNode(an.expr, v, ctx, false)
}
func VisitArrowFn(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrowFn)
	if !VisitNodes(an.params, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitSeqExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SeqExpr)
	return VisitNodes(an.elems, v, ctx, false)
}
func VisitClassExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.super, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitTplExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TplExpr)
	if !VisitNode(an.tag, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.elems, v, ctx, false)
}
func VisitYieldExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*YieldExpr)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitOptChainExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ChainExpr)
	return VisitNode(an.expr, v, ctx, false)
}
func VisitJsxElem(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxElem)
	if !VisitNode(an.open, v, ctx, false) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.children, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.close, v, ctx, false)
}
func VisitIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitImportCall(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportCall)
	return VisitNode(an.src, v, ctx, false)
}
func VisitMetaProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MetaProp)
	if !VisitNode(an.meta, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx, false)
}
func VisitSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Spread)
	return VisitNode(an.arg, v, ctx, false)
}

func VisitVarDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.init, v, ctx, false)
}
func VisitRestPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RestPat)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitArrPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrPat)
	return VisitNodes(an.elems, v, ctx, false)
}
func VisitAssignPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignPat)
	if !VisitNode(an.lhs, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx, false)
}
func VisitObjPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ObjPat)
	return VisitNodes(an.props, v, ctx, false)
}
func VisitProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Prop)
	if !VisitNode(an.key, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx, false)
}
func VisitSwitchCase(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchCase)
	if !VisitNode(an.test, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.cons, v, ctx, false)
}
func VisitCatch(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Catch)
	if !VisitNode(an.param, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx, false)
}
func VisitClassBody(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassBody)
	return VisitNodes(an.elems, v, ctx, false)
}
func VisitStaticBlock(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*StaticBlock)
	return VisitNodes(an.body, v, ctx, false)
}
func VisitMethod(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Method)
	if !VisitNode(an.key, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx, false)
}
func VisitField(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Field)
	if !VisitNode(an.key, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx, false)
}
func VisitSuper(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitImportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportSpec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.local, v, ctx, false)
}
func VisitExportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportSpec)
	if !VisitNode(an.id, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.local, v, ctx, false)
}

func VisitJsxIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxMemberExpr)
	if !VisitNode(an.obj, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx, false)
}
func VisitJsxNsName(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxNsName)
	if !VisitNode(an.ns, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.name, v, ctx, false)
}
func VisitJsxAttrSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxSpreadAttr)
	return VisitNode(an.arg, v, ctx, false)
}
func VisitJsxChildSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxSpreadChild)
	return VisitNode(an.expr, v, ctx, false)
}
func VisitJsxOpen(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxOpen)
	if !VisitNode(an.name, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.attrs, v, ctx, false)
}
func VisitJsxClose(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxClose)
	return VisitNode(an.name, v, ctx, false)
}
func VisitJsxEmpty(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxExprSpan(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxExprSpan)
	return VisitNode(an.expr, v, ctx, false)
}
func VisitJsxText(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxAttr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxAttr)
	if !VisitNode(an.name, v, ctx, false) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.val, v, ctx, false)
}

func NewListenerImpl() *NodeFns {
	n := DefaultListenerImpl
	return &n
}

var DefaultListenerImpl NodeFns = [N_L_END]NodeFn{}

func callNodeFn(node Node, typ NodeType, v interface{}, ctx *TraverseCtx) {
	vi := v.(*NodeFns)
	fn := vi[typ]
	if fn == nil {
		return
	}
	fn(node, v, ctx)
}

func ListenProg(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Prog)
	callNodeFn(n, N_PROP_BEFORE, v, ctx)
	VisitNodes(an.stmts, v, ctx, true)
	callNodeFn(n, N_PROP_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func ListenExprStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExprStmt)
	callNodeFn(n, N_STMT_EXPR_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	callNodeFn(n, N_STMT_EXPR_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenEmptyStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_STMT_EMPTY_BEFORE, v, ctx)
	callNodeFn(n, N_STMT_EMPTY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenVarDecStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDecStmt)
	callNodeFn(n, N_STMT_VAR_DEC_BEFORE, v, ctx)
	VisitNodes(an.decList, v, ctx, true)
	callNodeFn(n, N_STMT_VAR_DEC_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenFnDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	callNodeFn(n, N_STMT_FN_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNodes(an.params, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_FN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenBlockStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BlockStmt)
	callNodeFn(n, N_STMT_BLOCK_BEFORE, v, ctx)
	VisitNodes(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_BLOCK_AFTTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenDoWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*DoWhileStmt)
	callNodeFn(n, N_STMT_DO_WHILE_BEFORE, v, ctx)
	VisitNode(an.body, v, ctx, true)
	VisitNode(an.test, v, ctx, true)
	callNodeFn(n, N_STMT_DO_WHILE_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WhileStmt)
	callNodeFn(n, N_STMT_WHILE_BEFORE, v, ctx)
	VisitNode(an.test, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_WHILE_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenForStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForStmt)
	callNodeFn(n, N_STMT_FOR_BEFORE, v, ctx)
	VisitNode(an.init, v, ctx, true)
	VisitNode(an.test, v, ctx, true)
	VisitNode(an.update, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_FOR_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenForInOfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForInOfStmt)
	callNodeFn(n, N_STMT_FOR_IN_OF_BEFORE, v, ctx)
	VisitNode(an.left, v, ctx, true)
	VisitNode(an.right, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_FOR_IN_OF_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenIfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*IfStmt)
	callNodeFn(n, N_STMT_IF_BEFORE, v, ctx)
	VisitNode(an.test, v, ctx, true)
	VisitNode(an.cons, v, ctx, true)
	VisitNode(an.alt, v, ctx, true)
	callNodeFn(n, N_STMT_IF_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenSwitchStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchStmt)
	callNodeFn(n, N_STMT_SWITCH_BEFORE, v, ctx)
	VisitNode(an.test, v, ctx, true)
	VisitNodes(an.cases, v, ctx, true)
	callNodeFn(n, N_STMT_SWITCH_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenBrkStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BrkStmt)
	callNodeFn(n, N_STMT_BRK_BEFORE, v, ctx)
	VisitNode(an.label, v, ctx, true)
	callNodeFn(n, N_STMT_BRK_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenContStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ContStmt)
	callNodeFn(n, N_STMT_CONT_BEFORE, v, ctx)
	VisitNode(an.label, v, ctx, true)
	callNodeFn(n, N_STMT_CONT_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenLabelStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*LabelStmt)
	callNodeFn(n, N_STMT_LABEL_BEFORE, v, ctx)
	VisitNode(an.label, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_LABEL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenRetStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RetStmt)
	callNodeFn(n, N_STMT_RET_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_STMT_RET_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenThrowStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ThrowStmt)
	callNodeFn(n, N_STMT_THROW_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_STMT_THROW_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenTryStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TryStmt)
	callNodeFn(n, N_STMT_TRY_BEFORE, v, ctx)
	VisitNode(an.try, v, ctx, true)
	VisitNode(an.catch, v, ctx, true)
	VisitNode(an.fin, v, ctx, true)
	callNodeFn(n, N_STMT_TRY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenDebugStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_STMT_DEBUG_BEFORE, v, ctx)
	callNodeFn(n, N_STMT_DEBUG_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenWithStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WithStmt)
	callNodeFn(n, N_STMT_WITH_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_WITH_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenClassDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	callNodeFn(n, N_STMT_CLASS_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNode(an.super, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_STMT_CLASS_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenImportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportDec)
	callNodeFn(n, N_STMT_IMPORT_BEFORE, v, ctx)
	VisitNodes(an.specs, v, ctx, true)
	VisitNode(an.src, v, ctx, true)
	callNodeFn(n, N_STMT_IMPORT_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenExportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportDec)
	callNodeFn(n, N_STMT_EXPORT_BEFORE, v, ctx)
	VisitNodes(an.specs, v, ctx, true)
	VisitNode(an.src, v, ctx, true)
	callNodeFn(n, N_STMT_EXPORT_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func ListenNull(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_NULL_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_NULL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenBool(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_BOOL_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_BOOL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenNum(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_NUM_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_NUM_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenStr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_STR_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_STR_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenArr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_ARR_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_ARR_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenObj(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_OBJ_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_OBJ_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenReg(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_LIT_REGEXP_BEFORE, v, ctx)
	callNodeFn(n, N_LIT_REGEXP_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func ListenNewExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*NewExpr)
	callNodeFn(n, N_EXPR_NEW_BEFORE, v, ctx)
	VisitNode(an.callee, v, ctx, true)
	VisitNodes(an.args, v, ctx, true)
	callNodeFn(n, N_EXPR_NEW_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MemberExpr)
	callNodeFn(n, N_EXPR_MEMBER_BEFORE, v, ctx)
	VisitNode(an.obj, v, ctx, true)
	VisitNode(an.prop, v, ctx, true)
	callNodeFn(n, N_EXPR_MEMBER_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenCallExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CallExpr)
	callNodeFn(n, N_EXPR_CALL_BEFORE, v, ctx)
	VisitNode(an.callee, v, ctx, true)
	VisitNodes(an.args, v, ctx, true)
	callNodeFn(n, N_EXPR_CALL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenBinExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BinExpr)
	callNodeFn(n, N_EXPR_BIN_BEFORE, v, ctx)
	VisitNode(an.lhs, v, ctx, true)
	VisitNode(an.rhs, v, ctx, true)
	callNodeFn(n, N_EXPR_BIN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenUnaryExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UnaryExpr)
	callNodeFn(n, N_EXPR_UNARY_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_EXPR_UNARY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenUpdateExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UpdateExpr)
	callNodeFn(n, N_EXPR_UPDATE_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_EXPR_UPDATE_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenCondExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CondExpr)
	callNodeFn(n, N_EXPR_COND_BEFORE, v, ctx)
	VisitNode(an.test, v, ctx, true)
	VisitNode(an.cons, v, ctx, true)
	VisitNode(an.alt, v, ctx, true)
	callNodeFn(n, N_EXPR_COND_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenAssignExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignExpr)
	callNodeFn(n, N_EXPR_ASSIGN_BEFORE, v, ctx)
	VisitNode(an.lhs, v, ctx, true)
	VisitNode(an.rhs, v, ctx, true)
	callNodeFn(n, N_EXPR_ASSIGN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenFnExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	callNodeFn(n, N_EXPR_FN_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNodes(an.params, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_EXPR_FN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenThisExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_EXPR_THIS_BEFORE, v, ctx)
	callNodeFn(n, N_EXPR_THIS_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenParenExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ParenExpr)
	callNodeFn(n, N_EXPR_PAREN_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	callNodeFn(n, N_EXPR_PAREN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenArrowFn(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrowFn)
	callNodeFn(n, N_EXPR_ARROW_BEFORE, v, ctx)
	VisitNodes(an.params, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_EXPR_ARROW_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenSeqExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SeqExpr)
	callNodeFn(n, N_EXPR_SEQ_BEFORE, v, ctx)
	VisitNodes(an.elems, v, ctx, true)
	callNodeFn(n, N_EXPR_SEQ_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenClassExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	callNodeFn(n, N_EXPR_CLASS_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNode(an.super, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_EXPR_CLASS_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenTplExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TplExpr)
	callNodeFn(n, N_EXPR_TPL_BEFORE, v, ctx)
	VisitNode(an.tag, v, ctx, true)
	VisitNodes(an.elems, v, ctx, true)
	callNodeFn(n, N_EXPR_TPL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenYieldExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*YieldExpr)
	callNodeFn(n, N_EXPR_YIELD_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_EXPR_YIELD_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenOptChainExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ChainExpr)
	callNodeFn(n, N_EXPR_CHAIN_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	callNodeFn(n, N_EXPR_CHAIN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxElem(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxElem)
	callNodeFn(n, N_JSX_ELEM_BEFORE, v, ctx)
	VisitNode(an.open, v, ctx, true)
	VisitNodes(an.children, v, ctx, true)
	VisitNode(an.close, v, ctx, true)
	callNodeFn(n, N_JSX_ELEM_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_NAME_BEFORE, v, ctx)
	callNodeFn(n, N_NAME_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenImportCall(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportCall)
	callNodeFn(n, N_IMPORT_CALL_BEFORE, v, ctx)
	VisitNode(an.src, v, ctx, true)
	callNodeFn(n, N_IMPORT_CALL_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenMetaProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MetaProp)
	callNodeFn(n, N_META_PROP_BEFORE, v, ctx)
	VisitNode(an.meta, v, ctx, true)
	VisitNode(an.prop, v, ctx, true)
	callNodeFn(n, N_META_PROP_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Spread)
	callNodeFn(n, N_SPREAD_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_SPREAD_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func ListenVarDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDec)
	callNodeFn(n, N_VAR_DEC_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNode(an.init, v, ctx, true)
	callNodeFn(n, N_VAR_DEC_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenRestPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RestPat)
	callNodeFn(n, N_PAT_REST_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_PAT_REST_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenArrPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrPat)
	callNodeFn(n, N_PAT_ARRAY_BEFORE, v, ctx)
	VisitNodes(an.elems, v, ctx, true)
	callNodeFn(n, N_PAT_ARRAY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenAssignPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignPat)
	callNodeFn(n, N_PAT_ASSIGN_BEFORE, v, ctx)
	VisitNode(an.lhs, v, ctx, true)
	VisitNode(an.rhs, v, ctx, true)
	callNodeFn(n, N_PAT_ASSIGN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenObjPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ObjPat)
	callNodeFn(n, N_PAT_OBJ_BEFORE, v, ctx)
	VisitNodes(an.props, v, ctx, true)
	callNodeFn(n, N_PAT_OBJ_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Prop)
	callNodeFn(n, N_PROP_BEFORE, v, ctx)
	VisitNode(an.key, v, ctx, true)
	VisitNode(an.value, v, ctx, true)
	callNodeFn(n, N_PROP_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenSwitchCase(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchCase)
	callNodeFn(n, N_SWITCH_CASE_BEFORE, v, ctx)
	VisitNode(an.test, v, ctx, true)
	VisitNodes(an.cons, v, ctx, true)
	callNodeFn(n, N_SWITCH_CASE_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenCatch(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Catch)
	callNodeFn(n, N_CATCH_BEFORE, v, ctx)
	VisitNode(an.param, v, ctx, true)
	VisitNode(an.body, v, ctx, true)
	callNodeFn(n, N_CATCH_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenClassBody(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassBody)
	callNodeFn(n, N_ClASS_BODY_BEFORE, v, ctx)
	VisitNodes(an.elems, v, ctx, true)
	callNodeFn(n, N_ClASS_BODY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenStaticBlock(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*StaticBlock)
	callNodeFn(n, N_STATIC_BLOCK_BEFORE, v, ctx)
	VisitNodes(an.body, v, ctx, true)
	callNodeFn(n, N_STATIC_BLOCK_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenMethod(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Method)
	callNodeFn(n, N_METHOD_BEFORE, v, ctx)
	VisitNode(an.key, v, ctx, true)
	VisitNode(an.value, v, ctx, true)
	callNodeFn(n, N_METHOD_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenField(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Field)
	callNodeFn(n, N_FIELD_BEFORE, v, ctx)
	VisitNode(an.key, v, ctx, true)
	VisitNode(an.value, v, ctx, true)
	callNodeFn(n, N_FIELD_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenSuper(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_SUPER_BEFORE, v, ctx)
	callNodeFn(n, N_SUPER_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenImportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportSpec)
	callNodeFn(n, N_IMPORT_SPEC_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNode(an.local, v, ctx, true)
	callNodeFn(n, N_IMPORT_SPEC_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenExportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportSpec)
	callNodeFn(n, N_EXPORT_SPEC_BEFORE, v, ctx)
	VisitNode(an.id, v, ctx, true)
	VisitNode(an.local, v, ctx, true)
	callNodeFn(n, N_EXPORT_SPEC_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func ListenJsxIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_JSX_ID_BEFORE, v, ctx)
	callNodeFn(n, N_JSX_ID_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxMemberExpr)
	callNodeFn(n, N_JSX_MEMBER_BEFORE, v, ctx)
	VisitNode(an.obj, v, ctx, true)
	VisitNode(an.prop, v, ctx, true)
	callNodeFn(n, N_JSX_MEMBER_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxNsName(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxNsName)
	callNodeFn(n, N_JSX_NS_BEFORE, v, ctx)
	VisitNode(an.ns, v, ctx, true)
	VisitNode(an.name, v, ctx, true)
	callNodeFn(n, N_JSX_NS_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxAttrSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxSpreadAttr)
	callNodeFn(n, N_JSX_ATTR_SPREAD_BEFORE, v, ctx)
	VisitNode(an.arg, v, ctx, true)
	callNodeFn(n, N_JSX_ATTR_SPREAD_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxChildSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxSpreadChild)
	callNodeFn(n, N_JSX_CHILD_SPREAD_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	callNodeFn(n, N_JSX_CHILD_SPREAD_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxOpen(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxOpen)
	callNodeFn(n, N_JSX_OPEN_BEFORE, v, ctx)
	VisitNode(an.name, v, ctx, true)
	VisitNodes(an.attrs, v, ctx, true)
	callNodeFn(n, N_JSX_OPEN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxClose(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxClose)
	callNodeFn(n, N_JSX_CLOSE_BEFORE, v, ctx)
	VisitNode(an.name, v, ctx, true)
	callNodeFn(n, N_JSX_CLOSE_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxEmpty(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_JSX_EMPTY_BEFORE, v, ctx)
	callNodeFn(n, N_JSX_EMPTY_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxExprSpan(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxExprSpan)
	callNodeFn(n, N_JSX_EXPR_SPAN_BEFORE, v, ctx)
	VisitNode(an.expr, v, ctx, true)
	callNodeFn(n, N_JSX_EXPR_SPAN_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxText(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	callNodeFn(n, N_JSX_TXT_BEFORE, v, ctx)
	callNodeFn(n, N_JSX_TXT_AFTER, v, ctx)
	return TRAVERSE_CONT
}
func ListenJsxAttr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JsxAttr)
	callNodeFn(n, N_JSX_ATTR_BEFORE, v, ctx)
	VisitNode(an.name, v, ctx, true)
	VisitNode(an.val, v, ctx, true)
	callNodeFn(n, N_JSX_ATTR_AFTER, v, ctx)
	return TRAVERSE_CONT
}

func init() {
	DefaultListenerImpl[N_PROG] = ListenProg

	DefaultListenerImpl[N_STMT_EXPR] = ListenExprStmt
	DefaultListenerImpl[N_STMT_EMPTY] = ListenExprStmt
	DefaultListenerImpl[N_STMT_VAR_DEC] = ListenVarDecStmt
	DefaultListenerImpl[N_STMT_FN] = ListenFnDec
	DefaultListenerImpl[N_STMT_BLOCK] = ListenBlockStmt
	DefaultListenerImpl[N_STMT_DO_WHILE] = ListenDoWhileStmt
	DefaultListenerImpl[N_STMT_WHILE] = ListenDoWhileStmt
	DefaultListenerImpl[N_STMT_FOR] = ListenForStmt
	DefaultListenerImpl[N_STMT_FOR_IN_OF] = ListenForInOfStmt
	DefaultListenerImpl[N_STMT_IF] = ListenIfStmt
	DefaultListenerImpl[N_STMT_SWITCH] = ListenSwitchStmt
	DefaultListenerImpl[N_STMT_BRK] = ListenBrkStmt
	DefaultListenerImpl[N_STMT_CONT] = ListenContStmt
	DefaultListenerImpl[N_STMT_LABEL] = ListenLabelStmt
	DefaultListenerImpl[N_STMT_RET] = ListenRetStmt
	DefaultListenerImpl[N_STMT_THROW] = ListenThrowStmt
	DefaultListenerImpl[N_STMT_TRY] = ListenTryStmt
	DefaultListenerImpl[N_STMT_DEBUG] = ListenDebugStmt
	DefaultListenerImpl[N_STMT_WITH] = ListenWithStmt
	DefaultListenerImpl[N_STMT_CLASS] = ListenClassDec
	DefaultListenerImpl[N_STMT_IMPORT] = ListenImportDec
	DefaultListenerImpl[N_STMT_EXPORT] = ListenExportDec

	DefaultListenerImpl[N_LIT_NULL] = ListenNull
	DefaultListenerImpl[N_LIT_BOOL] = ListenBool
	DefaultListenerImpl[N_LIT_NUM] = ListenNum
	DefaultListenerImpl[N_LIT_STR] = ListenStr
	DefaultListenerImpl[N_LIT_ARR] = ListenArr
	DefaultListenerImpl[N_LIT_OBJ] = ListenObj
	DefaultListenerImpl[N_LIT_REGEXP] = ListenReg

	DefaultListenerImpl[N_EXPR_NEW] = ListenNewExpr
	DefaultListenerImpl[N_EXPR_MEMBER] = ListenMemberExpr
	DefaultListenerImpl[N_EXPR_CALL] = ListenCallExpr
	DefaultListenerImpl[N_EXPR_BIN] = ListenBinExpr
	DefaultListenerImpl[N_EXPR_UNARY] = ListenUnaryExpr
	DefaultListenerImpl[N_EXPR_UPDATE] = ListenUpdateExpr
	DefaultListenerImpl[N_EXPR_COND] = ListenCondExpr
	DefaultListenerImpl[N_EXPR_ASSIGN] = ListenAssignExpr
	DefaultListenerImpl[N_EXPR_FN] = ListenFnExpr
	DefaultListenerImpl[N_EXPR_THIS] = ListenThisExpr
	DefaultListenerImpl[N_EXPR_PAREN] = ListenParenExpr
	DefaultListenerImpl[N_EXPR_ARROW] = ListenArrowFn
	DefaultListenerImpl[N_EXPR_SEQ] = ListenSeqExpr
	DefaultListenerImpl[N_EXPR_CLASS] = ListenClassExpr
	DefaultListenerImpl[N_EXPR_TPL] = ListenTplExpr
	DefaultListenerImpl[N_EXPR_YIELD] = ListenYieldExpr
	DefaultListenerImpl[N_EXPR_CHAIN] = ListenOptChainExpr
	DefaultListenerImpl[N_JSX_ELEM] = ListenJsxElem
	DefaultListenerImpl[N_NAME] = ListenIdent
	DefaultListenerImpl[N_IMPORT_CALL] = ListenImportCall
	DefaultListenerImpl[N_META_PROP] = ListenMetaProp
	DefaultListenerImpl[N_SPREAD] = ListenSpread

	DefaultListenerImpl[N_VAR_DEC] = ListenVarDec
	DefaultListenerImpl[N_PAT_REST] = ListenRestPat
	DefaultListenerImpl[N_PAT_ARRAY] = ListenArrPat
	DefaultListenerImpl[N_PAT_ASSIGN] = ListenAssignPat
	DefaultListenerImpl[N_PAT_OBJ] = ListenObjPat
	DefaultListenerImpl[N_PROP] = ListenProp
	DefaultListenerImpl[N_SWITCH_CASE] = ListenSwitchCase
	DefaultListenerImpl[N_CATCH] = ListenCatch
	DefaultListenerImpl[N_ClASS_BODY] = ListenClassBody
	DefaultListenerImpl[N_STATIC_BLOCK] = ListenStaticBlock
	DefaultListenerImpl[N_METHOD] = ListenMethod
	DefaultListenerImpl[N_FIELD] = ListenField
	DefaultListenerImpl[N_SUPER] = ListenSuper
	DefaultListenerImpl[N_IMPORT_SPEC] = ListenImportSpec
	DefaultListenerImpl[N_EXPORT_SPEC] = ListenExportSpec

	DefaultListenerImpl[N_JSX_ID] = ListenJsxIdent
	DefaultListenerImpl[N_JSX_MEMBER] = ListenJsxMemberExpr
	DefaultListenerImpl[N_JSX_NS] = ListenJsxNsName
	DefaultListenerImpl[N_JSX_ATTR_SPREAD] = ListenJsxAttrSpread
	DefaultListenerImpl[N_JSX_CHILD_SPREAD] = ListenJsxChildSpread
	DefaultListenerImpl[N_JSX_OPEN] = ListenJsxOpen
	DefaultListenerImpl[N_JSX_CLOSE] = ListenJsxClose
	DefaultListenerImpl[N_JSX_EMPTY] = ListenJsxEmpty
	DefaultListenerImpl[N_JSX_EXPR_SPAN] = ListenJsxExprSpan
	DefaultListenerImpl[N_JSX_TXT] = ListenJsxText
	DefaultListenerImpl[N_JSX_ATTR] = ListenJsxAttr
}

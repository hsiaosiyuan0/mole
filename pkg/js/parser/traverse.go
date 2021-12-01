package parser

import "log"

type TraverseCtx struct {
	RecordDepth ContOrStop
	Depth       int

	RecordPath ContOrStop
	Path       []string
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

// the type of `v` is `func VisitorImpl`
type VisitFn = func(n Node, v interface{}, ctx *TraverseCtx) ContOrStop

// for the `func VisitFn` can be routed in O(1)
type VisitorImpl = [N_NODE_DEF_END]VisitFn

func VisitNode(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	if n == nil {
		return TRAVERSE_CONT
	}

	vi := v.(*VisitorImpl)
	typ := n.Type()
	fn := vi[typ]
	if fn == nil {
		log.Fatalf("Impl does not exist for NodeType %d", typ)
	}

	return fn(n, v, ctx)
}

func VisitNodes(ns []Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	vi := v.(*VisitorImpl)
	for _, n := range ns {
		if !VisitNode(n, vi, ctx) {
			return TRAVERSE_STOP
		}
	}
	return TRAVERSE_CONT
}

var Defaultfunc VisitorImpl = [N_NODE_DEF_END]VisitFn{
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
	return VisitNodes(an.stmts, v, ctx)
}

func VisitExprStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExprStmt)
	return VisitNode(an.expr, v, ctx)
}
func VisitEmptyStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitVarDecStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDecStmt)
	return VisitNodes(an.decList, v, ctx)
}
func VisitFnDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.params, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitBlockStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BlockStmt)
	return VisitNodes(an.body, v, ctx)
}
func VisitDoWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*DoWhileStmt)
	if !VisitNode(an.body, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.test, v, ctx)
}
func VisitWhileStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WhileStmt)
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitForStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForStmt)
	if !VisitNode(an.init, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.update, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitForInOfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ForInOfStmt)
	if !VisitNode(an.left, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.right, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitIfStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*IfStmt)
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.cons, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.alt, v, ctx)
}
func VisitSwitchStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchStmt)
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.cases, v, ctx)
}
func VisitBrkStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BrkStmt)
	return VisitNode(an.label, v, ctx)
}
func VisitContStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ContStmt)
	return VisitNode(an.label, v, ctx)
}
func VisitLabelStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*LabelStmt)
	if !VisitNode(an.label, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitRetStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RetStmt)
	return VisitNode(an.arg, v, ctx)
}
func VisitThrowStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ThrowStmt)
	return VisitNode(an.arg, v, ctx)
}
func VisitTryStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TryStmt)
	if !VisitNode(an.try, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.catch, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.fin, v, ctx)
}
func VisitDebugStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitWithStmt(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*WithStmt)
	if !VisitNode(an.expr, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitClassDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.super, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitImportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportDec)
	if !VisitNodes(an.specs, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.src, v, ctx)
}
func VisitExportDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportDec)
	if !VisitNodes(an.specs, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.src, v, ctx)
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
	if !VisitNode(an.callee, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.args, v, ctx)
}
func VisitMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MemberExpr)
	if !VisitNode(an.obj, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx)
}
func VisitCallExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CallExpr)
	if !VisitNode(an.callee, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.args, v, ctx)
}
func VisitBinExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*BinExpr)
	if !VisitNode(an.lhs, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx)
}
func VisitUnaryExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UnaryExpr)
	return VisitNode(an.arg, v, ctx)
}
func VisitUpdateExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*UpdateExpr)
	return VisitNode(an.arg, v, ctx)
}
func VisitCondExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*CondExpr)
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.cons, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.alt, v, ctx)
}
func VisitAssignExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignExpr)
	if !VisitNode(an.lhs, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx)
}
func VisitFnExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*FnDec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.params, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitThisExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitParenExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ParenExpr)
	return VisitNode(an.expr, v, ctx)
}
func VisitArrowFn(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrowFn)
	if !VisitNodes(an.params, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitSeqExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SeqExpr)
	return VisitNodes(an.elems, v, ctx)
}
func VisitClassExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassDec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNode(an.super, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitTplExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*TplExpr)
	if !VisitNode(an.tag, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.elems, v, ctx)
}
func VisitYieldExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*YieldExpr)
	return VisitNode(an.arg, v, ctx)
}
func VisitOptChainExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ChainExpr)
	return VisitNode(an.expr, v, ctx)
}
func VisitJsxElem(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXElem)
	if !VisitNode(an.open, v, ctx) {
		return TRAVERSE_STOP
	}
	if !VisitNodes(an.children, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.close, v, ctx)
}
func VisitIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitImportCall(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportCall)
	return VisitNode(an.src, v, ctx)
}
func VisitMetaProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*MetaProp)
	if !VisitNode(an.meta, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx)
}
func VisitSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Spread)
	return VisitNode(an.arg, v, ctx)
}

func VisitVarDec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*VarDec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.init, v, ctx)
}
func VisitRestPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*RestPat)
	return VisitNode(an.arg, v, ctx)
}
func VisitArrPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ArrPat)
	return VisitNodes(an.elems, v, ctx)
}
func VisitAssignPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*AssignPat)
	if !VisitNode(an.lhs, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.rhs, v, ctx)
}
func VisitObjPat(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ObjPat)
	return VisitNodes(an.props, v, ctx)
}
func VisitProp(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Prop)
	if !VisitNode(an.key, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx)
}
func VisitSwitchCase(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*SwitchCase)
	if !VisitNode(an.test, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.cons, v, ctx)
}
func VisitCatch(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Catch)
	if !VisitNode(an.param, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.body, v, ctx)
}
func VisitClassBody(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ClassBody)
	return VisitNodes(an.elems, v, ctx)
}
func VisitStaticBlock(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*StaticBlock)
	return VisitNodes(an.body, v, ctx)
}
func VisitMethod(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Method)
	if !VisitNode(an.key, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx)
}
func VisitField(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*Field)
	if !VisitNode(an.key, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.value, v, ctx)
}
func VisitSuper(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitImportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ImportSpec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.local, v, ctx)
}
func VisitExportSpec(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*ExportSpec)
	if !VisitNode(an.id, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.local, v, ctx)
}

func VisitJsxIdent(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxMemberExpr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXMemberExpr)
	if !VisitNode(an.obj, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.prop, v, ctx)
}
func VisitJsxNsName(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXNsName)
	if !VisitNode(an.ns, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.name, v, ctx)
}
func VisitJsxAttrSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXSpreadAttr)
	return VisitNode(an.arg, v, ctx)
}
func VisitJsxChildSpread(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXSpreadChild)
	return VisitNode(an.expr, v, ctx)
}
func VisitJsxOpen(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXOpen)
	if !VisitNode(an.name, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNodes(an.attrs, v, ctx)
}
func VisitJsxClose(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXClose)
	return VisitNode(an.name, v, ctx)
}
func VisitJsxEmpty(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxExprSpan(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXExprSpan)
	return VisitNode(an.expr, v, ctx)
}
func VisitJsxText(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	return TRAVERSE_CONT
}
func VisitJsxAttr(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
	an := n.(*JSXAttr)
	if !VisitNode(an.name, v, ctx) {
		return TRAVERSE_STOP
	}
	return VisitNode(an.val, v, ctx)
}

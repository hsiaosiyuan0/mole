package analysis

import (
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

const (
	N_CFG_DEBUG parser.NodeType = parser.N_NODE_DEF_END + 1 + iota
)

type AnalysisCtx struct {
	graphMap map[parser.Node]*Graph
	root     *Graph
	graph    *Graph

	stmtStack []*Block
	exprStack []*Block

	// map astNode to its basic block
	astNodeToBlock map[uint64]Block
}

func newAnalysisCtx() *AnalysisCtx {
	root := newGraph()
	a := &AnalysisCtx{
		graphMap:       map[parser.Node]*Graph{},
		root:           root,
		graph:          root,
		stmtStack:      make([]*Block, 0),
		exprStack:      make([]*Block, 0),
		astNodeToBlock: map[uint64]Block{},
	}
	a.stmtStack = append(a.stmtStack, a.graph.Head)
	return a
}

func (a *AnalysisCtx) GraphOf(node parser.Node) *Graph {
	g, ok := a.graphMap[node]
	if !ok {
		return nil
	}
	return g
}

func (a *AnalysisCtx) enterGraph(node parser.Node) {
	graph := newGraph()
	graph.Parent = a.graph

	a.graphMap[node] = graph
	a.graph = graph
}

func (a *AnalysisCtx) leaveGraph() {
	a.graph = a.graph.Parent
}

func (a *AnalysisCtx) lastExpr() *Block {
	return a.exprStack[len(a.exprStack)-1]
}

func (a *AnalysisCtx) lastStmt() *Block {
	return a.stmtStack[len(a.stmtStack)-1]
}

func (a *AnalysisCtx) newEnter(astNode parser.Node, info string) *Block {
	b := a.graph.newBasicBlk()
	b.Nodes = append(b.Nodes, newInfoNode(astNode, true, info))
	return b
}

func (a *AnalysisCtx) newExit(astNode parser.Node, info string) *Block {
	b := a.graph.newBasicBlk()
	b.Nodes = append(b.Nodes, newInfoNode(astNode, false, info))
	return b
}

func (a *AnalysisCtx) pushStmt(b *Block) {
	a.stmtStack = append(a.stmtStack, b)
}

func (a *AnalysisCtx) popStmt() *Block {
	cnt := len(a.stmtStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.stmtStack[cnt-1], a.stmtStack[:cnt-1]
	a.stmtStack = rest
	return last
}

func (a *AnalysisCtx) pushExpr(b *Block) {
	a.exprStack = append(a.exprStack, b)
}

func (a *AnalysisCtx) popExpr() *Block {
	cnt := len(a.exprStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.exprStack[cnt-1], a.exprStack[:cnt-1]
	a.exprStack = rest
	return last
}

func (a *AnalysisCtx) popExprsAndLink(cnt int) (*Block, *Block) {
	var head, tail, th, tt *Block
	if cnt == 1 {
		head = a.popExpr()
		tail = head
	} else if cnt > 1 {
		tail = a.popExpr()
		tt = tail
		for i := cnt - 2; i >= 0; i-- {
			th = a.popExpr()
			link(a, th, EK_JMP, ET_NONE, EK_JMP, ET_NONE, tt, LF_NONE)
			link(a, th, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, tt, LF_NONE)
			tt = th
		}
		head = th
	}
	return head, tail
}

type Analysis struct {
	WalkCtx *walk.WalkCtx
}

func NewAnalysis(root parser.Node, symtab *parser.SymTab) *Analysis {
	a := &Analysis{
		WalkCtx: walk.NewWalkCtx(root, symtab),
	}
	a.init()
	return a
}

func (a *Analysis) AnalysisCtx() *AnalysisCtx {
	return analysisCtx(a.WalkCtx.VisitorCtx())
}

func (a *Analysis) Graph() *Graph {
	return analysisCtx(a.WalkCtx.VisitorCtx()).graph
}

func isLoop(t parser.NodeType) bool {
	return t == parser.N_STMT_FOR || t == parser.N_STMT_WHILE || t == parser.N_STMT_DO_WHILE || t == parser.N_STMT_FOR_IN_OF
}

func isAtom(t parser.NodeType) bool {
	_, ok := walk.AtomNodeTypes[t]
	return ok
}

func isFn(t parser.NodeType) bool {
	return t == parser.N_STMT_FN || t == parser.N_EXPR_FN || t == parser.N_EXPR_ARROW
}

func pushAtomNode(node parser.Node, key string, ctx *walk.VisitorCtx) *Block {
	ac := analysisCtx(ctx)
	b := ac.graph.newBasicBlk()
	b.Nodes = append(b.Nodes, node)
	ac.pushExpr(b)
	return b
}

func handleBefore(node parser.Node, key string, ctx *walk.VisitorCtx) {
	ac := analysisCtx(ctx)

	astTyp := node.Type()
	pAstTyp := ctx.ParentNodeType()

	var blk *Block
	if isAtom(astTyp) {
		blk = pushAtomNode(node, key, ctx)
	} else {
		blk = ac.newEnter(node, "")
	}

	// record the label to its target scope
	if (astTyp.IsExpr() || astTyp.IsStmt() && astTyp == parser.N_STMT_EXPR) && pAstTyp != parser.N_STMT_LABEL {
		if cnt := len(ac.graph.hangingLabels); cnt > 0 {
			last, rest := ac.graph.hangingLabels[cnt-1], ac.graph.hangingLabels[:cnt-1]
			ac.graph.hangingLabels = rest
			ac.graph.labelLoop[last] = ctx.ScopeId()
		}
	}

	// stmt
	if astTyp.IsStmt() || astTyp == parser.N_PROG {

		if astTyp == parser.N_STMT_LABEL && node.(*parser.LabelStmt).Used() {
			ac.graph.hangingLabels = append(ac.graph.hangingLabels, node)
		} else if astTyp == parser.N_STMT_FN || astTyp == parser.N_EXPR_FN {
			enterFnGraph(node, ctx, true)
		}

		if pAstTyp == parser.N_STMT_IF || isLoop(pAstTyp) || pAstTyp == parser.N_STMT_LABEL ||
			pAstTyp == parser.N_STMT_TRY || pAstTyp == parser.N_CATCH || pAstTyp == parser.N_STMT_WITH ||
			pAstTyp == parser.N_STMT_EXPORT {
			// just pushing `enter` into stack without linking to the `prevStmt` to imitate a new branch
			// the forked branch will be linked to its source branch point in the post-process listener
			// of its astNode
			ac.pushStmt(blk)
		} else if astTyp != parser.N_STMT_FN && astTyp != parser.N_EXPR_FN {
			prev := ac.popStmt()
			scope := ctx.Scope()
			flag := LF_NONE
			if astTyp == parser.N_STMT_TRY && (scope.IsKind(parser.SPK_TRY) || scope.IsKind(parser.SPK_TRY_INDIRECT)) {
				flag |= LF_FORCE_SEP
			}
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, blk, flag)
			ac.pushStmt(grpBlock(ac, prev, blk))
		}
	} else if !isAtom(astTyp) {
		if isFn(astTyp) {
			enterFnGraph(node, ctx, false)
		} else if astTyp == parser.N_SWITCH_CASE {
			// here we push en empty block into the stack to start a new fork for the case clause
			// so its body will grow from this empty block. the new empty block does not link to
			// previous block yet at here, the connection will happen in the post-case handler
			b := ac.graph.newBasicBlk()
			ac.pushStmt(b)
		} else {
			ac.pushExpr(blk)
		}
	}
}

func enterFnGraph(node parser.Node, ctx *walk.VisitorCtx, stmt bool) {
	ac := analysisCtx(ctx)

	// reflects the graph entry in parent graph, parent graph only hold
	// the entry node to represent that there is a sub-graph fork from
	// that place
	b := ac.graph.newBasicBlk()
	b.Nodes = append(b.Nodes, newInfoNode(node, true, node.Type().String()))

	if stmt {
		// the info node should be connected to previous block if its origin is fnStmt
		// that's because the link strategy of stmt
		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, b, LF_NONE)
		ac.pushStmt(grpBlock(ac, prev, b))
	} else {
		// as opposed to the link strategy of stmt, link-to-previous process is not
		// needed in expr situation, the connection will be done in its parent expr
		// process
		ac.pushExpr(b)
	}

	// enter new graph
	ac.enterGraph(node)

	// use an empty block as header since we will push fn id and params back to header
	// the mimic header will be used the start of function body, so `pushStmt` is needed
	// below
	ac.graph.Head = ac.graph.newBasicBlk()
	ac.pushStmt(ac.graph.Head)
}

func handleAfter(node parser.Node, key string, ctx *walk.VisitorCtx) {
	ac := analysisCtx(ctx)

	astTyp := node.Type()
	pAstTyp := ctx.ParentNodeType()
	switch astTyp {

	case parser.N_EXPR_BIN:
		n := node.(*parser.BinExpr)

		rhs := ac.popExpr()
		lhs := ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, lhs, LF_NONE)
		lhs = grpBlock(ac, enter, lhs)

		op := n.Op()
		logic := true
		if op == parser.T_AND {
			lhs.newJmpOut(ET_JMP_F)
		} else if op == parser.T_OR || op == parser.T_NULLISH {
			lhs.newJmpOut(ET_JMP_T)
		} else {
			logic = false
		}

		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, rhs, LF_NONE)

		if op == parser.T_OR || op == parser.T_NULLISH {
			link(ac, lhs, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, rhs, LF_NONE)
		}

		vn := ac.graph.newGroupBlk()
		vn.Inlets = lhs.Inlets
		if logic {
			vn.Outlets = append(rhs.Outlets, lhs.xOutEdges()...)
		} else {
			vn.Outlets = rhs.Outlets
		}

		exit := ac.newExit(node, "")
		if logic && pAstTyp == parser.N_STMT_EXPR {
			link(ac, vn, EK_JMP, ET_JMP_T|ET_JMP_F, EK_JMP, ET_NONE, exit, LF_NONE)
		}
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, vn, exit))

	case parser.N_EXPR_UNARY:
		expr := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_NONE)

		prev := grpBlock(ac, enter, expr)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_COND:
		alt := ac.popExpr()
		cons := ac.popExpr()
		test := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		test.newJmpOut(ET_JMP_F)
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)

		prev := grpBlock(ac, enter, test)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, cons, LF_NONE)
		link(ac, prev, EK_JMP, ET_JMP_T, EK_SEQ, ET_NONE, cons, LF_NONE)
		link(ac, prev, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, alt, LF_NONE)

		vn := grpBlock(ac, prev, cons)
		vn.Outlets = append(vn.Outlets, alt.xOutEdges()...)

		link(ac, vn, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_FORCE_SEP)
		ac.pushExpr(grpBlock(ac, vn, exit))

	case parser.N_IMPORT_CALL:
		expr := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_NONE)

		prev := grpBlock(ac, enter, expr)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_UPDATE, parser.N_EXPR_PAREN:
		expr := ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_NONE)

		exit := ac.newExit(node, "")
		link(ac, expr, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, expr, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_LIT_ARR:
		n := node.(*parser.ArrLit)

		head, tail := ac.popExprsAndLink(len(n.Elems()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_PROP:
		n := node.(*parser.Prop)
		var key, val *Block
		if n.Val() != nil {
			val = ac.popExpr()
		}
		key = ac.popExpr()

		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, key, LF_NONE)

		exit := ac.newExit(node, "")
		if val != nil {
			link(ac, key, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, val, LF_NONE)

			link(ac, val, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, val, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, key, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_LIT_OBJ:
		n := node.(*parser.ObjLit)

		head, tail := ac.popExprsAndLink(len(n.Props()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_ASSIGN:
		rhs := ac.popExpr()
		lhs := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, lhs, LF_NONE)
		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, rhs, LF_NONE)

		link(ac, rhs, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, rhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_CALL:
		n := node.(*parser.CallExpr)

		head, tail := ac.popExprsAndLink(len(n.Args()))

		fn := ac.popExpr()
		if n.Optional() {
			fn.newJmpOut(ET_JMP_F)
		}

		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fn, LF_NONE)
			link(ac, fn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fn, LF_NONE)
			link(ac, fn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_NEW:
		n := node.(*parser.NewExpr)

		head, tail := ac.popExprsAndLink(len(n.Args()))

		fn := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fn, LF_NONE)

			prev := grpBlock(ac, enter, fn)
			link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, head, LF_NONE)
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fn, LF_NONE)

			prev := grpBlock(ac, enter, fn)
			link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_MEMBER:
		n := node.(*parser.MemberExpr)

		prop := ac.popExpr()
		obj := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		if n.Optional() {
			obj.newJmpOut(ET_JMP_F)
		}

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, obj, LF_NONE)

		prev := grpBlock(ac, enter, obj)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, prop, LF_NONE)

		prev = grpBlock(ac, prev, prop)
		link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_EXPR_SEQ:
		n := node.(*parser.SeqExpr)

		head, tail := ac.popExprsAndLink(len(n.Elems()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_PAT_ARRAY:
		n := node.(*parser.ArrPat)

		cnt := 0
		for _, el := range n.Elems() {
			if el != nil {
				cnt += 1
			}
		}
		head, tail := ac.popExprsAndLink(cnt)

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_PAT_ASSIGN:
		rhs := ac.popExpr()
		lhs := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, lhs, LF_NONE)

		lhs.newJmpOut(ET_JMP_F)
		link(ac, lhs, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, rhs, LF_NONE)
		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		link(ac, rhs, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, rhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_FORCE_SEP)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_PAT_OBJ:
		n := node.(*parser.ObjPat)

		head, tail := ac.popExprsAndLink(len(n.Props()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_PAT_REST, parser.N_SPREAD, parser.N_EXPR_YIELD, parser.N_EXPR_CHAIN:
		expr := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_NONE)

		prev := grpBlock(ac, enter, expr)
		link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_STMT_EXPR:
		expr := ac.popExpr()
		enter := ac.popStmt()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_FORCE_JOIN)

		prev := grpBlock(ac, enter, expr)
		exit := ac.newExit(node, "")

		link(ac, prev, EK_JMP, ET_JMP_T|ET_JMP_F, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_VAR_DEC:
		n := node.(*parser.VarDec)
		var id, init *Block
		if n.Init() != nil {
			init = ac.popExpr()
		}
		id = ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, id, LF_NONE)

		exit := ac.newExit(node, "")
		if init != nil {
			link(ac, id, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, init, LF_NONE)

			link(ac, init, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
			link(ac, init, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, id, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}
		vn := grpBlock(ac, enter, exit)

		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, vn, LF_NONE)
		ac.pushStmt(grpBlock(ac, prev, vn))

	case parser.N_STMT_VAR_DEC:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushStmt(grpBlock(ac, prev, exit))

		if pAstTyp == parser.N_STMT_FOR && key == "Init" ||
			pAstTyp == parser.N_STMT_FOR_IN_OF && key == "Left" ||
			pAstTyp == parser.N_STMT_EXPORT {
			ac.pushExpr(ac.popStmt())
		}

	case parser.N_PROG:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		// TODO:
		// forceSep := astTyp == parser.N_PROG && ac.graph.hasHangingThrow

		link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushStmt(grpBlock(ac, prev, exit))

	case parser.N_STMT_BLOCK:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushStmt(grpBlock(ac, prev, exit))

	case parser.N_STMT_IF:
		n := node.(*parser.IfStmt)
		var enter, test, cons, alt *Block
		if n.Alt() != nil {
			alt = ac.popStmt()
		}
		cons = ac.popStmt()
		test = ac.popExpr()
		enter = ac.popStmt()

		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, test, LF_NONE)
		test = grpBlock(ac, enter, test)

		test.newJmpOut(ET_JMP_F)
		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, cons, LF_NONE)
		link(ac, test, EK_JMP, ET_JMP_T, EK_SEQ, ET_NONE, cons, LF_NONE)
		if alt != nil {
			link(ac, test, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, alt, LF_NONE)
		}

		vn := ac.graph.newGroupBlk()
		vn.Inlets = enter.Inlets

		vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		vn.Outlets = append(vn.Outlets, cons.xOutEdges()...)
		if alt != nil {
			vn.Outlets = append(vn.Outlets, alt.xOutEdges()...)
		}

		exit := ac.newExit(node, "")
		link(ac, vn, EK_JMP, ET_JMP_T|ET_JMP_F, EK_JMP, ET_NONE, exit, LF_NONE)

		flag := LF_NONE
		if n.Alt() != nil {
			flag |= LF_FORCE_SEP
		}
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, flag)

		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_FOR:
		n := node.(*parser.ForStmt)
		var enter, init, test, update, body *Block
		body = ac.popStmt()
		if n.Update() != nil {
			update = ac.popExpr()
		}
		if n.Test() != nil {
			test = ac.popExpr()
			test.newJmpOut(ET_JMP_F)
			test.newLoopIn()
		}
		if n.Init() != nil {
			init = ac.popExpr()
		}
		enter = ac.popStmt()

		if init != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, init, LF_NONE)
			link(ac, init, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)
		}

		if test != nil {
			if init == nil {
				link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, test, LF_NONE)
			}
			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
		} else {
			body.newLoopIn()
			if init == nil {
				link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			} else {
				link(ac, init, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			}
		}

		if update != nil {
			if ac.graph.isLoopHasCont(node) {
				update.newLoopIn()
			}
			link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, update, LF_NONE)
			if test != nil {
				link(ac, update, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)
			} else {
				link(ac, update, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			}
			update.mrkSeqOutAsLoop()
		} else {
			if test != nil {
				link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)
			} else {
				link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			}
			body.mrkSeqOutAsLoop()
			if test == nil {
				body.addCutOutEdge()
			}
		}

		tailBlk := update
		if tailBlk == nil {
			tailBlk = test
		}
		if tailBlk == nil {
			tailBlk = body
		}
		resolveCont(ac, node, tailBlk)

		vn := ac.graph.newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		if tn := n.Test(); tn != nil && tn.Type().IsLit() {
			j := test.FindOutEdge(EK_JMP, ET_JMP_F, false)
			j.Kind = EK_SEQ
			j.Tag = ET_CUT
		}

		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_WHILE:
		n := node.(*parser.WhileStmt)
		body := ac.popStmt()
		test := ac.popExpr()
		enter := ac.popStmt()

		if !n.Test().Type().IsLit() {
			test.newJmpOut(ET_JMP_F)
		}

		test.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)

		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
		link(ac, test, EK_JMP, ET_JMP_T, EK_SEQ, ET_NONE, body, LF_NONE)

		link(ac, body, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, test, LF_NONE)
		link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)
		body.mrkSeqOutAsLoop()

		if n.Test().Type().IsLit() {
			body.addCutOutEdge()
		}

		resolveCont(ac, node, test)

		vn := ac.graph.newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_DO_WHILE:
		n := node.(*parser.DoWhileStmt)
		test := ac.popExpr()
		body := ac.popStmt()
		enter := ac.popStmt()

		test.newJmpOut(ET_JMP_F)
		if ac.graph.isLoopHasCont(n) {
			test.newJmpIn(ET_LOOP)
		}

		body.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, body, LF_NONE)
		link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)

		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
		test.mrkSeqOutAsLoop()
		link(ac, test, EK_JMP, ET_JMP_T, EK_JMP, ET_NONE, body, LF_NONE)
		test.mrkJmpOutAsLoop(ET_JMP_T)

		resolveCont(ac, node, test)

		vn := ac.graph.newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		if n.Test().Type().IsLit() {
			j := test.FindOutEdge(EK_JMP, ET_JMP_F, false)
			j.Kind = EK_SEQ
			j.Tag = ET_CUT
		}

		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_FOR_IN_OF:
		body := ac.popStmt()
		enter := ac.popStmt()
		rhs := ac.popExpr()
		lhs := ac.popExpr()

		lhs.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, lhs, LF_NONE)
		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, rhs, LF_NONE)

		resolveCont(ac, node, lhs)

		rhs.newJmpOut(ET_JMP_F)
		link(ac, rhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)

		body.mrkSeqOutAsLoop()
		link(ac, body, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, lhs, LF_NONE)

		exit := ac.newExit(node, "")
		link(ac, rhs, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_CATCH:
		_ = node.(*parser.Catch)
		id := ac.popExpr()
		enter := ac.popExpr()
		body := ac.popStmt()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, id, LF_NONE)
		link(ac, id, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
		link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_TRY:
		n := node.(*parser.TryStmt)

		var fin *Block
		if n.Fin() != nil {
			fin = ac.popStmt()
		}

		catch := ac.popStmt()
		try := ac.popStmt()
		enter := ac.popStmt()
		exit := ac.newExit(node, "")
		origExit := exit

		eb := []*Edge{}
		IterBlock(try.unwrapSeqIn(), func(blk *Block) {
			if edge := blk.FindOutEdge(EK_JMP, ET_JMP_E, false); edge != nil || blk.allInfoNode() || blk.throwLit() {
				return
			}
			edge := blk.newJmpOut(ET_JMP_E)
			eb = append(eb, edge)
		})
		try.Outlets = append(try.Outlets, eb...)

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, try, LF_NONE)
		link(ac, try, EK_JMP, ET_JMP_E, EK_SEQ, ET_NONE, catch, LF_NONE)

		throws := ac.graph.hangingThrow[n]
		if len(throws) > 0 {
			catch.newJmpIn(ET_JMP_U)
		}
		for _, blk := range throws {
			link(ac, blk, EK_JMP, ET_JMP_U, EK_JMP, ET_JMP_U, catch, LF_NONE)
		}

		// `ET_JMP_U` inside catch should jump to the fin while skipping the stmts after fin
		if edges := FindEdges(catch.Outlets, EK_JMP, ET_JMP_U); len(edges) > 0 && fin != nil {
			link(ac, catch, EK_JMP, ET_JMP_U, EK_JMP, ET_NONE, fin, LF_FORCE_SEP)
			for _, edge := range edges {
				edge.Tag |= ET_JMP_P
			}
			fin.newJmpOut(ET_JMP_P)
		}

		if edges := FindEdges(try.Outlets, EK_JMP, ET_JMP_U); len(edges) > 0 && fin != nil {
			link(ac, try, EK_JMP, ET_JMP_U, EK_JMP, ET_NONE, fin, LF_FORCE_SEP)
			for _, edge := range edges {
				edge.Tag |= ET_JMP_P
			}
			if edge := fin.FindOutEdge(EK_JMP, ET_JMP_U, false); edge == nil {
				fin.newJmpOut(ET_JMP_P)
			}
		}

		if fin != nil {
			link(ac, try, EK_JMP, ET_JMP_P, EK_SEQ, ET_NONE, fin, LF_NONE)
			link(ac, catch, EK_JMP, ET_JMP_P, EK_SEQ, ET_NONE, fin, LF_NONE)
			link(ac, fin, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
			exit = fin
		}

		exit.newJmpIn(ET_JMP_U)
		link(ac, catch, EK_SEQ, ET_NONE, EK_JMP, ET_JMP_U, exit, LF_FORCE_SEP)
		link(ac, try, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		vn := grpBlock(ac, enter, origExit)
		vn.Outlets = append(vn.Outlets, try.xJmpOutEdges()...)
		vn.Outlets = append(vn.Outlets, catch.xJmpOutEdges()...)
		if fin != nil {
			vn.Outlets = append(vn.Outlets, fin.xJmpOutEdges()...)
		}
		ac.pushStmt(vn)

	case parser.N_STMT_THROW:
		n := node.(*parser.ThrowStmt)
		prev := ac.popStmt()
		expr := ac.popExpr()
		exit := ac.newExit(node, "")

		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, LF_NONE)

		prev = grpBlock(ac, prev, expr)
		if edge := prev.FindOutEdge(EK_JMP, ET_JMP_E, false); edge == nil && !n.Arg().Type().IsLit() {
			prev.newJmpOut(ET_JMP_E)
		}

		exit.newJmpOut(ET_JMP_U)
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_JMP, ET_JMP_T|ET_JMP_F, EK_NONE, ET_NONE, exit, LF_NONE)
		exit.mrkSeqOutAsCut()

		vn := grpBlock(ac, prev, exit)
		vn.Outlets = append(vn.Outlets, expr.xJmpOutEdges()...)

		if n.Target() != nil {
			if n.Target().Type() == parser.N_PROG {
				ac.graph.hasHangingThrow = true
			} else {
				ac.root.addHangingThrow(n.Target(), vn)
			}
		}

		ac.pushStmt(vn)

	case parser.N_STMT_SWITCH:
		n := node.(*parser.SwitchStmt)

		var head, tail, th, tt, defBody, lastCase *Block
		xOutEdges := []*Edge{}
		cnt := len(n.Cases())
		if cnt == 1 {
			head = ac.popStmt()
			tail = head
		} else if cnt > 1 {
			tail = ac.popStmt()
			lastCase = tail
			tt = tail
			for i := cnt - 1; i > 0; i-- {
				th = ac.popStmt()

				thTest, thBody, def := SwitchCase(th)
				if def {
					link(ac, thTest, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, tt, LF_FORCE_SEP|LF_OVERWRITE)
					defBody = thBody
				} else {
					link(ac, thTest, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, tt, LF_NONE)
				}

				_, ttBody, def := SwitchCase(tt)
				if def && defBody != nil {
					ttBody = defBody
				}
				link(ac, th, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, ttBody, LF_NONE)
				xOutEdges = append(xOutEdges, th.xJmpOutEdges()...)

				tt = th
			}

			head = th
		}

		if defBody != nil && cnt > 1 {
			thTest, _, _ := SwitchCase(lastCase)
			defBody.newJmpIn(ET_JMP_F)
			link(ac, thTest, EK_JMP, ET_JMP_F, EK_JMP, ET_JMP_F, defBody, LF_NONE)
		}

		enter := ac.popStmt()

		test := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)

		exit := ac.newExit(node, "")
		_, hasBreak := ac.graph.hangingBrk[ctx.ScopeId()]

		flag := LF_NONE
		if hasBreak {
			flag |= LF_FORCE_SEP
		}

		if head != nil {
			link(ac, test, EK_JMP, ET_NONE, EK_JMP, ET_NONE, head, LF_NONE)
			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			link(ac, tail, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, flag)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, flag)
		} else {
			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, flag)
		}

		vn := grpBlock(ac, enter, exit)
		vn.Outlets = append(vn.Outlets, xOutEdges...)
		ac.pushStmt(vn)

	case parser.N_SWITCH_CASE:
		n := node.(*parser.SwitchCase)
		enter := ac.newEnter(node, "")

		var test *Block
		if n.Test() != nil {
			test = ac.popExpr()
			test.newJmpOut(ET_JMP_F)
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, LF_NONE)
		} else {
			// below code makes path grow from `test`, so when `test` is
			// nil which happens when node is the `default` clause, set it
			// to `enter` to unify below logic
			test = enter
		}

		var body *Block
		if len(n.Cons()) > 0 {
			body = ac.popStmt()
		}

		exit := ac.newExit(node, "")
		if body != nil {
			body.newJmpIn(ET_JMP_T)
			link(ac, test, EK_SEQ, ET_NONE, EK_JMP, ET_JMP_T, body, LF_NONE)
			link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			// discard the empty entry block which is ued for holding
			// the body of case clause
			ac.popStmt()

			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		vn := grpBlock(ac, enter, exit)
		if body != nil {
			vn.Outlets = append(vn.Outlets, body.xOutEdges()...)
		}
		ac.pushStmt(vn)

	case parser.N_STMT_LABEL:
		body := ac.popStmt()
		enter := ac.popStmt()
		lb := ac.popExpr()

		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, lb, LF_NONE)
		link(ac, lb, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, body, LF_NONE)
		exit := ac.newExit(node, "")
		link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_CONT:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		exit.newJmpOut(ET_LOOP)

		n := node.(*parser.ContStmt)
		var target *Block
		if n.Label() != nil {
			name := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, name, LF_NONE)
			link(ac, name, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

			target = ac.graph.newBasicBlk()
			ac.graph.addHangingCont(n.Target().(*parser.LabelStmt).Body(), target)
		} else {
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

			target = ac.graph.newBasicBlk()
			ac.graph.addHangingCont(n.Target(), target)
		}

		link(ac, exit, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, target, LF_NONE)
		exit.mrkSeqOutAsCut()

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	case parser.N_STMT_BRK:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")

		n := node.(*parser.BrkStmt)
		if n.Label() != nil {
			name := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, name, LF_NONE)
			link(ac, name, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
			exit.mrkSeqOutAsCut()
			exit.newJmpOut(ET_JMP_U)
			id := ac.graph.labelLoop[n.Target()]
			ac.graph.addHangingBrk(id, exit)
		} else {
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
			exit.mrkSeqOutAsCut()
			exit.newJmpOut(ET_JMP_U)
			filter := &[]parser.ScopeKind{parser.SPK_LOOP_DIRECT, parser.SPK_SWITCH}
			id := ctx.Scope().UpperScope(filter).Id
			ac.graph.addHangingBrk(id, exit)
		}

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	case parser.N_STMT_FN, parser.N_EXPR_FN, parser.N_EXPR_ARROW:
		var params []parser.Node
		var id parser.Node
		expRet := false

		if astTyp == parser.N_EXPR_ARROW {
			n := node.(*parser.ArrowFn)
			params = n.Params()
			expRet = n.ExpRet()
		} else {
			n := node.(*parser.FnDec)
			id = n.Id()
			params = n.Params()
			expRet = n.ExpRet()
		}

		paramsLen := len(params)
		head, tail := ac.popExprsAndLink(paramsLen)

		var fnName *Block
		if id != nil {
			fnName = ac.popExpr()
		}
		body := ac.graph.Head

		// the actual entry node
		ac.graph.Head = ac.newEnter(node, "")

		// - do the connection to link: head of fnDec -> fn id -> params -> fn body
		if head != nil {
			if fnName != nil {
				link(ac, ac.graph.Head, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fnName, LF_NONE)
				link(ac, fnName, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)
				link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			} else {
				link(ac, ac.graph.Head, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)
				link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			}
		} else {
			if fnName != nil {
				link(ac, ac.graph.Head, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, fnName, LF_NONE)
				link(ac, fnName, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			} else {
				link(ac, ac.graph.Head, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)
			}
		}

		// connect the tail of body to the exit node
		enter := ac.popStmt()
		exit := ac.newExit(node, "")
		if expRet {
			exit.newJmpIn(ET_JMP_U)
		}
		link(ac, enter, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.leaveGraph()

	case parser.N_STMT_RET:
		prev := ac.popStmt()
		origPrev := prev

		exit := ac.newExit(node, "")
		exit.newJmpOut(ET_JMP_U)

		n := node.(*parser.RetStmt)
		if n.Arg() != nil {
			arg := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, arg, LF_NONE)
			prev = arg
		}
		link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, exit, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		exit.mrkSeqOutAsCut()

		vn := grpBlock(ac, origPrev, exit)
		ac.pushStmt(vn)

	case parser.N_IMPORT_SPEC:
		n := node.(*parser.ImportSpec)

		exit := ac.newExit(node, "")

		var id *Block
		if n.Default() {
			ac.popExpr() // discard id
		} else if !n.NameSpace() {
			id = ac.popExpr()
		}

		local := ac.popExpr()
		enter := ac.popExpr()

		prev := enter
		if id != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, id, LF_NONE)
			prev = grpBlock(ac, enter, id)
		}

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, local, LF_NONE)
		prev = grpBlock(ac, prev, local)

		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)

		enter := ac.popStmt()
		exit := ac.newExit(node, "")

		src := ac.popExpr()
		head, tail := ac.popExprsAndLink(len(n.Specs()))

		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, src, LF_NONE)
			link(ac, src, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, src, LF_NONE)
			link(ac, src, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_EXPORT_SPEC:
		n := node.(*parser.ExportSpec)

		exit := ac.newExit(node, "")

		var id *Block
		if !n.NameSpace() {
			id = ac.popExpr()
		}

		local := ac.popExpr()
		enter := ac.popExpr()

		prev := enter
		if id != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, id, LF_NONE)
			prev = grpBlock(ac, enter, id)
		}

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, local, LF_NONE)
		prev = grpBlock(ac, prev, local)

		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_STMT_EXPORT:
		n := node.(*parser.ExportDec)

		enter := ac.popStmt()
		exit := ac.newExit(node, "")

		var src *Block
		if n.Src() != nil {
			src = ac.popExpr()
		}

		var head, tail *Block
		if len(n.Specs()) > 0 {
			head, tail = ac.popExprsAndLink(len(n.Specs()))
		} else if n.Dec() != nil {
			head = ac.popExpr()
			if head == nil {
				head = enter
				enter = ac.popStmt()
			}
			tail = head
		}

		if head != nil {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, head, LF_NONE)

			if src != nil {
				link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, src, LF_NONE)
				link(ac, src, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
			} else {
				link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
			}
		} else {
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, src, LF_NONE)
			link(ac, src, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)
		}

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_WITH:
		obj := ac.popExpr()
		body := ac.popStmt()
		enter := ac.popStmt()
		exit := ac.newExit(node, "")

		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, obj, LF_NONE)

		prev := grpBlock(ac, enter, obj)
		link(ac, prev, EK_JMP, ET_NONE, EK_JMP, ET_NONE, body, LF_NONE)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, LF_NONE)

		prev = grpBlock(ac, prev, body)
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, LF_NONE)

		ac.pushStmt(grpBlock(ac, prev, exit))

	default:
		if !isAtom(astTyp) {
			prev := ac.popStmt()
			exit := ac.newExit(node, "")
			link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, exit, LF_NONE)
			ac.pushStmt(grpBlock(ac, prev, exit))
		}
	}

	// resolve the hanging breaks inside the loop, link them to the loop exit
	if isLoop(astTyp) || astTyp == parser.N_STMT_SWITCH {
		id := ctx.ScopeId()
		loopBlk := ac.lastStmt()
		loopExit := loopBlk.OutSeqEdge().Src
		brkList := ac.graph.hangingBrk[id]
		for _, brk := range brkList {
			link(ac, brk, EK_JMP, ET_NONE, EK_JMP, ET_NONE, loopExit, LF_OVERWRITE)
		}
	}
}

func resolveCont(a *AnalysisCtx, loopNode parser.Node, loopTailBlk *Block) {
	for _, cont := range a.graph.hangingCont[loopNode] {
		// cont is a placeholder, so we need to find the source blk, then it to loopTailBlk
		edge, _ := cont.FindInEdge(EK_JMP, ET_LOOP, false)
		jmp := edge.Src
		link(a, jmp, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, loopTailBlk, LF_OVERWRITE)
	}
}

func (a *Analysis) init() {
	a.WalkCtx.Extra = newAnalysisCtx()

	walk.AddBeforeListener(&a.WalkCtx.Listeners, handleBefore)
	walk.AddAfterListener(&a.WalkCtx.Listeners, handleAfter)
}

func analysisCtx(ctx *walk.VisitorCtx) *AnalysisCtx {
	return ctx.WalkCtx.Extra.(*AnalysisCtx)
}

func (a *Analysis) Analyze() {
	walk.VisitNode(a.WalkCtx.Root, "", a.WalkCtx.VisitorCtx())
}

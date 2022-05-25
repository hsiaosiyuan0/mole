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
	b := newBasicBlk()
	b.Nodes = append(b.Nodes, newInfoNode(astNode, true, info))
	return b
}

func (a *AnalysisCtx) newExit(astNode parser.Node, info string) *Block {
	b := newBasicBlk()
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
			link(a, th, EK_NONE, ET_NONE, EK_NONE, ET_NONE, tt, false)
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
	b := newBasicBlk()
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

	// record the first basic block of the label stmt in `labelBlk`
	if isAtom(astTyp) && pAstTyp != parser.N_STMT_LABEL {
		if cnt := len(ac.graph.hangingLabels); cnt > 0 {
			blk.newLoopIn()
			last, rest := ac.graph.hangingLabels[cnt-1], ac.graph.hangingLabels[:cnt-1]
			ac.graph.hangingLabels = rest
			ac.graph.labelBlk[last] = blk
			ac.graph.labelLoop[last] = ctx.ScopeId()
		}

		if cnt := len(ac.graph.hangingLoops); cnt > 0 {
			last, rest := ac.graph.hangingLoops[cnt-1], ac.graph.hangingLoops[:cnt-1]
			ac.graph.hangingLoops = rest
			ac.graph.loopBlk[last] = blk
		}
	}

	// stmt
	if astTyp.IsStmt() || astTyp == parser.N_PROG {
		if astTyp == parser.N_STMT_LABEL && node.(*parser.LabelStmt).Used() {
			id := IdOfAstNode(node)
			ac.graph.hangingLabels = append(ac.graph.hangingLabels, id)
		} else if isLoop(astTyp) {
			id := ctx.ScopeId()
			ac.graph.hangingLoops = append(ac.graph.hangingLoops, id)
		} else if astTyp == parser.N_STMT_FN || astTyp == parser.N_EXPR_FN {
			enterFnGraph(node, ctx, true)
		}

		if pAstTyp == parser.N_STMT_IF || isLoop(pAstTyp) || pAstTyp == parser.N_STMT_LABEL {
			// just pushing `enter` into stack without linking to the `prevStmt` to imitate a new branch
			// the forked branch will be linked to its source branch point in the post-process listener
			// of its astNode
			ac.pushStmt(blk)
		} else if astTyp != parser.N_STMT_FN && astTyp != parser.N_EXPR_FN {
			prev := ac.popStmt()
			link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, blk, false)
			ac.pushStmt(grpBlock(ac, prev, blk))
		}
	} else if !isAtom(astTyp) {
		if isFn(astTyp) {
			enterFnGraph(node, ctx, false)
		} else if astTyp == parser.N_SWITCH_CASE {
			// here we push en empty block into the stack to start a new fork for the case clause
			// so its body will grow from this empty block. the new empty block does not link to
			// previous block yet at here, the connection will happen in the post-case handler
			b := newBasicBlk()
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
	b := newBasicBlk()
	b.Nodes = append(b.Nodes, newInfoNode(node, true, node.Type().String()))

	if stmt {
		// the info node should be connected to previous block if its origin is fnStmt
		// that's because the link strategy of stmt
		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, b, false)
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
	ac.graph.Head = newBasicBlk()
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
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, lhs, false)
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

		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, rhs, false)

		if op == parser.T_OR || op == parser.T_NULLISH {
			link(ac, lhs, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, rhs, false)
		}

		vn := newGroupBlk()
		vn.Inlets = lhs.Inlets
		if logic {
			vn.Outlets = append(rhs.Outlets, lhs.xOutEdges()...)
		} else {
			vn.Outlets = rhs.Outlets
		}

		exit := ac.newExit(node, "")
		link(ac, vn, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
		ac.pushExpr(grpBlock(ac, vn, exit))

	case parser.N_EXPR_UPDATE, parser.N_EXPR_PAREN:
		expr := ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, expr, false)
		exit := ac.newExit(node, "")
		link(ac, expr, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_LIT_ARR:
		n := node.(*parser.ArrLit)

		head, tail := ac.popExprsAndLink(len(n.Elems()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, head, false)
			link(ac, tail, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		} else {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		}

		vn := grpBlock(ac, enter, exit)
		ac.pushExpr(vn)

	case parser.N_PROP:
		n := node.(*parser.Prop)
		var key, val *Block
		if n.Val() != nil {
			val = ac.popExpr()
		}
		key = ac.popExpr()

		enter := ac.popExpr()
		link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, key, false)

		exit := ac.newExit(node, "")
		if val != nil {
			link(ac, key, EK_NONE, ET_NONE, EK_NONE, ET_NONE, val, false)
			link(ac, val, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		} else {
			link(ac, key, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		}

		vn := grpBlock(ac, enter, exit)
		ac.pushExpr(vn)

	case parser.N_LIT_OBJ:
		n := node.(*parser.ObjLit)

		head, tail := ac.popExprsAndLink(len(n.Props()))

		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, head, false)
			link(ac, tail, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		} else {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		}

		vn := grpBlock(ac, enter, exit)
		ac.pushExpr(vn)

	case parser.N_EXPR_ASSIGN:
		rhs := ac.popExpr()
		lhs := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, lhs, false)
		link(ac, lhs, EK_NONE, ET_NONE, EK_NONE, ET_NONE, rhs, false)
		link(ac, rhs, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		vn := grpBlock(ac, enter, exit)
		ac.pushExpr(vn)

	case parser.N_EXPR_CALL:
		n := node.(*parser.CallExpr)

		head, tail := ac.popExprsAndLink(len(n.Args()))

		fn := ac.popExpr()
		enter := ac.popExpr()
		exit := ac.newExit(node, "")
		if head != nil {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, fn, false)
			link(ac, fn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, head, false)
			link(ac, tail, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		} else {
			link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, fn, false)
			link(ac, fn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		}

		vn := grpBlock(ac, enter, exit)
		ac.pushExpr(vn)

	case parser.N_STMT_EXPR:
		expr := ac.popExpr()
		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, expr, false)

		exit := ac.newExit(node, "")
		link(ac, expr, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, prev, exit))

	case parser.N_VAR_DEC:
		n := node.(*parser.VarDec)
		var id, init *Block
		if n.Init() != nil {
			init = ac.popExpr()
		}
		id = ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, id, false)

		exit := ac.newExit(node, "")
		if init != nil {
			link(ac, id, EK_NONE, ET_NONE, EK_NONE, ET_NONE, init, false)
			link(ac, init, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		} else {
			link(ac, id, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		}
		vn := grpBlock(ac, enter, exit)

		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, vn, false)
		ac.pushStmt(grpBlock(ac, prev, vn))

	case parser.N_STMT_VAR_DEC:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, prev, exit))

		if pAstTyp == parser.N_STMT_FOR && key == "Init" ||
			pAstTyp == parser.N_STMT_FOR_IN_OF && key == "Left" {
			ac.pushExpr(ac.popStmt())
		}

	case parser.N_PROG, parser.N_STMT_BLOCK:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, exit, false)
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

		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, test, false)
		test = grpBlock(ac, enter, test)

		test.newJmpOut(ET_JMP_F)
		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, cons, false)
		link(ac, test, EK_JMP, ET_JMP_T, EK_SEQ, ET_NONE, cons, false)
		if alt != nil {
			link(ac, test, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, alt, false)
		}

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets

		vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		vn.Outlets = append(vn.Outlets, cons.xOutEdges()...)
		if alt != nil {
			vn.Outlets = append(vn.Outlets, alt.xOutEdges()...)
		}

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		vn = grpBlock(ac, vn, exit)
		ac.pushStmt(vn)

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
			link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, init, false)
			link(ac, init, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, false)
		}

		if test != nil {
			if init == nil {
				link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, test, false)
			}
			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
		} else {
			body.newLoopIn()
			if init == nil {
				link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
			} else {
				link(ac, init, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
			}
		}

		if update != nil {
			link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, update, false)
			if test != nil {
				link(ac, update, EK_NONE, ET_NONE, EK_NONE, ET_NONE, test, false)
			} else {
				link(ac, update, EK_NONE, ET_NONE, EK_NONE, ET_NONE, body, false)
			}
			update.mrkSeqOutAsLoop()
		} else {
			if test != nil {
				link(ac, body, EK_NONE, ET_NONE, EK_NONE, ET_NONE, test, false)
			} else {
				link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
			}
			body.mrkSeqOutAsLoop()
		}

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_WHILE:
		body := ac.popStmt()
		test := ac.popExpr()
		enter := ac.popStmt()

		test.newJmpOut(ET_JMP_F)
		test.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, false)

		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
		link(ac, test, EK_JMP, ET_JMP_T, EK_SEQ, ET_NONE, body, false)
		link(ac, body, EK_NONE, ET_NONE, EK_NONE, ET_NONE, test, false)
		body.mrkSeqOutAsLoop()

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_DO_WHILE:
		test := ac.popExpr()
		body := ac.popStmt()
		enter := ac.popStmt()

		test.newJmpOut(ET_JMP_F)
		body.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, body, false)
		link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, false)

		link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)
		test.mrkSeqOutAsLoop()
		link(ac, test, EK_JMP, ET_JMP_T, EK_JMP, ET_NONE, body, false)
		test.mrkJmpOutAsLoop(ET_JMP_T)

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_FOR_IN_OF:
		body := ac.popStmt()
		enter := ac.popStmt()
		rhs := ac.popExpr()
		lhs := ac.popExpr()

		lhs.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, lhs, false)
		link(ac, lhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, rhs, false)

		rhs.newJmpOut(ET_JMP_F)
		link(ac, rhs, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, body, false)

		body.mrkSeqOutAsLoop()
		link(ac, body, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, lhs, false)

		exit := ac.newExit(node, "")
		link(ac, rhs, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, exit, false)

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_SWITCH:
		n := node.(*parser.SwitchStmt)

		var head, tail, th, tt, defBody, lastCase *Block
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
					link(ac, thTest, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, tt, true)
					defBody = thBody
				} else {
					link(ac, thTest, EK_JMP, ET_JMP_F, EK_SEQ, ET_NONE, tt, false)
				}

				_, ttBody, def := SwitchCase(tt)
				if def && defBody != nil {
					ttBody = defBody
				}
				link(ac, th, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, ttBody, false)

				tt = th
			}

			head = th
		}

		if defBody != nil && cnt > 1 {
			thTest, _, _ := SwitchCase(lastCase)
			defBody.newJmpIn(ET_JMP_F)
			link(ac, thTest, EK_JMP, ET_JMP_F, EK_JMP, ET_JMP_F, defBody, false)
		}

		enter := ac.popStmt()

		test := ac.popExpr()
		link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, test, false)

		exit := ac.newExit(node, "")
		_, hasBreak := ac.graph.hangingBrk[ctx.ScopeId()]

		if head != nil {
			link(ac, test, EK_NONE, ET_NONE, EK_SEQ, ET_NONE, head, false)
			link(ac, tail, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, hasBreak)
		} else {
			link(ac, test, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, hasBreak)
		}

		vn := grpBlock(ac, enter, exit)
		ac.pushStmt(vn)

	case parser.N_SWITCH_CASE:
		n := node.(*parser.SwitchCase)
		enter := ac.newEnter(node, "")

		var test *Block
		if n.Test() != nil {
			test = ac.popExpr()
			test.newJmpOut(ET_JMP_F)
			link(ac, enter, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, test, false)
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
			link(ac, test, EK_SEQ, ET_NONE, EK_JMP, ET_JMP_T, body, false)
			link(ac, body, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
		} else {
			// discard the empty entry block which is ued for holding
			// the body of case clause
			ac.popStmt()

			link(ac, test, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
		}

		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_LABEL:
		body := ac.popStmt()
		enter := ac.popStmt()

		link(ac, enter, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, body, false)
		exit := ac.newExit(node, "")
		link(ac, body, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)
		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_CONT:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		exit.newJmpOut(ET_LOOP)

		n := node.(*parser.ContStmt)
		var target *Block
		if n.Label() != nil {
			name := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, name, false)
			link(ac, name, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
			target = ac.graph.labelBlk[IdOfAstNode(n.Target())]
		} else {
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
			id := ctx.Scope().UpperLoop().Id
			target = ac.graph.loopBlk[id]
		}

		link(ac, exit, EK_JMP, ET_LOOP, EK_JMP, ET_LOOP, target, false)
		exit.mrkSeqOutAsCut()

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	case parser.N_STMT_BRK:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")

		n := node.(*parser.BrkStmt)
		if n.Label() != nil {
			name := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, name, false)
			link(ac, name, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
			exit.mrkSeqOutAsCut()
			exit.newJmpOut(ET_JMP_U)
			id := ac.graph.labelLoop[IdOfAstNode(n.Target())]
			ac.graph.addHangingBrk(id, exit)
		} else {
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, exit, false)
			exit.mrkSeqOutAsCut()
			exit.newJmpOut(ET_JMP_U)
			filter := &[]parser.ScopeKind{parser.SPK_LOOP_DIRECT, parser.SPK_SWITCH}
			id := ctx.Scope().UpperScope(filter).Id
			ac.graph.addHangingBrk(id, exit)
		}

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	case parser.N_STMT_FN, parser.N_EXPR_FN:
		n := node.(*parser.FnDec)

		var head, tail, th, tt *Block
		paramsLen := len(n.Params())
		if paramsLen == 1 {
			head = ac.popExpr()
			tail = head
		} else if paramsLen > 1 {
			tail = ac.popExpr()
			tt = tail
			for i := paramsLen - 2; i >= 0; i-- {
				th = ac.popExpr()
				link(ac, th, EK_NONE, ET_NONE, EK_NONE, ET_NONE, tt, false)
				tt = th
			}
			head = th
		}

		var fnName *Block
		if n.Id() != nil {
			fnName = ac.popExpr()
		}
		body := ac.graph.Head

		// the actual entry node
		ac.graph.Head = ac.newEnter(node, "")

		// - do the connection to link: head of fnDec -> fn id -> params -> fn body
		if head != nil {
			if fnName != nil {
				link(ac, ac.graph.Head, EK_NONE, ET_NONE, EK_NONE, ET_NONE, fnName, false)
				link(ac, fnName, EK_NONE, ET_NONE, EK_NONE, ET_NONE, head, false)
				link(ac, tail, EK_NONE, ET_NONE, EK_NONE, ET_NONE, body, false)
			} else {
				link(ac, ac.graph.Head, EK_NONE, ET_NONE, EK_NONE, ET_NONE, head, false)
				link(ac, tail, EK_NONE, ET_NONE, EK_NONE, ET_NONE, body, false)
			}
		} else {
			if fnName != nil {
				link(ac, ac.graph.Head, EK_NONE, ET_NONE, EK_NONE, ET_NONE, fnName, false)
				link(ac, fnName, EK_NONE, ET_NONE, EK_NONE, ET_NONE, body, false)
			} else {
				link(ac, ac.graph.Head, EK_NONE, ET_NONE, EK_NONE, ET_NONE, body, false)
			}
		}

		// connect the tail of body to the exit node
		enter := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, enter, EK_NONE, ET_NONE, EK_NONE, ET_NONE, exit, false)

		ac.leaveGraph()

	case parser.N_STMT_RET:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		exit.newJmpOut(ET_JMP_U)

		n := node.(*parser.RetStmt)
		if n.Arg() != nil {
			arg := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, EK_SEQ, ET_NONE, arg, false)
			prev = arg
		}
		link(ac, prev, EK_NONE, ET_NONE, EK_SEQ, ET_NONE, exit, false)
		exit.mrkSeqOutAsCut()

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	default:
		if !isAtom(astTyp) {
			prev := ac.popStmt()
			exit := ac.newExit(node, "")
			link(ac, prev, EK_SEQ, ET_NONE, EK_NONE, ET_NONE, exit, false)
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
			if !loopExit.hasXOut(EK_JMP, ET_JMP_U) {
				loopExit.newJmpOut(ET_JMP_U)
			}
			link(ac, brk, EK_JMP, ET_NONE, EK_JMP, ET_NONE, loopExit, false)
		}
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

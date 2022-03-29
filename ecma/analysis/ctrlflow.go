package analysis

import (
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

const (
	N_CFG_DEBUG parser.NodeType = parser.N_NODE_DEF_END + 1 + iota
)

type AnalysisCtx struct {
	graph *Graph

	stmtStack []*Block
	exprStack []*Block

	// map astNode to its basic block
	astNodeToBlock map[uint64]Block
}

func newAnalysisCtx() *AnalysisCtx {
	a := &AnalysisCtx{
		graph:          newGraph(),
		stmtStack:      make([]*Block, 0),
		exprStack:      make([]*Block, 0),
		astNodeToBlock: map[uint64]Block{},
	}
	a.stmtStack = append(a.stmtStack, a.graph.Head)
	return a
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

func (a *Analysis) Graph() *Graph {
	return analysisCtx(a.WalkCtx.VisitorCtx()).graph
}

// below stmts and exprs have the condJmp(conditional-jump) semantic:
// - [x] logicAnd
// - [x] logicOr
// - [x] if
// - [x] for
// - [x] while
// - [x] doWhile
// - [] for-in-of

// - [ ] Loop
// - [x] Test

// below stmts have the unCondJmp(unconditional-jump) semantic
// - [ ] contine
// - [ ] break
// - [ ] return
// - [ ] callExpr

func isLoop(t parser.NodeType) bool {
	return t == parser.N_STMT_FOR || t == parser.N_STMT_WHILE || t == parser.N_STMT_DO_WHILE
}

func isAtom(t parser.NodeType) bool {
	_, ok := walk.AtomNodeTypes[t]
	return ok
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

	if astTyp.IsExpr() && pAstTyp != parser.N_STMT_LABEL {
		if cnt := len(ac.graph.hangingLabels); cnt > 0 {
			blk.newLoopIn()
			last, rest := ac.graph.hangingLabels[cnt-1], ac.graph.hangingLabels[:cnt-1]
			ac.graph.hangingLabels = rest
			ac.graph.labelAstMap[last] = blk
		}
	}

	if astTyp.IsStmt() || astTyp == parser.N_PROG {
		if astTyp == parser.N_STMT_LABEL && node.(*parser.LabelStmt).Used() {
			id := IdOfAstNode(node)
			ac.graph.labelAstMap[id] = newBasicBlk()
			ac.graph.hangingLabels = append(ac.graph.hangingLabels, id)
		}

		if pAstTyp == parser.N_STMT_IF || isLoop(pAstTyp) || pAstTyp == parser.N_STMT_LABEL {
			// just pushing `enter` into stack without linking to the `prevStmt` to imitate a new branch
			// the forked branch will be linked to its source branch point in the post-process listener
			// of its astNode
			ac.pushStmt(blk)
		} else {
			prev := ac.popStmt()
			link(ac, prev, EK_NONE, ET_NONE, blk)
			ac.pushStmt(grpBlock(ac, prev, blk))
		}
	} else if astTyp.IsExpr() && !isAtom(astTyp) || astTyp == parser.N_VAR_DEC {
		ac.pushExpr(blk)
	}
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
		link(ac, enter, EK_NONE, ET_NONE, lhs)
		lhs = grpBlock(ac, enter, lhs)

		op := n.Op()
		logic := true
		if op == parser.T_AND {
			lhs.newJmp(ET_JMP_F)
		} else if op == parser.T_OR {
			lhs.newJmp(ET_JMP_T)
		} else {
			logic = false
		}

		link(ac, lhs, EK_SEQ, ET_NONE, rhs)

		if op == parser.T_OR {
			link(ac, lhs, EK_JMP, ET_JMP_F, rhs)
		}

		vn := newGroupBlk()
		vn.Inlets = lhs.Inlets
		if logic {
			vn.Outlets = append(rhs.Outlets, lhs.xOutEdges()...)
		} else {
			vn.Outlets = rhs.Outlets
		}

		exit := ac.newExit(node, "")
		link(ac, vn, EK_SEQ, ET_NONE, exit)
		ac.pushExpr(grpBlock(ac, vn, exit))

	case parser.N_EXPR_UPDATE, parser.N_EXPR_PAREN:
		expr := ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_NONE, ET_NONE, expr)
		exit := ac.newExit(node, "")
		link(ac, expr, EK_NONE, ET_NONE, exit)
		ac.pushExpr(grpBlock(ac, enter, exit))

	case parser.N_STMT_EXPR:
		expr := ac.popExpr()
		prev := ac.popStmt()
		link(ac, prev, EK_SEQ, ET_NONE, expr)

		exit := ac.newExit(node, "")
		link(ac, expr, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, prev, exit))

	case parser.N_VAR_DEC:
		n := node.(*parser.VarDec)
		var id, init *Block
		if n.Init() != nil {
			init = ac.popExpr()
		}
		id = ac.popExpr()
		enter := ac.popExpr()
		link(ac, enter, EK_NONE, ET_NONE, id)

		exit := ac.newExit(node, "")
		if init != nil {
			link(ac, id, EK_NONE, ET_NONE, init)
			link(ac, init, EK_NONE, ET_NONE, exit)
		} else {
			link(ac, id, EK_NONE, ET_NONE, exit)
		}
		vn := grpBlock(ac, enter, exit)

		prev := ac.popStmt()
		link(ac, prev, EK_NONE, ET_NONE, vn)
		ac.pushStmt(grpBlock(ac, prev, vn))

	case parser.N_STMT_VAR_DEC:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, prev, exit))

		if pAstTyp == parser.N_STMT_FOR && key == "Init" {
			ac.pushExpr(ac.popStmt())
		}

	case parser.N_PROG, parser.N_STMT_BLOCK:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		link(ac, prev, EK_NONE, ET_NONE, exit)
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

		link(ac, enter, EK_NONE, ET_NONE, test)
		test = grpBlock(ac, enter, test)

		test.newJmp(ET_JMP_F)
		link(ac, test, EK_SEQ, ET_NONE, cons)
		link(ac, test, EK_JMP, ET_JMP_T, cons)
		if alt != nil {
			link(ac, test, EK_JMP, ET_JMP_F, alt)
		}

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets

		vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		vn.Outlets = append(vn.Outlets, cons.xOutEdges()...)
		if alt != nil {
			vn.Outlets = append(vn.Outlets, alt.xOutEdges()...)
		}

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, exit)
		vn = grpBlock(ac, vn, exit)
		ac.pushStmt(vn)

	case parser.N_STMT_FOR:
		n := node.(*parser.ForStmt)
		var enter, init, test, update, body *Block
		if n.Body() != nil {
			body = ac.popStmt()
		}
		if n.Update() != nil {
			update = ac.popExpr()
		}
		if n.Test() != nil {
			test = ac.popExpr()
			test.newJmp(ET_JMP_F)
			test.newLoopIn()
		}
		if n.Init() != nil {
			init = ac.popExpr()
		}
		enter = ac.popStmt()

		if init != nil {
			link(ac, enter, EK_NONE, ET_NONE, init)
			link(ac, init, EK_SEQ, ET_NONE, test)
		}

		if test != nil {
			if init == nil {
				link(ac, enter, EK_NONE, ET_NONE, test)
			}
			link(ac, test, EK_SEQ, ET_NONE, body)
		} else {
			body.newLoopIn()
			if init == nil {
				link(ac, enter, EK_SEQ, ET_NONE, body)
			} else {
				link(ac, init, EK_SEQ, ET_NONE, body)
			}
		}

		if update != nil {
			link(ac, body, EK_SEQ, ET_NONE, update)
			if test != nil {
				link(ac, update, EK_NONE, ET_NONE, test)
			} else {
				link(ac, update, EK_NONE, ET_NONE, body)
			}
			update.mrkSeqOutAsLoop()
		} else {
			if test != nil {
				link(ac, body, EK_NONE, ET_NONE, test)
			} else {
				link(ac, body, EK_SEQ, ET_NONE, body)
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
		link(ac, vn, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_WHILE:
		body := ac.popStmt()
		test := ac.popExpr()
		enter := ac.popStmt()

		test.newJmp(ET_JMP_F)
		test.newLoopIn()
		link(ac, enter, EK_SEQ, ET_NONE, test)

		link(ac, test, EK_SEQ, ET_NONE, body)
		link(ac, test, EK_JMP, ET_JMP_T, body)
		link(ac, body, EK_NONE, ET_NONE, test)
		body.mrkSeqOutAsLoop()

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_DO_WHILE:
		test := ac.popExpr()
		body := ac.popStmt()
		enter := ac.popStmt()

		test.newJmp(ET_JMP_F)
		body.newLoopIn()
		link(ac, enter, EK_NONE, ET_NONE, body)
		link(ac, body, EK_SEQ, ET_NONE, test)

		link(ac, test, EK_SEQ, ET_NONE, body)
		test.mrkSeqOutAsLoop()
		link(ac, test, EK_JMP, ET_JMP_T, body)
		test.mrkJmpOutAsLoop(ET_JMP_T)

		vn := newGroupBlk()
		vn.Inlets = enter.Inlets
		if test != nil {
			vn.Outlets = append(vn.Outlets, test.xJmpOutEdges()...)
		}
		vn.Outlets = append(vn.Outlets, body.xOutEdges()...)

		exit := ac.newExit(node, "")
		link(ac, vn, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, vn, exit))

	case parser.N_STMT_LABEL:
		body := ac.popStmt()
		enter := ac.popStmt()

		link(ac, enter, EK_NONE, ET_NONE, body)
		exit := ac.newExit(node, "")
		link(ac, body, EK_NONE, ET_NONE, exit)
		ac.pushStmt(grpBlock(ac, enter, exit))

	case parser.N_STMT_CONT:
		prev := ac.popStmt()
		exit := ac.newExit(node, "")
		exit.newJmp(ET_LOOP)

		n := node.(*parser.ContStmt)
		var target *Block
		if n.Label() != nil {
			name := ac.popExpr()
			link(ac, prev, EK_SEQ, ET_NONE, name)
			link(ac, name, EK_SEQ, ET_NONE, exit)
			target = ac.graph.labelAstMap[IdOfAstNode(n.Target())]
		} // TODO: no label

		link(ac, exit, EK_JMP, ET_LOOP, target)
		exit.mrkSeqOutAsCutted()

		vn := grpBlock(ac, prev, exit)
		ac.pushStmt(vn)

	default:
		if !isAtom(astTyp) {
			prev := ac.popStmt()
			exit := ac.newExit(node, "")
			link(ac, prev, EK_NONE, ET_NONE, exit)
			ac.pushStmt(grpBlock(ac, prev, exit))
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

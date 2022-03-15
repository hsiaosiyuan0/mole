package analysis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

type EdgeKind uint8

const (
	EK_NONE EdgeKind = 0
	EK_SEQ  EdgeKind = 1 << iota
	EK_JMP_FALSE
	EK_JMP_TRUE
	EK_LOOP
	EK_UNREACHABLE
)

func (k EdgeKind) String() string {
	if k&EK_JMP_TRUE != 0 {
		return "T"
	}
	if k&EK_JMP_FALSE != 0 {
		return "F"
	}
	if k&EK_LOOP != 0 {
		return "L"
	}
	return ""
}

func (k EdgeKind) DotColor() string {
	if k&EK_UNREACHABLE != 0 {
		return "red"
	}
	if k&EK_JMP_TRUE != 0 || k&EK_JMP_FALSE != 0 || k&EK_LOOP != 0 {
		return "orange"
	}
	return "black"
}

type Edge struct {
	Kind EdgeKind
	Src  *Node
	Dst  *Node
}

func (e *Edge) Key() string {
	from := "s"
	if e.Src != nil {
		from = e.Src.Id()
	}
	to := "e"
	if e.Dst != nil {
		to = e.Dst.Id()
	}
	return from + "_" + to
}

func (e *Edge) isSeq() bool {
	return e.Kind&EK_SEQ != 0
}

func (e *Edge) Dot() string {
	s := "initial"
	if e.Src != nil {
		s = e.Src.Id()
	}
	d := "final"
	if e.Dst != nil {
		d = e.Dst.Id()
	}
	c := e.Kind.DotColor()
	if e.Kind&EK_LOOP != 0 {
		return fmt.Sprintf("%s:s->%s:ne [xlabel=\"%s\",color=\"%s\"];\n", s, d, e.Kind.String(), c)
	}
	return fmt.Sprintf("%s->%s [xlabel=\"%s\",color=\"%s\"];\n", s, d, e.Kind.String(), c)
}

func nodesToString(nodes []parser.Node) string {
	var b strings.Builder
	for _, node := range nodes {
		b.WriteString(nodeToString(node) + "\\n")
	}
	return b.String()
}

func nodeToString(node parser.Node) string {
	switch node.Type() {
	case parser.N_NAME, parser.N_LIT_NUM:
		return fmt.Sprintf("%s(%s)", node.Type().String(), node.Loc().Text())
	case N_CFG_DEBUG:
		return node.(*DebugNode).String()
	}
	return node.Type().String()
}

const (
	N_CFG_DEBUG parser.NodeType = parser.N_NODE_DEF_END + 1 + iota
)

type DebugNode struct {
	astNode parser.Node
	enter   bool
	info    string
}

func newDebugNode(node parser.Node, enter bool, info string) *DebugNode {
	return &DebugNode{node, enter, info}
}

func (n *DebugNode) String() string {
	if n.info != "" {
		return n.info
	}

	typ := n.astNode.Type()
	enter := "enter"
	if !n.enter {
		enter = "exit"
	}
	switch typ {
	case parser.N_EXPR_BIN:
		return fmt.Sprintf("%s(%s):%s", typ.String(), n.astNode.(*parser.BinExpr).OpText(), enter)
	case parser.N_EXPR_UPDATE:
		return fmt.Sprintf("%s(%s):%s", typ.String(), n.astNode.(*parser.UpdateExpr).OpText(), enter)
	case parser.N_EXPR_UNARY:
		return fmt.Sprintf("%s(%s):%s", typ.String(), n.astNode.(*parser.UnaryExpr).OpText(), enter)
	}
	return fmt.Sprintf("%s:%s", typ.String(), enter)
}

func (n *DebugNode) Type() parser.NodeType {
	return N_CFG_DEBUG
}

func (n *DebugNode) Loc() *parser.Loc {
	return n.astNode.Loc()
}

type Node struct {
	id       string
	AstNodes []parser.Node
	In       []*Edge
	Out      []*Edge
}

func IdOfAstNode(node parser.Node) uint64 {
	pos := node.Loc().Begin()
	return uint64(pos.Line())<<32 | uint64(pos.Column())
}

func IdStrOfAstNode(node parser.Node) string {
	return strconv.FormatUint(IdOfAstNode(node), 10)
}

func (n *Node) Id() string {
	if n.id != "" {
		return n.id
	}
	if len(n.AstNodes) == 0 {
		return ""
	}
	n.id = IdStrOfAstNode(n.AstNodes[0])
	return n.id
}

func (n *Node) Dot() string {
	return fmt.Sprintf("%s[label=\"%s\"];\n", n.Id(), nodesToString(n.AstNodes))
}

func (n *Node) OutEdge(kind EdgeKind) *Edge {
	for _, edge := range n.Out {
		if edge.Kind&kind != 0 {
			return edge
		}
	}
	return nil
}

func (n *Node) seqOutEdge() *Edge {
	for _, edge := range n.Out {
		if edge.isSeq() {
			return edge
		}
	}
	panic("unreachable")
}

func (n *Node) seqInEdge() *Edge {
	for _, edge := range n.In {
		if edge.isSeq() {
			return edge
		}
	}
	panic("unreachable")
}

// since the node maybe virtual, use this method to unwrap the actual node
func (n *Node) unwrapOut() *Node {
	// the `0` access is safe since each node has at least one out edge
	return n.Out[0].Src
}

// since the node maybe virtual, use this method to unwrap the actual node
func (n *Node) unwrapIn() *Node {
	// the `0` access is safe since each node has at least one out edge
	return n.In[0].Dst
}

// the prefix `x` means the edge is danging, its out node is not resolved
func (n *Node) xOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range n.Out {
		if edge.Dst == nil {
			ret = append(ret, edge)
		}
	}
	return ret
}

func (n *Node) isVirtual() bool {
	return n.In[0].Dst != n
}

func (n *Node) xJmpOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range n.Out {
		if !edge.isSeq() && edge.Dst == nil {
			ret = append(ret, edge)
		}
	}
	return ret
}

func (n *Node) newJmp(k EdgeKind) {
	seq := n.seqOutEdge()
	edge := &Edge{k, seq.Src, nil}
	seq.Src.Out = append(seq.Src.Out, edge)
	if seq.Src != n {
		n.Out = append(n.Out, edge)
	}
}

func (n *Node) newLoopIn() {
	seq := n.seqInEdge()
	edge := &Edge{EK_LOOP, nil, seq.Dst}
	seq.Dst.In = append(seq.Dst.In, edge)
	// for n is virtual
	if seq.Dst != n {
		n.In = append(n.In, edge)
	}
}

func (n *Node) newLoopOut() {
	seq := n.seqOutEdge()
	edge := &Edge{EK_LOOP, seq.Src, nil}
	seq.Src.Out = append(seq.Src.Out, edge)
	// for n is virtual
	if seq.Src != n {
		n.Out = append(n.Out, edge)
	}
}

func (n *Node) mrkSeqOutAsLoop() {
	n.seqOutEdge().Kind |= EK_LOOP
}

func (n *Node) mrkSeqOutAsUnreachable() {
	n.seqOutEdge().Kind |= EK_UNREACHABLE
}

type Graph struct {
	Id     string
	Head   *Node
	Parent *Graph
	Subs   []*Graph

	// collect the heads of the loops, key is calculated via `IdOfAstNode`
	// LabeledLoops map[string]*Node
}

func newGraph() *Graph {
	g := &Graph{
		Subs: make([]*Graph, 0),
		// LabeledLoops: map[string]*Node{},
	}
	g.Head = g.newNode()
	return g
}

func (g *Graph) newNode() *Node {
	n := &Node{
		AstNodes: make([]parser.Node, 0),
		In:       make([]*Edge, 0),
		Out:      make([]*Edge, 0),
	}

	in := &Edge{Dst: n, Kind: EK_SEQ}
	n.In = append(n.In, in)

	out := &Edge{Src: n, Kind: EK_SEQ}
	n.Out = append(n.Out, out)
	return n
}

func (g *Graph) NodesEdges() (map[string]*Node, map[string]*Edge, map[parser.Node]*Node) {
	edges := g.Head.In
	whites := []*Node{g.Head}

	uniqueNodes := map[string]*Node{}
	uniqueEdges := map[string]*Edge{}

	// map astNode to the basic block which encapsulates it
	astNodeMap := map[parser.Node]*Node{}
	for len(whites) > 0 {
		cnt := len(whites)
		last, rest := whites[cnt-1], whites[:cnt-1]
		whites = rest

		if _, ok := uniqueNodes[last.Id()]; ok {
			continue
		}

		uniqueNodes[last.Id()] = last
		for _, astNode := range last.AstNodes {
			astNodeMap[astNode] = last
		}

		if last.Out != nil {
			edges = append(edges, last.Out...)
			for _, edge := range last.Out {
				if edge.Dst != nil {
					whites = append(whites, edge.Dst)
				}
			}
		}
	}

	for _, edge := range edges {
		key := edge.Key()
		if _, ok := uniqueEdges[key]; ok {
			continue
		}
		uniqueEdges[key] = edge
	}

	return uniqueNodes, uniqueEdges, astNodeMap
}

func (g *Graph) Dot() string {
	nodes, edges, _ := g.NodesEdges()

	var b strings.Builder

	b.WriteString(`digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
`)

	for _, node := range nodes {
		b.WriteString(node.Dot())
	}

	for _, edge := range edges {
		b.WriteString(edge.Dot())
	}

	b.WriteString("}\n\n")
	return b.String()
}

type AnalysisCtx struct {
	graph *Graph

	stmtStack []*Node
	exprStack []*Node

	// map astNode to its basic block
	astNodeInBlock map[uint64]*Node
}

func newAnalysisCtx() *AnalysisCtx {
	a := &AnalysisCtx{
		graph:          newGraph(),
		stmtStack:      make([]*Node, 0),
		exprStack:      make([]*Node, 0),
		astNodeInBlock: map[uint64]*Node{},
	}
	a.stmtStack = append(a.stmtStack, a.graph.Head)
	return a
}

// func (a *AnalysisCtx) enterGraph() {
// 	s := newGraph()
// 	s.Parent = a.graph
// 	a.graph.Subs = append(a.graph.Subs, s)
// 	a.graph = s
// }

// func (a *AnalysisCtx) leaveGraph() {
// 	a.graph = a.graph.Parent
// }

func (a *AnalysisCtx) lastExpr() *Node {
	return a.exprStack[len(a.exprStack)-1]
}

func (a *AnalysisCtx) lastStmt() *Node {
	return a.stmtStack[len(a.stmtStack)-1]
}

func (a *AnalysisCtx) newEnter(astNode parser.Node, info string) *Node {
	exit := a.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(astNode, true, info))
	return exit
}

func (a *AnalysisCtx) newExit(astNode parser.Node, info string) *Node {
	exit := a.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(astNode, false, info))
	return exit
}

func (a *AnalysisCtx) pushStmt(n *Node) {
	a.stmtStack = append(a.stmtStack, n)
	// a.setAstNodeMapInBasicBlk(n)
}

func (a *AnalysisCtx) popStmt() *Node {
	cnt := len(a.stmtStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.stmtStack[cnt-1], a.stmtStack[:cnt-1]
	a.stmtStack = rest
	// a.delAstNodeMapInBasicBlk(last)
	return last
}

func (a *AnalysisCtx) pushExpr(n *Node) {
	a.exprStack = append(a.exprStack, n)
	// for _, astNode := range n.AstNodes {
	// 	a.setAstNodeMap(astNode, n)
	// }
	// a.setAstNodeMapInBasicBlk(n)
}

func (a *AnalysisCtx) popExpr() *Node {
	cnt := len(a.exprStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.exprStack[cnt-1], a.exprStack[:cnt-1]
	a.exprStack = rest
	// a.delAstNodeMapInBasicBlk(last)
	return last
}

// func (a *AnalysisCtx) setAstNodeMap(astNode parser.Node, basicBlk *Node) {
// 	a.astNodeInBlock[IdOfAstNode(astNode)] = basicBlk
// }

// func (a *AnalysisCtx) delAstNodeMap(astNode parser.Node) {
// 	delete(a.astNodeInBlock, IdOfAstNode(astNode))
// }

func (a *AnalysisCtx) setAstNodeMapInBasicBlk(basicBlk *Node) {
	for _, astNode := range basicBlk.AstNodes {
		a.astNodeInBlock[IdOfAstNode(astNode)] = basicBlk
	}
}

func (a *AnalysisCtx) delAstNodeMapInBasicBlk(basicBlk *Node) {
	for _, astNode := range basicBlk.AstNodes {
		delete(a.astNodeInBlock, IdOfAstNode(astNode))
	}
}

func (a *AnalysisCtx) basicBlkOfAstNode(astNode parser.Node) *Node {
	return a.astNodeInBlock[IdOfAstNode(astNode)]
}

type Subgraph struct {
	Head   *Node
	Tail   *Node
	Parent *Graph
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

// func enterFn(node parser.Node, key string, ctx *walk.WalkCtx) {
// 	ctx.Extra.(*AnalysisCtx).enterGraph()
// }

// func leaveFn(node parser.Node, key string, ctx *walk.WalkCtx) {
// 	ctx.Extra.(*AnalysisCtx).leaveGraph()
// }

func link(a *AnalysisCtx, from *Node, kind EdgeKind, to *Node, override bool) {
	if from == nil || to == nil {
		return
	}

	dst := to.unwrapIn()

	// link `to` to `from` as well as merge the first node of
	// `to` into `from`
	if len(from.Out) == 1 && len(to.In) == 1 {
		from = from.unwrapOut()
		from.AstNodes = append(from.AstNodes, dst.AstNodes...)
		from.Out = dst.Out
		for _, edge := range from.Out {
			edge.Src = from
		}
		return
	}

	// this branch used to handle the connection between the tail of the node with the exit of the
	// node(sometimes used to debug), the tail of the node often has multiple outlets, the connection
	// is seq, the dest node has to be merged into the source node of the connection if the dest has
	// only one inlet as well as there is only one astNode in that dest node
	if kind == EK_SEQ && len(to.In) == 1 && len(to.AstNodes) == 1 {
		from = from.seqOutEdge().Src
		from.AstNodes = append(from.AstNodes, dst.AstNodes...)
		from.Out = dst.Out
		for _, edge := range from.Out {
			edge.Src = from
		}
		return
	}

	for _, edge := range from.Out {
		if (edge.Dst == nil || override) && (kind == EK_NONE || edge.Kind&kind != 0) {
			if kind != EK_NONE && edge.Src == dst {
				edge.Kind |= EK_LOOP
			}
			edge.Dst = dst

			if edge.Kind&EK_UNREACHABLE != 0 {
				for _, edge := range dst.Out {
					edge.Kind |= EK_UNREACHABLE
				}
			}
		}
	}

	for _, edge := range dst.In {
		if edge.Src == nil {
			edge.Src = from
		}
	}
}

func vnode(a *AnalysisCtx, from *Node, to *Node) *Node {
	if to == nil {
		return from
	}

	var vn *Node
	if from.isVirtual() {
		vn = from
	} else {
		vn = a.graph.newNode()
	}
	vn.In = from.In
	vn.Out = append(to.Out, from.xJmpOutEdges()...)
	return vn
}

func pushAtomNode(node parser.Node, key string, ctx *walk.VisitorCtx) {
	ac := analysisCtx(ctx)
	n := ac.graph.newNode()
	n.AstNodes = append(n.AstNodes, node)
	ac.pushExpr(n)
}

// below stmts and exprs have the condJmp(conditional-jump) semantic:
// - [x] logicAnd
// - [x] logicOr
// - [x] if
// - [x] for
// - [x] while
// - [x] doWhile

// - [x] Loop
// - [x] Test

// below stmts have the unCondJmp(unconditional-jump) semantic
// - [ ] contine
// - [ ] break
// - [ ] return
// - [ ] callExpr

func isLoopTyp(t parser.NodeType) bool {
	return t == parser.N_STMT_FOR || t == parser.N_STMT_WHILE || t == parser.N_STMT_DO_WHILE
}

func isAtom(t parser.NodeType) bool {
	_, ok := walk.AtomNodeTypes[t]
	return ok
}

func (a *Analysis) init() {
	a.WalkCtx.Extra = newAnalysisCtx()

	walk.AddBeforeListener(&a.WalkCtx.Listeners, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		ac := analysisCtx(ctx)

		astTyp := node.Type()
		pAstTyp := ctx.ParentNodeType()

		if isAtom(astTyp) && pAstTyp != parser.N_STMT_LABEL {
			pushAtomNode(node, key, ctx)
		}

		if astTyp.IsStmt() || astTyp == parser.N_PROG {
			enter := ac.newEnter(node, "")

			if pAstTyp == parser.N_STMT_IF || pAstTyp == parser.N_STMT_LABEL || isLoopTyp(pAstTyp) {
				// just push without connectting to the `prevStmt` to imitate a new branch
				ac.pushStmt(enter)
			} else {
				prev := ac.popStmt()
				link(ac, prev, EK_NONE, enter, false)
				ac.pushStmt(vnode(ac, prev, enter))
			}

			// }
			// if isLoopTyp(astTyp) {
			// 	if pAstTyp == parser.N_STMT_LABEL {
			// 		enter.newLoopIn()
			// 	}

			// 	// record the loops, key is the id of their parent label node
			// 	c := ctx
			// 	pn := c.ParentNode()
			// 	for pn != nil && pn.Type() == parser.N_STMT_LABEL {
			// 		ac.astNodeInBlock[IdOfAstNode(pn)] = enter
			// 		c = c.Parent
			// 		if c == nil {
			// 			break
			// 		}
			// 		pn = c.Node
			// 	}
		} else if (astTyp.IsExpr() || astTyp == parser.N_VAR_DEC) && !isAtom(astTyp) {
			ac.pushExpr(ac.newEnter(node, ""))
		}
	})

	walk.AddAfterListener(&a.WalkCtx.Listeners, func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		ac := analysisCtx(ctx)

		astTyp := node.Type()
		pAstTyp := ctx.ParentNodeType()
		switch astTyp {
		case parser.N_EXPR_BIN:
			n := node.(*parser.BinExpr)

			rhs := ac.popExpr()
			exit := ac.newExit(node, "")
			link(ac, rhs, EK_SEQ, exit, false)
			rhs = vnode(ac, rhs, exit)

			lhs := ac.popExpr()
			enter := ac.popExpr()
			link(ac, enter, EK_NONE, lhs, false)
			lhs = vnode(ac, enter, lhs)

			op := n.Op()
			logic := true
			if op == parser.T_AND {
				lhs.newJmp(EK_JMP_FALSE)
			} else if op == parser.T_OR {
				lhs.newJmp(EK_JMP_TRUE)
			} else {
				logic = false
			}

			link(ac, lhs, EK_SEQ, rhs, false)

			if op == parser.T_OR {
				link(ac, lhs, EK_JMP_FALSE, rhs, false)
			}

			vn := ac.graph.newNode()
			vn.In = lhs.In
			if logic {
				vn.Out = append(rhs.Out, lhs.xOutEdges()...)
			} else {
				vn.Out = rhs.Out
			}
			ac.pushExpr(vn)

		case parser.N_STMT_IF:
			n := node.(*parser.IfStmt)
			var enter, test, cons, alt *Node
			if n.Alt() != nil {
				alt = ac.popStmt()
			}
			cons = ac.popStmt()
			test = ac.popExpr()
			enter = ac.popStmt()

			link(ac, enter, EK_NONE, test, false)
			test = vnode(ac, enter, test)

			test.newJmp(EK_JMP_FALSE)
			link(ac, test, EK_SEQ, cons, false)
			link(ac, test, EK_JMP_TRUE, cons, false)
			if alt != nil {
				link(ac, test, EK_JMP_FALSE, alt, false)
			}

			vn := ac.graph.newNode()
			vn.In = enter.In

			vn.Out = append(vn.Out, test.xJmpOutEdges()...)
			vn.Out = append(vn.Out, cons.xOutEdges()...)
			if alt != nil {
				vn.Out = append(vn.Out, alt.xOutEdges()...)
			}

			exit := ac.newExit(node, "")
			link(ac, vn, EK_NONE, exit, false)
			vn = vnode(ac, vn, exit)
			ac.pushStmt(vn)

		case parser.N_STMT_EXPR:
			expr := ac.popExpr()
			prev := ac.popStmt()
			link(ac, prev, EK_SEQ, expr, false)

			exit := ac.newExit(node, "")
			link(ac, expr, EK_NONE, exit, false)
			ac.pushStmt(vnode(ac, prev, exit))

		case parser.N_EXPR_UPDATE, parser.N_EXPR_PAREN:
			expr := ac.popExpr()
			enter := ac.popExpr()
			link(ac, enter, EK_NONE, expr, false)
			exit := ac.newExit(node, "")
			link(ac, expr, EK_NONE, exit, false)
			ac.pushExpr(vnode(ac, enter, exit))

		case parser.N_VAR_DEC:
			n := node.(*parser.VarDec)
			var id, init *Node
			if n.Init() != nil {
				init = ac.popExpr()
			}
			id = ac.popExpr()
			enter := ac.popExpr()
			link(ac, enter, EK_NONE, id, false)

			exit := ac.newExit(node, "")
			if init != nil {
				link(ac, id, EK_NONE, init, false)
				link(ac, init, EK_NONE, exit, false)
			} else {
				link(ac, id, EK_NONE, exit, false)
			}
			vn := vnode(ac, enter, exit)

			prev := ac.popStmt()
			link(ac, prev, EK_NONE, vn, false)
			ac.pushStmt(vnode(ac, prev, vn))

		case parser.N_STMT_VAR_DEC:
			if pAstTyp == parser.N_STMT_FOR && key == "Init" {
				ac.pushExpr(ac.popStmt())
			}

		case parser.N_STMT_FOR:
			n := node.(*parser.ForStmt)
			var enter, init, test, update, body *Node
			if n.Body() != nil {
				body = ac.popStmt()
			}
			if n.Update() != nil {
				update = ac.popExpr()
			}
			if n.Test() != nil {
				test = ac.popExpr()
				test.newJmp(EK_JMP_FALSE)
				test.newLoopIn()
			}
			if n.Init() != nil {
				init = ac.popExpr()
			}
			enter = ac.popStmt()

			if init != nil {
				link(ac, enter, EK_NONE, init, false)
				link(ac, init, EK_SEQ, test, false)
			}

			if test != nil {
				if init == nil {
					link(ac, enter, EK_NONE, test, false)
				}
				link(ac, test, EK_SEQ, body, false)
			} else {
				body.newLoopIn()
				if init == nil {
					link(ac, enter, EK_SEQ, body, false)
				} else {
					link(ac, init, EK_SEQ, body, false)
				}
			}

			if update != nil {
				link(ac, body, EK_SEQ, update, false)
				if test != nil {
					link(ac, update, EK_NONE, test, false)
				} else {
					link(ac, update, EK_NONE, body, false)
				}
				update.mrkSeqOutAsLoop()
			} else {
				if test != nil {
					link(ac, body, EK_NONE, test, false)
				} else {
					link(ac, body, EK_SEQ, body, false)
				}
				body.mrkSeqOutAsLoop()
			}

			vn := ac.graph.newNode()
			vn.In = enter.In
			if test != nil {
				vn.Out = append(vn.Out, test.xJmpOutEdges()...)
			}
			vn.Out = append(vn.Out, body.xOutEdges()...)

			exit := ac.newExit(node, "")
			link(ac, vn, EK_NONE, exit, false)
			ac.pushStmt(vnode(ac, vn, exit))

		case parser.N_STMT_WHILE:
			body := ac.popStmt()
			test := ac.popExpr()
			enter := ac.popStmt()

			test.newJmp(EK_JMP_FALSE)
			test.newLoopIn()
			link(ac, enter, EK_SEQ, test, false)

			link(ac, test, EK_SEQ, body, false)
			link(ac, test, EK_JMP_TRUE, body, false)
			link(ac, body, EK_NONE, test, false)
			body.mrkSeqOutAsLoop()

			vn := ac.graph.newNode()
			vn.In = enter.In
			if test != nil {
				vn.Out = append(vn.Out, test.xJmpOutEdges()...)
			}
			vn.Out = append(vn.Out, body.xOutEdges()...)

			exit := ac.newExit(node, "")
			link(ac, vn, EK_NONE, exit, false)
			ac.pushStmt(vnode(ac, vn, exit))

		case parser.N_STMT_DO_WHILE:
			test := ac.popExpr()
			body := ac.popStmt()
			enter := ac.popStmt()

			test.newJmp(EK_JMP_FALSE)
			body.newLoopIn()
			link(ac, enter, EK_NONE, body, false)
			link(ac, body, EK_SEQ, test, false)

			link(ac, test, EK_SEQ, body, false)
			test.mrkSeqOutAsLoop()
			link(ac, test, EK_JMP_TRUE, body, false)

			vn := ac.graph.newNode()
			vn.In = enter.In
			if test != nil {
				vn.Out = append(vn.Out, test.xJmpOutEdges()...)
			}
			vn.Out = append(vn.Out, body.xOutEdges()...)

			exit := ac.newExit(node, "")
			link(ac, vn, EK_NONE, exit, false)
			ac.pushStmt(vnode(ac, vn, exit))

		case parser.N_STMT_LABEL:
			body := ac.popStmt()
			enter := ac.popStmt()

			link(ac, enter, EK_NONE, body, false)

			exit := ac.newExit(node, "")
			link(ac, body, EK_NONE, exit, false)

			ac.pushStmt(vnode(ac, enter, exit))

		case parser.N_STMT_CONT:
			prev := ac.popStmt()

			n := node.(*parser.ContStmt)
			var target *Node
			if n.Label() != nil {
				// graph does not need to include label name
				ac.popExpr()
				target = ac.basicBlkOfAstNode(n.Target())
			}

			exit := ac.newExit(node, "")
			link(ac, prev, EK_NONE, exit, false)

			exit.newLoopOut()
			link(ac, exit, EK_LOOP, target, false)

			vn := vnode(ac, prev, exit)
			vn.mrkSeqOutAsUnreachable()
			ac.pushStmt(vn)

		case parser.N_PROG, parser.N_STMT_BLOCK:
			prev := ac.popStmt()
			exit := ac.newExit(node, "")
			link(ac, prev, EK_NONE, exit, false)
			ac.pushStmt(vnode(ac, prev, exit))

		default:
			if !isAtom(astTyp) {
				prev := ac.popStmt()
				exit := ac.newExit(node, "")
				link(ac, prev, EK_NONE, exit, false)
				ac.pushStmt(vnode(ac, prev, exit))
			}
		}

		if astTyp.IsExpr() {
			if pAstTyp == parser.N_STMT_WHILE && key == "Test" {
				// record the loops, key is the id of their parent label node
				mapUpperLabelToBlkBlock(ctx.Parent, ac, ac.lastExpr())
			} else if pAstTyp == parser.N_STMT_FOR && key == "Init" {
				mapUpperLabelToBlkBlock(ctx.Parent, ac, ac.lastExpr())
			} else if pAstTyp == parser.N_STMT_DO_WHILE && key == "Body" {
				mapUpperLabelToBlkBlock(ctx.Parent, ac, ac.lastStmt())
			}
			// TODO: for_in
		}
	})
}

func mapUpperLabelToBlkBlock(c *walk.VisitorCtx, ac *AnalysisCtx, blkBlock *Node) {
	if c == nil {
		return
	}
	pn := c.ParentNode()
	for pn != nil && pn.Type() == parser.N_STMT_LABEL {
		ac.astNodeInBlock[IdOfAstNode(pn)] = blkBlock
		c = c.Parent
		if c == nil {
			break
		}
		pn = c.Node
	}
}

func analysisCtx(ctx *walk.VisitorCtx) *AnalysisCtx {
	return ctx.WalkCtx.Extra.(*AnalysisCtx)
}

func (a *Analysis) Analyze() {
	walk.VisitNode(a.WalkCtx.Root, "", a.WalkCtx.VisitorCtx())
}

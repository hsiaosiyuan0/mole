package analysis

import (
	"fmt"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

type IdGen struct {
	prefix string
	seed   int
}

func newIdGen(prefix string) *IdGen {
	return &IdGen{prefix: prefix, seed: 0}
}

func (i *IdGen) next() string {
	i.seed += 1
	return fmt.Sprintf("%s%d", i.prefix, i.seed)
}

type EdgeKind uint8

const (
	EK_NONE EdgeKind = 0
	EK_SEQ  EdgeKind = 1 << iota
	EK_JMP_FALSE
	EK_JMP_TRUE
	EK_LOOP
)

func (k EdgeKind) String() string {
	switch k {
	case EK_JMP_TRUE:
		return "JMP_TRUE"
	case EK_JMP_FALSE:
		return "JMP_FALSE"
	}
	return ""
}

type Edge struct {
	Kind EdgeKind
	Src  *Node
	Dst  *Node
}

func (e *Edge) isSeq() bool {
	return e.Kind&EK_JMP_TRUE == 0 && e.Kind&EK_JMP_FALSE == 0
}

func (e *Edge) Dot() string {
	s := "initial"
	if e.Src != nil {
		s = e.Src.Id
	}
	d := "final"
	if e.Dst != nil {
		d = e.Dst.Id
	}
	return fmt.Sprintf("%s->%s [ xlabel=\"%s\" ];\n", s, d, e.Kind.String())
}

func nodesToString(nodes []parser.Node) string {
	var b strings.Builder
	for _, node := range nodes {
		b.WriteString(nodeToString(node))
	}
	return b.String()
}

func nodeToString(node parser.Node) string {
	switch node.Type() {
	case parser.N_NAME:
		return fmt.Sprintf("%s(%s)\\n", node.Type().String(), node.Loc().Text())
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
		return fmt.Sprintf("%s(%s):%s\\n", typ.String(), n.astNode.(*parser.BinExpr).OpText(), enter)
	}
	return fmt.Sprintf("%s:%s\\n", typ.String(), enter)
}

func (n *DebugNode) Type() parser.NodeType {
	return N_CFG_DEBUG
}

func (n *DebugNode) Loc() *parser.Loc {
	return nil
}

type Node struct {
	Id       string
	AstNodes []parser.Node
	In       []*Edge
	Out      []*Edge
}

func (e *Node) Dot() string {
	return fmt.Sprintf("%s[label=\"%s\"];\n", e.Id, nodesToString(e.AstNodes))
}

func (n *Node) seqOutEdge() *Edge {
	for _, edge := range n.Out {
		if edge.isSeq() {
			return edge
		}
	}
	panic("unreachable")
}

// func (n *Node) forward(to *Node) {
// 	for _, edge := range n.Out {
// 		edge.Dst = to
// 	}
// }

// func (n *Node) backward(to *Node) {
// 	dst := n.In[0].Dst
// 	for _, edge := range to.Out {
// 		edge.Dst = dst
// 	}
// }

// func (n *Node) mergeBackword(to *Node) *Node {
// 	if len(to.Out) == 1 && len(n.In) == 1 {
// 		to.AstNodes = append(to.AstNodes, n.In[0].Dst.AstNodes...)
// 		to.Out = n.In[0].Dst.Out
// 		for _, edge := range to.Out {
// 			edge.Src = to
// 		}
// 	} else {
// 		n.backward(to)
// 	}
// 	return to
// }

// func (n *Node) OutTrue(create bool) *Edge {
// 	for _, edge := range n.Out {
// 		if edge.Kind&EK_TRUE != 0 {
// 			return edge
// 		}
// 	}
// 	if !create {
// 		return nil
// 	}
// 	edge := &Edge{}
// 	edge.Src = n
// 	edge.Kind |= EK_TRUE
// 	n.Out = append(n.Out, edge)
// 	return edge
// }

// func (n *Node) False(create bool) []*Edge {
// 	ret := []*Edge{}
// 	for _, edge := range n.Out {
// 		if edge.Kind&EK_FALSE != 0 {
// 			ret = append(ret, edge)
// 		}
// 	}
// 	if create && len(ret) == 0 {
// 		edge := &Edge{}
// 		edge.Src = n
// 		edge.Kind |= EK_FALSE
// 		n.Out = append(n.Out, edge)
// 		ret = append(ret, edge)
// 	}
// 	return ret
// }

func (n *Node) rareOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range n.Out {
		if edge.Dst == nil {
			ret = append(ret, edge)
		}
	}
	return ret
}

func (n *Node) rareJmpOutEdges() []*Edge {
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
	n.Out = append(n.Out, edge)
}

func (n *Node) jmpOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range n.Out {
		if !edge.isSeq() {
			ret = append(ret, edge)
		}
	}
	return ret
}

type Graph struct {
	Id     string
	Begin  *Node
	End    *Node
	Parent *Graph
	Subs   []*Graph
	IdGen  *IdGen
}

func newGraph() *Graph {
	g := &Graph{
		IdGen: newIdGen("n"),
		Subs:  make([]*Graph, 0),
	}
	g.Begin = g.newNode()
	g.End = g.Begin
	return g
}

func (g *Graph) newNode() *Node {
	n := &Node{
		Id:       g.IdGen.next(),
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

func (g *Graph) Dot() string {
	nodes := []*Node{}
	edges := g.Begin.In
	whites := []*Node{g.Begin}

	dupNodes := map[string]bool{}
	for len(whites) > 0 {
		cnt := len(whites)
		last, rest := whites[cnt-1], whites[:cnt-1]
		whites = rest

		if _, ok := dupNodes[last.Id]; ok {
			continue
		}

		nodes = append(nodes, last)
		dupNodes[last.Id] = true

		if last.Out != nil {
			edges = append(edges, last.Out...)
			for _, edge := range last.Out {
				if edge.Dst != nil {
					whites = append(whites, edge.Dst)
				}
			}
		}
	}

	var b strings.Builder

	b.WriteString(`digraph {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]

initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];

`)

	for _, node := range nodes {
		b.WriteString(node.Dot())
	}

	dup := map[string]bool{}
	for _, edge := range edges {
		s := edge.Dot()
		if _, ok := dup[s]; ok {
			continue
		}
		b.WriteString(s)
		dup[s] = true
	}

	b.WriteString("\n}\n\n")
	return b.String()
}

type AnalysisCtx struct {
	graph     *Graph
	exprStack []*Node
	stmtStack []*Node
}

func newAnalysisCtx() *AnalysisCtx {
	a := &AnalysisCtx{
		graph:     newGraph(),
		exprStack: make([]*Node, 0),
		stmtStack: make([]*Node, 0),
	}
	a.stmtStack = append(a.stmtStack, a.graph.Begin)
	return a
}

func (a *AnalysisCtx) enterGraph() {
	s := newGraph()
	s.IdGen = newIdGen(a.graph.Id + "_")
	s.Parent = a.graph
	a.graph.Subs = append(a.graph.Subs, s)
	a.graph = s
}

func (a *AnalysisCtx) leaveGraph() {
	a.graph = a.graph.Parent
}

func (a *AnalysisCtx) pushExpr(n *Node) {
	a.exprStack = append(a.exprStack, n)
}

func (a *AnalysisCtx) popExpr() *Node {
	cnt := len(a.exprStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.exprStack[cnt-1], a.exprStack[:cnt-1]
	a.exprStack = rest
	return last
}

// func (a *AnalysisCtx) tailExpr(expr *Node) {
// 	tail := a.popExpr()
// 	if tail == nil {
// 		a.pushExpr(expr)
// 	}

// 	link(a, tail, EK_NONE, expr)

// 	vn := a.graph.newNode()
// 	vn.In = tail.In
// 	vn.Out = append(expr.Out, tail.rareJmpOutEdges()...)
// 	a.pushExpr(vn)
// }

func (a *AnalysisCtx) pushStmt(n *Node) {
	a.stmtStack = append(a.stmtStack, n)
}

func (a *AnalysisCtx) popStmt() *Node {
	cnt := len(a.stmtStack)
	if cnt == 0 {
		return nil
	}
	last, rest := a.stmtStack[cnt-1], a.stmtStack[:cnt-1]
	a.stmtStack = rest
	return last
}

func (a *AnalysisCtx) tailStmt(stmt *Node) {
	tail := a.popStmt()
	if tail == nil {
		a.pushStmt(stmt)
	}

	link(a, tail, EK_NONE, stmt)

	vn := a.graph.newNode()
	vn.In = tail.In
	vn.Out = append(stmt.Out, tail.rareJmpOutEdges()...)
	a.pushStmt(vn)
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

func enterFn(node parser.Node, key string, ctx *walk.WalkCtx) {
	ctx.Extra.(*AnalysisCtx).enterGraph()
}

func leaveFn(node parser.Node, key string, ctx *walk.WalkCtx) {
	ctx.Extra.(*AnalysisCtx).leaveGraph()
}

func VisitBinExpr(node parser.Node, key string, ctx *walk.WalkCtx) {
	n := node.(*parser.BinExpr)
	ctx.PushVisitorCtx(n, key)
	defer ctx.PopVisitorCtx()

	walk.CallVisitor(walk.N_EXPR_BIN_BEFORE, n, key, ctx)
	defer walk.CallVisitor(walk.N_EXPR_BIN_AFTER, n, key, ctx)

	ac := ctx.Extra.(*AnalysisCtx)
	enter := ac.graph.newNode()
	enter.AstNodes = append(enter.AstNodes, newDebugNode(node, true, ""))

	walk.VisitNode(n.Lhs(), "Lhs", ctx)
	if ctx.Stopped() {
		return
	}

	lhs := ac.popExpr()
	link(ac, enter, EK_NONE, lhs)
	lhs = vnode(ac, enter, lhs)

	walk.VisitNode(n.Rhs(), "Rhs", ctx)
	if ctx.Stopped() {
		return
	}

	rhs := ac.popExpr()
	exit := ac.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(node, false, ""))
	link(ac, rhs, EK_NONE, exit)
	rhs = vnode(ac, rhs, exit)

	op := n.Op()
	if op == parser.T_AND {
		lhs.newJmp(EK_JMP_FALSE)
	} else if op == parser.T_OR {
		lhs.newJmp(EK_JMP_TRUE)
	}

	link(ac, lhs, EK_SEQ, rhs)

	if op == parser.T_OR {
		link(ac, lhs, EK_JMP_FALSE, rhs)
	}

	vn := ac.graph.newNode()
	vn.In = lhs.In
	vn.Out = append(rhs.Out, lhs.rareOutEdges()...)

	ac.pushExpr(vn)
}

func link(a *AnalysisCtx, from *Node, kind EdgeKind, to *Node) {
	// the dst of income edges of `to` should be equal
	dst := to.In[0].Dst

	// link `to` to `from` as well as merge the first node of
	// `to` into `from`
	if len(from.Out) == 1 && len(to.In) == 1 {
		from = from.Out[0].Src // `from` maybe virtual
		from.AstNodes = append(from.AstNodes, dst.AstNodes...)
		from.Out = dst.Out
		for _, edge := range from.Out {
			edge.Src = from
		}
		return
	}

	for _, edge := range from.Out {
		if kind == EK_NONE || edge.Kind&kind != 0 {
			edge.Dst = dst
		}
	}
	for _, edge := range dst.In {
		if edge.Src == nil {
			edge.Src = from
		}
	}
}

func vnode(a *AnalysisCtx, from *Node, to *Node) *Node {
	vn := a.graph.newNode()
	vn.In = from.In
	vn.Out = append(to.Out, from.rareJmpOutEdges()...)
	return vn
}

func VisitIdent(node parser.Node, key string, ctx *walk.WalkCtx) {
	ac := ctx.Extra.(*AnalysisCtx)
	n := ac.graph.newNode()
	n.AstNodes = append(n.AstNodes, node)
	ac.pushExpr(n)
	walk.CallListener(parser.N_NAME, node, key, ctx)
}

func VisitExprStmt(node parser.Node, key string, ctx *walk.WalkCtx) {
	n := node.(*parser.ExprStmt)
	ctx.PushVisitorCtx(n, key)
	defer ctx.PopVisitorCtx()

	walk.CallVisitor(walk.N_STMT_EXPR_BEFORE, n, key, ctx)
	defer walk.CallVisitor(walk.N_STMT_EXPR_AFTER, n, key, ctx)

	ac := ctx.Extra.(*AnalysisCtx)

	enter := ac.graph.newNode()
	enter.AstNodes = append(enter.AstNodes, newDebugNode(node, true, ""))

	walk.VisitNode(n.Expr(), "Expr", ctx)

	exit := ac.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(node, false, ""))

	expr := ac.popExpr()

	link(ac, enter, EK_NONE, expr)
	link(ac, expr, EK_NONE, exit)
	ac.pushStmt(vnode(ac, enter, exit))

	if ctx.Stopped() {
		return
	}
}

func VisitIfStmt(node parser.Node, key string, ctx *walk.WalkCtx) {
	n := node.(*parser.IfStmt)
	ctx.PushVisitorCtx(n, key)
	defer ctx.PopVisitorCtx()

	walk.CallVisitor(walk.N_STMT_IF_BEFORE, n, key, ctx)
	defer walk.CallVisitor(walk.N_STMT_IF_AFTER, n, key, ctx)

	walk.VisitNode(n.Test(), "Test", ctx)
	if ctx.Stopped() {
		return
	}

	ac := ctx.Extra.(*AnalysisCtx)
	enter := ac.graph.newNode()
	enter.AstNodes = append(enter.AstNodes, newDebugNode(node, true, ""))

	test := ac.popExpr()
	test.newJmp(EK_JMP_FALSE)
	link(ac, enter, EK_NONE, test)
	test = vnode(ac, enter, test)

	walk.VisitNode(n.Cons(), "Cons", ctx)
	if ctx.Stopped() {
		return
	}

	cons := ac.popStmt()

	walk.VisitNode(n.Alt(), "Alt", ctx)
	if ctx.Stopped() {
		return
	}

	var alt *Node
	if n.Alt() != nil {
		alt = ac.popStmt()
	}

	link(ac, test, EK_SEQ, cons)
	link(ac, test, EK_JMP_TRUE, cons)

	if alt != nil {
		link(ac, test, EK_JMP_FALSE, alt)
	}

	vn := ac.graph.newNode()
	vn.In = test.In

	vn.Out = append(vn.Out, test.rareJmpOutEdges()...)
	vn.Out = append(vn.Out, cons.rareOutEdges()...)
	if alt != nil {
		vn.Out = append(vn.Out, alt.rareOutEdges()...)
	}

	exit := ac.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(node, false, ""))
	link(ac, vn, EK_NONE, exit)
	vn = vnode(ac, vn, exit)

	ac.pushStmt(vn)
}

func linkStmts(ctx *walk.WalkCtx) {
	ac := ctx.Extra.(*AnalysisCtx)
	stmt := ac.popStmt()
	tail := ac.popStmt()
	link(ac, tail, EK_NONE, stmt)

	vn := ac.graph.newNode()
	vn.In = tail.In
	vn.Out = append(stmt.Out, tail.rareJmpOutEdges()...)
	ac.pushStmt(vn)
}

func VisitProg(node parser.Node, key string, ctx *walk.WalkCtx) {
	n := node.(*parser.Prog)
	ctx.PushVisitorCtx(n, key)
	defer ctx.PopVisitorCtx()

	ac := ctx.Extra.(*AnalysisCtx)
	enter := ac.graph.newNode()
	enter.AstNodes = append(enter.AstNodes, newDebugNode(node, true, ""))
	ac.tailStmt(enter)

	walk.CallVisitor(walk.N_PROG_BEFORE, n, key, ctx)
	defer walk.CallVisitor(walk.N_PROG_AFTER, n, key, ctx)

	walk.VisitNodesWithCb(n, n.Body(), "Body", ctx, linkStmts)

	exit := ac.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(node, false, ""))
	ac.tailStmt(exit)

	if ctx.Stopped() {
		return
	}
}

func VisitBlockStmt(node parser.Node, key string, ctx *walk.WalkCtx) {
	n := node.(*parser.BlockStmt)
	ctx.PushVisitorCtx(n, key)
	defer ctx.PopVisitorCtx()

	ctx.PushScope()
	defer ctx.PopScope()

	walk.CallVisitor(walk.N_STMT_BLOCK_BEFORE, n, key, ctx)
	defer walk.CallVisitor(walk.N_STMT_BLOCK_AFTER, n, key, ctx)

	ac := ctx.Extra.(*AnalysisCtx)
	enter := ac.graph.newNode()
	enter.AstNodes = append(enter.AstNodes, newDebugNode(node, true, ""))
	ac.pushStmt(enter)

	walk.VisitNodesWithCb(n, n.Body(), "Body", ctx, linkStmts)

	exit := ac.graph.newNode()
	exit.AstNodes = append(exit.AstNodes, newDebugNode(node, false, ""))

	end := ac.popStmt()
	link(ac, end, EK_NONE, exit)

	vn := ac.graph.newNode()
	vn.In = end.In
	vn.Out = append(exit.Out, end.rareJmpOutEdges()...)
	ac.pushStmt(vn)

	if ctx.Stopped() {
		return
	}
}

// below stmts and exprs have the condJmp(conditional-jump) semantic:
// - [x] logicAnd
// - [x] logicOr
// - [x] if
// - [ ] for
// - [ ] while
// - [ ] doWhile

// below stmts have the unCondJmp(unconditional-jump) semantic
// - [ ] contine
// - [ ] break
// - [ ] return
// - [ ] callExpr

func (a *Analysis) init() {
	a.WalkCtx.Extra = newAnalysisCtx()

	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_EXPR_BIN, VisitBinExpr)
	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_NAME, VisitIdent)
	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_STMT_EXPR, VisitExprStmt)
	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_STMT_IF, VisitIfStmt)
	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_PROG, VisitProg)
	walk.SetVisitor(&a.WalkCtx.Visitors, walk.N_STMT_BLOCK, VisitBlockStmt)

	walk.AddListener(&a.WalkCtx.Listeners, walk.N_PROG_AFTER, func(node parser.Node, key string, ctx *walk.WalkCtx) {
		ac := ctx.Extra.(*AnalysisCtx)
		fmt.Println(ac.graph.Dot())
	})
}

func (a *Analysis) Analyze() {
	walk.VisitNode(a.WalkCtx.Root, "", a.WalkCtx)
}

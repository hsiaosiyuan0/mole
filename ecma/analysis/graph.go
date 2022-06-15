package analysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/util"
)

type EdgeKind uint8

const (
	EK_NONE EdgeKind = iota
	EK_SEQ
	EK_JMP
)

type EdgeTag uint16

func (f EdgeTag) On(flag EdgeTag) EdgeTag {
	return f | flag
}

func (f EdgeTag) Off(flag EdgeTag) EdgeTag {
	return f & ^flag
}

const (
	ET_NONE  EdgeTag = 0
	ET_JMP_T EdgeTag = 1 << iota
	ET_JMP_F
	ET_JMP_E
	ET_JMP_U
	ET_JMP_P // jumps with the pair outlet
	ET_JMP_N // jumps if nil
	ET_LOOP
	ET_CUT
)

func (t EdgeTag) String() string {
	ts := []string{}
	if t&ET_JMP_T != 0 {
		ts = append(ts, "T")
	}
	if t&ET_JMP_F != 0 {
		ts = append(ts, "F")
	}
	if t&ET_LOOP != 0 {
		ts = append(ts, "L")
	}
	if t&ET_JMP_U != 0 {
		ts = append(ts, "U")
	}
	if t&ET_JMP_E != 0 {
		ts = append(ts, "E")
	}
	if t&ET_JMP_P != 0 {
		ts = append(ts, "P")
	}
	if t&ET_JMP_N != 0 {
		ts = append(ts, "N")
	}
	return strings.Join(ts, ",")
}

func (t EdgeTag) DotColor() string {
	if t&ET_CUT != 0 {
		return "red"
	}
	if t&ET_JMP_T != 0 || t&ET_JMP_F != 0 || t&ET_LOOP != 0 || t&ET_JMP_U != 0 ||
		t&ET_JMP_E != 0 || t&ET_JMP_P != 0 || t&ET_JMP_N != 0 {
		return "orange"
	}
	return "black"
}

type Edge struct {
	Kind EdgeKind
	Tag  EdgeTag
	Src  *Block
	Dst  *Block
}

func (e *Edge) Key() string {
	from := "s"
	if e.Src != nil {
		from = e.Src.DotId()
	}
	to := "e"
	if e.Dst != nil {
		to = e.Dst.DotId()
	}
	return from + "_" + to
}

func (e *Edge) Dot() string {
	s := "initial"
	if e.Src != nil {
		s = e.Src.DotId()
	}
	d := "final"
	if e.Dst != nil {
		d = e.Dst.DotId()
	}
	c := e.Tag.DotColor()
	fromCorner := ""
	toCorner := ""
	if e.Tag&ET_LOOP != 0 {
		fromCorner = ":s"
		toCorner = ":ne"
	}
	return fmt.Sprintf("%s%s->%s%s [xlabel=\"%s\",color=\"%s\"];\n", s, fromCorner, d, toCorner, e.Tag.String(), c)
}

type BlockKind uint8

const (
	BK_NONE BlockKind = iota
	BK_BASIC
	BK_GROUP
	BK_START
)

type Block struct {
	id      uint
	Kind    BlockKind
	Nodes   []parser.Node
	Inlets  map[*Edge]*Edge
	Outlets map[*Edge]*Edge

	graph *Graph
}

func assocAstNode(node parser.Node, b *Block) {
	if node.Type() == N_CFG_DEBUG {
		n := node.(*InfoNode)
		if n.enter {
			b.graph.astNodeToEntry[n.astNode] = b
		} else {
			b.graph.astNodeToExit[n.astNode] = b
		}
	} else {
		b.graph.astNodeToBlock[node] = b
	}
}

func (b *Block) addNode(node parser.Node) {
	assocAstNode(node, b)
	b.Nodes = append(b.Nodes, node)
}

func (b *Block) addNodes(nodes []parser.Node) {
	for _, node := range nodes {
		assocAstNode(node, b)
	}
	b.Nodes = append(b.Nodes, nodes...)
}

func (b *Block) Id() uint {
	return b.id
}

func (b *Block) DotId() string {
	return fmt.Sprintf("b%d", b.id)
}

func (b *Block) addOutlets(edges map[*Edge]*Edge) *Block {
	util.Merge(b.Outlets, edges)
	return b
}

func (b *Block) OutSeqEdge() *Edge {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_SEQ {
			return edge
		}
	}
	return nil
}

func (b *Block) OutJmpEdge(ET EdgeTag) *Edge {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_JMP && edge.Tag&ET != 0 {
			return edge
		}
	}
	return nil
}

func (b *Block) OutEdge(ET EdgeTag) *Edge {
	for _, edge := range b.Outlets {
		if edge.Tag&ET != 0 {
			return edge
		}
	}
	return nil
}

func (b *Block) InSeqEdge() *Edge {
	for _, edge := range b.Inlets {
		if edge.Kind == EK_SEQ {
			return edge
		}
	}
	return nil
}

func (b *Block) InJmpEdge(ET EdgeTag) *Edge {
	for _, edge := range b.Inlets {
		if edge.Tag&ET != 0 {
			return edge
		}
	}
	return nil
}

func FindEdge(edges map[*Edge]*Edge, k EdgeKind, t EdgeTag) *Edge {
	var e *Edge
	for _, edge := range edges {
		if edge.Kind == k && (edge.Tag == ET_NONE || t == ET_NONE || edge.Tag&t != 0) {
			e = edge
			break
		}
	}
	return e
}

func FindEdges(edges map[*Edge]*Edge, k EdgeKind, t EdgeTag) map[*Edge]*Edge {
	ret := map[*Edge]*Edge{}
	for _, edge := range edges {
		if edge.Kind == k && (edge.Tag == ET_NONE || edge.Tag&t != 0) {
			ret[edge] = edge
		}
	}
	return ret
}

func RemoveEdge(edges map[*Edge]*Edge, k EdgeKind, t EdgeTag) {
	ed := FindEdge(edges, k, t)
	if ed != nil {
		delete(edges, ed)
	}
}

func SwitchCase(clauseBlk *Block) (*Block, *Block, bool) {
	test := clauseBlk.seqInEdge().Dst
	body := test.seqOutEdge().Dst
	return test, body, len(test.Nodes) == 1
}

func (b *Block) FindInEdge(k EdgeKind, t EdgeTag, create bool) (*Edge, bool) {
	e := FindEdge(b.Inlets, k, t)

	new := false
	if e == nil && create {
		e = &Edge{k, t, nil, b}
		b.Inlets[e] = e
		new = true
	}
	return e, new
}

func (b *Block) FindHangingInEdge(src *Block, k EdgeKind, t EdgeTag, create bool) (*Edge, bool) {
	e := FindEdge(b.Inlets, k, t)

	new := false
	if (e == nil || e.Src != nil && e.Src != src) && create {
		e = &Edge{k, t, nil, b}
		b.Inlets[e] = e
		new = true
	}
	return e, new
}

func (b *Block) FindOutEdge(k EdgeKind, t EdgeTag, create bool) *Edge {
	e := FindEdge(b.Outlets, k, t)

	if e == nil && create {
		e = &Edge{k, t, b, nil}
		b.Outlets[e] = e
	}
	return e
}

func (b *Block) Dot() string {
	return fmt.Sprintf("%s[label=\"%s\"];\n", b.DotId(), nodesToString(b.Nodes))
}

func (b *Block) onlySeqIn() bool {
	return len(b.Inlets) == 1 && util.PickOne(b.Inlets).Kind == EK_SEQ
}

func (b *Block) onlySeqOut() bool {
	if len(b.Outlets) != 1 {
		return false
	}
	blk := util.PickOne(b.Outlets).Src
	return len(blk.Outlets) == 1 && util.PickOne(blk.Outlets).Kind == EK_SEQ
}

func (b *Block) hasEnter(node parser.Node) bool {
	for _, n := range b.Nodes {
		if n == node {
			return true
		}
	}
	return false
}

func (b *Block) allInfoNode() bool {
	for _, n := range b.Nodes {
		if n.Type() != N_CFG_DEBUG {
			return false
		}
	}
	return true
}

func (b *Block) throwLit() bool {
	if len(b.Nodes) < 2 {
		return false
	}
	c := len(b.Nodes)
	last := b.Nodes[c-1]
	return last.Type() == N_CFG_DEBUG && last.(*InfoNode).astNode.Type() == parser.N_STMT_THROW && b.Nodes[c-2].Type().IsLit()
}

func (b *Block) OnlyInfo() bool {
	for i := len(b.Nodes) - 1; i >= 0; i-- {
		if b.Nodes[i].Type() != N_CFG_DEBUG {
			return false
		}
	}
	return true
}

func (b *Block) NextBlk() *Block {
	if !b.OnlyInfo() {
		return b
	}
	next := b.OutSeqEdge().Dst
	if next == nil {
		return b
	}
	return b.OutSeqEdge().Dst.NextBlk()
}

func (b *Block) unwrapSeqIn() *Block {
	for _, edge := range b.Inlets {
		if edge.Kind == EK_SEQ {
			return edge.Dst
		}
	}

	// the start block have no seq-in edge
	return nil
}

func (b *Block) unwrapSeqOut() *Block {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_SEQ {
			return edge.Src
		}
	}
	// group-block and the end block may have no seq-out edge,
	// so return nil to indicate those cases
	return nil
}

func (b *Block) join(blk *Block) {
	to := blk.unwrapSeqIn()
	from := b.unwrapSeqOut()
	from.addNodes(to.Nodes)

	from.Outlets = to.Outlets
	isCut := from.IsInCut()
	for _, edge := range from.Outlets {
		edge.Src = from
		if isCut && edge.Dst == nil {
			edge.Tag |= ET_CUT
		}
	}
}

func (b *Block) xOutEdges() map[*Edge]*Edge {
	ret := map[*Edge]*Edge{}
	for _, edge := range b.Outlets {
		if edge.Dst == nil {
			ret[edge] = edge
		}
	}
	return ret
}

func (b *Block) xJmpOutEdges() map[*Edge]*Edge {
	ret := map[*Edge]*Edge{}
	for _, edge := range b.Outlets {
		if edge.Kind != EK_SEQ && edge.Dst == nil {
			ret[edge] = edge
		}
	}
	return ret
}

func (b *Block) seqOutEdge() *Edge {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_SEQ {
			return edge
		}
	}
	panic("unreachable")
}

func (b *Block) seqInEdge() *Edge {
	for _, edge := range b.Inlets {
		if edge.Kind == EK_SEQ {
			return edge
		}
	}
	panic("unreachable")
}

// add new jmp branch from the source node of seqOut
func (b *Block) newJmpOut(k EdgeTag) *Edge {
	seq := b.seqOutEdge()
	edge := &Edge{EK_JMP, k, seq.Src, nil}
	seq.Src.Outlets[edge] = edge

	if seq.Src.IsInCut() {
		edge.Tag |= ET_CUT
	}

	// if b is groupBlock, the new added jmp-edge should also be added
	// into `b` otherwise it can not be linked in future process
	if seq.Src != b {
		b.Outlets[edge] = edge
	}
	return edge
}

func (b *Block) hasXOut(k EdgeKind, t EdgeTag) bool {
	for _, edge := range b.Outlets {
		if edge.Kind == k && (t == ET_NONE || edge.Tag&t != 0) {
			return true
		}
	}
	return false
}

// add new loop branch to the dest node of seqIn
func (b *Block) newLoopIn() {
	seq := b.seqInEdge()
	edge := &Edge{EK_JMP, ET_LOOP, nil, seq.Dst}
	seq.Dst.Inlets[edge] = edge

	if seq.Dst != b {
		b.Inlets[edge] = edge
	}
}

func (b *Block) newJmpIn(t EdgeTag) {
	b.newIn(EK_JMP, t)
}

// stmt after try-catch has tow seq-in edges
func (b *Block) newIn(k EdgeKind, t EdgeTag) {
	seq := b.seqInEdge()
	edge := &Edge{k, t, nil, seq.Dst}
	seq.Dst.Inlets[edge] = edge

	if seq.Dst != b {
		b.Inlets[edge] = edge
	}
}

func (b *Block) mrkSeqOutAsJmp(t EdgeTag) {
	edge := b.seqOutEdge()
	edge.Kind = EK_JMP
	edge.Tag |= t
}

func (b *Block) mrkSeqOutAsLoop() {
	b.mrkSeqOutAsJmp(ET_LOOP)
}

func (b *Block) addCutOutEdge() {
	blk := util.PickOne(b.Outlets).Src
	edge := &Edge{EK_SEQ, ET_NONE, blk, nil}
	edge.Tag |= ET_CUT
	blk.Outlets[edge] = edge

	if blk != b {
		b.Outlets[edge] = edge
	}
}

func (b *Block) mrkJmpOutAsLoop(tag EdgeTag) {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_JMP && edge.Tag&tag != 0 {
			edge.Tag |= ET_LOOP
		}
	}
}

func (b *Block) mrkSeqOutAsCut() {
	b.seqOutEdge().Tag |= ET_CUT
}

func (b *Block) IsInCut() bool {
	for _, e := range b.Inlets {
		if e.Tag&ET_CUT == 0 {
			return false
		}
	}
	return true
}

func (b *Block) IsInCutPair() (bool, bool) {
	p := 0
	for _, e := range b.Inlets {
		if e.Tag&ET_CUT == 0 {
			if e.Tag&ET_JMP_P != 0 {
				p += 1
				continue
			}
			return false, false
		}
	}
	if p == 0 {
		return true, false
	}
	if p > 0 {
		for _, e := range b.Outlets {
			if e.Tag&ET_JMP_P != 0 {
				p -= 1
			}
		}
	}
	return p == 0, true
}

func (b *Block) IsCut() bool {
	cut := b.IsInCut()
	if cut {
		return true
	}
	n := 0
	for _, edge := range b.Outlets {
		if edge.Tag&ET_CUT == 0 && edge.Tag&ET_LOOP == 0 {
			n += 1
		}
	}
	return n == 0
}

func (b *Block) HasOutCut() bool {
	for _, edge := range b.Outlets {
		if edge.Tag&ET_CUT != 0 {
			return true
		}
	}
	return false
}

func (b *Block) IsOutCut(to *Block) bool {
	cutEdge := b.OutEdge(ET_CUT)
	return cutEdge != nil && cutEdge.Dst == to
}

func nodeToString(node parser.Node) string {
	switch node.Type() {
	case parser.N_NAME, parser.N_LIT_NUM:
		return fmt.Sprintf("%s(%s)", node.Type().String(), node.Loc().Text())
	case parser.N_JSX_ID:
		n := node.(*parser.JsxIdent)
		return fmt.Sprintf("%s(%s)", node.Type().String(), n.Text())
	case N_CFG_DEBUG:
		return node.(*InfoNode).String()
	}
	return node.Type().String()
}

func nodesToString(nodes []parser.Node) string {
	var b strings.Builder
	for _, node := range nodes {
		b.WriteString(nodeToString(node) + "\\n")
	}
	return b.String()
}

func IdOfAstNode(node parser.Node) string {
	pos := node.Loc().Begin()
	i := ""
	if node.Type() == N_CFG_DEBUG {
		if node.(*InfoNode).enter {
			i = "_0"
		} else {
			i = "_1"
		}
	}
	return fmt.Sprintf("loc%d_%d_%d%s", pos.Line, pos.Col, node.Type(), i)
}

type Graph struct {
	Id     string
	Head   *Block
	Parent *Graph

	// map the label to its target scope
	hangingLabels []parser.Node
	labelLoop     map[parser.Node]int

	// records basic block need to be resolved, key is the id of the scope which includes the basic block
	hangingBrk map[int][]*Block

	// cont jumps to tail of the loop, so when processing the labelled-cont, the jump target is unknown since
	// the tail of the loop has not been processed, for resolving this problem, a placeholder is introduced
	// here for imitating the jump target and the placeholder will be resolved when the tail of the loop
	// being processed
	//
	// the tail of the loop means:
	// - while-stmt, use the test part
	// - do-while, use the test part
	// - for-in, use the test part
	// - for, first use the update part, then the test part, then the head of the loop body
	//
	// key of below map is the loop which has unresolved cont-stmts
	hangingCont map[parser.Node][]*Block

	// map tryStmts to their unresolved throw-block
	hangingThrow    map[parser.Node][]*Block
	hasHangingThrow bool

	blkSeed uint

	// map astNode to its basic block
	astNodeToBlock map[parser.Node]*Block
	astNodeToEntry map[parser.Node]*Block
	astNodeToExit  map[parser.Node]*Block
}

func newGraph() *Graph {
	g := &Graph{
		labelLoop:    map[parser.Node]int{},
		hangingBrk:   map[int][]*Block{},
		hangingCont:  map[parser.Node][]*Block{},
		hangingThrow: map[parser.Node][]*Block{},

		astNodeToBlock: map[parser.Node]*Block{},
		astNodeToEntry: map[parser.Node]*Block{},
		astNodeToExit:  map[parser.Node]*Block{},
	}
	g.Head = g.newStartBlk()
	return g
}

func (g *Graph) BlkOfNode(node parser.Node) *Block {
	return g.astNodeToBlock[node]
}

func (g *Graph) EntryOfNode(node parser.Node) *Block {
	return g.astNodeToEntry[node]
}

func (g *Graph) ExitOfNode(node parser.Node) *Block {
	return g.astNodeToExit[node]
}

func (g *Graph) addHangingThrow(try parser.Node, blk *Block) {
	list := g.hangingThrow[try]
	if list == nil {
		list = make([]*Block, 0)
	}
	g.hangingThrow[try] = append(list, blk)
}

func (g *Graph) addHangingBrk(id int, blk *Block) {
	list := g.hangingBrk[id]
	if list == nil {
		list = make([]*Block, 0)
	}
	g.hangingBrk[id] = append(list, blk)
}

func (g *Graph) addHangingCont(loop parser.Node, blk *Block) {
	list := g.hangingCont[loop]
	if list == nil {
		list = make([]*Block, 0)
	}
	g.hangingCont[loop] = append(list, blk)
}

func (g *Graph) isLoopHasCont(loop parser.Node) bool {
	return len(g.hangingCont[loop]) > 0
}

func (g *Graph) newBasicBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_BASIC,
		Nodes:   []parser.Node{},
		Inlets:  map[*Edge]*Edge{},
		Outlets: map[*Edge]*Edge{},
		graph:   g,
	}
	g.blkSeed += 1
	in := &Edge{Kind: EK_SEQ, Src: nil, Dst: b}
	b.Inlets[in] = in

	out := &Edge{Kind: EK_SEQ, Src: b, Dst: nil}
	b.Outlets[out] = out
	return b
}

func (g *Graph) newStartBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_START,
		Inlets:  map[*Edge]*Edge{},
		Outlets: map[*Edge]*Edge{},
		graph:   g,
	}
	g.blkSeed += 1

	in := &Edge{Kind: EK_SEQ, Src: nil, Dst: b}
	b.Inlets[in] = in

	out := &Edge{Kind: EK_SEQ, Src: b, Dst: nil}
	b.Outlets[out] = out
	return b
}

func (g *Graph) newGroupBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_GROUP,
		Inlets:  nil,
		Outlets: map[*Edge]*Edge{},
		graph:   g,
	}
	g.blkSeed += 1
	return b
}

func (g *Graph) NodesEdges() (map[string]*Block, []string, map[string]*Edge, []string, map[parser.Node]*Block) {
	uniqueBlocks := map[string]*Block{}
	uniqueEdges := map[string]*Edge{}
	astNodeMap := map[parser.Node]*Block{}

	start := util.PickOne(g.Head.Inlets)
	uniqueEdges[start.Key()] = start

	blKKeys := []string{}
	edgeKeys := []string{
		start.Key(),
	}

	whites := []*Block{g.Head}
	for len(whites) > 0 {
		cnt := len(whites)
		last, rest := whites[cnt-1], whites[:cnt-1]
		whites = rest

		id := last.DotId()
		if _, ok := uniqueBlocks[id]; ok {
			continue
		}
		uniqueBlocks[id] = last
		blKKeys = append(blKKeys, id)

		// map astNodes to its basic block
		for _, astNode := range last.Nodes {
			astNodeMap[astNode] = last
		}

		for _, edge := range last.Outlets {
			ek := edge.Key()
			if _, ok := uniqueEdges[ek]; !ok {
				uniqueEdges[ek] = edge
				edgeKeys = append(edgeKeys, ek)
			}

			if edge.Dst != nil {
				whites = append(whites, edge.Dst)
			}
		}
	}

	sort.Strings(blKKeys)
	sort.Strings(edgeKeys)
	return uniqueBlocks, blKKeys, uniqueEdges, edgeKeys, astNodeMap
}

func (g *Graph) Dot() string {
	blocks, blKKeys, edges, edgeKeys, _ := g.NodesEdges()

	var b strings.Builder

	b.WriteString(`digraph G {
node[shape=box,style="rounded,filled",fillcolor=white,fontname="Consolas",fontsize=10];
edge[fontname="Consolas",fontsize=10]
initial[label="",shape=circle,style=filled,fillcolor=black,width=0.25,height=0.25];
final[label="",shape=doublecircle,style=filled,fillcolor=black,width=0.25,height=0.25];
`)

	for _, bk := range blKKeys {
		b.WriteString(blocks[bk].Dot())
	}

	for _, ek := range edgeKeys {
		b.WriteString(edges[ek].Dot())
	}

	b.WriteString("}\n\n")
	return b.String()
}

type InfoNode struct {
	astNode parser.Node
	enter   bool
	info    string
}

func newInfoNode(node parser.Node, enter bool, info string) *InfoNode {
	return &InfoNode{node, enter, info}
}

func (n *InfoNode) String() string {
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
	case parser.N_IMPORT_SPEC:
		node := n.astNode.(*parser.ImportSpec)
		if node.Default() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "Default", enter)
		}
		if node.NameSpace() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "Namespace", enter)
		}
	case parser.N_STMT_EXPORT:
		node := n.astNode.(*parser.ExportDec)
		if node.All() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "All", enter)
		}
		if node.Default() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "Default", enter)
		}
	case parser.N_EXPORT_SPEC:
		node := n.astNode.(*parser.ExportSpec)
		if node.NameSpace() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "Namespace", enter)
		}
	case parser.N_JSX_OPEN:
		node := n.astNode.(*parser.JsxOpen)
		if node.Closed() {
			return fmt.Sprintf("%s(%s):%s", typ.String(), "Closed", enter)
		}
	}
	return fmt.Sprintf("%s:%s", typ.String(), enter)
}

func (n *InfoNode) Type() parser.NodeType {
	return N_CFG_DEBUG
}

func (n *InfoNode) Loc() *parser.Loc {
	return n.astNode.Loc()
}

type LinkFlag uint16

const (
	LF_NONE      LinkFlag = 0
	LF_FORCE_SEP LinkFlag = 1 << iota
	LF_FORCE_JOIN
	LF_OVERWRITE
)

func link(a *AnalysisCtx, from *Block, fromKind EdgeKind, fromTag EdgeTag, toKind EdgeKind, toTag EdgeTag, to *Block, flag LinkFlag) {
	forceSep := flag&LF_FORCE_SEP != 0
	forceJoin := flag&LF_FORCE_JOIN != 0

	if from == nil || to == nil {
		return
	}

	var fromSeq *Block
	if edge := from.OutSeqEdge(); edge != nil {
		fromSeq = edge.Src
	}

	// if `from` has only one outlet then that outlet must be seq, merge the first node of `to` into `from`
	if forceJoin || fromKind == EK_SEQ && !forceSep && to.onlySeqIn() {

		// process `from` which maybe group blk
		if from.onlySeqOut() || forceJoin {
			from.join(to)
			return
		}

		// `to` is basic blk, try to merge it into `from`
		if fromSeq != nil && fromSeq.Kind == BK_BASIC && fromSeq.onlySeqOut() &&
			toKind == EK_SEQ && to.Kind == BK_BASIC && (forceJoin || fromSeq.Kind == BK_BASIC || fromSeq.IsInCut()) {
			fromSeq.join(to)
			return
		}
	}

	to = to.unwrapSeqIn()

	// process reaches here means the `from` maybe group block or basic block which has multiple outlets
	linkEdges(from.Outlets, fromKind, fromTag, toKind, toTag, to, flag)

	if from.IsOutCut(to) {
		for _, edge := range to.Inlets {
			if edge.Src == fromSeq {
				edge.Tag |= ET_CUT
			}
		}
		cut, pair := to.IsInCutPair()
		if cut {
			for _, edge := range to.Outlets {
				if pair {
					if edge.Kind == EK_SEQ {
						edge.Tag |= ET_CUT
					}
				} else {
					edge.Tag |= ET_CUT
				}
			}
		}
	}
}

func linkEdges(fromEdges map[*Edge]*Edge, fromKind EdgeKind, fromTag EdgeTag, toKind EdgeKind, toTag EdgeTag, to *Block, flag LinkFlag) {
	for _, edge := range fromEdges {
		if edge.Kind == fromKind && (fromTag == ET_NONE || edge.Tag&fromTag != 0) {
			if edge.Dst == nil || flag&LF_OVERWRITE != 0 {
				edge.Dst = to
			}
			toEdge, _ := to.FindHangingInEdge(edge.Src, toKind, toTag, true)
			toEdge.Src = edge.Src
			toEdge.Tag = edge.Tag
		}
	}
}

type BlockCb = func(blk *Block)

func IterBlock(blk *Block, cb BlockCb) {
	cb(blk)
	doneMap := map[*Block]bool{blk: true}

	edges := map[*Edge]*Edge{}
	util.Merge(edges, blk.Outlets)
	for {
		if len(edges) == 0 {
			break
		}

		edge := util.TakeOne(edges)
		if _, done := doneMap[edge.Dst]; !done && edge.Dst != nil {
			cb(edge.Dst)
			doneMap[edge.Dst] = true
			util.Merge(edges, edge.Dst.Outlets)
		}
	}
}

func grpBlock(a *AnalysisCtx, from *Block, to *Block) *Block {
	if to == nil {
		return from
	}

	x := from.xJmpOutEdges()

	var vn *Block
	if from.Kind == BK_GROUP {
		vn = from
		vn.Outlets = map[*Edge]*Edge{}
	} else {
		vn = a.graph.newGroupBlk()
		vn.Inlets = from.Inlets
	}

	vn.addOutlets(to.xOutEdges()).addOutlets(x)
	return vn
}

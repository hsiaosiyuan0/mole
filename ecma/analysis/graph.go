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

type EdgeTag uint8

const (
	ET_NONE  EdgeTag = 0
	ET_JMP_T EdgeTag = 1 << iota
	ET_JMP_F
	ET_JMP_E
	ET_JMP_U
	ET_JMP_P
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
	return strings.Join(ts, ",")
}

func (t EdgeTag) DotColor() string {
	if t&ET_CUT != 0 {
		return "red"
	}
	if t&ET_JMP_T != 0 || t&ET_JMP_F != 0 || t&ET_LOOP != 0 || t&ET_JMP_U != 0 || t&ET_JMP_E != 0 || t&ET_JMP_P != 0 {
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
	Inlets  []*Edge
	Outlets []*Edge
}

func (b *Block) Id() uint {
	return b.id
}

func (b *Block) DotId() string {
	return fmt.Sprintf("b%d", b.id)
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

func FindEdge(edges []*Edge, k EdgeKind, t EdgeTag) (*Edge, int) {
	var e *Edge
	idx := -1
	for i, edge := range edges {
		if edge.Kind == k && (edge.Tag == ET_NONE || edge.Tag&t != 0) {
			e = edge
			e.Tag = t
			idx = i
			break
		}
	}
	return e, idx
}

func FindEdges(edges []*Edge, k EdgeKind, t EdgeTag) []*Edge {
	ret := []*Edge{}
	for _, edge := range edges {
		if edge.Kind == k && (edge.Tag == ET_NONE || edge.Tag&t != 0) {
			ret = append(ret, edge)
		}
	}
	return ret
}

func RemoveEdge(edges *[]*Edge, k EdgeKind, t EdgeTag) {
	_, i := FindEdge(*edges, k, t)
	if i != -1 {
		util.RemoveAt(edges, i)
	}
}

func SwitchCase(clauseBlk *Block) (*Block, *Block, bool) {
	test := clauseBlk.seqInEdge().Dst
	body := test.seqOutEdge().Dst
	return test, body, len(test.Nodes) == 1
}

func (b *Block) FindInEdge(k EdgeKind, t EdgeTag, create bool) *Edge {
	e, _ := FindEdge(b.Inlets, k, t)

	if e == nil && create {
		e = &Edge{k, t, nil, b}
		b.Inlets = append(b.Inlets, e)
	}
	return e
}

func (b *Block) FindOutEdge(k EdgeKind, t EdgeTag, create bool) *Edge {
	e, _ := FindEdge(b.Outlets, k, t)

	if e == nil && create {
		e = &Edge{k, t, b, nil}
		b.Outlets = append(b.Outlets, e)
	}
	return e
}

func (b *Block) Dot() string {
	return fmt.Sprintf("%s[label=\"%s\"];\n", b.DotId(), nodesToString(b.Nodes))
}

func (b *Block) onlySeqIn() bool {
	return len(b.Inlets) == 1 && b.Inlets[0].Kind == EK_SEQ
}

func (b *Block) onlySeqOut() bool {
	if len(b.Outlets) != 1 || b.Outlets[0].Kind != EK_SEQ {
		return false
	}
	blk := b.Outlets[0].Src
	return len(blk.Outlets) == 1 && blk.Outlets[0].Kind == EK_SEQ
}

func (b *Block) hasEnter(node parser.Node) bool {
	for _, n := range b.Nodes {
		if n == node {
			return true
		}
	}
	return false
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
	isCut := from.HasOutCut()
	from.Nodes = append(from.Nodes, to.Nodes...)
	from.Outlets = to.Outlets
	for _, edge := range from.Outlets {
		edge.Src = from
		if isCut && edge.Dst == nil {
			edge.Tag |= ET_CUT
		}
	}
}

func (b *Block) xOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range b.Outlets {
		if edge.Dst == nil {
			ret = append(ret, edge)
		}
	}
	return ret
}

func (b *Block) xJmpOutEdges() []*Edge {
	ret := []*Edge{}
	for _, edge := range b.Outlets {
		if edge.Kind != EK_SEQ && edge.Dst == nil {
			ret = append(ret, edge)
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
	seq.Src.Outlets = append(seq.Src.Outlets, edge)

	if seq.Src.IsInCut() {
		edge.Tag |= ET_CUT
	}

	// if b is groupBlock, the new added jmp-edge should also be added
	// into `b` otherwise it can not be linked in future process
	if seq.Src != b {
		b.Outlets = append(b.Outlets, edge)
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
	seq.Dst.Inlets = append(seq.Dst.Inlets, edge)

	if seq.Dst != b {
		b.Inlets = append(b.Inlets, edge)
	}
}

func (b *Block) newJmpIn(t EdgeTag) {
	b.newIn(EK_JMP, t)
}

// stmt after try-catch has tow seq-in edges
func (b *Block) newIn(k EdgeKind, t EdgeTag) {
	seq := b.seqInEdge()
	edge := &Edge{k, t, nil, seq.Dst}
	seq.Dst.Inlets = append(seq.Dst.Inlets, edge)

	if seq.Dst != b {
		b.Inlets = append(b.Inlets, edge)
	}
}

func (b *Block) mrkSeqOutAsLoop() {
	edge := b.seqOutEdge()
	edge.Kind = EK_JMP
	edge.Tag |= ET_LOOP
}

func (b *Block) addCutOutEdge() {
	blk := b.Outlets[0].Src
	edge := &Edge{EK_SEQ, ET_NONE, blk, nil}
	edge.Tag |= ET_CUT
	blk.Outlets = append(b.Outlets, edge)

	if blk != b {
		b.Outlets = append(b.Outlets, edge)
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
	edge := b.InSeqEdge()
	return edge != nil && len(edge.Dst.Inlets) == 1 && edge.Tag&ET_CUT != 0
}

func (b *Block) HasOutCut() bool {
	return len(b.Outlets) > 0 && b.Outlets[0].Tag&ET_CUT != 0
}

func (b *Block) IsOutCut(to *Block) bool {
	cutEdge := b.OutEdge(ET_CUT)
	return cutEdge != nil && cutEdge.Dst == to
}

func nodeToString(node parser.Node) string {
	switch node.Type() {
	case parser.N_NAME, parser.N_LIT_NUM:
		return fmt.Sprintf("%s(%s)", node.Type().String(), node.Loc().Text())
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
	return fmt.Sprintf("loc%d_%d_%d%s", pos.Line(), pos.Column(), node.Type(), i)
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

	blkSeed uint
}

func newGraph() *Graph {
	g := &Graph{
		labelLoop:   map[parser.Node]int{},
		hangingBrk:  map[int][]*Block{},
		hangingCont: map[parser.Node][]*Block{},
	}
	g.Head = g.newStartBlk()
	return g
}

func (g *Graph) addHangingBrk(id int, blk *Block) {
	list := g.hangingBrk[id]
	if list == nil {
		list = make([]*Block, 0)
	}
	g.hangingBrk[id] = append(list, blk)
}

func (g *Graph) addHangingCont(loopNode parser.Node, blk *Block) {
	list := g.hangingCont[loopNode]
	if list == nil {
		list = make([]*Block, 0)
	}
	g.hangingCont[loopNode] = append(list, blk)
}

func (g *Graph) isLoopHasCont(loopNode parser.Node) bool {
	return len(g.hangingCont[loopNode]) > 0
}

func (g *Graph) newBasicBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_BASIC,
		Nodes:   make([]parser.Node, 0),
		Inlets:  make([]*Edge, 0),
		Outlets: make([]*Edge, 0),
	}
	g.blkSeed += 1
	b.Inlets = append(b.Inlets, &Edge{Kind: EK_SEQ, Src: nil, Dst: b})
	b.Outlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: b, Dst: nil})
	return b
}

func (g *Graph) newStartBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_START,
		Inlets:  nil,
		Outlets: make([]*Edge, 0),
	}
	g.blkSeed += 1
	b.Inlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: nil, Dst: b})
	b.Outlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: b, Dst: nil})
	return b
}

func (g *Graph) newGroupBlk() *Block {
	b := &Block{
		id:      g.blkSeed,
		Kind:    BK_GROUP,
		Inlets:  nil,
		Outlets: nil,
	}
	g.blkSeed += 1
	return b
}

func (g *Graph) NodesEdges() (map[string]*Block, []string, map[string]*Edge, []string, map[parser.Node]*Block) {
	uniqueBlocks := map[string]*Block{}
	uniqueEdges := map[string]*Edge{}
	astNodeMap := map[parser.Node]*Block{}

	start := g.Head.Inlets[0]
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
	}
	return fmt.Sprintf("%s:%s", typ.String(), enter)
}

func (n *InfoNode) Type() parser.NodeType {
	return N_CFG_DEBUG
}

func (n *InfoNode) Loc() *parser.Loc {
	return n.astNode.Loc()
}

func link(a *AnalysisCtx, from *Block, fromKind EdgeKind, fromTag EdgeTag, toKind EdgeKind, toTag EdgeTag, to *Block, forceSep bool, forceJoin bool) {
	if from == nil || to == nil {
		return
	}

	// if `from` has only one outlet then that outlet must be seq, merge the first node of `to` into `from`
	if (fromKind == EK_SEQ || fromKind == EK_NONE) && !forceSep && to.onlySeqIn() {

		// process `from` which maybe group blk
		if from.onlySeqOut() || forceJoin {
			from.join(to)
			return
		}

		// `to` is basic blk, try to merge it into `from`
		var fromSeq *Block
		if edge := from.OutSeqEdge(); edge != nil {
			fromSeq = edge.Src
		}
		if fromSeq != nil && fromSeq.Kind == BK_BASIC && fromSeq.onlySeqOut() &&
			toKind == EK_SEQ && to.Kind == BK_BASIC && (forceJoin || fromSeq.Kind == BK_BASIC || fromSeq.IsInCut()) {
			fromSeq.join(to)
			return
		}
	}

	to = to.unwrapSeqIn()

	// process reaches here means the `from` maybe group block or basic block which has multiple outlets
	linkEdges(from.Outlets, fromKind, fromTag, toKind, toTag, to)

	if from.IsOutCut(to) && len(to.Inlets) == 1 {
		for _, edge := range to.Outlets {
			edge.Tag |= ET_CUT
		}
	}
}

func linkEdges(fromEdges []*Edge, fromKind EdgeKind, fromTag EdgeTag, toKind EdgeKind, toTag EdgeTag, to *Block) {
	for _, edge := range fromEdges {
		if (fromKind == EK_NONE && edge.Dst == nil) || (edge.Kind == fromKind && (fromTag == ET_NONE || edge.Tag&fromTag != 0)) {
			edge.Dst = to
			toEdge := to.FindInEdge(edge.Kind, edge.Tag, true)
			toEdge.Src = edge.Src
		}
	}
}

type BlockCb = func(blk *Block)

func IterBlock(blk *Block, cb BlockCb) {
	cb(blk)
	doneMap := map[*Block]bool{blk: true}
	edges := blk.Outlets
	for {
		if len(edges) == 0 {
			break
		}
		first, rest := edges[0], edges[1:]
		if _, done := doneMap[first.Dst]; !done && first.Dst != nil {
			cb(first.Dst)
			doneMap[first.Dst] = true
			rest = append(rest, first.Dst.Outlets...)
		}
		edges = rest
	}
}

func grpBlock(a *AnalysisCtx, from *Block, to *Block) *Block {
	if to == nil {
		return from
	}

	var vn *Block
	if from.Kind == BK_GROUP {
		vn = from
	} else {
		vn = a.graph.newGroupBlk()
		vn.Inlets = from.Inlets
	}

	vn.Outlets = append(to.Outlets, from.xJmpOutEdges()...)
	return vn
}

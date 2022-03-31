package analysis

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
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
	ET_LOOP
	ET_CUT
)

func (t EdgeTag) String() string {
	if t&ET_JMP_T != 0 {
		return "T"
	}
	if t&ET_JMP_F != 0 {
		return "F"
	}
	if t&ET_LOOP != 0 {
		return "L"
	}
	return ""
}

func (t EdgeTag) DotColor() string {
	if t&ET_CUT != 0 {
		return "red"
	}
	if t&ET_JMP_T != 0 || t&ET_JMP_F != 0 || t&ET_LOOP != 0 {
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
	Kind    BlockKind
	Nodes   []parser.Node
	Inlets  []*Edge
	Outlets []*Edge
}

func (b *Block) Id() string {
	if len(b.Nodes) == 0 {
		return ""
	}
	return IdOfAstNode(b.Nodes[0])
}

func (b *Block) DotId() string {
	switch b.Kind {
	case BK_BASIC:
		return IdOfAstNode(b.Nodes[0])
	case BK_START:
		return "loc0"
	case BK_GROUP:
		return "v"
	}
	panic("unreachable")
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

func (b *Block) Dot() string {
	return fmt.Sprintf("%s[label=\"%s\"];\n", b.DotId(), nodesToString(b.Nodes))
}

func (b *Block) canBeJoined() bool {
	return len(b.Inlets) == 1 && b.Inlets[0].Kind&EK_SEQ != 0
}

func (b *Block) joinable() bool {
	return len(b.Outlets) == 1 && b.Outlets[0].Kind&EK_SEQ != 0
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
		if edge.Kind&EK_SEQ != 0 {
			return edge.Dst
		}
	}
	panic("unreachable")
}

func (b *Block) unwrapSeqOut() *Block {
	for _, edge := range b.Outlets {
		if edge.Kind == EK_SEQ {
			return edge.Src
		}
	}
	panic("unreachable")
}

func (b *Block) join(blk *Block) {
	to := blk.unwrapSeqIn()
	from := b.unwrapSeqOut()
	isCutted := from.IsOutCutted()
	from.Nodes = append(from.Nodes, to.Nodes...)
	from.Outlets = to.Outlets
	for _, edge := range from.Outlets {
		edge.Src = from
		if isCutted && edge.Dst == nil {
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
		if edge.Kind&EK_SEQ == 0 && edge.Dst == nil {
			ret = append(ret, edge)
		}
	}
	return ret
}

func (b *Block) seqOutEdge() *Edge {
	for _, edge := range b.Outlets {
		if edge.Kind&EK_SEQ != 0 {
			return edge
		}
	}
	panic("unreachable")
}

func (b *Block) seqInEdge() *Edge {
	for _, edge := range b.Inlets {
		if edge.Kind&EK_SEQ != 0 {
			return edge
		}
	}
	panic("unreachable")
}

// add new jmp branch from the source node of seqOut
func (b *Block) newJmp(k EdgeTag) {
	seq := b.seqOutEdge()
	edge := &Edge{EK_JMP, k, seq.Src, nil}
	seq.Src.Outlets = append(seq.Src.Outlets, edge)

	// if b is groupBlock, the new added jmp-edge should also be added
	// into `b` otherwise it can not be linked in future process
	if seq.Src != b {
		b.Outlets = append(b.Outlets, edge)
	}
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

func (b *Block) mrkSeqOutAsLoop() {
	edge := b.seqOutEdge()
	edge.Kind = EK_JMP
	edge.Tag |= ET_LOOP
}

func (b *Block) mrkJmpOutAsLoop(tag EdgeTag) {
	for _, edge := range b.Outlets {
		if edge.Kind&EK_JMP != 0 && edge.Tag&tag != 0 {
			edge.Tag |= ET_LOOP
		}
	}
}

func (b *Block) mrkSeqOutAsCutted() {
	b.seqOutEdge().Tag |= ET_CUT
}

func (b *Block) IsInCutted() bool {
	edge := b.InSeqEdge()
	return edge != nil && edge.Tag&ET_CUT != 0
}

func (b *Block) IsOutCutted() bool {
	return len(b.Outlets) > 0 && b.Outlets[0].Tag&ET_CUT != 0
}

func newBasicBlk() *Block {
	b := &Block{
		Kind:    BK_BASIC,
		Nodes:   make([]parser.Node, 0),
		Inlets:  make([]*Edge, 0),
		Outlets: make([]*Edge, 0),
	}
	b.Inlets = append(b.Inlets, &Edge{Kind: EK_SEQ, Src: nil, Dst: b})
	b.Outlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: b, Dst: nil})
	return b
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

func newStartBlk() *Block {
	b := &Block{
		Kind:    BK_START,
		Inlets:  nil,
		Outlets: make([]*Edge, 0),
	}
	b.Inlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: nil, Dst: b})
	b.Outlets = append(b.Outlets, &Edge{Kind: EK_SEQ, Src: b, Dst: nil})
	return b
}

func newGroupBlk() *Block {
	b := &Block{
		Kind:    BK_GROUP,
		Inlets:  nil,
		Outlets: nil,
	}
	return b
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
	Subs   []*Graph

	// map the labelled ast node to its basic block
	labelAstMap   map[string]*Block
	hangingLabels []string
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

func newGraph() *Graph {
	g := &Graph{
		Subs:          make([]*Graph, 0),
		labelAstMap:   map[string]*Block{},
		hangingLabels: make([]string, 0),
	}
	g.Head = newStartBlk()
	return g
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

func link(a *AnalysisCtx, from *Block, fromKind EdgeKind, fromTag EdgeTag, toKind EdgeKind, toTag EdgeTag, to *Block) {
	if from == nil || to == nil {
		return
	}

	// if `from` has only one outlet then that outlet must be seq, merge the first node of `to` into `from`
	if from.joinable() && to.canBeJoined() {
		from.join(to)
		return
	}

	// `from` maybe group block, so use `unwrapSeqOut` to get the inner seq block
	fromSeqEdge := from.OutSeqEdge()
	if fromSeqEdge != nil {
		fromSeq := fromSeqEdge.Src
		if fromSeq.Kind == BK_BASIC && len(fromSeq.Outlets) == 1 && fromKind == EK_SEQ && to.Kind == BK_BASIC && to.canBeJoined() {
			from.join(to)
			return
		}
	}

	to = to.unwrapSeqIn()

	// process reaches here means the `from` maybe group block or basic block which has multiple outlets
	isCutted := from.IsOutCutted()
	for _, edge := range from.Outlets {
		if (fromKind == EK_NONE && edge.Dst == nil) || (edge.Kind == fromKind && (fromTag == ET_NONE || edge.Tag&fromTag != 0)) {
			edge.Dst = to
		}
	}

	for _, edge := range to.Inlets {
		if (toKind == EK_NONE && edge.Src == nil) || (edge.Kind == toKind && (toTag == ET_NONE || edge.Tag&toTag != 0)) {
			edge.Src = from
		}
	}

	if isCutted && len(to.Inlets) == 1 {
		for _, edge := range to.Outlets {
			edge.Tag |= ET_CUT
		}
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
		vn = newGroupBlk()
		vn.Inlets = from.Inlets
	}

	vn.Outlets = append(to.Outlets, from.xJmpOutEdges()...)
	return vn
}

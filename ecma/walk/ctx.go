package walk

import (
	"log"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

type WalkCtx struct {
	Visitors         Visitors
	RaiseMissingImpl bool

	Path  bool        // whether to record the path of node
	Extra interface{} // attach biz extra

	Root   parser.Node    // the root node to start walking
	Symtab *parser.SymTab // the symtab associated with the Root

	vc   *VisitorCtx
	stop bool // whether to stop the walk
}

func NewWalkCtx(root parser.Node, symtab *parser.SymTab) *WalkCtx {
	c := &WalkCtx{}
	c.Visitors = DefaultVisitors
	c.Root = root
	c.Symtab = symtab
	c.vc = &VisitorCtx{c, c.vc, 0, nil, root}
	return c
}

func (c *WalkCtx) Enter(node parser.Node, path string) {
	vc := &VisitorCtx{c, c.vc, c.vc.Depth + 1, nil, node}
	if c.vc.Path != nil {
		vc.Path = append(vc.Path, path)
	}
	c.vc = vc
}

func (c *WalkCtx) Leave() {
	c.vc = c.vc.Parent
}

func (c *WalkCtx) Stop() {
	c.stop = true
}

func (c *WalkCtx) Stopped() bool {
	return c.stop
}

type VisitorCtx struct {
	WalkCtx *WalkCtx    // ctx of the entire walk
	Parent  *VisitorCtx // ctx of the parent node

	Depth int         // depth of current node
	Path  []string    // path of current node, if `WalkCtx.Path` is turned on
	Node  parser.Node // current node
}

func VisitNode(n parser.Node, ctx *WalkCtx) {
	if n == nil {
		return
	}
	CallVisitor(VisitorKind(n.Type()), n, ctx)
}

func VisitNodes(ns []parser.Node, ctx *WalkCtx) {
	for _, n := range ns {
		VisitNode(n, ctx)
		if ctx.stop {
			break
		}
	}
}

func CallVisitor(vk VisitorKind, n parser.Node, ctx *WalkCtx) {
	fns := ctx.Visitors[vk]
	if fns == nil {
		if ctx.RaiseMissingImpl {
			log.Fatalf("Missing Impl for NodeType %d with Kind %d", n.Type(), vk)
		}
		return
	}
	for _, fn := range fns {
		fn(n, ctx)
		if ctx.stop {
			break
		}
	}
}

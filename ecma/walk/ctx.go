package walk

import (
	"fmt"
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

	vc       *VisitorCtx // ctx of current running visitor
	scopeIds []int       // 1-based Id of the scope which current Node belongs to, 0 is reserved for the Global scope
	stop     bool        // whether to stop the walk
}

func NewWalkCtx(root parser.Node, symtab *parser.SymTab) *WalkCtx {
	c := &WalkCtx{}
	c.Visitors = DefaultVisitors
	c.Root = root
	c.Symtab = symtab
	c.vc = &VisitorCtx{c, c.vc, []string{}, root}
	c.scopeIds = []int{0}
	return c
}

type ConditionalNewScope interface {
	NewScope() bool
}

func (c *WalkCtx) ScopeId() int {
	return c.scopeIds[len(c.scopeIds)-1]
}

func (c *WalkCtx) Scope() *parser.Scope {
	return c.Symtab.Scopes[uint(c.ScopeId())]
}

func (c *WalkCtx) PushScope() {
	newScope := true
	if v, ok := c.vc.Node.(ConditionalNewScope); ok {
		newScope = v.NewScope()
	}
	if newScope {
		c.scopeIds = append(c.scopeIds, c.ScopeId()+1)
	}
}

func (c *WalkCtx) PopScope() {
	newScope := true
	if v, ok := c.vc.Node.(ConditionalNewScope); ok {
		newScope = v.NewScope()
	}
	if newScope {
		c.scopeIds = c.scopeIds[:len(c.scopeIds)-1]
	}
}

func (c *WalkCtx) PushVisitorCtx(node parser.Node, path string) {
	vc := &VisitorCtx{c, c.vc, nil, node}
	if c.Path {
		vc.Path = append(c.vc.Path, path)
	}
	c.vc = vc
}

func (c *WalkCtx) PopVisitorCtx() {
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

	Path []string    // path of current node, if `WalkCtx.Path` is turned on
	Node parser.Node // current node
}

func VisitNode(n parser.Node, key string, ctx *WalkCtx) {
	if n == nil {
		return
	}
	CallVisitor(VisitorKind(n.Type()), n, key, ctx)
}

func VisitNodes(n parser.Node, ns []parser.Node, key string, ctx *WalkCtx) {
	for i, n := range ns {
		VisitNode(n, fmt.Sprintf("%s[%d]", key, i), ctx)
		if ctx.stop {
			break
		}
	}
}

func CallVisitor(vk VisitorKind, n parser.Node, key string, ctx *WalkCtx) {
	fns := ctx.Visitors[vk]
	if fns == nil {
		if ctx.RaiseMissingImpl {
			log.Fatalf("Missing Impl for NodeType %d with Kind %d", n.Type(), vk)
		}
		return
	}
	for _, fn := range fns {
		fn(n, key, ctx)
		if ctx.stop {
			break
		}
	}
}

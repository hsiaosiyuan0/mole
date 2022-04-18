package walk

import (
	"fmt"
	"log"

	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

type WalkCtx struct {
	Visitors    Visitors
	Listeners   Listeners
	RaiseNoImpl bool

	Path  bool        // whether to record the path of node
	Extra interface{} // attach extra infos for biz logic

	Root   parser.Node    // the root node to start walking
	Symtab *parser.SymTab // the symtab associated with the Root

	vc          *VisitorCtx // ctx of current running visitor
	scopeIdSeed int         // the seed of scope id
	scopeIds    []int       // 1-based Id of the scope which current Node belongs to, 0 is reserved for the Global scope
	stop        bool        // whether to stop the walk
}

func NewWalkCtx(root parser.Node, symtab *parser.SymTab) *WalkCtx {
	c := &WalkCtx{}
	c.Visitors = DefaultVisitors
	c.Listeners = DefaultListeners
	c.Root = root
	c.Symtab = symtab
	c.vc = &VisitorCtx{c, c.vc, []string{}, root}
	c.scopeIds = []int{0}
	return c
}

func (c *WalkCtx) ScopeId() int {
	return c.scopeIds[len(c.scopeIds)-1]
}

func (c *WalkCtx) Scope() *parser.Scope {
	return c.Symtab.Scopes[c.ScopeId()]
}

func (c *WalkCtx) VisitorCtx() *VisitorCtx {
	return c.vc
}

type CondNewScope interface {
	NewScope() bool
}

func (c *WalkCtx) PushScope() {
	newScope := true
	if v, ok := c.vc.Node.(CondNewScope); ok {
		newScope = v.NewScope()
	}
	if newScope {
		c.scopeIdSeed += 1
		c.scopeIds = append(c.scopeIds, c.scopeIdSeed)
	}
}

func (c *WalkCtx) PopScope() {
	newScope := true
	if v, ok := c.vc.Node.(CondNewScope); ok {
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

func (c *VisitorCtx) ScopeId() int {
	return c.WalkCtx.ScopeId()
}

func (c *VisitorCtx) Scope() *parser.Scope {
	return c.WalkCtx.Scope()
}

func (vc *VisitorCtx) ParentNode() parser.Node {
	if vc.Parent != nil {
		return vc.Parent.Node
	}
	return nil
}

func (vc *VisitorCtx) ParentNodeType() parser.NodeType {
	if vc.Parent != nil {
		return vc.Parent.Node.Type()
	}
	return parser.N_ILLEGAL
}

func (vc *VisitorCtx) ParentIsJmp(key string) bool {
	pt := vc.ParentNodeType()
	if pt == parser.N_EXPR_BIN {
		pn := vc.Parent.Node.(*parser.BinExpr)
		op := pn.Op()
		return op == parser.T_AND || op == parser.T_OR
	} else if pt == parser.N_STMT_IF && key == "Test" {
		return true
	}
	return false
}

func VisitNode(n parser.Node, key string, ctx *VisitorCtx) {
	if n == nil {
		return
	}
	ctx.WalkCtx.PushVisitorCtx(n, key)
	CallVisitor(n.Type(), n, key, ctx.WalkCtx.VisitorCtx())
	ctx.WalkCtx.PopVisitorCtx()
}

func VisitNodes(n parser.Node, ns []parser.Node, key string, ctx *VisitorCtx) {
	for i, n := range ns {
		VisitNode(n, fmt.Sprintf("%s[%d]", key, i), ctx)
		if ctx.WalkCtx.stop {
			break
		}
	}
}

type VisitNodesCb = func(ctx *VisitorCtx)

func VisitNodesWithCb(n parser.Node, ns []parser.Node, key string, ctx *VisitorCtx, cb VisitNodesCb) {
	for i, n := range ns {
		VisitNode(n, fmt.Sprintf("%s[%d]", key, i), ctx)
		cb(ctx)
		if ctx.WalkCtx.stop {
			break
		}
	}
}

func CallVisitor(t parser.NodeType, n parser.Node, key string, ctx *VisitorCtx) {
	fn := ctx.WalkCtx.Visitors[t]
	if fn == nil {
		if ctx.WalkCtx.RaiseNoImpl {
			log.Fatalf("Missing visitor Impl for NodeType %d with Kind %d", n.Type(), t)
		}
		return
	}
	fn(n, key, ctx)
}

func CallListener(t parser.NodeType, n parser.Node, key string, ctx *VisitorCtx) {
	fns := ctx.WalkCtx.Listeners[t]
	for _, fn := range fns {
		fn(n, key, ctx)
		if ctx.WalkCtx.stop {
			break
		}
	}
}

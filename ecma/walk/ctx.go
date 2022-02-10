package walk

import "github.com/hsiaosiyuan0/mole/ecma/parser"

type WalkCtx struct {
	Path  bool        // whether to record the path of node
	Extra interface{} // attach biz extra

	Root   parser.Node    // the root node to start walking
	Symtab *parser.SymTab // the symtab associated with the Root
	vc     *VisitorCtx
}

func NewWalkCtx(root parser.Node, symtab *parser.SymTab) *WalkCtx {
	c := &WalkCtx{}
	c.Root = root
	c.Symtab = symtab
	return c
}

func (c *WalkCtx) Enter(node parser.Node, path string) {
	vc := &VisitorCtx{c, c.vc, c.vc.Depth + 1, nil, node, false}
	if c.vc.Path != nil {
		vc.Path = append(vc.Path, path)
	}
	c.vc = vc
}

func (c *WalkCtx) Leave() {
	c.vc = c.vc.Parent
}

type VisitorCtx struct {
	WalkCtx *WalkCtx    // ctx of the entire walk
	Parent  *VisitorCtx // ctx of the parent node

	Depth int         // depth of current node
	Path  []string    // path of current node, if `WalkCtx.Path` is turned on
	Node  parser.Node // current node

	stop bool // whether to stop the walk
}

func (vc *VisitorCtx) Stop() {
	vc.stop = true
}

func (vc *VisitorCtx) Stopped() bool {
	return vc.stop
}

package astutil

import (
	"github.com/hsiaosiyuan0/mole/ecma/exec"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/util"
)

func GetStaticPropertyName(node parser.Node) string {
	node = parser.UnParen(node)
	switch node.Type() {
	case parser.N_EXPR_MEMBER:
		n := node.(*parser.MemberExpr)
		if n.Prop().Type() == parser.N_NAME {
			return n.Prop().(*parser.Ident).Val()
		}
	case parser.N_EXPR_CHAIN:
		n := node.(*parser.ChainExpr)
		return GetStaticPropertyName(n.Expr())
	}
	return ""
}

func GetName(node parser.Node) string {
	if node.Type() != parser.N_NAME {
		return ""
	}
	return node.(*parser.Ident).Val()
}

func GetNodeNames(node parser.Node) []string {
	switch node.Type() {
	case parser.N_EXPR_FN, parser.N_STMT_FN:
		n := node.(*parser.FnDec)
		if n.Id() == nil {
			return nil
		}
		return []string{n.Id().(*parser.Ident).Val()}
	case parser.N_STMT_VAR_DEC:
		n := node.(*parser.VarDecStmt)
		if len(n.DecList()) != 1 {
			return nil
		}
		vd := n.DecList()[0].(*parser.VarDec)
		if vd.Id().Type() == parser.N_NAME {
			return []string{vd.Id().(*parser.Ident).Val()}
		}
	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)
		ret := []string{}
		for _, s := range n.Specs() {
			spec := s.(*parser.ImportSpec)
			ret = append(ret, spec.Local().(*parser.Ident).Val())
		}
		return ret
	case parser.N_STMT_EXPORT:
		n := node.(*parser.ExportDec)
		if n.Default() {
			return []string{"default"}
		}
		if n.All() {
			return []string{"#all"}
		}
		if dec := n.Dec(); dec != nil {
			return GetNodeNames(dec)
		}
	case parser.N_NAME:
		return []string{node.(*parser.Ident).Val()}
	case parser.N_JSX_ID:
		return []string{node.(*parser.JsxIdent).Val()}
	}
	return nil
}

func NamesInDecNode(node parser.Node) (ret []string, all bool) {
	ret = []string{}
	switch node.Type() {
	case parser.N_STMT_VAR_DEC:
		n := node.(*parser.VarDecStmt)
		for _, name := range n.Names() {
			ret = append(ret, name.(*parser.Ident).Val())
		}
	case parser.N_STMT_FN, parser.N_EXPR_FN:
		n := node.(*parser.FnDec)
		if n.Id() != nil {
			ret = append(ret, n.Id().(*parser.Ident).Val())
		}
	case parser.N_STMT_CLASS, parser.N_EXPR_CLASS:
		n := node.(*parser.ClassDec)
		if n.Id() != nil {
			ret = append(ret, n.Id().(*parser.Ident).Val())
		}
	case parser.N_STMT_IMPORT:
		n := node.(*parser.ImportDec)
		ret = []string{}
		for _, s := range n.Specs() {
			spec := s.(*parser.ImportSpec)
			if spec.Default() {
				ret = append(ret, spec.Local().(*parser.Ident).Val())
			} else if spec.NameSpace() {
				all = true
			} else {
				ret = append(ret, spec.Id().(*parser.Ident).Val())
			}
		}
	case parser.N_STMT_EXPORT:
		n := node.(*parser.ExportDec)
		if n.Default() {
			ret = append(ret, "default")
		} else if n.All() {
			all = true
		} else if n.Dec() != nil {
			ret, _ = NamesInDecNode(n.Dec())
		} else {
			for _, spec := range n.Specs() {
				sp := spec.(*parser.ExportSpec)
				if sp.NameSpace() {
					all = true
				} else {
					ret = append(ret, sp.Id().(*parser.Ident).Val())
				}
			}
		}
	}
	return
}

type SwitchBranch struct {
	negative bool
	test     parser.Node
	body     parser.Node
}

func IfStmtToSwitchBranches(node *parser.IfStmt) []*SwitchBranch {
	bs := NodeToSwitchBranches(node.Test())
	bs = append(bs, &SwitchBranch{false, node.Test(), node.Cons()})
	if node.Alt() != nil {
		switch node.Alt().Type() {
		case parser.N_STMT_IF:
			bs = append(bs, IfStmtToSwitchBranches(node.Alt().(*parser.IfStmt))...)
		default:
			nb := &SwitchBranch{true, node.Test(), node.Alt()}
			bs = append(bs, nb)
		}
	}
	return bs
}

func BinExprToSwitchBranch(node *parser.BinExpr) []*SwitchBranch {
	bs := []*SwitchBranch{{false, nil, node.Lhs()}}
	switch node.Op() {
	case parser.T_AND:
		bs = append(bs, &SwitchBranch{false, node.Lhs(), node.Rhs()})
	case parser.T_OR:
		bs = append(bs, &SwitchBranch{true, node.Lhs(), node.Rhs()})
	default:
		bs = append(bs, &SwitchBranch{true, nil, node.Rhs()})
	}
	return bs
}

func CondExprToSwitchBranches(node *parser.CondExpr) []*SwitchBranch {
	bs := []*SwitchBranch{}
	bs = append(bs, &SwitchBranch{false, node.Test(), node.Cons()})
	bs = append(bs, &SwitchBranch{true, node.Test(), node.Alt()})
	return bs
}

func NodeToSwitchBranches(node parser.Node) []*SwitchBranch {
	bs := []*SwitchBranch{}
	switch node.Type() {
	case parser.N_STMT_IF:
		bs = IfStmtToSwitchBranches(node.(*parser.IfStmt))
	case parser.N_EXPR_BIN:
		bs = BinExprToSwitchBranch(node.(*parser.BinExpr))
	case parser.N_EXPR_COND:
		bs = CondExprToSwitchBranches(node.(*parser.CondExpr))
	}
	return bs
}

func SelectTrueBranches(node parser.Node, vars map[string]interface{}, p *parser.Parser) []parser.Node {
	bs := NodeToSwitchBranches(node)
	tbs := []parser.Node{}
	for _, b := range bs {
		if b.test == nil {
			tbs = append(tbs, b.body)
		}
		v, err := exec.NewExprEvaluator(b.test, p).Exec(vars).GetResult()
		if err != nil {
			continue
		}
		bv := exec.ToBool(v)
		if (!b.negative && bv || b.negative && !bv) && b.body != nil {
			tbs = append(tbs, b.body)
		}
	}

	return tbs
}

// the minimal unit of the target nodes is expr
func CollectNodesInTrueBranches(node parser.Node, typ []parser.NodeType, vars map[string]interface{}, p *parser.Parser) []parser.Node {
	ret := []parser.Node{}
	wc := walk.NewWalkCtx(node, nil)

	walkTrueBranches := func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		subs := SelectTrueBranches(node, vars, p)
		for _, sub := range subs {
			walk.VisitNode(sub, key, ctx)
		}
	}

	walk.SetVisitor(&wc.Visitors, parser.N_STMT_IF, walkTrueBranches)
	walk.SetVisitor(&wc.Visitors, parser.N_EXPR_BIN, walkTrueBranches)
	walk.SetVisitor(&wc.Visitors, parser.N_EXPR_COND, walkTrueBranches)

	for _, t := range typ {
		walk.AddNodeAfterListener(&wc.Listeners, t, &walk.Listener{
			Id: "CollectNodesInTrueBranches",
			Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
				ret = append(ret, node)
			},
		})
	}

	walk.VisitNode(node, "", wc.VisitorCtx())

	return ret
}

func IsNodeContains(parent, sub parser.Node) bool {
	ctx := walk.NewWalkCtx(parent, nil)

	ok := false
	fn := &walk.Listener{
		Id: "IsNodeContains",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if node == sub {
				ok = true
				ctx.WalkCtx.Stop()
			}
		},
	}

	walk.AddBeforeListener(&ctx.Listeners, fn)
	walk.VisitNode(parent, "", ctx.VisitorCtx())
	return ok
}

func GetParent(ctx *walk.VisitorCtx, targetTyp, barrierTyp []parser.NodeType) (parser.Node, *walk.VisitorCtx) {
	for {
		pn := ctx.ParentNode()
		if pn == nil {
			break
		}

		pt := pn.Type()
		if util.Includes(targetTyp, pt) {
			return pn, ctx
		}

		if util.Includes(barrierTyp, pt) {
			break
		}

		ctx = ctx.Parent
	}
	return nil, nil
}

func IsPrimitive(node parser.Node) bool {
	nt := node.Type()
	return nt == parser.N_NAME || nt == parser.N_LIT_STR || nt == parser.N_LIT_NUM || nt == parser.N_LIT_BOOL || nt == parser.N_LIT_NULL || nt == parser.N_LIT_REGEXP
}

func IsPlainArr(node parser.Node) bool {
	if node.Type() != parser.N_LIT_ARR {
		return false
	}
	els := node.(*parser.ArrLit).Elems()
	for _, el := range els {
		if el != nil && !IsPlainObj(el) {
			return false
		}
	}
	return true
}

func IsPlainObjLit(node parser.Node) bool {
	if node.Type() != parser.N_LIT_OBJ {
		return false
	}
	props := node.(*parser.ObjLit).Props()
	for _, prop := range props {
		if prop.Type() == parser.N_PROP {
			if !IsPlainObj(prop.(*parser.Prop).Val()) {
				return false
			}
		}
	}
	return true
}

func IsPlainObj(node parser.Node) bool {
	return IsPrimitive(node) || IsPlainArr(node) || IsPlainObjLit(node)
}

func GetVarDec(name string, varDecStmt *parser.VarDecStmt) *parser.VarDec {
	for _, n := range varDecStmt.DecList() {
		dec := n.(*parser.VarDec)
		if dec.Id().Type() == parser.N_NAME && dec.Id().(*parser.Ident).Val() == name {
			return dec
		}
	}
	return nil
}

func IsIdentBoundInImportDec(ident *parser.Ident, scope *parser.Scope) *parser.ImportDec {
	name := ident.Val()
	ref := scope.BindingOf(name)
	dec := ref.Dec
	if dec != nil && dec.Type() == parser.N_STMT_IMPORT {
		return dec.(*parser.ImportDec)
	}
	return nil
}

func IsFn(node parser.Node) bool {
	typ := node.Type()
	return typ == parser.N_STMT_FN || typ == parser.N_EXPR_FN || typ == parser.N_EXPR_ARROW
}

type FnDepGraph struct {
	Nodes map[parser.Node]*FnDepNode
}

type FnDepNode struct {
	// the ast node defines this dep node
	Dec parser.Node

	// the refs captured by this fn
	Captures []*parser.Ref

	// the fns which this fn depends on
	Deps []*FnDepNode
}

func BuildFnDepGraph(node parser.Node, symtab *parser.SymTab) *FnDepGraph {
	ctx := walk.NewWalkCtx(node, symtab)

	graph := &FnDepGraph{
		Nodes: map[parser.Node]*FnDepNode{},
	}

	depNodeOfFn := func(node parser.Node) *FnDepNode {
		if graph.Nodes[node] != nil {
			return graph.Nodes[node]
		}
		fnNode := &FnDepNode{}
		graph.Nodes[node] = fnNode
		return fnNode
	}

	fnNodeStk := []*FnDepNode{}

	var fnNode *FnDepNode
	fnBefore := func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		newFn := depNodeOfFn(node)
		newFn.Dec = node
		fnNodeStk = append(fnNodeStk, newFn)

		if fnNode != nil {
			fnNode.Deps = append(fnNode.Deps, newFn)
		}
		fnNode = newFn
	}

	fnAfter := func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		cnt := len(fnNodeStk)
		if cnt > 0 {
			fnNode = fnNodeStk[cnt-1]
			fnNodeStk = fnNodeStk[:cnt-1]
		} else {
			fnNode = nil
		}
	}

	walk.AddListener(&ctx.Listeners, walk.N_STMT_FN_BEFORE, &walk.Listener{
		Id:     "N_STMT_FN_BEFORE",
		Handle: fnBefore,
	})
	walk.AddListener(&ctx.Listeners, walk.N_EXPR_FN_BEFORE, &walk.Listener{
		Id:     "N_EXPR_FN_BEFORE",
		Handle: fnBefore,
	})
	walk.AddListener(&ctx.Listeners, walk.N_EXPR_ARROW_BEFORE, &walk.Listener{
		Id:     "N_EXPR_ARROW_BEFORE",
		Handle: fnBefore,
	})

	walk.AddListener(&ctx.Listeners, walk.N_STMT_FN_AFTER, &walk.Listener{
		Id:     "N_STMT_FN_AFTER",
		Handle: fnAfter,
	})
	walk.AddListener(&ctx.Listeners, walk.N_EXPR_FN_AFTER, &walk.Listener{
		Id:     "N_EXPR_FN_AFTER",
		Handle: fnAfter,
	})
	walk.AddListener(&ctx.Listeners, walk.N_EXPR_ARROW_AFTER, &walk.Listener{
		Id:     "N_EXPR_ARROW_AFTER",
		Handle: fnAfter,
	})

	walk.AddListener(&ctx.Listeners, walk.N_NAME_AFTER, &walk.Listener{
		Id: "N_NAME_AFTER",
		Handle: func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			if fnNode == nil {
				return
			}

			pn := ctx.ParentNode()
			if pn != nil && pn.Type() == parser.N_EXPR_MEMBER && key == "Prop" && pn.(*parser.MemberExpr).Compute() == false {
				return
			}

			name := node.(*parser.Ident).Val()
			scope := ctx.Scope()
			if scope.Local(name) == nil {
				ref := scope.BindingOf(name)
				if ref != nil {
					fnNode.Captures = append(fnNode.Captures, ref)

					dec := ref.Dec
					if dec == nil {
						return
					}

					var fn parser.Node
					if dec.Type() == parser.N_STMT_VAR_DEC {
						varDec := GetVarDec(name, dec.(*parser.VarDecStmt))
						if varDec != nil {
							init := varDec.Init()
							if IsFn(init) {
								fn = init
							}
						}
					} else if IsFn(dec) {
						fn = dec
					}

					if fn != nil {
						dep := depNodeOfFn(fn)
						if !util.Includes(fnNode.Deps, dep) {
							fnNode.Deps = append(fnNode.Deps, dep)
						}
					}
				}
			}
		},
	})

	walk.VisitNode(node, "", ctx.VisitorCtx())

	return graph
}

func IsFnDepsOnNode(graph *FnDepGraph, fn parser.Node, target parser.Node) bool {
	f := graph.Nodes[fn]
	if f == nil {
		return false
	}

	visited := map[*FnDepNode]bool{}
	nodes := []*FnDepNode{f}
	for {
		if len(nodes) == 0 {
			break
		}

		node, rest := nodes[0], nodes[1:]
		nodes = rest
		if visited[node] {
			continue
		}
		visited[node] = true

		for _, cap := range node.Captures {
			if cap.Dec == target {
				return true
			}
		}

		nodes = append(nodes, node.Deps...)
	}

	return false
}

func IdOfLocalInImportDec(local string, n *parser.ImportDec) string {
	for _, s := range n.Specs() {
		spec := s.(*parser.ImportSpec)
		if spec.Local().(*parser.Ident).Val() == local {
			return spec.Id().(*parser.Ident).Val()
		}
	}
	return ""
}

package astutil

import (
	"container/list"

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
			return n.Prop().(*parser.Ident).Text()
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
	return node.(*parser.Ident).Text()
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

func SelectTrueBranches(node parser.Node, vars map[string]interface{}) []parser.Node {
	bs := NodeToSwitchBranches(node)
	tbs := []parser.Node{}
	for _, b := range bs {
		if b.test == nil {
			tbs = append(tbs, b.body)
		}
		v, err := exec.NewExprEvaluator(b.test).Exec(vars).GetResult()
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

func CollectNodesInTrueBranches(node parser.Node, typ []parser.NodeType, vars map[string]interface{}) []parser.Node {
	ret := []parser.Node{}
	nodes := list.New()
	nodes.PushBack(node)

	for node := nodes.Front(); node != nil; node = nodes.Front() {
		n := node.Value.(parser.Node)
		if util.Includes(typ, n.Type()) {
			ret = append(ret, n)
		}

		switch n.Type() {
		case parser.N_STMT_IF, parser.N_EXPR_BIN, parser.N_EXPR_COND:
			subs := SelectTrueBranches(n, vars)
			for _, sub := range subs {
				nodes.PushBack(sub)
			}
		case parser.N_EXPR_ASSIGN:
			nodes.PushBack(n.(*parser.AssignExpr).Rhs())
		case parser.N_STMT_VAR_DEC:
			for _, d := range n.(*parser.VarDecStmt).DecList() {
				nodes.PushBack(d.(*parser.VarDec).Init())
			}
		case parser.N_STMT_EXPR:
			nodes.PushBack(n.(*parser.ExprStmt).Expr())
		case parser.N_STMT_BLOCK:
			for _, sn := range n.(*parser.BlockStmt).Body() {
				nodes.PushBack(sn)
			}
		}
		nodes.Remove(node)
	}
	return ret
}

func IsNodeContains(parent, sub parser.Node) bool {
	ctx := walk.NewWalkCtx(parent, nil)
	ok := true
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

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

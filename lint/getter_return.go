package lint

import (
	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/analysis"
	"github.com/hsiaosiyuan0/mole/ecma/astutil"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
)

type GetterReturn struct{}

func (n *GetterReturn) Name() string {
	return "getter-return"
}

func (n *GetterReturn) Meta() *Meta {
	return &Meta{
		Lang: []string{RL_JS},
		Kind: RK_LINT_SEMANTIC,
		Docs: Docs{
			Desc: "enforce `return` statements in getters",
			Url:  "https://eslint.org/docs/rules/getter-return",
		},
	}
}

var getterReturnOpts = plugin.DefineOptions(GetterReturnOpts{})

func (n *GetterReturn) Options() *plugin.Options {
	return getterReturnOpts
}

func (n *GetterReturn) Validate() *validator.Validate {
	return nil
}

func (n *GetterReturn) Validates() map[int]plugin.Validate {
	return nil
}

type GetterReturnOpts struct {
	AllowImplicit bool `json:"allowImplicit"`
}

func isGetter(node parser.Node, ctx *walk.VisitorCtx) (bool, parser.Node) {
	parent := ctx.ParentNode()
	if parent == nil {
		return false, nil
	}

	pt := parent.Type()
	if pt == parser.N_METHOD && parent.(*parser.Method).PropKind() == parser.PK_GETTER {
		return true, parent.(*parser.Method).Key()
	} else if pt == parser.N_PROP && parent.(*parser.Prop).PropKind() == parser.PK_GETTER {
		return true, parent.(*parser.Prop).Key()
	}

	if pt == parser.N_PROP &&
		parent.(*parser.Prop).Key().Type() == parser.N_NAME &&
		parent.(*parser.Prop).Key().(*parser.Ident).Text() == "get" &&
		ctx.Parent.Parent.Node.Type() == parser.N_LIT_OBJ {

		// process `Object.defineProperty`
		if ctx.Parent.Parent.Parent.Node.Type() == parser.N_EXPR_CALL {
			// `ctx.Parent.Parent.Parent` can never be nil since `ctx.Parent.Parent` has been guarded as an ObjExpr and
			// `ObjExpr` never live alone

			callee := ctx.Parent.Parent.Parent.Node.(*parser.CallExpr).Callee()
			if astutil.GetStaticPropertyName(callee) == "defineProperty" {
				return true, parent.(*parser.Prop).Key()
			}
		}

		if ctx.Parent.Parent.Parent.Node.Type() == parser.N_PROP &&
			ctx.Parent.Parent.Parent.Parent.Node.Type() == parser.N_LIT_OBJ &&
			ctx.Parent.Parent.Parent.Parent.Parent.Node.Type() == parser.N_EXPR_CALL {

			callee := ctx.Parent.Parent.Parent.Parent.Parent.Node.(*parser.CallExpr).Callee()
			if astutil.GetStaticPropertyName(callee) == "defineProperties" {
				return true, parent.(*parser.Prop).Key()
			}
		}
	}

	return false, nil
}

func rets(node parser.Node) []parser.Node {
	switch node.Type() {
	case parser.N_EXPR_FN, parser.N_STMT_FN:
		return node.(*parser.FnDec).Rets()
	case parser.N_EXPR_ARROW:
		return node.(*parser.ArrowFn).Rets()
	}
	return nil
}

func (n *GetterReturn) Create(rc *RuleCtx) map[parser.NodeType]walk.ListenFn {
	fns := map[parser.NodeType]walk.ListenFn{}

	// `dup` holds this relation [fn => error]
	dup := map[parser.Node]bool{}

	interests := []parser.NodeType{parser.N_EXPR_FN, parser.N_STMT_FN, parser.N_EXPR_ARROW}
	handleFn := func(node parser.Node, key string, ctx *walk.VisitorCtx) {
		if _, ok := dup[node]; ok {
			return
		}

		getter, anchor := isGetter(node, ctx)
		if !getter {
			return
		}

		ac := analysis.AsAnalysisCtx(ctx)
		blk := ac.GraphOf(node).ExitOfNode(node)

		ok := true
		for _, edge := range blk.Inlets {
			if edge.Tag == analysis.ET_NONE ||
				(edge.Tag&analysis.ET_CUT == 0 &&
					edge.Tag&analysis.ET_JMP_U == 0 &&
					edge.Tag&analysis.ET_JMP_E == 0) {

				ok = false
				break
			}
		}

		if !ok {
			dup[node] = true
			rc.Report(anchor, "expected to return a value", DL_ERROR)
			return
		}

		allowImplicit := false
		opts := rc.Opts()
		if opts != nil {
			allowImplicit = opts[0].(*GetterReturnOpts).AllowImplicit
		}

		if !allowImplicit {
			for _, ret := range rets(node) {
				if ret.(*parser.RetStmt).Arg() == nil {
					dup[node] = true
					rc.Report(anchor, "expected to return a value", DL_ERROR)
				}
			}
		}
	}

	for _, nt := range interests {
		fns[walk.NodeAfterEvent(nt)] = handleFn
	}

	return fns
}

package linter

import (
	"github.com/hsiaosiyuan0/mole/ecma/analysis"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
)

type NoUnreachable struct{}

func (n *NoUnreachable) Name() string {
	return "no-unreachable"
}

func (n *NoUnreachable) Meta() *Meta {
	return &Meta{
		Lang: []string{RL_JS},
		Kind: RK_LINT_SEMANTIC,
		Docs: Docs{
			Desc: "disallow unreachable code after `return`, `throw`, `continue`, and `break` statements",
			Url:  "https://eslint.org/docs/rules/no-unreachable",
		},
	}
}

func (n *NoUnreachable) Create(rc *RuleCtx) map[parser.NodeType]walk.ListenFn {
	fns := map[parser.NodeType]walk.ListenFn{}

	// scopeId => has unreachable error, which means only one unreachable-error
	// in one scope
	dup := map[int]bool{}

	for nt := range walk.StmtNodeTypes {
		fns[walk.NodeAfterEvent(nt)] = func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			ac := analysis.AsAnalysisCtx(ctx)
			blk := ac.Graph().EntryOfNode(node)
			si := ctx.ScopeId()

			if blk != nil && blk.IsInCut() {
				if _, ok := dup[si]; !ok {
					rc.Report(node, "disallow unreachable code", DL_ERROR)
					dup[si] = true
				}
			}

		}
	}

	return fns
}

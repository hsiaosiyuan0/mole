package linter

import (
	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/analysis"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
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

func (n *NoUnreachable) Options() *plugin.Options {
	return nil
}

func (n *NoUnreachable) Validate() *validator.Validate {
	return nil
}

func (n *NoUnreachable) Validates() map[int]plugin.Validate {
	return nil
}

func (n *NoUnreachable) Create(rc *RuleCtx) map[parser.NodeType]walk.ListenFn {
	fns := map[parser.NodeType]walk.ListenFn{}

	// `dup` holds this relation [scopeId => has unreachable error], which means only one unreachable-error
	// can be reported in one scope
	dup := map[int]bool{}

	for nt := range walk.StmtNodeTypes {
		fns[walk.NodeAfterEvent(nt)] = func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			si := ctx.ScopeId()
			if _, ok := dup[si]; ok {
				return
			}

			ac := analysis.AsAnalysisCtx(ctx)
			blk := ac.Graph().EntryOfNode(node)

			if blk != nil && blk.IsInCut() {
				rc.Report(node, "disallow unreachable code", DL_ERROR)
				dup[si] = true
			}
		}
	}

	return fns
}

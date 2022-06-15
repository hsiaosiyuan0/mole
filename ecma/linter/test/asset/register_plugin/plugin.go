package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/linter"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
	"github.com/hsiaosiyuan0/mole/util"
)

func Register() []linter.RuleFact {
	return []linter.RuleFact{
		&NoAlert{},
	}
}

type NoAlert struct{}

func (n *NoAlert) Name() string {
	return "no-alert"
}

func (n *NoAlert) Meta() *linter.Meta {
	return &linter.Meta{
		Lang: []string{linter.RL_JS},
		Kind: linter.RK_LINT_SEMANTIC,
		Docs: linter.Docs{
			Desc: "",
			Url:  "",
		},
	}
}

func (n *NoAlert) Options() *plugin.Options {
	return nil
}

func (n *NoAlert) Validate() *validator.Validate {
	return nil
}

func (n *NoAlert) Validates() map[int]plugin.Validate {
	return nil
}

var forbids = []string{"alert", "confirm", "prompt"}

func (n *NoAlert) Create(rc *linter.RuleCtx) map[parser.NodeType]walk.ListenFn {
	return map[parser.NodeType]walk.ListenFn{
		walk.NodeBeforeEvent(parser.N_EXPR_CALL): func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			callee := node.(*parser.CallExpr).Callee()
			if callee.Type() == parser.N_NAME && util.Includes(forbids, callee.Loc().Text()) {
				rc.Report(node, "disallow the use of `alert`, `confirm`, and `prompt`", linter.DL_ERROR)
			}
		},
	}
}

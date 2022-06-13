package main

import (
	"fmt"

	"github.com/hsiaosiyuan0/mole/ecma/linter"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
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

func (n *NoAlert) Create(ctx linter.Unit) map[parser.NodeType]walk.ListenFn {
	return map[parser.NodeType]walk.ListenFn{
		walk.NodeBeforeEvent(parser.N_EXPR_CALL): func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			fmt.Println(node.(*parser.CallExpr).Callee().Loc().Text())
		},
	}
}

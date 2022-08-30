package lint

import (
	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
	"github.com/hsiaosiyuan0/mole/util"
)

type NoAlert struct{}

func (n *NoAlert) Name() string {
	return "no-alert"
}

func (n *NoAlert) Meta() *Meta {
	return &Meta{
		Lang: []string{RL_JS},
		Kind: RK_LINT_SEMANTIC,
		Docs: Docs{
			Desc: "disallow the use of `alert`, `confirm`, and `prompt`",
			Url:  "https://eslint.org/docs/rules/no-alert",
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

var alertForbids = []string{"alert", "confirm", "prompt"}

func (n *NoAlert) Create(rc *RuleCtx) map[parser.NodeType]walk.ListenFn {
	return map[parser.NodeType]walk.ListenFn{
		walk.NodeBeforeEvent(parser.N_EXPR_CALL): func(node parser.Node, key string, ctx *walk.VisitorCtx) {
			callee := node.(*parser.CallExpr).Callee()
			if callee.Type() == parser.N_NAME && util.Includes(alertForbids, callee.(*parser.Ident).Val()) {
				rc.Report(node, "disallow the use of `alert`, `confirm`, and `prompt`", DL_ERROR)
			}
		},
	}
}

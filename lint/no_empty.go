package lint

import (
	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
)

type NoEmpty struct{}

func (n *NoEmpty) Name() string {
	return "no-empty"
}

func (n *NoEmpty) Meta() *Meta {
	return &Meta{
		Lang: []string{RL_JS},
		Kind: RK_LINT_SEMANTIC,
		Docs: Docs{
			Desc: "disallow empty block statements",
			Url:  "https://eslint.org/docs/rules/no-empty",
		},
	}
}

type NoEmptyOpts struct {
	AllowEmptyCatch bool `json:"allowEmptyCatch"`
}

var noEmptyOpts = plugin.DefineOptions(NoEmptyOpts{})

func (n *NoEmpty) Options() *plugin.Options {
	return noEmptyOpts
}

func (n *NoEmpty) Validate() *validator.Validate {
	return nil
}

func (n *NoEmpty) Validates() map[int]plugin.Validate {
	return nil
}

func (n *NoEmpty) Create(rc *RuleCtx) map[parser.NodeType]walk.ListenFn {
	fns := map[parser.NodeType]walk.ListenFn{}

	interests := []parser.NodeType{
		parser.N_STMT_BLOCK, parser.N_STMT_SWITCH,
	}

	for _, nt := range interests {
		fns[walk.NodeAfterEvent(nt)] = func(node parser.Node, key string, ctx *walk.VisitorCtx) {

			nt := node.Type()

			switch nt {
			case parser.N_STMT_BLOCK:
				n := node.(*parser.BlockStmt)
				if len(n.Body()) > 0 {
					return
				}
			case parser.N_STMT_WHILE:
				n := node.(*parser.SwitchStmt)
				if len(n.Cases()) > 0 {
					return
				}
			}

			if rc.unit.HasCommentInSpan(node.Range()) {
				return
			}

			opts := rc.Opts()
			parent := ctx.ParentNode()
			if nt == parser.N_STMT_BLOCK &&
				parent != nil && parent.Type() == parser.N_CATCH &&
				opts != nil && opts[0].(*NoEmptyOpts).AllowEmptyCatch {
				return
			}

			rc.Report(node, "disallow empty block statements", DL_ERROR)
		}
	}

	return fns
}

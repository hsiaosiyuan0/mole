package lint

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/ecma/walk"
	"github.com/hsiaosiyuan0/mole/plugin"
)

type Kind uint16

const (
	RK_LINT_SYNTAX Kind = iota
	RK_LINT_SEMANTIC
	RK_LINT_STYLE
	RK_LINT_OTHERS
)

const (
	RL_JS  string = ".js"
	RL_JSX string = ".jsx"
)

type DiagLevel uint16

const (
	DL_OFF DiagLevel = iota
	DL_WARN
	DL_ERROR
)

type Diagnosis struct {
	Lang string
	Rule string

	Loc   *parser.Loc
	Level DiagLevel
	Msg   string
}

func (d *Diagnosis) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Line  uint32 `json:"line"`
		Col   uint32 `json:"col"`
		Lang  string `json:"lang"`
		Rule  string `json:"rule"`
		Level uint16 `json:"level"`
		Msg   string `json:"Msg"`
	}{
		Line:  d.Loc.Begin().Line,
		Col:   d.Loc.Begin().Col,
		Lang:  d.Lang,
		Rule:  d.Rule,
		Level: uint16(d.Level),
		Msg:   d.Msg,
	})
}

type Docs struct {
	Desc string
	Url  string
}

type Meta struct {
	Lang      []string
	Kind      Kind
	DiagLevel DiagLevel
	Docs      Docs
}

type RuleFact interface {
	Name() string
	Meta() *Meta

	Options() *plugin.Options
	Validate() *validator.Validate
	Validates() map[int]plugin.Validate

	Create(*RuleCtx) map[parser.NodeType]walk.ListenFn
}

type Rule struct {
	Proto     RuleFact
	Listeners map[parser.NodeType]*walk.Listener
}

type RuleCtx struct {
	unit     Unit
	ruleFact RuleFact
}

func (u *RuleCtx) Config() *Config {
	return u.unit.Config()
}

func (u *RuleCtx) Opts() []interface{} {
	c := u.Config().CfgOfRule(u.ruleFact.Name())
	if c == nil {
		return nil
	}
	return c.Opts
}

func (u *RuleCtx) Report(node parser.Node, msg string, level DiagLevel) {
	lang := u.unit.Lang()
	rule := u.ruleFact.Name()
	lvl := u.Config().LevelOfRule(rule, level)

	dig := &Diagnosis{
		Loc:   node.Loc().Clone(),
		Lang:  lang,
		Rule:  rule,
		Msg:   msg,
		Level: lvl,
	}
	u.unit.Report(dig)
}

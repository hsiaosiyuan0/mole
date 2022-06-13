package linter

import (
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

type Docs struct {
	Desc string
	Url  string
}

type Meta struct {
	Lang []string
	Kind Kind

	Docs Docs

	Options plugin.Options
}

type RuleFact interface {
	Name() string
	Meta() *Meta
	Create(Unit) map[parser.NodeType]walk.ListenFn
}

type Rule struct {
	Proto     RuleFact
	Listeners map[parser.NodeType]*walk.Listener
}

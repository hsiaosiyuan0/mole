package estree

import (
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type SrcLoc struct {
	Source string    `json:"source"`
	Start  *Position `json:"start"`
	End    *Position `json:"end"`
}

type Program struct {
	Type string      `json:"type"`
	Loc  *SrcLoc     `json:"loc"`
	Body interface{} `json:"body"`
}

func NewPosition(p *parser.Pos) *Position {
	return &Position{Line: p.Line(), Column: p.Column()}
}

func NewSrcLoc(s *parser.Loc) *SrcLoc {
	return &SrcLoc{
		Source: s.Source(),
		Start:  NewPosition(s.Begin()),
		End:    NewPosition(s.End()),
	}
}

func NewProgram(n *parser.Prog) *Program {
	return &Program{
		Type: "",
		Loc:  NewSrcLoc(n.Loc()),
		Body: n.Body(),
	}
}

package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	span "github.com/hsiaosiyuan0/mole/span"
	"github.com/hsiaosiyuan0/mole/util"
)

type Loc struct {
	Range span.Range
	Begin span.Pos
	End   span.Pos
}

func NewLocFromSpan(s *span.Source, start, end Node) *Loc {
	sr := start.Range()
	er := end.Range()

	return &Loc{
		Begin: s.OfstLineCol(sr.Lo),
		End:   s.OfstLineCol(er.Hi),
		Range: span.Range{
			Lo: sr.Lo,
			Hi: er.Hi,
		},
	}
}

func CalcLoc(node Node, s *span.Source) *Loc {
	if util.IsNilPtr(node) {
		return nil
	}
	return CalcLocInRng(node.Range(), s)
}

func CalcLocInRng(rng span.Range, s *span.Source) *Loc {
	begin, end := s.LineCol(rng)
	return &Loc{
		Begin: begin,
		End:   end,
		Range: rng,
	}
}

func (l *Loc) Before(b *Loc) bool {
	return l.Range.Hi < b.Range.Lo
}

func (l *Loc) Text(s *span.Source) string {
	return s.RngText(l.Range)
}

func NodeText(node Node, s *span.Source) string {
	ret := ""
	switch node.Type() {
	case N_NAME:
		id := node.(*Ident)
		ret = id.val
		if id.pvt {
			ret = "#" + ret
		}
	case N_LIT_STR:
		ret = node.(*StrLit).val
	case N_LIT_BOOL:
		if node.(*BoolLit).val {
			ret = "true"
		} else {
			ret = "false"
		}
	case N_LIT_NULL:
		ret = "null"
	}
	return s.RngText(node.Range())
}

func ParseBigint(t string) *big.Int {
	t = strings.ReplaceAll(t[:len(t)-1], "_", "")

	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		i := big.NewInt(0)
		i.SetString(t[2:], 16)
		return i
	}
	if strings.HasPrefix(t, "0o") || strings.HasPrefix(t, "0O") {
		i := big.NewInt(0)
		i.SetString(t[2:], 8)
		return i
	}
	if strings.HasPrefix(t, "0b") || strings.HasPrefix(t, "0B") {
		i := big.NewInt(0)
		i.SetString(t[2:], 2)
		return i
	}
	if strings.HasPrefix(t, "0") && len(t) > 1 {
		t = strings.TrimLeft(t, "0")
		i := big.NewInt(0)
		i.SetString(t, 8)
		return i
	}
	i := big.NewInt(0)
	i.SetString(t, 10)
	return i
}

func ParseFloat(t string) float64 {
	t = strings.ReplaceAll(t, "_", "")

	if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
		f, _ := strconv.ParseUint(t[2:], 16, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0o") || strings.HasPrefix(t, "0O") {
		f, _ := strconv.ParseUint(t[2:], 8, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0b") || strings.HasPrefix(t, "0B") {
		f, _ := strconv.ParseUint(t[2:], 2, 32)
		return float64(f)
	}
	if strings.HasPrefix(t, "0") && len(t) > 1 {
		t = strings.TrimLeft(t, "0")
		f, _ := strconv.ParseUint(t, 8, 32)
		return float64(f)
	}
	f, _ := strconv.ParseFloat(t, 64)
	return f
}

func NodeToFloat(node Node, s *span.Source) float64 {
	if node.Type() != N_LIT_NUM {
		return 0
	}

	t := NodeText(node, s)
	bi := t[len(t)-1] == 'n'

	if bi {
		i := ParseBigint(t)
		f := big.NewFloat(0)
		f.SetInt(i)
		max := big.NewFloat(0)
		max.SetUint64(math.MaxInt)
		c := f.Cmp(max)
		if c == -1 || c == 0 {
			ff, _ := f.Float64()
			return ff
		}
		return 0
	}
	return ParseFloat(t)
}

func FirstLoc(locs ...*Loc) *Loc {
	start := 0
	line := uint32(math.MaxUint32)
	col := uint32(math.MaxUint32)
	for i, loc := range locs {
		if loc == nil {
			continue
		}
		if loc.Begin.Line < line || (loc.Begin.Line == line && loc.Begin.Col < col) {
			line = loc.Begin.Line
			col = loc.Begin.Col
			start = i
		}
	}
	return locs[start]
}

func LastLoc(locs ...*Loc) *Loc {
	end := 0
	line := uint32(0)
	col := uint32(0)
	for i, loc := range locs {
		if loc == nil {
			continue
		}
		if loc.End.Line > line || (loc.End.Line == line && loc.End.Col > col) {
			line = loc.End.Line
			col = loc.End.Col
			end = i
		}
	}
	return locs[end]
}

func LocWithTypeInfo(node Node, includeParamProp bool, s *span.Source) *Loc {
	nw, ok := node.(NodeWithTypInfo)
	if !ok {
		return CalcLoc(node, s)
	}

	ti := nw.TypInfo()
	loc := CalcLoc(node, s)

	starLocList := []*Loc{CalcLoc(ti.TypParams(), s), CalcLoc(node, s)}
	if includeParamProp {
		starLocList = append(starLocList, []*Loc{CalcLocInRng(ti.BeginRng(), s)}...)
	}

	start := FirstLoc(starLocList...)
	loc.Begin = start.Begin
	loc.Range.Lo = start.Range.Lo

	end := LastLoc(CalcLoc(ti.TypAnnot(), s), CalcLocInRng(ti.Ques(), s), CalcLoc(node, s))
	loc.End = end.End
	loc.Range.Hi = end.Range.Hi

	return loc
}

func TplLocWithTag(n *TplExpr, s *span.Source) *Loc {
	loc := CalcLoc(n, s)
	if n.tag != nil {
		tl := CalcLoc(n.tag, s)
		loc.Begin = tl.Begin
		loc.Range.Lo = tl.Range.Lo
	}
	return loc
}

func TokText(t *Token, s *span.Source) string {
	if t.text != "" {
		return t.text
	}

	if t.IsKw() {
		return TokenKinds[t.value].Name
	}

	if !t.txt.Empty() {
		t.text = s.RngText(t.txt)
	} else {
		t.text = s.RngText(t.rng)
	}

	return t.text
}

func TokRawText(t *Token, s *span.Source) string {
	return s.RngText(t.rng)
}

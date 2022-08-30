package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	span "github.com/hsiaosiyuan0/mole/span"
)

func NodeText(node Node, s *span.Source) string {
	ret := ""
	switch node.Type() {
	case N_NAME:
		return node.(*Ident).val
	case N_LIT_STR:
		return node.(*StrLit).val
	case N_LIT_BOOL:
		if node.(*BoolLit).val {
			ret = "true"
		} else {
			ret = "false"
		}
		return ret
	case N_LIT_NULL:
		return "null"
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

func NodeIsBigint(node Node, s *span.Source) bool {
	t := NodeText(node, s)
	return t[len(t)-1] == 'n'
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

func TokText(t *Token, s *span.Source) string {
	if t.text != "" {
		return t.text
	}

	if t.IsKw() {
		return TokenKinds[t.value].Name
	}

	if t.value == T_EOF {
		return ""
	}

	if !t.txt.Empty() {
		t.text = s.RngText(t.txt)
	} else if t.rng.Valid() {
		t.text = s.RngText(t.rng)
	}
	return t.text
}

func TokRawText(t *Token, s *span.Source) string {
	return s.RngText(t.rng)
}

package span

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/util"
)

func TestEOL(t *testing.T) {
	s := NewSource("", "\u2028\u2029\u000a\u000d\u000a")
	AssertEqual(t, '\u2028', s.Read(), "\\u2028 should be EOL")
	AssertEqual(t, '\u2029', s.Read(), "\\u2029 should be EOL")
	AssertEqual(t, '\u000a', s.Read(), "\\u000a should be EOL")
	AssertEqual(t, '\u000a', s.Read(), "\\u000d\\u000a should be EOL")
}

func TestAhead(t *testing.T) {
	s := NewSource("", "hello world")
	AssertEqual(t, true, s.AheadIsCh('h'), "ahead should be h")
	AssertEqual(t, 'h', s.Read(), "next should be h")

	AssertEqual(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	AssertEqual(t, 'e', s.Peek(), "next should be e")
	AssertEqual(t, 'e', s.Read(), "next should be e")
	AssertEqual(t, 'l', s.Read(), "next should be l")

	AssertEqual(t, 'l', s.Read(), "next should be l")
	AssertEqual(t, 'o', s.Read(), "next should be o")
	AssertEqual(t, ' ', s.Read(), "next should be space")

	AssertEqual(t, true, s.AheadIsChOr('1', 'w'), "ahead maybe w")
	AssertEqual(t, 'w', s.Read(), "next should be w")
}

func TestPos(t *testing.T) {
	s := NewSource("", "hello world")
	AssertEqual(t, 'h', s.Read(), "next should be h")
	AssertEqual(t, uint32(1), s.Ofst(), "pos should be 1")

	AssertEqual(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	AssertEqual(t, uint32(1), s.Ofst(), "pos should be 1 after lookahead")

	AssertEqual(t, 'e', s.Read(), "next should be e")
	AssertEqual(t, 'l', s.Read(), "next should be l")
	AssertEqual(t, uint32(3), s.Ofst(), "pos should be 3")
}

func TestLineCol(t *testing.T) {
	s := NewSource("", "\"str1\"\r\n\"str2\"\n\n\"str3\"\n\n  \"str4\"\n\n\n\"multiline\n\n str\"\n")

	from, to := s.LineCol(Range{0, 6})
	AssertEqual(t, true, from.Line == 1 && from.Col == 0, "next should be h")
	AssertEqual(t, true, to.Line == 1 && to.Col == 6, "next should be h")

	from, to = s.LineCol(Range{8, 14})
	AssertEqual(t, true, from.Line == 2 && from.Col == 0, "next should be h")
	AssertEqual(t, true, to.Line == 2 && to.Col == 6, "next should be h")

	from, to = s.LineCol(Range{16, 22})
	AssertEqual(t, true, from.Line == 4 && from.Col == 0, "next should be h")
	AssertEqual(t, true, to.Line == 4 && to.Col == 6, "next should be h")

	from, to = s.LineCol(Range{26, 32})
	AssertEqual(t, true, from.Line == 6 && from.Col == 2, "next should be h")
	AssertEqual(t, true, to.Line == 6 && to.Col == 8, "next should be h")

	from, to = s.LineCol(Range{35, 52})
	AssertEqual(t, true, from.Line == 9 && from.Col == 0, "next should be h")
	AssertEqual(t, true, to.Line == 11 && to.Col == 5, "next should be h")
}

package parser

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/internal"
)

func TestEOL(t *testing.T) {
	s := NewSource("", "\u2028\u2029\u000a\u000d\u000a")
	AssertEqual(t, 1, s.line, "line should begin from 1")

	AssertEqual(t, '\u2028', s.Read(), "\\u2028 should be EOL")
	AssertEqual(t, 1, s.line, "\\u2028 should step line")

	AssertEqual(t, '\u2029', s.Read(), "\\u2029 should be EOL")
	AssertEqual(t, 1, s.line, "\\u2028 should step line")

	AssertEqual(t, EOL, s.Read(), "\\u000a should be EOL")
	AssertEqual(t, 2, s.line, "\\u000a should step line")

	AssertEqual(t, EOL, s.Read(), "\\u000d\\u000a should be EOL")
	AssertEqual(t, 3, s.line, "\\u000d should step line")
}

func TestAhead(t *testing.T) {
	s := NewSource("", "hello world")
	AssertEqual(t, true, s.AheadIsCh('h'), "ahead should be h")
	AssertEqual(t, 'h', s.Read(), "next should be h")

	AssertEqual(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	AssertEqual(t, 'e', s.Peek(), "next should be e")
	AssertEqual(t, 'e', s.Read(), "next should be e")
	AssertEqual(t, 'l', s.Read(), "next should be l")
	AssertEqual(t, 0, s.pl, "peek buf should be empty")

	AssertEqual(t, 'l', s.Read(), "next should be l")
	AssertEqual(t, 'o', s.Read(), "next should be o")
	AssertEqual(t, ' ', s.Read(), "next should be space")

	AssertEqual(t, true, s.AheadIsChOr('1', 'w'), "ahead maybe w")
	AssertEqual(t, 'w', s.Read(), "next should be w")
}

func TestPos(t *testing.T) {
	s := NewSource("", "hello world")
	AssertEqual(t, 'h', s.Read(), "next should be h")
	AssertEqual(t, 1, s.Ofst(), "pos should be 1")

	AssertEqual(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	AssertEqual(t, 1, s.Ofst(), "pos should be 1 after lookahead")

	AssertEqual(t, 'e', s.Read(), "next should be e")
	AssertEqual(t, 'l', s.Read(), "next should be l")
	AssertEqual(t, 3, s.Ofst(), "pos should be 3")
}

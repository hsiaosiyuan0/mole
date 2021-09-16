package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestEOL(t *testing.T) {
	s := NewSource("", "\u2028\u000a\u000d\u000a")
	assert.Equal(t, 1, s.line, "line should begin from 1")

	assert.Equal(t, EOL, s.Read(), "\\u2028 should be EOL")
	assert.Equal(t, 2, s.line, "\\u2028 should step line")

	assert.Equal(t, EOL, s.Read(), "\\u000a should be EOL")
	assert.Equal(t, 3, s.line, "\\u000a should step line")

	assert.Equal(t, EOL, s.Read(), "\\u000d\\u000a should be EOL")
	assert.Equal(t, 4, s.line, "\\u000d should step line")
}

func TestAhead(t *testing.T) {
	s := NewSource("", "hello world")
	assert.Equal(t, true, s.AheadIsCh('h'), "ahead should be h")
	assert.Equal(t, 'h', s.Read(), "next should be h")

	assert.Equal(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	assert.Equal(t, 'e', s.Peek(), "next should be e")
	assert.Equal(t, 'e', s.Read(), "next should be e")
	assert.Equal(t, 'l', s.Read(), "next should be l")
	assert.Equal(t, 0, s.peekedLen, "peek buf should be empty")

	assert.Equal(t, 'l', s.Read(), "next should be l")
	assert.Equal(t, 'o', s.Read(), "next should be o")
	assert.Equal(t, ' ', s.Read(), "next should be space")

	assert.Equal(t, true, s.AheadIsChOr('1', 'w'), "ahead maybe w")
	assert.Equal(t, 'w', s.Read(), "next should be w")
}

func TestPos(t *testing.T) {
	s := NewSource("", "hello world")
	assert.Equal(t, 'h', s.Read(), "next should be h")
	assert.Equal(t, 1, s.Ofst(), "pos should be 1")

	assert.Equal(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	assert.Equal(t, 1, s.Ofst(), "pos should be 1 after lookahead")

	assert.Equal(t, 'e', s.Read(), "next should be e")
	assert.Equal(t, 'l', s.Read(), "next should be l")
	assert.Equal(t, 3, s.Ofst(), "pos should be 3")
}

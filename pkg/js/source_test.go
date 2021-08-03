package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mlint/pkg/assert"
)

func TestEOL(t *testing.T) {
	s := NewSource("", "\u2028\u000a\u000d\u000a")
	assert.Equal(t, 1, s.line, "line should begin from 1")

	assert.Equal(t, EOL, s.Next(), "\\u2028 should be EOL")
	assert.Equal(t, 2, s.line, "\\u2028 should step line")

	assert.Equal(t, EOL, s.Next(), "\\u000a should be EOL")
	assert.Equal(t, 3, s.line, "\\u000a should step line")

	assert.Equal(t, EOL, s.Next(), "\\u000d\\u000a should be EOL")
	assert.Equal(t, 4, s.line, "\\u000d should step line")
}

func TestAhead(t *testing.T) {
	s := NewSource("", "hello world")
	assert.Equal(t, true, s.AheadIsCh('h'), "ahead should be h")
	assert.Equal(t, 'h', s.Next(), "next should be h")

	assert.Equal(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	assert.Equal(t, 'e', s.Next(), "next should be e")
	assert.Equal(t, 'l', s.Next(), "next should be l")
	assert.Equal(t, 0, len(s.peeked), "peek buf should be empty")

	assert.Equal(t, 'l', s.Next(), "next should be l")
	assert.Equal(t, 'o', s.Next(), "next should be o")
	assert.Equal(t, ' ', s.Next(), "next should be space")

	assert.Equal(t, true, s.AheadIsChOr('1', 'w'), "ahead maybe w")
	assert.Equal(t, 'w', s.Next(), "next should be w")
}

func TestPos(t *testing.T) {
	s := NewSource("", "hello world")
	assert.Equal(t, 'h', s.Next(), "next should be h")
	assert.Equal(t, 1, s.Pos(), "pos should be 1")

	assert.Equal(t, true, s.AheadIsChs2('e', 'l'), "ahead should be el")
	assert.Equal(t, 1, s.Pos(), "pos should be 1 after lookahead")

	assert.Equal(t, 'e', s.Next(), "next should be e")
	assert.Equal(t, 'l', s.Next(), "next should be l")
	assert.Equal(t, 3, s.Pos(), "pos should be 3")
}

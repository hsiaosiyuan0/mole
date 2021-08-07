package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mlint/pkg/assert"
)

func TestReadName(t *testing.T) {
	s := NewSource("", "\\u0074 t\\u0065st")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok t")
	assert.Equal(t, "t", tok.text, "should be t")
	assert.Equal(t, 6, tok.loc.hi, "should be t")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok test")
	assert.Equal(t, "test", tok.text, "should be test")
}

func TestReadNum(t *testing.T) {
	s := NewSource("", "1 23 1e1 .1e1 .1_1 1n 0b01 0B01 0o01 0O01 0x01 0X01 0x0_1")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 1")
	assert.Equal(t, "1", tok.Text(), "should be 1")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 23")
	assert.Equal(t, "23", tok.Text(), "should be 23")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 1e1")
	assert.Equal(t, "1e1", tok.Text(), "should be 1e1")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok .1e1")
	assert.Equal(t, ".1e1", tok.Text(), "should be .1e1")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok .1_1")
	assert.Equal(t, ".1_1", tok.Text(), "should be .1_1")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 1n")
	assert.Equal(t, "1n", tok.Text(), "should be 1n")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0b01")
	assert.Equal(t, "0b01", tok.Text(), "should be 0b01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0B01")
	assert.Equal(t, "0B01", tok.Text(), "should be 0B01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0o01")
	assert.Equal(t, "0o01", tok.Text(), "should be 0o01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0O01")
	assert.Equal(t, "0O01", tok.Text(), "should be 0O01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0x01")
	assert.Equal(t, "0x01", tok.Text(), "should be 0x01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0X01")
	assert.Equal(t, "0X01", tok.Text(), "should be 0X01")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok 0x0_1")
	assert.Equal(t, "0x0_1", tok.Text(), "should be 0x0_1")
}

func TestReadStr(t *testing.T) {
	s := NewSource("", `
  'h'
  'a\nb'
  't\u0065st'
  `)

	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok h")
	assert.Equal(t, "h", tok.Text(), "should be h")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok a\\nb")
	assert.Equal(t, "a\nb", tok.Text(), "should be a\\nb")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok test")
	assert.Equal(t, "test", tok.Text(), "should be test")
}

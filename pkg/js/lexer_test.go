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
  '\012'
  '\0012'
  '\251'
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

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok \\n")
	assert.Equal(t, "\n", tok.Text(), "should be \\n")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok \\u00012")
	assert.Equal(t, "\u00012", tok.Text(), "should be \\u00012")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok ©")
	assert.Equal(t, "©", tok.Text(), "should be ©")
}

func TestReadId(t *testing.T) {
	s := NewSource("", "if with")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok if")
	assert.Equal(t, "if", tok.Text(), "should be if")
	assert.Equal(t, T_IF, tok.value, "should be tok if")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok with")
	assert.Equal(t, "with", tok.Text(), "should be with")
	assert.Equal(t, T_WITH, tok.value, "should be tok with")
}

func TestReadSymbol(t *testing.T) {
	s := NewSource("", `
  { }
  ( )
  [ ]
  ... .
  ; , ? :
  ++ --
  ?. =>
  ?? < > <= >= == != === !==
  << >> >>> | ^ & || &&
  + - * / % ** ! ~
  = += -= ??= ||= &&= |= ^= &= <<= >>= >>>=
  *= /= %= **=
  `)
	l := NewLexer(s)
	assert.Equal(t, T_BRACE_L, l.Next().value, "should be tok {")
	assert.Equal(t, T_BRACE_R, l.Next().value, "should be tok }")
	assert.Equal(t, T_PAREN_L, l.Next().value, "should be tok (")
	assert.Equal(t, T_PAREN_R, l.Next().value, "should be tok )")
	assert.Equal(t, T_BRACKET_L, l.Next().value, "should be tok [")
	assert.Equal(t, T_BRACKET_R, l.Next().value, "should be tok ]")
	assert.Equal(t, T_DOT_TRI, l.Next().value, "should be tok ...")
	assert.Equal(t, T_DOT, l.Next().value, "should be tok .")
	assert.Equal(t, T_SEMI, l.Next().value, "should be tok ;")
	assert.Equal(t, T_COMMA, l.Next().value, "should be tok ,")
	assert.Equal(t, T_HOOK, l.Next().value, "should be tok ?")
	assert.Equal(t, T_COLON, l.Next().value, "should be tok :")
	assert.Equal(t, T_INC, l.Next().value, "should be tok ++")
	assert.Equal(t, T_DEC, l.Next().value, "should be tok --")
	assert.Equal(t, T_OPT_CHAIN, l.Next().value, "should be tok ?.")
	assert.Equal(t, T_ARROW, l.Next().value, "should be tok =>")
	assert.Equal(t, T_NULLISH, l.Next().value, "should be tok ??")
	assert.Equal(t, T_LT, l.Next().value, "should be tok <")
	assert.Equal(t, T_GT, l.Next().value, "should be tok >")
	assert.Equal(t, T_LE, l.Next().value, "should be tok <=")
	assert.Equal(t, T_GE, l.Next().value, "should be tok >=")
	assert.Equal(t, T_EQ, l.Next().value, "should be tok ==")
	assert.Equal(t, T_NE, l.Next().value, "should be tok !=")
	assert.Equal(t, T_EQ_S, l.Next().value, "should be tok ===")
	assert.Equal(t, T_NE_S, l.Next().value, "should be tok !==")
	assert.Equal(t, T_LSH, l.Next().value, "should be tok <<")
	assert.Equal(t, T_RSH, l.Next().value, "should be tok >>")
	assert.Equal(t, T_RSH_U, l.Next().value, "should be tok >>>")
	assert.Equal(t, T_BIT_OR, l.Next().value, "should be tok |")
	assert.Equal(t, T_BIT_XOR, l.Next().value, "should be tok ^")
	assert.Equal(t, T_BIT_AND, l.Next().value, "should be tok &")
	assert.Equal(t, T_OR, l.Next().value, "should be tok ||")
	assert.Equal(t, T_AND, l.Next().value, "should be tok &&")
	assert.Equal(t, T_ADD, l.Next().value, "should be tok +")
	assert.Equal(t, T_SUB, l.Next().value, "should be tok -")
	assert.Equal(t, T_MUL, l.Next().value, "should be tok *")
	assert.Equal(t, T_DIV, l.Next().value, "should be tok /")
	assert.Equal(t, T_MOD, l.Next().value, "should be tok %")
	assert.Equal(t, T_POW, l.Next().value, "should be tok **")
	assert.Equal(t, T_NOT, l.Next().value, "should be tok !")
	assert.Equal(t, T_BIT_NOT, l.Next().value, "should be tok ~")
	assert.Equal(t, T_ASSIGN, l.Next().value, "should be tok =")
	assert.Equal(t, T_ASSIGN_ADD, l.Next().value, "should be tok +=")
	assert.Equal(t, T_ASSIGN_SUB, l.Next().value, "should be tok -=")
	assert.Equal(t, T_ASSIGN_NULLISH, l.Next().value, "should be tok ??=")
	assert.Equal(t, T_ASSIGN_OR, l.Next().value, "should be tok ||=")
	assert.Equal(t, T_ASSIGN_AND, l.Next().value, "should be tok &&=")
	assert.Equal(t, T_ASSIGN_BIT_OR, l.Next().value, "should be tok |=")
	assert.Equal(t, T_ASSIGN_BIT_XOR, l.Next().value, "should be tok ^=")
	assert.Equal(t, T_ASSIGN_BIT_AND, l.Next().value, "should be tok &=")
	assert.Equal(t, T_ASSIGN_BIT_LSH, l.Next().value, "should be tok <==")
	assert.Equal(t, T_ASSIGN_BIT_RSH, l.Next().value, "should be tok >==")
	assert.Equal(t, T_ASSIGN_BIT_RSH_U, l.Next().value, "should be tok >>>=")
	assert.Equal(t, T_ASSIGN_MUL, l.Next().value, "should be tok *=")
	assert.Equal(t, T_ASSIGN_DIV, l.Next().value, "should be tok /=")
	assert.Equal(t, T_ASSIGN_MOD, l.Next().value, "should be tok %=")
	assert.Equal(t, T_ASSIGN_POW, l.Next().value, "should be tok **=")
}

func TestReadAsyncAwait(t *testing.T) {
	s := NewSource("", `async function(a = {}) {async await}
  async await
  `)
	l := NewLexer(s)

	assert.Equal(t, T_ASYNC, l.Next().value, "should be tok async")
	assert.Equal(t, T_FUNC, l.Next().value, "should be tok function")

	assert.Equal(t, T_PAREN_L, l.Next().value, "should be tok (")
	assert.Equal(t, T_NAME, l.Next().value, "should be tok name")
	assert.Equal(t, T_ASSIGN, l.Next().value, "should be tok =")
	assert.Equal(t, T_BRACE_L, l.Next().value, "should be tok {")
	assert.Equal(t, T_BRACE_R, l.Next().value, "should be tok }")
	assert.Equal(t, T_PAREN_R, l.Next().value, "should be tok )")

	assert.Equal(t, T_BRACE_L, l.Next().value, "should be tok {")
	assert.Equal(t, T_NAME, l.Next().value, "should be tok async:name")
	assert.Equal(t, T_AWAIT, l.Next().value, "should be tok await")
	assert.Equal(t, T_BRACE_R, l.Next().value, "should be tok }")

	l.popMode()
	assert.Equal(t, T_NAME, l.Next().value, "should be tok name")
	assert.Equal(t, T_NAME, l.Next().value, "should be tok name")
}

func TestReadRegexp(t *testing.T) {
	s := NewSource("", `
  /a/i
  a / /a/i
  `)
	_ = NewLexer(s)

	// assert.Equal(t, T_REGEXP, l.Next().value, "should be tok regexp")
	// assert.Equal(t, T_NAME, l.Next().value, "should be tok a")
	// assert.Equal(t, T_DIV, l.Next().value, "should be tok div")
	// assert.Equal(t, T_REGEXP, l.Next().value, "should be tok regexp")
}

package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestReadName(t *testing.T) {
	s := NewSource("", "\\u0074 t\\u0065st")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok t")
	assert.Equal(t, "t", tok.text, "should be t")
	assert.Equal(t, 6, tok.raw.hi, "should be t")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok test")
	assert.Equal(t, "test", tok.text, "should be test")
}

func TestReadId(t *testing.T) {
	s := NewSource("", "if with void")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok if")
	assert.Equal(t, "if", tok.Text(), "should be if")
	assert.Equal(t, T_IF, tok.value, "should be tok if")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok with")
	assert.Equal(t, "with", tok.Text(), "should be with")
	assert.Equal(t, T_WITH, tok.value, "should be tok with")

	tok = l.Next()
	assert.Equal(t, true, tok.IsLegal(), "should be ok void")
	assert.Equal(t, "void", tok.Text(), "should be void")
	assert.Equal(t, T_VOID, tok.value, "should be tok void")
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
  + - *
  } a / % ** ! ~
  = += -= ??= ||= &&= |= ^= &= <<= >>= >>>=
  *= a /= %= **=
  `)
	l := NewLexer(s)
	l.ver = ES12
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
	assert.Equal(t, T_LTE, l.Next().value, "should be tok <=")
	assert.Equal(t, T_GTE, l.Next().value, "should be tok >=")
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
	assert.Equal(t, T_BRACE_R, l.Next().value, "should be tok }")
	assert.Equal(t, T_NAME, l.Next().value, "should be tok name")
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
	assert.Equal(t, T_NAME, l.Next().value, "should be tok name")
	assert.Equal(t, T_ASSIGN_DIV, l.Next().value, "should be tok /=")
	assert.Equal(t, T_ASSIGN_MOD, l.Next().value, "should be tok %=")
	assert.Equal(t, T_ASSIGN_POW, l.Next().value, "should be tok **=")
}

func TestReadRegexp(t *testing.T) {
	s := NewSource("", `
  /a/ig
  a / /b/i
  `)
	l := NewLexer(s)

	tok := l.Next()
	assert.Equal(t, T_REGEXP, tok.value, "should be tok regexp /a/ig")
	assert.Equal(t, "a", tok.ext.(*TokExtRegexp).pattern.Text(), "should be tok regexp pattern /a/ig")
	assert.Equal(t, "ig", tok.ext.(*TokExtRegexp).flags.Text(), "should be tok regexp flags /a/ig")

	assert.Equal(t, T_NAME, l.Next().value, "should be tok a")
	assert.Equal(t, T_DIV, l.Next().value, "should be tok div")

	tok = l.Next()
	assert.Equal(t, T_REGEXP, tok.value, "should be tok regexp /b/i")
	assert.Equal(t, "b", tok.ext.(*TokExtRegexp).pattern.Text(), "should be tok regexp pattern /b/i")
	assert.Equal(t, "i", tok.ext.(*TokExtRegexp).flags.Text(), "should be tok regexp flags /b/i")
}
func TestReadTpl(t *testing.T) {
	s := NewSource("", "`abc`"+"`a${ {} }b${c}d`")
	l := NewLexer(s)

	tok := l.Next()
	assert.Equal(t, T_TPL_HEAD, tok.value, "should be tok str")
	assert.Equal(t, true, tok.ext.(*TokExtTplSpan).Plain, "should be tok str")
	assert.Equal(t, "abc", tok.ext.(*TokExtTplSpan).str, "should be tok str abc")

	tok = l.Next()
	assert.Equal(t, T_TPL_HEAD, tok.value, "should be tok tpl head")
	assert.Equal(t, "a", tok.ext.(*TokExtTplSpan).str, "should be tok tpl a")

	tok = l.Next()
	assert.Equal(t, T_BRACE_L, tok.value, "should be tok {")
	tok = l.Next()
	assert.Equal(t, T_BRACE_R, tok.value, "should be tok }")

	tok = l.Next()
	assert.Equal(t, T_TPL_SPAN, tok.value, "should be tok tpl span")
	assert.Equal(t, "b", tok.ext.(*TokExtTplSpan).str, "should be tok tpl b")

	tok = l.Next()
	assert.Equal(t, T_NAME, tok.value, "should be tok c")
	assert.Equal(t, "c", tok.Text(), "should be tok tpl c")

	tok = l.Next()
	assert.Equal(t, T_TPL_TAIL, tok.value, "should be tok tpl tail")
	assert.Equal(t, "d", tok.ext.(*TokExtTplSpan).str, "should be tok tpl d")
}

func TestReadNestTpl(t *testing.T) {
	s := NewSource("", "`a${ 1 + {{`c${d}e`}} }b`")
	l := NewLexer(s)

	tok := l.Next()
	assert.Equal(t, T_TPL_HEAD, tok.value, "should be tok tpl head")
	assert.Equal(t, "a", tok.ext.(*TokExtTplSpan).str, "should be tok tpl a")

	assert.Equal(t, T_NUM, l.Next().value, "should be tok 1")
	assert.Equal(t, T_ADD, l.Next().value, "should be tok +")

	assert.Equal(t, T_BRACE_L, l.Next().value, "should be tok {")
	assert.Equal(t, T_BRACE_L, l.Next().value, "should be tok {")

	tok = l.Next()
	assert.Equal(t, T_TPL_HEAD, tok.value, "should be tok tpl head c")
	assert.Equal(t, "c", tok.ext.(*TokExtTplSpan).str, "should be tok tpl c")

	assert.Equal(t, T_NAME, l.Next().value, "should be tok d")

	tok = l.Next()
	assert.Equal(t, T_TPL_TAIL, tok.value, "should be tok tpl tail e")

	tok = l.Next()
	assert.Equal(t, T_BRACE_R, tok.value, "should be tok }")

	tok = l.Next()
	assert.Equal(t, T_BRACE_R, tok.value, "should be tok }")

	tok = l.Next()
	assert.Equal(t, T_TPL_TAIL, tok.value, "should be tok tpl tail")
	assert.Equal(t, "b", tok.ext.(*TokExtTplSpan).str, "should be tok tpl b")
}

func TestReadTplOctalEscape(t *testing.T) {
	s := NewSource("", "`\\1`")
	l := NewLexer(s)
	tok := l.Next()
	assert.Equal(t, T_TPL_HEAD, tok.value, "should be tpl head")
	assert.Equal(t, true, tok.ext.(*TokExtTplSpan).IllegalEscape != nil, "should be tpl head")
	assert.Equal(t, 1, len(l.mode), "mode should be balanced")
}

func TestReadComment(t *testing.T) {
	s := NewSource("", `
  //
  `)
	l := NewLexer(s)
	l.Next()
	assert.Equal(t, "//", l.lastComment().Text(), "should be tok comment //")

	s = NewSource("", `
  // comment1
  `)
	l = NewLexer(s)
	l.Next()
	assert.Equal(t, "// comment1", l.lastComment().Text(), "should be tok // comment1")

	s = NewSource("", `
  /**/
  `)
	l = NewLexer(s)
	l.Next()
	assert.Equal(t, "/**/", l.lastComment().Text(), "should be tok /**/")

	s = NewSource("", `
  /* comment2 */
  `)
	l = NewLexer(s)
	l.Next()
	assert.Equal(t, "/* comment2 */", l.lastComment().Text(), "should be tok /* comment2 */")

	s = NewSource("", `/**

  comment 3
  **/
  `)
	l = NewLexer(s)
	l.Next()
	assert.Equal(t, `/**

  comment 3
  **/`, l.lastComment().Text(), "should be tok comment3")
}

func TestAfterLineTerminator(t *testing.T) {
	s := NewSource("", "a\n1")
	l := NewLexer(s)

	tok := l.Next()
	assert.Equal(t, T_NAME, tok.value, "should be tok a")
	assert.Equal(t, false, tok.afterLineTerminator, "mode should be afterLineTerminator false")

	tok = l.Next()
	assert.Equal(t, T_NUM, tok.value, "should be tok 1")
	assert.Equal(t, true, tok.afterLineTerminator, "mode should be afterLineTerminator true")
}

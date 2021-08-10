package js

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type LexerModeValue int

const (
	LM_NONE       LexerModeValue = 0
	LM_STRICT                    = 1
	LM_TEMPLATE                  = 1 << 2
	LM_IMPORT                    = 1 << 3
	LM_NAMED_LIST                = 1 << 4
	LM_ASYNC                     = 1 << 5
	LM_GENERATOR                 = 1 << 6
	LM_CLASS_BODY                = 1 << 7
	LM_CLASS_CTOR                = 1 << 8
	LM_FOR_OF                    = 1 << 9
	LM_NEW                       = 1 << 10
)

type LexerMode struct {
	value   LexerModeValue
	paren   int
	brace   int
	bracket int
}

const sizeOfPeekedTok = 5

type Lexer struct {
	*Source

	mode []*LexerMode

	peeked    [sizeOfPeekedTok]*Token
	peekedLen int
	peekedR   int
	peekedW   int

	prev *Token
}

func NewLexer(src *Source) *Lexer {
	lexer := &Lexer{Source: src, mode: make([]*LexerMode, 0)}
	lexer.mode = append(lexer.mode, &LexerMode{LM_NONE, 0, 0, 0})
	return lexer
}

func (l *Lexer) extMode(mode LexerModeValue, inherit bool) {
	if inherit {
		// only inherit the inheritable modes
		v := LM_NONE
		v |= l.curMode().value & LM_ASYNC
		v |= l.curMode().value & LM_STRICT
		mode |= v
	}
	l.mode = append(l.mode, &LexerMode{mode, 0, 0, 0})
}

func (l *Lexer) pushMode(mode LexerModeValue) {
	l.extMode(mode, true)
}

func (l *Lexer) popMode() *LexerMode {
	mLen := len(l.mode)
	if mLen == 1 {
		return l.mode[0]
	}
	m, last := l.mode[:mLen-1], l.mode[mLen-1]
	l.mode = m
	return last
}

func (l *Lexer) curMode() *LexerMode {
	return l.mode[len(l.mode)-1]
}

func (l *Lexer) isMode(mode LexerModeValue) bool {
	return l.curMode().value&mode > 0
}

func (l *Lexer) readTok() *Token {
	l.SkipSpace()
	if l.aheadIsIdStart() {
		return l.ReadName()
	} else if l.aheadIsNumStart() {
		return l.ReadNum()
	} else if l.aheadIsStrStart() {
		return l.ReadStr()
	} else if l.aheadIsTplStart() {
		return l.ReadTplSpan()
	}
	return l.ReadSymbol()
}

func (l *Lexer) PeekTok() *Token {
	if l.peekedLen == sizeOfPeekedTok {
		panic(l.error(fmt.Sprintf("peek buffer of lexer is full, max len is %d\n", l.peekedLen)))
	}

	tok := l.readTok()
	if tok.isEof() {
		return tok
	}

	l.peeked[l.peekedW] = tok
	l.peekedW += 1
	l.peekedLen += 1
	if l.peekedW == sizeOfPeekedTok {
		l.peekedW = 0
	}
	return tok
}

func (l *Lexer) nextTok() *Token {
	if l.peekedLen > 0 {
		tok := l.peeked[l.peekedR]
		l.peekedR += 1
		l.peekedLen -= 1
		if l.peekedR == sizeOfPeekedTok {
			l.peekedR = 0
		}
		return tok
	}
	return l.readTok()
}

func (l *Lexer) Next() *Token {
	tok := l.nextTok()
	if tok.value == T_NAME && tok.Text() == "async" && l.PeekTok().value == T_FUNC {
		tok.value = T_ASYNC
		tok.text = ""
		l.pushMode(LM_ASYNC)
	}
	l.prev = tok
	return tok
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-Template
func (l *Lexer) ReadTplSpan() *Token {
	l.Read() // consume `\`` or `}`

	// for inline spans can tell they are in template string
	l.pushMode(LM_TEMPLATE)
	defer l.popMode()

	tok := l.newToken()
	text, fin := l.readTplChs()
	if text == nil {
		return l.errToken(tok)
	}
	tok.text = string(text)
	if fin {
		l.popMode()
		return l.finToken(tok, T_TPL_TAIL)
	}
	return l.finToken(tok, T_TPL_SPAN)
}

func (l *Lexer) readTplChs() (text []rune, fin bool) {
	text = make([]rune, 0, 10)
	for {
		c := l.Peek()
		if c == '$' {
			l.Read()
			if l.AheadIsCh('{') {
				l.Read()
				l.pushMode(LM_TEMPLATE)
				break
			}
			text = append(text, c)
		} else if c == '\\' {
			l.Read()
			nc := l.Peek()
			if IsLineTerminator(nc) {
				l.readLineTerminator() // LineContinuation
			} else {
				r := l.readEscapeSeq()
				if r == utf8.RuneError || r == EOF {
					text = nil
					return
				}
				text = append(text, r)
			}
		} else if c == utf8.RuneError {
			text = nil
			return
		} else if c == '`' {
			l.Read()
			fin = true
			break
		} else {
			text = append(text, l.Read())
		}
	}
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-IdentifierName
func (l *Lexer) ReadName() *Token {
	tok := l.newToken()

	runes := make([]rune, 0, 10)
	r := l.readIdStart()
	if r == utf8.RuneError {
		return l.errToken(tok)
	}
	runes = append(runes, r)

	idPart, ok := l.readIdPart()
	if !ok {
		return l.errToken(tok)
	}
	runes = append(runes, idPart...)
	text := string(runes)

	if IsKeyword(text) {
		return l.finToken(tok, Keywords[text])
	} else if l.isMode(LM_STRICT) && IsStrictKeywords(text) {
		return l.finToken(tok, StrictKeywords[text])
	} else if l.isMode(LM_IMPORT) && text == "meta" {
		return l.finToken(tok, CtxKeywords[text])
	} else if l.isMode(LM_ASYNC) && text == "await" {
		return l.finToken(tok, T_AWAIT)
	} else if l.isMode(LM_FOR_OF) && text == "of" {
		return l.finToken(tok, T_OF)
	} else if l.isMode(LM_CLASS_BODY) && (text == "set" || text == "get") {
		return l.finToken(tok, CtxKeywords[text])
	} else if l.isMode(LM_NAMED_LIST) && (text == "as" || text == "from") {
		return l.finToken(tok, CtxKeywords[text])
	}
	tok.text = text
	return l.finToken(tok, T_NAME)
}

func (l *Lexer) ReadSymbol() *Token {
	c := l.Read()
	tok := l.newToken()
	val := tok.value
	switch c {
	case '{':
		val = T_BRACE_L
		if l.isMode(LM_ASYNC) && l.curMode().paren == 0 {
			// this branch means the brace_l of the function body is met,
			// skip push mode here to balance the pop of the brace_r of the
			// function body
		} else {
			l.pushMode(LM_NONE)
		}
	case '}':
		val = T_BRACE_R
		l.popMode()
	case '(':
		val = T_PAREN_L
		l.curMode().paren += 1
	case ')':
		val = T_PAREN_R
		l.curMode().paren -= 1
	case '[':
		val = T_BRACKET_L
	case ']':
		val = T_BRACKET_R
	case ';':
		val = T_SEMI
	case ',':
		val = T_COMMA
	case ':':
		val = T_COLON
	case '.':
		if l.AheadIsChs2('.', '.') {
			l.Read()
			l.Read()
			val = T_DOT_TRI
		} else {
			val = T_DOT
		}
	case '?':
		if l.AheadIsCh('.') {
			l.Read()
			val = T_OPT_CHAIN
		} else if l.AheadIsCh('?') {
			l.Read()
			if l.AheadIsCh('=') {
				l.Read()
				val = T_ASSIGN_NULLISH
			} else {
				val = T_NULLISH
			}
		} else {
			val = T_HOOK
		}
	case '+':
		if l.AheadIsCh('+') {
			l.Read()
			val = T_INC
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_ADD
		} else {
			val = T_ADD
		}
	case '-':
		if l.AheadIsCh('-') {
			l.Read()
			val = T_DEC
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_SUB
		} else {
			val = T_SUB
		}
	case '=':
		if l.AheadIsCh('>') {
			l.Read()
			val = T_ARROW
		} else if l.AheadIsChs2('=', '=') {
			l.Read()
			l.Read()
			val = T_EQ_S
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_EQ
		} else {
			val = T_ASSIGN
		}
	case '<':
		if l.AheadIsCh('=') {
			l.Read()
			val = T_LE
		} else if l.AheadIsCh('<') {
			l.Read()
			if l.AheadIsCh('=') {
				l.Read()
				val = T_ASSIGN_BIT_LSH
			} else {
				val = T_LSH
			}
		} else {
			val = T_LT
		}
	case '>':
		if l.AheadIsCh('=') {
			l.Read()
			val = T_GE
		} else if l.AheadIsCh('>') {
			l.Read()
			if l.AheadIsCh('>') {
				l.Read()
				if l.AheadIsCh('=') {
					l.Read()
					val = T_ASSIGN_BIT_RSH_U
				} else {
					val = T_RSH_U
				}
			} else if l.AheadIsCh('=') {
				l.Read()
				val = T_ASSIGN_BIT_RSH
			} else {
				val = T_RSH
			}
		} else {
			val = T_GT
		}
	case '*':
		if l.AheadIsCh('*') {
			l.Read()
			if l.AheadIsCh('=') {
				val = T_ASSIGN_POW
			} else {
				val = T_POW
			}
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_MUL
		} else {
			val = T_MUL
		}
	case '|':
		if l.AheadIsCh('|') {
			l.Read()
			if l.AheadIsCh('=') {
				l.Read()
				val = T_ASSIGN_OR
			} else {
				val = T_OR
			}
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_BIT_OR
		} else {
			val = T_BIT_OR
		}
	case '&':
		if l.AheadIsCh('&') {
			l.Read()
			if l.AheadIsCh('=') {
				l.Read()
				val = T_ASSIGN_AND
			} else {
				val = T_AND
			}
		} else if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_BIT_AND
		} else {
			val = T_BIT_AND
		}
	case '%':
		if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_MOD
		} else {
			val = T_MOD
		}
	case '!':
		if l.AheadIsCh('=') {
			l.Read()
			if l.AheadIsCh('=') {
				l.Read()
				val = T_NE_S
			} else {
				val = T_NE
			}
		} else {
			val = T_NOT
		}
	case '~':
		val = T_BIT_NOT
	case '^':
		if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_BIT_XOR
		} else {
			val = T_BIT_XOR
		}
	case '/':
		if l.AheadIsCh('=') {
			l.Read()
			val = T_ASSIGN_DIV
		} else {
			if l.prev == nil || TokenKinds[l.prev.value].BeforeExpr {
				return l.readRegexp()
			} else {
				val = T_DIV
			}
		}
	}
	tok.value = val
	return tok
}

// here is an assertion, for any valid regexp, the backslash is always escaped if it appears
// at any point of the content of the regexp
// base on above assertion, here we read the regexp roughly by stepping the content until the
// close backslash is matched as well as no validation is applied on that content
func (l *Lexer) readRegexp() *Token {
	tok := l.newToken()
	pattern := l.NewOpenRange()
	escaped := false
	for {
		c := l.Peek()
		if IsLineTerminator(c) || c == utf8.RuneError {
			return l.errToken(tok)
		} else if c == '\\' {
			escaped = true
		} else if !escaped && c == '/' {
			break
		}
		l.Read()
	}
	pattern.hi = l.Pos()
	l.Read() // consume the end `/`

	flags := l.NewOpenRange()
	i := 0
	for {
		if l.aheadIsIdPart() {
			l.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		flags = nil
	} else {
		flags.hi = l.Pos()
	}

	tok.ext = &TokExtRegexp{pattern, flags}
	return l.finToken(tok, T_REGEXP)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#sec-literals-string-literals
func (l *Lexer) ReadStr() *Token {
	tok := l.newToken()
	open := l.Read()
	text := make([]rune, 0, 10)
	for {
		c := l.Read()
		if c == utf8.RuneError || c == EOF {
			return l.errToken(tok)
		} else if c == '\\' {
			nc := l.Peek()
			if IsLineTerminator(nc) {
				l.readLineTerminator() // LineContinuation
			} else {
				r := l.readEscapeSeq()
				if r == utf8.RuneError || r == EOF {
					return l.errToken(tok)
				}
				text = append(text, r)
			}
		} else if c == open {
			break
		} else {
			text = append(text, c)
		}
	}
	tok.ext = &TokExtStr{open}
	tok.text = string(text)
	return l.finToken(tok, T_STRING)
}

func (l *Lexer) readEscapeSeq() rune {
	c := l.Read()
	switch c {
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'v':
		return '\v'
	case '0', '1', '2', '3', '4', '5', '6', '7':
		return l.readOctalEscapeSeq(c)
	case 'x':
		return l.readHexEscapeSeq()
	case 'u':
		return l.readUnicodeEscapeSeq()
	}
	return c
}

// https://tc39.es/ecma262/multipage/additional-ecmascript-features-for-web-browsers.html#prod-annexB-LegacyOctalEscapeSequence
// TODO: disabled in strict mode
func (l *Lexer) readOctalEscapeSeq(first rune) rune {
	octal := make([]rune, 0, 3)
	octal = append(octal, first)
	zeroToThree := first >= '0' && first <= '3'
	i := 1
	if first != '0' && l.isMode(LM_TEMPLATE) {
		// octal escape sequences are not allowed in template strings
		return utf8.RuneError
	}
	for {
		if !zeroToThree && i == 2 || zeroToThree && i == 3 {
			break
		}
		c := l.Peek()
		if !IsOctalDigit(c) {
			break
		}
		octal = append(octal, l.Read())
		i += 1
	}
	r, err := strconv.ParseInt(string(octal[:]), 8, 32)
	if err != nil {
		return utf8.RuneError
	}
	return rune(r)
}

func (l *Lexer) readHexEscapeSeq() rune {
	hex := [2]rune{}
	c := l.Read()
	for i := 0; i < 2; i++ {
		if IsHexDigit(c) {
			hex[i] = c
		} else {
			return utf8.RuneError
		}
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return utf8.RuneError
	}
	return rune(r)
}

func (l *Lexer) readLineTerminator() {
	c := l.Read()
	if c == '\r' {
		l.ReadIfNextIs('\n')
	}
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-NumericLiteral
func (l *Lexer) ReadNum() *Token {
	tok := l.newToken()
	c := l.Read()
	if c == '0' {
		switch l.Peek() {
		case 'b', 'B':
			return l.readBinaryNum(tok)
		case 'o', 'O':
			return l.readOctalNum(tok)
		case 'x', 'X':
			return l.readHexNum(tok)
		}
	}
	return l.readDecimalNum(tok, c)
}

func (l *Lexer) readDecimalNum(tok *Token, first rune) *Token {
	isFractionOpt := first == '.'
	if first != '.' && first != '0' {
		l.readDecimalDigits(true)
	}

	if first != '.' && l.AheadIsCh('.') || first == '.' {
		// read the fraction part
		if err := l.readDecimalDigits(isFractionOpt); err != nil {
			return l.errToken(tok)
		}
	}

	if l.AheadIsChOr('e', 'E') {
		if err := l.readExpPart(); err != nil {
			return l.errToken(tok)
		}
	}

	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readExpPart() error {
	l.Read() // consume `e` or `E`
	if l.AheadIsChOr('+', '-') {
		l.Read()
	}
	return l.readDecimalDigits(false)
}

func (l *Lexer) readDecimalDigits(opt bool) error {
	err := l.errCharError()
	i := 0
	for {
		c := l.Peek()
		if IsDecimalDigit(c) || i != 0 && c == '_' {
			l.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 && !opt {
		return err
	}
	return nil
}

func (l *Lexer) readBinaryNum(tok *Token) *Token {
	l.Read()
	i := 0
	for {
		c := l.Peek()
		if c == '0' || c == '1' || c == '_' {
			l.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(tok)
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readOctalNum(tok *Token) *Token {
	l.Read()
	i := 0
	for {
		c := l.Peek()
		if c >= '0' && c <= '7' || c == '_' {
			l.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(tok)
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readHexNum(tok *Token) *Token {
	l.Read()
	i := 0
	for {
		c := l.Peek()
		if IsHexDigit(c) || c == '_' {
			l.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(tok)
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readIdStart() rune {
	c := l.Read()
	return l.readUnicodeEscape(c)
}

func (l *Lexer) readIdPart() ([]rune, bool) {
	runes := make([]rune, 0, 10)
	for {
		c := l.Peek()
		if IsIdPart(c) {
			c := l.readUnicodeEscape(l.Read())
			if c == utf8.RuneError {
				return nil, false
			}
			runes = append(runes, c)
		} else {
			break
		}
	}
	return runes, true
}

func (l *Lexer) readUnicodeEscape(c rune) rune {
	if c == '\\' && l.AheadIsCh('u') {
		l.Read()
		return l.readUnicodeEscapeSeq()
	}
	return c
}

func (l *Lexer) readUnicodeEscapeSeq() rune {
	if l.AheadIsCh('{') {
		return l.readCodepoint()
	}
	return l.readHex4Digits()
}

func (l *Lexer) readCodepoint() rune {
	hex := make([]byte, 0, 4)
	l.Read() // consume `{`
	for {
		if l.AheadIsChThenConsume('}') {
			break
		} else if l.AheadIsEof() {
			return utf8.RuneError
		} else {
			c := l.Read()
			if c == utf8.RuneError || !IsHexDigit(c) {
				return utf8.RuneError
			}
			hex = append(hex, byte(c))
		}
	}
	r, err := strconv.ParseInt(string(hex), 16, 32)
	if err != nil {
		return utf8.RuneError
	}
	return rune(r)
}

func (l *Lexer) readHex4Digits() rune {
	hex := [4]byte{0}
	for i := 0; i < 4; i++ {
		c := l.Read()
		if c == utf8.RuneError || !IsHexDigit(c) {
			return utf8.RuneError
		}
		hex[i] = byte(c)
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return utf8.RuneError
	}
	return rune(r)
}

func (l *Lexer) error(msg string) *LexerError {
	return NewLexerError(msg, l.path, l.line, l.Pos()-1)
}

func (l *Lexer) errCharError() *LexerError {
	return l.error("unexpected character")
}

func (l *Lexer) aheadIsIdStart() bool {
	return IsIdStart(l.Peek())
}

func (l *Lexer) aheadIsIdPart() bool {
	return IsIdPart(l.Peek())
}

func (l *Lexer) aheadIsNumStart() bool {
	v := l.Peek()
	if IsDecimalDigit(v) {
		return true
	}
	return v == '.' && IsDecimalDigit(l.peek())
}

func (l *Lexer) aheadIsStrStart() bool {
	v := l.Peek()
	return v == '\'' || v == '"'
}

func (l *Lexer) aheadIsTplStart() bool {
	return l.Peek() == '`' || l.isMode(LM_TEMPLATE) && l.AheadIsCh('}')
}

func (l *Lexer) newToken() *Token {
	return &Token{
		value: T_ILLEGAL,
		loc:   l.NewOpenRange(),
	}
}

func (l *Lexer) finToken(tok *Token, value TokenValue) *Token {
	tok.value = value
	tok.loc.hi = l.Pos()
	tok.afterLineTerminator = l.metLineTerminator
	return tok
}

func (l *Lexer) errToken(tok *Token) *Token {
	tok.loc.hi = l.Pos()
	tok.ext = l.errCharError()
	return tok
}

func IsIdStart(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '$' || c == '_' ||
		unicode.In(c, unicode.Upper, unicode.Lower,
			unicode.Title, unicode.Modi,
			unicode.Other_Lowercase,
			unicode.Other_Uppercase,
			unicode.Other_ID_Start) ||
		c == '\\'
}

func IsIdPart(c rune) bool {
	return IsIdStart(c) || c >= '0' && c <= '9' || c == 0x200C || c == 0x200D
}

func IsOctalDigit(c rune) bool {
	return c >= '0' && c <= '7'
}

func IsHexDigit(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'F'
}

func IsDecimalDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsSingleEscapeChar(c rune) bool {
	return c == '\'' || c == '"' || c == '\\' || c == 'b' ||
		c == 'f' || c == 'n' || c == 'r' || c == 't' || c == 'v'
}

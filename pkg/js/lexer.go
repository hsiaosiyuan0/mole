package js

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	*Source
}

func NewLexer(src *Source) *Lexer {
	return &Lexer{Source: src}
}

func (l *Lexer) Next() *Token {
	l.SkipSpace()
	if l.aheadIsIdStart() {
		return l.ReadName()
	} else if l.aheadIsNumStart() {
		return l.ReadNum()
	} else if l.aheadIsStrStart() {
		return l.ReadStr()
	}
	return nil
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
	tok.text = string(runes)

	return l.finToken(tok, T_NAME)
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
				l.readLineTerminator()
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
	tok.ext = &TokenExtStr{open}
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
	err := l.unexpectedCharError()
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
		if IsId(c) {
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

func (l *Lexer) unexpectedCharError() *LexerError {
	return l.error("unexpected chart")
}

func (l *Lexer) aheadIsIdStart() bool {
	return IsIdStart(l.Peek())
}

func (l *Lexer) aheadIsNumStart() bool {
	v := l.Peek()
	return IsDecimalDigit(v) || v == '.'
}

func (l *Lexer) aheadIsStrStart() bool {
	v := l.Peek()
	return v == '\'' || v == '"'
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
	return l.errToken(tok)
}

func (l *Lexer) errToken(tok *Token) *Token {
	tok.loc.hi = l.Pos()
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

func IsId(c rune) bool {
	return IsIdStart(c) || IsIdPart(c)
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

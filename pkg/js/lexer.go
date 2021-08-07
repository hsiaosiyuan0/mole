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

func (l *Lexer) Next() (*Token, error) {
	l.SkipSpace()
	if l.aheadIsIdStart() {
		return l.ReadName()
	} else if l.aheadIsNumStart() {
		return l.ReadNum()
	}
	return nil, nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-IdentifierName
func (l *Lexer) ReadName() (*Token, error) {
	tok := l.newToken()

	runes := make([]rune, 0, 10)
	r, err := l.readIdStart()
	if err != nil {
		return tok, err
	}
	runes = append(runes, r)

	idPart, err := l.readIdPart()
	if err != nil {
		return tok, err
	}
	runes = append(runes, idPart...)
	tok.text = string(runes)

	return l.finToken(tok, T_NAME), nil
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-NumericLiteral
func (l *Lexer) ReadNum() (*Token, error) {
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

func (l *Lexer) readDecimalNum(tok *Token, first rune) (*Token, error) {
	isFractionOpt := first == '.'
	if first != '.' && first != '0' {
		l.readDecimalDigits(true)
	}

	if first != '.' && l.AheadIsCh('.') || first == '.' {
		// read the fraction part
		if err := l.readDecimalDigits(isFractionOpt); err != nil {
			return tok, err
		}
	}

	if l.AheadIsChOr('e', 'E') {
		if err := l.readExpPart(); err != nil {
			return tok, err
		}
	}

	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM), nil
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

func (l *Lexer) readBinaryNum(tok *Token) (*Token, error) {
	l.Read()
	err := l.unexpectedCharError()
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
		return tok, err
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM), nil
}

func (l *Lexer) readOctalNum(tok *Token) (*Token, error) {
	l.Read()
	err := l.unexpectedCharError()
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
		return tok, err
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM), nil
}

func (l *Lexer) readHexNum(tok *Token) (*Token, error) {
	l.Read()
	err := l.unexpectedCharError()
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
		return tok, err
	}
	l.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM), nil
}

func (l *Lexer) readIdStart() (rune, error) {
	c := l.Read()
	return l.readUnicodeEscape(c)
}

func (l *Lexer) readIdPart() ([]rune, error) {
	runes := make([]rune, 0, 10)
	for {
		c := l.Peek()
		if IsId(c) {
			c := l.Read() // advance
			c, err := l.readUnicodeEscape(c)
			if err != nil {
				return nil, err
			}
			runes = append(runes, c)
		} else {
			break
		}
	}
	return runes, nil
}

func (l *Lexer) readUnicodeEscape(c rune) (rune, error) {
	if c == '\\' && l.AheadIsCh('u') {
		l.Read()
		return l.readUnicodeEscapeSeq()
	}
	return c, nil
}

func (l *Lexer) readUnicodeEscapeSeq() (rune, error) {
	if l.AheadIsCh('{') {
		return l.readCodepoint()
	}
	return l.readHex4Digits()
}

func (l *Lexer) readCodepoint() (rune, error) {
	deformedHexErr := l.unexpectedCharError()
	hex := make([]byte, 0, 4)
	l.Read() // consume `{`
	for {
		if l.AheadIsChThenConsume('}') {
			break
		} else if l.AheadIsEof() {
			return 0, deformedHexErr
		} else {
			c := l.Read()
			if c == utf8.RuneError || !IsHexDigit(c) {
				return 0, deformedHexErr
			}
			hex = append(hex, byte(c))
		}
	}
	r, err := strconv.ParseInt(string(hex), 16, 32)
	if err != nil {
		return 0, deformedHexErr
	}
	return rune(r), nil
}

func (l *Lexer) readHex4Digits() (rune, error) {
	deformedHexErr := l.unexpectedCharError()
	hex := [4]byte{0}
	for i := 0; i < 4; i++ {
		c := l.Read()
		if c == utf8.RuneError || !IsHexDigit(c) {
			return 0, deformedHexErr
		}
		hex[i] = byte(c)
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return 0, deformedHexErr
	}
	return rune(r), nil
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

func (l *Lexer) newToken() *Token {
	return &Token{
		value: T_ILLEGAL,
		loc:   l.NewOpenRange(),
	}
}

func (l *Lexer) finToken(tok *Token, value TokenValue) *Token {
	tok.value = value
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
	return c >= 0 && c <= '7'
}

func IsHexDigit(c rune) bool {
	return c >= '0' && c <= '9' || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'F'
}

func IsDecimalDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

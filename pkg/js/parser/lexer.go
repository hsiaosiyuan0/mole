package parser

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type LexerModeValue int

const (
	LM_NONE       LexerModeValue = 0
	LM_STRICT                    = 1 << 0
	LM_TEMPLATE                  = 1 << 1 // for inline spans can tell they are in template string
	LM_ASYNC                     = 1 << 2
	LM_GENERATOR                 = 1 << 3
	LM_CLASS_BODY                = 1 << 4
	LM_CLASS_CTOR                = 1 << 5
	LM_NEW                       = 1 << 6
)

type LexerMode struct {
	value   LexerModeValue
	paren   int
	brace   int
	bracket int
}

const sizeOfPeekedTok = 5

type Lexer struct {
	src *Source

	ver  ESVersion
	mode []*LexerMode

	peeked [sizeOfPeekedTok]Token
	pl     int
	pr     int
	pw     int

	beginStmt bool
	notIn     bool

	prtVal   TokenValue  // the value of prev read token
	prtRng   SourceRange // the source range of prev read token
	prtBegin Pos         // the begin position of perv read token
	prtEnd   Pos         // the end position of prev read token

	pptVal      TokenValue // prev peek
	pptAfterEOL bool

	lastCommentLine int
	comments        map[int][]Token
}

func NewLexer(src *Source) *Lexer {
	lexer := &Lexer{src: src, mode: make([]*LexerMode, 0)}
	lexer.mode = append(lexer.mode, &LexerMode{LM_NONE, 0, 0, 0})
	lexer.comments = make(map[int][]Token)
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

func (l *Lexer) addMode(mode LexerModeValue) {
	cur := l.mode[len(l.mode)-1]
	cur.value |= mode
}

func (l *Lexer) isMode(mode LexerModeValue) bool {
	return l.curMode().value&mode > 0
}

func (l *Lexer) readTokWithComment() *Token {
	if l.src.SkipSpace().AheadIsEOF() {
		tok := l.newToken()
		tok.value = T_EOF
		return tok
	}

	if l.aheadIsIdStart() {
		return l.ReadName()
	} else if l.aheadIsNumStart() {
		return l.ReadNum()
	} else if l.aheadIsStrStart() {
		return l.ReadStr()
	} else if l.aheadIsTplStart() {
		return l.ReadTplSpan()
	} else if l.aheadIsPvt() {
		return l.ReadNumPvt()
	}
	return l.ReadSymbol()
}

func (l *Lexer) lastComment() *Token {
	line := l.comments[l.lastCommentLine]
	if line == nil {
		return nil
	}
	if len(line) == 0 {
		return nil
	}
	return &line[len(line)-1]
}

func (l *Lexer) readTok() *Token {
	prt := T_ILLEGAL
	var prtExt interface{}
	for {
		tok := l.readTokWithComment()
		if tok.value != T_COMMENT {
			if !tok.afterLineTerminator && prt == T_COMMENT && prtExt == true {
				tok.afterLineTerminator = true
			}
			return tok
		}
		l.lastCommentLine = tok.begin.line
		line := l.comments[tok.begin.line]
		if line == nil {
			line = make([]Token, 0)
		}
		line = append(line, *tok)
		l.comments[tok.begin.line] = line
		prt = tok.value
		prtExt = tok.ext
	}
}

func (l *Lexer) PeekGrow() *Token {
	if l.pl == sizeOfPeekedTok {
		panic(l.error(fmt.Sprintf("peek buffer of lexer is full, max len is %d\n", l.pl)))
	}

	tok := l.readTok()
	l.prtVal = tok.value
	l.pptAfterEOL = tok.afterLineTerminator

	l.beginStmt = false
	if tok.isEof() {
		return tok
	}

	l.pw += 1
	if l.pw == sizeOfPeekedTok {
		l.pw = 0
	}
	l.pl += 1
	return tok
}

// the line and column in Source maybe moved forward then their actual position
// that's because Lexer will reads tokens in buffer, so here firstly return Loc from
// the foremost peeked token otherwise return from Source if peeked buffer is empty
func (l *Lexer) Loc() *Loc {
	loc := NewLoc()
	loc.src = l.src
	if l.pl > 0 {
		tok := l.peeked[l.pr]
		p := tok.begin
		loc.begin.line = p.line
		loc.begin.col = p.col
		loc.rng.start = tok.raw.lo
	} else {
		loc.begin.line = l.src.line
		loc.begin.col = l.src.col
		loc.rng.start = l.src.Pos()
	}
	return loc
}

func (l *Lexer) FinLoc(loc *Loc) *Loc {
	if l.prtVal != T_ILLEGAL {
		p := l.prtEnd
		loc.end.line = p.line
		loc.end.col = p.col
		loc.rng.end = l.prtRng.hi
	} else {
		loc.end.line = l.src.line
		loc.end.col = l.src.col
		loc.rng.end = l.src.pos
	}
	return loc
}

func (l *Lexer) Peek() *Token {
	if l.pl > 0 {
		return &l.peeked[l.pr]
	}

	return l.PeekGrow()
}

func (l *Lexer) PeekStmtBegin() *Token {
	l.beginStmt = true
	tok := l.Peek()
	l.beginStmt = false
	return tok
}

func (l *Lexer) nextTok() *Token {
	if l.pl > 0 {
		tok := &l.peeked[l.pr]
		l.pr += 1
		if l.pr == sizeOfPeekedTok {
			l.pr = 0
		}
		l.pl -= 1
		return tok
	}
	return l.readTok()
}

func (l *Lexer) Next() *Token {
	tok := l.nextTok()
	l.prtVal = tok.value
	l.prtRng = tok.raw
	l.prtBegin = tok.begin
	l.prtEnd = tok.end
	return tok
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-Template
func (l *Lexer) ReadTplSpan() *Token {
	c := l.src.Read() // consume `\`` or `}`
	head := c == '`'
	if head {
		l.pushMode(LM_TEMPLATE)
	} else {
		l.popMode()
	}

	tok := l.newToken()
	text, fin, line, col, ofst, pos, err, _ := l.readTplChs()
	if text == nil {
		l.popMode()
		return l.errToken(tok, err)
	}

	tok.text = string(text)
	tok.value = T_ILLEGAL
	tok.raw.hi = ofst
	tok.len = pos - tok.len // tok.len stores the begin pos
	tok.end.line = line
	tok.end.col = col
	tok.afterLineTerminator = l.src.metLineTerminator
	ext := &TokExtTplSpan{false}
	tok.ext = ext

	if head {
		tok.value = T_TPL_HEAD
	} else {
		tok.value = T_TPL_SPAN
	}

	if fin {
		l.popMode()
		if head {
			ext.Plain = true
			return tok
		}
		tok.value = T_TPL_TAIL
		return tok
	}
	return tok
}

func (l *Lexer) readTplChs() (text []rune, fin bool, line, col, ofst, pos int, err string, legacyOctalEscapeSeq bool) {
	text = make([]rune, 0, 10)
	for {
		line = l.src.line
		col = l.src.col
		ofst = l.src.Ofst()
		pos = l.src.Pos()

		c := l.src.Peek()
		if c == '$' {
			l.src.Read()
			if l.src.AheadIsCh('{') {
				l.src.Read()
				l.pushMode(LM_TEMPLATE)
				break
			}
			text = append(text, c)
		} else if c == '\\' {
			l.src.Read()
			nc := l.src.Peek()
			if IsLineTerminator(nc) {
				l.readLineTerminator() // LineContinuation
			} else {
				r, e, lo := l.readEscapeSeq()
				if !legacyOctalEscapeSeq && lo {
					legacyOctalEscapeSeq = lo
				}
				if e != "" {
					err = e
					text = nil
					return
				}
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
			l.src.Read()
			fin = true
			break
		} else {
			text = append(text, l.src.Read())
		}
	}
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-IdentifierName
func (l *Lexer) ReadName() *Token {
	tok := l.newToken()

	runes := make([]rune, 0, 10)
	r, escapeInStart, err := l.readIdStart()
	if r == utf8.RuneError || err != "" {
		return l.errToken(tok, err)
	}
	runes = append(runes, r)

	col := l.src.col
	idPart, escapeInPart, err := l.readIdPart()
	if err != "" {
		tok.begin.col = col
		return l.errToken(tok, err)
	}
	runes = append(runes, idPart...)
	text := string(runes)

	containsEscape := escapeInStart || escapeInPart
	tok.ext = &TokExtIdent{containsEscape}
	if IsKeyword(text) {
		if containsEscape {
			return l.errToken(tok, ERR_ESCAPE_IN_KEYWORD)
		}
		return l.finToken(tok, Keywords[text])
	} else if l.isMode(LM_STRICT) && IsStrictKeyword(text) {
		return l.finToken(tok, StrictKeywords[text])
	} else if l.isMode(LM_ASYNC) && text == "await" {
		return l.finToken(tok, T_AWAIT)
	}
	tok.text = text
	return l.finToken(tok, T_NAME)
}

func (l *Lexer) aheadIsRegexp(afterLineTerminator bool) bool {
	if l.beginStmt {
		return true
	}
	prev := l.prtVal
	if prev == T_ILLEGAL {
		prev = l.pptVal
	}
	if prev == T_ILLEGAL {
		return true
	}

	if l.notIn && prev == T_IN {
		prev = T_NAME
	}
	be := TokenKinds[prev].BeforeExpr
	return be || afterLineTerminator
}

func (l *Lexer) ReadSymbol() *Token {
	tok := l.newToken()
	c := l.src.Read()
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
		if l.src.AheadIsChs2('.', '.') {
			l.src.Read()
			l.src.Read()
			val = T_DOT_TRI
		} else {
			val = T_DOT
		}
	case '?':
		if l.src.AheadIsCh('.') {
			l.src.Read()
			val = T_OPT_CHAIN
		} else if l.src.AheadIsCh('?') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				l.src.Read()
				val = T_ASSIGN_NULLISH
			} else {
				val = T_NULLISH
			}
		} else {
			val = T_HOOK
		}
	case '+':
		if l.src.AheadIsCh('+') {
			l.src.Read()
			val = T_INC
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_ADD
		} else {
			val = T_ADD
		}
	case '-':
		if l.src.AheadIsCh('-') {
			l.src.Read()
			val = T_DEC
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_SUB
		} else {
			val = T_SUB
		}
	case '=':
		if l.src.AheadIsCh('>') {
			l.src.Read()
			val = T_ARROW
		} else if l.src.AheadIsChs2('=', '=') {
			l.src.Read()
			l.src.Read()
			val = T_EQ_S
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_EQ
		} else {
			val = T_ASSIGN
		}
	case '<':
		if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_LTE
		} else if l.src.AheadIsCh('<') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				l.src.Read()
				val = T_ASSIGN_BIT_LSH
			} else {
				val = T_LSH
			}
		} else {
			val = T_LT
		}
	case '>':
		if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_GTE
		} else if l.src.AheadIsCh('>') {
			l.src.Read()
			if l.src.AheadIsCh('>') {
				l.src.Read()
				if l.src.AheadIsCh('=') {
					l.src.Read()
					val = T_ASSIGN_BIT_RSH_U
				} else {
					val = T_RSH_U
				}
			} else if l.src.AheadIsCh('=') {
				l.src.Read()
				val = T_ASSIGN_BIT_RSH
			} else {
				val = T_RSH
			}
		} else {
			val = T_GT
		}
	case '*':
		if l.src.AheadIsCh('*') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				val = T_ASSIGN_POW
			} else {
				val = T_POW
			}
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_MUL
		} else {
			val = T_MUL
		}
	case '|':
		if l.src.AheadIsCh('|') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				l.src.Read()
				val = T_ASSIGN_OR
			} else {
				val = T_OR
			}
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_BIT_OR
		} else {
			val = T_BIT_OR
		}
	case '&':
		if l.src.AheadIsCh('&') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				l.src.Read()
				val = T_ASSIGN_AND
			} else {
				val = T_AND
			}
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_BIT_AND
		} else {
			val = T_BIT_AND
		}
	case '%':
		if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_MOD
		} else {
			val = T_MOD
		}
	case '!':
		if l.src.AheadIsCh('=') {
			l.src.Read()
			if l.src.AheadIsCh('=') {
				l.src.Read()
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
		if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_BIT_XOR
		} else {
			val = T_BIT_XOR
		}
	case '/':
		if l.src.AheadIsCh('/') {
			return l.readSinglelineComment(tok)
		} else if l.src.AheadIsCh('*') {
			return l.readMultilineComment(tok)
		} else if l.aheadIsRegexp(l.src.metLineTerminator) {
			return l.readRegexp(tok)
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_DIV
		} else {
			val = T_DIV
		}

	}

	if val == T_DOT_TRI && l.ver < ES6 {
		return l.errToken(tok, ERR_MSG_UNEXPECTED_TOKEN)
	}

	return l.finToken(tok, val)
}

func (l *Lexer) readMultilineComment(tok *Token) *Token {
	l.src.Read() // consume `*`
	multiline := false
	for {
		c := l.src.Read()
		if c == '*' {
			if l.src.AheadIsCh('/') {
				l.src.Read()
				break
			}
		} else if c == EOL {
			multiline = true
		} else if c == EOF {
			return l.errToken(tok, ERR_UNTERMINATED_COMMENT)
		}
	}
	l.finToken(tok, T_COMMENT)
	tok.ext = multiline
	return tok
}

func (l *Lexer) readSinglelineComment(tok *Token) *Token {
	l.src.Read() // consume `/`
	for {
		c := l.src.Peek()
		if c == EOF || IsLineTerminator(c) {
			break
		}
		l.src.Read()
	}
	return l.finToken(tok, T_COMMENT)
}

// here is an assertion, for any valid regexp, the backslash is always escaped if it appears
// at any point of the content of the regexp
// base on above assertion, here we read the regexp roughly by stepping the content until the
// close backslash is matched as well as no validation is applied on that content
func (l *Lexer) readRegexp(tok *Token) *Token {
	pattern := l.src.NewOpenRange()
	for {
		c := l.src.Peek()
		if c == utf8.RuneError {
			return l.errToken(tok, "")
		} else if IsLineTerminator(c) {
			return l.errToken(tok, ERR_UNTERMINATED_REGEXP)
		}
		if c == '\\' {
			l.src.Read()
			nc := l.src.Peek()
			if !IsLineTerminator(nc) {
				l.src.Read()
			}
			continue
		}
		if c == '/' {
			break
		} else if c == EOF {
			tok.begin.col += 1
			return l.errToken(tok, ERR_UNTERMINATED_REGEXP)
		}
		l.src.Read()
	}
	pattern.hi = l.src.Ofst()
	l.src.Read() // consume the end `/`

	flags := l.src.NewOpenRange()
	i := 0
	for {
		if l.aheadIsIdPart(false) {
			col := l.src.col
			_, _, err := l.readIdPart()
			if err != "" {
				tok.begin.col = col
				return l.errToken(nil, err)
			}
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		flags = nil
	} else {
		flags.hi = l.src.Ofst()
	}

	tok.ext = &TokExtRegexp{pattern, flags}
	return l.finToken(tok, T_REGEXP)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#sec-literals-string-literals
func (l *Lexer) ReadStr() *Token {
	tok := l.newToken()
	open := l.src.Read()
	text := make([]rune, 0, 10)
	legacyOctalEscapeSeq := false
	for {
		c := l.src.Read()
		if c == utf8.RuneError || c == EOF {
			return l.errToken(tok, ERR_UNTERMINATED_STR)
		} else if c == '\\' {
			nc := l.src.Peek()
			if IsLineTerminator(nc) {
				l.readLineTerminator() // LineContinuation
			} else {
				r, err, lo := l.readEscapeSeq()
				if !legacyOctalEscapeSeq && lo {
					legacyOctalEscapeSeq = lo
				}
				if err != "" {
					return l.errToken(tok, err)
				}
				// allow `utf8.RuneError` to represent "Unicode replacement character"
				// in string literal
				if r == EOF {
					return l.errToken(tok, ERR_UNTERMINATED_STR)
				}
				text = append(text, r)
			}
		} else if IsLineTerminator(c) {
			return l.errToken(tok, ERR_UNTERMINATED_STR)
		} else if c == open {
			break
		} else {
			text = append(text, c)
		}
	}
	tok.ext = &TokExtStr{open, legacyOctalEscapeSeq}
	tok.text = string(text)
	return l.finToken(tok, T_STRING)
}

func (l *Lexer) readEscapeSeq() (r rune, errMsg string, octalEscapeSeq bool) {
	c := l.src.Read()
	switch c {
	case 'b':
		r = '\b'
		return
	case 'f':
		r = '\f'
		return
	case 'n':
		r = '\n'
		return
	case 'r':
		r = '\r'
		return
	case 't':
		r = '\t'
		return
	case 'v':
		r = '\v'
		return
	case '0', '1', '2', '3', '4', '5', '6', '7':
		octalEscapeSeq = true
		r, errMsg = l.readOctalEscapeSeq(c)
		return
	case 'x':
		r, errMsg = l.readHexEscapeSeq()
		return
	case 'u':
		r, errMsg = l.readUnicodeEscapeSeq(false)
		return
	}
	r = c
	return
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-LegacyOctalEscapeSequence
// TODO: disabled in strict mode
func (l *Lexer) readOctalEscapeSeq(first rune) (rune, string) {
	octal := make([]rune, 0, 3)
	octal = append(octal, first)
	zeroToThree := first >= '0' && first <= '3'
	i := 1
	if l.isMode(LM_TEMPLATE) {
		return utf8.RuneError, ERR_MSG_LEGACY_OCTAL_ESCAPE_IN_TPL
	}
	for {
		if !zeroToThree && i == 2 || zeroToThree && i == 3 {
			break
		}
		c := l.src.Peek()
		if !IsOctalDigit(c) {
			break
		}
		octal = append(octal, l.src.Read())
		i += 1
	}
	r, err := strconv.ParseInt(string(octal[:]), 8, 32)
	if err != nil {
		return utf8.RuneError, ""
	}
	return rune(r), ""
}

func (l *Lexer) readHexEscapeSeq() (rune, string) {
	hex := [2]rune{}
	hex[0] = l.src.Read()
	hex[1] = l.src.Read()
	if !IsHexDigit(hex[0]) || !IsHexDigit(hex[1]) {
		return utf8.RuneError, ""
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return utf8.RuneError, ""
	}
	return rune(r), ""
}

func (l *Lexer) readLineTerminator() {
	c := l.src.Read()
	if c == '\r' {
		l.src.ReadIfNextIs('\n')
	}
}

func (l *Lexer) ReadNumPvt() *Token {
	l.src.Read()
	tok := l.ReadName()
	tok.raw.lo -= 1
	if tok.value != T_NAME {
		return tok
	}
	tok.value = T_NAME_PVT
	return tok
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-NumericLiteral
func (l *Lexer) ReadNum() *Token {
	tok := l.newToken()
	c := l.src.Read()
	if c == '0' {
		switch l.src.Peek() {
		case 'b', 'B':
			return l.readBinaryNum(tok)
		case 'o', 'O':
			return l.readOctalNum(tok, 0)
		case 'x', 'X':
			return l.readHexNum(tok)
		}
		nc := l.src.Peek()
		if IsDecimalDigit(nc) {
			if IsOctalDigit(nc) {
				if l.isMode(LM_STRICT) {
					return l.errToken(tok, ERR_MSG_LEGACY_OCTAL_IN_STRICT_MODE)
				} else {
					return l.readOctalNum(tok, 1)
				}
			} else {
				return l.errToken(tok, ERR_MSG_INVALID_NUMBER)
			}
		}
	}
	return l.readDecimalNum(tok, c)
}

func (l *Lexer) readDecimalNum(tok *Token, first rune) *Token {
	if first != '.' && first != '0' {
		c := l.src.Peek()
		if c != 'e' && c != 'E' && c != 'n' && IsIdStart(c) {
			tok = l.newToken()
			return l.errToken(tok, ERR_MSG_IDENT_AFTER_NUMBER)
		}
		l.readDecimalDigits(true)
	}

	if first != '.' && l.src.AheadIsCh('.') || first == '.' {
		if l.src.AheadIsCh('.') {
			l.src.Read()
		}
		// read the fraction part
		if err := l.readDecimalDigits(true); err != nil {
			if IsIdStart(l.src.Peek()) {
				return l.errToken(nil, ERR_MSG_IDENT_AFTER_NUMBER)
			}
			return l.errToken(tok, ERR_MSG_INVALID_NUMBER)
		}
	}

	if l.src.AheadIsChOr('e', 'E') {
		if err := l.readExpPart(); err != nil {
			return l.errToken(tok, ERR_MSG_INVALID_NUMBER)
		}
	}

	l.src.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readExpPart() error {
	l.src.Read() // consume `e` or `E`
	if l.src.AheadIsChOr('+', '-') {
		l.src.Read()
	}
	return l.readDecimalDigits(false)
}

func (l *Lexer) readDecimalDigits(opt bool) error {
	err := l.errCharError()
	i := 0
	for {
		c := l.src.Peek()
		if IsDecimalDigit(c) || i != 0 && c == '_' {
			l.src.Read()
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
	l.src.Read()
	i := 0
	for {
		c := l.src.Peek()
		if c == '0' || c == '1' || c == '_' {
			l.src.Read()
			i += 1
		} else if IsIdStart(c) {
			return l.errToken(nil, ERR_MSG_IDENT_AFTER_NUMBER)
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(tok, ERR_MSG_INVALID_NUMBER)
	}
	l.src.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readOctalNum(tok *Token, i int) *Token {
	l.src.Read()
	for {
		c := l.src.Peek()
		if c >= '0' && c <= '7' || c == '_' {
			l.src.Read()
			i += 1
		} else if IsIdStart(c) {
			return l.errToken(nil, ERR_MSG_IDENT_AFTER_NUMBER)
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(tok, ERR_MSG_INVALID_NUMBER)
	}
	l.src.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readHexNum(tok *Token) *Token {
	l.src.Read()
	i := 0
	for {
		c := l.src.Peek()
		if IsHexDigit(c) || c == '_' {
			l.src.Read()
			i += 1
		} else if IsIdStart(c) {
			return l.errToken(nil, ERR_MSG_IDENT_AFTER_NUMBER)
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(nil, "Expected number in radix 16")
	}
	l.src.ReadIfNextIs('n')
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readIdStart() (r rune, containsEscape bool, errMsg string) {
	c := l.src.Read()
	return l.readUnicodeEscape(c, true)
}

func (l *Lexer) readIdPart() (rs []rune, containsEscape bool, errMsg string) {
	runes := make([]rune, 0, 10)
	for {
		c := l.src.Peek()
		if IsIdStart(c) || IsIdPart(c) {
			c, escape, err := l.readUnicodeEscape(l.src.Read(), true)
			if escape && !containsEscape {
				containsEscape = escape
			}
			if err != "" {
				return nil, escape, err
			} else if c == '\\' {
				return nil, escape, ERR_EXPECTING_UNICODE_ESCAPE
			}
			runes = append(runes, c)
		} else {
			break
		}
	}
	return runes, containsEscape, ""
}

func (l *Lexer) readUnicodeEscape(c rune, id bool) (r rune, containsEscape bool, errMsg string) {
	if c == '\\' {
		if l.src.AheadIsCh('u') {
			l.src.Read()
			containsEscape = true
			r, errMsg = l.readUnicodeEscapeSeq(id)
			return
		} else {
			return utf8.RuneError, false, ERR_EXPECTING_UNICODE_ESCAPE
		}
	}
	return c, false, ""
}

func (l *Lexer) readUnicodeEscapeSeq(id bool) (rune, string) {
	if l.src.AheadIsCh('{') {
		return l.readCodepoint()
	}
	return l.readHex4Digits(id)
}

func (l *Lexer) readCodepoint() (rune, string) {
	hex := make([]byte, 0, 4)
	l.src.Read() // consume `{`
	for {
		if l.src.ReadIfNextIs('}') {
			break
		} else if l.src.AheadIsEOF() {
			return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
		} else {
			c := l.src.Read()
			if c == utf8.RuneError || !IsHexDigit(c) {
				return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
			}
			hex = append(hex, byte(c))
		}
	}
	r, err := strconv.ParseInt(string(hex), 16, 32)
	if err != nil {
		return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
	}
	return rune(r), ""
}

func (l *Lexer) readHex4Digits(id bool) (rune, string) {
	hex := [4]byte{0}
	for i := 0; i < 4; i++ {
		c := l.src.Peek()
		if c == utf8.RuneError || !IsHexDigit(c) {
			return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
		}
		hex[i] = byte(l.src.Read())
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
	}
	rr := rune(r)
	if id {
		if !IsIdStart(rr) && !IsIdPart(rr) || rr == '\\' {
			return utf8.RuneError, ERR_INVALID_UNICODE_ESCAPE
		}
	}
	return rr, ""
}

func (l *Lexer) error(msg string) *LexerError {
	return NewLexerError(msg, l.src.path, l.src.line, l.src.Ofst())
}

func (l *Lexer) errCharError() *LexerError {
	return l.error("unexpected character")
}

func (l *Lexer) aheadIsIdStart() bool {
	return IsIdStart(l.src.Peek())
}

func (l *Lexer) aheadIsPvt() bool {
	return l.src.AheadIsCh('#')
}

func (l *Lexer) aheadIsIdPart(permitBackslash bool) bool {
	c := l.src.Peek()
	return IsIdStart(c) && (permitBackslash || c != '\\') || IsIdPart(c)
}

func (l *Lexer) aheadIsNumStart() bool {
	v := l.src.Peek()
	if IsDecimalDigit(v) {
		return true
	}
	return v == '.' && IsDecimalDigit(l.src.peekGrow())
}

func (l *Lexer) aheadIsStrStart() bool {
	v := l.src.Peek()
	return v == '\'' || v == '"'
}

func (l *Lexer) aheadIsTplStart() bool {
	return l.src.Peek() == '`' || l.isMode(LM_TEMPLATE) && l.src.AheadIsCh('}')
}

func (l *Lexer) newToken() *Token {
	tok := &l.peeked[l.pw]
	tok.value = T_ILLEGAL
	tok.text = ""
	l.src.openRange(&tok.raw)
	tok.begin.line = l.src.line
	tok.begin.col = l.src.col
	tok.end.line = l.src.line
	tok.end.col = l.src.col
	tok.len = l.src.Pos()
	return tok
}

func (l *Lexer) finToken(tok *Token, value TokenValue) *Token {
	tok.value = value
	tok.raw.hi = l.src.Ofst()
	tok.len = l.src.Pos() - tok.len // tok.len stores the begin pos
	tok.end.line = l.src.line
	tok.end.col = l.src.col
	tok.afterLineTerminator = l.src.metLineTerminator
	return tok
}

func (l *Lexer) errToken(tok *Token, msg string) *Token {
	if tok == nil {
		tok = l.newToken()
	}
	tok.raw.hi = l.src.Ofst()
	if msg != "" {
		tok.ext = msg
	} else {
		tok.ext = l.errCharError()
	}
	return tok
}

func IsIdStart(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '$' || c == '_' ||
		unicode.In(c, unicode.Upper, unicode.Lower,
			unicode.Title, unicode.Modi,
			unicode.Lo,
			unicode.Other_Lowercase,
			unicode.Other_Uppercase,
			unicode.Other_ID_Start) ||
		c == '\\'
}

func IsIdPart(c rune) bool {
	return c >= '0' && c <= '9' || c == 0x200C || c == 0x200D ||
		unicode.In(c,
			unicode.Pc,
			unicode.Mark,
			unicode.Other_ID_Continue)
}

func IsOctalDigit(c rune) bool {
	return c >= '0' && c <= '7'
}

func IsHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

func IsDecimalDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsSingleEscapeChar(c rune) bool {
	return c == '\'' || c == '"' || c == '\\' || c == 'b' ||
		c == 'f' || c == 'n' || c == 'r' || c == 't' || c == 'v'
}

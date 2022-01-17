package parser

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	span "github.com/hsiaosiyuan0/mole/span"
)

type LexerModeKind int

const (
	LM_NONE            LexerModeKind = 0
	LM_STRICT          LexerModeKind = 1 << iota
	LM_TEMPLATE        LexerModeKind = 1 << iota // for inline spans can tell they are in template string
	LM_TEMPLATE_TAGGED LexerModeKind = 1 << iota
	LM_ASYNC           LexerModeKind = 1 << iota
	LM_GENERATOR       LexerModeKind = 1 << iota
	LM_CLASS_BODY      LexerModeKind = 1 << iota
	LM_CLASS_CTOR      LexerModeKind = 1 << iota
	LM_NEW             LexerModeKind = 1 << iota
	LM_JSX             LexerModeKind = 1 << iota
	LM_JSX_CHILD       LexerModeKind = 1 << iota
	LM_JSX_ATTR        LexerModeKind = 1 << iota
)

type LexerMode struct {
	kind    LexerModeKind
	paren   int
	brace   int
	bracket int
}

const TOKENS_BUF_LEN = 5

type TokensBuf struct {
	buf [TOKENS_BUF_LEN]Token
	len int // len of the alive cells in `buf`
	r   int // the index in `buf` to perform next read
	w   int // the index in `buf` to perform next write
}

func (b *TokensBuf) incW() int {
	w := b.w + 1
	if w == TOKENS_BUF_LEN {
		return 0
	}
	return w
}

func (b *TokensBuf) incR() int {
	r := b.r + 1
	if r == TOKENS_BUF_LEN {
		return 0
	}
	return r
}

func (b *TokensBuf) writable() bool {
	return b.len < TOKENS_BUF_LEN
}

func (b *TokensBuf) readable() bool {
	return b.len > 0
}

func (b *TokensBuf) write(tok *Token) {
	if !b.writable() {
		panic("no place in tokens buf")
	}
	// just used to satisfy the static-checker
	_ = tok
	b.len += 1
	b.w = b.incW()
}

func (b *TokensBuf) newToken() *Token {
	return &b.buf[b.w]
}

func (b *TokensBuf) read() *Token {
	if !b.readable() {
		panic("tokens buf is empty")
	}
	t := &b.buf[b.r]
	b.len -= 1
	b.r = b.incR()
	return t
}

func (b *TokensBuf) cur() *Token {
	return &b.buf[b.r]
}

// consider the ambiguities introduced by TS grammar like `<`:
//
// ```ts
// a < b   // binExpr
// a<b>()  // callExpr
// ```
//
// for dealing with above problem, parser should store its state when
// the first `<` is met before performing the `argList` processing,
// the stored state will be restored if the `argList` processing was
// failed or be discarded if `argList` was succeeded
type LexerState struct {
	mode []LexerMode

	tb TokensBuf

	beginStmt bool

	prtVal   TokenValue // the value of prev read token
	prtRng   span.Range // the source range of prev read token
	prtBegin Pos        // the begin position of perv read token
	prtEnd   Pos        // the end position of prev read token

	pptVal TokenValue // prev peek

	// always save loc of the previous whitespace being skipped
	// by `skipSpace` in jsx mode
	prevWs Token
}

func newLexerState() LexerState {
	return LexerState{
		mode: []LexerMode{{LM_NONE, 0, 0, 0}},
	}
}

type Lexer struct {
	src  *span.Source
	ver  ESVersion
	feat Feature

	lastCmtLine int
	lastCmtCol  int
	comments    map[int]map[int]*span.Range

	state LexerState
	ss    []LexerState // state stack
}

func NewLexer(src *span.Source) *Lexer {
	lexer := &Lexer{
		src:      src,
		state:    newLexerState(),
		ss:       make([]LexerState, 0),
		comments: make(map[int]map[int]*span.Range),
	}
	return lexer
}

// read a token, named `next` to indicate it will move the cursor
func (l *Lexer) Next() *Token {
	tok := l.readTok()
	l.state.prtVal = tok.value
	l.state.prtRng = tok.raw
	l.state.prtBegin = tok.begin
	l.state.prtEnd = tok.end
	return tok
}

// for tokens like `in` and `of`, they are firstly read
// as names and then switched to keywords by the parser
// according to its context, so it's necessary to revise
// the `prtVal` of lexer to the corresponding of that
// keywords for satisfying the further lookheads
func (l *Lexer) NextRevise(v TokenValue) *Token {
	tok := l.readTok()
	l.state.prtVal = v
	l.state.prtRng = tok.raw
	l.state.prtBegin = tok.begin
	l.state.prtEnd = tok.end
	return tok
}

func (l *Lexer) Peek() *Token {
	if !l.state.tb.readable() {
		return l.advance()
	}
	return l.state.tb.cur()
}

func (l *Lexer) PeekStmtBegin() *Token {
	l.state.beginStmt = true
	tok := l.Peek()
	l.state.beginStmt = false
	return tok
}

func (l *Lexer) PeekGrow() *Token {
	return l.advance()
}

// guard the peeked buffer has at least 2 tokens, return
// the 2nd if the guarding is succeeded otherwise return
// the `EOF_TOK`
func (l *Lexer) Peek2nd() *Token {
	if l.state.tb.len < 2 {
		l.advance()
	}
	if l.state.tb.len < 2 {
		l.advance()
	}
	if l.state.tb.len >= 2 {
		return &l.state.tb.buf[l.state.tb.incR()]
	}
	return l.finToken(l.newToken(), T_EOF)
}

// the line and column in Source maybe moved forward then their actual position
// that's because Lexer will reads tokens in buffer, so here firstly return Loc from
// the foremost peeked token otherwise return from Source if peeked buffer is empty
func (l *Lexer) Loc() *Loc {
	loc := NewLoc()
	loc.src = l.src
	if l.state.tb.readable() {
		tok := l.state.tb.cur()
		p := tok.begin
		loc.begin.line = p.line
		loc.begin.col = p.col
		loc.rng.start = tok.raw.Lo
	} else {
		loc.begin.line = l.src.Line()
		loc.begin.col = l.src.Col()
		loc.rng.start = l.src.Pos()
	}
	return loc
}

func (l *Lexer) FinLoc(loc *Loc) *Loc {
	if l.state.prtVal != T_ILLEGAL {
		p := l.state.prtEnd
		loc.end.line = p.line
		loc.end.col = p.col
		loc.rng.end = l.state.prtRng.Hi
	} else {
		loc.end.line = l.src.Line()
		loc.end.col = l.src.Col()
		loc.rng.end = l.src.Pos()
	}
	return loc
}

func (l *Lexer) PrevTok() TokenValue {
	return l.state.prtVal
}

func (l *Lexer) PrevTokBegin() *Pos {
	return &l.state.prtBegin
}

func (l *Lexer) PushState() {
	l.ss = append(l.ss, l.state)
}

func (l *Lexer) DiscardState() {
	last := len(l.ss) - 1
	l.ss = l.ss[:last]
}

func (l *Lexer) PopState() {
	last := len(l.ss) - 1
	rest, state := l.ss[:last], l.ss[last]
	l.ss = rest
	l.state = state
}

func (l *Lexer) PushMode(mode LexerModeKind, inherit bool) {
	if inherit {
		// only inherit the inheritable modes
		v := LM_NONE
		v |= l.CurMode().kind & LM_STRICT
		v |= l.CurMode().kind & LM_ASYNC
		mode |= v
	}
	l.state.mode = append(l.state.mode, LexerMode{mode, 0, 0, 0})
}

func (l *Lexer) PopMode() *LexerMode {
	mLen := len(l.state.mode)
	if mLen == 1 {
		l.state.mode[0] = LexerMode{LM_NONE, 0, 0, 0}
		return &l.state.mode[0]
	}
	rest, last := l.state.mode[:mLen-1], l.state.mode[mLen-1]
	l.state.mode = rest
	return &last
}

func (l *Lexer) CurMode() *LexerMode {
	return &l.state.mode[len(l.state.mode)-1]
}

func (l *Lexer) AddMode(mode LexerModeKind) {
	cur := &l.state.mode[len(l.state.mode)-1]
	cur.kind |= mode
}

func (l *Lexer) IsMode(mode LexerModeKind) bool {
	return l.CurMode().kind&mode > 0
}

func (l *Lexer) skipSpace() *span.Source {
	if l.feat&FEAT_JSX == 0 {
		l.src.SkipSpace()
		return l.src
	}

	prevWs := &l.state.prevWs
	prevWs.raw.Src = l.src
	prevWs.begin.line = l.src.Line()
	prevWs.begin.col = l.src.Col()
	prevWs.raw.Lo = l.src.Ofst()
	prevWs.len = l.src.Pos()
	l.src.SkipSpace()
	l.finToken(prevWs, T_ILLEGAL)
	return l.src
}

func (l *Lexer) readTokWithComment() *Token {
	if l.skipSpace().AheadIsEOF() {
		return l.finToken(l.newToken(), T_EOF)
	}

	if !l.IsMode(LM_JSX) && !l.IsMode(LM_JSX_CHILD) {
		if l.aheadIsIdStart() {
			return l.readName()
		} else if l.aheadIsNumStart() {
			return l.readNum()
		} else if l.aheadIsStrStart() {
			return l.readStr()
		} else if l.aheadIsTplStart() {
			return l.readTplSpan()
		} else if l.aheadIsPvt() {
			return l.readNamePvt()
		}
		return l.readSymbol()
	}

	tok := l.newToken()
	c := l.src.Peek()
	// `{` used to enter attribute value or child expr
	if c == '{' {
		l.src.Read()
		l.PushMode(LM_NONE, true)
		return l.finToken(tok, T_BRACE_L)
	} else if c == '}' {
		l.src.Read()
		l.PopMode()
		return l.finToken(tok, T_BRACE_R)
	} else if c == '<' {
		l.src.Read()
		return l.finToken(tok, T_LT)
	}

	// in pair of open and close tags
	if l.IsMode(LM_JSX) {
		switch c {
		case ':':
			l.src.Read()
			return l.finToken(tok, T_COLON)
		case '/':
			l.src.Read()
			return l.finToken(tok, T_DIV)
		case '>':
			l.src.Read()
			return l.finToken(tok, T_GT)
		case '=':
			l.src.Read()
			return l.finToken(tok, T_ASSIGN)
		case '"', '\'':
			if l.IsMode(LM_JSX_ATTR) {
				return l.readStr()
			}
		}
		if l.aheadIsIdStart() {
			return l.readName()
		}
		return l.errToken(l.newToken(), ERR_UNEXPECTED_TOKEN)
	}

	// in child
	return l.readJSXTxt()
}

func (l *Lexer) lastComment() *span.Range {
	line := l.comments[l.lastCmtLine]
	if line == nil {
		return nil
	}
	if len(line) == 0 {
		return nil
	}
	return line[l.lastCmtCol]
}

func (l *Lexer) appendCmt(tok *Token) {
	l.lastCmtLine = tok.begin.line
	line := l.comments[tok.begin.line]
	if line == nil {
		line = make(map[int]*span.Range, 0)
		l.comments[tok.begin.line] = line
	}
	line[tok.begin.col] = tok.raw.Clone()
	l.lastCmtCol = tok.begin.col
}

func (l *Lexer) lexTok() *Token {
	prt := T_ILLEGAL
	var prtExt interface{}
	for {
		tok := l.readTokWithComment()
		if tok.value != T_COMMENT {
			if !tok.afterLineTerm && prt == T_COMMENT && prtExt == true {
				tok.afterLineTerm = true
			}
			return tok
		}
		l.appendCmt(tok)
		prt = tok.value
		prtExt = tok.ext
	}
}

func (l *Lexer) advance() *Token {
	tok := l.lexTok()

	// only update `prtVal` for telling whether ahead is regexp or not
	// other prt fields such as `prtRng` will be updated after `Next` is
	// called
	l.state.prtVal = tok.value

	if tok.value == T_EOF {
		return tok
	}

	l.state.tb.write(tok)
	return tok
}

func (l *Lexer) readTok() *Token {
	if l.state.tb.readable() {
		return l.state.tb.read()
	}
	return l.lexTok()
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-Template
func (l *Lexer) readTplSpan() *Token {
	span := l.newToken()
	c := l.src.Read() // consume `\`` or `}`
	head := c == '`'
	if head {
		l.PushMode(LM_TEMPLATE, true)
	} else {
		l.PopMode()
	}

	ext := &TokExtTplSpan{IllegalEscape: nil}

	// record the begin info of the internal string
	strBeginPos := l.src.Pos()
	l.src.OpenRange(&ext.strRng)
	ext.strBegin.line = l.src.Line()
	ext.strBegin.col = l.src.Col()

	// below ugly assignment is account for taking the in-place optimization(RVO), which can
	// use stack allocation instead of putting the values on heap, albeit it has not been proved
	text, fin, line, col, ofst, pos, err, _, errEscape, errEscapeLine, errEscapeCol := l.readTplChs()

	if err == ERR_UNTERMINATED_TPL {
		l.PopMode()
		return l.errToken(span, err)
	}

	if errEscape != "" {
		loc := NewLoc()
		loc.begin.line = errEscapeLine
		loc.begin.col = errEscapeCol
		ext.IllegalEscape = &IllegalEscapeInfo{errEscape, loc}
	}

	ext.str = string(text)
	ext.strRng.Hi = ofst
	ext.strLen = pos - strBeginPos
	ext.strEnd.line = line
	ext.strEnd.col = col

	span.afterLineTerm = l.src.MetLineTerm()
	span.ext = ext
	span = l.finToken(span, T_TPL_SPAN)

	if head {
		span.value = T_TPL_HEAD
	}

	if fin {
		l.PopMode()
		if head {
			ext.Plain = true
			return span
		}
		span.value = T_TPL_TAIL
		return span
	}
	return span
}

func (l *Lexer) readTplChs() (text []rune, fin bool, line, col, ofst, pos int,
	err string, legacyOctalEscapeSeq bool, errEscape string, errEscapeLine, errEscapeCol int) {
	text = make([]rune, 0, 10)
	for {
		line = l.src.Line()
		col = l.src.Col()
		ofst = l.src.Ofst()
		pos = l.src.Pos()

		c := l.src.Peek()
		if c == '$' {
			l.src.Read()
			if l.src.AheadIsCh('{') {
				l.src.Read()
				l.PushMode(LM_TEMPLATE, true)
				break
			}
			text = append(text, c)
		} else if c == '\\' {
			l.src.Read()
			nc := l.src.Peek()

			// LineContinuation
			if span.IsLineTerminator(nc) {
				l.readLineTermin()
				continue
			}

			// since the bad escape sequence is permitted in tagged template
			// here advance the cursor if `errEscape` already occurred
			if errEscape != "" {
				l.src.Read()
				continue
			}

			r, e, lo := l.readEscapeSeq()
			if !legacyOctalEscapeSeq && lo {
				legacyOctalEscapeSeq = lo
			}

			// just records the first occurred escape error
			if errEscape == "" {
				if e != "" {
					errEscape = e
				} else if r == utf8.RuneError {
					errEscape = ERR_BAD_ESCAPE_SEQ
				}
				errEscapeLine = line
				errEscapeCol = col
			}

			if r == span.EOF {
				err = ERR_UNTERMINATED_TPL
				return
			}
			text = append(text, r)
		} else if c == utf8.RuneError {
			errEscape = ERR_BAD_RUNE
		} else if c == span.EOF {
			err = ERR_UNTERMINATED_TPL
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
func (l *Lexer) readName() *Token {
	tok := l.newToken()

	runes := make([]rune, 0, 10)
	r, escapeInStart, err := l.readIdStart()
	if r == utf8.RuneError || err != "" {
		return l.errToken(tok, err)
	}
	runes = append(runes, r)

	col := l.src.Col()
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
	} else if l.IsMode(LM_STRICT) && IsStrictKeyword(text) {
		return l.finToken(tok, StrictKeywords[text])
	} else if text == "await" {
		if l.feat&FEAT_MODULE != 0 || (l.feat&FEAT_ASYNC_AWAIT != 0 && l.IsMode(LM_ASYNC)) {
			if containsEscape {
				return l.errToken(tok, ERR_ESCAPE_IN_KEYWORD)
			}
			return l.finToken(tok, T_AWAIT)
		}
	}

	tok.text = text
	return l.finToken(tok, T_NAME)
}

func (l *Lexer) aheadIsRegexp(afterLineTerm bool) bool {
	if l.IsMode(LM_JSX) || l.IsMode(LM_JSX_ATTR) {
		return false
	}

	if l.state.beginStmt {
		return true
	}

	// base on prev read firstly
	prev := l.state.prtVal
	if prev == T_ILLEGAL {
		prev = l.state.pptVal
	}

	// then try to base on prev peeked
	if prev == T_ILLEGAL {
		return true
	}

	be := TokenKinds[prev].BeforeExpr
	return be || afterLineTerm
}

func (l *Lexer) readSymbol() *Token {
	tok := l.newToken()
	c := l.src.Read()
	val := tok.value
	switch c {
	case '{':
		val = T_BRACE_L
		l.PushMode(LM_NONE, true)
	case '}':
		val = T_BRACE_R
		l.PopMode()
	case '(':
		val = T_PAREN_L
		l.CurMode().paren += 1
	case ')':
		val = T_PAREN_R
		l.CurMode().paren -= 1
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
		// a?.1:.0
		if l.src.AheadIsCh('.') && !IsDecimalDigit(l.src.Ahead2nd()) {
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
				l.src.Read()
				val = T_ASSIGN_POW
			} else {
				val = T_POW
			}
			if l.feat&FEAT_POW == 0 {
				return l.errToken(tok, ERR_UNEXPECTED_TOKEN)
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
		} else if l.aheadIsRegexp(l.src.MetLineTerm()) {
			return l.readRegexp(tok)
		} else if l.src.AheadIsCh('=') {
			l.src.Read()
			val = T_ASSIGN_DIV
		} else {
			val = T_DIV
		}

	}

	if val == T_DOT_TRI && (l.feat&FEAT_SPREAD == 0 || l.feat&FEAT_BINDING_REST_ELEM == 0) {
		return l.errToken(tok, ERR_UNEXPECTED_TOKEN)
	} else if val == T_OPT_CHAIN && l.feat&FEAT_OPT_EXPR == 0 {
		return l.errToken(tok, ERR_UNEXPECTED_TOKEN)
	} else if val == T_NULLISH && l.feat&FEAT_NULLISH == 0 {
		return l.errToken(tok, ERR_UNEXPECTED_TOKEN)
	} else if (val == T_ASSIGN_NULLISH || val == T_ASSIGN_AND || val == T_ASSIGN_OR) && l.feat&FEAT_LOGIC_ASSIGN == 0 {
		return l.errToken(tok, ERR_UNEXPECTED_TOKEN)
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
		} else if span.IsLineTerminator(c) {
			multiline = true
		} else if c == span.EOF {
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
		if c == span.EOF || span.IsLineTerminator(c) {
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
		} else if span.IsLineTerminator(c) {
			return l.errToken(nil, ERR_UNTERMINATED_REGEXP)
		}
		if c == '\\' {
			l.src.Read()
			nc := l.src.Peek()
			if nc == 'u' {
				l.src.Read()
				col := l.src.Col()
				_, errMsg := l.readUnicodeEscapeSeq(false)
				if errMsg != "" {
					tok.begin.col = col
					return l.errToken(tok, errMsg)
				}
			} else if !span.IsLineTerminator(nc) {
				l.src.Read()
			}
			continue
		}
		if c == '/' {
			break
		} else if c == span.EOF {
			tok.begin.col += 1
			return l.errToken(tok, ERR_UNTERMINATED_REGEXP)
		}
		l.src.Read()
	}
	pattern.Hi = l.src.Ofst()
	l.src.Read() // consume the end `/`

	flags := l.src.NewOpenRange()
	i := 0
	var err string
	var fs []rune
	if l.aheadIsIdPart(false) {
		col := l.src.Col()
		fs, _, err = l.readIdPart()
		if err != "" {
			tok.begin.col = col
			return l.errToken(nil, err)
		}
		i = len(fs)
	}

	if l.feat&FEAT_CHK_REGEXP_FLAGS != 0 {
		for _, f := range fs {
			if !l.isLegalFlag(f) {
				return l.errToken(tok, ERR_INVALID_REGEXP_FLAG)
			}
		}
	}

	if i == 0 {
		flags = nil
	} else {
		flags.Hi = l.src.Ofst()
	}

	tok.ext = &TokExtRegexp{pattern, flags}
	return l.finToken(tok, T_REGEXP)
}

func (l *Lexer) isLegalFlag(f rune) bool {
	switch f {
	case 'g', 'i', 'm':
		return true
	case 'd':
		return l.feat&FEAT_REGEXP_HAS_INDICES != 0
	case 'u':
		return l.feat&FEAT_REGEXP_UNICODE != 0
	case 'y':
		return l.feat&FEAT_REGEXP_STICKY != 0
	case 's':
		return l.feat&FEAT_REGEXP_DOT_ALL != 0
	}
	return false
}

func (l *Lexer) IsLineTerminator(c rune) bool {
	return c == 0x0a || c == 0x0d || ((c == 0x2028 || c == 0x2029) && l.feat&FEAT_JSON_SUPER_SET == 0)
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#sec-literals-string-literals
func (l *Lexer) readStr() *Token {
	tok := l.newToken()
	open := l.src.Read()
	text := make([]rune, 0, 10)
	legacyOctalEscapeSeq := false

	for {
		c := l.src.Read()
		if c == utf8.RuneError || c == span.EOF {
			return l.errToken(tok, ERR_UNTERMINATED_STR)
		} else if c == '\\' {
			nc := l.src.Peek()
			if span.IsLineTerminator(nc) {
				l.readLineTermin() // LineContinuation
			} else {
				line := l.src.Line()
				col := l.src.Col()
				r, err, lo := l.readEscapeSeq()
				if !legacyOctalEscapeSeq && lo {
					legacyOctalEscapeSeq = lo
				}
				if err != "" {
					tok.begin.line = line
					tok.begin.col = col
					return l.errToken(tok, err)
				}
				// allow `utf8.RuneError` to represent "Unicode replacement character"
				// in string literal
				if r == span.EOF {
					return l.errToken(tok, ERR_UNTERMINATED_STR)
				}
				text = append(text, r)
			}
		} else if l.IsLineTerminator(c) {
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
func (l *Lexer) readOctalEscapeSeq(first rune) (rune, string) {
	octal := make([]rune, 0, 3)
	octal = append(octal, first)
	zeroToThree := first >= '0' && first <= '3'
	i := 1
	if l.IsMode(LM_TEMPLATE) {
		return utf8.RuneError, ERR_TPL_LEGACY_OCTAL_ESCAPE_IN
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
		return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
	}
	r, err := strconv.ParseInt(string(hex[:]), 16, 32)
	if err != nil {
		return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
	}
	return rune(r), ""
}

func (l *Lexer) readLineTermin() {
	c := l.src.Read()
	if c == '\r' {
		l.src.AdvanceIf('\n')
	}
}

func (l *Lexer) readNamePvt() *Token {
	l.src.Read()
	if !l.aheadIsIdStart() {
		return l.errToken(l.newToken(), ERR_UNEXPECTED_CHAR)
	}
	tok := l.readName()
	tok.raw.Lo -= 1
	tok.begin.col -= 1
	if tok.value != T_NAME {
		return tok
	}
	tok.value = T_NAME_PVT
	if l.feat&FEAT_CLASS_PRV == 0 {
		tok.value = T_ILLEGAL
		return l.errToken(tok, ERR_UNEXPECTED_CHAR)
	}
	return tok
}

// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-NumericLiteral
func (l *Lexer) readNum() *Token {
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
				if l.IsMode(LM_STRICT) {
					return l.errToken(tok, ERR_LEGACY_OCTAL_IN_STRICT_MODE)
				} else {
					return l.readOctalNum(tok, 1)
				}
			} else {
				return l.errToken(tok, ERR_INVALID_NUMBER)
			}
		}
	}
	return l.readDecimalNum(tok, c)
}

func (l *Lexer) readDecimalNum(tok *Token, first rune) *Token {
	if first != '.' && first != '0' {
		c := l.src.Peek()
		if c != 'e' && c != 'E' && c != 'n' && IsIdStart(c) {
			return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
		}
		if err := l.readDecimalDigits(true); err != "" {
			if err == ERR_NUM_SEP_END {
				tok := l.newToken()
				tok.begin.col -= 1
				return l.errToken(tok, err)
			}
			return l.errToken(nil, err)
		}
	}

	float := false
	if first != '.' && l.src.AheadIsCh('.') || first == '.' {
		if l.src.AheadIsCh('.') {
			l.src.Read()
			float = true
		}
		// read the fraction part
		if err := l.readDecimalDigits(true); err != "" {
			if err == ERR_NUM_SEP_END {
				tok := l.newToken()
				tok.begin.col -= 1
				return l.errToken(tok, err)
			}
			return l.errToken(nil, err)
		}
	}

	exp := false
	if l.src.AheadIsChOr('e', 'E') {
		if err := l.readExpPart(); err != "" {
			return l.errToken(nil, err)
		}
		exp = true
	}

	if first != '.' && !float && !exp {
		if tok := l.bigintSuffix(); tok != nil {
			return l.errToken(tok, ERR_IDENT_AFTER_NUMBER)
		}
	}
	if IsIdStart(l.src.Peek()) {
		return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
	}
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) bigintSuffix() *Token {
	col := l.src.Col()
	if l.src.AdvanceIf('n') && l.feat&FEAT_BIGINT == 0 {
		tok := l.newToken()
		tok.begin.col = col
		return tok
	}
	return nil
}

func (l *Lexer) readExpPart() string {
	l.src.Read() // consume `e` or `E`
	if l.src.AheadIsChOr('+', '-') {
		l.src.Read()
	}
	return l.readDecimalDigits(false)
}

func (l *Lexer) readDecimalDigits(opt bool) string {
	i := 0

	var last rune
	for {
		c := l.src.Peek()
		if IsDecimalDigit(c) || c == '_' {
			if c == '_' && l.feat&FEAT_NUM_SEP == 0 {
				return ERR_INVALID_NUMBER
			}
			if i == 0 && c == '_' {
				return ERR_NUM_SEP_BEGIN
			} else {
				if last == '_' && c == '_' {
					return ERR_NUM_SEP_DUP
				}
				last = c
			}
			l.src.Read()
			i += 1
		} else {
			break
		}
	}
	if last == '_' {
		return ERR_NUM_SEP_END
	}
	if i == 0 && !opt {
		if IsIdStart(l.src.Peek()) {
			return ERR_IDENT_AFTER_NUMBER
		}
		return ERR_INVALID_NUMBER
	}
	return ""
}

func (l *Lexer) readBinaryNum(tok *Token) *Token {
	l.src.Read()
	i := 0
	var last rune
	for {
		c := l.src.Peek()
		if c == '0' || c == '1' || c == '_' {
			if c == '_' {
				if l.feat&FEAT_NUM_SEP == 0 {
					return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
				}
				if i == 0 {
					return l.errToken(nil, ERR_NUM_SEP_BEGIN)
				} else if last == '_' && c == '_' {
					return l.errToken(nil, ERR_NUM_SEP_DUP)
				}
			}
			last = c
			l.src.Read()
			i += 1
		} else if IsIdStart(c) {
			if c == 'n' && l.feat&FEAT_BIGINT != 0 {
				break
			}
			return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
		} else {
			break
		}
	}
	if last == '_' {
		tok := l.newToken()
		tok.begin.col -= 1
		return l.errToken(tok, ERR_NUM_SEP_END)
	}
	if i == 0 {
		tok.begin.line = l.src.Line()
		tok.begin.col = l.src.Col()
		return l.errToken(tok, fmt.Sprintf(ERR_TPL_EXPECT_NUM_RADIX, "2"))
	}

	if tok := l.bigintSuffix(); tok != nil {
		return l.errToken(tok, ERR_IDENT_AFTER_NUMBER)
	}
	if IsIdStart(l.src.Peek()) {
		return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
	}
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readOctalNum(tok *Token, i int) *Token {
	legacy := i == 1
	l.src.Read()
	var last rune
	for {
		c := l.src.Peek()
		if c >= '0' && c <= '7' || c == '_' {
			if c == '_' {
				if l.feat&FEAT_NUM_SEP == 0 {
					return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
				}
				if i == 0 {
					return l.errToken(nil, ERR_NUM_SEP_BEGIN)
				} else if legacy {
					return l.errToken(nil, ERR_NUM_SEP_IN_LEGACY_OCTAL)
				} else if last == '_' && c == '_' {
					return l.errToken(nil, ERR_NUM_SEP_DUP)
				}
			}
			last = c
			l.src.Read()
			i += 1
		} else {
			break
		}
	}
	if last == '_' {
		tok := l.newToken()
		tok.begin.col -= 1
		return l.errToken(tok, ERR_NUM_SEP_END)
	}
	if i == 0 {
		tok.begin.line = l.src.Line()
		tok.begin.col = l.src.Col()
		return l.errToken(tok, fmt.Sprintf(ERR_TPL_EXPECT_NUM_RADIX, "8"))
	}

	if !legacy {
		if tok := l.bigintSuffix(); tok != nil {
			return l.errToken(tok, ERR_IDENT_AFTER_NUMBER)
		}
	}
	if IsIdStart(l.src.Peek()) {
		return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
	}
	return l.finToken(tok, T_NUM)
}

func (l *Lexer) readHexNum(tok *Token) *Token {
	l.src.Read()
	i := 0
	var last rune
	for {
		c := l.src.Peek()
		if IsHexDigit(c) || c == '_' {
			if c == '_' {
				if l.feat&FEAT_NUM_SEP == 0 {
					return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
				}
				if i == 0 {
					return l.errToken(nil, ERR_NUM_SEP_BEGIN)
				} else if last == '_' && c == '_' {
					return l.errToken(nil, ERR_NUM_SEP_DUP)
				}
			}
			last = c
			l.src.Read()
			i += 1
		} else {
			break
		}
	}
	if i == 0 {
		return l.errToken(nil, "Expected number in radix 16")
	}
	if last == '_' {
		tok := l.newToken()
		tok.begin.col -= 1
		return l.errToken(tok, ERR_NUM_SEP_BEGIN)
	}

	if tok := l.bigintSuffix(); tok != nil {
		return l.errToken(tok, ERR_IDENT_AFTER_NUMBER)
	}
	if IsIdStart(l.src.Peek()) {
		return l.errToken(nil, ERR_IDENT_AFTER_NUMBER)
	}
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
		if l.src.AdvanceIf('}') {
			break
		} else if l.src.AheadIsEOF() {
			return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
		} else {
			c := l.src.Peek()
			if c == utf8.RuneError || !IsHexDigit(c) {
				return utf8.RuneError, ERR_BAD_ESCAPE_SEQ
			}
			hex = append(hex, byte(l.src.Read()))
		}
	}
	r, err := strconv.ParseInt(string(hex), 16, 32)
	if r > unicode.MaxRune {
		return utf8.RuneError, ERR_CODEPOINT_OUT_OF_BOUNDS
	}
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

func (l *Lexer) readJSXTxt() *Token {
	tok := l.newToken()

	var preWs string
	prevWs := &l.state.prevWs
	if prevWs.len != 0 {
		tok.begin.line = prevWs.begin.line
		tok.begin.col = prevWs.begin.col
		preWs = prevWs.Text()
		tok.raw.Lo = prevWs.raw.Lo
	}

	rs := make([]rune, 0)

	i := 0
	entity := make([]rune, 0, MaxHTMLEntityName)

	// loc of the first `&`
	ampLine := 0
	ampCol := 0
	for {
		c := l.src.Peek()
		line := l.src.Line()
		col := l.src.Col()
		if c == '{' || c == '<' || c == span.EOF {
			if i == 0 && c == span.EOF {
				l.src.Read()
				return l.finToken(tok, T_EOF)
			}
			break
		} else if c == '&' || i > 0 {
			if c == '&' {
				ampLine = line
				ampCol = col
			}

			l.src.Read()

			i += 1
			entity = append(entity, c)
			if c == ';' {
				key := string(entity[0:i])
				if ed, ok := HTMLEntities[key]; ok {
					rs = append(rs, ed.CodePoints...)
				} else {
					tok.begin.line = ampLine
					tok.begin.col = ampCol
					return l.errToken(tok, fmt.Sprintf(ERR_TPL_JSX_UNDEFINED_HTML_ENTITY, key))
				}
				i = 0
			}
		} else {
			rs = append(rs, l.src.Read())
		}
	}

	tok.ext = preWs + string(rs)
	return l.finToken(tok, T_JSX_TXT)
}

func (l *Lexer) error(msg string) *LexerError {
	return newLexerError(msg, l.src.Path, l.src.Line(), l.src.Ofst())
}

func (l *Lexer) errCharError() *LexerError {
	return l.error(ERR_UNEXPECTED_CHAR)
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
	return v == '.' && IsDecimalDigit(l.src.Ahead2nd())
}

func (l *Lexer) aheadIsStrStart() bool {
	v := l.src.Peek()
	return v == '\'' || v == '"'
}

func (l *Lexer) aheadIsTplStart() bool {
	return l.src.Peek() == '`' || l.IsMode(LM_TEMPLATE) && l.src.AheadIsCh('}')
}

func (l *Lexer) newToken() *Token {
	tok := l.state.tb.newToken()
	tok.value = T_ILLEGAL
	tok.text = ""
	l.src.OpenRange(&tok.raw)

	line := l.src.Line()
	col := l.src.Col()
	tok.begin.line = line
	tok.begin.col = col
	tok.end.line = line
	tok.end.col = col
	tok.len = l.src.Pos()
	return tok
}

func (l *Lexer) finToken(tok *Token, value TokenValue) *Token {
	tok.value = value
	tok.raw.Hi = l.src.Ofst()
	tok.len = l.src.Pos() - tok.len // tok.len stores the begin pos
	tok.end.line = l.src.Line()
	tok.end.col = l.src.Col()
	tok.afterLineTerm = l.src.MetLineTerm()
	return tok
}

func (l *Lexer) errToken(tok *Token, msg string) *Token {
	if tok == nil {
		tok = l.newToken()
	}
	tok.raw.Hi = l.src.Ofst()
	if msg != "" {
		tok.ext = msg
	} else {
		tok.ext = l.errCharError()
	}
	return tok
}

func IsIdStart(c rune) bool {
	if c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '$' || c == '_' ||
		c == '\\' {
		return true
	}
	if unicode.In(c, unicode.Upper, unicode.Lower,
		unicode.Title, unicode.Modi,
		unicode.Other_Lowercase,
		unicode.Other_Uppercase,
		unicode.Other_ID_Start) {
		return true
	}

	// CJK Unified Ideographs Extension D(U+2B740 to U+2B81F) is in `Lo`
	// but not permitted as the id start
	if unicode.In(c, unicode.Lo) && c <= 0x2B81D {
		return true
	}
	return false
}

func IsIdPart(c rune) bool {
	return c >= '0' && c <= '9' || c == 0x200C || c == 0x200D ||
		unicode.In(c,
			unicode.Pc,
			unicode.Mark,
			unicode.Nd,
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

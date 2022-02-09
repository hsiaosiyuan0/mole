package parser

import "github.com/hsiaosiyuan0/mole/span"

type Pos struct {
	line uint32
	col  uint32
}

func (p *Pos) Line() uint32 {
	return p.line
}

func (p *Pos) Column() uint32 {
	return p.col
}

func (p *Pos) Clone() *Pos {
	return &Pos{
		line: p.line,
		col:  p.col,
	}
}

type Token struct {
	value TokenValue
	text  string
	raw   span.Range
	begin Pos
	end   Pos

	// len of the codepoints in token. during the token is processing,
	// it will store the begin pos of that token
	len uint32

	// whether the token is after a line terminator or not
	afterLineTerm bool

	ext interface{}
}

func (t *Token) Begin() *Pos {
	return &t.begin
}

func (t *Token) End() *Pos {
	return &t.end
}

func (t *Token) Len() uint32 {
	return t.len
}

func (t *Token) IsLit(in bool) bool {
	v := t.value
	return v == T_VOID || v == T_NULL || v == T_TRUE || v == T_FALSE || v == T_STRING ||
		v == T_NUM || (in && v == T_IN)
}

func (t *Token) IsCtxKw() bool {
	v := t.value
	return (v > T_CTX_KEYWORD_BEGIN && v < T_CTX_KEYWORD_END) ||
		(v > T_CTX_KEYWORD_STRICT_BEGIN && v < T_CTX_KEYWORD_STRICT_END)
}

func (t *Token) IsKw() bool {
	v := t.value
	return (v > T_KEYWORD_BEGIN && v < T_KEYWORD_END) ||
		(v > T_CTX_KEYWORD_BEGIN && v < T_CTX_KEYWORD_END) ||
		(v > T_CTX_KEYWORD_STRICT_BEGIN && v < T_CTX_KEYWORD_STRICT_END) ||
		v == T_VOID || v == T_NULL || v == T_TRUE || v == T_FALSE || v == T_TYPE_OF ||
		v == T_DELETE || v == T_IN || v == T_INSTANCE_OF
}

func (t *Token) CanBePropKey() (string, bool, bool) {
	v := t.value
	if v == T_NAME {
		return t.text, false, true
	}
	if v == T_NUM {
		return t.raw.Text(), false, true
	}
	if t.IsKw() {
		return TokenKinds[v].Name, true, true
	}
	return "", false, false
}

func (t *Token) IsLegal() bool {
	return t.value != T_ILLEGAL
}

func (t *Token) RawText() string {
	return t.raw.Text()
}

func (t *Token) Text() string {
	if t.text != "" || t.value == T_STRING {
		return t.text
	}
	if name, _, ok := t.CanBePropKey(); ok {
		return name
	}
	return t.RawText()
}

func (t *Token) IsBin(notIn bool, ts bool) TokenValue {
	bin := t.value > T_BIN_OP_BEGIN && t.value < T_BIN_OP_END
	if bin {
		return t.value
	}
	if !notIn && IsName(t, "in", false) {
		return T_IN
	}
	if ts && IsName(t, "as", false) {
		return T_TS_AS
	}
	return T_ILLEGAL
}

func (t *Token) IsUnary() bool {
	return t.value > T_UNARY_OP_BEGIN && t.value < T_UNARY_OP_END
}

func (t *Token) Val() TokenValue {
	return t.value
}

func (t *Token) Kind() *TokenKind {
	return TokenKinds[t.value]
}

func (t *Token) IsPlainTpl() bool {
	if t.value != T_TPL_HEAD {
		return false
	}
	return t.ext.(*TokExtTplSpan).Plain
}

func (t *Token) HasLegacyOctalEscapeSeq() bool {
	if t.value == T_STRING {
		return t.ext.(*TokExtStr).LegacyOctalEscapeSeq
	}
	return false
}

func (t *Token) ContainsEscape() bool {
	if _, ok := t.ext.(*TokExtIdent); ok {
		return t.ext.(*TokExtIdent).ContainsEscape
	}
	return false
}

func (t *Token) ErrMsg() string {
	if msg, ok := t.ext.(string); ok {
		return msg
	}
	if msg, ok := t.ext.(*LexerError); ok {
		return msg.Error()
	}
	return "Unexpected character"
}

type TokExtStr struct {
	Open                 rune
	LegacyOctalEscapeSeq bool
}

type TokExtIdent struct {
	ContainsEscape bool
}

type IllegalEscapeInfo struct {
	Err string
	Loc *Loc
}

type TokExtTplSpan struct {
	// store the internal string
	str      string
	strLen   uint32
	strRng   span.Range
	strBegin Pos
	strEnd   Pos

	Plain bool

	// from ES2018 and later the tagged template can contain
	// illegal escape sequence, here records loc of the illegal
	// sequence to report lexer error under ES2018
	// see more details about that from: ES2018 revision of illegal escape sequences
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Template_literals
	IllegalEscape *IllegalEscapeInfo
}

type TokExtRegexp struct {
	pattern *span.Range
	flags   *span.Range
}

func (t *TokExtRegexp) Pattern() string {
	if t.pattern == nil {
		return ""
	}
	return t.pattern.Text()
}

func (t *TokExtRegexp) Flags() string {
	if t.flags == nil {
		return ""
	}
	return t.flags.Text()
}

type TokenValue int

const (
	T_ILLEGAL TokenValue = iota
	T_EOF
	T_COMMENT

	// literals
	T_NULL
	T_TRUE
	T_FALSE
	T_NUM
	T_STRING

	T_TPL_HEAD
	T_TPL_SPAN
	T_TPL_TAIL

	T_NAME
	T_NAME_PVT

	// keywords
	// https://tc39.es/ecma262/multipage/ecmascript-language-lexical-grammar.html#prod-ReservedWord
	T_KEYWORD_BEGIN
	T_BREAK
	T_CASE
	T_CATCH
	T_CLASS
	T_CONTINUE
	T_DEBUGGER
	T_DEFAULT
	T_DO
	T_ELSE
	T_ENUM
	T_EXPORT
	T_EXTENDS
	T_FINALLY
	T_FOR
	T_FUNC
	T_IF
	T_IMPORT
	T_NEW
	T_RETURN
	T_SUPER
	T_SWITCH
	T_THIS
	T_THROW
	T_TRY
	T_VAR
	T_WHILE
	T_WITH
	T_KEYWORD_END

	// contextual keywords
	T_CTX_KEYWORD_BEGIN
	// token treated as keyword in strict mode
	T_CTX_KEYWORD_STRICT_BEGIN
	T_LET
	T_CONST
	T_STATIC
	T_IMPLEMENTS
	T_INTERFACE
	T_PACKAGE
	T_PRIVATE
	T_PROTECTED
	T_PUBLIC
	T_YIELD
	T_CTX_KEYWORD_STRICT_END
	T_AS
	T_ASYNC
	T_AWAIT
	T_FROM
	T_GET
	T_META
	T_OF
	T_SET
	T_TARGET
	T_CTX_KEYWORD_END

	T_REGEXP
	T_BACK_QUOTE
	T_BRACE_L
	T_BRACE_R
	T_PAREN_L
	T_PAREN_R
	T_BRACKET_L
	T_BRACKET_R
	T_DOT
	T_DOT_TRI
	T_SEMI
	T_COMMA
	T_HOOK
	T_COLON
	T_INC
	T_DEC
	T_OPT_CHAIN
	T_ARROW

	T_BIN_OP_BEGIN
	T_NULLISH

	// relational
	T_LT
	T_GT
	T_LTE
	T_GTE

	// equality
	T_EQ
	T_NE
	T_EQ_S
	T_NE_S

	// bitwise
	T_LSH
	T_RSH
	T_RSH_U
	T_BIT_OR
	T_BIT_XOR
	T_BIT_AND

	T_OR
	T_AND

	T_INSTANCE_OF
	T_IN

	T_ADD
	T_SUB
	T_MUL
	T_DIV
	T_MOD
	T_POW
	T_BIN_OP_END

	// unary
	T_UNARY_OP_BEGIN
	T_TYPE_OF
	T_VOID
	T_DELETE
	T_NOT
	T_BIT_NOT
	T_UNARY_OP_END

	// assignment
	T_ASSIGN_BEGIN
	T_ASSIGN
	T_ASSIGN_ADD
	T_ASSIGN_SUB
	T_ASSIGN_NULLISH
	T_ASSIGN_OR
	T_ASSIGN_AND
	T_ASSIGN_BIT_OR
	T_ASSIGN_BIT_XOR
	T_ASSIGN_BIT_AND
	T_ASSIGN_BIT_LSH
	T_ASSIGN_BIT_RSH
	T_ASSIGN_BIT_RSH_U
	T_ASSIGN_MUL
	T_ASSIGN_DIV
	T_ASSIGN_MOD
	T_ASSIGN_POW
	T_ASSIGN_END

	T_JSX_TXT
	T_TS_AS
	T_TS_NO_NULL

	T_AT

	T_TOKEN_DEF_END
)

type TokenKind struct {
	Value      TokenValue
	Name       string
	Pcd        int
	RightAssoc bool

	// reference [acorn](https://github.com/acornjs/acorn/blob/master/acorn/src/tokentype.js)
	// `beforeExpr` is attached to each token to indicate that the slashes after those
	// tokens would be the beginning of regexps if the value of `beforeExpr` are `true`, works at
	// tokenizing phase
	BeforeExpr bool

	// `StartExpr` is attached to each token to indicate the token itself is the beginning of expr
	// it's used when parsing the argument of `yield`
	StartExpr bool
}

// order should be as same as `TokenValue`
var TokenKinds = [T_TOKEN_DEF_END + 1]*TokenKind{
	{T_ILLEGAL, "T_ILLEGAL", 0, false, false, false},
	{T_EOF, "EOF", 0, false, false, false},
	{T_COMMENT, "comment", 0, false, false, false},

	// literals
	{T_NULL, "null", 0, false, true, true},
	{T_TRUE, "true", 0, false, true, true},
	{T_FALSE, "false", 0, false, true, true},
	{T_NUM, "number", 0, false, false, true},
	{T_STRING, "string", 0, false, false, true},

	{T_TPL_HEAD, "template head", 0, false, true, true},
	{T_TPL_SPAN, "template span", 0, false, false, false},
	{T_TPL_TAIL, "template tail", 0, false, true, false},

	{T_NAME, "identifier", 0, false, false, true},
	{T_NAME_PVT, "private identifier", 0, false, false, false},

	// keywords
	{T_KEYWORD_BEGIN, "keyword begin", 0, false, false, false},
	{T_BREAK, "break", 0, false, false, false},
	{T_CASE, "case", 0, false, true, false},
	{T_CATCH, "catch", 0, false, false, false},
	{T_CLASS, "class", 0, false, false, true},
	{T_CONTINUE, "continue", 0, false, false, false},
	{T_DEBUGGER, "debugger", 0, false, false, false},
	{T_DEFAULT, "default", 0, false, true, false},
	{T_DO, "do", 0, false, true, false},
	{T_ELSE, "else", 0, false, true, false},
	{T_ENUM, "enum", 0, false, false, false},
	{T_EXPORT, "export", 0, false, false, false},
	{T_EXTENDS, "extends", 0, false, true, false},
	{T_FINALLY, "finally", 0, false, false, false},
	{T_FOR, "for", 0, false, false, false},
	{T_FUNC, "function", 0, false, false, true},
	{T_IF, "if", 0, false, false, false},
	{T_IMPORT, "import", 0, false, false, false},
	{T_NEW, "new", 0, true, true, true},
	{T_RETURN, "return", 0, false, true, false},
	{T_SUPER, "super", 0, false, false, true},
	{T_SWITCH, "switch", 0, false, false, false},
	{T_THIS, "this", 0, false, false, true},
	{T_THROW, "throw", 0, false, true, false},
	{T_TRY, "try", 0, false, false, false},
	{T_VAR, "var", 0, false, false, false},
	{T_WHILE, "while", 0, false, false, false},
	{T_WITH, "with", 0, false, false, false},
	{T_KEYWORD_END, "keyword end", 0, false, false, false},

	// contextual keywords
	{T_CTX_KEYWORD_BEGIN, "contextual keyword begin", 0, false, false, false},
	{T_CTX_KEYWORD_STRICT_BEGIN, "contextual keyword strict begin", 0, false, false, false},
	{T_LET, "let", 0, false, false, false},
	{T_CONST, "const", 0, false, false, false},
	{T_STATIC, "static", 0, false, false, false},
	{T_IMPLEMENTS, "implements", 0, false, false, false},
	{T_INTERFACE, "interface", 0, false, false, false},
	{T_PACKAGE, "package", 0, false, false, false},
	{T_PRIVATE, "private", 0, false, false, false},
	{T_PROTECTED, "protected", 0, false, false, false},
	{T_PUBLIC, "public", 0, false, false, false},
	{T_YIELD, "yield", 0, true, true, true},
	{T_CTX_KEYWORD_STRICT_END, "contextual keyword strict end", 0, false, false, false},
	{T_AS, "as", 0, false, false, false},
	{T_ASYNC, "async", 0, false, false, false},
	{T_AWAIT, "await", 0, true, false, false},
	{T_FROM, "from", 0, false, false, false},
	{T_GET, "get", 0, false, false, false},
	{T_META, "meta", 0, false, false, false},
	{T_OF, "of", 0, false, true, false},
	{T_SET, "set", 0, false, false, false},
	{T_TARGET, "target", 0, false, false, false},
	{T_CTX_KEYWORD_END, "contextual keyword end", 0, false, false, false},

	{T_REGEXP, "regexp", 0, false, false, true},
	{T_BACK_QUOTE, "`", 0, false, false, true},
	{T_BRACE_L, "{", 0, false, true, true},
	{T_BRACE_R, "}", 0, false, false, false},
	{T_PAREN_L, "(", 0, false, true, true},
	{T_PAREN_R, ")", 0, false, false, false},
	{T_BRACKET_L, "[", 0, false, true, true},
	{T_BRACKET_R, "]", 0, false, false, false},
	{T_DOT, ".", 0, false, false, false},
	{T_DOT_TRI, "...", 0, false, false, false},
	{T_SEMI, ";", 0, false, true, false},
	{T_COMMA, ",", 0, false, true, false},
	{T_HOOK, "?", 0, false, true, false},
	{T_COLON, ":", 0, false, true, false},
	{T_INC, "++", 0, true, false, true},
	{T_DEC, "--", 0, true, false, true},
	{T_OPT_CHAIN, "?.", 0, false, false, false},
	{T_ARROW, "=>", 0, false, true, false},

	{T_BIN_OP_BEGIN, "binary operator begin", 0, false, false, false},
	{T_NULLISH, "??", 6, false, true, false},

	// relational
	{T_LT, "<", 12, false, true, false},
	{T_GT, ">", 12, false, true, false},
	{T_LTE, "<=", 12, false, true, false},
	{T_GTE, ">=", 12, false, true, false},

	// equality
	{T_EQ, "==", 11, false, true, false},
	{T_NE, "!=", 11, false, true, false},
	{T_EQ_S, "===", 11, false, true, false},
	{T_NE_S, "!==", 11, false, true, false},

	// bitwise
	{T_LSH, "<<", 13, false, true, false},
	{T_RSH, ">>", 13, false, true, false},
	{T_RSH_U, ">>>", 13, false, true, false},
	{T_BIT_OR, "|", 8, false, true, false},
	{T_BIT_XOR, "^", 9, false, true, false},
	{T_BIT_AND, "&", 10, false, true, false},

	{T_OR, "||", 6, false, true, false},
	{T_AND, "&&", 7, false, true, false},

	{T_INSTANCE_OF, "instanceof", 12, false, true, false},
	{T_IN, "in", 12, false, true, false},

	{T_ADD, "+", 14, false, true, true},
	{T_SUB, "-", 14, false, true, true},
	{T_MUL, "*", 15, false, true, false},
	{T_DIV, "/", 15, false, true, false},
	{T_MOD, "%", 15, false, true, false},
	{T_POW, "**", 16, true, true, false},
	{T_BIN_OP_END, "binary operator end", 0, false, false, false},

	// unary
	{T_UNARY_OP_BEGIN, "unary operator being", 0, false, false, false},
	{T_TYPE_OF, "typeof", 0, true, false, true},
	{T_VOID, "void", 0, true, false, true},
	{T_DELETE, "delete", 0, true, false, true},
	{T_NOT, "!", 0, true, true, true},
	{T_BIT_NOT, "~", 0, true, true, true},
	{T_UNARY_OP_END, "unary operator end", 0, false, false, false},

	// assignment
	{T_ASSIGN_BEGIN, "assignment begin", 0, false, false, false},
	{T_ASSIGN, "=", 0, true, true, false},
	{T_ASSIGN_ADD, "+=", 0, true, true, false},
	{T_ASSIGN_SUB, "-=", 0, true, true, false},
	{T_ASSIGN_NULLISH, "??=", 0, true, true, false},
	{T_ASSIGN_OR, "||=", 0, true, true, false},
	{T_ASSIGN_AND, "&&=", 0, true, true, false},
	{T_ASSIGN_BIT_OR, "|=", 0, true, true, false},
	{T_ASSIGN_BIT_XOR, "^=", 0, true, true, false},
	{T_ASSIGN_BIT_AND, "&=", 0, true, true, false},
	{T_ASSIGN_BIT_LSH, "<<=", 0, true, true, false},
	{T_ASSIGN_BIT_RSH, ">>=", 0, true, true, false},
	{T_ASSIGN_BIT_RSH_U, ">>>=", 0, true, true, false},
	{T_ASSIGN_MUL, "*=", 0, true, true, false},
	{T_ASSIGN_DIV, "/=", 0, true, true, false},
	{T_ASSIGN_MOD, "%=", 0, true, true, false},
	{T_ASSIGN_POW, "**=", 0, true, true, false},
	{T_ASSIGN_END, "assignment end", 0, false, false, false},

	{T_JSX_TXT, "jsx text", 0, false, false, false},
	{T_TS_AS, "as", 12, false, true, false},
	{T_TS_NO_NULL, "!", 0, false, false, false},

	{T_AT, "@", 0, false, false, true},

	{T_TOKEN_DEF_END, "token end def", 0, false, false, false},
}

var Keywords = make(map[string]TokenValue)
var CtxKeywords = make(map[string]TokenValue)
var StrictKeywords = make(map[string]TokenValue)

func init() {
	for i := T_KEYWORD_BEGIN + 1; i < T_KEYWORD_END; i++ {
		Keywords[TokenKinds[i].Name] = i
	}

	// although below tokens are not keyword in strictly, put them
	// in the keywords map just for convenience
	Keywords[TokenKinds[T_VOID].Name] = T_VOID
	Keywords[TokenKinds[T_NULL].Name] = T_NULL
	Keywords[TokenKinds[T_TRUE].Name] = T_TRUE
	Keywords[TokenKinds[T_FALSE].Name] = T_FALSE
	Keywords[TokenKinds[T_TYPE_OF].Name] = T_TYPE_OF
	Keywords[TokenKinds[T_DELETE].Name] = T_DELETE
	Keywords[TokenKinds[T_INSTANCE_OF].Name] = T_INSTANCE_OF

	for i := T_CTX_KEYWORD_BEGIN + 1; i < T_CTX_KEYWORD_END; i++ {
		CtxKeywords[TokenKinds[i].Name] = i
	}
	for i := T_CTX_KEYWORD_STRICT_BEGIN + 1; i < T_CTX_KEYWORD_STRICT_END; i++ {
		StrictKeywords[TokenKinds[i].Name] = i
	}
}

func IsKeyword(str string) bool {
	_, ok := Keywords[str]
	return ok
}

func IsCtxKeyword(str string) bool {
	_, ok := CtxKeywords[str]
	return ok
}

func IsStrictKeyword(str string) bool {
	_, ok := StrictKeywords[str]
	return ok
}

func IsName(tok *Token, name string, canContainsEscape bool) bool {
	matched := tok.value == T_NAME && tok.Text() == name
	if !matched {
		return false
	}
	if !canContainsEscape {
		return !tok.ContainsEscape()
	}
	return true
}

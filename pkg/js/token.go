package js

type Token struct {
	value TokenValue
	text  string
	loc   *SourceRange
}

func (t *Token) Text() string {
	if t.text != "" {
		return t.text
	}
	loc := t.loc
	return loc.src.code[loc.lo:loc.hi]
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

	T_NAME
	T_NAME_PRIVATE

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
	T_LET
	T_STATIC
	T_IMPLEMENTS
	T_INTERFACE
	T_PACKAGE
	T_PRIVATE
	T_PROTECTED
	T_PUBLIC
	T_AS
	T_ASYNC
	T_FROM
	T_GET
	T_META
	T_OF
	T_SET
	T_TARGET
	T_YIELD
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
	T_LE
	T_GE

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
	T_ASSIGN_COALESCE
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

	T_TOKEN_DEF_END
)

type TokenKind struct {
	Value TokenValue
	Name  string

	// here is a manner taken from [acorn](https://github.com/acornjs/acorn/blob/master/acorn/src/tokentype.js)
	// it attaches a `beforeExpr` attribute to each token to indicate that the slashes after those tokens would be
	// the beginning of regexp if the value of their `beforeExpr` attributes are `true`, it works at tokenizing phase
	// therefore it can obey the definition of regexp to produce RegExprLiteral tokens
	BeforeExpr bool
}

// order should be as same as `TokenValue`
var TokenKinds = [T_TOKEN_DEF_END - 1]*TokenKind{
	{T_ILLEGAL, "T_ILLEGAL", false},
	{T_EOF, "end of script", false},
	{T_COMMENT, "comment", false},

	// literals
	{T_NULL, "null", false},
	{T_TRUE, "true", false},
	{T_FALSE, "false", false},
	{T_NUM, "number", false},
	{T_STRING, "string", false},
	{T_TPL_HEAD, "${", false},

	{T_NAME, "identifer", false},
	{T_NAME_PRIVATE, "private identifer", false},

	// keywords
	{T_KEYWORD_BEGIN, "keyword begin", false},
	{T_BREAK, "break", false},
	{T_CASE, "case", false},
	{T_CATCH, "catch", false},
	{T_CLASS, "class", false},
	{T_CONTINUE, "continue", false},
	{T_DEBUGGER, "debugger", false},
	{T_DEFAULT, "default", false},
	{T_DO, "do", false},
	{T_ELSE, "else", false},
	{T_ENUM, "enum", false},
	{T_EXPORT, "export", false},
	{T_EXTENDS, "extends", false},
	{T_FINALLY, "finally", false},
	{T_FOR, "for", false},
	{T_FUNC, "function", false},
	{T_IF, "if", false},
	{T_IMPORT, "import", false},
	{T_NEW, "new", false},
	{T_RETURN, "return", false},
	{T_SUPER, "super", false},
	{T_SWITCH, "switch", false},
	{T_THIS, "this", false},
	{T_THROW, "throw", false},
	{T_TRY, "try", false},
	{T_VAR, "var", false},
	{T_WHILE, "while", false},
	{T_WITH, "with", false},
	{T_CTX_KEYWORD_END, "keyword end", false},

	// contextual keywords
	{T_CTX_KEYWORD_BEGIN, "contextual keyword begin", false},
	{T_LET, "let", false},
	{T_STATIC, "static", false},
	{T_IMPLEMENTS, "implements", false},
	{T_INTERFACE, "interface", false},
	{T_PACKAGE, "package", false},
	{T_PRIVATE, "private", false},
	{T_PROTECTED, "protected", false},
	{T_PUBLIC, "public", false},
	{T_AS, "as", false},
	{T_ASYNC, "async", false},
	{T_FROM, "from", false},
	{T_GET, "get", false},
	{T_META, "meta", false},
	{T_OF, "of", false},
	{T_SET, "set", false},
	{T_TARGET, "target", false},
	{T_YIELD, "yield", false},
	{T_CTX_KEYWORD_END, "contextual keyword end", false},

	{T_REGEXP, "regexp", false},
	{T_BACK_QUOTE, "`", false},
	{T_BRACE_L, "{", false},
	{T_BRACE_R, "}", false},
	{T_PAREN_L, "(", false},
	{T_PAREN_R, ")", false},
	{T_BRACKET_L, "[", false},
	{T_BRACKET_R, "]", false},
	{T_DOT, ".", false},
	{T_DOT_TRI, "...", false},
	{T_SEMI, ";", false},
	{T_COMMA, ",", false},
	{T_HOOK, "?", false},
	{T_COLON, ":", false},
	{T_INC, "++", false},
	{T_DEC, "--", false},
	{T_OPT_CHAIN, "?.", false},
	{T_ARROW, "=>", false},

	{T_BIN_OP_BEGIN, "binary operator begin", false},
	{T_NULLISH, "??", false},

	// relational
	{T_LE, "<", false},
	{T_GE, ">", false},
	{T_LET, "<=", false},
	{T_GET, ">=", false},

	// equality
	{T_EQ, "==", false},
	{T_NE, "!=", false},
	{T_EQ_S, "===", false},
	{T_NE_S, "!==", false},

	// bitwise
	{T_LSH, "<<", false},
	{T_RSH, ">>", false},
	{T_RSH_U, ">>>", false},
	{T_BIT_OR, "|", false},
	{T_BIT_XOR, "^", false},
	{T_BIT_AND, "&", false},

	{T_OR, "||", false},
	{T_AND, "&&", false},

	{T_INSTANCE_OF, "instanceof", false},
	{T_IN, "in", false},

	{T_ADD, "+", false},
	{T_SUB, "-", false},
	{T_MUL, "*", false},
	{T_DIV, "/", false},
	{T_MOD, "%", false},
	{T_POW, "**", false},
	{T_BIN_OP_END, "binary operator end", false},

	// unary
	{T_UNARY_OP_BEGIN, "unary operator being", false},
	{T_TYPE_OF, "typeof", false},
	{T_VOID, "void", false},
	{T_DELETE, "delete", false},
	{T_NOT, "!", false},
	{T_BIT_NOT, "~", false},
	{T_UNARY_OP_END, "unary operator end", false},

	// assignment
	{T_ASSIGN_BEGIN, "assignment begin", false},
	{T_ASSIGN, "=", false},
	{T_ASSIGN_ADD, "+=", false},
	{T_ASSIGN_SUB, "-=", false},
	{T_ASSIGN_COALESCE, "??=", false},
	{T_ASSIGN_OR, "||=", false},
	{T_ASSIGN_AND, "&&=", false},
	{T_ASSIGN_BIT_OR, "|=", false},
	{T_ASSIGN_BIT_XOR, "^=", false},
	{T_ASSIGN_BIT_AND, "&=", false},
	{T_ASSIGN_BIT_LSH, "<<=", false},
	{T_ASSIGN_BIT_RSH, ">>=", false},
	{T_ASSIGN_BIT_RSH_U, ">>>=", false},
	{T_ASSIGN_MUL, "*=", false},
	{T_ASSIGN_DIV, "/=", false},
	{T_ASSIGN_MOD, "%=", false},
	{T_ASSIGN_POW, "**=", false},
}

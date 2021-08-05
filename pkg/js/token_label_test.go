package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mlint/pkg/assert"
)

func TestLabel(t *testing.T) {
	assert.Equal(t, "end of script", TokenLabel[T_EOF], "should be end of script")
	assert.Equal(t, "comment", TokenLabel[T_COMMENT], "should be comment")

	// literals
	assert.Equal(t, "null", TokenLabel[T_NULL], "should be null")
	assert.Equal(t, "true", TokenLabel[T_TRUE], "should be true")
	assert.Equal(t, "false", TokenLabel[T_FALSE], "should be false")
	assert.Equal(t, "number", TokenLabel[T_NUM], "should be number")
	assert.Equal(t, "string", TokenLabel[T_STRING], "should be string")
	assert.Equal(t, "${", TokenLabel[T_TPL_HEAD], "should be ${")

	assert.Equal(t, "identifer", TokenLabel[T_NAME], "should be identifer")
	assert.Equal(t, "private identifer", TokenLabel[T_NAME_PRIVATE], "should be private identifer")

	// keywords
	assert.Equal(t, "break", TokenLabel[T_BREAK], "should be break")
	assert.Equal(t, "case", TokenLabel[T_CASE], "should be case")
	assert.Equal(t, "catch", TokenLabel[T_CATCH], "should be catch")
	assert.Equal(t, "class", TokenLabel[T_CLASS], "should be class")
	assert.Equal(t, "continue", TokenLabel[T_CONTINUE], "should be continue")
	assert.Equal(t, "debugger", TokenLabel[T_DEBUGGER], "should be debugger")
	assert.Equal(t, "default", TokenLabel[T_DEFAULT], "should be default")
	assert.Equal(t, "do", TokenLabel[T_DO], "should be do")
	assert.Equal(t, "else", TokenLabel[T_ELSE], "should be else")
	assert.Equal(t, "enum", TokenLabel[T_ENUM], "should be enum")
	assert.Equal(t, "export", TokenLabel[T_EXPORT], "should be export")
	assert.Equal(t, "extends", TokenLabel[T_EXTENDS], "should be extends")
	assert.Equal(t, "finally", TokenLabel[T_FINALLY], "should be finally")
	assert.Equal(t, "for", TokenLabel[T_FOR], "should be for")
	assert.Equal(t, "function", TokenLabel[T_FUNC], "should be function")
	assert.Equal(t, "if", TokenLabel[T_IF], "should be if")
	assert.Equal(t, "import", TokenLabel[T_IMPORT], "should be import")
	assert.Equal(t, "new", TokenLabel[T_NEW], "should be new")
	assert.Equal(t, "return", TokenLabel[T_RETURN], "should be return")
	assert.Equal(t, "super", TokenLabel[T_SUPER], "should be super")
	assert.Equal(t, "switch", TokenLabel[T_SWITCH], "should be switch")
	assert.Equal(t, "this", TokenLabel[T_THIS], "should be this")
	assert.Equal(t, "throw", TokenLabel[T_THROW], "should be throw")
	assert.Equal(t, "try", TokenLabel[T_TRY], "should be try")
	assert.Equal(t, "var", TokenLabel[T_VAR], "should be var")
	assert.Equal(t, "while", TokenLabel[T_WHILE], "should be while")
	assert.Equal(t, "with", TokenLabel[T_WITH], "should be with")

	// contextual keywords
	assert.Equal(t, "let", TokenLabel[T_LET], "should be let")
	assert.Equal(t, "static", TokenLabel[T_STATIC], "should be static")
	assert.Equal(t, "implements", TokenLabel[T_IMPLEMENTS], "should be implements")
	assert.Equal(t, "interface", TokenLabel[T_INTERFACE], "should be interface")
	assert.Equal(t, "package", TokenLabel[T_PACKAGE], "should be package")
	assert.Equal(t, "private", TokenLabel[T_PRIVATE], "should be private")
	assert.Equal(t, "protected", TokenLabel[T_PROTECTED], "should be protected")
	assert.Equal(t, "public", TokenLabel[T_PUBLIC], "should be public")
	assert.Equal(t, "as", TokenLabel[T_AS], "should be as")
	assert.Equal(t, "async", TokenLabel[T_ASYNC], "should be async")
	assert.Equal(t, "from", TokenLabel[T_FROM], "should be from")
	assert.Equal(t, "get", TokenLabel[T_GET], "should be get")
	assert.Equal(t, "meta", TokenLabel[T_META], "should be meta")
	assert.Equal(t, "of", TokenLabel[T_OF], "should be of")
	assert.Equal(t, "set", TokenLabel[T_SET], "should be set")
	assert.Equal(t, "target", TokenLabel[T_TARGET], "should be target")
	assert.Equal(t, "yield", TokenLabel[T_YIELD], "should be yield")

	assert.Equal(t, "regexp", TokenLabel[T_REGEXP], "should be regexp")
	assert.Equal(t, "`", TokenLabel[T_BACK_QUOTE], "should be `")
	assert.Equal(t, "{", TokenLabel[T_BRACE_L], "should be `{")
	assert.Equal(t, "}", TokenLabel[T_BRACE_R], "should be }")
	assert.Equal(t, "(", TokenLabel[T_PAREN_L], "should be (")
	assert.Equal(t, ")", TokenLabel[T_PAREN_R], "should be )")
	assert.Equal(t, "[", TokenLabel[T_BRACKET_L], "should be [")
	assert.Equal(t, "]", TokenLabel[T_BRACKET_R], "should be ]")
	assert.Equal(t, ".", TokenLabel[T_DOT], "should be .")
	assert.Equal(t, "...", TokenLabel[T_DOT_TRI], "should be ...")
	assert.Equal(t, ";", TokenLabel[T_SEMI], "should be ;")
	assert.Equal(t, ",", TokenLabel[T_COMMA], "should be ,")
	assert.Equal(t, "?", TokenLabel[T_HOOK], "should be ?")
	assert.Equal(t, ":", TokenLabel[T_COLON], "should be :")
	assert.Equal(t, "++", TokenLabel[T_INC], "should be ++")
	assert.Equal(t, "--", TokenLabel[T_DEC], "should be --")
	assert.Equal(t, "?.", TokenLabel[T_OPT_CHAIN], "should be ?.")
	assert.Equal(t, "=>", TokenLabel[T_ARROW], "should be =>")

	assert.Equal(t, "??", TokenLabel[T_COALESCE], "should be ??")

	// relational
	assert.Equal(t, "<", TokenLabel[T_LT], "should be <")
	assert.Equal(t, ">", TokenLabel[T_GT], "should be >")
	assert.Equal(t, "<=", TokenLabel[T_LE], "should be <=")
	assert.Equal(t, ">=", TokenLabel[T_GE], "should be >=")

	// equality
	assert.Equal(t, "==", TokenLabel[T_EQ], "should be ==")
	assert.Equal(t, "!=", TokenLabel[T_NE], "should be !=")
	assert.Equal(t, "===", TokenLabel[T_EQ_S], "should be ===")
	assert.Equal(t, "!==", TokenLabel[T_NE_S], "should be !==")

	// bitwise
	assert.Equal(t, "<<", TokenLabel[T_LSH], "should be <<")
	assert.Equal(t, ">>", TokenLabel[T_RSH], "should be >>")
	assert.Equal(t, ">>>", TokenLabel[T_RSH_U], "should be >>>")
	assert.Equal(t, "|", TokenLabel[T_BIT_OR], "should be |")
	assert.Equal(t, "^", TokenLabel[T_BIT_XOR], "should be ^")
	assert.Equal(t, "&", TokenLabel[T_BIT_AND], "should be &")

	assert.Equal(t, "||", TokenLabel[T_OR], "should be ||")
	assert.Equal(t, "&&", TokenLabel[T_AND], "should be &&")

	assert.Equal(t, "instanceof", TokenLabel[T_INSTANCE_OF], "should be instanceof")
	assert.Equal(t, "in", TokenLabel[T_IN], "should be in")

	// unary
	assert.Equal(t, "+", TokenLabel[T_ADD], "should be +")
	assert.Equal(t, "-", TokenLabel[T_SUB], "should be -")
	assert.Equal(t, "*", TokenLabel[T_MUL], "should be *")
	assert.Equal(t, "/", TokenLabel[T_DIV], "should be /")
	assert.Equal(t, "%", TokenLabel[T_MOD], "should be %")
	assert.Equal(t, "**", TokenLabel[T_POW], "should be **")

	// assignment
	assert.Equal(t, "=", TokenLabel[T_ASSIGN], "should be =")
	assert.Equal(t, "+=", TokenLabel[T_ASSIGN_ADD], "should be +=")
	assert.Equal(t, "-=", TokenLabel[T_ASSIGN_SUB], "should be -=")
	assert.Equal(t, "??=", TokenLabel[T_ASSIGN_COALESCE], "should be ??=")
	assert.Equal(t, "||=", TokenLabel[T_ASSIGN_OR], "should be ||=")
	assert.Equal(t, "&&=", TokenLabel[T_ASSIGN_AND], "should be &&=")
	assert.Equal(t, "|=", TokenLabel[T_ASSIGN_BIT_OR], "should be |=")
	assert.Equal(t, "^=", TokenLabel[T_ASSIGN_BIT_XOR], "should be ^=")
	assert.Equal(t, "&=", TokenLabel[T_ASSIGN_BIT_AND], "should be &=")
	assert.Equal(t, "<<=", TokenLabel[T_ASSIGN_BIT_LSH], "should be <<=")
	assert.Equal(t, ">>=", TokenLabel[T_ASSIGN_BIT_RSH], "should be >>=")
	assert.Equal(t, ">>>=", TokenLabel[T_ASSIGN_BIT_RSH_U], "should be >>>=")
	assert.Equal(t, "*=", TokenLabel[T_ASSIGN_MUL], "should be *=")
	assert.Equal(t, "/=", TokenLabel[T_ASSIGN_DIV], "should be /=")
	assert.Equal(t, "%=", TokenLabel[T_ASSIGN_MOD], "should be %=")
	assert.Equal(t, "**=", TokenLabel[T_ASSIGN_POW], "should be **=")
}

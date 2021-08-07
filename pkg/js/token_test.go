package js

import (
	"testing"

	"github.com/hsiaosiyuan0/mlint/pkg/assert"
)

func TestLabel(t *testing.T) {
	assert.Equal(t, "end of script", TokenKinds[T_EOF].Name, "should be end of script")
	assert.Equal(t, "comment", TokenKinds[T_COMMENT].Name, "should be comment")

	// literals
	assert.Equal(t, "null", TokenKinds[T_NULL].Name, "should be null")
	assert.Equal(t, "true", TokenKinds[T_TRUE].Name, "should be true")
	assert.Equal(t, "false", TokenKinds[T_FALSE].Name, "should be false")
	assert.Equal(t, "number", TokenKinds[T_NUM].Name, "should be number")
	assert.Equal(t, "string", TokenKinds[T_STRING].Name, "should be string")
	assert.Equal(t, "${", TokenKinds[T_TPL_HEAD].Name, "should be ${")

	assert.Equal(t, "identifer", TokenKinds[T_NAME].Name, "should be identifer")
	assert.Equal(t, "private identifer", TokenKinds[T_NAME_PRIVATE].Name, "should be private identifer")

	// keywords
	assert.Equal(t, "break", TokenKinds[T_BREAK].Name, "should be break")
	assert.Equal(t, "case", TokenKinds[T_CASE].Name, "should be case")
	assert.Equal(t, "catch", TokenKinds[T_CATCH].Name, "should be catch")
	assert.Equal(t, "class", TokenKinds[T_CLASS].Name, "should be class")
	assert.Equal(t, "continue", TokenKinds[T_CONTINUE].Name, "should be continue")
	assert.Equal(t, "debugger", TokenKinds[T_DEBUGGER].Name, "should be debugger")
	assert.Equal(t, "default", TokenKinds[T_DEFAULT].Name, "should be default")
	assert.Equal(t, "do", TokenKinds[T_DO].Name, "should be do")
	assert.Equal(t, "else", TokenKinds[T_ELSE].Name, "should be else")
	assert.Equal(t, "enum", TokenKinds[T_ENUM].Name, "should be enum")
	assert.Equal(t, "export", TokenKinds[T_EXPORT].Name, "should be export")
	assert.Equal(t, "extends", TokenKinds[T_EXTENDS].Name, "should be extends")
	assert.Equal(t, "finally", TokenKinds[T_FINALLY].Name, "should be finally")
	assert.Equal(t, "for", TokenKinds[T_FOR].Name, "should be for")
	assert.Equal(t, "function", TokenKinds[T_FUNC].Name, "should be function")
	assert.Equal(t, "if", TokenKinds[T_IF].Name, "should be if")
	assert.Equal(t, "import", TokenKinds[T_IMPORT].Name, "should be import")
	assert.Equal(t, "new", TokenKinds[T_NEW].Name, "should be new")
	assert.Equal(t, "return", TokenKinds[T_RETURN].Name, "should be return")
	assert.Equal(t, "super", TokenKinds[T_SUPER].Name, "should be super")
	assert.Equal(t, "switch", TokenKinds[T_SWITCH].Name, "should be switch")
	assert.Equal(t, "this", TokenKinds[T_THIS].Name, "should be this")
	assert.Equal(t, "throw", TokenKinds[T_THROW].Name, "should be throw")
	assert.Equal(t, "try", TokenKinds[T_TRY].Name, "should be try")
	assert.Equal(t, "var", TokenKinds[T_VAR].Name, "should be var")
	assert.Equal(t, "while", TokenKinds[T_WHILE].Name, "should be while")
	assert.Equal(t, "with", TokenKinds[T_WITH].Name, "should be with")

	// contextual keywords
	assert.Equal(t, "let", TokenKinds[T_LET].Name, "should be let")
	assert.Equal(t, "static", TokenKinds[T_STATIC].Name, "should be static")
	assert.Equal(t, "implements", TokenKinds[T_IMPLEMENTS].Name, "should be implements")
	assert.Equal(t, "interface", TokenKinds[T_INTERFACE].Name, "should be interface")
	assert.Equal(t, "package", TokenKinds[T_PACKAGE].Name, "should be package")
	assert.Equal(t, "private", TokenKinds[T_PRIVATE].Name, "should be private")
	assert.Equal(t, "protected", TokenKinds[T_PROTECTED].Name, "should be protected")
	assert.Equal(t, "public", TokenKinds[T_PUBLIC].Name, "should be public")
	assert.Equal(t, "as", TokenKinds[T_AS].Name, "should be as")
	assert.Equal(t, "async", TokenKinds[T_ASYNC].Name, "should be async")
	assert.Equal(t, "from", TokenKinds[T_FROM].Name, "should be from")
	assert.Equal(t, "get", TokenKinds[T_GET].Name, "should be get")
	assert.Equal(t, "meta", TokenKinds[T_META].Name, "should be meta")
	assert.Equal(t, "of", TokenKinds[T_OF].Name, "should be of")
	assert.Equal(t, "set", TokenKinds[T_SET].Name, "should be set")
	assert.Equal(t, "target", TokenKinds[T_TARGET].Name, "should be target")
	assert.Equal(t, "yield", TokenKinds[T_YIELD].Name, "should be yield")

	assert.Equal(t, "regexp", TokenKinds[T_REGEXP].Name, "should be regexp")
	assert.Equal(t, "`", TokenKinds[T_BACK_QUOTE].Name, "should be `")
	assert.Equal(t, "{", TokenKinds[T_BRACE_L].Name, "should be `{")
	assert.Equal(t, "}", TokenKinds[T_BRACE_R].Name, "should be }")
	assert.Equal(t, "(", TokenKinds[T_PAREN_L].Name, "should be (")
	assert.Equal(t, ")", TokenKinds[T_PAREN_R].Name, "should be )")
	assert.Equal(t, "[", TokenKinds[T_BRACKET_L].Name, "should be [")
	assert.Equal(t, "]", TokenKinds[T_BRACKET_R].Name, "should be ]")
	assert.Equal(t, ".", TokenKinds[T_DOT].Name, "should be .")
	assert.Equal(t, "...", TokenKinds[T_DOT_TRI].Name, "should be ...")
	assert.Equal(t, ";", TokenKinds[T_SEMI].Name, "should be ;")
	assert.Equal(t, ",", TokenKinds[T_COMMA].Name, "should be ,")
	assert.Equal(t, "?", TokenKinds[T_HOOK].Name, "should be ?")
	assert.Equal(t, ":", TokenKinds[T_COLON].Name, "should be :")
	assert.Equal(t, "++", TokenKinds[T_INC].Name, "should be ++")
	assert.Equal(t, "--", TokenKinds[T_DEC].Name, "should be --")
	assert.Equal(t, "?.", TokenKinds[T_OPT_CHAIN].Name, "should be ?.")
	assert.Equal(t, "=>", TokenKinds[T_ARROW].Name, "should be =>")

	assert.Equal(t, "??", TokenKinds[T_NULLISH].Name, "should be ??")

	// relational
	assert.Equal(t, "<", TokenKinds[T_LT].Name, "should be <")
	assert.Equal(t, ">", TokenKinds[T_GT].Name, "should be >")
	assert.Equal(t, "<=", TokenKinds[T_LE].Name, "should be <=")
	assert.Equal(t, ">=", TokenKinds[T_GE].Name, "should be >=")

	// equality
	assert.Equal(t, "==", TokenKinds[T_EQ].Name, "should be ==")
	assert.Equal(t, "!=", TokenKinds[T_NE].Name, "should be !=")
	assert.Equal(t, "===", TokenKinds[T_EQ_S].Name, "should be ===")
	assert.Equal(t, "!==", TokenKinds[T_NE_S].Name, "should be !==")

	// bitwise
	assert.Equal(t, "<<", TokenKinds[T_LSH].Name, "should be <<")
	assert.Equal(t, ">>", TokenKinds[T_RSH].Name, "should be >>")
	assert.Equal(t, ">>>", TokenKinds[T_RSH_U].Name, "should be >>>")
	assert.Equal(t, "|", TokenKinds[T_BIT_OR].Name, "should be |")
	assert.Equal(t, "^", TokenKinds[T_BIT_XOR].Name, "should be ^")
	assert.Equal(t, "&", TokenKinds[T_BIT_AND].Name, "should be &")

	assert.Equal(t, "||", TokenKinds[T_OR].Name, "should be ||")
	assert.Equal(t, "&&", TokenKinds[T_AND].Name, "should be &&")

	assert.Equal(t, "instanceof", TokenKinds[T_INSTANCE_OF].Name, "should be instanceof")
	assert.Equal(t, "in", TokenKinds[T_IN].Name, "should be in")

	// unary
	assert.Equal(t, "+", TokenKinds[T_ADD].Name, "should be +")
	assert.Equal(t, "-", TokenKinds[T_SUB].Name, "should be -")
	assert.Equal(t, "*", TokenKinds[T_MUL].Name, "should be *")
	assert.Equal(t, "/", TokenKinds[T_DIV].Name, "should be /")
	assert.Equal(t, "%", TokenKinds[T_MOD].Name, "should be %")
	assert.Equal(t, "**", TokenKinds[T_POW].Name, "should be **")

	// assignment
	assert.Equal(t, "=", TokenKinds[T_ASSIGN].Name, "should be =")
	assert.Equal(t, "+=", TokenKinds[T_ASSIGN_ADD].Name, "should be +=")
	assert.Equal(t, "-=", TokenKinds[T_ASSIGN_SUB].Name, "should be -=")
	assert.Equal(t, "??=", TokenKinds[T_ASSIGN_COALESCE].Name, "should be ??=")
	assert.Equal(t, "||=", TokenKinds[T_ASSIGN_OR].Name, "should be ||=")
	assert.Equal(t, "&&=", TokenKinds[T_ASSIGN_AND].Name, "should be &&=")
	assert.Equal(t, "|=", TokenKinds[T_ASSIGN_BIT_OR].Name, "should be |=")
	assert.Equal(t, "^=", TokenKinds[T_ASSIGN_BIT_XOR].Name, "should be ^=")
	assert.Equal(t, "&=", TokenKinds[T_ASSIGN_BIT_AND].Name, "should be &=")
	assert.Equal(t, "<<=", TokenKinds[T_ASSIGN_BIT_LSH].Name, "should be <<=")
	assert.Equal(t, ">>=", TokenKinds[T_ASSIGN_BIT_RSH].Name, "should be >>=")
	assert.Equal(t, ">>>=", TokenKinds[T_ASSIGN_BIT_RSH_U].Name, "should be >>>=")
	assert.Equal(t, "*=", TokenKinds[T_ASSIGN_MUL].Name, "should be *=")
	assert.Equal(t, "/=", TokenKinds[T_ASSIGN_DIV].Name, "should be /=")
	assert.Equal(t, "%=", TokenKinds[T_ASSIGN_MOD].Name, "should be %=")
	assert.Equal(t, "**=", TokenKinds[T_ASSIGN_POW].Name, "should be **=")
}

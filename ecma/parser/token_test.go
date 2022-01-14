package parser

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/fuzz"
)

func TestLabel(t *testing.T) {
	AssertEqual(t, "EOF", TokenKinds[T_EOF].Name, "should be end of script")
	AssertEqual(t, "comment", TokenKinds[T_COMMENT].Name, "should be comment")

	// literals
	AssertEqual(t, "null", TokenKinds[T_NULL].Name, "should be null")
	AssertEqual(t, "true", TokenKinds[T_TRUE].Name, "should be true")
	AssertEqual(t, "false", TokenKinds[T_FALSE].Name, "should be false")
	AssertEqual(t, "number", TokenKinds[T_NUM].Name, "should be number")
	AssertEqual(t, "string", TokenKinds[T_STRING].Name, "should be string")

	AssertEqual(t, "template span", TokenKinds[T_TPL_SPAN].Name, "should be template span")
	AssertEqual(t, "template tail", TokenKinds[T_TPL_TAIL].Name, "should be template tail")

	AssertEqual(t, "identifier", TokenKinds[T_NAME].Name, "should be identifier")
	AssertEqual(t, "private identifier", TokenKinds[T_NAME_PVT].Name, "should be private identifier")

	// keywords
	AssertEqual(t, "break", TokenKinds[T_BREAK].Name, "should be break")
	AssertEqual(t, "case", TokenKinds[T_CASE].Name, "should be case")
	AssertEqual(t, "catch", TokenKinds[T_CATCH].Name, "should be catch")
	AssertEqual(t, "class", TokenKinds[T_CLASS].Name, "should be class")
	AssertEqual(t, "continue", TokenKinds[T_CONTINUE].Name, "should be continue")
	AssertEqual(t, "debugger", TokenKinds[T_DEBUGGER].Name, "should be debugger")
	AssertEqual(t, "default", TokenKinds[T_DEFAULT].Name, "should be default")
	AssertEqual(t, "do", TokenKinds[T_DO].Name, "should be do")
	AssertEqual(t, "else", TokenKinds[T_ELSE].Name, "should be else")
	AssertEqual(t, "enum", TokenKinds[T_ENUM].Name, "should be enum")
	AssertEqual(t, "export", TokenKinds[T_EXPORT].Name, "should be export")
	AssertEqual(t, "extends", TokenKinds[T_EXTENDS].Name, "should be extends")
	AssertEqual(t, "finally", TokenKinds[T_FINALLY].Name, "should be finally")
	AssertEqual(t, "for", TokenKinds[T_FOR].Name, "should be for")
	AssertEqual(t, "function", TokenKinds[T_FUNC].Name, "should be function")
	AssertEqual(t, "if", TokenKinds[T_IF].Name, "should be if")
	AssertEqual(t, "import", TokenKinds[T_IMPORT].Name, "should be import")
	AssertEqual(t, "new", TokenKinds[T_NEW].Name, "should be new")
	AssertEqual(t, "return", TokenKinds[T_RETURN].Name, "should be return")
	AssertEqual(t, "super", TokenKinds[T_SUPER].Name, "should be super")
	AssertEqual(t, "switch", TokenKinds[T_SWITCH].Name, "should be switch")
	AssertEqual(t, "this", TokenKinds[T_THIS].Name, "should be this")
	AssertEqual(t, "throw", TokenKinds[T_THROW].Name, "should be throw")
	AssertEqual(t, "try", TokenKinds[T_TRY].Name, "should be try")
	AssertEqual(t, "var", TokenKinds[T_VAR].Name, "should be var")
	AssertEqual(t, "while", TokenKinds[T_WHILE].Name, "should be while")
	AssertEqual(t, "with", TokenKinds[T_WITH].Name, "should be with")

	// contextual keywords
	AssertEqual(t, "let", TokenKinds[T_LET].Name, "should be let")
	AssertEqual(t, "const", TokenKinds[T_CONST].Name, "should be const")
	AssertEqual(t, "static", TokenKinds[T_STATIC].Name, "should be static")
	AssertEqual(t, "implements", TokenKinds[T_IMPLEMENTS].Name, "should be implements")
	AssertEqual(t, "interface", TokenKinds[T_INTERFACE].Name, "should be interface")
	AssertEqual(t, "package", TokenKinds[T_PACKAGE].Name, "should be package")
	AssertEqual(t, "private", TokenKinds[T_PRIVATE].Name, "should be private")
	AssertEqual(t, "protected", TokenKinds[T_PROTECTED].Name, "should be protected")
	AssertEqual(t, "public", TokenKinds[T_PUBLIC].Name, "should be public")
	AssertEqual(t, "as", TokenKinds[T_AS].Name, "should be as")
	AssertEqual(t, "async", TokenKinds[T_ASYNC].Name, "should be async")
	AssertEqual(t, "from", TokenKinds[T_FROM].Name, "should be from")
	AssertEqual(t, "get", TokenKinds[T_GET].Name, "should be get")
	AssertEqual(t, "meta", TokenKinds[T_META].Name, "should be meta")
	AssertEqual(t, "of", TokenKinds[T_OF].Name, "should be of")
	AssertEqual(t, "set", TokenKinds[T_SET].Name, "should be set")
	AssertEqual(t, "target", TokenKinds[T_TARGET].Name, "should be target")
	AssertEqual(t, "yield", TokenKinds[T_YIELD].Name, "should be yield")

	AssertEqual(t, "regexp", TokenKinds[T_REGEXP].Name, "should be regexp")
	AssertEqual(t, "`", TokenKinds[T_BACK_QUOTE].Name, "should be `")
	AssertEqual(t, "{", TokenKinds[T_BRACE_L].Name, "should be `{")
	AssertEqual(t, "}", TokenKinds[T_BRACE_R].Name, "should be }")
	AssertEqual(t, "(", TokenKinds[T_PAREN_L].Name, "should be (")
	AssertEqual(t, ")", TokenKinds[T_PAREN_R].Name, "should be )")
	AssertEqual(t, "[", TokenKinds[T_BRACKET_L].Name, "should be [")
	AssertEqual(t, "]", TokenKinds[T_BRACKET_R].Name, "should be ]")
	AssertEqual(t, ".", TokenKinds[T_DOT].Name, "should be .")
	AssertEqual(t, "...", TokenKinds[T_DOT_TRI].Name, "should be ...")
	AssertEqual(t, ";", TokenKinds[T_SEMI].Name, "should be ;")
	AssertEqual(t, ",", TokenKinds[T_COMMA].Name, "should be ,")
	AssertEqual(t, "?", TokenKinds[T_HOOK].Name, "should be ?")
	AssertEqual(t, ":", TokenKinds[T_COLON].Name, "should be :")
	AssertEqual(t, "++", TokenKinds[T_INC].Name, "should be ++")
	AssertEqual(t, "--", TokenKinds[T_DEC].Name, "should be --")
	AssertEqual(t, "?.", TokenKinds[T_OPT_CHAIN].Name, "should be ?.")
	AssertEqual(t, "=>", TokenKinds[T_ARROW].Name, "should be =>")

	AssertEqual(t, "??", TokenKinds[T_NULLISH].Name, "should be ??")

	// relational
	AssertEqual(t, "<", TokenKinds[T_LT].Name, "should be <")
	AssertEqual(t, ">", TokenKinds[T_GT].Name, "should be >")
	AssertEqual(t, "<=", TokenKinds[T_LTE].Name, "should be <=")
	AssertEqual(t, ">=", TokenKinds[T_GTE].Name, "should be >=")

	// equality
	AssertEqual(t, "==", TokenKinds[T_EQ].Name, "should be ==")
	AssertEqual(t, "!=", TokenKinds[T_NE].Name, "should be !=")
	AssertEqual(t, "===", TokenKinds[T_EQ_S].Name, "should be ===")
	AssertEqual(t, "!==", TokenKinds[T_NE_S].Name, "should be !==")

	// bitwise
	AssertEqual(t, "<<", TokenKinds[T_LSH].Name, "should be <<")
	AssertEqual(t, ">>", TokenKinds[T_RSH].Name, "should be >>")
	AssertEqual(t, ">>>", TokenKinds[T_RSH_U].Name, "should be >>>")
	AssertEqual(t, "|", TokenKinds[T_BIT_OR].Name, "should be |")
	AssertEqual(t, "^", TokenKinds[T_BIT_XOR].Name, "should be ^")
	AssertEqual(t, "&", TokenKinds[T_BIT_AND].Name, "should be &")

	AssertEqual(t, "||", TokenKinds[T_OR].Name, "should be ||")
	AssertEqual(t, "&&", TokenKinds[T_AND].Name, "should be &&")

	AssertEqual(t, "instanceof", TokenKinds[T_INSTANCE_OF].Name, "should be instanceof")
	AssertEqual(t, "in", TokenKinds[T_IN].Name, "should be in")

	// unary
	AssertEqual(t, "+", TokenKinds[T_ADD].Name, "should be +")
	AssertEqual(t, "-", TokenKinds[T_SUB].Name, "should be -")
	AssertEqual(t, "*", TokenKinds[T_MUL].Name, "should be *")
	AssertEqual(t, "/", TokenKinds[T_DIV].Name, "should be /")
	AssertEqual(t, "%", TokenKinds[T_MOD].Name, "should be %")
	AssertEqual(t, "**", TokenKinds[T_POW].Name, "should be **")

	// assignment
	AssertEqual(t, "=", TokenKinds[T_ASSIGN].Name, "should be =")
	AssertEqual(t, "+=", TokenKinds[T_ASSIGN_ADD].Name, "should be +=")
	AssertEqual(t, "-=", TokenKinds[T_ASSIGN_SUB].Name, "should be -=")
	AssertEqual(t, "??=", TokenKinds[T_ASSIGN_NULLISH].Name, "should be ??=")
	AssertEqual(t, "||=", TokenKinds[T_ASSIGN_OR].Name, "should be ||=")
	AssertEqual(t, "&&=", TokenKinds[T_ASSIGN_AND].Name, "should be &&=")
	AssertEqual(t, "|=", TokenKinds[T_ASSIGN_BIT_OR].Name, "should be |=")
	AssertEqual(t, "^=", TokenKinds[T_ASSIGN_BIT_XOR].Name, "should be ^=")
	AssertEqual(t, "&=", TokenKinds[T_ASSIGN_BIT_AND].Name, "should be &=")
	AssertEqual(t, "<<=", TokenKinds[T_ASSIGN_BIT_LSH].Name, "should be <<=")
	AssertEqual(t, ">>=", TokenKinds[T_ASSIGN_BIT_RSH].Name, "should be >>=")
	AssertEqual(t, ">>>=", TokenKinds[T_ASSIGN_BIT_RSH_U].Name, "should be >>>=")
	AssertEqual(t, "*=", TokenKinds[T_ASSIGN_MUL].Name, "should be *=")
	AssertEqual(t, "/=", TokenKinds[T_ASSIGN_DIV].Name, "should be /=")
	AssertEqual(t, "%=", TokenKinds[T_ASSIGN_MOD].Name, "should be %=")
	AssertEqual(t, "**=", TokenKinds[T_ASSIGN_POW].Name, "should be **=")
}

package parser

import "fmt"

type SourceError struct {
	msg  string
	file string
	line int
	col  int
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("%s at %sL%d:%d\n", e.msg, e.file, e.line, e.col)
}

func NewSourceError(msg, file string, line, col int) *SourceError {
	return &SourceError{
		file: file,
		line: line,
		col:  col,
	}
}

type LexerError struct {
	msg  string
	file string
	line int
	col  int
}

func NewLexerError(msg, file string, line, col int) *LexerError {
	return &LexerError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}

type ParserError struct {
	msg  string
	file string
	line int
	col  int
}

func NewParserError(msg, file string, line, col int) *ParserError {
	return &ParserError{
		msg:  msg,
		file: file,
		line: line,
		col:  col,
	}
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s at %s(%d:%d)", e.msg, e.file, e.line, e.col)
}

const (
	ERR_UNEXPECTED_TOKEN                      = "Unexpected token"
	ERR_UNEXPECTED_TOKEN_TYPE                 = "Unexpected token `%s`"
	ERR_UNTERMINATED_COMMENT                  = "Unterminated comment"
	ERR_UNTERMINATED_REGEXP                   = "Unterminated regular expression"
	ERR_UNTERMINATED_STR                      = "Unterminated string constant"
	ERR_IDENT_AFTER_NUMBER                    = "Identifier directly after number"
	ERR_INVALID_NUMBER                        = "Invalid number"
	ERR_LEGACY_OCTAL_IN_STRICT_MODE           = "Octal literals are not allowed in strict mode"
	ERR_LEGACY_OCTAL_ESCAPE_IN_TPL            = "Octal escape sequences are not allowed in template strings"
	ERR_LEGACY_OCTAL_ESCAPE_IN_STRICT_MODE    = "Octal escape sequences are not allowed in strict mode"
	ERR_EXPECTING_UNICODE_ESCAPE              = "Expecting Unicode escape sequence \\uXXXX"
	ERR_BAD_ESCAPE_SEQ                        = "Bad character escape sequence"
	ERR_BAD_RUNE                              = "Bad character"
	ERR_UNTERMINATED_TPL                      = "Unterminated template"
	ERR_INVALID_UNICODE_ESCAPE                = "Invalid Unicode escape"
	ERR_ILLEGAL_RETURN                        = "Illegal return"
	ERR_ILLEGAL_BREAK                         = "Illegal break"
	ERR_DUP_LABEL                             = "Label `%s` already declared"
	ERR_UNDEF_LABEL                           = "Undefined label `%s`"
	ERR_ILLEGAL_CONTINUE                      = "Illegal continue"
	ERR_MULTI_DEFAULT                         = "Multiple default clauses"
	ERR_ILLEGAL_LEXICAL_DEC                   = "Illegal lexical declaration"
	ERR_ASSIGN_TO_RVALUE                      = "Assigning to rvalue"
	ERR_DUP_BINDING                           = "Must have a single binding"
	ERR_RESERVED_WORD_IN_STRICT_MODE          = "Unexpected strict mode reserved word"
	ERR_STRICT_DIRECTIVE_AFTER_NOT_SIMPLE     = "Illegal 'use strict' directive in function with non-simple parameter list"
	ERR_DUP_PARAM_NAME                        = "Parameter name clash"
	ERR_REST_TRAILING_COMMA                   = "Unexpected trailing comma after rest element"
	ERR_TRAILING_COMMA                        = "Unexpected trailing comma"
	ERR_REST_ELEM_MUST_LAST                   = "Rest element must be last element"
	ERR_DELETE_LOCAL_IN_STRICT                = "Deleting local variable in strict mode"
	ERR_REDEF_PROP                            = "Redefinition of property"
	ERR_ILLEGAL_NEWLINE_AFTER_THROW           = "Illegal newline after throw"
	ERR_CONST_DEC_INIT_REQUIRED               = "Const declarations require an initialization value"
	ERR_GETTER_SHOULD_NO_PARAM                = "Getter must not have any formal parameters"
	ERR_SETTER_SHOULD_ONE_PARAM               = "Setter must have exactly one formal parameter"
	ERR_ESCAPE_IN_KEYWORD                     = "Keyword must not contain escaped characters"
	ERR_ID_DUP_DEF                            = "Identifier `%s` has already been declared"
	ERR_WITH_STMT_IN_STRICT                   = "Strict mode code may not include a with statement"
	ERR_CLASS_NAME_REQUIRED                   = "Class name is required"
	ERR_SHORTHAND_PROP_ASSIGN_NOT_IN_DESTRUCT = "Shorthand property assignments are valid only in destructuring patterns"
	ERR_REST_ARG_NOT_SIMPLE                   = "Invalid rest operator's argument"
	ERR_INVALID_PAREN_ASSIGN_PATTERN          = "Invalid parenthesized assignment pattern"
	ERR_OBJ_PATTERN_CANNOT_FN                 = "Object pattern can't contain getter or setter"
	ERR_REST_CANNOT_SET_DEFAULT               = "Rest elements cannot have a default value"
	ERR_MALFORMED_ARROW_PARAM                 = "Malformed arrow function parameter list"
)
